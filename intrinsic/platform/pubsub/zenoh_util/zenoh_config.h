// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_PLATFORM_PUBSUB_ZENOH_UTIL_ZENOH_CONFIG_H_
#define INTRINSIC_PLATFORM_PUBSUB_ZENOH_UTIL_ZENOH_CONFIG_H_

#include <fstream>
#include <ios>
#include <iostream>
#include <string>

#include "absl/flags/declare.h"
#include "absl/flags/flag.h"
#include "absl/log/log.h"
#include "absl/strings/string_view.h"
#include "intrinsic/platform/pubsub/zenoh_util/zenoh_helpers.h"

ABSL_DECLARE_FLAG(std::string, zenoh_router);

namespace intrinsic {

// The ZenohConfig object holds the two state variables that are potentially
// adjusted when the config is generated: the "connect" endpoint, which can be
// used to adjust the IP:port of the Zenoh router that this Zenoh peer tries to
// connect to, and the "listen" endpoint, which can be used to limit which
// interface, if any, that Zenoh will listen on for incoming connections. It
// is important in some test environments to be able to disable listening.
// This is detected and handled appropriately in the GetZenohPeerConfig()
// implementation.

class ZenohConfig {
 public:
  std::string GenerateJsonString() const;

  void set_connect_endpoint(absl::string_view endpoint) {
    connect_endpoint_ = endpoint;
  }

  void set_listen_endpoint(absl::string_view endpoint) {
    listen_endpoint_ = endpoint;
  }

 private:
  std::string connect_endpoint_ =
      "tcp/zenoh-router.app-intrinsic-base.svc.cluster.local:7447";
  std::string listen_endpoint_ = "tcp/0.0.0.0:0";
};

// GetZenohPeerConfig() uses a ZenohConfig instance to generate a JSON string
// that is appropriate for the runtime environment and which uses the
// requested connect and listen endpoints in env vars or absl flags.

std::string GetZenohPeerConfig(absl::string_view router_override = "");

}  // namespace intrinsic

#endif  // INTRINSIC_PLATFORM_PUBSUB_ZENOH_UTIL_ZENOH_CONFIG_H_
