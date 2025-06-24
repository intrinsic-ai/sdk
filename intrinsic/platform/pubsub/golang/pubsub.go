// Copyright 2023 Intrinsic Innovation LLC

// Package pubsub is a wrapper around the imw library.
//
// This package exposes the same pubsub interface as defined in
// intrinsic/platform/pubsub/pubsub.h
//
// To avoid circular dependencies, we split the core functionality of the Go
// bindings into two distinct packages:
//
//	pubsub:
//	  The concrete implementation of the bindings that call into the
//	  the imw library.
//
//	pubsubinterface:
//	  The high-level interface to pubsub exposed via the concrete
//	  implementation in pubsub.
//
// When using this suite of packages, you will likely need to use this package
// (pubsub) to instantiate a pubsub instance, but pubsubinterface should
// be used for type-level constructs.
package pubsub

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/cgo"
	"strings"
	"sync"
	"time"
	"unsafe"

	"flag"
	log "github.com/golang/glog"
	"google.golang.org/protobuf/proto"
	anypb "google.golang.org/protobuf/types/known/anypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	pubsubpb "intrinsic/platform/pubsub/adapters/pubsub_go_proto"
	"intrinsic/platform/pubsub/golang/kvstore"
	"intrinsic/platform/pubsub/golang/pubsubinterface"
	"intrinsic/util/path_resolver/pathresolver"
)

/*
#include <stdlib.h>  // for C.free
#include "intrinsic/platform/pubsub/golang/pubsub_c.h"

// We forward declare functions that will be implemented below so we can take their address from Go code.
void intrinsic_ImwSubscriptionCallback(void*, void*, size_t, void*);
void intrinsic_ImwQueryStaticCallback(void*, void*, size_t, void*);
void intrinsic_ImwQueryDoneStaticCallback(void*, void*);
*/
import "C"

var zenohRouter = flag.String("zenoh_router", "", "Override the default Zenoh connection to PROTOCOL/HOSTNAME:PORT")

const highConsistencyTimeout = 30 * time.Second
const highConsistencyGetTimeout = 10 * time.Second
const defaultKeyPrefix = "kv_store"
const replicationKeyPrefix = "kv_store_repl"

// NewPubSub creates a new PubSub adapter if possible. Returns either a valid handle
// or an error, but not both. The caller is responsible for freeing up resources
// after use by calling Close() on the returned handle.
func NewPubSub() (*Handle, error) {
	zh, err := getZenohHandle()
	if err != nil {
		return nil, err
	}
	result := &Handle{
		zenohHandle: zh,
	}

	zenohConfig, err := getZenohPeerConfig()
	if err != nil {
		return nil, err
	}

	if err := result.zenohHandle.ImwInit(zenohConfig); err != nil {
		return nil, err
	}

	return result, nil
}

func addTopicPrefix(topic string) string {
	if topic[0] == '/' {
		return "in" + topic
	}
	return "in/" + topic
}

func errorFromImwRet(imwRet C.int) error {
	switch imwRet {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("imw returned an error")
	case 2:
		return fmt.Errorf("imw is not initialized")
	default:
		return fmt.Errorf("unknown error from imw")
	}
}

// Handle represents an instance of a Fast DDS pubsub adapter, as exposed via
// the interface defined in pubsubinterface.PubSub.
type Handle struct {
	mutex sync.Mutex

	zenohHandle zenohHandle
}

const zenohConfigPath = "intrinsic/platform/pubsub/zenoh_util/peer_config.json"

func isRunningUnderTest() bool {
	if os.Getenv("TEST_TMPDIR") == "" && os.Getenv("PORTSERVER_ADDRESS") == "" {
		return true
	}

	return false
}

func isRunningInKubernetes() bool {
	if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
		return false
	}

	return true
}

func getZenohPeerConfig() (string, error) {
	var path string
	var err error

	// If we're running in k8s, then we can assume the config file is
	// available via base layer, rather than needing to find it in
	// runfiles.
	if isRunningInKubernetes() {
		path = zenohConfigPath
	} else {
		path, err = pathresolver.ResolveRunfilesPath(zenohConfigPath)
		if err != nil {
			return "", err
		}
	}

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	config := string(b)

	if len(config) != 0 {
		// This logic all seems pretty fragile, and would benefit from some more structured handling.
		if isRunningUnderTest() {
			// Remove listen endpoints when running in test. (go/forge-limits#ipv4)
			config = strings.ReplaceAll(config, `"tcp/0.0.0.0:0"`, "")
		} else if isRunningInKubernetes() {
			if allowedIP := os.Getenv("ALLOWED_PUBSUB_IPv4"); allowedIP != "" {
				config = strings.ReplaceAll(config, "0.0.0.0", allowedIP)
			}
		}
	}

	if *zenohRouter != "" {
		// replace router endpoint
		config = strings.Replace(config, "tcp/zenoh-router.app-intrinsic-base.svc.cluster.local:7447", *zenohRouter, 1)
	}

	return config, nil
}

func topicConfigToZenohQos(config pubsubinterface.TopicConfig) (string, error) {
	switch config.Qos {
	case pubsubinterface.Sensor:
		return "Sensor", nil
	case pubsubinterface.HighReliability:
		return "HighReliability", nil
	default:
		return "UNKNOWN", fmt.Errorf("unknown QOS setting %v", config.Qos)
	}
}

// NewPublisher creates a publisher for a topic
func (ps *Handle) NewPublisher(topic string, config pubsubinterface.TopicConfig) (pubsubinterface.Publisher, error) {
	topicQos, err := topicConfigToZenohQos(config)
	if err != nil {
		return nil, err
	}

	publisher := &publisherHandle{topicName: topic, zenohHandle: ps.zenohHandle}
	if err := ps.zenohHandle.ImwCreatePublisher(addTopicPrefix(topic), topicQos); err != nil {
		return nil, err
	}
	return publisher, nil
}

// NewSubscription will create a subscription to the given topic, using the exemplar proto as the
// type expected to be called by the msg_callback.
func (ps *Handle) NewSubscription(topic string, config pubsubinterface.TopicConfig, exemplar proto.Message, msgCallback func(proto.Message), errCallback func(string, error)) (pubsubinterface.Subscription, error) {
	topicQos, err := topicConfigToZenohQos(config)
	if err != nil {
		return nil, err
	}

	subscription := &subscriptionHandle{
		topicName:   topic,
		zenohHandle: ps.zenohHandle,
		exemplar:    exemplar,
		callback: func(sub *subscriptionHandle, topic string, bytes []byte) {
			packet := &pubsubpb.PubSubPacket{}
			if err := proto.Unmarshal(bytes, packet); err != nil {
				log.Errorf("Failed to unmarshal packet: %v", err)
				return
			}

			msg := sub.exemplar.ProtoReflect().New().Interface()
			if err := packet.GetPayload().UnmarshalTo(msg); err != nil {
				errCallback(string(bytes), err)
				return
			}

			msgCallback(msg)
		},
	}
	subh := cgo.NewHandle(subscription)
	subscription.subHandle = subh

	if err := ps.zenohHandle.ImwCreateSubscription(addTopicPrefix(topic), subscription, topicQos); err != nil {
		return nil, err
	}

	return subscription, nil
}

// NewRawSubscription will create a raw subscription to the given topic, passing the full packet to callback.
func (ps *Handle) NewRawSubscription(topic string, config pubsubinterface.TopicConfig, callback func(*pubsubpb.PubSubPacket)) (pubsubinterface.Subscription, error) {
	topicQos, err := topicConfigToZenohQos(config)
	if err != nil {
		return nil, err
	}

	subscription := &subscriptionHandle{
		topicName:   topic,
		zenohHandle: ps.zenohHandle,
		callback: func(sub *subscriptionHandle, topic string, bytes []byte) {
			packet := &pubsubpb.PubSubPacket{}
			if err := proto.Unmarshal(bytes, packet); err != nil {
				log.Errorf("Failed to unmarshal packet: %v", err)
				return
			}

			callback(packet)
		},
	}
	subh := cgo.NewHandle(subscription)
	subscription.subHandle = subh

	if err := ps.zenohHandle.ImwCreateSubscription(addTopicPrefix(topic), subscription, topicQos); err != nil {
		return nil, err
	}

	return subscription, nil
}

// KVStore returns an interface to the KVStore associated with this pubsub
// instance.
func (ps *Handle) KVStore() kvstore.KVStore {
	return &kvStoreHandle{ps: ps, zenohHandle: ps.zenohHandle}
}

// KVStoreReplicated returns an interface to the replicated KVStore associated with this pubsub
// instance.
func (ps *Handle) KVStoreReplicated() kvstore.KVStore {
	return &kvStoreHandle{ps: ps, zenohHandle: ps.zenohHandle, keyPrefix: replicationKeyPrefix}
}

// KVStoreWithPrefix returns an interface to the KVStore associated with this pubsub
// but default key prefixes are overridden with the given prefix.
func (ps *Handle) KVStoreWithPrefix(prefix string) kvstore.KVStore {
	return &kvStoreHandle{ps: ps, zenohHandle: ps.zenohHandle, keyPrefix: prefix}
}

// Close the PubSub connection, unsubscribe from all topics, and free the
// associated resources.
func (ps *Handle) Close() {
	ps.zenohHandle.Destroy()
}

type subscriptionHandle struct {
	topicName   string
	zenohHandle zenohHandle
	callback    func(sub *subscriptionHandle, topic string, bytes []byte)

	callbackPtr unsafe.Pointer
	exemplar    proto.Message

	subHandle cgo.Handle
}

func (s *subscriptionHandle) TopicName() string { return s.topicName }

func (s *subscriptionHandle) Close() {
	inKeyExprString := C.CString(addTopicPrefix(s.topicName))
	defer C.free(unsafe.Pointer(inKeyExprString))

	if err := s.zenohHandle.ImwDestroySubscription(addTopicPrefix(s.topicName), s); err != nil {
		panic(err)
	}
	s.subHandle.Delete()
}

type publisherHandle struct {
	topicName   string
	zenohHandle zenohHandle
}

func (p *publisherHandle) TopicName() string { return p.topicName }

func (p *publisherHandle) Publish(msg proto.Message) error {
	packet := &pubsubpb.PubSubPacket{
		PublishTime: timestamppb.New(time.Now()),
		Payload:     &anypb.Any{},
	}

	if err := packet.GetPayload().MarshalFrom(msg); err != nil {
		return err
	}

	bytes, err := proto.Marshal(packet)
	if err != nil {
		return err
	}
	return p.zenohHandle.ImwPublish(addTopicPrefix(p.topicName), bytes)
}

func (p *publisherHandle) HasMatchingSubscribers() (bool, error) {
	return p.zenohHandle.ImwPublisherHasMatchingSubscribers(addTopicPrefix(p.topicName))
}

func (p *publisherHandle) Close() {
	if err := p.zenohHandle.ImwDestroyPublisher(addTopicPrefix(p.topicName)); err != nil {
		panic(err)
	}
}

type kvStoreHandle struct {
	ps          *Handle
	zenohHandle zenohHandle
	keyPrefix   string
}

func (kv *kvStoreHandle) addKeyPrefix(key string) string {
	keyPrefix := defaultKeyPrefix
	if kv.keyPrefix != "" {
		keyPrefix = kv.keyPrefix
	}
	if key[0] == '/' {
		return keyPrefix + key
	}
	return keyPrefix + "/" + key
}

func (kv *kvStoreHandle) Set(key string, value proto.Message, highConsistency bool) error {
	valueAny, err := anypb.New(value)
	if err != nil {
		return err
	}

	valueBytes, err := proto.Marshal(valueAny)
	if err != nil {
		return err
	}

	if err := kv.zenohHandle.ImwSet(kv.addKeyPrefix(key), valueBytes); err != nil {
		return err
	}

	if highConsistency {
		ctx, cancel := context.WithTimeout(context.Background(), highConsistencyTimeout)
		defer cancel()

	loopUntilWritten:
		for {
			select {
			case <-ctx.Done():
				// timeout
				return fmt.Errorf("timeout waiting for high consistency: %w", kvstore.ErrDeadlineExceeded)
			default:
				timeout := highConsistencyGetTimeout
				_, err := kv.Get(key, &timeout)
				if err != nil {
					// Small wait before retrying.
					time.Sleep(100 * time.Millisecond)
					continue
				}
				break loopUntilWritten
			}
		}
	}
	return nil
}

type queryHandle struct {
	query func(keyexpr string, bytes []byte)
	done  func(keyexpr string)

	handle cgo.Handle
}

func (q *queryHandle) Close() {
	q.handle.Delete()
}

func (kv *kvStoreHandle) GetAll(key string, valueCallback func(*anypb.Any), ondoneCallback func(string)) (kvstore.KVQuery, error) {
	var qh *queryHandle
	qh = &queryHandle{
		query: func(keyexpr string, bytes []byte) {
			value := &anypb.Any{}
			if err := proto.Unmarshal(bytes, value); err != nil {
				panic(err)
			}
			valueCallback(value)
		},
		done: func(keyexpr string) {
			ondoneCallback(keyexpr)
		},
	}
	qh.handle = cgo.NewHandle(qh)

	if err := kv.zenohHandle.ImwQuery(kv.addKeyPrefix(key), qh); err != nil {
		return nil, err
	}

	return qh, nil
}

//export intrinsic_ImwQueryStaticCallback
func intrinsic_ImwQueryStaticCallback(keyexpr unsafe.Pointer, responseBytes unsafe.Pointer, responseBytesLen C.size_t, userContext unsafe.Pointer) {
	if userContext == nil {
		return
	}
	h := *(*cgo.Handle)(userContext)
	callbacks := h.Value().(*queryHandle)
	callbacks.query(C.GoString((*C.char)(keyexpr)), C.GoBytes(responseBytes, C.int(responseBytesLen)))
}

//export intrinsic_ImwQueryDoneStaticCallback
func intrinsic_ImwQueryDoneStaticCallback(keyexpr unsafe.Pointer, userContext unsafe.Pointer) {
	if userContext == nil {
		return
	}
	h := *(*cgo.Handle)(userContext)
	callbacks := h.Value().(*queryHandle)
	callbacks.done(C.GoString((*C.char)(keyexpr)))
}

func (kv *kvStoreHandle) Get(key string, timeout *time.Duration) (*anypb.Any, error) {

	resultCh := make(chan *anypb.Any)

	valueCallback := func(value *anypb.Any) {
		resultCh <- value
	}

	done := make(chan any)
	onDoneCallback := func(string) {
		done <- true
	}

	errCh := make(chan error)

	go func() {
		query, err := kv.GetAll(key, valueCallback, onDoneCallback)
		if err != nil {
			errCh <- err
		}
		_ = <-done
		query.Close()
	}()

	timeoutCh := make(chan time.Time)

	if timeout != nil {
		go func() {
			timeoutCh <- <-time.After(*timeout)
		}()
	}

	select {
	case result := <-resultCh:
		return result, nil
	case _ = <-done:
		return nil, fmt.Errorf("%q not found: %w", key, kvstore.ErrNotFound)
	case err := <-errCh:
		return nil, err
	case _ = <-timeoutCh:
		return nil, fmt.Errorf("timeout waiting for %q: %w", key, kvstore.ErrDeadlineExceeded)
	}
}
func (kv *kvStoreHandle) Delete(key string) error {
	prefixedKey := kv.addKeyPrefix(key)

	cPrefixedKey := C.CString(prefixedKey)
	defer C.free(unsafe.Pointer(cPrefixedKey))

	if res := C.ZenohHandleImwDeleteKeyExpr(kv.zenohHandle.Ptr(), cPrefixedKey); res != 0 {
		return errorFromImwRet(res)
	}
	return nil
}

type zenohHandle interface {
	Destroy()
	ImwInit(config string) error
	ImwFini() error
	ImwCreatePublisher(keyExpr string, qos string) error
	ImwDestroyPublisher(keyExpr string) error
	ImwPublish(keyExpr string, bytes []byte) error
	ImwPublisherHasMatchingSubscribers(keyExpr string) (bool, error)
	ImwCreateSubscription(keyExpr string, sub *subscriptionHandle, qos string) error
	ImwDestroySubscription(keyExpr string, sub *subscriptionHandle) error
	ImwSet(keyExpr string, value []byte) error
	ImwQuery(keyExpr string, query *queryHandle) error
	Ptr() unsafe.Pointer
}

type zenohHandleImpl struct {
	ptr unsafe.Pointer
}

var globalZenohHandle *zenohHandleImpl
var zenohHandleRefCount int64 = 0
var zenohHandleMutex sync.Mutex

func getZenohHandle() (zenohHandle, error) {
	zenohHandleMutex.Lock()
	defer zenohHandleMutex.Unlock()

	if zenohHandleRefCount == 0 {
		ptr := C.NewZenohHandle()
		if ptr == nil {
			return nil, fmt.Errorf("something went wrong")
		}

		globalZenohHandle = &zenohHandleImpl{
			ptr: ptr,
		}

		// Unconditionally close the zenoh handle when the pointer to
		// it is garbage collected. This case can occur if the refcount
		// never goes to zero before a program terminates.
		runtime.AddCleanup(globalZenohHandle, func(ptr unsafe.Pointer) {
			_ = C.ZenohHandleImwFini(ptr)
		}, globalZenohHandle.ptr)
	} else if globalZenohHandle == nil {
		panic(fmt.Errorf("reference count is nonzero, but globalZenohHandle is nil"))
	}

	return globalZenohHandle, nil
}

func (z *zenohHandleImpl) Destroy() {
	zenohHandleMutex.Lock()
	defer zenohHandleMutex.Unlock()

	zenohHandleRefCount--
	if zenohHandleRefCount == 0 {
		z.ImwFini()
		C.DestroyZenohHandle(z.ptr)
		globalZenohHandle = nil
	}
}

// String type is no bueno here, pass a struct
func (z *zenohHandleImpl) ImwInit(config string) error {
	configString := C.CString(config)
	defer C.free(unsafe.Pointer(configString))
	if res := C.ZenohHandleImwInit(z.ptr, configString); res != 0 {
		return errorFromImwRet(res)
	}

	return nil
}

func (z *zenohHandleImpl) ImwFini() error {
	if res := C.ZenohHandleImwFini(z.ptr); res != 0 {
		return errorFromImwRet(res)
	}

	return nil
}

func (z *zenohHandleImpl) ImwCreatePublisher(keyExpr string, qos string) error {
	keyExprString := C.CString(keyExpr)
	defer C.free(unsafe.Pointer(keyExprString))
	qosString := C.CString(qos)
	defer C.free(unsafe.Pointer(qosString))

	if res := C.ZenohHandleImwCreatePublisher(z.ptr, keyExprString, qosString); res != 0 {
		return errorFromImwRet(res)
	}
	return nil
}

func (z *zenohHandleImpl) ImwDestroyPublisher(keyExpr string) error {
	keyExprString := C.CString(keyExpr)
	defer C.free(unsafe.Pointer(keyExprString))

	if res := C.ZenohHandleImwDestroyPublisher(z.ptr, keyExprString); res != 0 {
		return errorFromImwRet(res)
	}
	return nil
}

func (z *zenohHandleImpl) ImwPublish(keyExpr string, bytes []byte) error {
	keyExprString := C.CString(keyExpr)
	defer C.free(unsafe.Pointer(keyExprString))

	// Avoid the malloc = copy here!
	bytesArray := C.CBytes(bytes)
	defer C.free(bytesArray)

	if res := C.ZenohHandleImwPublish(z.ptr, keyExprString, bytesArray, C.size_t(len(bytes))); res != 0 {
		return errorFromImwRet(res)
	}
	return nil
}

func (z *zenohHandleImpl) ImwPublisherHasMatchingSubscribers(keyExpr string) (bool, error) {
	keyExprString := C.CString(keyExpr)
	defer C.free(unsafe.Pointer(keyExprString))

	cbool := C.bool(false)
	if res := C.ZenohHandleImwPublisherHasMatchingSubscribers(z.ptr, keyExprString, &cbool); res != 0 {
		return false, errorFromImwRet(res)
	}

	return bool(cbool), nil
}

//export intrinsic_ImwSubscriptionCallback
func intrinsic_ImwSubscriptionCallback(keyexpr unsafe.Pointer, bytes unsafe.Pointer, bytesLen C.size_t, userContext unsafe.Pointer) {
	if userContext == nil {
		return
	}
	h := *(*cgo.Handle)(userContext)
	sub := h.Value().(*subscriptionHandle)
	sub.callback(sub, C.GoString((*C.char)(keyexpr)), C.GoBytes(bytes, C.int(bytesLen)))
}

func (z *zenohHandleImpl) ImwCreateSubscription(keyExpr string, sub *subscriptionHandle, qos string) error {
	inKeyExprString := C.CString(keyExpr)
	defer C.free(unsafe.Pointer(inKeyExprString))

	qosString := C.CString(qos)
	defer C.free(unsafe.Pointer(qosString))

	if res := C.ZenohHandleImwCreateSubscription(z.ptr, inKeyExprString, C.zenoh_handle_imw_subscription_callback_fn(C.intrinsic_ImwSubscriptionCallback), qosString, unsafe.Pointer(&sub.subHandle)); res != 0 {
		return errorFromImwRet(res)
	}

	return nil
}

func (z *zenohHandleImpl) ImwDestroySubscription(keyExpr string, sub *subscriptionHandle) error {
	inKeyExprString := C.CString(keyExpr)
	defer C.free(unsafe.Pointer(inKeyExprString))

	if res := C.ZenohHandleImwDestroySubscription(z.ptr, inKeyExprString, C.zenoh_handle_imw_subscription_callback_fn(C.intrinsic_ImwSubscriptionCallback), unsafe.Pointer(&sub.subHandle)); res != 0 {
		return errorFromImwRet(res)
	}

	return nil
}

func (z *zenohHandleImpl) ImwSet(keyExpr string, value []byte) error {
	keyExprString := C.CString(keyExpr)
	defer C.free(unsafe.Pointer(keyExprString))

	bytesArray := C.CBytes(value)
	defer C.free(bytesArray)

	if res := C.ZenohHandleImwSet(z.ptr, keyExprString, bytesArray, C.size_t(len(value))); res != 0 {
		return errorFromImwRet(res)
	}
	return nil
}

func (z *zenohHandleImpl) ImwQuery(keyExpr string, query *queryHandle) error {
	keyExprString := C.CString(keyExpr)
	defer C.free(unsafe.Pointer(keyExprString))

	if res := C.ZenohHandleImwQuery(z.ptr, keyExprString, C.zenoh_handle_imw_query_callback_fn(C.intrinsic_ImwQueryStaticCallback), C.zenoh_handle_imw_query_on_done_fn(C.intrinsic_ImwQueryDoneStaticCallback), nil, 0, unsafe.Pointer(&query.handle), C.uint64_t(0), false); res != 0 {
		return errorFromImwRet(res)
	}

	return nil
}

func (z *zenohHandleImpl) Ptr() unsafe.Pointer {
	return z.ptr
}
