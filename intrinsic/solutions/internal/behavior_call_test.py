# Copyright 2023 Intrinsic Innovation LLC

import datetime

from absl.testing import absltest
from google.protobuf import text_format
from intrinsic.executive.proto import behavior_call_pb2
from intrinsic.solutions.internal import actions
from intrinsic.solutions.internal import behavior_call
from intrinsic.solutions.testing import compare


def _create_behavior_call_proto(index: int) -> behavior_call_pb2.BehaviorCall:
  proto = behavior_call_pb2.BehaviorCall(skill_id=f'ai.intrinsic.skill-{index}')
  return proto


class BehaviorCallActionTest(absltest.TestCase):
  """Tests behavior_call.Action."""

  def test_init_from_proto(self):
    """Tests if BehaviorCallProto object can be constructed from proto."""
    empty_action = behavior_call.Action()
    compare.assertProto2Equal(
        self, empty_action.proto, behavior_call_pb2.BehaviorCall()
    )

    proto = _create_behavior_call_proto(123)
    action: actions.ActionBase = behavior_call.Action(proto)
    compare.assertProto2Equal(self, action.proto, proto)
    proto.skill_id = 'ai.intrinsic.different_name'
    compare.assertProto2Equal(
        self, action.proto, _create_behavior_call_proto(123)
    )

  def test_init_from_id(self):
    """Tests if Action object can be constructed from skill ID string."""
    proto = _create_behavior_call_proto(234)
    action = behavior_call.Action(skill_id=proto.skill_id)
    compare.assertProto2Equal(self, action.proto, proto)

  def test_set_proto(self):
    """Tests if proto can properly be read and set."""
    proto = _create_behavior_call_proto(123)
    action = behavior_call.Action()
    compare.assertProto2Equal(
        self, action.proto, behavior_call_pb2.BehaviorCall()
    )
    action.proto = proto
    compare.assertProto2Equal(self, action.proto, proto)

  def test_timeouts(self):
    """Tests if timeouts are transferred to proto."""
    proto = _create_behavior_call_proto(123)
    action = behavior_call.Action(proto)
    compare.assertProto2Equal(self, action.proto, proto)
    action.execute_timeout = datetime.timedelta(seconds=5)
    action.project_timeout = datetime.timedelta(seconds=10)
    proto.skill_execution_options.execute_timeout.FromTimedelta(
        datetime.timedelta(seconds=5)
    )
    proto.skill_execution_options.project_timeout.FromTimedelta(
        datetime.timedelta(seconds=10)
    )
    compare.assertProto2Equal(self, action.proto, proto)

  def test_str(self):
    """Tests if Action conversion to string works."""
    self.assertEqual(repr(behavior_call.Action()), r"""Action(skill_id='')""")
    self.assertEqual(str(behavior_call.Action()), r"""Action(skill_id='')""")

    proto = text_format.Parse(
        r"""
            skill_id: "ai.intrinsic.my_custom_action"
            resources {
              key: "device"
              value {
                handle: "SomeSpeaker"
              }
            }
        """,
        behavior_call_pb2.BehaviorCall(),
    )
    self.assertEqual(
        repr(behavior_call.Action(proto)),
        r"""Action(skill_id='ai.intrinsic.my_custom_action')."""
        r"""require(device={handle: "SomeSpeaker"})""",
    )
    self.assertEqual(
        str(behavior_call.Action(proto)),
        r"""Action(skill_id='ai.intrinsic.my_custom_action')."""
        r"""require(device={handle: "SomeSpeaker"})""",
    )

  def test_incomplete_abstract_class(self):
    """Tests incomplete actions are rejected."""

    class IncompleteAction(actions.ActionBase):
      # NOT defined: property proto
      pass

    with self.assertRaises(TypeError):
      _ = IncompleteAction()

  def test_require_resources(self):
    proto = behavior_call_pb2.BehaviorCall(
        skill_id='ai.intrinsic.my_custom_action'
    )
    proto.resources['robot'].handle = 'my_robot'

    action = behavior_call.Action(
        skill_id='ai.intrinsic.my_custom_action'
    ).require(robot='my_robot')
    compare.assertProto2Equal(self, action.proto, proto)


if __name__ == '__main__':
  absltest.main()
