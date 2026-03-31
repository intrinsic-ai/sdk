// Copyright 2023 Intrinsic Innovation LLC

// Package main provides a demonstration of the spinner utility, iterating through
// all available styles, colors, and positional settings.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"intrinsic/tools/inctl/util"
)

func main() {
	ctx := context.Background()

	styles := []struct {
		name  string
		style util.SpinnerStyle
	}{
		{"StyleSpinningDots", util.StyleSpinningDots},
		{"StyleFillUp", util.StyleFillUp},
		{"StyleWaveThrough", util.StyleWaveThrough},
		{"StyleBouncingWave", util.StyleBouncingWave},
		{"StyleBouncingLine", util.StyleBouncingLine},
		{"StyleSweepingDots", util.StyleSweepingDots},
		{"StyleLine", util.StyleLine},
		{"StyleArc", util.StyleArc},
		{"StyleBarThrough", util.StyleBarThrough},
		{"StyleDotsDrop", util.StyleDotsDrop},
		{"StyleSimpleDots", util.StyleSimpleDots},
		{"StyleStar", util.StyleStar},
		{"StyleFlip", util.StyleFlip},
		{"StyleNoise", util.StyleNoise},
		{"StyleCircleHalves", util.StyleCircleHalves},
		{"StyleToggle", util.StyleToggle},
		{"StyleBouncingBall", util.StyleBouncingBall},
	}

	colors := []struct {
		name string
		col  util.SpinnerColor
	}{
		{"ColorDefault", util.ColorDefault},
		{"ColorRGB (Rainbow)", util.ColorRGB},
		{"ColorRed", util.ColorRed},
		{"ColorGreen", util.ColorGreen},
		{"ColorBlue", util.ColorBlue},
	}

	fmt.Println("=== Intrinsic Spinner Utility Demo ===")

	// 1. Front position demo
	frontSpinner := util.NewSpinner(ctx, os.Stdout, 80*time.Millisecond, util.PositionFront, util.StyleSpinningDots, util.ColorDefault, util.DirectionForward)
	frontSpinner.Start("Demonstrating PositionFront for 5 seconds...")
	time.Sleep(3 * time.Second)
	frontSpinner.Stop("")

	// 2. Back position demo
	backSpinner := util.NewSpinner(ctx, os.Stdout, 80*time.Millisecond, util.PositionBack, util.StyleSpinningDots, util.ColorDefault, util.DirectionReverse)
	backSpinner.Start("Demonstrating PositionBack for 5 seconds...")
	time.Sleep(3 * time.Second)
	backSpinner.Stop("")

	// 3. Keep at front, loop styles (outer) and colors (inner)
	for _, s := range styles {
		for _, c := range colors {
			desc := fmt.Sprintf("Style: %-20s | Color: %s", s.name, c.name)
			spinner := util.NewSpinner(ctx, os.Stdout, 80*time.Millisecond, util.PositionFront, s.style, c.col, util.DirectionForward)
			spinner.Start(desc)

			time.Sleep(3 * time.Second)
			spinner.Stop("")
		}
	}

	fmt.Println("=== Demo Complete ===")
}
