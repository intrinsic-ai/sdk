// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_ICON_RELEASE_PORTABLE_INIT_INTRINSIC_H_
#define INTRINSIC_ICON_RELEASE_PORTABLE_INIT_INTRINSIC_H_

// Initializes an application by parsing the command-line flags.
//
void InitXfa(const char* usage, int argc, char* argv[]);
// Function alias to migrate consumers in src SoT.
extern void (&InitIntrinsic)(const char* usage, int argc, char* argv[]);

#endif  // INTRINSIC_ICON_RELEASE_PORTABLE_INIT_INTRINSIC_H_
