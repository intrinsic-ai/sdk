# Copyright 2023 Intrinsic Innovation LLC

load("@com_google_protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")
load("@com_google_protobuf//bazel:proto_library.bzl", "proto_library")
load("@com_google_protobuf//bazel:py_proto_library.bzl", "py_proto_library")
load("@rules_cc//cc:cc_binary.bzl", "cc_binary")
load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_python//python:defs.bzl", "py_binary", "py_library")
load("//bazel:container.bzl", "container_image")
load("//bazel:python_oci_image.bzl", "python_oci_image")
load("//intrinsic/tools/inbuild/integration_tests:inbuild.bzl", "inbuild_service_bundle", "inbuild_skill_bundle", "inbuild_skill_generate_config", "inbuild_skill_generate_entrypoint_cc", "inbuild_skill_generate_entrypoint_py")

proto_library(
    name = "inbuild_skill_proto",
    srcs = ["inbuild_skill.proto"],
    visibility = ["//intrinsic:internal_api_users"],
)

proto_library(
    name = "inbuild_service_proto",
    srcs = ["inbuild_service.proto"],
    visibility = ["//intrinsic:internal_api_users"],
)

###################
# Python Skill
###################
py_proto_library(
    name = "inbuild_skill_py_pb2",
    visibility = ["//intrinsic:internal_api_users"],
    deps = [":inbuild_skill_proto"],
)

inbuild_skill_generate_entrypoint_py(
    name = "skill_service_py",
    manifest = "inbuild_skill_py.manifest.pbtxt",
)

inbuild_skill_generate_config(
    name = "skill_service_config_py",
    manifest = "inbuild_skill_py.manifest.pbtxt",
    proto = ":inbuild_skill_proto",
)

py_library(
    name = "inbuild_skill_py_lib",
    srcs = ["inbuild_skill.py"],
    deps = [
        ":inbuild_skill_py_pb2",
        "//intrinsic/skills/python:skill_interface",
        "//intrinsic/util:decorators",
    ],
)

py_binary(
    name = "inbuild_skill_py_bin",
    srcs = [":skill_service_py"],
    main = ":skill_service_py.py",
    deps = [
        ":inbuild_skill_py_lib",
        ":inbuild_skill_py_pb2",
        "//intrinsic/skills/generator:app",
        "//intrinsic/skills/internal:runtime_data_py",
        "//intrinsic/skills/internal:single_skill_factory_py",
        "//intrinsic/skills/internal:skill_init_py",
        "//intrinsic/skills/internal:skill_service_config_utils_py",
        "//intrinsic/skills/proto:skill_service_config_py_pb2",
        "//intrinsic/util/status:status_specs_py",
        "@com_google_absl_py//absl/flags",
    ],
)

python_oci_image(
    name = "inbuild_skill_py_image",
    base = "@distroless_python3",
    binary = ":inbuild_skill_py_bin",
    data_path = "/",
    files = [
        ":skill_service_config_py",
    ],
    symlinks = {
        "/skills/skill_service": "/intrinsic/tools/inbuild/integration_tests/inbuild_skill_py_bin",
        "/skills/skill_service_config.proto.bin": "/intrinsic/tools/inbuild/integration_tests/skill_service_config_py.pbbin",
    },
)

inbuild_skill_bundle(
    name = "inbuild_skill_py",
    manifest = "inbuild_skill_py.manifest.pbtxt",
    oci_image = "inbuild_skill_py_image.tar",
    proto = ":inbuild_skill_proto",
    visibility = ["//intrinsic:internal_api_users"],
)

###################
# C++ Skill
###################
cc_proto_library(
    name = "inbuild_skill_cc_proto",
    visibility = ["//intrinsic:internal_api_users"],
    deps = [":inbuild_skill_proto"],
)

cc_library(
    name = "inbuild_skill_cc_lib",
    srcs = ["inbuild_skill.cc"],
    hdrs = ["inbuild_skill.h"],
    deps = [
        ":inbuild_skill_cc_proto",
        "//intrinsic/skills/cc:execute_context",
        "//intrinsic/skills/cc:skill_interface",
        "//intrinsic/skills/cc:skill_interface_utils",
        "//intrinsic/skills/proto:equipment_cc_proto",
        "//intrinsic/skills/proto:skill_service_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/status:statusor",
        "@com_google_protobuf//:protobuf",
    ],
)

inbuild_skill_generate_entrypoint_cc(
    name = "skill_service_cc",
    cc_header = "intrinsic/tools/inbuild/integration_tests/inbuild_skill.h",
    manifest = "inbuild_skill_cc.manifest.pbtxt",
)

inbuild_skill_generate_config(
    name = "skill_service_config_cc",
    manifest = "inbuild_skill_cc.manifest.pbtxt",
    proto = ":inbuild_skill_proto",
)

cc_binary(
    name = "inbuild_skill_cc_bin",
    srcs = [":skill_service_cc"],
    deps = [
        ":inbuild_skill_cc_lib",
        "//intrinsic/icon/release/portable:init_intrinsic",
        "//intrinsic/skills/internal:runtime_data",
        "//intrinsic/skills/internal:single_skill_factory",
        "//intrinsic/skills/internal:skill_init",
        "//intrinsic/skills/internal:skill_service_config_utils",
        "//intrinsic/util/grpc",
        "//intrinsic/util/status:status_specs",
        "@abseil-cpp//absl/flags:flag",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/time",
    ],
)

container_image(
    name = "inbuild_skill_cc_image",
    base = "@distroless_base",
    data_path = "/",
    files = [
        ":inbuild_skill_cc_bin",
        ":skill_service_config_cc",
    ],
    symlinks = {
        "/skills/skill_service": "/intrinsic/tools/inbuild/integration_tests/inbuild_skill_cc_bin",
        "/skills/skill_service_config.proto.bin": "/intrinsic/tools/inbuild/integration_tests/skill_service_config_cc.pbbin",
    },
)

inbuild_skill_bundle(
    name = "inbuild_skill_cc",
    manifest = "inbuild_skill_cc.manifest.pbtxt",
    oci_image = "inbuild_skill_cc_image.tar",
    proto = ":inbuild_skill_proto",
    visibility = ["//intrinsic:internal_api_users"],
)

###################
# Python Service
###################

py_proto_library(
    name = "inbuild_service_py_pb2",
    visibility = ["//intrinsic:internal_api_users"],
    deps = [":inbuild_service_proto"],
)

py_binary(
    name = "inbuild_service_py_bin",
    srcs = ["inbuild_service.py"],
    main = "inbuild_service.py",
    deps = [
        ":inbuild_service_py_pb2",
        "//intrinsic/resources/proto:runtime_context_py_pb2",
    ],
)

python_oci_image(
    name = "inbuild_service_py_image",
    base = "@distroless_python3",
    binary = ":inbuild_service_py_bin",
    data_path = "/",
)

inbuild_service_bundle(
    name = "inbuild_service_py",
    default_config = "inbuild_service_default_config.textproto",
    manifest = "inbuild_service_py.manifest.textproto",
    oci_image = "inbuild_service_py_image.tar",
    proto = ":inbuild_service_proto",
    visibility = ["//intrinsic:internal_api_users"],
)

###################
# C++ Service
###################

cc_proto_library(
    name = "inbuild_service_cc_proto",
    visibility = ["//intrinsic:internal_api_users"],
    deps = [":inbuild_service_proto"],
)

cc_binary(
    name = "inbuild_service_cc_bin",
    srcs = ["inbuild_service.cc"],
    deps = [
        ":inbuild_service_cc_proto",
        "//intrinsic/icon/release:file_helpers",
        "//intrinsic/icon/release/portable:init_intrinsic",
        "//intrinsic/resources/proto:runtime_context_cc_proto",
        "//intrinsic/util/proto:any",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings:string_view",
        "@abseil-cpp//absl/time",
    ],
)

container_image(
    name = "inbuild_service_cc_image",
    base = "@distroless_base",
    cmd = ["intrinsic/tools/inbuild/integration_tests/inbuild_service_cc_bin"],
    data_path = "/",
    files = [":inbuild_service_cc_bin"],
)

inbuild_service_bundle(
    name = "inbuild_service_cc",
    default_config = "inbuild_service_default_config.textproto",
    manifest = "inbuild_service_cc.manifest.textproto",
    oci_image = "inbuild_service_cc_image.tar",
    proto = ":inbuild_service_proto",
    visibility = ["//intrinsic:internal_api_users"],
)
