// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SCENE_VALIDATE_GO_SCENE_OBJECT_VALIDATION_C_H_
#define INTRINSIC_SCENE_VALIDATE_GO_SCENE_OBJECT_VALIDATION_C_H_

#include "absl/base/attributes.h"

#ifdef __cplusplus
extern "C" {
#endif

// Returns 0 on success, absl::StatusCode on failure.
// If failure, *error_message is allocated and must be freed by caller.
ABSL_ATTRIBUTE_UNUSED int intrinsic_scene_object_go_ValidateSceneObject(
    const char* proto_data, int proto_len, char** error_message);

#ifdef __cplusplus
}  // extern "C"
#endif

#endif  // INTRINSIC_SCENE_VALIDATE_GO_SCENE_OBJECT_VALIDATION_C_H_
