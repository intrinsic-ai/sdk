# Copyright 2023 Intrinsic Innovation LLC

"""API to work with/visualize structured data.

This class provides access to LogItem protos logged by the Intrinsic Data
Logger.

Throughout this file, you'll see that the current time is obtained with:

  now = datetime.datetime.now(datetime.timezone.utc)

(See
https://docs.python.org/3/library/datetime.html#aware-and-naive-objects for
explanation about "aware" and "naive" in the following section.)

The fact the timezone is specified is important because this makes the created
datetime object "aware" of the timezone and ensures that the conversion to
google.protobuf.Timestamp is correct since FromDatetime assumes a UTC timezone
for "naive" datetime objects.
https://github.com/protocolbuffers/protobuf/blob/10307b5b1d7ca5cec4e7e18c1c12eb989c5435e9/python/google/protobuf/internal/well_known_types.py#L265

Note the fact that the timezone is UTC specifically is not important. As long as
the datetime object is "aware", the conversion to google.protobuf.Timestamp will
take it into account.
"""

from collections.abc import Iterable
import dataclasses
import datetime
import logging
import re
from typing import Any, Callable, Dict, Optional, Type, Union

from google.protobuf import empty_pb2
from google.protobuf import json_format
from google.protobuf import message as proto_message
import grpc
from intrinsic.icon.proto import part_status_pb2
from intrinsic.icon.proto import streaming_output_pb2
from intrinsic.icon.python import icon_logging
from intrinsic.logging.proto import bag_metadata_pb2
from intrinsic.logging.proto import log_item_pb2
from intrinsic.logging.proto import logger_service_pb2
from intrinsic.logging.proto import logger_service_pb2_grpc
from intrinsic.util.grpc import error_handling
import pandas as pd

# Used to transform arbitrary event source strings into valid Python names.
_REGEX_INVALID_PYTHON_VAR_CHARS = r'\W|^(?=\d)'


def _is_timezone_aware(dt: datetime.datetime) -> bool:
  """Checks whether the given datetime is timezone aware.

  Follows the definition here:
    https://docs.python.org/3/library/datetime.html#determining-if-an-object-is-aware-or-naive

  Args:
    dt: a datetime object

  Returns:
    True if dt is timezone "aware".
  """
  return dt.tzinfo is not None and dt.tzinfo.utcoffset(dt) is not None


def _interpret_as_utc(dt: datetime.datetime) -> datetime.datetime:
  """Forcibly interprets the given datetime as UTC.

  Note this is different than finding the equivalent time in UTC. Rather than
  adjusting the time to account for the clock offset, this function OVERRIDES
  the timezone with UTC.

  This is primarily useful for sticking a UTC timezone on to a datetime that
  doesn't have a timezone specified. Be very careful if attempting to use it in
  other scenarios.

  Args:
    dt: a datetime object

  Returns:
    A new datetime object with the same date and time values, and the timezone
    replaced with UTC.
  """
  return dt.combine(dt.date(), dt.time(), datetime.timezone.utc)


class DataSource:
  """Provides an API for interacting with LogItem payloads of a given type."""

  def __init__(self, log_items: list[log_item_pb2.LogItem]):
    self._log_items: list[log_item_pb2.LogItem] = log_items

  @property
  def log_items(self) -> list[log_item_pb2.LogItem]:
    """The list of log items backing this source."""
    return self._log_items

  @property
  def num_events(self) -> int:
    """How many log items we have of this source."""
    return len(self._log_items)

  def _get_data_frame(
      self,
      payload_accessor: Callable[
          [log_item_pb2.LogItem], log_item_pb2.LogItem.Payload
      ],
      every_n: int = 1,
  ) -> pd.DataFrame:
    """Returns a Pandas Dataframe with any payload proto.

    Args:
      payload_accessor: function to extract the relevant fields from a LogItem
        proto object.
      every_n: Sample rate, only every nth sample is returned.

    Returns:
      Pandas Dataframe with two columns; "payload" with protos matching
      `payload_accessor` and "time" as datetime, sorted by "time".
    """
    df = pd.DataFrame(columns=['payload', 'time'])

    payload = []
    times = []
    for log_item in self._log_items[::every_n]:
      payload.append(payload_accessor(log_item))
      times.append(log_item.metadata.acquisition_time.ToDatetime())
    df[df.columns[0]] = payload
    df[df.columns[1]] = times
    return df.sort_values(by=['time'])


def _get_part_status(log_item: log_item_pb2.LogItem, part_name: str):
  return log_item.payload.icon_robot_status.status_map[part_name]


class PartStatusSource(DataSource):
  """Data source for the status of a single ICON part.

  This class wraps a single ICON part a part status LogItem and allows it to
  access the fields as pandas.DataFrame to make plotting and data processing
  simple.
  """

  def __init__(self, log_items: list[log_item_pb2.LogItem], part_name: str):
    """Init the PartStatusSource.

    Args:
      log_items: A list of LogItems containing a part_status.
      part_name: The name of the part this PartStatusSource should wrap.
    """
    super().__init__(log_items)
    self._part_name = part_name

  def _get_used_fields(self) -> list[str]:
    """Returns all used fields in the part_status payload."""
    first_item = _get_part_status(self._log_items[0], self._part_name)
    return [
        field_descriptor.name for field_descriptor, _ in first_item.ListFields()
    ]

  def __dir__(self):
    return [
        'get_' + field for field in self._get_used_fields()
    ] + _list_public_methods(self)

  class _CallablePayloadMethod:
    """Helper class to create a callable get_<field_name> method."""

    def __init__(
        self,
        log_items: list[log_item_pb2.LogItem],
        payload_accessor: Callable[
            [part_status_pb2.PartStatus],
            Union[proto_message.Message, list[proto_message.Message]],
        ],
        part_name: str,
    ):
      """Init the _CallablePayloadMethod.

      Args:
        log_items: A list of log items.
        payload_accessor: A function to access the payload in PartStatus of
          `part_name`.
        part_name: The name of the part which should be accessed.
      """
      self._payload_accessor = payload_accessor
      self._log_items = log_items
      self._part_name = part_name

    def _get_data_frame(
        self,
        payload_accessor: Callable[
            [part_status_pb2.PartStatus],
            Union[proto_message.Message, list[proto_message.Message]],
        ],
        every_n: int = 1,
    ) -> pd.DataFrame:
      """Returns a Pandas Dataframe with any payload proto.

      Args:
        payload_accessor: A function to access the payload in PartStatus of
          self.part_name.
        every_n: Sample rate, only every nth sample is returned.

      Returns:
        Pandas Dataframe with all proto fields returned by `payload_accessor` as
        columns indexed by the timestamp_ns in seconds.
      """
      items = []
      for log_item in self._log_items[::every_n]:
        part_status = _get_part_status(log_item, self._part_name)
        item = json_format.MessageToDict(
            payload_accessor(part_status),
            always_print_fields_with_no_presence=True,
            preserving_proto_field_name=True,
        )
        item['time_s'] = part_status.timestamp_ns * 1e-9
        item['skill_log_id'] = log_item.context.skill_id
        icon_action_id = None
        if log_item.context.HasField('icon_action_id'):
          icon_action_id = log_item.context.icon_action_id
        item['icon_action_id'] = icon_action_id
        items.append(item)
      df = pd.json_normalize(items, sep='#')
      df.set_index('time_s', inplace=True)
      df.columns = df.columns.str.split('#', expand=True)
      return df

    def _get_repeated_data_frame(
        self,
        payload_accessor: Callable[
            [part_status_pb2.PartStatus],
            Union[proto_message.Message, list[proto_message.Message]],
        ],
        every_n: int = 1,
    ) -> pd.DataFrame:
      """Returns a Pandas Dataframe with any repeated payload proto.

      Args:
        payload_accessor: A function to access the payload in PartStatus of
          `part_name`.
        every_n: Sample rate, only every nth sample is returned.

      Returns:
        Pandas Dataframe with all proto fields returned by `payload_accessor` as
        columns indexed by the timestamp_ns in seconds.
      """

      def make_repeated_payload_accessor(
          index: int,
      ) -> Callable[[part_status_pb2.PartStatus], proto_message.Message]:
        """Create a payload accessor for a single item in a repeated field."""

        def repeated_payload_accessor(log_item: part_status_pb2.PartStatus):
          return payload_accessor(log_item)[index]

        return repeated_payload_accessor

      n_dof = len(payload_accessor(self._first_part_status()))
      df_list = []
      for i in range(n_dof):
        df_i = self._get_data_frame(make_repeated_payload_accessor(i), every_n)
        if i > 0:
          df_i.drop(columns=['skill_log_id', 'icon_action_id'], inplace=True)
        df_list.append(df_i.add_suffix('#' + str(i)))
      df = pd.concat(df_list, axis=1)
      df.columns = df.columns.str.split('#', expand=True)
      return df

    def __call__(self, *args, **kwargs):
      first_log_item = self._payload_accessor(self._first_part_status())
      every_n = kwargs.get('every_n') or 1
      if isinstance(first_log_item, Iterable):
        return self._get_repeated_data_frame(self._payload_accessor, every_n)
      return self._get_data_frame(self._payload_accessor, every_n)

    def _first_part_status(self):
      return _get_part_status(self._log_items[0], self._part_name)

  def __getattr__(
      self, name: str, *args, **kwargs
  ) -> Callable[[], pd.DataFrame]:
    prefix = 'get_'
    field = name
    if name.startswith(prefix):
      field = name[len(prefix) :]
    if field not in self._get_used_fields():
      raise AttributeError(field + ' is empty or no field of the part_status.')

    def payload_accessor(log_item: part_status_pb2.PartStatus):
      return getattr(log_item, field)

    return PartStatusSource._CallablePayloadMethod(
        self._log_items, payload_accessor, self._part_name
    )

  def has_field(self, field_name: str) -> bool:
    """Checks if the PartStatus usues the given field."""
    return field_name in self._get_used_fields()


class RobotStatusSource(DataSource):
  """Data source for the ICON robot status.

  This class wraps LogItems and has attributes for every part in the log item.
  """

  def _part_names(self):
    return self._log_items[0].payload.icon_robot_status.status_map.keys()

  def __getattr__(self, part_name: str) -> PartStatusSource:
    if part_name in self._part_names():
      return PartStatusSource(self._log_items, part_name)
    raise AttributeError(part_name + ' is not a valid part name.')

  def __getitem__(self, part_name: str) -> PartStatusSource:
    try:
      return self.__getattr__(part_name)
    except AttributeError as e:
      raise KeyError(str(e)) from None

  def __dir__(self):
    return self._part_names()

  def _get_single_part(
      self, part_selector: Callable[[part_status_pb2.PartStatus], bool]
  ) -> PartStatusSource:
    """Gets the logs for a single part.

    Only works if the logs contain exactly one part matching the part_selector.

    Args:
      part_selector: A function to select the part.

    Returns:
      The selected part.

    Raises:
      ValueError: If no or more than one part matching the part_selector were
      found.
    """
    robot_status_item = self.log_items[0].payload.icon_robot_status
    selected_parts = [
        part
        for part in self._part_names()
        if part_selector(robot_status_item.status_map[part])
    ]
    if len(selected_parts) != 1:
      raise ValueError(
          f'Found {len(selected_parts)} parts matching the selector, expected'
          ' one.'
      )
    return self.__getattr__(selected_parts[0])

  def get_single_arm_part(self) -> PartStatusSource:
    """Gets the logs for the arm part.

    Only works if the logs contain exactly one arm-part, which is usually the
    case.

    Returns:
      The arm part logs.

    Raises:
      ValueError: If no or more than one arm part were found.
    """
    return self._get_single_part(lambda part_status: part_status.joint_states)

  def get_single_ft_sensor_part(self) -> PartStatusSource:
    """Gets the logs for the force-torque sensor part.

    Only works if the logs contain exactly one force-torque sensor part, which
    is usually the case.

    Returns:
      The force-torque sensor part logs.

    Raises:
      ValueError: If no or more than one force-torque sensor part were found.
    """
    return self._get_single_part(
        lambda part_status: part_status.HasField('wrench_at_tip')
    )


class StreamingOutputSource(DataSource):
  """Data source for streamed action outputs."""

  def get_payload(
      self, class_to_unpack_to: Type[proto_message.Message], every_n: int = 1
  ) -> pd.DataFrame:
    """Returns the payload as pandas DataFrame."""
    proto_items, timestamps = icon_logging.unpack_streaming_outputs(
        self._log_items, class_to_unpack_to
    )
    items = []
    contexts = [log_item.context for log_item in self._log_items]
    for timestamp, log_item, context in list(
        zip(timestamps, proto_items, contexts)
    )[::every_n]:
      item = json_format.MessageToDict(
          log_item,
          always_print_fields_with_no_presence=True,
          preserving_proto_field_name=True,
      )
      item['time_s'] = timestamp
      item['skill_log_id'] = context.skill_id
      icon_action_id = None
      if context.HasField('icon_action_id'):
        icon_action_id = context.icon_action_id
      item['icon_action_id'] = icon_action_id

      items.append(item)
    df = pd.json_normalize(items, sep='#')
    df.set_index('time_s', inplace=True)
    df.columns = df.columns.str.split('#', expand=True)
    return df


def _data_source_factory(data: list[log_item_pb2.LogItem]) -> DataSource:
  """Create a DataSource given a list of LogItems.

  Depending on payload type, instantiate different classes with different
  abilities (draw plots, show images etc.).

  Args:
    data: The list of LogItems.

  Returns:
    Data source for the given data.
  """
  if not data:
    return DataSource([])
  # We assume that all log items with the same event source will have the same
  # payload type. That's why we only check data[0] below.
  # Other payloads.
  payload = data[0].payload
  if payload.HasField('icon_robot_status'):
    return RobotStatusSource(data)
  elif payload.HasField('any'):
    if payload.any.type_url.endswith(
        streaming_output_pb2.StreamingOutputWithMetadata.DESCRIPTOR.full_name
    ):
      return StreamingOutputSource(data)
  return DataSource(data)


def format_event_source(event_source: str) -> str:
  if not event_source:
    return ''
  # Replace characters which are not valid in Python names
  # This is to enable structured_logs.event_source for any event_source
  # since event sources can be arbitrary strings.
  formatted = re.sub(_REGEX_INVALID_PYTHON_VAR_CHARS, '_', event_source)
  # Avoid starting with an underscore.
  formatted = formatted.lstrip('_')
  return formatted


# Typically, the user should make sure to use "aware" datetime objects to avoid
# ambiguity. This can be done by simply including some timezone in the call, for
# example:
#
#   now = datetime.datetime.now(datetime.timezone.utc)
#
# See https://docs.python.org/3/library/datetime.html#aware-and-naive-objects
#
# If "aware" objects are not used, expect that the time will be interpreted as
# UTC.
@dataclasses.dataclass
class EventSourceWindow:
  start_time: datetime.datetime = datetime.datetime.fromtimestamp(0)
  end_time: datetime.datetime = datetime.datetime.now(datetime.timezone.utc)


class EventSourceReader:
  """Reader of a particular event source string."""

  def __init__(
      self, stub: logger_service_pb2_grpc.DataLoggerStub, event_source: str
  ):
    self._stub: logger_service_pb2_grpc.DataLoggerStub = stub
    self._event_source: str = event_source
    self._data: list[log_item_pb2.LogItem] = []
    self._cursor: str = None

  # The "*," makes the parameters keyword-only, which improves
  # Jupyter discoverability.
  @error_handling.retry_on_grpc_transient_errors
  def read(
      self,
      *,
      seconds_to_read: Optional[int] = None,
      time_window: Optional[EventSourceWindow] = None,
      sampling_period_ms: int = 0,
      max_num_items: int = 10000,
      filter_labels: Optional[dict[str, str]] = None,
  ) -> DataSource:
    """Read the last `seconds_to_read` of onprem logs for this event source.

    The time range to read the data can be either defined by the last seconds in
    the past using seconds_to_read or by the start and endtime using
    time_window.

    Args:
      seconds_to_read: How many seconds into the past we want to read. Use this
        or time_window.
      time_window: The start and end time of the data to read. Use this or
        seconds_to_read.
      sampling_period_ms: An optional downsampling parameter representing the
        minimum time in milliseconds between successive samples.
      max_num_items: The maximum number of returned items.
      filter_labels: Dictionary of label to value to filter the query.

    When specifying time_window, the user should typically make sure to use
    "aware" datetime objects to avoid ambiguity. This can be done by simply
    including some timezone in the call, for example:

      now = datetime.datetime.now(datetime.timezone.utc)

    See https://docs.python.org/3/library/datetime.html#aware-and-naive-objects

    If "aware" objects are not used, expect that the time will be interpreted as
    UTC.

    Returns:
      The DataSource for the read items.
    """
    if seconds_to_read is not None and time_window is not None:
      raise AttributeError('Only seconds_to_read or time_window can be used.')

    if seconds_to_read is not None:
      now = datetime.datetime.now(datetime.timezone.utc)
      used_time_window = EventSourceWindow(
          start_time=now - datetime.timedelta(seconds=seconds_to_read),
          end_time=now,
      )
    elif time_window is not None:
      used_time_window = time_window
    else:
      raise ValueError('seconds_to_read or time_window need to be defined.')

    return self._read_time_window(
        window=used_time_window,
        sampling_period_ms=sampling_period_ms,
        max_num_items=max_num_items,
        filter_labels=filter_labels,
    )

  def _read_time_window(
      self,
      *,
      window: EventSourceWindow,
      sampling_period_ms: int = 0,
      max_num_items: int = 10000,
      filter_labels: Optional[dict[str, str]] = None,
  ):
    """Read the onprem logs for a given time window for this event source.

    Args:
      window: The time window to read the logs.
      sampling_period_ms: An optional downsampling parameter representing the
        minimum time in milliseconds between successive samples.
      max_num_items: The maximum number of returned items.
      filter_labels: Dictionary of label to value to filter the query.

    Returns:
      The DataSource for the read items.
    """
    # If the datetimes in the window are naive (i.e. do not specify a timezone),
    # we assume they are utc, and make this explicit.
    tz_aware_window = window
    if not _is_timezone_aware(tz_aware_window.start_time):
      tz_aware_window.start_time = _interpret_as_utc(tz_aware_window.start_time)
    if not _is_timezone_aware(tz_aware_window.end_time):
      tz_aware_window.end_time = _interpret_as_utc(tz_aware_window.end_time)

    get_request = logger_service_pb2.GetLogItemsRequest()
    get_request.max_num_items = max_num_items

    if self._cursor:  # This is the cursor from the previous response
      get_request.cursor = self._cursor
    else:  # New query
      query_proto = get_request.get_query
      query_proto.event_source = self._event_source
      query_proto.start_time.FromDatetime(tz_aware_window.start_time)
      query_proto.end_time.FromDatetime(tz_aware_window.end_time)
      if sampling_period_ms > 0:
        query_proto.downsampler_options.sampling_interval_time.FromMilliseconds(
            sampling_period_ms
        )
      if filter_labels is not None:
        query_proto.filter_labels.update(filter_labels)

    response = self._stub.GetLogItems(get_request)
    for item in response.log_items:
      self._data.append(item)

    if len(self._data) >= max_num_items:
      logging.warning(
          'max_num_items exceedeed. Use a bigger value to get all logs.'
      )

    self._cursor = response.next_page_cursor

    # Delete old items that have fallen out of the current window.
    first_item_index = 0
    for item in self._data:
      if (
          item.metadata.acquisition_time.ToDatetime(datetime.timezone.utc)
          >= tz_aware_window.start_time
      ):
        # Since, http://cl/365776235, the log items are returned in order of
        # acquisition_time so break as soon as we have found the first item
        # that falls inside the requested window.
        break
      else:
        first_item_index += 1
    # Deletes the range [0, first_item_index).
    del self._data[:first_item_index]
    return _data_source_factory(self._data)

  @error_handling.retry_on_grpc_transient_errors
  def peek(self) -> log_item_pb2.LogItem:
    """Returns the most recent LogItem for this event source."""
    request = logger_service_pb2.GetMostRecentItemRequest(
        event_source=self._event_source
    )
    response = self._stub.GetMostRecentItem(request)
    return response.item


def _list_public_methods(instance: object) -> list[str]:
  """Returns all public methods of the given instance.

  Args:
    instance: Any class instance.
  """
  return [
      method for method in dir(instance.__class__) if not method.startswith('_')
  ]


class StructuredLogs:
  """Wrapper for interacting with structured logged data.

  This class handles reading of onprem logs. Class attributes, triggered through
  auto complete, give access to items from a given event source.
  Alternatively, you can use 'query' to directly query for the underlying
  protos.
  """

  class LogOptions:
    """Wrapper for LogOptions proto."""

    def __init__(self):
      self._log_options = logger_service_pb2.LogOptions()

    def set_event_source(
        self, event_source: str
    ) -> 'StructuredLogs.LogOptions':
      """Sets the event source for the log options.

      Args:
        event_source: The event source to set the log options for, must be a
          valid regex. Used to determine which event sources the options should
          be applied to.

      Returns:
        The log options object.
      """
      self._log_options.event_source = event_source
      return self

    def set_log_options_precedence_value(
        self, log_options_precedence_value: int
    ) -> 'StructuredLogs.LogOptions':
      """Sets the log options precedence value for the log options.

      If there are multiple log options that match a given event source,
      the log options with the highest precedence value is used by the logging
      service.

      If there are multiple log options with the same precedence value, then a
      random one amongst the equals is used.

      It is generally recommended to have higher precedence value for more
      specific regexes, and lower precedence value for more general regexes.

      The larger the value, higher in the precedence order.

      Args:
        log_options_precedence_value: The log options precedence value to set.

      Returns:
        The log options object.
      """
      self._log_options.log_options_precedence_value = (
          log_options_precedence_value
      )
      return self

    def set_sync_active(self, sync_active: bool) -> 'StructuredLogs.LogOptions':
      self._log_options.sync_active = sync_active
      return self

    def set_max_buffer_byte_size(
        self, max_buffer_byte_size: int
    ) -> 'StructuredLogs.LogOptions':
      self._log_options.max_buffer_byte_size = max_buffer_byte_size
      return self

    def set_token_bucket_options(
        self, refresh: int, burst: int
    ) -> 'StructuredLogs.LogOptions':
      param = logger_service_pb2.TokenBucketOptions(
          refresh=refresh, burst=burst
      )
      self._log_options.logging_budget.CopyFrom(param)
      return self

    def set_priority(self, priority: int) -> 'StructuredLogs.LogOptions':
      self._log_options.priority = priority
      return self

    def set_retain_on_disk(
        self, retain_on_disk: bool
    ) -> 'StructuredLogs.LogOptions':
      self._log_options.retain_on_disk = retain_on_disk
      return self

    def set_retain_buffer_on_disk(
        self, retain_buffer_on_disk: bool
    ) -> 'StructuredLogs.LogOptions':
      self._log_options.retain_buffer_on_disk = retain_buffer_on_disk
      return self

    @property
    def log_options(self) -> logger_service_pb2.LogOptions:
      return self._log_options

  def __init__(self, stub: logger_service_pb2_grpc.DataLoggerStub):
    self._stub: logger_service_pb2_grpc.DataLoggerStub = stub
    self._cached_event_sources: list[str] = []

  def __getattr__(self, event_source: str) -> EventSourceReader:
    return self.get_event_source(event_source)

  def get_event_source(self, event_source: str) -> EventSourceReader:
    """Gets an EventSourceReader for the given event source.

    Will first check the cached event sources, and if the event source is not
    found there, will query the logger service for the list of event sources and
    refresh the cache.

    Args:
      event_source: The event source to get.

    Returns:
      An EventSourceReader for the given event source.

    Raises:
      AttributeError: If the event source is not found.
    """
    # The user might have specified a differently formatted event source, even
    # though we need an exact match with the event source to build a reader so
    # we need to handle that case by checking against the list of "true" event
    # sources obtained from the logger service.
    #
    # The most common occurrence of this is from the fact that this method is
    # called from __getattr__, which cannot support some expressions of event
    # sources, such as "event_sources.foo".
    #
    # e.g.: "event_sources.foo" would become "event_sources_foo", breaking the
    # match.
    #
    # Since listing event sources is expensive, we first check against the
    # cached list of event sources.
    for source in self._cached_event_sources:
      # We do our comparison by re-creating the mangling.
      if format_event_source(source) == event_source:
        try:
          # The event source string to use for the reader must be an exact match
          # of the event source name in the logger service.
          event_source_reader = EventSourceReader(self._stub, source)
          # We peek to check if the event source exists, since the cache might
          # have gotten stale.
          #
          # We peek instead of listing the event sources because it is cheaper.
          event_source_reader.peek()
          return event_source_reader
        except grpc.RpcError as exc:
          logging.warning(
              'Failed to read from event source %s: %s, refreshing event source'
              ' cache',
              source,
              exc,
          )

    event_sources = self.get_event_sources()
    for source in event_sources:
      if format_event_source(source) == event_source:
        return EventSourceReader(self._stub, source)
    raise AttributeError(
        f'Event source "{event_source}" not found. Available sources'
        f' ["{event_sources}"]'
    )

  def __dir__(self) -> list[str]:
    return [
        format_event_source(source) for source in self.get_event_sources()
    ] + _list_public_methods(self)

  @classmethod
  def connect(cls, grpc_channel: grpc.Channel) -> 'StructuredLogs':
    """Connect to a running data logger service.

    To allow retrieving large LogItems (blobs) remove the max gRPC response size
    limit.

    Args:
      grpc_channel: Channel to the executive gRPC service.

    Returns:
      A newly created instance of the DataLogger wrapper class.
    """
    return cls(logger_service_pb2_grpc.DataLoggerStub(grpc_channel))
  @classmethod
  def for_solution(cls, solution: Any) -> 'StructuredLogs':
    """Connect to the data logger service of a running solution.

    Args:
      solution: The running solution.

    Returns:
      A newly created instance of the DataLogger wrapper class.
    """
    return cls.connect(solution.grpc_channel)

  @error_handling.retry_on_grpc_transient_errors
  def get_event_sources(self) -> list[str]:
    """Returns all event sources logged. Mainly useful for debugging."""
    self._cached_event_sources = list(
        self._stub.ListLogSources(empty_pb2.Empty()).event_sources
    )
    return self._cached_event_sources

  @error_handling.retry_on_grpc_transient_errors
  def set_log_options(
      self,
      log_options: Dict[str, LogOptions],
  ) -> None:
    """Configures log options for an event_source."""
    log_options_request = logger_service_pb2.SetLogOptionsRequest(
        log_options={e: i.log_options for e, i in log_options.items()}
    )
    self._stub.SetLogOptions(log_options_request)

  @error_handling.retry_on_grpc_transient_errors
  def get_log_options(
      self,
      *,
      event_source: Optional[str] = None,
      key: Optional[str] = None,
  ) -> LogOptions:
    """Returns the log options for an event source.

    Args:
      event_source: The event source to get the log options for.
      key: The key to get the log options for in log options map. This is used
        to identify the log options for setting/overwriting and fetching, and is
        not used for matching against event sources.

    Returns:
      The log options for the event source or key.

    Raises:
      ValueError: If neither event_source nor key is provided, or both are
      provided.
    """
    if event_source is None and key is None:
      raise ValueError('Either event_source or key must be provided.')
    if event_source is not None and key is not None:
      raise ValueError('Only one of event_source or key can be provided.')

    log_options_request: logger_service_pb2.GetLogOptionsRequest = (
        logger_service_pb2.GetLogOptionsRequest(
            event_source=event_source,
            key=key,
        )
    )

    ret = self._stub.GetLogOptions(log_options_request).log_options
    return (
        self.LogOptions()
        .set_event_source(ret.event_source)
        .set_log_options_precedence_value(ret.log_options_precedence_value)
        .set_sync_active(ret.sync_active)
        .set_max_buffer_byte_size(ret.max_buffer_byte_size)
        .set_token_bucket_options(
            ret.logging_budget.refresh,
            ret.logging_budget.burst,
        )
        .set_priority(ret.priority)
        .set_retain_on_disk(ret.retain_on_disk)
        .set_retain_buffer_on_disk(ret.retain_buffer_on_disk)
    )

  @error_handling.retry_on_grpc_transient_errors
  def query(
      self,
      event_source: str,
      seconds_to_read: int = 1200,
      max_num_items: int = 10000,
  ) -> list[log_item_pb2.LogItem]:
    """Queries the data logs.

    Args:
      event_source: The topic to read.
      seconds_to_read: Only considers recent logs within this timeframe
      max_num_items: Return at most this many items from the start of the time
        range.

    Returns:
      Log items from the given event source
    """
    now = datetime.datetime.now(datetime.timezone.utc)
    window_start = now - datetime.timedelta(seconds=seconds_to_read)
    get_request = logger_service_pb2.GetLogItemsRequest()
    get_request.max_num_items = max_num_items
    get_request.get_query.event_source = event_source
    get_request.get_query.start_time.FromDatetime(window_start)
    get_request.get_query.end_time.FromDatetime(now)
    get_response = self._stub.GetLogItems(get_request)
    if get_response.HasField('truncation_cause'):
      logging.warning(get_response.truncation_cause)
    return get_response.log_items

  @error_handling.retry_on_grpc_transient_errors
  def query_for_time_range(
      self,
      event_source: str,
      start_time: datetime.datetime,
      end_time: datetime.datetime,
      max_num_items: int = 10000,
  ) -> list[log_item_pb2.LogItem]:
    """Queries the data logs for a given time range.

    Args:
      event_source: The topic to read.
      start_time: Beginning of window to query data for.
      end_time: End of window to query data for.
      max_num_items: Return at most this many items from the start of the time
        range.

    Returns:
      Log items from the given event source within the specified time range.
    """
    get_request = logger_service_pb2.GetLogItemsRequest()
    get_request.max_num_items = max_num_items
    get_request.get_query.event_source = event_source
    get_request.get_query.start_time.FromDatetime(start_time)
    get_request.get_query.end_time.FromDatetime(end_time)
    get_response = self._stub.GetLogItems(get_request)
    if get_response.HasField('truncation_cause'):
      logging.warning(get_response.truncation_cause)
    return get_response.log_items

  @error_handling.retry_on_grpc_transient_errors
  def log(self, request: logger_service_pb2.LogRequest) -> None:
    """Logs a LogRequest to the cloud.

    Args:
      request: a fully populated log request.
    """

    self._stub.Log(request)

  @error_handling.retry_on_grpc_transient_errors
  def sync_and_rotate_logs(
      self, event_sources: Optional[list[str]] = None
  ) -> logger_service_pb2.SyncResponse:
    """Syncs remaining logs to GCS and rotates log files.

    If no event source is specified, all logs should be synced and rotated.

    Args:
      event_sources: event sources to sync, as a list of regex patterns.

    Returns:
      A SyncAndRotateLogsResponse instance representing the response.
    """
    sync_request = logger_service_pb2.SyncRequest()
    if event_sources is None:
      sync_request.sync_all = True
    else:
      sync_request.sync_all = False
      for event_source in event_sources:
        sync_request.event_sources.append(event_source)
    return self._stub.SyncAndRotateLogs(sync_request)

  @error_handling.retry_on_grpc_unavailable
  def create_local_recording(
      self,
      start_time: datetime.datetime,
      end_time: datetime.datetime,
      description: str,
      event_sources_to_record: list[str],
  ) -> bag_metadata_pb2.BagMetadata:
    """Starts a local recording in the logging service.

    This invokes CreateLocalRecording on the on-prem logging service, which
    creates a local recording that is tracked and robustly uploaded to the
    cloud.

    Args:
      start_time: Beginning of window to record data for.
      end_time: End of window to record data for.
      description: A human-readable description of the recording.
      event_sources_to_record: Event sources to include in the recording, as a
        list of regex patterns.

    Returns:
      A CreateLocalRecordingResponse instance representing the response from
      calling the on-prem logging service, which contains the metadata of the
      created recording.
    """
    create_request = logger_service_pb2.CreateLocalRecordingRequest(
        description=description,
    )
    create_request.start_time.FromDatetime(start_time)
    create_request.end_time.FromDatetime(end_time)
    for event_source in event_sources_to_record:
      create_request.event_sources_to_record.append(event_source)
    return self._stub.CreateLocalRecording(create_request).bag

  @error_handling.retry_on_grpc_unavailable
  def list_local_recordings(
      self,
      start_time: Optional[datetime.datetime],
      end_time: Optional[datetime.datetime],
      only_summary_metadata: bool,
      bag_ids: list[str],
  ) -> list[bag_metadata_pb2.BagMetadata]:
    """Calls ListLocalRecordings on the on-prem logging service.

    This lists local recordings from structured logging data that were
    previously created with CreateLocalRecording.

    Args:
      start_time: Beginning of window to record data for. If None, defaults to
        the beginning of time.
      end_time: End of window to record data for. If None, defaults to the
        current time.
      only_summary_metadata: Whether to only return summary metadata.
      bag_ids: Bag IDs to filter by, if empty, all bags in query range are
        returned.

    Returns:
      A ListLocalRecordingsgResponse instance representing the response.
    """
    list_request = logger_service_pb2.ListLocalRecordingsRequest(
        only_summary_metadata=only_summary_metadata,
    )
    if start_time is not None:
      list_request.start_time.FromDatetime(start_time)
    if end_time is not None:
      list_request.end_time.FromDatetime(end_time)
    for bag_id in bag_ids:
      list_request.bag_ids.append(bag_id)
    return list(self._stub.ListLocalRecordings(list_request).bags)
