# Copyright 2023 Intrinsic Innovation LLC

load("@aspect_rules_js//js:proto.bzl", "js_proto_toolchain")
load("@aspect_rules_ts//ts:defs.bzl", "ts_config")
load("@npm//:defs.bzl", "npm_link_all_packages")
load("@npm//intrinsic/frontend/sdk:@bufbuild/protoc-gen-es/package_json.bzl", "bin")

package(default_visibility = ["//visibility:public"])

exports_files(
    srcs = ["LICENSE"],
    visibility = ["//intrinsic/frontend/sdk:__subpackages__"],
)

npm_link_all_packages(name = "node_modules")

bin.protoc_gen_es_binary(
    name = "protoc_gen_es",
)

js_proto_toolchain(
    name = "gen_es_toolchain",
    out_dts_extension = "_pb.d.ts",
    out_js_extension = "_pb.js",
    plugin_bin = ":protoc_gen_es",
    plugin_name = "es",
    plugin_options = [
        "keep_empty_files=true",
        "target=js+dts",
    ],
    runtime = "//intrinsic/frontend/sdk:node_modules/@bufbuild/protobuf",
    visibility = ["//visibility:public"],
)

ts_config(
    name = "tsconfig",
    src = "tsconfig.json",
    visibility = ["//visibility:public"],
)
