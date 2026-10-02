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

package seekable_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"

	"intrinsic/util/go/seekable"
)

type unseekableReader struct {
	r io.Reader
}

func (u *unseekableReader) Read(p []byte) (int, error) {
	return u.r.Read(p)
}

type errorReader struct {
	err error
}

func (e *errorReader) Read(p []byte) (int, error) {
	return 0, e.err
}

type op struct {
	isSeek     bool
	isZeroRead bool
	seekOffset int64
	seekWhence int
	want       string
	wantEOF    bool
}

func TestFrom(t *testing.T) {
	tests := []struct {
		name     string
		reader   io.Reader
		wantErr  bool
		wantSame bool
	}{
		{
			name:     "already seekable returns same instance",
			reader:   bytes.NewReader([]byte("test")),
			wantSame: true,
		},
		{
			name:    "nil reader returns error",
			reader:  nil,
			wantErr: true,
		},
		{
			name:   "unseekable reader wraps successfully",
			reader: &unseekableReader{r: bytes.NewReader([]byte("test"))},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rs, cleanup, err := seekable.From(tc.reader)
			if (err != nil) != tc.wantErr {
				t.Fatalf("From() err = %v, wantErr %v", err, tc.wantErr)
			}
			if tc.wantErr {
				if cleanup != nil {
					t.Error("From() returned non-nil cleanup on error")
				}
				return
			}
			if cleanup == nil {
				t.Fatal("From() returned nil cleanup on success")
			}
			t.Cleanup(func() { _ = cleanup() })

			if tc.wantSame && rs != tc.reader {
				t.Errorf("From() returned different instance, got %T, want %T", rs, tc.reader)
			}
		})
	}
}

func TestSeekAndRead(t *testing.T) {
	tests := []struct {
		data string
		name string
		ops  []op
	}{
		{
			name: "read and rewind to start",
			data: "hello world",
			ops: []op{
				{want: "hello"},
				{isSeek: true, seekOffset: 0, seekWhence: io.SeekStart},
				{want: "hello world"},
			},
		},
		{
			name: "seek from end",
			data: "0123456789",
			ops: []op{
				{isSeek: true, seekOffset: -3, seekWhence: io.SeekEnd},
				{want: "789"},
			},
		},
		{
			name: "seek to exact end",
			data: "hello world",
			ops: []op{
				{isSeek: true, seekOffset: 0, seekWhence: io.SeekEnd},
				{wantEOF: true},
			},
		},
		{
			name: "seek relative to current position forward",
			data: "0123456789",
			ops: []op{
				{want: "01"},
				{isSeek: true, seekOffset: 3, seekWhence: io.SeekCurrent},
				{want: "56"},
			},
		},
		{
			name: "seek relative to current position backward",
			data: "abcdefghij",
			ops: []op{
				{want: "abcdef"},
				{isSeek: true, seekOffset: -3, seekWhence: io.SeekCurrent},
				{want: "defghij"},
			},
		},
		{
			name: "query current position without moving",
			data: "hello world",
			ops: []op{
				{want: "hello"},
				{isSeek: true, seekOffset: 0, seekWhence: io.SeekCurrent},
				{want: " world"},
			},
		},
		{
			name: "zero-length read buffer",
			data: "hello world",
			ops: []op{
				{isZeroRead: true},
				{want: "hello"},
			},
		},
		{
			name: "seek forward within cached data",
			data: "abcdefghijklmnopqrstuvwxyz",
			ops: []op{
				{want: "abcdefghijklmnopqrst"},
				{isSeek: true, seekOffset: 0, seekWhence: io.SeekStart},
				{isSeek: true, seekOffset: 10, seekWhence: io.SeekStart},
				{want: "klmno"},
			},
		},
		{
			name: "seek forward beyond recorded data",
			data: "abcdefghijklmnopqrstuvwxyz",
			ops: []op{
				{isSeek: true, seekOffset: 10, seekWhence: io.SeekStart},
				{want: "klmno"},
			},
		},
		{
			name: "seek across multiple chunks (>32KB stream)",
			data: strings.Repeat("0123456789abcdef", 4500), // 72,000 bytes
			ops: []op{
				{isSeek: true, seekOffset: 50_000, seekWhence: io.SeekStart},
				{want: "0123456789abcdef"},
				{isSeek: true, seekOffset: -16, seekWhence: io.SeekEnd},
				{want: "0123456789abcdef"},
				{wantEOF: true},
			},
		},
		{
			name: "seek beyond EOF clamps to end",
			data: "01234",
			ops: []op{
				{isSeek: true, seekOffset: 100, seekWhence: io.SeekStart},
				{wantEOF: true},
			},
		},
		{
			name: "multiple rewinds and full replays",
			data: "abcdefghijklmnopqrstuvwxyz",
			ops: []op{
				{want: "abcdefghijklmnopqrstuvwxyz"},
				{isSeek: true, seekOffset: 0, seekWhence: io.SeekStart},
				{want: "abcdefghijklmnopqrstuvwxyz"},
				{isSeek: true, seekOffset: 0, seekWhence: io.SeekStart},
				{want: "abcdefghij"},
			},
		},
		{
			name: "empty stream",
			data: "",
			ops: []op{
				{wantEOF: true},
				{isSeek: true, seekOffset: 0, seekWhence: io.SeekStart},
				{wantEOF: true},
				{isSeek: true, seekOffset: 0, seekWhence: io.SeekEnd},
				{wantEOF: true},
			},
		},
		{
			name: "read in single byte increments with seeking",
			data: "01234",
			ops: []op{
				{want: "0"},
				{want: "1"},
				{isSeek: true, seekOffset: 0, seekWhence: io.SeekStart},
				{want: "0"},
				{want: "1"},
				{want: "2"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name+"/unseekable", func(t *testing.T) {
			raw := &unseekableReader{r: bytes.NewReader([]byte(tc.data))}
			rs, cleanup, err := seekable.From(raw)
			if err != nil {
				t.Fatalf("From failed: %v", err)
			}
			t.Cleanup(func() { _ = cleanup() })

			runOps(t, rs, tc.ops)
		})

		t.Run(tc.name+"/seekable", func(t *testing.T) {
			rs, cleanup, err := seekable.From(bytes.NewReader([]byte(tc.data)))
			if err != nil {
				t.Fatalf("From failed: %v", err)
			}
			t.Cleanup(func() { _ = cleanup() })

			runOps(t, rs, tc.ops)
		})
	}
}

func runOps(t *testing.T, rs io.ReadSeeker, ops []op) {
	t.Helper()

	for i, o := range ops {
		if o.isSeek {
			if _, err := rs.Seek(o.seekOffset, o.seekWhence); err != nil {
				t.Fatalf("op %d (seek offset %d, whence %d) failed: %v", i, o.seekOffset, o.seekWhence, err)
			}
			continue
		}

		if o.isZeroRead {
			n, err := rs.Read([]byte{})
			if n != 0 || err != nil {
				t.Fatalf("op %d (zero-length read) got (%d, %v), want (0, nil)", i, n, err)
			}
			continue
		}

		if o.wantEOF {
			buf := make([]byte, 1)
			n, err := rs.Read(buf)
			if err != io.EOF {
				t.Fatalf("op %d (read at EOF) got err %v (n=%d), want io.EOF", i, err, n)
			}
			continue
		}

		buf := make([]byte, len(o.want))
		n, err := io.ReadFull(rs, buf)
		if err != nil {
			t.Fatalf("op %d (read %d bytes) failed: %v", i, len(o.want), err)
		}
		if got := string(buf[:n]); got != o.want {
			t.Fatalf("op %d got %q, want %q", i, got, o.want)
		}
	}
}

func TestErrors(t *testing.T) {
	errSource := errors.New("underlying read failed")

	tests := []struct {
		closeFirst bool
		isClose    bool
		isRead     bool
		name       string
		offset     int64
		reader     io.Reader
		wantErr    error
		whence     int
	}{
		{
			name:   "negative seek from start",
			offset: -1,
			whence: io.SeekStart,
		},
		{
			name:   "negative seek from current",
			offset: -5,
			whence: io.SeekCurrent,
		},
		{
			name:   "negative seek from end",
			offset: -100,
			whence: io.SeekEnd,
		},
		{
			name:   "invalid whence",
			offset: 0,
			whence: 999,
		},
		{
			name:       "seek after cleanup",
			closeFirst: true,
			offset:     0,
			whence:     io.SeekStart,
		},
		{
			name:       "read after cleanup",
			closeFirst: true,
			isRead:     true,
		},
		{
			name:       "double cleanup succeeds",
			closeFirst: true,
			isClose:    true,
		},
		{
			name:    "underlying source read error",
			isRead:  true,
			reader:  &errorReader{err: errSource},
			wantErr: errSource,
		},
		{
			name:    "underlying source seek error",
			offset:  10,
			reader:  &errorReader{err: errSource},
			wantErr: errSource,
			whence:  io.SeekStart,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.reader
			if r == nil {
				r = &unseekableReader{r: bytes.NewReader([]byte("hello world"))}
			}
			rs, cleanup, err := seekable.From(r)
			if err != nil {
				t.Fatalf("From failed: %v", err)
			}

			if tc.closeFirst {
				if err := cleanup(); err != nil {
					t.Fatalf("cleanup failed: %v", err)
				}
			} else {
				t.Cleanup(func() { _ = cleanup() })
			}

			if tc.isClose {
				if err := cleanup(); err != nil {
					t.Errorf("second cleanup failed: %v", err)
				}
				return
			}

			if tc.isRead {
				buf := make([]byte, 5)
				_, err := rs.Read(buf)
				if tc.wantErr != nil {
					if !errors.Is(err, tc.wantErr) {
						t.Errorf("Read err = %v, want %v", err, tc.wantErr)
					}
				} else if err == nil {
					t.Error("Read succeeded, want error")
				}
				return
			}

			_, err = rs.Seek(tc.offset, tc.whence)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("Seek err = %v, want %v", err, tc.wantErr)
				}
			} else if err == nil {
				t.Errorf("Seek(%d, %d) succeeded, want error", tc.offset, tc.whence)
			}
		})
	}
}

func TestConcurrent(t *testing.T) {
	raw := &unseekableReader{r: bytes.NewReader([]byte("abcdefghijklmnopqrstuvwxyz0123456789"))}
	rs, cleanup, err := seekable.From(raw)
	if err != nil {
		t.Fatalf("From failed: %v", err)
	}
	t.Cleanup(func() { _ = cleanup() })

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := make([]byte, 4)
			for j := 0; j < 20; j++ {
				_, _ = rs.Seek(0, io.SeekStart)
				_, _ = rs.Read(buf)
				_, _ = rs.Seek(2, io.SeekCurrent)
			}
		}()
	}
	wg.Wait()
}
