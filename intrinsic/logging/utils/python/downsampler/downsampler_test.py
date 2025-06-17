# Copyright 2023 Intrinsic Innovation LLC

"""Tests for the downsampler pybind module."""

import datetime

from absl.testing import absltest
from intrinsic.logging.proto import log_item_pb2
from intrinsic.logging.utils.python.downsampler import downsampler as downsampler_py
from pybind11_abseil import status as absl_status


def create_log_item(
    event_source: str, acquisition_time: datetime.datetime
) -> log_item_pb2.LogItem:
  """Helper to create a LogItem for testing."""
  item = log_item_pb2.LogItem()
  item.metadata.event_source = event_source
  item.metadata.acquisition_time.FromDatetime(acquisition_time)
  return item


class TypeEqualityTest(absltest.TestCase):

  def test_downsampler_options(self):
    options_1 = downsampler_py.DownsamplerOptions(
        sampling_interval_time=datetime.timedelta(seconds=1),
        sampling_interval_count=10,
    )
    options_2 = downsampler_py.DownsamplerOptions(
        sampling_interval_time=datetime.timedelta(seconds=1),
        sampling_interval_count=10,
    )
    self.assertEqual(options_1, options_2)

  def test_downsampler_event_source_state(self):
    state_1 = downsampler_py.DownsamplerEventSourceState(
        last_use_time=datetime.datetime.fromtimestamp(
            100, tz=datetime.timezone.utc
        ),
        count_since_last_use=1,
    )
    state_2 = downsampler_py.DownsamplerEventSourceState(
        last_use_time=datetime.datetime.fromtimestamp(
            100, tz=datetime.timezone.utc
        ),
        count_since_last_use=1,
    )
    self.assertEqual(state_1, state_2)

  def test_downsampler_state(self):
    state_1 = downsampler_py.DownsamplerState(
        event_source_states={
            "test_source": downsampler_py.DownsamplerEventSourceState(
                last_use_time=datetime.datetime.fromtimestamp(
                    100, tz=datetime.timezone.utc
                ),
                count_since_last_use=1,
            )
        }
    )
    state_2 = downsampler_py.DownsamplerState(
        event_source_states={
            "test_source": downsampler_py.DownsamplerEventSourceState(
                last_use_time=datetime.datetime.fromtimestamp(
                    100, tz=datetime.timezone.utc
                ),
                count_since_last_use=1,
            )
        }
    )
    self.assertEqual(state_1, state_2)


class DownsamplerTest(absltest.TestCase):

  def test_first_item_should_not_downsample(self):
    downsampler = downsampler_py.Downsampler(
        downsampler_py.DownsamplerOptions(sampling_interval_count=2)
    )
    item = create_log_item(
        "test_source",
        datetime.datetime.fromtimestamp(0, tz=datetime.timezone.utc),
    )
    self.assertFalse(downsampler.should_downsample(item))
    self.assertFalse(downsampler.should_downsample(item))
    downsampler.register_ingest(item)
    self.assertTrue(downsampler.should_downsample(item))

  def test_no_options_should_not_downsample(self):
    downsampler = downsampler_py.Downsampler(
        downsampler_py.DownsamplerOptions()
    )
    item = create_log_item(
        "test_source",
        datetime.datetime.fromtimestamp(0, tz=datetime.timezone.utc),
    )
    self.assertFalse(downsampler.should_downsample(item))
    self.assertFalse(downsampler.should_downsample(item))
    downsampler.register_ingest(item)
    self.assertFalse(downsampler.should_downsample(item))

  def test_time_interval_downsampling_zero(self):
    downsampler = downsampler_py.Downsampler(
        downsampler_py.DownsamplerOptions(
            sampling_interval_time=datetime.timedelta(0)
        )
    )
    item = create_log_item(
        "test_source",
        datetime.datetime.fromtimestamp(0, tz=datetime.timezone.utc),
    )
    self.assertFalse(downsampler.should_downsample(item))
    self.assertFalse(downsampler.should_downsample(item))
    downsampler.register_ingest(item)
    self.assertFalse(downsampler.should_downsample(item))

  def test_time_interval_downsampling_works(self):
    downsampler = downsampler_py.Downsampler(
        downsampler_py.DownsamplerOptions(
            sampling_interval_time=datetime.timedelta(seconds=2)
        )
    )
    item = create_log_item(
        "test_source",
        datetime.datetime.fromtimestamp(0, tz=datetime.timezone.utc),
    )

    # First message should go through.
    self.assertFalse(downsampler.should_downsample(item))

    # It should continue to be passed until a use is registered.
    self.assertFalse(downsampler.should_downsample(item))
    downsampler.register_ingest(item)
    self.assertTrue(downsampler.should_downsample(item))

    # Within interval.
    item = create_log_item(
        "test_source",
        datetime.datetime.fromtimestamp(1, tz=datetime.timezone.utc),
    )
    self.assertTrue(downsampler.should_downsample(item))

    item = create_log_item(
        "test_source",
        datetime.datetime.fromtimestamp(2, tz=datetime.timezone.utc),
    )
    self.assertFalse(downsampler.should_downsample(item))
    downsampler.register_ingest(item)
    self.assertTrue(downsampler.should_downsample(item))

    # Outside interval.
    item = create_log_item(
        "test_source",
        datetime.datetime.fromtimestamp(10, tz=datetime.timezone.utc),
    )
    self.assertFalse(downsampler.should_downsample(item))

    # Again, it should continue to be passed until a use is registered.
    self.assertFalse(downsampler.should_downsample(item))
    downsampler.register_ingest(item)
    self.assertTrue(downsampler.should_downsample(item))

  def test_count_interval_downsampling_zero(self):
    downsampler = downsampler_py.Downsampler(
        downsampler_py.DownsamplerOptions(sampling_interval_count=0)
    )
    item = create_log_item(
        "test_source",
        datetime.datetime.fromtimestamp(0, tz=datetime.timezone.utc),
    )
    self.assertFalse(downsampler.should_downsample(item))
    self.assertFalse(downsampler.should_downsample(item))
    downsampler.register_ingest(item)
    self.assertFalse(downsampler.should_downsample(item))

  def test_count_interval_downsampling_works(self):
    downsampler = downsampler_py.Downsampler(
        downsampler_py.DownsamplerOptions(sampling_interval_count=3)
    )
    item = create_log_item(
        "test_source",
        datetime.datetime.fromtimestamp(0, tz=datetime.timezone.utc),
    )
    self.assertFalse(downsampler.should_downsample(item))  # 1
    self.assertFalse(downsampler.should_downsample(item))  # 2
    downsampler.register_ingest(item)
    self.assertTrue(downsampler.should_downsample(item))  # 1
    self.assertTrue(downsampler.should_downsample(item))  # 2
    self.assertFalse(downsampler.should_downsample(item))  # 3

  def test_combined_downsampling(self):
    downsampler = downsampler_py.Downsampler(
        downsampler_py.DownsamplerOptions(
            sampling_interval_time=datetime.timedelta(seconds=2),
            sampling_interval_count=5,
        )
    )
    item = create_log_item(
        "test_source",
        datetime.datetime.fromtimestamp(0, tz=datetime.timezone.utc),
    )

    # First message should go through.
    self.assertFalse(downsampler.should_downsample(item))

    # It should continue to be passed until a use is registered.
    self.assertFalse(downsampler.should_downsample(item))
    downsampler.register_ingest(item)
    self.assertTrue(downsampler.should_downsample(item))  # 1

    # Within time interval, but out of count interval.
    item = create_log_item(
        "test_source",
        datetime.datetime.fromtimestamp(1, tz=datetime.timezone.utc),
    )
    self.assertTrue(downsampler.should_downsample(item))  # 2
    self.assertTrue(downsampler.should_downsample(item))  # 3
    self.assertTrue(downsampler.should_downsample(item))  # 4
    self.assertTrue(downsampler.should_downsample(item))  # 5
    downsampler.register_ingest(item)

    # Out of time interval, but not count interval.
    item = create_log_item(
        "test_source",
        datetime.datetime.fromtimestamp(5, tz=datetime.timezone.utc),
    )
    self.assertTrue(downsampler.should_downsample(item))  # 1
    self.assertTrue(downsampler.should_downsample(item))  # 2
    self.assertTrue(downsampler.should_downsample(item))  # 3
    self.assertTrue(downsampler.should_downsample(item))  # 4
    self.assertFalse(downsampler.should_downsample(item))  # 5

    # Again, it should continue to be passed until a use is registered.
    self.assertFalse(downsampler.should_downsample(item))
    downsampler.register_ingest(item)
    self.assertTrue(downsampler.should_downsample(item))

  def test_tracks_different_event_sources_separately(self):
    downsampler = downsampler_py.Downsampler(
        downsampler_py.DownsamplerOptions(sampling_interval_count=3)
    )

    a = create_log_item(
        "test_source_a",
        datetime.datetime.fromtimestamp(0, tz=datetime.timezone.utc),
    )
    b = create_log_item(
        "test_source_b",
        datetime.datetime.fromtimestamp(0, tz=datetime.timezone.utc),
    )

    self.assertFalse(downsampler.should_downsample(a))  # A: 1
    downsampler.register_ingest(a)  # A: 0
    self.assertTrue(downsampler.should_downsample(a))  # A: 1

    self.assertFalse(downsampler.should_downsample(b))  # B: 1
    downsampler.register_ingest(b)  # B: 0
    self.assertTrue(downsampler.should_downsample(b))  # B: 1

    self.assertTrue(downsampler.should_downsample(a))  # A: 2
    self.assertFalse(downsampler.should_downsample(a))  # A: 3

    self.assertTrue(downsampler.should_downsample(b))  # B: 2
    self.assertFalse(downsampler.should_downsample(b))  # B: 3

    downsampler.register_ingest(a)  # A: 0
    self.assertTrue(downsampler.should_downsample(a))  # A: 1
    self.assertFalse(downsampler.should_downsample(b))  # B: 4

  def test_get_and_set_event_source_state_works(self):
    downsampler = downsampler_py.Downsampler(
        downsampler_py.DownsamplerOptions(sampling_interval_count=3)
    )
    item = create_log_item(
        "test_source",
        datetime.datetime.fromtimestamp(0, tz=datetime.timezone.utc),
    )

    self.assertFalse(downsampler.should_downsample(item))  # 1
    downsampler.register_ingest(item)
    self.assertTrue(downsampler.should_downsample(item))  # 1

    # Get ok.
    state = downsampler.get_event_source_state("test_source")
    self.assertEqual(state.count_since_last_use, 1)
    self.assertEqual(
        state.last_use_time,
        datetime.datetime.fromtimestamp(0, tz=datetime.timezone.utc),
    )

    # Get non-existent.
    with self.assertRaises((absl_status.StatusNotOk, RuntimeError)):
      downsampler.get_event_source_state("non_existent_source")

    # Tracks use.
    item = create_log_item(
        "test_source",
        datetime.datetime.fromtimestamp(1, tz=datetime.timezone.utc),
    )
    downsampler.register_ingest(item)
    state = downsampler.get_event_source_state("test_source")
    self.assertEqual(state.count_since_last_use, 0)
    self.assertEqual(
        state.last_use_time,
        datetime.datetime.fromtimestamp(1, tz=datetime.timezone.utc),
    )

    # Restores across instances.
    downsampler_2 = downsampler_py.Downsampler(
        downsampler_py.DownsamplerOptions(sampling_interval_count=3)
    )
    downsampler_2.set_event_source_state("test_source", state)
    self.assertTrue(downsampler_2.should_downsample(item))  # 1
    self.assertTrue(downsampler_2.should_downsample(item))  # 2
    self.assertFalse(downsampler_2.should_downsample(item))  # 3

    # Also works if downsampler noops.
    downsampler_3 = downsampler_py.Downsampler(
        downsampler_py.DownsamplerOptions()
    )
    self.assertFalse(downsampler_3.should_downsample(item))  # 1
    with self.assertRaises((absl_status.StatusNotOk, RuntimeError)):
      downsampler_3.get_event_source_state("test_source")
    downsampler_3.set_event_source_state("test_source", state)
    self.assertFalse(downsampler_3.should_downsample(item))  # 1
    state = downsampler_3.get_event_source_state("test_source")
    self.assertEqual(state.count_since_last_use, 1)
    self.assertEqual(
        state.last_use_time,
        datetime.datetime.fromtimestamp(1, tz=datetime.timezone.utc),
    )

  def test_get_and_set_state_works(self):
    downsampler = downsampler_py.Downsampler(
        downsampler_py.DownsamplerOptions(sampling_interval_count=2)
    )
    state = downsampler.get_state()
    self.assertEmpty(state.event_source_states)

    item_a = create_log_item(
        "source_a",
        datetime.datetime.fromtimestamp(10, tz=datetime.timezone.utc),
    )
    item_b = create_log_item(
        "source_b",
        datetime.datetime.fromtimestamp(20, tz=datetime.timezone.utc),
    )

    # Doesn't track unseen event sources.
    self.assertFalse(downsampler.should_downsample(item_a))
    state = downsampler.get_state()
    self.assertEmpty(state.event_source_states)

    # Tracks registrations.
    downsampler.register_ingest(item_a)
    state = downsampler.get_state()
    self.assertDictEqual(
        state.event_source_states,
        {
            "source_a": downsampler_py.DownsamplerEventSourceState(
                last_use_time=datetime.datetime.fromtimestamp(
                    10, tz=datetime.timezone.utc
                ),
                count_since_last_use=0,
            )
        },
    )

    self.assertFalse(downsampler.should_downsample(item_b))
    downsampler.register_ingest(item_b)
    self.assertTrue(downsampler.should_downsample(item_b))
    state = downsampler.get_state()
    self.assertDictEqual(
        state.event_source_states,
        {
            "source_a": downsampler_py.DownsamplerEventSourceState(
                last_use_time=datetime.datetime.fromtimestamp(
                    10, tz=datetime.timezone.utc
                ),
                count_since_last_use=0,
            ),
            "source_b": downsampler_py.DownsamplerEventSourceState(
                last_use_time=datetime.datetime.fromtimestamp(
                    20, tz=datetime.timezone.utc
                ),
                count_since_last_use=1,
            ),
        },
    )

    # Restores across instances.
    downsampler_2 = downsampler_py.Downsampler(
        downsampler_py.DownsamplerOptions(sampling_interval_count=2)
    )
    downsampler_2.set_state(state)  # Sets B: 1

    item_a_2 = create_log_item(
        "source_a",
        datetime.datetime.fromtimestamp(11, tz=datetime.timezone.utc),
    )
    item_b_2 = create_log_item(
        "source_b",
        datetime.datetime.fromtimestamp(21, tz=datetime.timezone.utc),
    )

    self.assertTrue(downsampler_2.should_downsample(item_a_2))  # A: 1
    self.assertFalse(downsampler_2.should_downsample(item_a_2))  # A: 2
    self.assertFalse(downsampler_2.should_downsample(item_b_2))  # B: 2

    # Also works if downsampler noops.
    downsampler_3 = downsampler_py.Downsampler(
        downsampler_py.DownsamplerOptions()
    )
    downsampler_3.set_state(state)
    self.assertFalse(downsampler_3.should_downsample(item_a_2))  # A: 1
    self.assertFalse(downsampler_3.should_downsample(item_a_2))  # A: 2

  def test_reset(self):
    downsampler = downsampler_py.Downsampler(
        downsampler_py.DownsamplerOptions(
            sampling_interval_time=datetime.timedelta(seconds=2)
        )
    )
    a = create_log_item(
        "test_source_a",
        datetime.datetime.fromtimestamp(0, tz=datetime.timezone.utc),
    )
    b = create_log_item(
        "test_source_b",
        datetime.datetime.fromtimestamp(0, tz=datetime.timezone.utc),
    )

    self.assertFalse(downsampler.should_downsample(a))
    downsampler.register_ingest(a)
    self.assertTrue(downsampler.should_downsample(a))

    self.assertFalse(downsampler.should_downsample(b))
    downsampler.register_ingest(b)
    self.assertTrue(downsampler.should_downsample(b))

    downsampler.reset()

    self.assertFalse(downsampler.should_downsample(a))
    self.assertFalse(downsampler.should_downsample(b))


if __name__ == "__main__":
  absltest.main()
