# Copyright 2023 Intrinsic Innovation LLC

"""Tests for the downsampler pybind module."""

import datetime

from absl.testing import absltest

from intrinsic.logging.proto import downsampler_pb2
from intrinsic.logging.utils.python.downsampler.downsampler import DownsamplerEventSourceState
from intrinsic.logging.utils.python.downsampler.downsampler import DownsamplerOptions
from intrinsic.logging.utils.python.downsampler.downsampler import DownsamplerState

from pybind11_abseil import status  # isort: skip

DownsamplerStateProto = downsampler_pb2.DownsamplerState
DownsamplerOptionsProto = downsampler_pb2.DownsamplerOptions


class DownsamplerOptionsTest(absltest.TestCase):

  def test_roundtrip_works(self):
    options = DownsamplerOptions(
        sampling_interval_time=datetime.timedelta(seconds=1),
        sampling_interval_count=10,
    )
    options_proto = DownsamplerOptions.to_proto(options)
    self.assertEqual(DownsamplerOptions.from_proto(options_proto), options)

    nullopt_options = DownsamplerOptions(
        sampling_interval_time=None, sampling_interval_count=None
    )
    nullopt_options_proto = DownsamplerOptions.to_proto(nullopt_options)
    self.assertEqual(
        DownsamplerOptions.from_proto(nullopt_options_proto), nullopt_options
    )

  def test_downsampler_options_errors_on_invalid_sampling_interval_time(self):
    options = DownsamplerOptions(sampling_interval_time=datetime.timedelta.max)
    with self.assertRaisesRegex(
        (status.StatusNotOk, RuntimeError), "invalid sampling_interval_time"
    ):
      DownsamplerOptions.to_proto(options)

    options_proto = DownsamplerOptionsProto(
        sampling_interval_time=datetime.timedelta(seconds=315576000001)
    )
    with self.assertRaisesRegex(
        (status.StatusNotOk, RuntimeError), "invalid sampling_interval_time"
    ):
      DownsamplerOptions.from_proto(options_proto)


class DownsamplerEventSourceStateTest(absltest.TestCase):

  def test_downsampler_event_source_state_roundtrip_works(self):
    state = DownsamplerEventSourceState(
        last_use_time=datetime.datetime.fromtimestamp(
            123, tz=datetime.timezone.utc
        ),
        count_since_last_use=42,
    )
    state_proto = DownsamplerEventSourceState.to_proto(state)
    self.assertEqual(DownsamplerEventSourceState.from_proto(state_proto), state)

  def test_downsampler_event_source_state_errors_on_invalid_last_use_time(self):
    invalid_state = DownsamplerEventSourceState(
        last_use_time=datetime.datetime.max, count_since_last_use=42
    )
    with self.assertRaisesRegex(
        (status.StatusNotOk, RuntimeError), "invalid last_use_time"
    ):
      DownsamplerEventSourceState.to_proto(invalid_state)

    # from_proto is covered by datetime checks.
    # i.e., datetime.fromtimestamp(253402300800, tz=datetime.timezone.utc errors


class DownsamplerStateTest(absltest.TestCase):

  def test_downsampler_state_roundtrip_works(self):
    state = DownsamplerState(
        event_source_states={
            "event_source_1": DownsamplerEventSourceState(
                last_use_time=datetime.datetime.fromtimestamp(
                    123, tz=datetime.timezone.utc
                ),
                count_since_last_use=42,
            ),
            "event_source_2": DownsamplerEventSourceState(
                last_use_time=datetime.datetime.fromtimestamp(
                    456, tz=datetime.timezone.utc
                ),
                count_since_last_use=84,
            ),
        }
    )
    state_proto = DownsamplerState.to_proto(state)
    self.assertEqual(DownsamplerState.from_proto(state_proto), state)

    empty_state = DownsamplerState(event_source_states={})
    empty_state_proto = DownsamplerState.to_proto(empty_state)
    self.assertEqual(
        DownsamplerState.from_proto(empty_state_proto), empty_state
    )

  def test_downsampler_state_errors_on_invalid_last_use_time(self):
    invalid_state = DownsamplerState(
        event_source_states={
            "event_source_1": DownsamplerEventSourceState(
                last_use_time=datetime.datetime.max,
                count_since_last_use=42,
            ),
        }
    )
    with self.assertRaisesRegex(
        (status.StatusNotOk, RuntimeError), "invalid last_use_time"
    ):
      DownsamplerState.to_proto(invalid_state)

    invalid_state_proto = DownsamplerStateProto()
    invalid_state_proto.event_source_states[
        "event_source_1"
    ].last_use_time.seconds = 253402300800
    with self.assertRaisesRegex(
        (status.StatusNotOk, RuntimeError), "invalid last_use_time"
    ):
      DownsamplerState.from_proto(invalid_state_proto)


if __name__ == "__main__":
  absltest.main()
