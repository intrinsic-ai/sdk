// Copyright 2023 Intrinsic Innovation LLC

package pubsub

func init() {
	PubsubCmd.AddCommand(
		NewServiceDeleteCmd(
			"hub-service-delete",
			"Deletes the PubSub Hub service from the currently running solution.",
			hubServicePackage,
			hubServiceName,
		),
	)
}
