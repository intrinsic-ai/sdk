// Copyright 2023 Intrinsic Innovation LLC

package pubsub

func init() {
	PubsubCmd.AddCommand(
		NewServiceDeleteCmd(
			"stop-forwarding",
			"Stops forwarding of PubSub topics and KV store paths.",
			forwardingServicePackage,
			forwardingServiceName,
		),
	)
}
