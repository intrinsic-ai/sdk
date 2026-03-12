# Copyright 2023 Intrinsic Innovation LLC

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _non_module_deps_impl(ctx):
    http_archive(
        name = "raft",
        build_file = "//third_party/raft:BUILD.raft.bazel",
        sha256 = "6353114607b5dbc085ebc486e0770502649ced20b993c93e834946f7b2688743",
        strip_prefix = "RAFT-e7f94cbf18c9bb7402968fa2583ec90cf3a58fb0",
        urls = ["https://github.com/albertodallolio/RAFT/archive/e7f94cbf18c9bb7402968fa2583ec90cf3a58fb0.tar.gz"],
    )

    # Build rocksdb_sys from an archive because crate_universe rules don't
    # correctly express dependencies on non-system Clang. Filed
    # https://github.com/bazelbuild/rules_rust/issues/3431 upstream.
    http_archive(
        name = "rocksdb_sys",
        build_file = "//third_party/rocksdb-sys:BUILD.rocksdb-sys.bazel",
        integrity = "sha256-eHitrOJSfCy67p76H0+j0iGeG6XkEoNoCYRfENa4BQs=",
        strip_prefix = "rust-rocksdb-0.23.0",
        type = "tar.gz",
        urls = ["https://github.com/rust-rocksdb/rust-rocksdb/archive/refs/tags/v0.23.0.tar.gz"],
    )

    http_archive(
        name = "zenohd",
        build_file = "//third_party/zenohd:BUILD.zenohd.bazel",
        integrity = "sha256-V3NQiPmp5CbIwBG7ZdcxbSTaXWvaJjiEsTTRR10kSeI=",
        strip_prefix = "zenoh-1.3.3",
        type = "tar.gz",
        urls = ["https://github.com/eclipse-zenoh/zenoh/archive/refs/tags/1.3.3.tar.gz"],
    )

    # Downloading from crates.io because this crate is found alongside many crates
    # in a cargo workspace in the Zenoh repo; if we download from crates.io, we can
    # get it just by itself. Setting up a multi-crate repo using rules_rust seems
    # to be challenging.
    http_archive(
        name = "zenoh_plugin_storage_manager",
        build_file = "//third_party/zenohd:BUILD.zenoh_plugin_storage_manager.bazel",
        integrity = "sha256-IuUOduV6RIZMWmX4gwZv+QVobXTlG/WxEbj1dpKtLgw=",
        strip_prefix = "zenoh-plugin-storage-manager-1.3.3",
        type = "tar.gz",
        urls = ["https://crates.io/api/v1/crates/zenoh-plugin-storage-manager/1.3.3/download"],
    )

    # Same for zenoh-plugin-rest: downloading from crates.io seems like the
    # easiest way to make this work.
    http_archive(
        name = "zenoh_plugin_rest",
        build_file = "//third_party/zenohd:BUILD.zenoh_plugin_rest.bazel",
        integrity = "sha256-/2J04/E1lmYsktQ8w9KZQmx+O+A11Re08simzemLhR8=",
        strip_prefix = "zenoh-plugin-rest-1.3.3",
        type = "tar.gz",
        urls = ["https://crates.io/api/v1/crates/zenoh-plugin-rest/1.3.3/download"],
    )

    http_archive(
        name = "zenoh_backend_filesystem",
        build_file = "//third_party/zenohd:BUILD.zenoh_backend_filesystem.bazel",
        integrity = "sha256-G97XnocB8dpm8x6+WIWS31R2WAbKgWxQA2VU+jNh1NU=",
        strip_prefix = "zenoh-backend-filesystem-1.3.3",
        type = "tar.gz",
        urls = ["https://github.com/eclipse-zenoh/zenoh-backend-filesystem/archive/refs/tags/1.3.3.tar.gz"],
    )

    # Although it would be ideal to run cbindgen on zenoh-c to generate
    # for the architecture on hand, unfortunately this is a very deep
    # rabbit hole currently without a reasonably simple solution
    # in the rules_rust ecosystem, since cbindgen relies on Cargo to
    # generate its metadata analysis. To sidestep the need to run
    # cbindgen during the build, we run cbindgen on a x86-64 machine
    # outside of Bazel. In the future, if we want to run on multiple
    # architectures, we may need to revisit this approach.
    # Hopefully at that point there will be a simpler approach for dealing
    # with cbindgen in Bazel.
    http_archive(
        name = "zenoh_c",
        build_file = "//third_party/zenoh_c:BUILD.zenoh_c.bazel",
        integrity = "sha256-bfu+ev8hD46EzSGjeV7Bm3+gSYYiABwvNAeiJuPN89g=",
        strip_prefix = "zenoh-c-1.3.3",
        urls = ["https://github.com/eclipse-zenoh/zenoh-c/archive/refs/tags/1.3.3.tar.gz"],
        patches = [
            "//third_party/zenoh_c:modify_generated_headers.patch",
            "//third_party/zenoh_c:add_opaque_types_mod.patch",
            "//third_party/zenoh_c:add_zenoh_configure_h.patch",
            "//third_party/zenoh_c:add_zenoh_opaque_h.patch",
        ],
    )

    http_archive(
        name = "zenoh_cpp",
        build_file = "//third_party/zenoh_cpp:BUILD.zenoh_cpp.bazel",
        integrity = "sha256-kAZnjaMtZBmKa2+v03fEvLsJNiT5iEos2YkpIDiCxEQ=",
        strip_prefix = "zenoh-cpp-1.3.3",
        urls = ["https://github.com/eclipse-zenoh/zenoh-cpp/archive/refs/tags/1.3.3.tar.gz"],
    )

    http_archive(
        name = "wstunnel",
        build_file = "//third_party/wstunnel:BUILD.wstunnel.bazel",
        sha256 = "4362bb70883404f6ab78a82c862be3542718cca711807ad0d86acec629615b3f",
        strip_prefix = "wstunnel-9.4.1",
        type = "tar.gz",
        urls = ["https://github.com/erebe/wstunnel/archive/refs/tags/v9.4.1.tar.gz"],
    )

    http_archive(
        name = "libvpx",
        build_file = Label("//third_party/libvpx:BUILD.libvpx.bazel"),
        sha256 = "901747254d80a7937c933d03bd7c5d41e8e6c883e0665fadcb172542167c7977",
        strip_prefix = "libvpx-1.14.1",
        urls = ["https://github.com/webmproject/libvpx/archive/refs/tags/v1.14.1.tar.gz"],
    )

    http_archive(
        name = "libsvtav1",
        build_file = Label("//third_party/libsvtav1:BUILD.libsvtav1.bazel"),
        sha256 = "d02b54685542de0236bce4be1b50912aba68aff997c43b350d84a518df0cf4e5",
        strip_prefix = "SVT-AV1-v2.2.1",
        urls = ["https://gitlab.com/AOMediaCodec/SVT-AV1/-/archive/v2.2.1/SVT-AV1-v2.2.1.tar.gz"],
    )

    http_archive(
        name = "ffmpeg",
        build_file = Label("//third_party/ffmpeg:BUILD.ffmpeg.bazel"),
        sha256 = "5eb46d18d664a0ccadf7b0adee03bd3b7fa72893d667f36c69e202a807e6d533",
        strip_prefix = "FFmpeg-n7.0.2",
        urls = ["https://github.com/FFmpeg/FFmpeg/archive/refs/tags/n7.0.2.tar.gz"],
    )

    # Chrome for Testing, required by Cypress.
    # Find new versions at https://googlechromelabs.github.io/chrome-for-testing/

    http_archive(
        name = "chrome_linux",
        build_file_content = """filegroup(
    name = "all",
    srcs = glob(["**"]),
    visibility = ["//visibility:public"],
    )""",
        sha256 = "5023ec2b8995b74caa5de0e22d5e30f871c3ecce67a6e52d3c9f9dfed423ed01",
        strip_prefix = "chrome-linux64",
        urls = [
            "https://storage.googleapis.com/chrome-for-testing-public/141.0.7390.54/linux64/chrome-linux64.zip",
        ],
    )

non_module_deps_ext = module_extension(implementation = _non_module_deps_impl)
