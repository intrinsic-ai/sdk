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

"""Tests extended status matcher."""

from absl.testing import absltest

from intrinsic.util.status import status_exception
from intrinsic.util.status import status_matcher


class StatusMatcherTest(absltest.TestCase):

  def test_matcher_component_and_code(self):
    with self.assertRaisesWithPredicateMatch(
        status_exception.ExtendedStatusError,
        status_matcher.matches(component="ai.testing.my_component", code=123),
    ):
      raise status_exception.ExtendedStatusError("ai.testing.my_component", 123)

  def test_matcher_component_and_code_mismatch(self):
    self.assertFalse(
        status_matcher.matches(component="ai.testing.my_component", code=123)(
            status_exception.ExtendedStatusError(
                "ai.testing.other_component", 123
            )
        )
    )
    self.assertFalse(
        status_matcher.matches(component="ai.testing.my_component", code=123)(
            status_exception.ExtendedStatusError("ai.testing.my_component", 432)
        )
    )

  def test_matcher_title(self):
    with self.assertRaisesWithPredicateMatch(
        status_exception.ExtendedStatusError,
        status_matcher.matches(title="My Title"),
    ):
      raise status_exception.ExtendedStatusError(
          "ai.testing.my_component", 123, title="My Title"
      )

  def test_matcher_title_mismatch(self):
    self.assertFalse(
        status_matcher.matches(
            component="ai.testing.my_component", code=123, title="My title"
        )(
            status_exception.ExtendedStatusError(
                "ai.testing.other_component", 123, title="Other title"
            )
        )
    )

  def test_matcher_user_message(self):
    with self.assertRaisesWithPredicateMatch(
        status_exception.ExtendedStatusError,
        status_matcher.matches(user_message="Ext message"),
    ):
      raise status_exception.ExtendedStatusError(
          "ai.testing.my_component", 123, user_message="Ext message"
      )

  def test_matcher_user_message_mismatch(self):
    self.assertFalse(
        status_matcher.matches(
            user_message="Ext message",
        )(
            status_exception.ExtendedStatusError(
                "ai.testing.my_component",
                123,
                user_message="Other message",
            )
        )
    )

  def test_matcher_user_message_regex(self):
    with self.assertRaisesWithPredicateMatch(
        status_exception.ExtendedStatusError,
        status_matcher.matches(user_message_regex=r"Ext \d+"),
    ):
      raise status_exception.ExtendedStatusError(
          "ai.testing.my_component", 123, user_message="Ext 353"
      )

  def test_matcher_user_message_regex_mismatch(self):
    self.assertFalse(
        status_matcher.matches(
            user_message_regex=r"Ext \d+",
        )(
            status_exception.ExtendedStatusError(
                "ai.testing.my_component",
                123,
                user_message="Ext message",
            )
        )
    )

  def test_matcher_debug_message(self):
    with self.assertRaisesWithPredicateMatch(
        status_exception.ExtendedStatusError,
        status_matcher.matches(debug_message="Int message"),
    ):
      raise status_exception.ExtendedStatusError(
          "ai.testing.my_component", 123, debug_message="Int message"
      )

  def test_matcher_debug_message_mismatch(self):
    self.assertFalse(
        status_matcher.matches(
            debug_message="Int message",
        )(
            status_exception.ExtendedStatusError(
                "ai.testing.my_component",
                123,
                debug_message="Other message",
            )
        )
    )

  def test_matcher_debug_message_regex(self):
    with self.assertRaisesWithPredicateMatch(
        status_exception.ExtendedStatusError,
        status_matcher.matches(debug_message_regex=r"Int \d+"),
    ):
      raise status_exception.ExtendedStatusError(
          "ai.testing.my_component", 123, debug_message="Int 353"
      )

  def test_matcher_debug_message_regex_mismatch(self):
    self.assertFalse(
        status_matcher.matches(
            debug_message_regex=r"Int \d+",
        )(
            status_exception.ExtendedStatusError(
                "ai.testing.my_component",
                123,
                debug_message="Int message",
            )
        )
    )


if __name__ == "__main__":
  absltest.main()
