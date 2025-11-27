# Copyright 2023 Intrinsic Innovation LLC

import enum

from absl.testing import absltest
from intrinsic.executive.proto import test_message_pb2
from intrinsic.solutions import utils


class UtilsTest(absltest.TestCase):
  """Tests functions in utils."""

  def test_is_iterable(self):
    """Tests is_iterable."""

    class NonIterableTestClass:

      def __init__(self):
        pass

    class IterableTestClass:

      def __init__(self):
        pass

      def __iter__(self):
        for i in range(0, 5):
          yield i

    self.assertTrue(utils.is_iterable(list()))
    self.assertTrue(utils.is_iterable({}))
    self.assertTrue(utils.is_iterable('foo'))
    self.assertTrue(utils.is_iterable(IterableTestClass()))
    self.assertFalse(utils.is_iterable(5))
    self.assertFalse(utils.is_iterable(NonIterableTestClass()))


class PrefixOptionsTest(absltest.TestCase):
  """Tests for utils.PrefixOptions."""

  def test_init(self):
    options = utils.PrefixOptions(
        intrinsic_prefix='my_intrinsic',
        world_prefix='my_world',
        skill_prefix='my_skills',
    )

    self.assertEqual(options.intrinsic_prefix, 'my_intrinsic')
    self.assertEqual(options.world_prefix, 'my_world')
    self.assertEqual(options.skill_prefix, 'my_skills')


class ProtoEnumTest(absltest.TestCase):
  """Tests utils.protoenum decorator."""

  def test_maps_all_values(self):
    """Tests that all proto enum values are mapped indeed."""

    @utils.protoenum(proto_enum_type=test_message_pb2.TestEnum)
    class TestEnum:
      pass

    self.assertEqual(
        TestEnum.TEST_ENUM_UNSPECIFIED.value,
        test_message_pb2.TestEnum.TEST_ENUM_UNSPECIFIED,
    )
    self.assertEqual(
        TestEnum.TEST_ENUM_1.value,
        test_message_pb2.TestEnum.TEST_ENUM_1,
    )
    self.assertEqual(
        TestEnum.TEST_ENUM_2.value,
        test_message_pb2.TestEnum.TEST_ENUM_2,
    )
    self.assertEqual(
        TestEnum.TEST_ENUM_3.value,
        test_message_pb2.TestEnum.TEST_ENUM_3,
    )

  def test_does_not_map_unspecified_if_mapped_to_none(self):
    """Tests that the unspecified value does not appear in enum when mapped to None."""

    @utils.protoenum(
        proto_enum_type=test_message_pb2.TestEnum,
        unspecified_proto_enum_map_to_none=test_message_pb2.TestEnum.TEST_ENUM_UNSPECIFIED,
    )
    class TestEnum:
      pass

    with self.assertRaises(AttributeError):
      _ = TestEnum.TEST_ENUM_UNSPECIFIED

  def test_strip_prefix(self):
    """Tests that prefixes are properly stripped from member names."""

    @utils.protoenum(
        proto_enum_type=test_message_pb2.TestEnum,
        unspecified_proto_enum_map_to_none=test_message_pb2.TestEnum.TEST_ENUM_UNSPECIFIED,
        strip_prefix='TEST_',
    )
    class TestEnum:
      pass

    self.assertEqual(
        TestEnum.ENUM_1.value,
        test_message_pb2.TestEnum.TEST_ENUM_1,
    )
    self.assertEqual(
        TestEnum.ENUM_2.value,
        test_message_pb2.TestEnum.TEST_ENUM_2,
    )
    self.assertEqual(
        TestEnum.ENUM_3.value,
        test_message_pb2.TestEnum.TEST_ENUM_3,
    )

  def test_from_proto(self):
    """Tests from_proto method of wrapped enum."""

    @utils.protoenum(
        proto_enum_type=test_message_pb2.TestEnum,
        unspecified_proto_enum_map_to_none=test_message_pb2.TestEnum.TEST_ENUM_UNSPECIFIED,
    )
    class TestEnum(enum.Enum):
      pass

    self.assertIsNone(TestEnum.from_proto(None))

    self.assertIsNone(
        TestEnum.from_proto(test_message_pb2.TestEnum.TEST_ENUM_UNSPECIFIED)
    )

    self.assertEqual(
        TestEnum.from_proto(test_message_pb2.TestEnum.TEST_ENUM_1),
        TestEnum.TEST_ENUM_1,
    )
    self.assertEqual(
        TestEnum.from_proto(test_message_pb2.TestEnum.TEST_ENUM_2),
        TestEnum.TEST_ENUM_2,
    )
    self.assertEqual(
        TestEnum.from_proto(test_message_pb2.TestEnum.TEST_ENUM_3),
        TestEnum.TEST_ENUM_3,
    )

  def test_from_proto_with_alias(self):
    """Tests from_proto method of wrapped enum with aliases."""

    @utils.protoenum(
        proto_enum_type=test_message_pb2.TestAliasedEnum,
        unspecified_proto_enum_map_to_none=test_message_pb2.TestAliasedEnum.TEST_ALIASED_ENUM_UNSPECIFIED,
    )
    class TestAliasedEnum(enum.Enum):
      pass

    self.assertIsNone(TestAliasedEnum.from_proto(None))

    self.assertIsNone(
        TestAliasedEnum.from_proto(
            test_message_pb2.TestAliasedEnum.TEST_ALIASED_ENUM_UNSPECIFIED
        )
    )

    self.assertEqual(
        TestAliasedEnum.from_proto(
            test_message_pb2.TestAliasedEnum.TEST_ALIASED_ENUM_1A
        ),
        TestAliasedEnum.TEST_ALIASED_ENUM_1A,
    )
    self.assertEqual(
        TestAliasedEnum.from_proto(
            test_message_pb2.TestAliasedEnum.TEST_ALIASED_ENUM_1B
        ),
        TestAliasedEnum.TEST_ALIASED_ENUM_1A,
    )
    self.assertEqual(
        TestAliasedEnum.from_proto(
            test_message_pb2.TestAliasedEnum.TEST_ALIASED_ENUM_1B
        ),
        TestAliasedEnum.TEST_ALIASED_ENUM_1B,
    )
    self.assertEqual(
        TestAliasedEnum.from_proto(
            test_message_pb2.TestAliasedEnum.TEST_ALIASED_ENUM_2
        ),
        TestAliasedEnum.TEST_ALIASED_ENUM_2,
    )


if __name__ == '__main__':
  absltest.main()
