# Copyright 2023 Intrinsic Innovation LLC

"""Tests for code_execution_client."""

from unittest import mock

from absl.testing import absltest
from google.protobuf import descriptor_pb2
import grpc

from intrinsic.executive.proto import behavior_call_pb2
from intrinsic.executive.proto import code_execution_info_service_pb2
from intrinsic.executive.proto import code_execution_pb2
from intrinsic.solutions import behavior_tree
from intrinsic.solutions import proto_building
from intrinsic.solutions.internal import code_execution_client


def _create_python_script(*, function_body: str) -> behavior_tree.PythonScript:
  fds = descriptor_pb2.FileDescriptorSet()
  file = fds.file.add(name="node.proto", package="my_pkg")
  file.message_type.add(name="Params")
  file.message_type.add(name="ReturnValue")

  sig = proto_building.Signature(
      parameter_message_full_name="my_pkg.Params",
      return_value_message_full_name="my_pkg.ReturnValue",
      file_descriptor_set=fds,
  )
  return behavior_tree.PythonScript(
      signature_with_args=sig.with_args(),
      function_body=function_body,
  )


class CodeExecutionClientTest(absltest.TestCase):

  def setUp(self):
    super().setUp()
    self._code_execution_info_service = mock.MagicMock()
    self._code_execution_info = code_execution_client.InfoClient(
        self._code_execution_info_service
    )

  def test_connect(self):
    channel = mock.MagicMock(spec=grpc.Channel)
    client = code_execution_client.InfoClient.connect(channel)
    self.assertIsInstance(client, code_execution_client.InfoClient)

  def test_preview_executed_code_with_python_script(self):
    script = _create_python_script(
        function_body="  return node_pb2.ReturnValue()"
    )
    sig = script.signature_with_args
    self._code_execution_info_service.GetPythonTemplate.return_value = (
        code_execution_info_service_pb2.GetPythonTemplateResponse(
            code_template="def compute():\n  return node_pb2.ReturnValue()"
        )
    )

    preview = self._code_execution_info.preview_executed_code(script)

    self.assertEqual(preview, "def compute():\n  return node_pb2.ReturnValue()")
    self._code_execution_info_service.GetPythonTemplate.assert_called_once_with(
        code_execution_info_service_pb2.GetPythonTemplateRequest(
            parameter_message_full_name=sig.parameter_message_full_name,
            return_value_message_full_name=sig.return_value_message_full_name,
            file_descriptor_set=sig.file_descriptor_set,
            python_code=code_execution_pb2.PythonCode(
                function_body="  return node_pb2.ReturnValue()"
            ),
        )
    )

  def test_preview_executed_code_with_task(self):
    script = _create_python_script(
        function_body="  return node_pb2.ReturnValue()"
    )
    sig = script.signature_with_args
    task = behavior_tree.Task(script)
    self._code_execution_info_service.GetPythonTemplate.return_value = (
        code_execution_info_service_pb2.GetPythonTemplateResponse(
            code_template="def compute():\n  return node_pb2.ReturnValue()"
        )
    )

    preview = self._code_execution_info.preview_executed_code(task)

    self.assertEqual(preview, "def compute():\n  return node_pb2.ReturnValue()")
    self._code_execution_info_service.GetPythonTemplate.assert_called_once_with(
        code_execution_info_service_pb2.GetPythonTemplateRequest(
            parameter_message_full_name=sig.parameter_message_full_name,
            return_value_message_full_name=sig.return_value_message_full_name,
            file_descriptor_set=sig.file_descriptor_set,
            python_code=code_execution_pb2.PythonCode(
                function_body="  return node_pb2.ReturnValue()"
            ),
        )
    )

  def test_preview_executed_code_invalid_task(self):
    task = behavior_tree.Task(
        behavior_call_pb2.BehaviorCall(skill_id="ai.intrinsic.some_skill")
    )

    with self.assertRaisesRegex(ValueError, "does not contain a PythonScript"):
      self._code_execution_info.preview_executed_code(task)

  def test_preview_executed_code_invalid_type(self):
    with self.assertRaises(TypeError):
      self._code_execution_info.preview_executed_code("invalid_script_object")


if __name__ == "__main__":
  absltest.main()
