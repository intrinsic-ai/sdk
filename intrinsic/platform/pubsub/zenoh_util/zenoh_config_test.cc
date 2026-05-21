// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/platform/pubsub/zenoh_util/zenoh_config.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <nlohmann/json.hpp>
#include <string>

#include "absl/flags/flag.h"
#include "intrinsic/util/testing/gtest_wrapper.h"

using ::testing::HasSubstr;

namespace intrinsic {

TEST(ZenohConfigTest, TestRouterFlag) {
  const std::string test_router_endpoint("tcp/foo.bar.baz:12345");
  absl::SetFlag(&FLAGS_zenoh_router, test_router_endpoint);
  std::string config(GetZenohPeerConfig());
  EXPECT_THAT(config, HasSubstr(test_router_endpoint));
}

TEST(ZenohConfigTest, TestGenerateJsonString) {
  ZenohConfig config;
  config.set_connect_endpoint("tcp/foo:123");
  config.set_listen_endpoint({});
  std::string json_str = config.GenerateJsonString();
  auto j = nlohmann::json::parse(json_str);
  EXPECT_EQ(j["connect"]["endpoints"][0], "tcp/foo:123");
  EXPECT_EQ(j["listen"]["endpoints"].size(), 0);
  EXPECT_EQ(j["mode"], "peer");
}

TEST(ZenohConfigTest, TestDefaultJsonString) {
  ZenohConfig config;
  std::string json_str = config.GenerateJsonString();
  auto j = nlohmann::json::parse(json_str);
  EXPECT_EQ(j["connect"]["endpoints"][0],
            "tcp/zenoh-router.app-intrinsic-base.svc.cluster.local:7447");
  EXPECT_EQ(j["listen"]["endpoints"][0], "tcp/0.0.0.0:0");
  EXPECT_EQ(j["mode"], "peer");
}

TEST(ZenohConfigTest, TestMatchesPeerConfigJsonVerbatim) {
  const char* kPeerConfigJson = R"json({
  "connect": {
    "endpoints": ["tcp/zenoh-router.app-intrinsic-base.svc.cluster.local:7447"]
  },
  "imw": {
    "introspection": {
      "enable": true
    }
  },
  "listen": {
    "endpoints": ["tcp/0.0.0.0:0"]
  },
  "mode": "peer",
  "plugins": {},
  "scouting": {
    "multicast": { "enabled": false },
    "gossip": {
      "enabled": true,
      "multihop": false,
      "autoconnect": { "peer": ["peer", "router"] }
    }
  },
  "transport": {
    "shared_memory": {
      "enabled": false
    }
  }
})json";

  auto expected_json = nlohmann::json::parse(kPeerConfigJson);
  ZenohConfig config;
  auto actual_json = nlohmann::json::parse(config.GenerateJsonString());
  EXPECT_EQ(actual_json, expected_json);
}

TEST(ZenohConfigTest, TestCustomEndpointsStringMatch) {
  ZenohConfig config;
  config.set_connect_endpoint("tcp/custom-router:1234");
  config.set_listen_endpoint("tcp/custom-listen:5678");

  const char* kExpectedString = R"json({
  "connect": {
    "endpoints": [
      "tcp/custom-router:1234"
    ]
  },
  "imw": {
    "introspection": {
      "enable": true
    }
  },
  "listen": {
    "endpoints": [
      "tcp/custom-listen:5678"
    ]
  },
  "mode": "peer",
  "plugins": {},
  "scouting": {
    "gossip": {
      "autoconnect": {
        "peer": [
          "peer",
          "router"
        ]
      },
      "enabled": true,
      "multihop": false
    },
    "multicast": {
      "enabled": false
    }
  },
  "transport": {
    "shared_memory": {
      "enabled": false
    }
  }
})json";

  EXPECT_EQ(config.GenerateJsonString(), kExpectedString);
}

}  // namespace intrinsic
