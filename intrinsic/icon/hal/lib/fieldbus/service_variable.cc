// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/icon/hal/lib/fieldbus/service_variable.h"

#include <cstddef>
#include <cstdint>
#include <functional>

#include "absl/status/status.h"

namespace intrinsic::fieldbus {

ServiceVariable::ServiceVariable(
    std::function<absl::Status(uint8_t*, std::size_t)> service_variable_read,
    std::function<absl::Status(const uint8_t*, std::size_t)>
        service_variable_write)
    : service_variable_read_(service_variable_read),
      service_variable_write_(service_variable_write) {}

}  // namespace intrinsic::fieldbus
