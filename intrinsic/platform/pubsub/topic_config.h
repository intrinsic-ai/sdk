// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_PLATFORM_PUBSUB_TOPIC_CONFIG_H_
#define INTRINSIC_PLATFORM_PUBSUB_TOPIC_CONFIG_H_

#include <string>

namespace intrinsic {

struct TopicConfig {
  enum TopicQoS {
    Sensor = 0,
    HighReliability = 1,
  };

  TopicQoS topic_qos = HighReliability;
};

std::string PubSubQoSToZenohQos(const TopicConfig::TopicQoS& qos);

}  // namespace intrinsic

#endif  // INTRINSIC_PLATFORM_PUBSUB_TOPIC_CONFIG_H_
