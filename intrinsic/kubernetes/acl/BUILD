# Copyright 2023 Intrinsic Innovation LLC

# An acl engine using Authzed SpiceDB.

load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = [
    "//visibility:public",
])

go_library(
    name = "cookies",
    srcs = ["cookies.go"],
    deps = [
        "//intrinsic/stats/go:telemetry",
        "@com_github_golang_glog//:go_default_library",
        "@io_opencensus_go//trace:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
    ],
)

go_library(
    name = "jwt",
    srcs = ["jwt.go"],
    deps = ["@com_github_pkg_errors//:go_default_library"],
)

go_library(
    name = "org",
    srcs = ["org.go"],
    deps = ["@com_github_golang_glog//:go_default_library"],
)

go_library(
    name = "identity",
    srcs = ["identity.go"],
    deps = [
        ":cookies",
        ":jwt",
        ":org",
        "//intrinsic/stats/go:telemetry",
        "@com_github_golang_glog//:go_default_library",
        "@io_opencensus_go//trace:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
