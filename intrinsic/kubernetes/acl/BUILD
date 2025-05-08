# Copyright 2023 Intrinsic Innovation LLC

# An acl engine using Authzed SpiceDB.

load("//bazel:go_macros.bzl", "go_library", "go_test")

package(default_visibility = [
    "//visibility:public",
])

go_library(
    name = "cookies",
    srcs = ["cookies.go"],
    deps = [
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
    ],
)

go_test(
    name = "cookies_test",
    srcs = ["cookies_test.go"],
    library = ":cookies",
    deps = [
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
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

go_test(
    name = "org_test",
    srcs = ["org_test.go"],
    library = ":org",
    deps = ["@com_github_google_go_cmp//cmp:go_default_library"],
)

go_test(
    name = "jwt_test",
    srcs = ["jwt_test.go"],
    library = ":jwt",
    deps = [
        "//intrinsic/kubernetes/acl/testing:jwttesting",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_google_go_cmp//cmp/cmpopts:go_default_library",
    ],
)

go_library(
    name = "identity",
    srcs = [
        "identity.go",
    ],
    deps = [
        ":cookies",
        ":org",
        "//intrinsic/frontend/go:origin",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
    ],
)

go_test(
    name = "identity_test",
    srcs = [
        "identity_test.go",
    ],
    library = ":identity",
    deps = [
        ":cookies",
        ":org",
        "@org_golang_google_grpc//metadata:go_default_library",
    ],
)
