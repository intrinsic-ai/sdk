// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/logging/utils/downsampler/downsampler.h"

#include <optional>
#include <string>

#include "absl/container/flat_hash_map.h"
#include "absl/time/time.h"
#include "intrinsic/logging/proto/downsampler.pb.h"
#include "intrinsic/logging/proto/logger_service.pb.h"
#include "intrinsic/logging/utils/downsampler/proto_conversion.h"
#include "pybind11/operators.h"
#include "pybind11/pybind11.h"
#include "pybind11/stl.h"
#include "pybind11_abseil/absl_casters.h"
#include "pybind11_abseil/status_casters.h"
#include "pybind11_protobuf/native_proto_caster.h"

using DownsamplerOptionsProto =
    ::intrinsic_proto::data_logger::DownsamplerOptions;
using DownsamplerEventSourceStateProto =
    ::intrinsic_proto::data_logger::DownsamplerEventSourceState;
using DownsamplerStateProto = ::intrinsic_proto::data_logger::DownsamplerState;
namespace py = pybind11;

namespace intrinsic::data_logger {

PYBIND11_MODULE(downsampler, m) {
  py::google::ImportStatusModule();
  pybind11_protobuf::ImportNativeProtoCasters();

  py::class_<DownsamplerOptions>(m, "DownsamplerOptions")
      .def(py::init<std::optional<absl::Duration>, std::optional<int>>(),
           py::arg("sampling_interval_time") = py::none(),
           py::arg("sampling_interval_count") = py::none())
      .def_static(
          "from_proto",
          py::overload_cast<const DownsamplerOptionsProto&>(
              &::intrinsic_proto::data_logger::FromProto),
          py::arg("options"),
          "Creates a DownsamplerOptions from a DownsamplingOptionsProto.")
      .def_static("to_proto",
                  py::overload_cast<const DownsamplerOptions&>(
                      &intrinsic::data_logger::ToProto),
                  py::arg("options"),
                  "Converts a DownsamplerOptions to a DownsamplerOptionsProto.")
      .def(py::self == py::self)
      .def_readwrite("sampling_interval_time",
                     &DownsamplerOptions::sampling_interval_time,
                     "Optional: Time-based sampling interval (absl::Duration).")
      .def_readwrite("sampling_interval_count",
                     &DownsamplerOptions::sampling_interval_count,
                     "Optional: Count-based sampling interval (int).");

  py::class_<DownsamplerEventSourceState>(m, "DownsamplerEventSourceState")
      .def(py::init<absl::Time, int>(), py::arg("last_use_time"),
           py::arg("count_since_last_use"))
      .def_static("from_proto",
                  py::overload_cast<const DownsamplerEventSourceStateProto&>(
                      &::intrinsic_proto::data_logger::FromProto),
                  py::arg("proto"),
                  "Creates a DownsamplerEventSourceState from a "
                  "DownsamplerEventSourceStateProto.")
      .def_static("to_proto",
                  py::overload_cast<const DownsamplerEventSourceState&>(
                      &intrinsic::data_logger::ToProto),
                  py::arg("state"),
                  "Converts a DownsamplerEventSourceState to a "
                  "DownsamplerEventSourceStateProto.")
      .def(py::self == py::self)
      .def_readwrite("last_use_time",
                     &DownsamplerEventSourceState::last_use_time)
      .def_readwrite("count_since_last_use",
                     &DownsamplerEventSourceState::count_since_last_use);

  py::class_<DownsamplerState>(m, "DownsamplerState")
      .def(py::init<
               absl::flat_hash_map<std::string, DownsamplerEventSourceState>>(),
           py::arg("event_source_states"))
      .def_static("from_proto",
                  py::overload_cast<const DownsamplerStateProto&>(
                      &::intrinsic_proto::data_logger::FromProto),
                  py::arg("proto"),
                  "Creates a DownsamplerState from a DownsamplerStateProto.")
      .def_static("to_proto",
                  py::overload_cast<const DownsamplerState&>(
                      &intrinsic::data_logger::ToProto),
                  py::arg("state"),
                  "Converts a DownsamplerState to a DownsamplerStateProto.")
      .def(py::self == py::self)
      .def_readwrite("event_source_states",
                     &DownsamplerState::event_source_states);

  py::class_<Downsampler>(m, "Downsampler")
      .def(py::init<DownsamplerOptions>(), py::arg("options"),
           "Constructs a Downsampler with DownsamplerOptions.")
      .def("should_downsample", &Downsampler::ShouldDownsample, py::arg("item"),
           "Determines whether a LogItem should be downsampled.")
      .def("register_ingest", &Downsampler::RegisterIngest, py::arg("item"),
           "Registers the ingestion of a LogItem.")
      .def("set_event_source_state", &Downsampler::SetEventSourceState,
           py::arg("event_source"), py::arg("state"),
           "Sets the state of the downsampler for the given event source.")
      .def("get_event_source_state", &Downsampler::GetEventSourceState,
           py::arg("event_source"),
           "Gets the state of the downsampler for the given event source.",
           py::return_value_policy::copy)
      .def("get_state", &Downsampler::GetState,
           "Gets the state of the downsampler.")
      .def("set_state", &Downsampler::SetState, py::arg("state"),
           "Sets the state of the downsampler.")
      .def("reset", &Downsampler::Reset,
           "Resets the downsampler to its initial state.");
}

}  // namespace intrinsic::data_logger
