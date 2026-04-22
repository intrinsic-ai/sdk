// Copyright 2023 Intrinsic Innovation LLC

#include <signal.h>
#include <stddef.h>

#include <cstring>

#include "absl/debugging/failure_signal_handler.h"
#include "absl/debugging/symbolize.h"
#include "absl/flags/flag.h"
#include "absl/flags/parse.h"
#include "absl/flags/usage.h"
#include "absl/log/globals.h"
#include "absl/log/initialize.h"
#include "absl/log/log.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "intrinsic/icon/utils/log.h"

ABSL_FLAG(bool, sleep, false,
          "(optional) Do nothing and sleep indefinitely. This is used to "
          "download the image and avoid crash looping.");

void InitIntrinsic(const char* usage, int argc, char* argv[]) {
  if (usage != nullptr && strlen(usage) > 0) {
    absl::SetProgramUsageMessage(usage);
  }
  absl::ParseCommandLine(argc, argv);
  absl::InitializeLog();
  absl::SetStderrThreshold(absl::LogSeverityAtLeast::kInfo);

  // Provide stack traces on SIGSEGV and other signals.
  absl::InitializeSymbolizer(argv[0]);
  absl::FailureSignalHandlerOptions options;
  options.call_previous_handler = true;
  absl::InstallFailureSignalHandler(options);
  // Restore the default SIGTERM handler. This is to avoid getting confusing
  // stack traces printed when we purposefully kill subprocesses.
  struct sigaction action;
  action.sa_handler = SIG_DFL;        // Set to Default
  sigemptyset(&action.sa_mask);       // Don't block any other signals
  action.sa_flags = 0;                // No special flags
  sigaction(SIGTERM, &action, NULL);  // Apply the change

  intrinsic::RtLogInitForThisThread();

  if (absl::GetFlag(FLAGS_sleep)) {
    LOG(INFO) << "Started with --sleep=true. Sleeping indefinitely...";
    absl::SleepFor(absl::InfiniteDuration());
  }

  LOG(INFO) << "********* Process Begin *********";
}
