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

package gpio

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	codespb "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/prototext"

	pb "intrinsic/hardware/gpio/v1/gpio_service_go_proto"
	signalpb "intrinsic/hardware/gpio/v1/signal_go_proto"
)

// tryConvertSignalValue attempts to convert value into a SignalValue message,
// first by attempting to unmarshal it as pbtxt, then by comparing it to
// various versions of "true" and "false".
func tryConvertSignalValue(value string) (*signalpb.SignalValue, error) {
	signalValue := &signalpb.SignalValue{}

	// First try to convert the value by unmarshaling pbtxt.
	if err := prototext.Unmarshal([]byte(value), signalValue); err == nil { // If NO error
		return signalValue, nil
	}

	if strings.ToLower(value) == "true" {
		signalValue.Value = &signalpb.SignalValue_BoolValue{BoolValue: true}
		return signalValue, nil
	}

	if strings.ToLower(value) == "false" {
		signalValue.Value = &signalpb.SignalValue_BoolValue{BoolValue: false}
		return signalValue, nil
	}

	return nil, fmt.Errorf("failed to convert value argument %q to a SignalValue", value)
}

func writeSignal(ctx context.Context, serverAddress, k8sContext, resourceInstanceName, signalName, value string) error {
	val, err := tryConvertSignalValue(value)
	if err != nil {
		return err
	}

	c, err := makeConnectionManager(ctx, serverAddress, k8sContext)
	if err != nil {
		return err
	}
	defer c.close()

	rctx := metadata.AppendToOutgoingContext(ctx, "x-resource-instance-name", resourceInstanceName)

	initialReq := &pb.OpenWriteSessionRequest{
		InitialSessionData: &pb.OpenWriteSessionRequest_InitialSessionData{
			SignalNames: []string{signalName},
		},
	}
	session, err := c.client().OpenWriteSession(rctx)
	if err != nil {
		return errors.Wrap(err, "open")
	}
	if err := session.Send(initialReq); err != nil {
		return errors.Wrap(err, "failed to send initial req")
	}
	resp, err := session.Recv()
	if err != nil {
		return errors.Wrap(err, "failed to receive")
	}

	if got, want := resp.Status.Code, int32(codespb.OK); got != want {
		return fmt.Errorf("initial session failed, got %v, want %v", got, want)
	}

	defer func() {
		if err := session.CloseSend(); err != nil {
			fmt.Printf("Failed to close write session: %v", err)
		}
	}()

	writeReq := &pb.OpenWriteSessionRequest{
		ActionRequest: &pb.OpenWriteSessionRequest_WriteSignals{
			WriteSignals: &pb.WriteSignalsRequest{
				SignalValues: &signalpb.SignalValueSet{
					Values: map[string]*signalpb.SignalValue{
						signalName: val,
					},
				},
			},
		},
	}

	if err := session.Send(writeReq); err != nil {
		return errors.Wrap(err, "failed to send write req")
	}

	writeResp, err := session.Recv()
	if err != nil {
		return errors.Wrap(err, "failed to receive write resp")
	}
	if got, want := writeResp.Status.Code, int32(codespb.OK); got != want {
		return fmt.Errorf("write failed, got %v, want %v", got, want)
	}
	return nil
}

var (
	flagSignalName string
	flagValue      string
)

var gpioWriteSignalsCmd = &cobra.Command{
	Use:   "write-signal",
	Short: "Sets the value of a signal",
	RunE: func(cmd *cobra.Command, args []string) error {
		return writeSignal(cmd.Context(), flagServerAddress, flagK8sContext, flagResourceInstanceName, flagSignalName, flagValue)
	},
}

func init() {
	gpioWriteSignalsCmd.PersistentFlags().StringVar(&flagSignalName, "signal_name", "", "signal name to set")
	gpioWriteSignalsCmd.PersistentFlags().StringVar(&flagValue, "value", "false", "value to set")
	gpioCmd.AddCommand(gpioWriteSignalsCmd)
}
