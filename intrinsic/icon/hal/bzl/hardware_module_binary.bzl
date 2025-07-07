# Copyright 2023 Intrinsic Innovation LLC

"""
Create an executable for a given plugin implementation.

The rule adds necessary driver code around the plugin interface implementation
to guarantee a guided execution of the plugin. This further lets a plugin developer focus on
the essential implementation of the interface rather than duplicating boilerplate code.
"""

def hardware_module_binary(
        name,
        hardware_module_lib,
        **kwargs):
    """Creates a binary for a hardware module.

    This can be run directly, as a standard hardware module, or as a resource.

    Args:
      name: The name of the binary.
      hardware_module_lib: The C++ library that defines the hardware module to
          generate an image for.
      **kwargs: Additional arguments to pass to cc_binary.
    """
    native.cc_binary(
        name = name,
        srcs = [Label("//intrinsic/icon/hal:hardware_module_main")],
        deps = [hardware_module_lib] + [
            "@abseil-cpp//absl/base:nullability",
            "@abseil-cpp//absl/container:flat_hash_set",
            "@abseil-cpp//absl/flags:flag",
            "@abseil-cpp//absl/log",
            "@abseil-cpp//absl/log:check",
            "@abseil-cpp//absl/log:flags",
            "@abseil-cpp//absl/status",
            "@abseil-cpp//absl/status:statusor",
            "@abseil-cpp//absl/strings",
            "@abseil-cpp//absl/time",
            Label("//intrinsic/icon/control:realtime_clock_interface"),
            Label("//intrinsic/icon/hal:hardware_module_config_cc_proto"),
            Label("//intrinsic/icon/hal:hardware_module_health_service"),
            Label("//intrinsic/icon/hal:hardware_module_init_context"),
            Label("//intrinsic/icon/hal:hardware_module_main_util"),
            Label("//intrinsic/icon/hal:hardware_module_registry"),
            Label("//intrinsic/icon/hal:hardware_module_runtime"),
            Label("//intrinsic/icon/hal:hardware_module_util"),
            Label("//intrinsic/icon/hal:module_config"),
            Label("//intrinsic/icon/interprocess/shared_memory_manager"),
            Label("//intrinsic/icon/release/portable:init_xfa_absl"),
            Label("//intrinsic/icon/release:file_helpers"),
            Label("//intrinsic/icon/utils:shutdown_signals"),
            Label("//intrinsic/logging:data_logger_client"),
            Label("//intrinsic/resources/proto:resource_registry_cc_proto"),
            Label("//intrinsic/resources/proto:runtime_context_cc_proto"),
            Label("//intrinsic/util/proto:any"),
            Label("//intrinsic/util/proto:get_text_proto"),
            Label("//intrinsic/util/status:status_builder"),
            Label("//intrinsic/util/status:status_macros"),
            Label("//intrinsic/util/thread:util"),
            Label("//intrinsic/util:memory_lock"),
        ],
        **kwargs
    )
