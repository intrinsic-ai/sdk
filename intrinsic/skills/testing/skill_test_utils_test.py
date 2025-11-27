# Copyright 2023 Intrinsic Innovation LLC

from absl.testing import absltest
from absl.testing import parameterized
import grpc

from intrinsic.assets import id_utils
from intrinsic.assets.services.examples.calcserver import calc_server
from intrinsic.assets.services.examples.calcserver import calc_server_pb2
from intrinsic.assets.services.examples.calcserver import calc_server_pb2_grpc
from intrinsic.logging.proto import log_item_pb2

# isort: off
# isort: on
from intrinsic.skills.proto import skill_manifest_pb2
from intrinsic.skills.python import skill_logging_context
from intrinsic.skills.testing import echo_skill_pb2
from intrinsic.skills.testing import skill_test_utils
from intrinsic.solutions.testing import compare

_MANIFEST_PATH = 'intrinsic/skills/testing/echo_skill_py_manifest.pbbin'


class SkillTestUtilsTest(parameterized.TestCase):

  def test_make_grpc_server_with_channel(self):
    server, channel = skill_test_utils.make_grpc_server_with_channel()
    try:
      servicer = calc_server.CalculatorServiceServicer(
          calc_server_pb2.CalculatorConfig()
      )
      calc_server_pb2_grpc.add_CalculatorServicer_to_server(servicer, server)
      server.start()
      stub = calc_server_pb2_grpc.CalculatorStub(channel)
      request = calc_server_pb2.CalculatorRequest(
          operation=calc_server_pb2.CALCULATOR_OPERATION_ADD, x=4, y=2
      )
      response = stub.Calculate(request)
      self.assertEqual(response.result, 6)
    finally:
      server.stop(None)

  def test_make_grpc_server_with_resource_handle(self):
    server, handle = skill_test_utils.make_grpc_server_with_resource_handle(
        'bar'
    )
    try:
      servicer = calc_server.CalculatorServiceServicer(
          calc_server_pb2.CalculatorConfig()
      )
      calc_server_pb2_grpc.add_CalculatorServicer_to_server(servicer, server)
      server.start()
      channel = grpc.insecure_channel(handle.connection_info.grpc.address)
      stub = calc_server_pb2_grpc.CalculatorStub(channel)
      request = calc_server_pb2.CalculatorRequest(
          operation=calc_server_pb2.CALCULATOR_OPERATION_ADD, x=1, y=9
      )
      response = stub.Calculate(request)
      self.assertEqual(response.result, 10)
      self.assertEqual(handle.name, 'bar')
    finally:
      server.stop(None)

  def test_get_skill_manifest(self):
    manifest = skill_test_utils.get_skill_manifest(_MANIFEST_PATH)
    self.assertIsInstance(manifest, skill_manifest_pb2.SkillManifest)
    self.assertEqual(manifest.id.name, 'echo')
    self.assertEqual(
        manifest.options.python_config.create_skill,
        'intrinsic.skills.testing.echo_skill.EchoSkill',
    )
    self.assertEqual(
        manifest.parameter.message_full_name,
        'intrinsic_proto.skills.EchoSkillParams',
    )
    self.assertEqual(
        manifest.return_type.message_full_name,
        'intrinsic_proto.skills.EchoSkillReturn',
    )

  def test_make_skill_logging_context_from_manifest(self):
    manifest = skill_test_utils.get_skill_manifest(_MANIFEST_PATH)
    logging_context = (
        skill_test_utils.make_test_skill_logging_context_from_manifest(manifest)
    )
    self.assertIsInstance(
        logging_context, skill_logging_context.SkillLoggingContext
    )
    self.assertEqual(
        logging_context.skill_id,
        id_utils.id_from(manifest.id.package, manifest.id.name),
    )

  @parameterized.named_parameters([
      dict(
          testcase_name='load_manifest',
          manifest=skill_test_utils.get_skill_manifest(_MANIFEST_PATH),
      )
  ])
  def test_get_skill_manifest_works_before_parsing_flags(
      self, manifest: skill_manifest_pb2.SkillManifest
  ):
    """Tests that get_skill_manifest works before parsing flags.

    parameterized.named_parameters is executed before absltest.main(). This
    test checks that get_skill_manifest can be called before absltest.main().

    Args:
      manifest: The manifest to test.
    """
    compare.assertProto2Equal(
        self, manifest, skill_test_utils.get_skill_manifest(_MANIFEST_PATH)
    )


if __name__ == '__main__':
  absltest.main()
