# Copyright 2023 Intrinsic Innovation LLC

"""Main function for the "approximate function" CustomCalculation Service."""

from collections.abc import Sequence
from concurrent import futures

from absl import app
from absl import logging
import grpc

from intrinsic.assets.dependencies import utils
from intrinsic.assets.services.examples.calcserver import approximate_function
from intrinsic.assets.services.examples.calcserver import approximate_function_pb2
from intrinsic.assets.services.examples.calcserver import calc_server_pb2_grpc
from intrinsic.assets.services.examples.calcserver import two_dimensional_function_data_pb2
from intrinsic.resources.proto import runtime_context_pb2

_DEFAULT_DEGREE = 2


def _get_runtime_context() -> runtime_context_pb2.RuntimeContext:
  with open("/etc/intrinsic/runtime_config.pb", "rb") as f:
    return runtime_context_pb2.RuntimeContext.FromString(f.read())


def main(argv: Sequence[str]) -> None:
  del argv  # unused

  context = _get_runtime_context()
  port = context.port
  config = approximate_function_pb2.ApproximateFunctionConfig()
  if not context.config.Unpack(config):
    raise RuntimeError("Failed to unpack config")

  # Retrieve the function data.
  data = None
  try:
    data_any = utils.get_data_payload(
        dep=config.data,
        iface=f"data://{two_dimensional_function_data_pb2.TwoDimensionalFunctionData.DESCRIPTOR.full_name}",
    )
  except utils.MissingInterfaceError as e:
    logging.error("No data provided to approximate function: %s", e)
  else:
    data = two_dimensional_function_data_pb2.TwoDimensionalFunctionData()
    if not data_any.Unpack(data):
      raise RuntimeError("Failed to unpack data")

  # Resolve the polynomial degree.
  degree = config.degree if config.HasField("degree") else _DEFAULT_DEGREE

  logging.info("Starting ApproximateFunction on port %d", port)
  logging.info("Degree: %d", degree)
  if data is not None:
    data_strings = [
        f"({x}, {y}, {f})" for x, y, f in zip(data.x, data.y, data.f)
    ]
    data_msg = "\n" + "\n".join(data_strings)
  else:
    data_msg = "None"
  logging.info("Data: %s", data_msg)

  service = approximate_function.ApproximateFunction(data, degree)
  server = grpc.server(
      futures.ThreadPoolExecutor(max_workers=4),
      options=(("grpc.so_reuseport", 0),),
  )
  calc_server_pb2_grpc.add_CustomCalculationServicer_to_server(service, server)
  added_port = server.add_insecure_port(f"[::]:{port}")
  if added_port != port:
    raise RuntimeError(f"Failed to use port {port}")

  server.start()

  logging.info("ApproximateFunction listening on port %d", port)

  server.wait_for_termination()


if __name__ == "__main__":
  app.run(main)
