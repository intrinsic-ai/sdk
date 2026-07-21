// Copyright 2023 Intrinsic Innovation LLC

package sceneobjectvalidation

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
	"unsafe"

	epb "intrinsic/scene/proto/v1/entity_go_proto"
	sopb "intrinsic/scene/proto/v1/scene_object_go_proto"
)

/*
#include <stdlib.h>
#include "intrinsic/scene/validate/go/scene_object_validation_c.h"
*/
import "C"

func TestValidateSceneObject(t *testing.T) {
	tests := []struct {
		desc       string
		object     *sopb.SceneObject
		wantErr    bool
		wantCode   int
		wantSubstr string
	}{
		{
			desc: "valid scene object",
			object: &sopb.SceneObject{
				Name: "valid_object",
				Entities: []*epb.Entity{
					{
						Name:       "root",
						EntityType: &epb.Entity_Link{},
					},
				},
			},
			wantErr: false,
		},
		{
			desc: "invalid: empty entity name",
			object: &sopb.SceneObject{
				Name: "invalid_object",
				Entities: []*epb.Entity{
					{
						Name:       "",
						EntityType: &epb.Entity_Link{},
					},
				},
			},
			wantErr:    true,
			wantCode:   3, // InvalidArgument
			wantSubstr: "no name",
		},
		{
			desc: "invalid: duplicate entity name",
			object: &sopb.SceneObject{
				Name: "invalid_object",
				Entities: []*epb.Entity{
					{
						Name:       "root",
						EntityType: &epb.Entity_Link{},
					},
					{
						Name:       "root",
						EntityType: &epb.Entity_Link{},
					},
				},
			},
			wantErr:    true,
			wantCode:   3, // InvalidArgument
			wantSubstr: "Duplicate entity name",
		},
		{
			desc: "invalid: root entity is not a link",
			object: &sopb.SceneObject{
				Name: "invalid_object",
				Entities: []*epb.Entity{
					{
						Name:       "root",
						EntityType: &epb.Entity_Frame{},
					},
				},
			},
			wantErr:    true,
			wantCode:   3, // InvalidArgument
			wantSubstr: "must be a link",
		},
	}

	errRegex := regexp.MustCompile(`^validation failed with code (\d+): (.*)$`)

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := ValidateSceneObject(tc.object)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ValidateSceneObject() error = %v, wantErr %v", err, tc.wantErr)
			}
			if tc.wantErr {
				errStr := err.Error()
				matches := errRegex.FindStringSubmatch(errStr)
				if len(matches) != 3 {
					t.Fatalf("error message %q did not match expected pattern", errStr)
				}

				code, parseErr := strconv.Atoi(matches[1])
				if parseErr != nil {
					t.Fatalf("failed to parse code from error: %v", parseErr)
				}
				msg := matches[2]

				if code != tc.wantCode {
					t.Errorf("Error code = %d, want %d", code, tc.wantCode)
				}
				if !strings.Contains(strings.ToLower(msg), strings.ToLower(tc.wantSubstr)) {
					t.Errorf("Error message = %q, want substring (case-insensitive) %q", msg, tc.wantSubstr)
				}
			}
		})
	}
}

func TestValidateSceneObject_CAPI(t *testing.T) {
	dummy := byte(0)
	nonNullPtr := (*C.char)(unsafe.Pointer(&dummy))

	tests := []struct {
		desc        string
		protoData   *C.char
		protoLen    C.int
		wantCode    int
		wantSubstr  string
		checkErrMsg bool
	}{
		{
			desc:        "nil proto_data with positive length",
			protoData:   nil,
			protoLen:    10,
			wantCode:    3, // InvalidArgument
			wantSubstr:  "proto_data is null",
			checkErrMsg: true,
		},
		{
			desc:        "nil proto_data with zero length",
			protoData:   nil,
			protoLen:    0,
			wantCode:    3, // InvalidArgument
			checkErrMsg: false,
		},
		{
			desc:        "non-nil proto_data with negative length",
			protoData:   nonNullPtr,
			protoLen:    -1,
			wantCode:    3, // InvalidArgument
			checkErrMsg: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			var c_err_msg *C.char
			var p_err_msg **C.char
			if tc.checkErrMsg {
				p_err_msg = &c_err_msg
			}
			defer func() {
				if c_err_msg != nil {
					C.free(unsafe.Pointer(c_err_msg))
				}
			}()

			status_code := C.intrinsic_scene_object_go_ValidateSceneObject(tc.protoData, tc.protoLen, p_err_msg)
			if int(status_code) != tc.wantCode {
				t.Errorf("Expected status code %d, got %d", tc.wantCode, status_code)
			}
			if tc.checkErrMsg && tc.wantSubstr != "" {
				if c_err_msg == nil {
					t.Fatal("Expected error message, got nil")
				}
				err_str := C.GoString(c_err_msg)
				if !strings.Contains(strings.ToLower(err_str), strings.ToLower(tc.wantSubstr)) {
					t.Errorf("Error message = %q, want substring %q", err_str, tc.wantSubstr)
				}
			}
		})
	}
}
