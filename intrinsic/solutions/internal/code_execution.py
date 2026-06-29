# Copyright 2023 Intrinsic Innovation LLC

"""Python wrappers around code execution."""

from __future__ import annotations

import abc
import textwrap
import uuid

from intrinsic.executive.proto import code_execution_pb2
from intrinsic.solutions import blackboard_value
from intrinsic.solutions import proto_building
from intrinsic.solutions.internal import skill_utils

_DEFAULT_SCRIPT_NODE_PROTO_FILE = "node.proto"
# Expected to be filled with a UUID.
_DEFAULT_RETURN_VALUE_KEY_PATTERN = "pynode_%s"


class CodeExecution(abc.ABC):
  """Corresponds to a CodeExecution proto."""

  @property
  @abc.abstractmethod
  def proto(self) -> code_execution_pb2.CodeExecution:
    """Returns the proto representation."""

  @property
  @abc.abstractmethod
  def result(self) -> blackboard_value.BlackboardValue | None:
    """Returns a reference to the return value of the code execution."""


class PythonScript(CodeExecution):
  """Represents a Python script node in a behavior tree."""

  _signature_with_args: proto_building.SignatureWithArgs
  _function_body: str
  _return_value_key: str

  def __init__(
      self,
      signature_with_args: proto_building.SignatureWithArgs,
      *,
      function_body: str,
      return_value_key: str | None = None,
      create_unique_signature: bool = True,
  ):
    """Creates a PythonScript.

    The signature of the new script node is (by default) initialized to a unique
    copy of the given signature by renaming the "main" proto file of the file
    descriptor set (in which the parameter message and the return value message
    are defined) to `<random path>/node.proto`, so the given code can assume the
    module name `node_pb2` to be defined.

    The given code to be executed is the body of a function and must not contain
    a function signature. It can assume the following template:

    ```
    import numpy as np
    from intrinsic.skills.python import basic_compute_context
    from <random module path> import node_pb2

    def compute(
        params: node_pb2.Params,
        context: basic_compute_context.BasicComputeContext,
    ) -> node_pb2.ReturnValue:
      {function_body}
    ```

    Args:
      signature_with_args: Signature and arguments for the script node. By
        default, a unique copy of the given object will be created and used -
        the given object itself will not be stored or modified.
      function_body: Function body of the Python code to execute. Can be passed
        with arbitrary indentation (as long at it is consistent for all lines).
      return_value_key: Optional blackboard key under which to store the return
        value. If not provided, a unique key will be generated if the signature
        defines a return value message.
      create_unique_signature: If True (default), a unique copy of the
        signature will be created and used. If False (advanced), the signature
        will be used as is.
    """
    if create_unique_signature:
      self._signature_with_args = signature_with_args.unique_copy(
          _DEFAULT_SCRIPT_NODE_PROTO_FILE
      )
    else:
      self._signature_with_args = signature_with_args
    if not function_body.strip():
      raise ValueError("function_body must not be empty")
    # Normalize indentation. The code execution service expects indented code.
    self._function_body = textwrap.indent(
        textwrap.dedent(function_body).strip(), "  "
    )

    if self._signature_with_args.return_value_message_full_name:
      if return_value_key is None:
        self._return_value_key = (
            _DEFAULT_RETURN_VALUE_KEY_PATTERN % uuid.uuid4().hex
        )
      else:
        if not return_value_key:
          raise ValueError("return_value_key must not be empty")
        self._return_value_key = return_value_key
    else:
      if return_value_key is not None:
        raise ValueError(
            "return_value_key provided but signature does not define return"
            " value"
        )
      self._return_value_key = ""

  @property
  def proto(self) -> code_execution_pb2.CodeExecution:
    result = code_execution_pb2.CodeExecution(
        python_code=code_execution_pb2.PythonCode(
            function_body=self._function_body
        ),
        return_value_key=self._return_value_key,
        parameter_message_full_name=(
            self._signature_with_args.parameter_message_full_name
        ),
        return_value_message_full_name=(
            self._signature_with_args.return_value_message_full_name
        ),
    )

    if self._signature_with_args.file_descriptor_set.file:
      result.file_descriptor_set.CopyFrom(
          self._signature_with_args.file_descriptor_set
      )

    if self._signature_with_args.params_message is not None:
      result.parameters.proto.Pack(self._signature_with_args.params_message)

    for (
        path,
        cel_expression,
    ) in self._signature_with_args.blackboard_params.items():
      assignment = result.parameters.assign.add()
      assignment.path = path
      assignment.cel_expression = cel_expression

    return result

  @property
  def result(self) -> blackboard_value.BlackboardValue | None:
    if not self._signature_with_args.return_value_message_full_name:
      return None

    msg = skill_utils.create_message_from_file_descriptor_set(
        self._signature_with_args.file_descriptor_set,
        self._signature_with_args.return_value_message_full_name,
    )
    return blackboard_value.BlackboardValue(
        msg.DESCRIPTOR.fields_by_name,
        self._return_value_key,
        type(msg),
        None,
    )

  def unique_copy(self) -> PythonScript:
    """Creates a unique copy of the PythonScript.

    Returns a copy that uses a unique copy of the underlying signature with an
    equivalent file descriptor set whose "main" proto file (in which the
    parameter message and the return value message are defined) is renamed to
    `<random path>/node.proto`.

    This is useful, e.g., to "reuse" a single PythonScript inside a behavior
    tree since using the same PythonScript instance more than once in a single
    behavior tree is not allowed. Each PythonScript in a behavior tree must have
    a different file descriptor set with a different parameter and/or return
    value message name.
    """
    # The constructor internally creates a unique copy of the signature so we
    # don't really have to do anything here.
    return PythonScript(
        signature_with_args=self._signature_with_args,
        function_body=self._function_body,
        return_value_key=self._return_value_key,
        create_unique_signature=True,
    )

  @classmethod
  def _create_from_proto(
      cls, proto_object: code_execution_pb2.CodeExecution
  ) -> PythonScript:
    """Creates a PythonScript instance from a proto."""
    # Reconstruct Signature
    signature = proto_building.Signature(
        parameter_message_full_name=proto_object.parameter_message_full_name,
        return_value_message_full_name=proto_object.return_value_message_full_name,
        file_descriptor_set=proto_object.file_descriptor_set,
    )

    # Reconstruct SignatureWithArgs
    if proto_object.parameter_message_full_name:
      params_message = skill_utils.create_message_from_file_descriptor_set(
          proto_object.file_descriptor_set,
          proto_object.parameter_message_full_name,
      )
      if proto_object.parameters.proto.type_url:
        proto_object.parameters.proto.Unpack(params_message)

      signature_with_args = proto_building.SignatureWithArgs(
          signature=signature,
          params_message=params_message,
          blackboard_params={
              assignment.path: assignment.cel_expression
              for assignment in proto_object.parameters.assign
          },
      )
    else:
      signature_with_args = signature.with_args()

    return_value_key = (
        proto_object.return_value_key if proto_object.return_value_key else None
    )

    return cls(
        signature_with_args=signature_with_args,
        function_body=proto_object.python_code.function_body,
        return_value_key=return_value_key,
        create_unique_signature=False,
    )


def create_from_proto(
    proto_object: code_execution_pb2.CodeExecution,
) -> CodeExecution:
  """Creates a CodeExecution instance from a proto."""
  match proto_object.WhichOneof("code"):
    case "python_code":
      return PythonScript._create_from_proto(proto_object)
    case _:
      raise ValueError(
          f"Unsupported code execution type: {proto_object.WhichOneof('code')}"
      )
