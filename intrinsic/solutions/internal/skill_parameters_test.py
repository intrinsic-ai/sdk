# Copyright 2023 Intrinsic Innovation LLC

from absl.testing import absltest
from absl.testing import parameterized

from intrinsic.solutions.internal import skill_parameters
from intrinsic.solutions.testing import skill_test_utils
from intrinsic.solutions.testing import test_skill_params_pb2
from intrinsic.util.path_resolver import path_resolver

_TEST_MESSAGE = test_skill_params_pb2.TestMessage()


class SkillParametersTest(parameterized.TestCase):

  def setUp(self):
    super().setUp()

    self._utils = skill_test_utils.SkillTestUtils(
        path_resolver.resolve_runfiles_path(
            'intrinsic/solutions/testing/test_skill_params_proto_descriptors_transitive_set_sci.proto.bin'
        )
    )

  @parameterized.named_parameters(
      (
          'primitive_field',
          'my_required_int32',
          _TEST_MESSAGE,
          False,
      ),
      (
          'optional_primitive_field',
          'my_double',
          _TEST_MESSAGE,
          True,
      ),
      (
          'message_field',
          'sub_message',
          _TEST_MESSAGE,
          False,
      ),
      (
          'optional_message_field',
          'optional_sub_message',
          _TEST_MESSAGE,
          True,
      ),
      (
          'oneof_field',
          'my_oneof_double',
          _TEST_MESSAGE,
          True,
      ),
      (
          'repeated_field',
          'my_repeated_doubles',
          _TEST_MESSAGE,
          False,
      ),
      (
          'map_field',
          'int32_string_map',
          _TEST_MESSAGE,
          False,
      ),
  )
  def test_is_optional_in_python_signature(
      self, field_name, test_message, expected_result
  ):
    skill_info = self._utils.create_test_skill_info(
        skill_id='ai.intrinsic.my_skill',
        parameter_defaults=test_message,
    )
    skill_params = skill_parameters.SkillParameters(
        msg_descriptor=test_message.DESCRIPTOR,
        file_descriptor_set=(
            skill_info.parameter_description.parameter_descriptor_fileset
        ),
    )

    self.assertEqual(
        skill_params.is_optional_in_python_signature(field_name),
        expected_result,
    )


if __name__ == '__main__':
  absltest.main()
