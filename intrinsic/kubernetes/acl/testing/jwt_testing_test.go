// Copyright 2023 Intrinsic Innovation LLC

package jwttesting_test

import (
	"testing"
	"time"

	"intrinsic/kubernetes/acl/testing/jwttesting"
)

func TestMintToken(t *testing.T) {
	opts := jwttesting.Options{
		jwttesting.WithEmail("doe@gmail.com"),
		jwttesting.WithExpiresAt(time.Unix(420, 0)),
		jwttesting.WithIssuedAt(time.Unix(123, 456)),
		jwttesting.WithSigningKey("hello Intrinsic"),
		jwttesting.WithSubject("on the significance of subjects"),
		jwttesting.WithUID("snowflake"),
	}
	if got, want := jwttesting.MintToken(t, opts...), "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImRvZUBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsInVpZCI6InNub3dmbGFrZSIsImV4cCI6NDIwLCJpYXQiOjEyMywic3ViIjoib24gdGhlIHNpZ25pZmljYW5jZSBvZiBzdWJqZWN0cyJ9.YTsLPYu6_FEgxd2IlPHSybdgzY-PeCBNyuOE-Bkz9eM"; got != want {
		t.Errorf("jwttesting.MintToken(t, opts...) = %q, want %q", got, want)
	}
}

func TestMustMintToken(t *testing.T) {
	opts := jwttesting.Options{
		jwttesting.WithEmail("edo@google.com"),
		jwttesting.WithExpiresAt(time.Unix(420, 0)),
		jwttesting.WithIssuedAt(time.Unix(456, 123)),
		jwttesting.WithSigningKey("goodbye Intrinsic"),
		jwttesting.WithAudience("intergalactic robots federation"),
	}
	if got, want := jwttesting.MustMintToken(opts...), "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImVkb0Bnb29nbGUuY29tIiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJ1aWQiOiIiLCJhdWQiOiJpbnRlcmdhbGFjdGljIHJvYm90cyBmZWRlcmF0aW9uIiwiZXhwIjo0MjAsImlhdCI6NDU2fQ.g7xKejLnzRDhnYGwDKN8FtKDZqjJ09yUZbBTTrfsfdk"; got != want {
		t.Errorf("jwttesting.MustMintToken(opts...) = %q, want %q", got, want)
	}
}
