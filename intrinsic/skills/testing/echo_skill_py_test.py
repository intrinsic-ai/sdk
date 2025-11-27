# Copyright 2023 Intrinsic Innovation LLC

from absl.testing import absltest

from intrinsic.skills.testing import echo_skill
from intrinsic.skills.testing import echo_skill_pb2
from intrinsic.skills.testing import skill_test_utils

_MANIFEST_PATH = "intrinsic/skills/testing/echo_skill_py_manifest.pbbin"


class EchoSkillTest(absltest.TestCase):

  @classmethod
  def setUpClass(cls):
    super().setUpClass()
    cls._manifest = skill_test_utils.get_skill_manifest(_MANIFEST_PATH)

  def test_execute_echos_params(self):
    foo = "bunny"
    echo = echo_skill.EchoSkill()

    result = echo.execute(
        request=skill_test_utils.make_test_execute_request(
            params=echo_skill_pb2.EchoSkillParams(foo=foo)
        ),
        context=skill_test_utils.make_test_execute_context(),
    )

    self.assertEqual(result.foo, foo)

  def test_preview_echos_params(self):
    foo = "fighter"
    echo = echo_skill.EchoSkill()

    result = echo.preview(
        request=skill_test_utils.make_test_preview_request(
            params=echo_skill_pb2.EchoSkillParams(foo=foo)
        ),
        context=skill_test_utils.make_test_preview_context(),
    )

    self.assertEqual(result.foo, foo)

  def test_no_required_equipment(self):
    self.assertEmpty(self._manifest.dependencies.required_equipment)

  def test_name(self):
    self.assertEqual(self._manifest.id.name, "echo")

  def test_default_parameters_have_correct_type(self):
    self.assertEqual(
        self._manifest.parameter.default_value.type_url,
        "type.googleapis.com/intrinsic_proto.skills.EchoSkillParams",
    )

  def test_correct_parameter_descriptor(self):
    self.assertEqual(
        self._manifest.parameter.message_full_name,
        echo_skill_pb2.EchoSkillParams.DESCRIPTOR.full_name,
    )

  def test_correct_return_value_descriptor(self):
    self.assertEqual(
        self._manifest.return_type.message_full_name,
        echo_skill_pb2.EchoSkillReturn.DESCRIPTOR.full_name,
    )


if __name__ == "__main__":
  absltest.main()
