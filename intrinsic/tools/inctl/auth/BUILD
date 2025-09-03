# Copyright 2023 Intrinsic Innovation LLC

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "auth",
    srcs = [
        "auth.go",
        "cloud.go",
        "perrpccreds.go",
        "token_source.go",
        "tokenshttp.go",
    ],
    deps = [
        "//intrinsic/config:environments",
        "//intrinsic/kubernetes/acl:cookies",
        "//intrinsic/kubernetes/acl:identity",
        "//intrinsic/kubernetes/acl:jwt",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_spf13_viper//:go_default_library",
        "@io_opencensus_go//plugin/ocgrpc:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
    ],
)

go_library(
    name = "authtest",
    testonly = True,
    srcs = ["authtest.go"],
    deps = [":auth"],
)
