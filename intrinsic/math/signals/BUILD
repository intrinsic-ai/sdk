# Copyright 2023 Intrinsic Innovation LLC

load("@rules_cc//cc:cc_library.bzl", "cc_library")
load("@rules_cc//cc:cc_test.bzl", "cc_test")
load("//intrinsic/icon/testing:cc_test_and_malloc_test.bzl", "cc_test_and_malloc_test")

package(default_visibility = ["//visibility:public"])

cc_library(
    name = "adaptive_quadrature",
    srcs = ["adaptive_quadrature.cc"],
    hdrs = ["adaptive_quadrature.h"],
    deps = [
        ":nested_quadrature",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/algorithm:container",
        "@abseil-cpp//absl/functional:function_ref",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_test(
    name = "adaptive_quadrature_test",
    srcs = ["adaptive_quadrature_test.cc"],
    deps = [
        ":adaptive_quadrature",
        ":nested_quadrature",
        "//intrinsic/util/status:status_macros",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/algorithm:container",
        "@abseil-cpp//absl/functional:any_invocable",
        "@abseil-cpp//absl/functional:function_ref",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_library(
    name = "butter_filter2",
    srcs = [
        "butter_filter2.cc",
    ],
    hdrs = [
        "butter_filter2.h",
    ],
    visibility = ["//intrinsic/icon:rtcl_users"],
    deps = [
        "//intrinsic/icon/utils:log",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/log:check",
    ],
)

cc_test_and_malloc_test(
    name = "butter_filter2_test",
    srcs = [
        "butter_filter2_test.cc",
    ],
    deps = [
        ":butter_filter2",
        "//intrinsic/eigenmath",
    ],
)

cc_library(
    name = "common_iir_filters",
    srcs = [
        "common_iir_filters.cc",
    ],
    hdrs = [
        "common_iir_filters.h",
    ],
    deps = [
        ":iir_filt",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_test(
    name = "common_iir_filters_test",
    srcs = [
        "common_iir_filters_test.cc",
    ],
    deps = [
        ":common_iir_filters",
        ":iir_filt",
        "//intrinsic/icon/utils:realtime_status_matchers",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
    ],
)

cc_library(
    name = "dft_utils",
    srcs = ["dft_utils.cc"],
    hdrs = ["dft_utils.h"],
    deps = [
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_macro",
        "//intrinsic/icon/utils:realtime_status_or",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_test_and_malloc_test(
    name = "dft_utils_test",
    srcs = ["dft_utils_test.cc"],
    deps = [
        ":dft_utils",
        "//intrinsic/icon/testing:malloc_test",
        "//intrinsic/icon/utils:realtime_status_matchers",
        "//intrinsic/math:complex_matchers",
    ],
)

cc_library(
    name = "fft",
    srcs = ["fft.cc"],
    hdrs = ["fft.h"],
    copts = ["-fexceptions"],  # pocketfft can throw exceptions
    features = ["-use_header_modules"],  # Incompatible with -fexceptions.
    deps = [
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/types:span",
        "@pocketfft",
    ],
)

cc_test(
    name = "fft_test",
    srcs = ["fft_test.cc"],
    deps = [
        ":fft",
        "//intrinsic/math:complex_matchers",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_library(
    name = "gauss_kronrod_quadrature",
    srcs = ["gauss_kronrod_quadrature.cc"],
    hdrs = ["gauss_kronrod_quadrature.h"],
    deps = [
        ":nested_quadrature",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/functional:function_ref",
        "@abseil-cpp//absl/log",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_test(
    name = "gauss_kronrod_quadrature_test",
    srcs = ["gauss_kronrod_quadrature_test.cc"],
    deps = [
        ":gauss_kronrod_quadrature",
        ":nested_quadrature",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//third_party/imported/cpp_libraries/math:mathutil",
        "@abseil-cpp//absl/functional:any_invocable",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_library(
    name = "gauss_legendre_quadrature",
    srcs = ["gauss_legendre_quadrature.cc"],
    hdrs = ["gauss_legendre_quadrature.h"],
    deps = [
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/functional:function_ref",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_test(
    name = "gauss_legendre_quadrature_test",
    srcs = ["gauss_legendre_quadrature_test.cc"],
    deps = [
        ":gauss_legendre_quadrature",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "high_frequency_monitor",
    srcs = ["high_frequency_monitor.cc"],
    hdrs = ["high_frequency_monitor.h"],
    deps = [
        ":dft_utils",
        ":rms_monitor",
        ":sliding_dft",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_macro",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_test_and_malloc_test(
    name = "high_frequency_monitor_test",
    srcs = ["high_frequency_monitor_test.cc"],
    deps = [
        ":high_frequency_monitor",
        "//intrinsic/icon/testing:malloc_test",
        "//intrinsic/icon/utils:realtime_status_matchers",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_library(
    name = "iir_filt",
    hdrs = [
        "iir_filt.h",
    ],
    deps = [
        "//intrinsic/eigenmath",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_macro",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/math:almost_equals",
        "//intrinsic/math/proto:iir_filter_config_cc_proto",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/algorithm:container",
        "@abseil-cpp//absl/log:check",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_test_and_malloc_test(
    name = "iir_filt_test",
    srcs = [
        "iir_filt_test.cc",
    ],
    deps = [
        ":butter_filter2",
        ":iir_filt",
        "//intrinsic/eigenmath",
        "//intrinsic/eigenmath:matchers",
        "//intrinsic/eigenmath:random_number_utils",
        "//intrinsic/icon/testing:malloc_test",
        "//intrinsic/icon/utils:realtime_status_matchers",
        "//intrinsic/math/proto:iir_filter_config_cc_proto",
        "//intrinsic/math/signals/testing:iir_filter_test_utils",
        "@abseil-cpp//absl/random",
        "@abseil-cpp//absl/random:distributions",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_library(
    name = "iir_filtfilt",
    hdrs = [
        "iir_filtfilt.h",
    ],
    deps = [
        ":butter_filter2",
        ":iir_filt",
        "//intrinsic/util/status:status_macros",
        "@abseil-cpp//absl/algorithm:container",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_test(
    name = "iir_filtfilt_test",
    srcs = [
        "iir_filtfilt_test.cc",
    ],
    deps = [
        ":butter_filter2",
        ":iir_filt",
        ":iir_filtfilt",
        "//intrinsic/eigenmath",
        "//intrinsic/eigenmath:matchers",
        "//intrinsic/eigenmath:random_number_utils",
        "//intrinsic/math/signals/testing:iir_filter_test_utils",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "@abseil-cpp//absl/random",
        "@abseil-cpp//absl/random:distributions",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_library(
    name = "nested_quadrature",
    hdrs = ["nested_quadrature.h"],
    deps = [
        "@abseil-cpp//absl/functional:function_ref",
        "@abseil-cpp//absl/status:statusor",
    ],
)

cc_library(
    name = "numerical_differentiation",
    hdrs = ["numerical_differentiation.h"],
    deps = [
        "//intrinsic/icon/utils:fixed_str_cat",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_or",
        "//third_party/imported/cpp_libraries/math:mathutil",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
    ],
)

cc_test(
    name = "numerical_differentiation_test",
    srcs = ["numerical_differentiation_test.cc"],
    deps = [
        ":numerical_differentiation",
        "//intrinsic/eigenmath",
        "//intrinsic/eigenmath:matchers",
        "//intrinsic/icon/testing:malloc_test",
        "//intrinsic/icon/utils:realtime_status_matchers",
        "//intrinsic/util/testing:gtest_wrapper_main",
        "//third_party/imported/cpp_libraries/math:mathutil",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_library(
    name = "numerical_integration",
    hdrs = ["numerical_integration.h"],
    deps = [
        "//intrinsic/icon/utils:realtime_status_macro",
        "//intrinsic/icon/utils:realtime_status_or",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
    ],
)

cc_test(
    name = "numerical_integration_test",
    srcs = ["numerical_integration_test.cc"],
    deps = [
        ":numerical_integration",
        "//intrinsic/eigenmath",
        "//intrinsic/eigenmath:matchers",
        "//intrinsic/icon/testing:malloc_test",
        "//intrinsic/icon/utils:realtime_status_matchers",
        "//intrinsic/icon/utils:realtime_status_or",
        "//intrinsic/util/testing:gtest_wrapper_main",
    ],
)

cc_library(
    name = "phase_locked_loop",
    srcs = ["phase_locked_loop.cc"],
    hdrs = ["phase_locked_loop.h"],
    deps = ["@abseil-cpp//absl/time"],
)

cc_test_and_malloc_test(
    name = "phase_locked_loop_test",
    srcs = ["phase_locked_loop_test.cc"],
    deps = [
        ":phase_locked_loop",
        "//intrinsic/kinematics/testing:trajectory_plot",
        "@abseil-cpp//absl/random",
        "@abseil-cpp//absl/random:distributions",
        "@abseil-cpp//absl/time",
    ],
)

cc_library(
    name = "planck_taper_window",
    srcs = ["planck_taper_window.cc"],
    hdrs = ["planck_taper_window.h"],
)

cc_test(
    name = "planck_taper_window_test",
    srcs = ["planck_taper_window_test.cc"],
    deps = [
        ":planck_taper_window",
        "//intrinsic/util/testing:gtest_wrapper_main",
    ],
)

cc_library(
    name = "rms_monitor",
    srcs = ["rms_monitor.cc"],
    hdrs = ["rms_monitor.h"],
    deps = [
        "//intrinsic/icon/utils:realtime_status",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status:status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_test_and_malloc_test(
    name = "rms_monitor_test",
    srcs = ["rms_monitor_test.cc"],
    deps = [
        ":rms_monitor",
        "//intrinsic/icon/testing:malloc_test",
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_matchers",
        "@abseil-cpp//absl/status",
    ],
)

cc_library(
    name = "sliding_dft",
    hdrs = ["sliding_dft.h"],
    deps = [
        "//intrinsic/icon/utils:realtime_status",
        "//intrinsic/icon/utils:realtime_status_or",
        "//third_party/imported/cpp_libraries/math:mathutil",
        "@abseil-cpp//absl/memory",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_test_and_malloc_test(
    name = "sliding_dft_test",
    srcs = ["sliding_dft_test.cc"],
    deps = [
        ":sliding_dft",
        "//intrinsic/icon/testing:malloc_test",
        "//intrinsic/icon/utils:realtime_status_macro",
        "//intrinsic/icon/utils:realtime_status_matchers",
        "//intrinsic/math:complex_matchers",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_library(
    name = "sliding_window",
    hdrs = ["sliding_window.h"],
    deps = [
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/types:span",
    ],
)

cc_test_and_malloc_test(
    name = "sliding_window_test",
    srcs = ["sliding_window_test.cc"],
    deps = [
        ":sliding_window",
        "//intrinsic/icon/testing:malloc_test",
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/status:statusor",
        "@abseil-cpp//absl/strings",
        "@abseil-cpp//absl/types:span",
    ],
)
