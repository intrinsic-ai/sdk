// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/platform/pubsub/zenoh_util/zenoh_config.h"

#include <iostream>
#include <nlohmann/json.hpp>
#include <string>

#include "absl/flags/flag.h"

ABSL_FLAG(std::string, zenoh_router, "",
          "Override the default Zenoh connection to PROTOCOL/HOSTNAME:PORT");

namespace intrinsic {

/*****************************************************************************
The default values for the ZenohConfig will produce the following JSON
output from ZenohConfig::GenerateJsonString(). This JSON is created
programmatically to allow for adjustments to the connect and listen
endpoints.

{
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
}
*****************************************************************************/

std::string ZenohConfig::GenerateJsonString() const {
  nlohmann::json j;

  nlohmann::json connect_endpoints = nlohmann::json::array();
  if (!connect_endpoint_.empty()) {
    connect_endpoints.push_back(connect_endpoint_);
  }
  j["connect"]["endpoints"] = connect_endpoints;

  j["imw"]["introspection"]["enable"] = true;

  nlohmann::json listen_endpoints = nlohmann::json::array();
  if (!listen_endpoint_.empty()) {
    listen_endpoints.push_back(listen_endpoint_);
  }
  j["listen"]["endpoints"] = listen_endpoints;

  j["mode"] = "peer";
  j["plugins"] = nlohmann::json::object();
  j["scouting"]["multicast"]["enabled"] = false;
  j["scouting"]["gossip"]["enabled"] = true;
  j["scouting"]["gossip"]["multihop"] = false;
  j["scouting"]["gossip"]["autoconnect"]["peer"] = {"peer", "router"};
  j["transport"]["shared_memory"]["enabled"] = false;

  return j.dump(2);
}

// GetZenohPeerConfig is a global function for compatibility with existing
// code. It uses the ZenohConfig class to generate a JSON string, adjusting the
// endpoints as requested by env vars.
//
// The router_override parameter is used because when this C++ code is called
// from Go binaries via CGO, the Abseil flags (like FLAGS_zenoh_router) are
// not parsed and remain empty. The Go binary must parse its own flags and
// pass the override value explicitly.
std::string GetZenohPeerConfig(absl::string_view router_override) {
  ZenohConfig config;

  if (RunningUnderTest()) {
    config.set_listen_endpoint({});
  } else if (const char* allowed_ip = getenv("ALLOWED_PUBSUB_IPv4");
             allowed_ip != nullptr) {
    config.set_listen_endpoint(absl::StrFormat("tcp/%s:0", allowed_ip));
  } else if (RunningInKubernetes()) {
    config.set_listen_endpoint({});
  }

  // If requested by the zenoh_router flag or override parameter, try to alter
  // the default router connection endpoint defaults
  if (!router_override.empty()) {
    config.set_connect_endpoint(router_override);
  } else if (!absl::GetFlag(FLAGS_zenoh_router).empty()) {
    config.set_connect_endpoint(absl::GetFlag(FLAGS_zenoh_router));
  }
  return config.GenerateJsonString();
}

}  // namespace intrinsic
