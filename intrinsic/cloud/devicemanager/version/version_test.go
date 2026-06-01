// Copyright 2023 Intrinsic Innovation LLC

package version

import "testing"

const (
	inversionOS   = "xfa.20241221.RC01"                      // xfa.YYYYMMDD.RCXX
	uiOS          = "20241221.RC01"                          // YYYYMMDD.RCXX
	apiOS         = "0.0.1+xfa.20241221.RC01"                // 0.0.1+xfa.YYYYMMDD.RCXX
	inversionBase = "20250121.rc00"                          // YYYYMMDD.rcXX
	uiBase        = "20250121.RC00"                          // YYYYMMDD.RCXX
	apiBase       = "0.0.1+intrinsic.platform.20250121.RC00" // 0.0.1+intrinsic.platform.YYYYMMDD.RCXX
)

var (
	osI2U = TranslateOSInversionToUI
	osU2I = TranslateOSUIToInversion
	osA2U = TranslateOSAPIToUI
	osU2A = TranslateOSUIToAPI
	osA2I = TranslateOSAPIToInversion
	osI2A = TranslateOSInversionToAPI

	baseI2U = TranslateBaseInversionToUI
	baseU2I = TranslateBaseUIToInversion
	baseA2U = TranslateBaseAPIToUI
	baseU2A = TranslateBaseUIToAPI
	baseA2I = TranslateBaseAPIToInversion
	baseI2A = TranslateBaseInversionToAPI
)

func TestInverse(t *testing.T) {
	if got, want := osU2I(osI2U(inversionOS)), inversionOS; got != want {
		t.Errorf("TranslateOSUIToInversion is not the inverse of TranslateOSInversionToUI. got: %q want: %q", got, want)
	}
	if got, want := osA2I(osI2A(inversionOS)), inversionOS; got != want {
		t.Errorf("TranslateOSAPIToInversion is not the inverse of TranslateOSInversionToAPI. got: %q want: %q", got, want)
	}
	if got, want := osU2A(osA2U(apiOS)), apiOS; got != want {
		t.Errorf("TranslateOSAPIToUI is not the inverse of TranslateOSUIToAPI. got: %q want: %q", got, want)
	}
	if got, want := baseU2I(baseI2U(inversionBase)), inversionBase; got != want {
		t.Errorf("TranslateBaseUIToInversion is not the inverse of TranslateBaseInversionToUI. got: %q want: %q", got, want)
	}
	if got, want := baseA2I(baseI2A(inversionBase)), inversionBase; got != want {
		t.Errorf("TranslateBaseAPIToInversion is not the inverse of TranslateBaseInversionToAPI. got: %q want: %q", got, want)
	}
	if got, want := baseA2U(baseU2A(uiBase)), uiBase; got != want {
		t.Errorf("TranslateBaseAPIToUI is not the inverse of TranslateBaseUIToAPI. got: %q want: %q", got, want)
	}
}

func TestOS(t *testing.T) {
	if got, want := osU2I(uiOS), inversionOS; got != want {
		t.Errorf("TranslateOSUIToInversion(%q)=%q, want %q", uiOS, got, want)
	}
	if got, want := osI2U(inversionOS), uiOS; got != want {
		t.Errorf("TranslateOSInversionToUI(%q)=%q, want %q", inversionOS, got, want)
	}
	if got, want := osA2U(apiOS), uiOS; got != want {
		t.Errorf("TranslateOSAPIToUI(%q)=%q, want %q", apiOS, got, want)
	}
	if got, want := osU2A(uiOS), apiOS; got != want {
		t.Errorf("TranslateOSUIToAPI(%q)=%q, want %q", uiOS, got, want)
	}
	if got, want := osA2I(apiOS), inversionOS; got != want {
		t.Errorf("TranslateOSAPIToInversion(%q)=%q, want %q", apiOS, got, want)
	}
	if got, want := osI2A(inversionOS), apiOS; got != want {
		t.Errorf("TranslateOSInversionToAPI(%q)=%q, want %q", inversionOS, got, want)
	}
}

func TestBase(t *testing.T) {
	if got, want := baseU2I(uiBase), inversionBase; got != want {
		t.Errorf("TranslateBaseUIToInversion(%q)=%q, want %q", uiBase, got, want)
	}
	if got, want := baseI2U(inversionBase), uiBase; got != want {
		t.Errorf("TranslateBaseInversionToUI(%q)=%q, want %q", inversionBase, got, want)
	}
	if got, want := baseA2U(apiBase), uiBase; got != want {
		t.Errorf("TranslateBaseAPIToUI(%q)=%q, want %q", apiBase, got, want)
	}
	if got, want := baseU2A(uiBase), apiBase; got != want {
		t.Errorf("TranslateBaseUIToAPI(%q)=%q, want %q", uiBase, got, want)
	}
	if got, want := baseA2I(apiBase), inversionBase; got != want {
		t.Errorf("TranslateBaseAPIToInversion(%q)=%q, want %q", apiBase, got, want)
	}
	if got, want := baseI2A(inversionBase), apiBase; got != want {
		t.Errorf("TranslateBaseInversionToAPI(%q)=%q, want %q", inversionBase, got, want)
	}
}

func TestEmptyBase(t *testing.T) {
	if got, want := TranslateBaseInversionToUI(""), ""; got != want {
		t.Errorf("TranslateBaseInversionToUI(%q)=%q, want %q", want, got, want)
	}
	if got, want := TranslateBaseUIToInversion(""), ""; got != want {
		t.Errorf("TranslateBaseUIToInversion(%q)=%q, want %q", want, got, want)
	}
	if got, want := TranslateBaseAPIToUI(""), ""; got != want {
		t.Errorf("TranslateBaseAPIToUI(%q)=%q, want %q", want, got, want)
	}
	if got, want := TranslateBaseInversionToAPI(""), ""; got != want {
		t.Errorf("TranslateBaseInversionToAPI(%q)=%q, want %q", want, got, want)
	}
	if got, want := TranslateBaseUIToAPI(""), ""; got != want {
		t.Errorf("TranslateBaseUIToAPI(%q)=%q, want %q", want, got, want)
	}
	if got, want := TranslateBaseAPIToInversion(""), ""; got != want {
		t.Errorf("TranslateBaseAPIToInversion(%q)=%q, want %q", want, got, want)
	}
}

func TestEmptyOS(t *testing.T) {
	if got, want := TranslateOSInversionToUI(""), ""; got != want {
		t.Errorf("TranslateOSInversionToUI(%q)=%q, want %q", want, got, want)
	}
	if got, want := TranslateOSUIToInversion(""), ""; got != want {
		t.Errorf("TranslateBasTranslateOSUIToInversioneInversionToUI(%q)=%q, want %q", want, got, want)
	}
	if got, want := TranslateOSAPIToUI(""), ""; got != want {
		t.Errorf("TranslateOSAPIToUI(%q)=%q, want %q", want, got, want)
	}
	if got, want := TranslateOSUIToAPI(""), ""; got != want {
		t.Errorf("TranslateOSUIToAPI(%q)=%q, want %q", want, got, want)
	}
	if got, want := TranslateOSAPIToInversion(""), ""; got != want {
		t.Errorf("TranslateOSAPIToInversion(%q)=%q, want %q", want, got, want)
	}
	if got, want := TranslateOSInversionToAPI(""), ""; got != want {
		t.Errorf("TranslateOSInversionToAPI(%q)=%q, want %q", want, got, want)
	}
}
