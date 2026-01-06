// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/platform/pubsub/topic_config.h"

#include <string>

namespace intrinsic {

std::string PubSubQoSToZenohQos(const TopicConfig::TopicQoS& qos) {
  return qos == TopicConfig::TopicQoS::Sensor ? "Sensor" : "HighReliability";
}

}  // namespace intrinsic
