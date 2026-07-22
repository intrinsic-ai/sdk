// Copyright 2023 Intrinsic Innovation LLC

package icon

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	typespb "intrinsic/icon/proto/v1/types_go_proto"

	anypb "google.golang.org/protobuf/types/known/anypb"
)

// ActionID is an int64 used to identify a server-side action instance.
type ActionID int64

// ActionHandle is a reference to an action instance. Client-side code should
// treat this as an opaque handle. The default zero-valued ActionHandle is
// considered invalid. Handles are created client-side before the corresponding
// action instance is created server-side, allowing for cyclic graphs of action
// nodes. ActionHandle objects may be safely compared for equality (==, !=),
// copied by value, and used as keys in a map.
type ActionHandle struct {
	id ActionID
}

// ID returns the associated action ID, which is used by the server to identify
// the referenced action instance.
func (h ActionHandle) ID() ActionID {
	return h.id
}

// IsZero reports whether h is zero-valued, which represents an invalid action
// instance reference.
func (h ActionHandle) IsZero() bool {
	return h.id == 0
}

// SlotData is a wrapper that can hold either
//   - a string (a single ICON part name)
//   - a map[string]string (a map from ICON slot names to ICON part names)
//
// Only one of the se two fields may be set at a time.
// [FromSlotPartMap] and [FromPartName] enforce this.
type SlotData struct {
	// If this is the empty string, that counts as "unset".
	partName    string
	slotPartMap map[string]string
}

// FromSlotPartMap makes a SlotData object that holds `slotPartMap`.
func FromSlotPartMap(slotPartMap map[string]string) SlotData {
	return SlotData{partName: "", slotPartMap: slotPartMap}
}

// FromPartName makes a SlotData object that holds `partName`.
func FromPartName(partName string) SlotData {
	return SlotData{partName: partName, slotPartMap: nil}
}

// ActionDescription is used to construct a new Action instance via
// session.AddAction.
type ActionDescription struct {
	// Handle identifies this action, created with session.MakeActionHandle().
	Handle     ActionHandle
	ActionType string
	Params     proto.Message
	SlotData   SlotData
	Reactions  []*Reaction
}

// proto converts ad to an ActionInstance proto.
func (ad *ActionDescription) proto() (*typespb.ActionInstance, error) {
	var anyParams *anypb.Any
	if ad.Params != nil {
		var err error
		if anyParams, err = anypb.New(ad.Params); err != nil {
			return nil, err
		}
	}

	if (ad.SlotData.partName == "") && (ad.SlotData.slotPartMap == nil) {
		return nil, fmt.Errorf("an ActionDescription should have either a part name or a SlotPartMap, but we have neither")
	}
	if (ad.SlotData.partName != "") && (ad.SlotData.slotPartMap != nil) {
		return nil, fmt.Errorf("an ActionDescription should have either a part name or a SlotPartMap, but we have both")
	}

	if ad.SlotData.partName != "" {
		return &typespb.ActionInstance{
			ActionInstanceId: int64(ad.Handle.ID()),
			ActionTypeName:   ad.ActionType,
			FixedParameters:  anyParams,
			SlotData: &typespb.ActionInstance_PartName{
				PartName: ad.SlotData.partName,
			},
		}, nil
	}

	// SlotPartMap case
	return &typespb.ActionInstance{
		ActionInstanceId: int64(ad.Handle.ID()),
		ActionTypeName:   ad.ActionType,
		FixedParameters:  anyParams,
		SlotData: &typespb.ActionInstance_SlotPartMap{SlotPartMap: &typespb.SlotPartMap{
			SlotNameToPartName: ad.SlotData.slotPartMap,
		}},
	}, nil
}
