# Copyright 2023 Intrinsic Innovation LLC

import datetime
import enum
import inspect
import textwrap
from unittest import mock

from absl.testing import absltest
from absl.testing import parameterized
from google.protobuf import any_pb2
from google.protobuf import descriptor_pb2
from google.protobuf import text_format
import grpc

from intrinsic.assets.configuration import asset_configuration_client
from intrinsic.assets.proto import id_pb2
from intrinsic.executive.proto import behavior_call_pb2
from intrinsic.executive.proto import test_message_pb2
from intrinsic.math.proto import point_pb2
from intrinsic.math.proto import pose_pb2
from intrinsic.math.proto import quaternion_pb2
from intrinsic.math.python import data_types
from intrinsic.math.python import proto_conversion as math_proto_conversion
from intrinsic.resources.proto import resource_handle_pb2
from intrinsic.skills.client import skill_registry_client
from intrinsic.skills.proto import equipment_pb2
from intrinsic.solutions import blackboard_value
from intrinsic.solutions import cel
from intrinsic.solutions import provided
from intrinsic.solutions.internal import skill_generation
from intrinsic.solutions.internal import skill_providing
from intrinsic.solutions.internal import skill_utils
from intrinsic.solutions.testing import compare
from intrinsic.solutions.testing import skill_test_utils
from intrinsic.solutions.testing import test_skill_params_pb2
from intrinsic.util.path_resolver import path_resolver
from third_party.ros2.ros_interfaces.jazzy.geometry_msgs.msg import point_pb2 as ros_point_pb2
from third_party.ros2.ros_interfaces.jazzy.geometry_msgs.msg import pose_pb2 as ros_pose_pb2
from third_party.ros2.ros_interfaces.jazzy.geometry_msgs.msg import quaternion_pb2 as ros_quaternion_pb2

_DEFAULT_TEST_MESSAGE = test_skill_params_pb2.TestMessage(
    my_double=2.5,
    my_float=-1.5,
    my_int32=5,
    my_int64=9,
    my_uint32=11,
    my_uint64=21,
    my_bool=False,
    my_string='bar',
    sub_message=test_skill_params_pb2.SubMessage(name='baz'),
    optional_sub_message=test_skill_params_pb2.SubMessage(name='quz'),
    my_repeated_doubles=[-5.5, 10.5],
    repeated_submessages=[
        test_skill_params_pb2.SubMessage(name='foo'),
        test_skill_params_pb2.SubMessage(name='bar'),
    ],
    my_required_int32=42,
    my_oneof_double=1.5,
    pose=pose_pb2.Pose(
        position=point_pb2.Point(),
        orientation=quaternion_pb2.Quaternion(x=0.5, y=0.5, z=0.5, w=0.5),
    ),
    ros_pose=ros_pose_pb2.Pose(
        position=ros_point_pb2.Point(),
        orientation=ros_quaternion_pb2.Quaternion(x=0.5, y=0.5, z=0.5, w=0.5),
    ),
    foo=test_skill_params_pb2.TestMessage.Foo(
        bar=test_skill_params_pb2.TestMessage.Foo.Bar(test='test')
    ),
    enum_v=test_skill_params_pb2.TestMessage.THREE,
    string_int32_map={'foo': 1},
    int32_string_map={3: 'foobar'},
    string_message_map={
        'bar': test_skill_params_pb2.TestMessage.MessageMapValue(value='baz')
    },
    executive_test_message=test_message_pb2.TestMessage(int32_value=123),
    non_unique_field_name=test_skill_params_pb2.TestMessage.SomeType(
        non_unique_field_name=test_skill_params_pb2.TestMessage.AnotherType()
    ),
)


def _test_skill_params_file_descriptor_set() -> (
    descriptor_pb2.FileDescriptorSet
):
  return skill_test_utils.read_file_descriptor_set(
      path_resolver.resolve_runfiles_path(
          'intrinsic/solutions/testing/test_skill_params_proto_descriptors_transitive_set_sci.proto.bin'
      )
  )


def _create_skill_registry_with_mock_stub():
  skill_registry_stub = mock.MagicMock()
  skill_registry = skill_registry_client.SkillRegistryClient(
      skill_registry_stub
  )

  return (skill_registry, skill_registry_stub)


class SkillsTest(parameterized.TestCase):
  """Tests public methods of the skills wrapper class."""

  def assertSignature(self, actual, expected):
    actual = str(actual)
    # Insert newlines after commas, else the diff printed by assertEqual is
    # unreadable. We don't need to assert on newlines in the signature exactly.
    actual_with_newlines = actual.replace(',', ',\n')
    expected_with_newlines = expected.replace(',', ',\n')
    self.assertEqual(
        actual_with_newlines,
        expected_with_newlines,
        # In case of failure print code for updating the test. Not a great
        # solution, but our signatures can be very long and this can save
        # **a lot** of time when making signature changes.
        'Signature not as expected. If you made changes and want to update the '
        "test expectation you can use the following code:\n\n'"
        + actual.replace('\n', '\\n').replace(', ', ", '\n'")
        + "'\n\nSignature not as expected. Diff with inserted newlines:\n\n",
    )

  @parameterized.parameters(
      {'parameter': {'my_double': 2}},
      {'parameter': {'my_float': 1}},
      {'parameter': {'my_int32': 1.0}},
      {'parameter': {'my_int64': 2.0}},
      {'parameter': {'my_uint32': 1.1}},
      {'parameter': {'my_uint64': 2.1}},
      {'parameter': {'my_bool': 'foo'}},
      {'parameter': {'my_string': -1}},
      {'parameter': {'sub_message': data_types.Pose3()}},
      {'parameter': {'pose': test_skill_params_pb2.SubMessage()}},
      {'parameter': {'my_repeated_doubles': 0.0}},
      {
          'parameter': {
              'my_repeated_doubles': [test_skill_params_pb2.SubMessage()]
          }
      },
      {'parameter': {'my_repeated_doubles': [data_types.Pose3()]}},
      {'parameter': {'my_repeated_doubles': [True, False]}},
      {'parameter': {'my_repeated_doubles': [1.0, False]}},
      {'parameter': {'repeated_submessages': [1.0, False]}},
      {
          'parameter': {
              'repeated_submessages': [data_types.Pose3(), data_types.Pose3()]
          }
      },
      {'parameter': {'repeated_submessages': {'foo': 1}}},
      {'parameter': {'my_double': [1, 2]}},
  )
  def test_gen_skill_param_message_type_mismatch(self, parameter):
    skill = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
        default_params=_DEFAULT_TEST_MESSAGE,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([skill]),
        skill_test_utils.create_asset_configuration_client(),
    )

    with self.assertRaises(TypeError):
      skills.ai.intrinsic.my_skill(**parameter)

  def test_list_skills(self):
    registry_skills = [
        skill_test_utils.create_legacy_process('ai.ai_pbt'),
        skill_test_utils.create_legacy_process('ai.intr.intr_pbt'),
        skill_test_utils.create_legacy_process('foo.bar.bar_pbt'),
        skill_test_utils.create_legacy_process('foo.foo_pbt'),
        skill_test_utils.create_legacy_process('foo.foo'),
        skill_test_utils.create_legacy_process('global_pbt'),
    ]
    assets = [
        # Regular skills
        skill_test_utils.create_skill_asset('ai.intr.intr_skill_a'),
        skill_test_utils.create_skill_asset('ai.intr.intr_skill_b'),
        skill_test_utils.create_skill_asset('foo.bar.bar_skill'),
        skill_test_utils.create_skill_asset('foo.bar.bar'),
        # Process assets
        skill_test_utils.create_process_asset('ai.intr.intr_process_a'),
        skill_test_utils.create_process_asset('ai.intr.intr_process_b'),
    ]
    skills = skill_providing.Skills(
        skill_test_utils.create_skill_registry_for_skill_infos(registry_skills),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets(assets),
        skill_test_utils.create_asset_configuration_client(),
    )

    self.assertEqual(
        dir(skills),
        [
            'ai',
            'foo',
            'global_pbt',
        ],
    )
    self.assertEqual(dir(skills.ai), ['ai_pbt', 'intr'])
    self.assertEqual(
        dir(skills.ai.intr),
        [
            'intr_pbt',
            'intr_process_a',
            'intr_process_b',
            'intr_skill_a',
            'intr_skill_b',
        ],
    )
    self.assertEqual(dir(skills.foo), ['bar', 'foo', 'foo_pbt'])
    self.assertEqual(dir(skills.foo.bar), ['bar', 'bar_pbt', 'bar_skill'])

  def test_skills_attribute_access(self):
    """Tests attribute-based access (skills.<skill_id>)."""
    registry_skills = [
        skill_test_utils.create_legacy_process('ai.intr.intr_pbt'),
        skill_test_utils.create_legacy_process('ai.ai_pbt'),
        skill_test_utils.create_legacy_process('foo.foo_pbt'),
        skill_test_utils.create_legacy_process('global_pbt'),
    ]
    assets = [
        # Regular skills
        skill_test_utils.create_skill_asset('ai.intr.intr_skill'),
        skill_test_utils.create_skill_asset('foo.bar.bar_skill'),
        # Process assets
        skill_test_utils.create_process_asset('ai.intr.intr_process'),
        skill_test_utils.create_process_asset('foo.bar.bar_process'),
    ]
    skills = skill_providing.Skills(
        skill_test_utils.create_skill_registry_for_skill_infos(registry_skills),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets(assets),
        skill_test_utils.create_asset_configuration_client(),
    )

    # Id notation: skills.<skill_id>
    self.assertIsInstance(skills.ai, provided.SkillPackage)
    self.assertIsInstance(skills.ai.intr, provided.SkillPackage)
    self.assertIsInstance(skills.foo, provided.SkillPackage)
    self.assertIsInstance(skills.foo.bar, provided.SkillPackage)

    self.assertIsInstance(
        skills.ai.intr.intr_pbt(), skill_generation.GeneratedSkill
    )
    self.assertIsInstance(skills.ai.ai_pbt(), skill_generation.GeneratedSkill)
    self.assertIsInstance(skills.foo.foo_pbt(), skill_generation.GeneratedSkill)
    self.assertIsInstance(skills.global_pbt(), skill_generation.GeneratedSkill)
    self.assertIsInstance(
        skills.ai.intr.intr_skill(), skill_generation.GeneratedSkill
    )
    self.assertIsInstance(
        skills.foo.bar.bar_skill(), skill_generation.GeneratedSkill
    )
    self.assertIsInstance(
        skills.ai.intr.intr_process(), skill_generation.GeneratedSkill
    )
    self.assertIsInstance(
        skills.foo.bar.bar_process(), skill_generation.GeneratedSkill
    )

    with self.assertRaises(AttributeError):
      _ = skills.skill5

    with self.assertRaises(AttributeError):
      _ = skills.foo.skill5

  def test_skills_with_same_name_but_different_packages(self):
    assets = [
        skill_test_utils.create_skill_asset('ai.intr.my_skill'),
        skill_test_utils.create_skill_asset('com.foo.my_skill'),
    ]

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets(assets),
        skill_test_utils.create_asset_configuration_client(),
    )

    self.assertIsInstance(
        skills.ai.intr.my_skill(), skill_generation.GeneratedSkill
    )
    self.assertIsInstance(
        skills.com.foo.my_skill(), skill_generation.GeneratedSkill
    )

  def test_asset_skill_takes_precedence_over_registry_skill(self):
    registry_pbt = skill_test_utils.create_legacy_process(
        'ai.intr.same_name', description='registry skill description'
    )
    skill_asset = skill_test_utils.create_skill_asset(
        'ai.intr.same_name', description='asset skill description'
    )

    skills = skill_providing.Skills(
        skill_test_utils.create_skill_registry_for_skill_infos([registry_pbt]),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([skill_asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    self.assertEqual(
        skills.ai.intr.same_name.info.description, 'asset skill description'
    )

  def test_skills_dict_access(self):
    """Tests id-string-based access via __getitem__ (skills['<skill_id>'])."""
    registry_skills = [
        skill_test_utils.create_legacy_process('global_pbt'),
    ]
    assets = [
        skill_test_utils.create_skill_asset('ai.intr.skill_two'),
        skill_test_utils.create_skill_asset('ai.intr.skill_one'),
        skill_test_utils.create_process_asset('ai.intr.process'),
    ]
    skills = skill_providing.Skills(
        skill_test_utils.create_skill_registry_for_skill_infos(registry_skills),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets(assets),
        skill_test_utils.create_asset_configuration_client(),
    )

    self.assertIsInstance(
        skills['ai.intr.skill_one'](), skill_generation.GeneratedSkill
    )
    self.assertIsInstance(
        skills['ai.intr.skill_two'](), skill_generation.GeneratedSkill
    )
    self.assertIsInstance(
        skills['ai.intr.process'](), skill_generation.GeneratedSkill
    )
    self.assertIsInstance(
        skills['global_pbt'](), skill_generation.GeneratedSkill
    )

    with self.assertRaises(KeyError):
      _ = skills['skill5']

  def test_skills_get_skill_ids(self):
    registry_skills = [
        skill_test_utils.create_legacy_process('global_pbt'),
    ]
    assets = [
        skill_test_utils.create_skill_asset('ai.intr.skill_two'),
        skill_test_utils.create_skill_asset('ai.intr.skill_one'),
        skill_test_utils.create_process_asset('ai.intr.process'),
    ]
    skills = skill_providing.Skills(
        skill_test_utils.create_skill_registry_for_skill_infos(registry_skills),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets(assets),
        skill_test_utils.create_asset_configuration_client(),
    )

    skill_ids = skills.get_skill_ids()

    self.assertCountEqual(
        list(skill_ids),
        [
            'ai.intr.skill_one',
            'ai.intr.skill_two',
            'ai.intr.process',
            'global_pbt',
        ],
    )

  def test_skills_get_skill_classes(self):
    registry_skills = [
        skill_test_utils.create_legacy_process('global_pbt'),
    ]
    assets = [
        skill_test_utils.create_skill_asset('ai.intr.skill_two'),
        skill_test_utils.create_skill_asset('ai.intr.skill_one'),
        skill_test_utils.create_process_asset('ai.intr.process'),
    ]
    skills = skill_providing.Skills(
        skill_test_utils.create_skill_registry_for_skill_infos(registry_skills),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets(assets),
        skill_test_utils.create_asset_configuration_client(),
    )

    skill_classes = skills.get_skill_classes()

    self.assertLen(skill_classes, 4)
    for skill_class in skill_classes:
      self.assertIsInstance(skill_class(), skill_generation.GeneratedSkill)

  def test_skills_get_skill_ids_and_classes(self):
    registry_skills = [
        skill_test_utils.create_legacy_process('global_pbt'),
    ]
    assets = [
        skill_test_utils.create_skill_asset('ai.intr.skill_two'),
        skill_test_utils.create_skill_asset('ai.intr.skill_one'),
        skill_test_utils.create_process_asset('ai.intr.process'),
    ]
    skills = skill_providing.Skills(
        skill_test_utils.create_skill_registry_for_skill_infos(registry_skills),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets(assets),
        skill_test_utils.create_asset_configuration_client(),
    )

    ids_and_classes = skills.get_skill_ids_and_classes()

    self.assertCountEqual(
        [skill_id for skill_id, _ in ids_and_classes],
        [
            'ai.intr.skill_one',
            'ai.intr.skill_two',
            'ai.intr.process',
            'global_pbt',
        ],
    )
    for _, skill_class in ids_and_classes:
      self.assertIsInstance(skill_class(), skill_generation.GeneratedSkill)

  def test_type_url_areas(self):
    registry_skills = [
        skill_test_utils.create_legacy_process('ai.intr.legacy_pbt')
    ]
    assets = [
        skill_test_utils.create_skill_asset('ai.intr.skill'),
        skill_test_utils.create_process_asset('ai.intr.process'),
    ]
    skills = skill_providing.Skills(
        skill_test_utils.create_skill_registry_for_skill_infos(registry_skills),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets(assets),
        skill_test_utils.create_asset_configuration_client(),
    )

    self.assertEqual(skills.ai.intr.legacy_pbt.info.type_url_area, 'skills')
    self.assertEqual(skills.ai.intr.skill.info.type_url_area, 'assets')
    self.assertEqual(skills.ai.intr.process.info.type_url_area, 'assets')

  def test_skill_type(self):
    registry_skills = [
        skill_test_utils.create_legacy_process('ai.intr.legacy_pbt')
    ]
    assets = [
        skill_test_utils.create_skill_asset('ai.intr.skill'),
        skill_test_utils.create_process_asset('ai.intr.process'),
    ]
    skills = skill_providing.Skills(
        skill_test_utils.create_skill_registry_for_skill_infos(registry_skills),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets(assets),
        skill_test_utils.create_asset_configuration_client(),
    )

    self.assertEqual(
        skills.ai.intr.legacy_pbt.info.skill_type, provided.SkillType.PROCESS
    )
    self.assertEqual(
        skills.ai.intr.skill.info.skill_type, provided.SkillType.REGULAR_SKILL
    )
    self.assertEqual(
        skills.ai.intr.process.info.skill_type, provided.SkillType.PROCESS
    )

  def test_gen_skill(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
        resource_selectors={'my_resource_slot': ['my_capability']},
    )

    resource_registry = (
        skill_test_utils.create_resource_registry_with_single_handle(
            'my_resource', 'my_capability'
        )
    )

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        resource_registry,
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    expected_repeated_doubles = [2.1, 3.1]
    expected_repeated_submessages = [
        test_skill_params_pb2.SubMessage(name='foo'),
        test_skill_params_pb2.SubMessage(name='bar'),
    ]

    parameters = test_skill_params_pb2.TestMessage(
        my_double=1.1,
        my_float=2.0,
        my_int32=1,
        my_int64=2,
        my_string='foo',
        my_uint32=10,
        my_uint64=20,
        my_bool=True,
        sub_message=test_skill_params_pb2.SubMessage(name='bar'),
        pose=math_proto_conversion.pose_to_proto(data_types.Pose3()),
        my_repeated_doubles=expected_repeated_doubles,
        repeated_submessages=expected_repeated_submessages,
        my_oneof_double=2.0,
    )
    skill = skills.ai.intrinsic.my_skill(
        my_double=parameters.my_double,
        my_float=parameters.my_float,
        my_int32=parameters.my_int32,
        my_int64=parameters.my_int64,
        my_string=parameters.my_string,
        my_uint32=parameters.my_uint32,
        my_uint64=parameters.my_uint64,
        my_bool=parameters.my_bool,
        sub_message=parameters.sub_message,
        pose=data_types.Pose3(),
        my_repeated_doubles=expected_repeated_doubles,
        repeated_submessages=expected_repeated_submessages,
        my_oneof_double=parameters.my_oneof_double,
        my_resource_slot=(
            provided.ResourceHandle.create('my_resource', ['my_capability'])
        ),
    )

    expected_proto = behavior_call_pb2.BehaviorCall(
        skill_id='ai.intrinsic.my_skill',
        return_value_name=skill.proto.return_value_name,
    )
    expected_proto.resources['my_resource_slot'].handle = 'my_resource'
    expected_proto.parameters.Pack(
        parameters,
        type_url_prefix='type.intrinsic.ai/assets/ai.intrinsic.my_skill',
    )

    compare.assertProto2Equal(self, expected_proto, skill.proto)

    skill_str = (
        'skills.ai.intrinsic.my_skill('
        'my_double=1.1, '
        'my_float=2.0, '
        'my_int32=1, '
        'my_int64=2, '
        'my_uint32=10, '
        'my_uint64=20, '
        'my_bool=True, '
        "my_string='foo', "
        'sub_message=name: "bar"\n, '
        'my_repeated_doubles=[2.1, 3.1], '
        'repeated_submessages=[name: "foo"\n, name: "bar"\n], '
        'my_oneof_double=2.0, '
        'pose=position {\n}\norientation {\n  w: 1\n}\n, '
        'my_resource_slot={handle: "my_resource"})'
    )
    self.assertEqual(str(skill), skill_str)

  def test_gen_skill_uses_defaults(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
        default_params=_DEFAULT_TEST_MESSAGE,
        resource_selectors={'my_resource_slot': ['my_capability']},
    )

    resource_registry = (
        skill_test_utils.create_resource_registry_with_single_handle(
            'my_resource', 'my_capability'
        )
    )

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        resource_registry,
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    skill = skills.ai.intrinsic.my_skill()

    expected_proto = behavior_call_pb2.BehaviorCall(
        skill_id='ai.intrinsic.my_skill',
        return_value_name=skill.proto.return_value_name,
    )
    expected_proto.resources['my_resource_slot'].handle = 'my_resource'
    expected_proto.parameters.Pack(
        _DEFAULT_TEST_MESSAGE,
        type_url_prefix='type.intrinsic.ai/assets/ai.intrinsic.my_skill',
    )

    compare.assertProto2Equal(self, expected_proto, skill.proto)

    skill_str = (
        'skills.ai.intrinsic.my_skill('
        'my_double=2.5, '
        'my_float=-1.5, '
        'my_int32=5, '
        'my_int64=9, '
        'my_uint32=11, '
        'my_uint64=21, '
        'my_bool=False, '
        "my_string='bar', "
        'sub_message=name: "baz"\n, '
        'optional_sub_message=name: "quz"\n, '
        'my_repeated_doubles=[-5.5, 10.5], '
        'repeated_submessages=[name: "foo"\n, name: "bar"\n], '
        'my_required_int32=42, '
        'my_oneof_double=1.5, '
        'pose=position {\n}\norientation {\n  x: 0.5\n  y: 0.5\n  z: 0.5\n  w:'
        ' 0.5\n}\n, '
        'foo=bar {\n  test: "test"\n}\n, '
        'enum_v=3, '
        'executive_test_message=int32_value: 123\n, '
        'string_int32_map={"foo": 1}, '
        "int32_string_map={3: 'foobar'}, "
        'string_message_map={"bar": value: "baz"\n}, '
        'non_unique_field_name=non_unique_field_name {\n}\n, '
        'ros_pose=position {\n}\norientation {\n  x: 0.5\n  y: 0.5\n  z:'
        ' 0.5\n  w:'
        ' 0.5\n}\n, '
        'my_resource_slot={handle: "my_resource"})'
    )
    self.assertEqual(str(skill), skill_str)

  def test_gen_skill_nested_map(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
        resource_selectors={'my_resource_slot': ['my_capability']},
    )

    resource_registry = (
        skill_test_utils.create_resource_registry_with_single_handle(
            'my_resource', 'my_capability'
        )
    )

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        resource_registry,
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    parameters = test_skill_params_pb2.TestMessage(
        executive_test_message=test_message_pb2.TestMessage(
            string_int32_map={'foo': 2}
        )
    )
    skill = skills.ai.intrinsic.my_skill(
        executive_test_message=parameters.executive_test_message,
        my_resource_slot=provided.ResourceHandle.create(
            'my_resource', ['my_capability']
        ),
    )

    expected_proto = behavior_call_pb2.BehaviorCall(
        skill_id='ai.intrinsic.my_skill',
        return_value_name=skill.proto.return_value_name,
    )
    expected_proto.resources['my_resource_slot'].handle = 'my_resource'
    expected_proto.parameters.Pack(
        parameters,
        type_url_prefix='type.intrinsic.ai/assets/ai.intrinsic.my_skill',
    )

    compare.assertProto2Equal(self, expected_proto, skill.proto)

    skill_str = (
        'skills.ai.intrinsic.my_skill('
        'executive_test_message=string_int32_map {\n  key: "foo"\n  value:'
        ' 2\n}\n, '
        'my_resource_slot={handle: "my_resource"})'
    )
    self.assertEqual(str(skill), skill_str)

  @parameterized.parameters(
      {
          'value_specification': blackboard_value.BlackboardValue(
              {}, 'test', None, None
          )
      },
      {'value_specification': cel.CelExpression('test')},
  )
  def test_gen_skill_with_blackboard_parameter(self, value_specification):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
        default_params=_DEFAULT_TEST_MESSAGE,
    )

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    skill = skills.ai.intrinsic.my_skill(
        my_oneof_double=value_specification,
    )

    expected_proto = behavior_call_pb2.BehaviorCall(
        skill_id='ai.intrinsic.my_skill',
        return_value_name=skill.proto.return_value_name,
    )
    expected_proto.parameters.Pack(
        _DEFAULT_TEST_MESSAGE,
        type_url_prefix='type.intrinsic.ai/assets/ai.intrinsic.my_skill',
    )
    expected_proto.assignments.append(
        behavior_call_pb2.BehaviorCall.ParameterAssignment(
            parameter_path='my_oneof_double', cel_expression='test'
        )
    )

    compare.assertProto2Equal(self, expected_proto, skill.proto)

  @parameterized.parameters(
      {
          'value_specification': blackboard_value.BlackboardValue(
              {}, 'test', None, None
          )
      },
      {'value_specification': cel.CelExpression('test')},
  )
  def test_gen_skill_with_nested_blackboard_parameter(
      self, value_specification
  ):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    my_skill = skills.ai.intrinsic.my_skill

    skill = my_skill(
        my_oneof_double=value_specification,
        executive_test_message=my_skill.intrinsic_proto.executive.TestMessage(
            message_list=[
                my_skill.intrinsic_proto.executive.TestMessage(
                    int32_value=value_specification
                ),
                my_skill.intrinsic_proto.executive.TestMessage(
                    message_value=value_specification
                ),
                my_skill.intrinsic_proto.executive.TestMessage(
                    foo_msg=value_specification
                ),
                my_skill.intrinsic_proto.executive.TestMessage(
                    message_list=[
                        my_skill.intrinsic_proto.executive.TestMessage(
                            message_value=value_specification
                        )
                    ]
                ),
                my_skill.intrinsic_proto.executive.TestMessage(
                    message_list=[
                        my_skill.intrinsic_proto.executive.TestMessage(
                            message_list=value_specification
                        )
                    ]
                ),
                my_skill.intrinsic_proto.executive.TestMessage(
                    message_list=[value_specification, value_specification]
                ),
                value_specification,
                my_skill.intrinsic_proto.executive.TestMessage(
                    int32_list=[value_specification, value_specification]
                ),
                # The following is NOT supported, we cannot set a map value from
                # a CEL expression (we cannot address the field with proto_path)
                # my_skill.intrinsic_proto.executive.TestMessage(
                #   string_int32_map={'foo_key': value_specification}
                # ),
            ]
        ),
    )

    expected_proto = behavior_call_pb2.BehaviorCall(
        skill_id='ai.intrinsic.my_skill',
        return_value_name=skill.proto.return_value_name,
        assignments=[
            behavior_call_pb2.BehaviorCall.ParameterAssignment(
                parameter_path='my_oneof_double', cel_expression='test'
            ),
            behavior_call_pb2.BehaviorCall.ParameterAssignment(
                parameter_path=(
                    'executive_test_message.message_list[0].int32_value'
                ),
                cel_expression='test',
            ),
            behavior_call_pb2.BehaviorCall.ParameterAssignment(
                parameter_path=(
                    'executive_test_message.message_list[1].message_value'
                ),
                cel_expression='test',
            ),
            behavior_call_pb2.BehaviorCall.ParameterAssignment(
                parameter_path='executive_test_message.message_list[2].foo_msg',
                cel_expression='test',
            ),
            behavior_call_pb2.BehaviorCall.ParameterAssignment(
                parameter_path='executive_test_message.message_list[3].message_list[0].message_value',
                cel_expression='test',
            ),
            behavior_call_pb2.BehaviorCall.ParameterAssignment(
                parameter_path='executive_test_message.message_list[4].message_list[0].message_list',
                cel_expression='test',
            ),
            behavior_call_pb2.BehaviorCall.ParameterAssignment(
                parameter_path=(
                    'executive_test_message.message_list[5].message_list[0]'
                ),
                cel_expression='test',
            ),
            behavior_call_pb2.BehaviorCall.ParameterAssignment(
                parameter_path=(
                    'executive_test_message.message_list[5].message_list[1]'
                ),
                cel_expression='test',
            ),
            behavior_call_pb2.BehaviorCall.ParameterAssignment(
                parameter_path='executive_test_message.message_list[6]',
                cel_expression='test',
            ),
            behavior_call_pb2.BehaviorCall.ParameterAssignment(
                parameter_path=(
                    'executive_test_message.message_list[7].int32_list[0]'
                ),
                cel_expression='test',
            ),
            behavior_call_pb2.BehaviorCall.ParameterAssignment(
                parameter_path=(
                    'executive_test_message.message_list[7].int32_list[1]'
                ),
                cel_expression='test',
            ),
        ],
    )

    expected_parameters = test_skill_params_pb2.TestMessage(
        executive_test_message=test_message_pb2.TestMessage(
            message_list=[
                test_message_pb2.TestMessage(),
                test_message_pb2.TestMessage(),
                test_message_pb2.TestMessage(),
                test_message_pb2.TestMessage(
                    message_list=[test_message_pb2.TestMessage()]
                ),
                test_message_pb2.TestMessage(
                    message_list=[test_message_pb2.TestMessage()]
                ),
                test_message_pb2.TestMessage(
                    message_list=[
                        test_message_pb2.TestMessage(),
                        test_message_pb2.TestMessage(),
                    ]
                ),
                test_message_pb2.TestMessage(),
                test_message_pb2.TestMessage(int32_list=[0, 0]),
            ]
        )
    )
    expected_proto.parameters.Pack(
        expected_parameters,
        type_url_prefix='type.intrinsic.ai/assets/ai.intrinsic.my_skill',
    )
    compare.assertProto2Equal(self, expected_proto, skill.proto)

  @parameterized.parameters(
      {
          'value_specification': blackboard_value.BlackboardValue(
              {}, 'test', None, None
          )
      },
      {'value_specification': cel.CelExpression('test')},
  )
  def test_gen_skill_with_blackboard_parameter_list(self, value_specification):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    my_skill = skills.ai.intrinsic.my_skill

    skill = my_skill(repeated_submessages=[value_specification])

    expected_proto = behavior_call_pb2.BehaviorCall(
        skill_id='ai.intrinsic.my_skill',
        return_value_name=skill.proto.return_value_name,
        assignments=[
            behavior_call_pb2.BehaviorCall.ParameterAssignment(
                parameter_path='repeated_submessages[0]',
                cel_expression='test',
            ),
        ],
    )

    expected_parameters = test_skill_params_pb2.TestMessage(
        repeated_submessages=[test_skill_params_pb2.SubMessage()]
    )
    expected_proto.parameters.Pack(
        expected_parameters,
        type_url_prefix='type.intrinsic.ai/assets/ai.intrinsic.my_skill',
    )
    compare.assertProto2Equal(self, expected_proto, skill.proto)

  def test_gen_skill_with_map_parameter(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    expected_parameters = test_skill_params_pb2.TestMessage(
        string_int32_map={'foo': 1}
    )

    skill = skills.ai.intrinsic.my_skill(string_int32_map={'foo': 1})

    actual_parameters = test_skill_params_pb2.TestMessage()
    skill.proto.parameters.Unpack(actual_parameters)

    compare.assertProto2Equal(self, expected_parameters, actual_parameters)

  def test_gen_skill_with_message_map_parameter_from_alias(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    expected_parameters = test_skill_params_pb2.TestMessage(
        string_message_map={
            'foo': test_skill_params_pb2.TestMessage.MessageMapValue(
                value='bar'
            )
        }
    )

    my_skill = skills.ai.intrinsic.my_skill

    skill = my_skill(
        string_message_map={
            'foo': (
                my_skill.intrinsic_proto.test_data.TestMessage.MessageMapValue(
                    value='bar'
                )
            )
        }
    )

    actual_parameters = test_skill_params_pb2.TestMessage()
    skill.proto.parameters.Unpack(actual_parameters)

    compare.assertProto2Equal(self, expected_parameters, actual_parameters)

  def test_gen_skill_with_message_map_parameter_from_actual_type(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    expected_parameters = test_skill_params_pb2.TestMessage(
        string_message_map={
            'foo': test_skill_params_pb2.TestMessage.MessageMapValue(
                value='bar'
            )
        }
    )

    skill = skills.ai.intrinsic.my_skill(
        string_message_map={
            'foo': test_skill_params_pb2.TestMessage.MessageMapValue(
                value='bar'
            )
        }
    )

    actual_parameters = test_skill_params_pb2.TestMessage()
    skill.proto.parameters.Unpack(actual_parameters)

    compare.assertProto2Equal(self, expected_parameters, actual_parameters)

  def test_gen_skill_with_return_value_key(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        return_value_message=test_skill_params_pb2.SubMessage,
    )

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    my_skill = skills.ai.intrinsic.my_skill

    self.assertRegex(my_skill().proto.return_value_name, '^my_skill_.*$')
    self.assertRegex(
        my_skill(return_value_key=None).proto.return_value_name, '^my_skill_.*$'
    )
    self.assertEqual(
        my_skill(return_value_key='foo').proto.return_value_name, 'foo'
    )

  def test_gen_skill_fails_for_set_instead_of_dict(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    with self.assertRaisesRegex(TypeError, 'Got set where expected dict'):
      # In the following, the map parameter should be initialized as {'foo': 1},
      # with a colon instead of a comma.
      skills.ai.intrinsic.my_skill(string_int32_map={'foo', 1})

  def test_gen_skill_map_rejects_blackboard_value(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    with self.assertRaisesRegex(
        TypeError, 'Cannot set field .* from blackboard'
    ):
      skills.ai.intrinsic.my_skill(
          string_int32_map={
              'foo': blackboard_value.BlackboardValue({}, 'test', None, None),
          }
      )

  def test_gen_skill_with_invalid_parameter(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    # Normal assignment to known field should work
    m = skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.SubMessage(
        name='hello'
    )
    expected_proto = test_skill_params_pb2.SubMessage(name='hello')
    compare.assertProto2Equal(self, expected_proto, m.wrapped_message)

    # Assignment to unknown field should fail
    with self.assertRaisesRegex(KeyError, 'not_a_field.*does not exist'):
      skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.SubMessage(
          not_a_field='hello'
      )

    blackboard_test_value = blackboard_value.BlackboardValue(
        {}, 'test', None, None
    )
    # Assigning blackboard value to unknown field should fail
    with self.assertRaisesRegex(KeyError, 'not_a_field.*does not exist'):
      skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.SubMessage(
          not_a_field=blackboard_test_value
      )

  def test_compatible_resources(self):
    assets = [
        skill_test_utils.create_skill_asset(
            'ai.intrinsic.skill_one',
            resource_selectors={'slot_one': ['capability_a', 'capability_b']},
        ),
        skill_test_utils.create_skill_asset(
            'ai.intrinsic.skill_two',
            resource_selectors={
                'slot_one': ['capability_a'],
                'slot_two': ['capability_b'],
            },
        ),
        skill_test_utils.create_skill_asset(
            'ai.intrinsic.skill_three', resource_selectors={}
        ),
        skill_test_utils.create_skill_asset(
            'ai.intrinsic.skill_four',
            resource_selectors={
                'slot_one': ['capability_not_matched_by_any_resource']
            },
        ),
    ]

    resource_registry = skill_test_utils.create_resource_registry_with_handles([
        text_format.Parse(
            """name: 'a_resource'
               resource_data { key: 'capability_a' }""",
            resource_handle_pb2.ResourceHandle(),
        ),
        text_format.Parse(
            """name: 'a_b_resource'
               resource_data { key: 'capability_a' }
               resource_data { key: 'capability_b' }""",
            resource_handle_pb2.ResourceHandle(),
        ),
        text_format.Parse(
            """name: 'b_resource'
               resource_data { key: 'capability_b' }""",
            resource_handle_pb2.ResourceHandle(),
        ),
    ])

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        resource_registry,
        skill_test_utils.create_installed_assets(assets),
        skill_test_utils.create_asset_configuration_client(),
    )

    self.assertCountEqual(
        dir(skills.ai.intrinsic.skill_one.compatible_resources['slot_one']),
        ['a_b_resource'],
    )
    self.assertCountEqual(
        dir(skills.ai.intrinsic.skill_two.compatible_resources['slot_one']),
        ['a_b_resource', 'a_resource'],
    )
    self.assertCountEqual(
        dir(skills.ai.intrinsic.skill_two.compatible_resources['slot_two']),
        ['b_resource', 'a_b_resource'],
    )
    self.assertCountEqual(
        dir(skills.ai.intrinsic.skill_three.compatible_resources), []
    )
    self.assertCountEqual(
        dir(skills.ai.intrinsic.skill_four.compatible_resources['slot_one']),
        [],
    )

  def test_gen_skill_incompatible_resources(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
        resource_selectors={
            'my_resource_slot': ['capability_a', 'capability_b']
        },
    )

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    resource_a = provided.ResourceHandle.create('resource_a', ['capability_a'])

    with self.assertRaises(TypeError):
      skills.ai.intrinsic.my_skill(my_resource_slot=resource_a)

  def test_skill_class_name(self):
    legacy_pbt = skill_test_utils.create_legacy_process('global_pbt')
    asset = skill_test_utils.create_skill_asset('ai.intrinsic.my_skill')
    skills = skill_providing.Skills(
        skill_test_utils.create_skill_registry_for_skill_info(legacy_pbt),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    self.assertEqual(skills.ai.intrinsic.my_skill.__name__, 'my_skill')
    self.assertEqual(skills.ai.intrinsic.my_skill.__qualname__, 'my_skill')
    self.assertEqual(
        skills.ai.intrinsic.my_skill.__module__,
        'intrinsic.solutions.skills.ai.intrinsic',
    )

    self.assertEqual(skills.global_pbt.__name__, 'global_pbt')
    self.assertEqual(skills.global_pbt.__qualname__, 'global_pbt')
    self.assertEqual(
        skills.global_pbt.__module__,
        'intrinsic.solutions.skills',
    )

  def test_skill_uses_recommended_config(self):
    default_params = test_skill_params_pb2.TestMessage(my_double=2.5)
    recommended_params = test_skill_params_pb2.TestMessage(
        my_double=2.5, my_float=99.9
    )
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
        default_params=default_params,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client({
            'ai.intrinsic.my_skill': [(default_params, recommended_params)],
        }),
    )

    skill = skills.ai.intrinsic.my_skill()

    expected_proto = behavior_call_pb2.BehaviorCall(
        skill_id='ai.intrinsic.my_skill',
        return_value_name=skill.proto.return_value_name,
    )
    expected_proto.parameters.Pack(
        recommended_params,
        type_url_prefix='type.intrinsic.ai/assets/ai.intrinsic.my_skill',
    )

    compare.assertProto2Equal(self, expected_proto, skill.proto)

  def test_skill_skips_recommended_config(self):
    default_params = test_skill_params_pb2.TestMessage(my_double=2.5)
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
        default_params=default_params,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        # Configure no recommendations for my_skill. This causes an error if the
        # asset configuration client is called in any form.
        skill_test_utils.create_asset_configuration_client({
            'ai.intrinsic.my_skill': [],
        }),
    )

    skill = skills.ai.intrinsic.my_skill(
        _with_recommended_config=False,
    )

    expected_proto = behavior_call_pb2.BehaviorCall(
        skill_id='ai.intrinsic.my_skill',
        return_value_name=skill.proto.return_value_name,
    )
    expected_proto.parameters.Pack(
        default_params,
        type_url_prefix='type.intrinsic.ai/assets/ai.intrinsic.my_skill',
    )

    compare.assertProto2Equal(self, expected_proto, skill.proto)

  @parameterized.named_parameters(
      {
          'testcase_name': 'from_registry',
          'is_legacy_pbt': True,
      },
      {
          'testcase_name': 'from_installed_assets',
          'is_legacy_pbt': False,
      },
  )
  def test_process_skips_recommended_config(self, is_legacy_pbt):
    default_params = test_skill_params_pb2.TestMessage(my_double=2.5)
    legacy_pbts = []
    assets = []
    if is_legacy_pbt:
      legacy_pbts = [
          skill_test_utils.create_legacy_process(
              'ai.intrinsic.my_process',
              parameter_message=test_skill_params_pb2.TestMessage,
              default_params=default_params,
          )
      ]
    else:
      assets = [
          skill_test_utils.create_process_asset(
              'ai.intrinsic.my_process',
              parameter_message=test_skill_params_pb2.TestMessage,
              default_params=default_params,
          )
      ]
    skills = skill_providing.Skills(
        skill_test_utils.create_skill_registry_for_skill_infos(legacy_pbts),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets(assets),
        # Configure no recommendations for my_process. This causes an error if
        # the asset configuration client is called in any form.
        skill_test_utils.create_asset_configuration_client({
            'ai.intrinsic.my_process': [],
        }),
    )

    process = skills.ai.intrinsic.my_process()

    expected_proto = behavior_call_pb2.BehaviorCall(
        skill_id='ai.intrinsic.my_process',
        return_value_name=process.proto.return_value_name,
    )

    expected_proto.parameters.Pack(
        default_params,
        type_url_prefix=(
            f'type.intrinsic.ai/{process.info.type_url_area}/ai.intrinsic.my_process'
        ),
    )

    compare.assertProto2Equal(self, expected_proto, process.proto)

  def test_skill_recommended_config_unavailable(self):
    class FakeUnavailableGrpcError(grpc.RpcError, grpc.Call):

      def code(self) -> grpc.StatusCode:
        return grpc.StatusCode.UNAVAILABLE

      def details(self) -> str:
        return 'Asset configuration service is temporarily down.'

    stub = mock.MagicMock()
    stub.RecommendAssetConfiguration.side_effect = FakeUnavailableGrpcError()
    asset_config_client = asset_configuration_client.AssetConfigurationClient(
        stub
    )

    default_params = test_skill_params_pb2.TestMessage(my_double=2.5)
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
        default_params=default_params,
    )

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        asset_config_client,
    )

    skill = skills.ai.intrinsic.my_skill()

    stub.RecommendAssetConfiguration.assert_called_once()
    expected_proto = behavior_call_pb2.BehaviorCall(
        skill_id='ai.intrinsic.my_skill',
        return_value_name=skill.proto.return_value_name,
    )
    expected_proto.parameters.Pack(
        default_params,
        type_url_prefix='type.intrinsic.ai/assets/ai.intrinsic.my_skill',
    )

    compare.assertProto2Equal(self, expected_proto, skill.proto)

  def test_skill_recommended_config_override_precedence(self):
    default_params = test_skill_params_pb2.TestMessage(
        my_string='default-value'
    )
    recommended_params = test_skill_params_pb2.TestMessage(
        my_string='recommended-config'
    )
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
        default_params=default_params,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client({
            'ai.intrinsic.my_skill': [(default_params, recommended_params)],
        }),
    )

    # User overrides in the constructor to 'user-override'
    skill = skills.ai.intrinsic.my_skill(
        my_string='user-override',
    )

    expected_params = test_skill_params_pb2.TestMessage(
        my_string='user-override'
    )

    expected_proto = behavior_call_pb2.BehaviorCall(
        skill_id='ai.intrinsic.my_skill',
        return_value_name=skill.proto.return_value_name,
    )
    expected_proto.parameters.Pack(
        expected_params,
        type_url_prefix='type.intrinsic.ai/assets/ai.intrinsic.my_skill',
    )

    compare.assertProto2Equal(self, expected_proto, skill.proto)

  def test_skill_constructor_fails_for_unknown_arguments(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    with self.assertRaisesRegex(
        NameError, r'Unknown argument\(s\): key_that_does_not_exist'
    ):
      skills.ai.intrinsic.my_skill(key_that_does_not_exist=None)

    with self.assertRaisesRegex(
        NameError,
        r'Unknown argument\(s\): non_existent_key_with_non_none_value',
    ):
      skills.ai.intrinsic.my_skill(non_existent_key_with_non_none_value=42)

  def test_skill_signature_without_default_values(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
        return_value_message=test_skill_params_pb2.TestMessage,
    )

    # pyformat: disable
    expected_signature = (
        '(*, my_double: Union[float,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, my_float:'
        ' Union[float, intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, my_int32:'
        ' Union[int, intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, my_int64:'
        ' Union[int, intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, my_uint32:'
        ' Union[int, intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, my_uint64:'
        ' Union[int, intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, my_bool:'
        ' Union[bool, intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, my_string:'
        ' Union[str, intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, sub_message:'
        ' Union[intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.SubMessage,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' optional_sub_message:'
        ' Union[intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.SubMessage,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' my_repeated_doubles: Union[Sequence[Union[float,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression]],'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' repeated_submessages:'
        ' Union[Sequence[Union[intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.SubMessage,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression]],'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' my_required_int32: Union[int,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' my_oneof_double: Union[float,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' my_oneof_sub_message:'
        ' Union[intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.SubMessage,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, pose:'
        ' Union[intrinsic.math.python.pose3.Pose3,'
        ' intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.Pose,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, ros_pose:'
        ' Union[intrinsic.math.python.pose3.Pose3,'
        ' intrinsic.solutions.skills.ai.intrinsic.my_skill.geometry_msgs.msg.pb.jazzy.Pose,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, foo:'
        ' Union[intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.TestMessage.Foo,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, enum_v:'
        ' Union[intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.TestMessage.TestEnum,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' string_int32_map: Union[dict[str, int],'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' int32_string_map: Union[dict[int, str],'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' string_message_map: Union[dict[str,'
        ' intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.TestMessage.MessageMapValue],'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' executive_test_message:'
        ' Union[intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.executive.TestMessage,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' non_unique_field_name:'
        ' Union[intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.TestMessage.SomeType,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' return_value_key: Optional[str] = None, _with_recommended_config:'
        ' bool = True)'
    )
    # pyformat: enable

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    my_skill = skills.ai.intrinsic.my_skill()
    signature = inspect.signature(my_skill.__init__)
    self.assertSignature(signature, expected_signature)

  @parameterized.named_parameters(
      {
          'testcase_name': 'PoseSkill',
          'parameter_message': test_skill_params_pb2.PoseSkill,
          'expected_signature': (
              # pyformat: disable
              '(*, param_pose: Union[intrinsic.math.python.pose3.Pose3,'
              ' intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.Pose,'
              ' intrinsic.solutions.blackboard_value.BlackboardValue,'
              ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
              ' _with_recommended_config: bool = True)'
              # pyformat: enable
          ),
      },
      {
          'testcase_name': 'JointMotionTargetSkill',
          'parameter_message': test_skill_params_pb2.JointMotionTargetSkill,
          'expected_signature': (
              # pyformat: disable
              '(*, param_joint_motion_target:'
              ' Union[intrinsic.world.python.object_world_resources.JointConfiguration,'
              ' intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.icon.JointVec,'
              ' intrinsic.solutions.blackboard_value.BlackboardValue,'
              ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
              ' _with_recommended_config: bool = True)'
              # pyformat: enable
          ),
      },
      {
          'testcase_name': 'CollisionSettingsSkill',
          'parameter_message': test_skill_params_pb2.CollisionSettingsSkill,
          'expected_signature': (
              # pyformat: disable
              '(*, param_collision_settings:'
              ' Union[intrinsic.solutions.worlds.CollisionSettings,'
              ' intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.world.CollisionSettings,'
              ' intrinsic.solutions.blackboard_value.BlackboardValue,'
              ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
              ' _with_recommended_config: bool = True)'
              # pyformat: enable
          ),
      },
      {
          'testcase_name': 'CartesianMotionTargetSkill',
          'parameter_message': test_skill_params_pb2.CartesianMotionTargetSkill,
          'expected_signature': (
              # pyformat: disable
              '(*, target:'
              ' Union[intrinsic.solutions.worlds.CartesianMotionTarget,'
              ' intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.motion_planning.CartesianMotionTarget,'
              ' intrinsic.solutions.blackboard_value.BlackboardValue,'
              ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
              ' _with_recommended_config: bool = True)'
              # pyformat: enable
          ),
      },
      {
          'testcase_name': 'DurationSkill',
          'parameter_message': test_skill_params_pb2.DurationSkill,
          'expected_signature': (
              # pyformat: disable
              '(*, param_duration: Union[datetime.timedelta, float, int,'
              ' intrinsic.solutions.skills.ai.intrinsic.my_skill.google.protobuf.Duration,'
              ' intrinsic.solutions.blackboard_value.BlackboardValue,'
              ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
              ' _with_recommended_config: bool = True)'
              # pyformat: enable
          ),
      },
      {
          'testcase_name': 'ObjectReferenceSkill',
          'parameter_message': test_skill_params_pb2.ObjectReferenceSkill,
          'expected_signature': (
              # pyformat: disable
              '(*, param_object:'
              ' Union[intrinsic.world.python.object_world_resources.TransformNode,'
              ' intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.world.ObjectReference,'
              ' intrinsic.solutions.blackboard_value.BlackboardValue,'
              ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
              ' param_frame:'
              ' Union[intrinsic.world.python.object_world_resources.TransformNode,'
              ' intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.world.FrameReference,'
              ' intrinsic.solutions.blackboard_value.BlackboardValue,'
              ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
              ' param_transform_node:'
              ' Union[intrinsic.world.python.object_world_resources.TransformNode,'
              ' intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.world.TransformNodeReference,'
              ' intrinsic.solutions.blackboard_value.BlackboardValue,'
              ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
              ' param_object_or_entity:'
              ' Union[intrinsic.world.python.object_world_resources.WorldObject,'
              ' intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.world.ObjectOrEntityReference,'
              ' intrinsic.solutions.blackboard_value.BlackboardValue,'
              ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
              ' _with_recommended_config: bool = True)'
              # pyformat: enable
          ),
      },
      {
          'testcase_name': 'PoseEstimatorSkill',
          'parameter_message': test_skill_params_pb2.PoseEstimatorSkill,
          'expected_signature': (
              # pyformat: disable
              '(*, pose_estimator:'
              ' Union[intrinsic.solutions.pose_estimation.PoseEstimatorId,'
              ' intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.perception.v1.PoseEstimatorId,'
              ' intrinsic.solutions.blackboard_value.BlackboardValue,'
              ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
              ' _with_recommended_config: bool = True)'
              # pyformat: enable
          ),
      },
  )
  def test_skill_signature_for_types_with_auto_conversion(
      self, parameter_message, expected_signature
  ):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=parameter_message,
    )

    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    my_skill = skills.ai.intrinsic.my_skill

    signature = inspect.signature(my_skill)
    self.assertSignature(signature, expected_signature)

  def test_skill_class_docstring(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
        return_value_message=test_skill_params_pb2.TestMessage,
        description='This is an awesome skill.',
        resource_selectors={'a': ['some-type-a'], 'b': ['some-type-b']},
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    docstring = """\
Skill class for ai.intrinsic.my_skill.

This is an awesome skill."""

    self.assertEqual(skills.ai.intrinsic.my_skill.__doc__, docstring)

  def test_skill_init_docstring(self):
    asset = skill_test_utils.create_skill_asset_with_file_descriptor_set(
        'ai.intrinsic.my_skill',
        file_descriptor_set=_test_skill_params_file_descriptor_set(),
        parameter_message_full_name='intrinsic_proto.test_data.TestMessage',
        return_value_message_full_name='intrinsic_proto.test_data.TestMessage',
        description='Not to be included in init docstring',
        resource_selectors={'a': ['some-type-a'], 'b': ['some-type-b']},
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    docstring = """\
Initializes an instance of the skill ai.intrinsic.my_skill.

This method accepts the following proto messages:
  - my_skill.geometry_msgs.msg.pb.jazzy.Pose
  - my_skill.intrinsic_proto.Pose
  - my_skill.intrinsic_proto.executive.TestMessage
  - my_skill.intrinsic_proto.test_data.SubMessage
  - my_skill.intrinsic_proto.test_data.TestMessage.Foo
  - my_skill.intrinsic_proto.test_data.TestMessage.Int32StringMapEntry
  - my_skill.intrinsic_proto.test_data.TestMessage.SomeType
  - my_skill.intrinsic_proto.test_data.TestMessage.StringInt32MapEntry
  - my_skill.intrinsic_proto.test_data.TestMessage.StringMessageMapEntry

This method accepts the following proto enums:
  - my_skill.intrinsic_proto.test_data.TestMessage.TestEnum

Args:
    a:
        Resource with capability some-type-a
    b:
        Resource with capability some-type-b
    enum_v:
        enum_v comment
    executive_test_message:
        executive_test_message comment
    foo:
        foo comment
    int32_string_map:
        int32_string_map comment
    my_bool:
        my_bool comment
    my_double:
        my_double comment
    my_float:
        my_float comment
    my_int32:
        my_int32 comment
    my_int64:
        my_int64 comment
    my_oneof_double:
        my_oneof_double comment
    my_oneof_sub_message:
        my_oneof_sub_message comment
    my_repeated_doubles:
        my_repeated_doubles comment
    my_required_int32:
        my_required_int32 comment (leading). my_required_int32 comment
        (trailing).
    my_string:
        my_string comment
    my_uint32:
        my_uint32 comment
    my_uint64:
        my_uint64 comment
    non_unique_field_name:
        non_unique_field_name comment
    optional_sub_message:
        optional_sub_message comment
    pose:
        Special intrinsic-type support Handled differently than a normal
        protobuf submessage.
    repeated_submessages:
        repeated_submessages comment
    return_value_key:
        Blackboard key where to store the return value
    ros_pose:
        ros_pose comment
    string_int32_map:
        string_int32_map comment
    string_message_map:
        string_message_map comment
    sub_message:
        sub_message comment
    _with_recommended_config:
        Whether to update the parameters with the recommended configuration for
        this Skill.
        **Warning**: This field will be removed in a future release. Do not use
        this field in scripts.
        Default value: True

Returns:
    enum_v:
        enum_v comment
    executive_test_message:
        executive_test_message comment
    foo:
        foo comment
    int32_string_map:
        int32_string_map comment
    my_bool:
        my_bool comment
    my_double:
        my_double comment
    my_float:
        my_float comment
    my_int32:
        my_int32 comment
    my_int64:
        my_int64 comment
    my_oneof_double:
        my_oneof_double comment
    my_oneof_sub_message:
        my_oneof_sub_message comment
    my_repeated_doubles:
        my_repeated_doubles comment
    my_required_int32:
        my_required_int32 comment (leading). my_required_int32 comment
        (trailing).
    my_string:
        my_string comment
    my_uint32:
        my_uint32 comment
    my_uint64:
        my_uint64 comment
    non_unique_field_name:
        non_unique_field_name comment
    optional_sub_message:
        optional_sub_message comment
    pose:
        Special intrinsic-type support Handled differently than a normal
        protobuf submessage.
    repeated_submessages:
        repeated_submessages comment
    ros_pose:
        ros_pose comment
    string_int32_map:
        string_int32_map comment
    string_message_map:
        string_message_map comment
    sub_message:
        sub_message comment"""

    self.assertEqual(skills.ai.intrinsic.my_skill.__init__.__doc__, docstring)

  def test_skill_repr(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
        default_params=_DEFAULT_TEST_MESSAGE,
        resource_selectors={'a': ['some-type-a'], 'b': ['some-type-b']},
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    skill_repr = (
        'skills.ai.intrinsic.my_skill('
        'my_double=2.5, '
        'my_float=1.25, '
        'my_int32=5, '
        'my_int64=9, '
        'my_uint32=11, '
        'my_uint64=21, '
        'my_bool=True, '
        "my_string='bar', "
        'sub_message=name: "baz"\n, '
        'optional_sub_message=name: "quz"\n, '
        'my_repeated_doubles=[-5.5, 10.5], '
        'repeated_submessages=[name: "foo"\n, name: "bar"\n], '
        'my_required_int32=42, '
        'my_oneof_double=1.5, '
        'pose=position {\n}\n'
        'orientation {\n  x: 0.5\n  y: 0.5\n  z: 0.5\n  w: 0.5\n}\n, '
        'foo=bar {\n  test: "test"\n}\n, '
        'enum_v=3, '
        'executive_test_message=int32_value: 123\n, '
        'string_int32_map={"foo": 1}, '
        "int32_string_map={3: 'foobar'}, "
        'string_message_map={"bar": value: "baz"\n}, '
        'non_unique_field_name=non_unique_field_name {\n}\n, '
        'ros_pose=position {\n}\n'
        'orientation {\n  x: 0.5\n  y: 0.5\n  z: 0.5\n  w: 0.5\n}\n, '
        'a={handle: "resource_a"}, '
        'b={handle: "resource_b"})'
    )

    resource_a = provided.ResourceHandle.create('resource_a', ['some-type-a'])
    resource_b = provided.ResourceHandle.create('resource_b', ['some-type-b'])

    skill = skills.ai.intrinsic.my_skill(
        my_float=1.25, my_bool=True, a=resource_a, b=resource_b
    )
    self.assertEqual(repr(skill), skill_repr)

  def test_ambiguous_parameter_and_resource_name(self):
    """Tests ambiguous parameter name and resource slot are handled properly."""
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.ResourceConflict,
        default_params=test_skill_params_pb2.ResourceConflict(a='bar'),
        # alias chosen to match a field name from ResourceConflict
        resource_selectors={'a': ['some-type-a']},
    )
    resource_registry = skill_test_utils.create_resource_registry_with_handles([
        text_format.Parse(
            """name: 'some-resource1'
               resource_data { key: 'some-type-a' }""",
            resource_handle_pb2.ResourceHandle(),
        ),
        text_format.Parse(
            """name: 'some-resource2'
               resource_data { key: 'some-type-a' }""",
            resource_handle_pb2.ResourceHandle(),
        ),
    ])
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        resource_registry,
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    self.assertEqual(
        skills.ai.intrinsic.my_skill.__init__.__doc__,
        textwrap.dedent("""\
            Initializes an instance of the skill ai.intrinsic.my_skill.

            Args:
                a:
                a_resource:
                    Resource with capability some-type-a
                _with_recommended_config:
                    Whether to update the parameters with the recommended configuration for
                    this Skill.
                    **Warning**: This field will be removed in a future release. Do not use
                    this field in scripts.
                    Default value: True"""),
    )

    resource_a = provided.ResourceHandle.create('resource_a', ['some-type-a'])

    skill = skills.ai.intrinsic.my_skill(a='foo', a_resource=resource_a)
    self.assertEqual(
        repr(skill),
        """skills.ai.intrinsic.my_skill(a='foo', a_resource={handle: "resource_a"})""",
    )

    with self.assertRaises(TypeError):
      skills.ai.intrinsic.my_skill(a=resource_a)

    with self.assertRaisesRegex(KeyError, '.*more than one compatible.*'):
      skills.ai.intrinsic.my_skill(a='foo')

  def test_resource_default_value(self):
    """Tests default resource is used properly."""
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        resource_selectors={'a': ['some-type-a']},
    )
    resource_registry = (
        skill_test_utils.create_resource_registry_with_single_handle(
            'some-resource', 'some-type-a'
        )
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        resource_registry,
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    self.assertEqual(
        skills.ai.intrinsic.my_skill.__init__.__doc__,
        textwrap.dedent("""\
      Initializes an instance of the skill ai.intrinsic.my_skill.

      Args:
          a:
              Resource with capability some-type-a
              Default resource: some-resource
          _with_recommended_config:
              Whether to update the parameters with the recommended configuration for
              this Skill.
              **Warning**: This field will be removed in a future release. Do not use
              this field in scripts.
              Default value: True"""),
    )

    # Ensure that no "resource not found" exception is thrown
    try:
      skills.ai.intrinsic.my_skill()
    except KeyError:
      self.fail('Instantiating skill failed with default resource')

  def test_non_resource_as_resource_is_rejected(self):
    """Tests non-resource passed for resource parameter is rejected."""
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        resource_selectors={'a': ['some-type-a']},
    )
    resource_registry = (
        skill_test_utils.create_resource_registry_with_single_handle(
            'some-resource', 'some-type-a'
        )
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        resource_registry,
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    class BogusObject:

      def __init__(self):
        pass

    with self.assertRaisesRegex(TypeError, '.* not a ResourceHandle'):
      skills.ai.intrinsic.my_skill(a=BogusObject())

  def test_timeouts(self):
    """Tests if timeouts are transferred to proto."""
    asset = skill_test_utils.create_skill_asset('ai.intrinsic.my_skill')
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    skill = skills.ai.intrinsic.my_skill()
    skill.execute_timeout = datetime.timedelta(seconds=5)
    skill.project_timeout = datetime.timedelta(seconds=10)

    expected_proto = text_format.Parse(
        f"""skill_id: 'ai.intrinsic.my_skill'
            return_value_name: "{skill.proto.return_value_name}"
            skill_execution_options {{
              execute_timeout {{
                seconds: 5
              }}
              project_timeout {{
                seconds: 10
              }}
            }}
        """,
        behavior_call_pb2.BehaviorCall(),
    )
    compare.assertProto2Equal(self, skill.proto, expected_proto)

  def test_nested_message_classes(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    sub_message = (
        skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.SubMessage(
            name='nested_message_classes_test_name'
        )
    )

    skill_with_nested_class_generated_param = skills.ai.intrinsic.my_skill(
        sub_message=sub_message
    )
    action_proto = skill_with_nested_class_generated_param.proto

    test_message = test_skill_params_pb2.TestMessage()
    action_proto.parameters.Unpack(test_message)
    self.assertEqual(
        test_message.sub_message.name, sub_message.wrapped_message.name
    )

  def test_nested_message_list_with_blackboard_value(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    sub_message = (
        skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.SubMessage(
            name=cel.CelExpression('test')
        )
    )

    skill_with_nested_class_generated_param = skills.ai.intrinsic.my_skill(
        repeated_submessages=[sub_message]
    )
    action_proto = skill_with_nested_class_generated_param.proto

    test_message = test_skill_params_pb2.TestMessage()
    action_proto.parameters.Unpack(test_message)
    self.assertLen(test_message.repeated_submessages, 1)
    self.assertEqual(
        action_proto.assignments[0].parameter_path,
        'repeated_submessages[0].name',
    )
    self.assertEqual(action_proto.assignments[0].cel_expression, 'test')

  def test_skill_info_id_version_properties(self):
    info = skill_generation.SkillInfoImpl(
        id_version=id_pb2.IdVersion(
            id=id_pb2.Id(package='ai.intrinsic', name='my_skill'),
            version='0.0.1',
        ),
        description='',
        parameter_message_full_name='',
        return_value_message_full_name='',
        file_descriptor_set=descriptor_pb2.FileDescriptorSet(),
        default_params=None,
        resource_selectors={},
        proto_comments={},
        skill_type=provided.SkillType.REGULAR_SKILL,
        type_url_area='someTypeUrlArea',
    )
    self.assertEqual(info.id, 'ai.intrinsic.my_skill')
    self.assertEqual(info.id_version, 'ai.intrinsic.my_skill.0.0.1')
    self.assertEqual(info.version, '0.0.1')
    self.assertEqual(info.skill_name, 'my_skill')
    self.assertEqual(info.package_name, 'ai.intrinsic')
  def test_skill_info_invalid_package(self):
    info = skill_generation.SkillInfoImpl(
        id_version=id_pb2.IdVersion(
            id=id_pb2.Id(package='no_period', name='my_skill'),
            version='0.0.1',
        ),
        description='',
        parameter_message_full_name='',
        return_value_message_full_name='',
        file_descriptor_set=descriptor_pb2.FileDescriptorSet(),
        default_params=None,
        resource_selectors={},
        proto_comments={},
        skill_type=provided.SkillType.REGULAR_SKILL,
        type_url_area='someTypeUrlArea',
    )
    self.assertEqual(info.id, 'no_period.my_skill')
    self.assertEqual(info.id_version, 'no_period.my_skill.0.0.1')
    self.assertEqual(info.version, '0.0.1')
    self.assertEqual(info.skill_name, 'my_skill')
    self.assertEqual(info.package_name, 'no_period')
  def test_skill_info_empty_package(self):
    info = skill_generation.SkillInfoImpl(
        id_version=id_pb2.IdVersion(
            id=id_pb2.Id(package='', name='my_skill'),
            version='0.0.1',
        ),
        description='',
        parameter_message_full_name='',
        return_value_message_full_name='',
        file_descriptor_set=descriptor_pb2.FileDescriptorSet(),
        default_params=None,
        resource_selectors={},
        proto_comments={},
        skill_type=provided.SkillType.REGULAR_SKILL,
        type_url_area='someTypeUrlArea',
    )
    self.assertEqual(info.id, 'my_skill')
    self.assertEqual(info.id_version, 'my_skill.0.0.1')
    self.assertEqual(info.version, '0.0.1')
    self.assertEqual(info.skill_name, 'my_skill')
    self.assertEqual(info.package_name, '')

  def test_skill_info_parameter_and_return_value_properties(self):
    file_descriptor_set = descriptor_pb2.FileDescriptorSet(
        file=[
            descriptor_pb2.FileDescriptorProto(
                name='my_skill.proto',
                package='my_package',
                message_type=[
                    descriptor_pb2.DescriptorProto(name='MySkillParams'),
                    descriptor_pb2.DescriptorProto(name='MySkillResult'),
                ],
            )
        ]
    )
    defaults_any = any_pb2.Any(
        type_url='type.intrinsic.ai/assets/ai.intrinsic.my_skill/my_package.MySkillParams'
    )

    info = skill_generation.SkillInfoImpl(
        id_version=id_pb2.IdVersion(
            id=id_pb2.Id(package='ai.intrinsic', name='my_skill'),
            version='0.0.1',
        ),
        description='',
        parameter_message_full_name='my_package.MySkillParams',
        return_value_message_full_name='my_package.MySkillResult',
        file_descriptor_set=file_descriptor_set,
        default_params=defaults_any,
        resource_selectors={},
        proto_comments={},
        skill_type=provided.SkillType.REGULAR_SKILL,
        type_url_area='someTypeUrlArea',
    )

    self.assertEqual(
        info.parameter_message_full_name, 'my_package.MySkillParams'
    )
    self.assertEqual(
        info.return_value_message_full_name, 'my_package.MySkillResult'
    )
    self.assertEqual(info.file_descriptor_set, file_descriptor_set)
    self.assertEqual(info.default_params, defaults_any)

  def test_construct_skill_info_incomplete_fileset(self):
    file_descriptor_set = descriptor_pb2.FileDescriptorSet()
    test_skill_params_pb2.TestMessage.DESCRIPTOR.file.CopyToProto(
        file_descriptor_set.file.add()
    )

    with self.assertRaises(TypeError):
      skill_generation.SkillInfoImpl(
          id_version=id_pb2.IdVersion(
              id=id_pb2.Id(package='ai.intrinsic', name='my_skill'),
              version='0.0.1',
          ),
          description='',
          parameter_message_full_name='',
          return_value_message_full_name='',
          file_descriptor_set=file_descriptor_set,
          default_params=None,
          resource_selectors={},
          proto_comments={},
          skill_type=provided.SkillType.REGULAR_SKILL,
          type_url_area='someTypeUrlArea',
      )

  def test_skill_info_resource_selectors(self):
    resource_selectors = {
        'robot': equipment_pb2.ResourceSelector(capability_names=['IconApi'])
    }
    info = skill_generation.SkillInfoImpl(
        id_version=id_pb2.IdVersion(
            id=id_pb2.Id(package='ai.intrinsic', name='my_skill'),
            version='0.0.1',
        ),
        description='',
        parameter_message_full_name='',
        return_value_message_full_name='',
        file_descriptor_set=descriptor_pb2.FileDescriptorSet(),
        default_params=None,
        resource_selectors=resource_selectors,
        proto_comments={},
        skill_type=provided.SkillType.REGULAR_SKILL,
        type_url_area='someTypeUrlArea',
    )

    self.assertEqual(info.resource_selectors, resource_selectors)

  def test_skill_info_type_and_type_url(self):
    info = skill_generation.SkillInfoImpl(
        id_version=id_pb2.IdVersion(
            id=id_pb2.Id(package='ai.intrinsic', name='my_skill'),
            version='0.0.1',
        ),
        description='',
        parameter_message_full_name='',
        return_value_message_full_name='',
        file_descriptor_set=descriptor_pb2.FileDescriptorSet(),
        default_params=None,
        resource_selectors={},
        proto_comments={},
        skill_type=provided.SkillType.REGULAR_SKILL,
        type_url_area='someTypeUrlArea',
    )

    self.assertEqual(info.skill_type, provided.SkillType.REGULAR_SKILL)
    self.assertEqual(info.type_url_area, 'someTypeUrlArea')

  def test_skill_info_proto_comments(self):
    info = skill_generation.SkillInfoImpl(
        id_version=id_pb2.IdVersion(
            id=id_pb2.Id(package='ai.intrinsic', name='my_skill'),
            version='0.0.1',
        ),
        description='',
        parameter_message_full_name='',
        return_value_message_full_name='',
        file_descriptor_set=descriptor_pb2.FileDescriptorSet(),
        default_params=None,
        resource_selectors={},
        proto_comments={
            'intrinsic_proto.MyMessage': 'MyMessage comment\n',
            'intrinsic_proto.MyMessage.my_field': 'my_field comment\n',
        },
        skill_type=provided.SkillType.REGULAR_SKILL,
        type_url_area='someTypeUrlArea',
    )

    self.assertEqual(
        info.get_proto_comment('intrinsic_proto.MyMessage'),
        'MyMessage comment\n',
    )
    self.assertEqual(
        info.get_proto_comment('intrinsic_proto.MyMessage.my_field'),
        'my_field comment\n',
    )
    self.assertEqual(info.get_proto_comment('non_existing.Name'), '')

  def test_result_access(self):
    """Tests if BlackboardValue gets created when accessing result."""
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        return_value_message=test_skill_params_pb2.TestMessage,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    skill = skills.ai.intrinsic.my_skill()
    self.assertIsInstance(skill.result, blackboard_value.BlackboardValue)
    self.assertContainsSubset(
        _DEFAULT_TEST_MESSAGE.DESCRIPTOR.fields_by_name.keys(),
        dir(skill.result),
        'Missing attributes in BlackboardValue',
    )
    self.assertEqual(skill.result.value_access_path(), skill.result_key)

  def test_gen_message_wrapper(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    test_message = (
        skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.SubMessage(
            name='bar'
        )
    )

    self.assertEqual('name: "bar"\n', str(test_message.wrapped_message))

  def test_wrapper_class_name(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    my_skill = skills.ai.intrinsic.my_skill

    # test_data is the subpackage of the proto file and should be represented by
    # a simple namespace class.
    self.assertEqual(my_skill.intrinsic_proto.test_data.__name__, 'test_data')
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.__qualname__,
        'my_skill.intrinsic_proto.test_data',
    )
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.__module__,
        'intrinsic.solutions.skills.ai.intrinsic',
    )

    # SubMessage is a top-level message in proto file.
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.SubMessage.__name__, 'SubMessage'
    )
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.SubMessage.__qualname__,
        'my_skill.intrinsic_proto.test_data.SubMessage',
    )
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.SubMessage.__module__,
        'intrinsic.solutions.skills.ai.intrinsic',
    )

    # Foo (=TestMessage.Foo) is a nested message in proto file.
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.TestMessage.Foo.__name__, 'Foo'
    )
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.TestMessage.Foo.__qualname__,
        'my_skill.intrinsic_proto.test_data.TestMessage.Foo',
    )
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.TestMessage.Foo.__module__,
        'intrinsic.solutions.skills.ai.intrinsic',
    )

  def test_wrapper_access(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
        return_value_message=test_skill_params_pb2.TestMessageReturn,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    my_skill = skills.ai.intrinsic.my_skill

    # Id notation: skill.<full message name>
    self.assertIsInstance(
        my_skill.intrinsic_proto.test_data.TestMessageReturn(),
        skill_utils.MessageWrapper,
    )
    self.assertIsInstance(
        my_skill.intrinsic_proto.test_data.ReturnValue(),
        skill_utils.MessageWrapper,
    )
    self.assertIsInstance(
        my_skill.intrinsic_proto.test_data.SubMessage(),
        skill_utils.MessageWrapper,
    )
    self.assertIsInstance(
        my_skill.intrinsic_proto.test_data.TestMessage(),
        skill_utils.MessageWrapper,
    )
    self.assertIsInstance(
        my_skill.intrinsic_proto.test_data.TestMessage.Foo(),
        skill_utils.MessageWrapper,
    )
    self.assertIsInstance(
        my_skill.intrinsic_proto.test_data.TestMessage.Foo.Bar(),
        skill_utils.MessageWrapper,
    )
    self.assertIsInstance(
        my_skill.intrinsic_proto.executive.TestMessage(),
        skill_utils.MessageWrapper,
    )
    self.assertIsInstance(
        my_skill.google.protobuf.Any(),
        skill_utils.MessageWrapper,
    )

  def test_message_wrapper_class_docstring(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    self.assertEqual(
        skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.TestMessage.__doc__,
        'Proto message wrapper class for'
        ' intrinsic_proto.test_data.TestMessage.',
    )

  def test_message_wrapper_init_docstring(self):
    asset = skill_test_utils.create_skill_asset_with_file_descriptor_set(
        'ai.intrinsic.my_skill',
        file_descriptor_set=_test_skill_params_file_descriptor_set(),
        parameter_message_full_name='intrinsic_proto.test_data.TestMessage',
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    docstring = """\
Initializes an instance of my_skill.intrinsic_proto.test_data.TestMessage.

This method accepts the following proto messages:
  - my_skill.geometry_msgs.msg.pb.jazzy.Pose
  - my_skill.intrinsic_proto.Pose
  - my_skill.intrinsic_proto.executive.TestMessage
  - my_skill.intrinsic_proto.test_data.SubMessage
  - my_skill.intrinsic_proto.test_data.TestMessage.Foo
  - my_skill.intrinsic_proto.test_data.TestMessage.Int32StringMapEntry
  - my_skill.intrinsic_proto.test_data.TestMessage.SomeType
  - my_skill.intrinsic_proto.test_data.TestMessage.StringInt32MapEntry
  - my_skill.intrinsic_proto.test_data.TestMessage.StringMessageMapEntry

This method accepts the following proto enums:
  - my_skill.intrinsic_proto.test_data.TestMessage.TestEnum

Fields:
    enum_v:
        enum_v comment
    executive_test_message:
        executive_test_message comment
    foo:
        foo comment
    int32_string_map:
        int32_string_map comment
    my_bool:
        my_bool comment
    my_double:
        my_double comment
    my_float:
        my_float comment
    my_int32:
        my_int32 comment
    my_int64:
        my_int64 comment
    my_oneof_double:
        my_oneof_double comment
    my_oneof_sub_message:
        my_oneof_sub_message comment
    my_repeated_doubles:
        my_repeated_doubles comment
    my_required_int32:
        my_required_int32 comment (leading). my_required_int32 comment
        (trailing).
    my_string:
        my_string comment
    my_uint32:
        my_uint32 comment
    my_uint64:
        my_uint64 comment
    non_unique_field_name:
        non_unique_field_name comment
    optional_sub_message:
        optional_sub_message comment
    pose:
        Special intrinsic-type support Handled differently than a normal
        protobuf submessage.
    repeated_submessages:
        repeated_submessages comment
    ros_pose:
        ros_pose comment
    string_int32_map:
        string_int32_map comment
    string_message_map:
        string_message_map comment
    sub_message:
        sub_message comment"""

    my_skill = skills.ai.intrinsic.my_skill

    self.assertEqual(
        my_skill.intrinsic_proto.test_data.TestMessage.__init__.__doc__,
        docstring,
    )

  def test_message_wrapper_signature(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    # pyformat: disable
    expected_signature = (
        '(*, my_double: Union[float,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, my_float:'
        ' Union[float, intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, my_int32:'
        ' Union[int, intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, my_int64:'
        ' Union[int, intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, my_uint32:'
        ' Union[int, intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, my_uint64:'
        ' Union[int, intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, my_bool:'
        ' Union[bool, intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, my_string:'
        ' Union[str, intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, sub_message:'
        ' Union[intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.SubMessage,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' optional_sub_message:'
        ' Union[intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.SubMessage,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' my_repeated_doubles: Union[Sequence[Union[float,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression]],'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' repeated_submessages:'
        ' Union[Sequence[Union[intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.SubMessage,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression]],'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' my_required_int32: Union[int,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' my_oneof_double: Union[float,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' my_oneof_sub_message:'
        ' Union[intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.SubMessage,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, pose:'
        ' Union[intrinsic.math.python.pose3.Pose3,'
        ' intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.Pose,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, ros_pose:'
        ' Union[intrinsic.math.python.pose3.Pose3,'
        ' intrinsic.solutions.skills.ai.intrinsic.my_skill.geometry_msgs.msg.pb.jazzy.Pose,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, foo:'
        ' Union[intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.TestMessage.Foo,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None, enum_v:'
        ' Union[intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.TestMessage.TestEnum,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' string_int32_map: Union[dict[str, int],'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' int32_string_map: Union[dict[int, str],'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' string_message_map: Union[dict[str,'
        ' intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.TestMessage.MessageMapValue],'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' executive_test_message:'
        ' Union[intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.executive.TestMessage,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None,'
        ' non_unique_field_name:'
        ' Union[intrinsic.solutions.skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.TestMessage.SomeType,'
        ' intrinsic.solutions.blackboard_value.BlackboardValue,'
        ' intrinsic.solutions.cel.CelExpression, NoneType] = None)'
    )
    # pyformat: enable

    my_skill = skills.ai.intrinsic.my_skill
    signature = inspect.signature(
        my_skill.intrinsic_proto.test_data.TestMessage
    )
    self.assertSignature(signature, expected_signature)

  def test_message_wrapper_explicit_none(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    my_skill = skills.ai.intrinsic.my_skill

    m = my_skill.intrinsic_proto.test_data.TestMessage(
        sub_message=None,
        my_repeated_doubles=None,
        string_int32_map=None,
    )

    expected_proto = test_skill_params_pb2.TestMessage()
    compare.assertProto2Equal(self, expected_proto, m.wrapped_message)

  def test_message_wrapper_params(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    expected_test_message = test_skill_params_pb2.TestMessage(
        my_double=2,
        my_float=-1.8,
        my_int32=5,
        my_uint32=11,
        my_bool=True,
        my_string='bar',
        repeated_submessages=[
            test_skill_params_pb2.SubMessage(),
            test_skill_params_pb2.SubMessage(name='foo'),
            test_skill_params_pb2.SubMessage(),
            test_skill_params_pb2.SubMessage(),
        ],
    )

    test_message_wrapper = (
        skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.TestMessage(
            my_double=2.0,
            my_float=-1.8,
            my_int32=5,
            my_int64=blackboard_value.BlackboardValue(
                {}, 'my_int64', None, None
            ),
            my_uint32=11,
            my_uint64=cel.CelExpression('my_uint64'),
            my_bool=True,
            my_string='bar',
            repeated_submessages=[
                blackboard_value.BlackboardValue({}, 'test', None, None),
                test_skill_params_pb2.SubMessage(name='foo'),
                blackboard_value.BlackboardValue({}, 'bar', None, None),
                cel.CelExpression('fax'),
            ],
        )
    )

    self.assertContainsSubset(
        {
            'repeated_submessages[0]': 'test',
            'repeated_submessages[2]': 'bar',
            'repeated_submessages[3]': 'fax',
            'my_int64': 'my_int64',
            'my_uint64': 'my_uint64',
        },
        test_message_wrapper.blackboard_params,
    )
    compare.assertProto2Equal(
        self, test_message_wrapper.wrapped_message, expected_test_message
    )

  def test_message_wrapper_attribute_assignment(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )
    expected_test_message = test_skill_params_pb2.TestMessage(
        my_double=2,
        my_float=-1.8,
        my_int32=5,
        my_uint32=11,
        my_bool=True,
        my_string='bar',
        repeated_submessages=[
            test_skill_params_pb2.SubMessage(),
            test_skill_params_pb2.SubMessage(name='foo'),
            test_skill_params_pb2.SubMessage(),
            test_skill_params_pb2.SubMessage(),
        ],
    )

    test_message_wrapper = (
        skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.TestMessage()
    )
    test_message_wrapper.my_double = 2.0
    test_message_wrapper.my_float = -1.8
    test_message_wrapper.my_int32 = 5
    test_message_wrapper.my_int64 = blackboard_value.BlackboardValue(
        {}, 'my_int64', None, None
    )
    test_message_wrapper.my_uint32 = 11
    test_message_wrapper.my_uint64 = cel.CelExpression('my_uint64')
    test_message_wrapper.my_bool = True
    test_message_wrapper.my_string = 'bar'
    test_message_wrapper.repeated_submessages = [
        blackboard_value.BlackboardValue({}, 'test', None, None),
        test_skill_params_pb2.SubMessage(name='foo'),
        blackboard_value.BlackboardValue({}, 'bar', None, None),
        cel.CelExpression('fax'),
    ]

    self.assertContainsSubset(
        {
            'repeated_submessages[0]': 'test',
            'repeated_submessages[2]': 'bar',
            'repeated_submessages[3]': 'fax',
            'my_int64': 'my_int64',
            'my_uint64': 'my_uint64',
        },
        test_message_wrapper.blackboard_params,
    )
    compare.assertProto2Equal(
        self, test_message_wrapper.wrapped_message, expected_test_message
    )

  def test_message_wrapper_attribute_assignment_raises(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    test_message_wrapper = (
        skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.TestMessage()
    )
    with self.assertRaisesRegex(AttributeError, 'not_a_field'):
      test_message_wrapper.not_a_field = 2.0

    test_message_wrapper = (
        skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.TestMessage()
    )
    with self.assertRaisesRegex(AttributeError, 'not_a_field'):
      test_message_wrapper.not_a_field = cel.CelExpression('some_key')

  def test_message_wrapper_to_any(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.TestMessage,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    expected_test_message = test_skill_params_pb2.TestMessage(
        my_int32=5,
        repeated_submessages=[
            test_skill_params_pb2.SubMessage(),
            test_skill_params_pb2.SubMessage(name='foo'),
        ],
    )
    test_message_wrapper = (
        skills.ai.intrinsic.my_skill.intrinsic_proto.test_data.TestMessage(
            my_int32=5,
            repeated_submessages=[
                test_skill_params_pb2.SubMessage(),
                test_skill_params_pb2.SubMessage(name='foo'),
            ],
        )
    )

    any_msg = test_message_wrapper.to_any()
    unpacked_msg = test_skill_params_pb2.TestMessage()
    any_msg.Unpack(unpacked_msg)

    compare.assertProto2Equal(self, unpacked_msg, expected_test_message)
    self.assertEqual(
        any_msg.type_url,
        'type.intrinsic.ai/assets/ai.intrinsic.my_skill/intrinsic_proto.test_data.TestMessage',
    )

  def test_enum_wrapper_class(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.VariousEnumsMessage,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    my_skill = skills.ai.intrinsic.my_skill

    global_enum = my_skill.intrinsic_proto.test_data.GlobalEnum
    self.assertTrue(issubclass(global_enum, enum.IntEnum))
    self.assertEqual(global_enum.__name__, 'GlobalEnum')
    self.assertEqual(
        global_enum.__qualname__,
        'my_skill.intrinsic_proto.test_data.GlobalEnum',
    )
    self.assertEqual(
        global_enum.__module__,
        'intrinsic.solutions.skills.ai.intrinsic',
    )

    test_enum = my_skill.intrinsic_proto.test_data.TestMessage.TestEnum
    self.assertTrue(issubclass(test_enum, enum.IntEnum))
    self.assertEqual(test_enum.__name__, 'TestEnum')
    self.assertEqual(
        test_enum.__qualname__,
        'my_skill.intrinsic_proto.test_data.TestMessage.TestEnum',
    )
    self.assertEqual(
        test_enum.__module__,
        'intrinsic.solutions.skills.ai.intrinsic',
    )

  def test_enum_values(self):
    asset = skill_test_utils.create_skill_asset(
        'ai.intrinsic.my_skill',
        parameter_message=test_skill_params_pb2.VariousEnumsMessage,
    )
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets([asset]),
        skill_test_utils.create_asset_configuration_client(),
    )

    my_skill = skills.ai.intrinsic.my_skill

    # Test global enum.
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.GlobalEnum.GLOBAL_ENUM_UNSPECIFIED, 0
    )
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.GlobalEnum.GLOBAL_ENUM_ONE, 1
    )
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.GlobalEnum.GLOBAL_ENUM_TWO, 2
    )

    # Test global enum shortcuts.
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.GLOBAL_ENUM_UNSPECIFIED, 0
    )
    self.assertEqual(my_skill.intrinsic_proto.test_data.GLOBAL_ENUM_ONE, 1)
    self.assertEqual(my_skill.intrinsic_proto.test_data.GLOBAL_ENUM_TWO, 2)

    # Test enum which is nested in a message.
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.TestMessage.TestEnum.UNKNOWN, 0
    )
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.TestMessage.TestEnum.ONE, 1
    )
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.TestMessage.TestEnum.THREE, 3
    )
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.TestMessage.TestEnum.FIVE, 5
    )

    # Test shortcuts for an enum which is nested in a message.
    self.assertEqual(my_skill.intrinsic_proto.test_data.TestMessage.UNKNOWN, 0)
    self.assertEqual(my_skill.intrinsic_proto.test_data.TestMessage.ONE, 1)
    self.assertEqual(my_skill.intrinsic_proto.test_data.TestMessage.THREE, 3)
    self.assertEqual(my_skill.intrinsic_proto.test_data.TestMessage.FIVE, 5)

    # Test enum which is nested in the skills param message.
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.VariousEnumsMessage.VariousEnumsEnum.VARIOUS_ENUMS_ENUM_UNSPECIFIED,
        0,
    )
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.VariousEnumsMessage.VariousEnumsEnum.VARIOUS_ENUMS_ENUM_ONE,
        1,
    )
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.VariousEnumsMessage.VariousEnumsEnum.VARIOUS_ENUMS_ENUM_TWO,
        2,
    )

    # Test shortcuts for enum which is nested in the skills param message.
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.VariousEnumsMessage.VARIOUS_ENUMS_ENUM_UNSPECIFIED,
        0,
    )
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.VariousEnumsMessage.VARIOUS_ENUMS_ENUM_ONE,
        1,
    )
    self.assertEqual(
        my_skill.intrinsic_proto.test_data.VariousEnumsMessage.VARIOUS_ENUMS_ENUM_TWO,
        2,
    )

    # Test special shortcuts directly on skill class for enum which is nested in
    # the skills param message.
    self.assertEqual(my_skill.VARIOUS_ENUMS_ENUM_UNSPECIFIED, 0)
    self.assertEqual(my_skill.VARIOUS_ENUMS_ENUM_ONE, 1)
    self.assertEqual(my_skill.VARIOUS_ENUMS_ENUM_TWO, 2)

  def test_skills_len(self):
    assets = [
        skill_test_utils.create_skill_asset('ai.intr.intr_skill_one'),
        skill_test_utils.create_skill_asset('ai.intr.intr_skill_two'),
        skill_test_utils.create_process_asset('ai.intr.intr_process'),
    ]
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets(assets),
        skill_test_utils.create_asset_configuration_client(),
    )

    self.assertLen(skills, 3)

  def test_skills_contains(self):
    assets = [
        skill_test_utils.create_skill_asset('ai.intr.intr_skill_one'),
        skill_test_utils.create_skill_asset('ai.intr.intr_skill_two'),
        skill_test_utils.create_process_asset('ai.intr.intr_process'),
    ]
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets(assets),
        skill_test_utils.create_asset_configuration_client(),
    )

    self.assertIn('ai.intr.intr_skill_one', skills)
    self.assertIn('ai.intr.intr_skill_two', skills)
    self.assertIn('ai.intr.intr_process', skills)

  def test_skills_iter(self):
    assets = [
        skill_test_utils.create_skill_asset('ai.intr.intr_skill_one'),
        skill_test_utils.create_skill_asset('ai.intr.intr_skill_two'),
        skill_test_utils.create_process_asset('ai.intr.intr_process'),
    ]
    skills = skill_providing.Skills(
        skill_test_utils.create_empty_skill_registry(),
        skill_test_utils.create_empty_resource_registry(),
        skill_test_utils.create_installed_assets(assets),
        skill_test_utils.create_asset_configuration_client(),
    )

    iterated_skills = []
    for skill in skills:
      iterated_skills.append(skill)

    self.assertEqual(
        iterated_skills,
        [
            skills['ai.intr.intr_process'],
            skills['ai.intr.intr_skill_one'],
            skills['ai.intr.intr_skill_two'],
        ],
    )


if __name__ == '__main__':
  absltest.main()
