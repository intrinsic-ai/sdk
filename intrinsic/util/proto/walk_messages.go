// Copyright 2023 Intrinsic Innovation LLC

// Package walkmessages provides a utility for walking through proto messages recursively.
package walkmessages

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// FProcessMessage is a function that takes a message and processes it.
//
// Used in WalkRecursively. If a non-nil message is returned, the return value replaces the
// original value in the message that is being walked.
//
// The function also indicated whether to enter into the message recursively.
type fProcessMessage func(proto.Message) (proto.Message, bool, error)

// Recursively walks through a proto message, executing a function for each message it finds.
//
// The function returns whether to enter into the message recursively.
//
// The input message may be mutated, and the processed message is returned.
func Recursively(msg proto.Message, f fProcessMessage) (proto.Message, error) {
	msgOut, shouldEnter, err := f(msg)
	if err != nil {
		return nil, err
	}
	if msgOut == nil { // No changes made. Use original message.
		msgOut = msg
	}
	if !shouldEnter {
		return msgOut, nil
	}

	msgOutR := msgOut.ProtoReflect()
	for i := 0; i < msgOutR.Descriptor().Fields().Len(); i++ {
		field := msgOutR.Descriptor().Fields().Get(i)

		// Skip unspecified fields.
		if !msgOutR.Has(field) {
			continue
		}

		// Skip non-message/group types.
		if field.Kind() != protoreflect.MessageKind && field.Kind() != protoreflect.GroupKind {
			continue
		}

		valueR := msgOutR.Get(field)
		if field.IsList() { // Walk through lists.
			for i := 0; i < valueR.List().Len(); i++ {
				msgItem := valueR.List().Get(i).Message().Interface()
				if msgItemOut, err := Recursively(msgItem, f); err != nil {
					return nil, err
				} else if msgItemOut != nil { // Item was changed; update the parent.
					valueR.List().Set(i, protoreflect.ValueOfMessage(msgItemOut.ProtoReflect()))
				}
			}
		} else if field.IsMap() { // Walk through maps.
			var err error
			valueR.Map().Range(func(key protoreflect.MapKey, value protoreflect.Value) bool {
				msgItem := value.Message().Interface()
				var msgItemOut proto.Message
				if msgItemOut, err = Recursively(msgItem, f); err != nil {
					return false
				} else if msgItemOut != nil { // Item was changed; update the parent.
					valueR.Map().Set(key, protoreflect.ValueOfMessage(msgItemOut.ProtoReflect()))
				}
				return true
			})
			if err != nil {
				return nil, err
			}
		} else if valueROut, err := Recursively(valueR.Message().Interface(), f); err != nil {
			return nil, err
		} else if valueROut != nil { // Field was changed; update the parent.
			msgOutR.Set(field, protoreflect.ValueOfMessage(valueROut.ProtoReflect()))
		}
	}

	return msgOut, nil
}
