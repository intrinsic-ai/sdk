# Copyright 2023 Intrinsic Innovation LLC

# An acl engine using Authzed SpiceDB.
load("//bazel:go_macros.bzl", "go_library")

package(default_visibility = [
    "//visibility:public",
])

go_library(
    name = "cookies",
    srcs = ["cookies.go"],
    importpath = "intrinsic/kubernetes/acl/cookies",
    visibility = [
        "//intrinsic/frontend/experiments/mapping_provider:__subpackages__",
        "//intrinsic/kubernetes/acl:__subpackages__",
        "//intrinsic/kubernetes/datasilo/examples:__subpackages__",
        "//intrinsic/kubernetes/vmpool/manager/internal:__subpackages__",
        "//intrinsic/skills/tools/skill/cmd:__subpackages__",
        "//intrinsic/storage/content_addressable_storage/internal/service:__subpackages__",
        "//intrinsic/tools/inctl/auth:__subpackages__",
        "//intrinsic/tools/inctl/util:__subpackages__",
    ],
    deps = [
        "@com_github_golang_glog//:go_default_library",
        "@io_opencensus_go//trace:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
    ],
)

go_library(
    name = "jwt",
    srcs = ["jwt.go"],
    importpath = "intrinsic/kubernetes/acl/jwt",
    deps = [
        "@com_github_golang_jwt_jwt_v5//:jwt",
        "@com_github_pkg_errors//:go_default_library",
    ],
)

go_library(
    name = "org",
    srcs = ["org.go"],
    importpath = "intrinsic/kubernetes/acl/org",
    deps = [
        "@com_github_golang_glog//:go_default_library",
        "@com_github_rs_xid//:go_default_library",
    ],
)

go_library(
    name = "identity",
    srcs = ["identity.go"],
    importpath = "intrinsic/kubernetes/acl/identity",
    deps = [
        ":cookies",
        ":jwt",
        ":org",
        "//intrinsic/frontend/go:origin",
        "//intrinsic/stats/go:telemetry",
        "@com_github_golang_glog//:go_default_library",
        "@io_opencensus_go//trace:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
