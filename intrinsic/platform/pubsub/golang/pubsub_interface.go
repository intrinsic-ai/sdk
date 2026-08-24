// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	// NewSubscription creates a subscription to the given topic, using the exemplar proto as the
	// type expected to be called by the msgCallback.
	// The errCallback is invoked when unmarshaling the payload fails; its first argument receives
	// the raw packet bytes (as a string) and its second argument receives the unmarshal error.
	NewSubscription(topic string, config TopicConfig, exemplar proto.Message,
		msgCallback func(proto.Message), errCallback func(string, error)) (Subscription, error)

	// NewRawSubscription creates a subscription to the given topic, passing the full packet to callback.
	NewRawSubscription(topic string, config TopicConfig, callback func(*pubsubpb.PubSubPacket)) (Subscription, error)

	// NewPublisher creates a new publisher used for publishing messages.
	NewPublisher(topic string, config TopicConfig) (Publisher, error)

	// DeclareLivelinessToken declares a liveliness token on the given key expression.
	DeclareLivelinessToken(keyExpr string) error

	// DropLivelinessToken drops the previously declared liveliness token on the given key expression.
	DropLivelinessToken(keyExpr string) error

	// NewLivelinessSubscription subscribes to liveliness notifications matching
	// the given key expression.
	//
	// Parameters:
	//  - keyExpr: key expression to subscribe to. May contain wildcards.
	//  - notifyAboutExistingTokens: whether to receive notifications about tokens
	//      that were declared before the subscription was created.
	//  - msgCallback: called when liveliness status changes.
	//
	// Callback parameters:
	//  - key expression whose liveliness token's status changed.
	//  - boolean value for the liveliness status (true means alive).
	NewLivelinessSubscription(keyExpr string, notifyAboutExistingTokens bool, msgCallback func(string, bool)) (LivelinessSubscription, error)

	// LivelinessGet searches for currently available liveliness tokens matching the given key expression.
	// Returns immediately. Callbacks are called in different goroutines.
	//
	// Parameters:
	//  - keyExpr: key expression for matching liveliness tokens. May contain wildcards.
	//  - callback: called when a liveliness token is found.
	//     - Takes the key expression that the token was declared on.
	//  - onDone: called when all tokens matching `keyExpr` have been found.
	//     - Takes the original key expression (same as `keyExpr`).
	//
	// WARNING: both callbacks (`callback` and `onDone`) must return as quickly
	// as possible. Long-running callbacks may affect results of this
	// `LivelinessGet` call (some callbacks may be skipped). They may also affect
	// other clients of the PubSub API (e.g. PubSub may start dropping messages).
	// If the total duration of callbacks invoked by the same `LivelinessGet`
	// query exceeds 10 minutes, then no more callbacks will be called after 10
	// minutes.
	LivelinessGet(keyExpr string, callback func(string), onDone func(string)) (LivelinessGetQuery, error)

	// LivelinessGetAllSynchronous searches for currently available liveliness tokens matching the given
	// key expression.
	// Blocks until all matching tokens are found.
	//
	// Parameters:
	//  - keyExpr: key expression for matching liveliness tokens. May contain wildcards.
	//
	// On success, LivelinessGetAllSynchronous returns a list of all found liveliness tokens
	// in _arbitrary_ order.
	LivelinessGetAllSynchronous(keyExpr string) ([]string, error)
}

// Subscription is a handle for a created PubSub subscription
type Subscription interface {
	// TopicName returns the name of the topic for the subscription.
	TopicName() string
	// Close closes out the subscription
	Close()
}

// LivelinessSubscription represents a subscription to liveliness
// notifications.
type LivelinessSubscription interface {
	// Close closes the subscription
	Close()
}

// LivelinessQuery is a handle for the LivelinessGet query.
// Keep it alive until the query completes.
type LivelinessGetQuery interface {

	// Closes the query and frees up underlying resources.
	// Call this method after the query's onDone callback returns.
	Close()
}

// Publisher is a handle for a created PubSub publisher
type Publisher interface {
	// Publish publishes the message.
	// The message is wrapped into anypb.Any before publishing.
	Publish(msg proto.Message) error

	// Publish publishes the given Any proto as is.
	PublishAny(msg *anypb.Any) error

	// TopicName returns the name of the topic for the subscription
	TopicName() string
	// Close closes out the Publisher
	Close()

	// HasMatchingSubscribers returns true if there are subscribers for this topic.
	HasMatchingSubscribers() (bool, error)
}
