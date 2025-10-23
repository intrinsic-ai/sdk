# Copyright 2023 Intrinsic Innovation LLC

"""Supports cooperative cancellation of skills by the skill service."""

import abc
from collections.abc import Callable
import threading
import time

_inner_wait_timeout_seconds = 0.001


class CallbackAlreadyRegisteredError(RuntimeError):
  """A callback was registered after a callback had already been registered."""


class SkillAlreadyCancelledError(RuntimeError):
  """A cancelled skill was cancelled again."""


class SkillReadyForCancellationError(RuntimeError):
  """An invalid action occurred on a skill that is ready for cancellation."""


class SkillCanceller(abc.ABC):
  """Supports cooperative cancellation of skills by the skill service.

  When a cancellation request is received, the skill should:
  1) stop as soon as possible and leave resources in a safe and recoverable
     state;
  2) raise skill_interface.SkillCancelledError.

  The skill must call `ready` once it is ready to be cancelled.

  A skill can implement cancellation in one of three ways:
  1) Poll `cancelled`, and safely cancel if and when it becomes true.
  2) Call `wait` and interrupt the skill if it returns true. `wait` is blocking,
     so the typical pattern is to call it from a thread. If the skill finishes
     without cancellation, the wait should be interrupted using `stop_wait`. For
     example:
     ```
     canceller = execute_context.canceller
     canceller.ready()

     def cancel_thread():
       if canceller.wait(timeout=float('inf')):
         cancel_work()  # Do work to safely cancel the skill.
     cancel_thread = threading.Thread(target=cancel_thread, daemon=True)
     cancel_thread.start()

     try:
       do_work()  # Do work that may be interrupted by cancellation.
     finally:
       canceller.stop_wait()
       cancel_thread.join()
     ```
  3) Register a callback via `register_callback`. This callback will be invoked
     when the skill receives a cancellation request.

  Attributes:
    cancelled: True if the skill has received a cancellation request.
  """

  @property
  @abc.abstractmethod
  def cancelled(self) -> bool:
    pass

  @abc.abstractmethod
  def ready(self) -> None:
    """Signals that the skill is ready to be cancelled."""

  @abc.abstractmethod
  def register_callback(self, callback: Callable[[], None]):
    """Sets a callback that will be invoked when cancellation is requested.

    If a callback will be used, it must be registered before calling `ready`.
    Only one callback may be registered, and the callback will be called at most
    once.

    Args:
      callback: The cancellation callback. Will be called at most once.

    Raises:
      CallbackAlreadyRegisteredError: If a callback was already registered.
      Exception: Any exception raised when calling the cancellation callback.
      SkillReadyForCancellationError: If the skill has already signaled that it
        is ready for cancellation.

      Raising an error will indicate that the skill could not be cancelled and
      the skill will be considered to be in an error state. Only raise an error
      if, after cancellation, the skill was not able to leave resources in a
      safe and recoverable state.
    """

  @abc.abstractmethod
  def wait(self, timeout: float) -> bool:
    """Waits for the skill to be cancelled.

    Args:
      timeout: The maximum number of seconds to wait for cancellation.

    Returns:
      cancelled: True if the skill was cancelled.
    """

  @abc.abstractmethod
  def stop_wait(self) -> None:
    """Unblocks `wait` if it is waiting."""


class SkillCancellationManager(SkillCanceller):
  """A SkillCanceller used by the skill service to cancel skills.

  Attributes:
    cancelled: True if the skill has received a cancellation request.
  """

  @property
  def cancelled(self) -> bool:
    return self._cancelled.is_set()

  def __init__(self, ready_timeout: float) -> None:
    """Initializes the instance.

    Args:
      ready_timeout: The maximum number of seconds to wait for the skill to be
        ready for cancellation before timing out.
    """
    self._ready_timeout = ready_timeout

    self._lock = threading.Lock()
    self._ready = threading.Event()
    self._cancelled = threading.Event()
    self._stop_wait = False
    self._callback = None

  def cancel(self) -> None:
    """Sets the cancelled flag and calls the callback (if set).

    Raises:
      Exception: Any exception raised when calling the callback.
      SkillAlreadyCancelledError: If the skill was already cancelled.
      TimeoutError: If we timeout while waiting for the skill to be ready for
        cancellation.
    """
    self.wait_for_ready()

    with self._lock:
      if self.cancelled:
        raise SkillAlreadyCancelledError("The skill was already cancelled.")
      self._cancelled.set()

      callback = self._callback

    if callback is not None:
      callback()

  def ready(self) -> None:
    """Signals that the skill is ready to be cancelled."""
    self._ready.set()

  def register_callback(self, callback: Callable[[], None]) -> None:
    """Sets a callback that will be invoked when cancellation is requested.

    If a callback will be used, it must be registered before calling `ready`.
    Only one callback may be registered, and the callback will be called at most
    once.

    Args:
      callback: The cancellation callback. Will be called at most once.

    Raises:
      CallbackAlreadyRegisteredError: If a callback was already registered.
      Exception: Any exception raised when calling the cancellation callback.
      SkillReadyForCancellationError: If the skill is already ready for
        cancellation.
    """
    with self._lock:
      if self._ready.is_set():
        raise SkillReadyForCancellationError(
            "A callback cannot be registered after the skill is ready for"
            " cancellation."
        )
      if self._callback is not None:
        raise CallbackAlreadyRegisteredError(
            "A callback was already registered."
        )

      self._callback = callback

  def wait(self, timeout: float) -> bool:
    """Waits for the skill to be cancelled.

    Args:
      timeout: The maximum number of seconds to wait for cancellation.

    Returns:
      cancelled: True if the skill was cancelled.
    """
    timeout_time = time.perf_counter() + timeout
    while time.perf_counter() < timeout_time and not self._stop_wait:
      if self._cancelled.wait(_inner_wait_timeout_seconds):
        break

    return self._cancelled.is_set()

  def stop_wait(self) -> None:
    """Unblocks `wait` if it is waiting."""
    self._stop_wait = True

  def wait_for_ready(self) -> None:
    """Waits for the skill to be ready for cancellation.

    Raises:
      TimeoutError: If we timeout while waiting for the skill to be ready for
        cancellation.
    """
    if not self._ready.wait(self._ready_timeout):
      raise TimeoutError(
          "Timed out waiting for the skill to be ready for cancellation."
      )
