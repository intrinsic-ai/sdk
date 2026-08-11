# Copyright 2023 Intrinsic Innovation LLC

"""Provides utility functions for connecting to gRPC services defined by GrpcConnection messages."""

from typing import Any
from typing import Sequence

from absl import logging
import grpc

from intrinsic.assets.proto.v1 import grpc_connection_pb2
from intrinsic.util.grpc import interceptor


def connect(
    connection: grpc_connection_pb2.GrpcConnection,
    grpc_options: Sequence[tuple[str, Any]] | None = None,
) -> grpc.Channel:
  """Creates a gRPC channel for communicating with the specified provider.

  The returned channel will be intercepted to include any needed metadata for
  communicating with the provider.

  Args:
    connection: The GrpcConnection message containing address and metadata.
    grpc_options: Additional options for the gRPC channel.

  Returns:
    A gRPC channel.
  """
  metadata = connection.metadata
  headers = [f"{m.key}={m.value}" for m in metadata]
  logging.info(
      'Connecting to address "%s" with headers injected through'
      " HeaderAdderInterceptor: [%s]",
      connection.address,
      ", ".join(headers),
  )

  channel = grpc.intercept_channel(
      grpc.insecure_channel(connection.address, grpc_options),
      interceptor.HeaderAdderInterceptor(
          lambda: [(m.key, m.value) for m in metadata]
      ),
  )

  return channel
