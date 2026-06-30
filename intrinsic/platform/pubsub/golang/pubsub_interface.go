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
