# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "icon",
    srcs = [
        "action.go",
        "client.go",
        "condition.go",
        "event.go",
        "reaction.go",
        "reactiondata.go",
        "serverconfig.go",
        "session.go",
        "state_variable_path.go",
        "stream.go",
    ],
    importpath = "intrinsic/icon/go/icon",
    deps = [
        ":intsequence",
        "//intrinsic/icon/proto/v1:condition_types_go_proto",
        "//intrinsic/icon/proto/v1:service_go_proto",
        "//intrinsic/icon/proto/v1:types_go_proto",
        "//intrinsic/logging/proto:context_go_proto",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//types/known/anypb",
        "@org_golang_google_protobuf//types/known/emptypb",
    ],
)

go_library(
    name = "intsequence",
    srcs = [
        "intsequence.go",
    ],
    importpath = "intrinsic/icon/go/intsequence",
)

go_library(
    name = "sessionutil",
    srcs = [
        "sessionutil.go",
    ],
    importpath = "intrinsic/icon/go/sessionutil",
    deps = [
        ":icon",
    ],
)
