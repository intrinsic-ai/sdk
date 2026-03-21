// Copyright 2023 Intrinsic Innovation LLC

package util

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSpinner_BasicLifecycle(t *testing.T) {
	out := &bytes.Buffer{}
	ctx := context.Background()

	// Use a very fast interval for testing so ticks happen almost instantly
	spinner := NewSpinner(ctx, out, 1*time.Millisecond, PositionBack, StyleSpinningDots, ColorDefault, DirectionForward)

	spinner.Start("Loading")
	time.Sleep(5 * time.Millisecond) // let it spin a few times

	spinner.Stop("Done")

	output := out.String()
	// Should contain the start message
	assert.Contains(t, output, "Loading")
	// Should contain the stop message
	assert.Contains(t, output, "Done\n")
	// Should contain clear escape sequences (\r)
	assert.Contains(t, output, "\r")
}

func TestSpinner_UpdateMessage(t *testing.T) {
	out := &bytes.Buffer{}
	ctx := context.Background()

	spinner := NewSpinner(ctx, out, 5*time.Millisecond, PositionBack, StyleSpinningDots, ColorDefault, DirectionForward)
	spinner.Start("Step 1")
	time.Sleep(10 * time.Millisecond)

	spinner.UpdateMessage("Step 2")
	time.Sleep(10 * time.Millisecond)

	spinner.Stop("Finished")

	output := out.String()
	assert.Contains(t, output, "Step 1")
	assert.Contains(t, output, "Step 2")
	assert.Contains(t, output, "Finished")
}

func TestSpinner_Interrupt(t *testing.T) {
	out := &bytes.Buffer{}
	ctx := context.Background()

	spinner := NewSpinner(ctx, out, 5*time.Millisecond, PositionBack, StyleSpinningDots, ColorDefault, DirectionForward)
	spinner.Start("Working")
	time.Sleep(10 * time.Millisecond)

	spinner.Interrupt("Warning: Something happened")
	time.Sleep(10 * time.Millisecond)

	spinner.Stop("")

	output := out.String()
	assert.Contains(t, output, "Warning: Something happened")
	assert.Contains(t, output, "Working")

	// Ensure the interrupt message is printed followed by a newline (so it isn't cleared by \r)
	assert.True(t, strings.Contains(output, "Warning: Something happened\n"))
}
