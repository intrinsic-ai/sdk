// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/lib/fieldbus/async_device_request.h"

#include <string>

namespace intrinsic::fieldbus {

std::string ToString(RequestType type) {
  switch (type) {
    case RequestType::kNormalOperation:
      return "kNormalOperation";
    case RequestType::kActivate:
      return "kActivate";
    case RequestType::kDeactivate:
      return "kDeactivate";
    case RequestType::kEnableMotion:
      return "kEnableMotion";
    case RequestType::kDisableMotion:
      return "kDisableMotion";
    case RequestType::kClearFaults:
      return "kClearFaults";
  }
}
std::string ToString(RequestStatus status) {
  switch (status) {
    case RequestStatus::kDone:
      return "kDone";
    case RequestStatus::kProcessing:
      return "kProcessing";
  }
}

}  // namespace intrinsic::fieldbus
