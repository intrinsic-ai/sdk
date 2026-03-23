<!--
Copyright 2023 Intrinsic Innovation LLC
-->
# Golang resolvers for pb.Any types

We use gRPC gateway to generate golang code that translates between protobuf and
JSON. gRPC gateway must know the definition of a message to generate code.
Sometimes our gRPC services use messages with
[google.protobuf.Any](https://protobuf.dev/programming-guides/proto3/#any)
types. We don't know the definition of a message inside a `google.protobuf.Any`
(`Any`) at build time, so at runtime we must get a
[Descriptor](https://protobuf.dev/programming-guides/techniques/#self-describing-messages)
for that message type. The only information we have about the message type at
runtime is a string in the `type_url` field of the `Any`.

This folder contains golang
[MessageTypeResolver](https://pkg.go.dev/google.golang.org/protobuf/reflect/protoregistry#MessageTypeResolver)s
for `Any` types. Each resolver accepts a type URL and returns a descriptor for
that message.

## AnyResolver (public)

`AnyResolver` resolves `type_url` to message types using a combination of other
resolvers.

## GreedyResolver (private)

Given an ordered list of resolvers, `GreedyResolver` returns the result from the
first resolver that can resolve the type.

## ProtoRegistryResolver (private)

`ProtoRegisryResolver` resolves `type_url` to message types by querying the
ProtoRegistry service if the `type_url` begins with `type.intrinsic.ai/`.

## InstalledAssetsResolver (private)

The `InstalledAssetsResolver` resolves `type_url` to message types using
isolated file descriptor sets from all installed assets. The behavior is
undefined if multiple installed assets have different and incompatible
definitions for the same protobuf message.
