# Copyright 2026 Intrinsic Innovation LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""In-memory fakes for the PubSub key-value store.

These fakes mirror `intrinsic.platform.pubsub.python.pubsub.PubSub` and
`KeyValueStore` (including `google.protobuf.Any` packing, raw binary `bytes`
storage with cross-encoding checks, Zenoh key/timeout validation, store-prefixed
`list_all_keys` results, `StatusNotOk` error codes, and Zenoh `*`/`**`/`$*`
wildcard deletions) so clients and services can be tested without a live Zenoh
session.
"""

import math
import random
import threading
from typing import TypeVar

from google.protobuf import any_pb2
from google.protobuf import message
from pybind11_abseil import status

_DEFAULT_STORE = "kv_store"
_REPLICATION_STORE = "kv_store_repl"
_T = TypeVar("_T")


def _validate_zenoh_keyexpr(keyexpr: str) -> None:
  """Validates `keyexpr` against `ValidZenohKeyexpr` / `ValidZenohKey` rules."""
  if not keyexpr:
    raise status.BuildStatusNotOk(
        status.StatusCode.INVALID_ARGUMENT, "Keyexpr must not be empty"
    )
  if keyexpr.startswith("/"):
    raise status.BuildStatusNotOk(
        status.StatusCode.INVALID_ARGUMENT, "Keyexpr must not start with /"
    )
  if keyexpr.endswith("/"):
    raise status.BuildStatusNotOk(
        status.StatusCode.INVALID_ARGUMENT, "Keyexpr must not end with /"
    )
  if any(not part for part in keyexpr.split("/")):
    raise status.BuildStatusNotOk(
        status.StatusCode.INVALID_ARGUMENT,
        "Keyexpr must not contain empty parts",
    )


def _validate_non_negative_int_timeout(timeout: int) -> None:
  """Validates an integer `timeout` (matching pybind11 `int` bindings)."""
  if isinstance(timeout, bool) or not isinstance(timeout, int):
    raise TypeError(f"Expected int timeout, got {type(timeout)}.")
  if timeout < 0:
    raise status.BuildStatusNotOk(
        status.StatusCode.INVALID_ARGUMENT, "Timeout must be zero or positive"
    )


def _validate_non_negative_float_timeout(timeout: float) -> None:
  """Validates a `float` `timeout` (matching pybind11 `double` bindings)."""
  if isinstance(timeout, bool) or not isinstance(timeout, (int, float)):
    raise TypeError(f"Expected float timeout, got {type(timeout)}.")
  if not math.isfinite(timeout) or timeout < 0:
    raise status.BuildStatusNotOk(
        status.StatusCode.INVALID_ARGUMENT, "Timeout must be zero or positive"
    )


def _validate_verification_options(options: object) -> None:
  """Validates `options.timeout` when a verification options object is given."""
  timeout = getattr(options, "timeout", None)
  if timeout is not None:
    if isinstance(timeout, bool) or not isinstance(timeout, (int, float)):
      raise TypeError(f"Expected numeric options.timeout, got {type(timeout)}.")
    if not math.isfinite(timeout) or timeout <= 0:
      raise status.BuildStatusNotOk(
          status.StatusCode.INVALID_ARGUMENT,
          "SetWithVerificationOptions::timeout must be positive",
      )


def _matches_segment(pattern: str, segment: str) -> bool:
  """Returns whether `segment` matches a single Zenoh keyexpr segment `pattern`.

  `*` and `$*` match any complete segment, while `$*` inside a segment (such as
  `frame_$*2`) matches any zero-or-more-character substring within that segment.

  Args:
    pattern: Single `/`-delimited segment from a Zenoh key expression.
    segment: Single `/`-delimited segment from a stored key.

  Returns:
    `True` if `segment` matches `pattern`.
  """
  if pattern in ("*", "$*"):
    return True
  if "$*" not in pattern:
    return pattern == segment
  parts = pattern.split("$*")
  prefix = parts[0]
  suffix = parts[-1]
  if (
      len(segment) < len(prefix) + len(suffix)
      or not segment.startswith(prefix)
      or not segment.endswith(suffix)
  ):
    return False
  pos = len(prefix)
  end = len(segment) - len(suffix)
  for part in parts[1:-1]:
    if not part:
      continue
    pos = segment.find(part, pos, end)
    if pos == -1:
      return False
    pos += len(part)
  return True


def _matches_segments(
    pattern_segments: list[str], key_segments: list[str]
) -> bool:
  """Returns whether `key_segments` matches `pattern_segments` with `**` globs.

  Each `**` segment matches zero or more complete key segments, while all other
  segments are matched individually via `_matches_segment`. Tracks the most
  recent `**` position iteratively in `O(len(pattern_segments) *
  len(key_segments))` worst-case time without recursion.

  Args:
    pattern_segments: `/`-split segments of the Zenoh key expression.
    key_segments: `/`-split segments of the stored key.

  Returns:
    `True` if `key_segments` matches `pattern_segments`.
  """
  pattern_idx = 0
  key_idx = 0
  wildcard_pattern_idx = -1
  wildcard_key_idx = 0
  while key_idx < len(key_segments):
    if (
        pattern_idx < len(pattern_segments)
        and pattern_segments[pattern_idx] == "**"
    ):
      wildcard_pattern_idx = pattern_idx
      wildcard_key_idx = key_idx
      pattern_idx += 1
    elif pattern_idx < len(pattern_segments) and _matches_segment(
        pattern_segments[pattern_idx], key_segments[key_idx]
    ):
      pattern_idx += 1
      key_idx += 1
    elif wildcard_pattern_idx != -1:
      wildcard_key_idx += 1
      key_idx = wildcard_key_idx
      pattern_idx = wildcard_pattern_idx + 1
    else:
      return False
  while (
      pattern_idx < len(pattern_segments)
      and pattern_segments[pattern_idx] == "**"
  ):
    pattern_idx += 1
  return pattern_idx == len(pattern_segments)


class FakeKeyValueStore:
  """Thread-safe dict-backed stand-in for `pubsub.KeyValueStore`.

  Protobuf values are stored packed in `google.protobuf.Any`, exactly like the
  real store, so readers unpack them with `get(key).Unpack(message)`. Raw binary
  payloads written via `set_raw` or `set_raw_with_verification` share the same
  keyspace and are read via `get_raw`, matching C++ `KeyValueStore`
  cross-encoding checks.
  """

  def __init__(self, prefix_override: str | None = None):
    self._key_prefix = (
        _DEFAULT_STORE if prefix_override is None else prefix_override
    )
    self._lock = threading.Lock()
    self._entries: dict[str, any_pb2.Any | bytes] = {}

  def _validate_prefixed_key(self, key: str) -> None:
    """Validates `key` and its store-prefixed Zenoh path."""
    _validate_zenoh_keyexpr(key)
    _validate_zenoh_keyexpr(f"{self._key_prefix}/{key}")

  def set(
      self,
      key: str,
      value: message.Message,
      high_consistency: bool | None = None,
  ) -> None:
    """Stores `value` at `key`, packing non-`Any` messages into an `Any`.

    Args:
      key: Key to write.
      value: Protobuf message to store. `Any` values are copied as-is.
      high_consistency: Accepted for signature compatibility; ignored.

    Raises:
      status.StatusNotOk: With `StatusCode.INVALID_ARGUMENT` if `key` or the
        configured store prefix is not a valid Zenoh key expression.
    """
    del high_consistency  # Unused.
    self._validate_prefixed_key(key)
    packed = any_pb2.Any()
    if isinstance(value, any_pb2.Any):
      packed.CopyFrom(value)
    else:
      packed.Pack(value)
    with self._lock:
      self._entries[key] = packed

  def set_with_verification(
      self,
      key: str,
      value: message.Message,
      options: object = None,
  ) -> None:
    """Stores `value` at `key`, mirroring `KeyValueStore.set_with_verification`.

    Args:
      key: Key to write.
      value: Protobuf message to store. `Any` values are copied as-is.
      options: Optional `SetWithVerificationOptions`-compatible object; if it
        exposes a `.timeout` attribute, `timeout` must be positive.

    Raises:
      TypeError: If `options.timeout` is not a numeric value.
      status.StatusNotOk: With `StatusCode.INVALID_ARGUMENT` if `key`, the store
        prefix, or `options.timeout` is invalid.
    """
    _validate_verification_options(options)
    self.set(key, value)

  def _set_raw_bytes(self, key: str, raw: bytes) -> None:
    """Validates `key` and stores normalized immutable `raw` bytes."""
    self._validate_prefixed_key(key)
    with self._lock:
      self._entries[key] = raw

  def set_raw(
      self,
      key: str,
      value: bytes | bytearray | memoryview,
  ) -> None:
    """Stores raw binary `bytes` at `key`, mirroring `KeyValueStore.set_raw`.

    Args:
      key: Key to write.
      value: Buffer-protocol payload (`bytes`, `bytearray`, `memoryview`, etc.).

    Raises:
      TypeError: If `value` does not support the Python buffer protocol.
      status.StatusNotOk: With `StatusCode.INVALID_ARGUMENT` if `key` or the
        configured store prefix is not a valid Zenoh key expression.
    """
    if type(value) is bytes:  # pylint: disable=unidiomatic-typecheck
      raw = value
    else:
      with memoryview(value) as view:
        raw = bytes(view)
    self._set_raw_bytes(key, raw)

  def set_raw_with_verification(
      self,
      key: str,
      value: bytes | bytearray | memoryview,
      options: object = None,
  ) -> None:
    """Stores raw binary `bytes` at `key`, mirroring `set_raw_with_verification`.

    Args:
      key: Key to write.
      value: Buffer-protocol payload (`bytes`, `bytearray`, `memoryview`, etc.).
      options: Optional `SetWithVerificationOptions`-compatible object; if it
        exposes a `.timeout` attribute, `timeout` must be positive.

    Raises:
      TypeError: If `value` does not support the Python buffer protocol or
        `options.timeout` is not a numeric value.
      status.StatusNotOk: With `StatusCode.INVALID_ARGUMENT` if `key`, the store
        prefix, or `options.timeout` is invalid.
    """
    if type(value) is bytes:  # pylint: disable=unidiomatic-typecheck
      _validate_verification_options(options)
      raw = value
    else:
      with memoryview(value) as view:
        _validate_verification_options(options)
        raw = bytes(view)
    self._set_raw_bytes(key, raw)

  def _get_typed(
      self, key: str, expected_type: type[_T], mismatch_suffix: str
  ) -> _T:
    """Validates `key`, fetches its entry, and checks its stored encoding."""
    self._validate_prefixed_key(key)
    with self._lock:
      entry = self._entries.get(key)
    if entry is None:
      raise status.BuildStatusNotOk(
          status.StatusCode.NOT_FOUND,
          f"Key '{key}' not found in fake KV store.",
      )
    if not isinstance(entry, expected_type):
      raise status.BuildStatusNotOk(
          status.StatusCode.INVALID_ARGUMENT,
          f"Key '{self._key_prefix}/{key}' {mismatch_suffix}",
      )
    return entry

  def get(self, key: str, timeout: int = 10) -> any_pb2.Any:
    """Returns a copy of the `Any` stored at `key`.

    Args:
      key: Key to read.
      timeout: Integer query timeout in seconds (must be zero or positive).

    Returns:
      The stored message packed in an `Any`.

    Raises:
      TypeError: If `timeout` is not an `int`.
      status.StatusNotOk: With `StatusCode.INVALID_ARGUMENT` if `key`, the store
        prefix, or `timeout` is invalid or `key` stores raw bytes, or
        `StatusCode.NOT_FOUND` if `key` was never set.
    """
    _validate_non_negative_int_timeout(timeout)
    entry = self._get_typed(
        key, any_pb2.Any, "stores raw bytes, not a Protobuf message"
    )
    result = any_pb2.Any()
    result.CopyFrom(entry)
    return result

  def get_raw(self, key: str, timeout: float = 10.0) -> bytes:
    """Returns the raw `bytes` payload stored at `key`.

    Args:
      key: Key to read.
      timeout: Query timeout in seconds (must be finite and zero or positive).

    Returns:
      The stored raw bytes.

    Raises:
      TypeError: If `timeout` is not an `int` or `float`.
      status.StatusNotOk: With `StatusCode.INVALID_ARGUMENT` if `key`, the store
        prefix, or `timeout` is invalid or `key` stores a Protobuf message, or
        `StatusCode.NOT_FOUND` if `key` was never set.
    """
    _validate_non_negative_float_timeout(timeout)
    return self._get_typed(
        key, bytes, "stores a Protobuf message, not raw bytes"
    )

  def delete(self, key: str) -> None:
    """Removes `key` (or all keys matching a Zenoh `*`/`**`/`$*` keyexpr)."""
    self._validate_prefixed_key(key)
    with self._lock:
      if "*" not in key:
        self._entries.pop(key, None)
        return
      pattern_segments = key.split("/")
      matching = [
          stored_key
          for stored_key in self._entries
          if _matches_segments(pattern_segments, stored_key.split("/"))
      ]
      for stored_key in matching:
        self._entries.pop(stored_key, None)

  def list_all_keys(self, timeout: int = 10) -> list[str]:
    """Returns all store-prefixed keys (`<store>/<key>`) in unspecified order.

    Args:
      timeout: Query timeout in seconds (must be zero or positive).

    Returns:
      Store-prefixed keys matching `<prefix>/**`, in randomized order.

    Raises:
      TypeError: If `timeout` is not an `int`.
      status.StatusNotOk: With `StatusCode.UNIMPLEMENTED` on replicated stores
        (`kv_store_repl`), or `StatusCode.INVALID_ARGUMENT` if `timeout < 0` or
        the store prefix is not a valid Zenoh key segment.
    """
    _validate_non_negative_int_timeout(timeout)
    if self._key_prefix == _REPLICATION_STORE:
      raise status.BuildStatusNotOk(
          status.StatusCode.UNIMPLEMENTED,
          "ListAllKeys is not supported for replicated KV store; use "
          "ListAllOnpremKeys or ListAllGlobalKeys instead.",
      )
    _validate_zenoh_keyexpr(f"{self._key_prefix}/**")
    with self._lock:
      keys = [f"{self._key_prefix}/{key}" for key in self._entries]
    random.shuffle(keys)
    return keys


class FakePubSub:
  """Thread-safe stand-in for `pubsub.PubSub` handing out `FakeKeyValueStore`s.

  Matches `pubsub.PubSub.key_value_store`: `prefix_override=None` selects the
  default `"kv_store"` namespace, whereas an explicit string selects that exact
  store prefix (and an empty `""` prefix fails Zenoh key validation on store
  operations just like live Zenoh).
  """

  def __init__(self, participant_name: str = "", config: str = ""):
    del participant_name, config  # Unused.
    self._lock = threading.Lock()
    self._stores: dict[str, FakeKeyValueStore] = {}

  @staticmethod
  def destroy_session_when_unused() -> None:
    """No-op stand-in for `pubsub.PubSub.destroy_session_when_unused`."""

  def key_value_store(
      self, prefix_override: str | None = None
  ) -> FakeKeyValueStore:
    """Returns the (lazily created) fake store for `prefix_override`.

    Args:
      prefix_override: Store name. `None` maps to the default `"kv_store"`
        namespace; explicit strings map to their own store prefix.

    Returns:
      The `FakeKeyValueStore` shared by all callers using the same store prefix.
    """
    name = _DEFAULT_STORE if prefix_override is None else prefix_override
    with self._lock:
      store = self._stores.get(name)
      if store is None:
        store = FakeKeyValueStore(name)
        self._stores[name] = store
      return store

  def replication_key_value_store(self) -> FakeKeyValueStore:
    """Returns the fake replicated store (`kv_store_repl`)."""
    return self.key_value_store(_REPLICATION_STORE)
