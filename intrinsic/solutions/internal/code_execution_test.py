# Copyright 2023 Intrinsic Innovation LLC

"""Tests for code_execution.py."""

from unittest import mock
import uuid

from absl.testing import absltest
from google.protobuf import descriptor_pb2
from google.protobuf import text_format

from intrinsic.executive.proto import code_execution_pb2
from intrinsic.solutions import cel
from intrinsic.solutions import proto_building
from intrinsic.solutions.internal import code_execution
from intrinsic.solutions.testing import compare


class PythonScriptTest(absltest.TestCase):

  def test_python_script_empty_function_body(self):
    sig_with_args = proto_building.Signature(
        parameter_message_full_name='',
        return_value_message_full_name='',
        file_descriptor_set=None,
    ).with_args()

    with self.assertRaisesRegex(ValueError, 'function_body'):
      code_execution.PythonScript(
          signature_with_args=sig_with_args,
          function_body='',
      )
    with self.assertRaisesRegex(ValueError, 'function_body'):
      code_execution.PythonScript(
          signature_with_args=sig_with_args,
          function_body='   \n  ',
      )

  def test_python_script_empty_return_value_key(self):
    sig_with_args = proto_building.Signature(
        parameter_message_full_name='',
        return_value_message_full_name='gen.sbl.ReturnValue',
        file_descriptor_set=text_format.Parse(
            """
            file {
                name: "gen/sbl/signature.proto"
                package: "gen.sbl"
                message_type {
                    name: "ReturnValue"
                    field {
                        name: "int_out"
                        number: 3
                        label: LABEL_OPTIONAL
                        type: TYPE_INT64
                    }
                }
                syntax: "proto3"
            }""",
            descriptor_pb2.FileDescriptorSet(),
        ),
    ).with_args()

    with self.assertRaisesRegex(ValueError, 'return_value_key.*empty'):
      code_execution.PythonScript(
          signature_with_args=sig_with_args,
          function_body='pass',
          return_value_key='',
      )

  def test_python_script_return_value_key_but_no_return_value_message(
      self,
  ):
    sig_with_args = proto_building.Signature(
        parameter_message_full_name='',
        return_value_message_full_name='',
        file_descriptor_set=None,
    ).with_args()

    with self.assertRaisesRegex(ValueError, 'return_value_key provided'):
      code_execution.PythonScript(
          signature_with_args=sig_with_args,
          function_body='pass',
          return_value_key='some_key',
      )
    with self.assertRaisesRegex(ValueError, 'return_value_key provided'):
      code_execution.PythonScript(
          signature_with_args=sig_with_args,
          function_body='pass',
          return_value_key='',
      )

  def test_python_script_proto(self):
    sig_with_args = proto_building.Signature(
        parameter_message_full_name='gen.sbl.Params',
        return_value_message_full_name='gen.sbl.ReturnValue',
        file_descriptor_set=text_format.Parse(
            """
            file {
                name: "gen/sbl/signature.proto"
                package: "gen.sbl"
                message_type {
                    name: "Params"
                    field {
                        name: "int_in"
                        number: 1
                        label: LABEL_OPTIONAL
                        type: TYPE_INT64
                    }
                    field {
                        name: "str_in"
                        number: 2
                        label: LABEL_OPTIONAL
                        type: TYPE_STRING
                    }
                }
                message_type {
                    name: "ReturnValue"
                    field {
                        name: "int_out"
                        number: 3
                        label: LABEL_OPTIONAL
                        type: TYPE_INT64
                    }
                }
                syntax: "proto3"
            }""",
            descriptor_pb2.FileDescriptorSet(),
        ),
    ).with_args(int_in=42, str_in=cel.CelExpression('some.cel.expression'))

    # Mock the UUID used in the creation of the unique signature copy
    with mock.patch(
        'uuid.uuid4',
        return_value=uuid.UUID('11111111-aaaa-2222-bbbb-333333333333'),
    ):
      script = code_execution.PythonScript(
          signature_with_args=sig_with_args,
          function_body="""
                result = node_pb2.ReturnValue(
                    int_out = 2 * params.int_in
                )
                return result

                """,
          return_value_key='my_result',
      )

    proto = script.proto

    compare.assertProto2Equal(
        self,
        proto,
        """
        python_code {
            function_body: "  result = node_pb2.ReturnValue(\\n      int_out = 2 * params.int_in\\n  )\\n  return result"
        }
        parameters {
            proto {
                type_url: "type.googleapis.com/gen.sbl_11111111aaaa2222bbbb333333333333.Params"
                # field 1 with int value 42
                value: "\\010*"
            }
            assign {
                path: "str_in"
                cel_expression: "some.cel.expression"
            }
        }
        return_value_key: "my_result"
        parameter_message_full_name: "gen.sbl_11111111aaaa2222bbbb333333333333.Params"
        return_value_message_full_name: "gen.sbl_11111111aaaa2222bbbb333333333333.ReturnValue"
        file_descriptor_set {
            file {
                name: "gen/sbl_11111111aaaa2222bbbb333333333333/node.proto"
                package: "gen.sbl_11111111aaaa2222bbbb333333333333"
                message_type {
                name: "Params"
                    field {
                        name: "int_in"
                        number: 1
                        label: LABEL_OPTIONAL
                        type: TYPE_INT64
                    }
                    field {
                        name: "str_in"
                        number: 2
                        label: LABEL_OPTIONAL
                        type: TYPE_STRING
                    }
                }
                message_type {
                name: "ReturnValue"
                    field {
                        name: "int_out"
                        number: 3
                        label: LABEL_OPTIONAL
                        type: TYPE_INT64
                    }
                }
                syntax: "proto3"
            }
        }
        """,
    )

  def test_python_script_proto_auto_generated_return_value_key(self):
    sig_with_args = proto_building.Signature(
        parameter_message_full_name='',
        return_value_message_full_name='gen.sbl.ReturnValue',
        file_descriptor_set=text_format.Parse(
            """
            file {
                name: "gen/sbl/signature.proto"
                package: "gen.sbl"
                message_type {
                    name: "ReturnValue"
                    field {
                        name: "int_out"
                        number: 3
                        label: LABEL_OPTIONAL
                        type: TYPE_INT64
                    }
                }
                syntax: "proto3"
            }""",
            descriptor_pb2.FileDescriptorSet(),
        ),
    ).with_args()

    # Mock the UUIDs used for:
    # - creation of the unique signature copy
    # - return value key
    with mock.patch(
        'uuid.uuid4',
        side_effect=[
            uuid.UUID('11111111-aaaa-2222-bbbb-333333333333'),
            uuid.UUID('77777777-eeee-8888-ffff-999999999999'),
        ],
    ):
      script = code_execution.PythonScript(
          signature_with_args=sig_with_args,
          function_body='return node_pb2.ReturnValue(int_out=42)',
          return_value_key=None,
      )

    proto = script.proto

    compare.assertProto2Equal(
        self,
        proto,
        """
            python_code {
                function_body: "  return node_pb2.ReturnValue(int_out=42)"
            }
            return_value_key: "pynode_77777777eeee8888ffff999999999999"
            return_value_message_full_name: "gen.sbl_11111111aaaa2222bbbb333333333333.ReturnValue"
            file_descriptor_set {
                file {
                    name: "gen/sbl_11111111aaaa2222bbbb333333333333/node.proto"
                    package: "gen.sbl_11111111aaaa2222bbbb333333333333"
                    message_type {
                        name: "ReturnValue"
                        field {
                            name: "int_out"
                            number: 3
                            label: LABEL_OPTIONAL
                            type: TYPE_INT64
                        }
                    }
                    syntax: "proto3"
                }
            }
            """,
    )

  def test_python_script_result(self):
    sig_with_args = proto_building.Signature(
        parameter_message_full_name='',
        return_value_message_full_name='gen.sbl.ReturnValue',
        file_descriptor_set=text_format.Parse(
            """
            file {
                name: "gen/sbl/signature.proto"
                package: "gen.sbl"
                message_type {
                    name: "ReturnValue"
                    field {
                        name: "int_out"
                        number: 3
                        label: LABEL_OPTIONAL
                        type: TYPE_INT64
                    }
                }
                syntax: "proto3"
            }""",
            descriptor_pb2.FileDescriptorSet(),
        ),
    ).with_args()

    # Mock the UUIDs used for:
    # - creation of the unique signature copy
    # - return value key
    with mock.patch(
        'uuid.uuid4',
        side_effect=[
            uuid.UUID('11111111-aaaa-2222-bbbb-333333333333'),
            uuid.UUID('77777777-eeee-8888-ffff-999999999999'),
        ],
    ):
      script = code_execution.PythonScript(
          signature_with_args=sig_with_args,
          function_body='return node_pb2.ReturnValue(int_out=42)',
          return_value_key=None,
      )

    result = script.result

    self.assertTrue(result.is_toplevel_value)
    self.assertEqual(
        result.value_access_path(),
        'pynode_77777777eeee8888ffff999999999999',
    )
    self.assertEqual(
        result.value_type.DESCRIPTOR.full_name,
        'gen.sbl_11111111aaaa2222bbbb333333333333.ReturnValue',
    )
    self.assertEqual(
        result.int_out.value_access_path(),
        'pynode_77777777eeee8888ffff999999999999.int_out',
    )

  def test_python_script_result_no_return_value_message(self):
    sig_with_args = proto_building.Signature(
        parameter_message_full_name='',
        return_value_message_full_name='',
        file_descriptor_set=None,
    ).with_args()

    script = code_execution.PythonScript(
        signature_with_args=sig_with_args,
        function_body='return node_pb2.ReturnValue(int_out=42)',
        return_value_key=None,
    )

    self.assertIsNone(script.result)

  def test_python_script_unique_copy(self):
    sig_with_args = proto_building.Signature(
        parameter_message_full_name='gen.sbl.Params',
        return_value_message_full_name='gen.sbl.ReturnValue',
        file_descriptor_set=text_format.Parse(
            """
            file {
                name: "gen/sbl/signature.proto"
                package: "gen.sbl"
                message_type {
                    name: "Params"
                    field {
                        name: "int_in"
                        number: 1
                        label: LABEL_OPTIONAL
                        type: TYPE_INT64
                    }
                }
                message_type {
                    name: "ReturnValue"
                    field {
                        name: "int_out"
                        number: 3
                        label: LABEL_OPTIONAL
                        type: TYPE_INT64
                    }
                }
                syntax: "proto3"
            }""",
            descriptor_pb2.FileDescriptorSet(),
        ),
    ).with_args(int_in=42)

    script = code_execution.PythonScript(
        signature_with_args=sig_with_args,
        function_body='return node_pb2.ReturnValue(int_out=params.int_in)',
        return_value_key='my_result',
    )

    # Mock the UUID for the unique copy creation
    with mock.patch(
        'uuid.uuid4',
        return_value=uuid.UUID('99999999-9999-9999-9999-999999999999'),
    ):
      copy = script.unique_copy()

    # Verify they are different instances
    self.assertIsNot(script, copy)
    compare.assertProto2Equal(
        self,
        copy.proto,
        """
        python_code {
            function_body: "  return node_pb2.ReturnValue(int_out=params.int_in)"
        }
        parameters {
            proto {
                type_url: "type.googleapis.com/gen.sbl_99999999999999999999999999999999.Params"
                # field 1 with int value 42
                value: "\\010*"
            }
        }
        return_value_key: "my_result"
        parameter_message_full_name: "gen.sbl_99999999999999999999999999999999.Params"
        return_value_message_full_name: "gen.sbl_99999999999999999999999999999999.ReturnValue"
        file_descriptor_set {
            file {
                name: "gen/sbl_99999999999999999999999999999999/node.proto"
                package: "gen.sbl_99999999999999999999999999999999"
                message_type {
                name: "Params"
                    field {
                        name: "int_in"
                        number: 1
                        label: LABEL_OPTIONAL
                        type: TYPE_INT64
                    }
                }
                message_type {
                name: "ReturnValue"
                    field {
                        name: "int_out"
                        number: 3
                        label: LABEL_OPTIONAL
                        type: TYPE_INT64
                    }
                }
                syntax: "proto3"
            }
        }
        """,
    )

  def test_python_script_create_from_proto(self):
    proto = text_format.Parse(
        """
        python_code {
            function_body: "  result = script_node_pb2.ReturnValue(int_out=params.int_in)"
        }
        parameters {
            proto {
                type_url: "type.googleapis.com/my_package.Params"
                # field 1 with int value 42
                value: "\\010*"
            }
            assign {
                path: "str_in"
                cel_expression: "some.cel.expression"
            }
        }
        return_value_key: "my_result"
        parameter_message_full_name: "my_package.Params"
        return_value_message_full_name: "my_package.ReturnValue"
        file_descriptor_set {
            file {
                name: "my_package/script_node.proto"
                package: "my_package"
                message_type {
                name: "Params"
                    field {
                        name: "int_in"
                        number: 1
                        label: LABEL_OPTIONAL
                        type: TYPE_INT64
                    }
                    field {
                        name: "str_in"
                        number: 2
                        label: LABEL_OPTIONAL
                        type: TYPE_STRING
                    }
                }
                message_type {
                name: "ReturnValue"
                    field {
                        name: "int_out"
                        number: 3
                        label: LABEL_OPTIONAL
                        type: TYPE_INT64
                    }
                }
                syntax: "proto3"
            }
        }
        """,
        code_execution_pb2.CodeExecution(),
    )

    script = code_execution.create_from_proto(proto)

    self.assertIsInstance(script, code_execution.PythonScript)
    # Verify roundtrip is lossless
    compare.assertProto2Equal(self, script.proto, proto)

  def test_python_script_create_from_proto_no_params(self):
    proto = text_format.Parse(
        """
        python_code {
          function_body: "  pass"
        }
        """,
        code_execution_pb2.CodeExecution(),
    )

    script = code_execution.create_from_proto(proto)

    self.assertIsInstance(script, code_execution.PythonScript)
    # Verify roundtrip is lossless
    compare.assertProto2Equal(self, script.proto, proto)

  # This test documents behavior which is not absolutely required. We can change
  # it if necessary.
  def test_python_script_create_from_proto_normalizes_code(self):
    proto = text_format.Parse(
        """
        python_code {
          function_body: "    print('hello')\\n"
        }
        """,
        code_execution_pb2.CodeExecution(),
    )

    script = code_execution.create_from_proto(proto)

    self.assertIsInstance(script, code_execution.PythonScript)
    # Verify roundtrip is "functionally lossless"
    compare.assertProto2Equal(
        self,
        script.proto,
        """
        python_code {
          # Indentation has changed (4 -> 2) and trailing newline is gone
          function_body: "  print('hello')"
        }
        """,
    )


class GetFunctionBodyAsStrTest(absltest.TestCase):

  def test_simple_function(self):
    def sample_func(a, b):
      # Calculate sum
      c = a + b
      return c

    body = code_execution.get_function_body_as_str(sample_func)
    self.assertEqual(body, '# Calculate sum\nc = a + b\nreturn c')

  def test_type_annotations(self):
    def sample_func(
        a: int = 1, b: str = 'test', c: dict[str, list[int]] = {}
    ) -> str:
      return f'{b}_{a}'

    body = code_execution.get_function_body_as_str(sample_func)
    self.assertEqual(body, "return f'{b}_{a}'")

  def test_single_line(self):
    # fmt: off
    def sample_func(): return 42
    # fmt: on

    body = code_execution.get_function_body_as_str(sample_func)
    self.assertEqual(body, 'return 42')

  def test_multiline_header(self):
    def sample_func(
        parameter_one: int,
        parameter_two: str,
        parameter_three: bool,
        parameter_four: float,
    ) -> str:
      print(42)
      return """A long string value"""

    body = code_execution.get_function_body_as_str(sample_func)
    self.assertEqual(body, 'print(42)\nreturn """A long string value"""')

  def test_class_method(self):
    class Dummy:

      def method(self):
        val = 42
        return val

    body = code_execution.get_function_body_as_str(Dummy.method)
    self.assertEqual(body, 'val = 42\nreturn val')

    body = code_execution.get_function_body_as_str(Dummy().method)
    self.assertEqual(body, 'val = 42\nreturn val')

  def test_builtin(self):
    with self.assertRaises(ValueError):
      code_execution.get_function_body_as_str(len)

  def test_lambda(self):
    with self.assertRaises(ValueError):
      code_execution.get_function_body_as_str(lambda x: x + 1)


if __name__ == '__main__':
  absltest.main()
