// Copyright 2023 Intrinsic Innovation LLC

// Package pubsubinterface provides type level info for the pubsub package.
//
// This package provides the types used by the PubSub interface. Please see
// intrinsic/platform/pubsub/golang/fast_dds.go for more details.
package pubsubinterface

import (
	"google.golang.org/protobuf/proto"
	anypb "google.golang.org/protobuf/types/known/anypb"
	pubsubpb "intrinsic/platform/pubsub/adapters/pubsub_go_proto"
)

// SubscriptionCallbacks is a struct that contains all of the information needed
// to run user-provided callbacks from the PubSub interface.
type SubscriptionCallbacks struct {
	msgCb func(proto.Message)
	errCb func(string, error)
	ex    proto.Message // exemplar for marshalling protos to C++
}

// NewCallbacks builds a SubscriptionCallbacks struct from the provided callbacks and exemplar proto.
func NewCallbacks(msgCb func(proto.Message), errCb func(string, error), exemplar proto.Message) *SubscriptionCallbacks {
	return &SubscriptionCallbacks{
		msgCb: msgCb,
		errCb: errCb,
		ex:    exemplar,
	}
}

// RunMessageCallback runs the message callback with the given message m
func (cbs *SubscriptionCallbacks) RunMessageCallback(m proto.Message) {
	cbs.msgCb(m)
}

// GetDefaultMessageProto duplicates the exemplar proto used to create the
// SubscriptionCallbacks and returns it with all fields containing their
// default values.
func (cbs *SubscriptionCallbacks) GetDefaultMessageProto() proto.Message {
	return cbs.ex.ProtoReflect().New().Interface()
}

// RunErrorCallback runs the error callback with the given packet string and
// error provided.
func (cbs *SubscriptionCallbacks) RunErrorCallback(packet string, e error) {
	cbs.errCb(packet, e)
}

// RawSubscriptionCallbacks is a struct that contains all of the information needed
// to run user-provided callbacks from the PubSub interface for a raw subscription.
type RawSubscriptionCallbacks struct {
	rawCb func(*pubsubpb.PubSubPacket)
}

// NewRawCallbacks builds a RawSubscriptionCallbacks struct from the provided raw callback.
func NewRawCallbacks(cb func(*pubsubpb.PubSubPacket)) *RawSubscriptionCallbacks {
	return &RawSubscriptionCallbacks{
		rawCb: cb,
	}
}

// RunRawCallback runs the message callback with the given message m
func (cbs *RawSubscriptionCallbacks) RunRawCallback(m *pubsubpb.PubSubPacket) {
	cbs.rawCb(m)
}

// KVQueryCallbacks is a struct that contains all of the information needed
// to run user-provided callbacks from the PubSub interface for a KV query.
type KVQueryCallbacks struct {
	valueCb  func(*anypb.Any)
	onDoneCb func(string)
}

// NewKVQueryCallbacks builds a KVQueryCallbacks struct from the provided callbacks.
func NewKVQueryCallbacks(valueCb func(*anypb.Any), onDoneCb func(string)) *KVQueryCallbacks {
	return &KVQueryCallbacks{
		valueCb:  valueCb,
		onDoneCb: onDoneCb,
	}
}

// RunKVQueryValueCallback runs the value callback with the given value.
func (cbs *KVQueryCallbacks) RunKVQueryValueCallback(v *anypb.Any) {
	cbs.valueCb(v)
}

// RunKVQueryOnDoneCallback runs the ondone callback.
func (cbs *KVQueryCallbacks) RunKVQueryOnDoneCallback(k string) {
	cbs.onDoneCb(k)
}

// TopicQos denotes the QoS to be used for the topic for PubSub
type TopicQos int

const (
	// Sensor signifies best effort QoS
	Sensor TopicQos = 0
	// HighReliability signifies reliable QoS
	HighReliability = 1
)

// TopicConfig contains the configuration for the topic for PubSub
type TopicConfig struct {
	Qos TopicQos
}

// PubSub is the main interface
//
// Currently the only implementation is that provided by the pubsub package.
type PubSub interface {
	// Frees the resources and unsubscribes from all topics.
	Close()

	// CreateSubscription creates a subscription to the given topic, using the exemplar proto as the
	// type expected to be called by the msg_callback.
	NewSubscription(topic string, config TopicConfig, exemplar proto.Message,
		msgCallback func(proto.Message), errCallback func(string, error)) (Subscription, error)

	// NewPublisher creates a new publisher used for publishing messages.
	NewPublisher(topic string, config TopicConfig) (Publisher, error)
}

// Subscription is a handle for a created PubSub subscription
type Subscription interface {
	// TopicName returns the name of the topic for the subscription.
	TopicName() string
	// Close closes out the subscription
	Close()
}

// Publisher is a handle for a created PubSub publisher
type Publisher interface {
	// Publish publishes the message
	Publish(msg proto.Message) error
	// TopicName returns the name of the topic for the subscription
	TopicName() string
	// Close closes out the Publisher
	Close()

	// HasMatchingSubscribers returns true if there are subscribers for this topic.
	HasMatchingSubscribers() (bool, error)
}
