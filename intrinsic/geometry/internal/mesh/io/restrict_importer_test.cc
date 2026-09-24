// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#include "intrinsic/geometry/internal/mesh/io/restrict_importer.h"

#include <assimp/BaseImporter.h>
#include <gmock/gmock.h>
#include <gtest/gtest.h>

#include <assimp/Importer.hpp>
#include <set>
#include <string>

namespace intrinsic::geo {
namespace {

using ::testing::Gt;

struct RestrictImporterTestCase {
  std::string input_extension;
  std::string expected_extension;
};

class RestrictImporterParameterizedTest
    : public ::testing::TestWithParam<RestrictImporterTestCase> {};

TEST_P(RestrictImporterParameterizedTest, RestrictsToSupportedExtension) {
  const auto& [input_extension, expected_extension] = GetParam();

  Assimp::Importer importer;
  ASSERT_THAT(importer.GetImporterCount(), Gt(0));

  RestrictImporterToExtension(importer, input_extension);

  const size_t count = importer.GetImporterCount();
  EXPECT_THAT(count, Gt(0));

  for (size_t i = 0; i < count; ++i) {
    Assimp::BaseImporter* imp = importer.GetImporter(i);
    ASSERT_NE(imp, nullptr);
    std::set<std::string> extensions;
    imp->GetExtensionList(extensions);
    EXPECT_TRUE(extensions.contains(expected_extension));
  }
}

INSTANTIATE_TEST_SUITE_P(ValidExtensions, RestrictImporterParameterizedTest,
                         ::testing::Values(
                             RestrictImporterTestCase{
                                 .input_extension = "glb",
                                 .expected_extension = "glb",
                             },
                             RestrictImporterTestCase{
                                 .input_extension = ".obj",
                                 .expected_extension = "obj",
                             },
                             RestrictImporterTestCase{
                                 .input_extension = ".ply",
                                 .expected_extension = "ply",
                             }));

TEST(RestrictImporterTest, RemovesAllForUnknownExtension) {
  Assimp::Importer importer;
  RestrictImporterToExtension(importer, "nonexistent_extension");
  EXPECT_EQ(importer.GetImporterCount(), 0);
}

}  // namespace
}  // namespace intrinsic::geo
