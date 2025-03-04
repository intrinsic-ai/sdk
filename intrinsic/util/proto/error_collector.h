// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_UTIL_PROTO_ERROR_COLLECTOR_H_
#define INTRINSIC_UTIL_PROTO_ERROR_COLLECTOR_H_

#include <string>
#include <vector>

#include "absl/strings/str_format.h"
#include "absl/strings/str_join.h"
#include "absl/strings/string_view.h"
#include "google/protobuf/io/tokenizer.h"

namespace intrinsic {

// A simple error collector that collects all errors and warnings into a single
// string.
class SimpleErrorCollector : public google::protobuf::io::ErrorCollector {
 public:
  void RecordError(int line, int column, absl::string_view message) override {
    errors_.push_back(absl::StrFormat("Error in line %d (column %d): %s", line,
                                      column, message));
  }

  void RecordWarning(int line, int column, absl::string_view message) override {
    errors_.push_back(absl::StrFormat("Warning in line %d (column %d): %s",
                                      line, column, message));
  }

  std::string str() const { return absl::StrJoin(errors_, "\n"); }

 private:
  std::vector<std::string> errors_;
};

}  // namespace intrinsic

#endif  // INTRINSIC_UTIL_PROTO_ERROR_COLLECTOR_H_
