// Copyright 2023 Intrinsic Innovation LLC

package icon

import (
	"fmt"
	"strings"
)

const (
	stateVariablePathPrefix    = "@"
	stateVariablePathSeparator = "."
	armPartNodeName            = "ArmPart"
	ftPartNodeName             = "ForceTorqueSensorPart"
	gripperPartNodeName        = "GripperPart"
	rangefinderPartNodeName    = "RangefinderPart"
	adioPartNodeName           = "ADIOPart"
	safetyNodeName             = "Safety"
)

// stateVariablePathNode represents one node of a StateVariablePath consisting of a node name and an optional index.
type stateVariablePathNode struct {
	name  string
	index *uint64
}

// stateVariablePathBuilder is a struct that helps building state variable paths by adding nodes that construct in the end the complete path.
type stateVariablePathBuilder struct {
	nodes []stateVariablePathNode
}

// newStateVariablePathBuilder creates a new [stateVariablePathBuilder].
func newStateVariablePathBuilder() *stateVariablePathBuilder {
	return &stateVariablePathBuilder{}
}

// addNodeWithIndex adds a single node with the given `nodeName` and an array index to the path.
func (p *stateVariablePathBuilder) addNodeWithIndex(nodeName string, index uint64) *stateVariablePathBuilder {
	p.nodes = append(p.nodes, stateVariablePathNode{
		name:  nodeName,
		index: &index,
	})
	return p
}

// addNode adds one node with the given `nodeNames` to the path.
func (p *stateVariablePathBuilder) addNode(nodeName string) *stateVariablePathBuilder {
	p.nodes = append(p.nodes, stateVariablePathNode{
		name:  nodeName,
		index: nil,
	})

	return p
}

// addNodes adds multiple nodes with the given `nodeNames` to the path.
func (p *stateVariablePathBuilder) addNodes(nodeNames ...string) *stateVariablePathBuilder {
	for _, n := range nodeNames {
		p = p.addNode(n)
	}
	return p
}

// build builds the final state variable path string based on the previous add-calls.
func (p *stateVariablePathBuilder) build() string {
	nodeStrings := []string{}
	for _, node := range p.nodes {
		nodeString := node.name
		if node.index != nil {
			nodeString += fmt.Sprintf("[%d]", *node.index)
		}
		nodeStrings = append(nodeStrings, nodeString)
	}
	return stateVariablePathPrefix + strings.Join(nodeStrings, stateVariablePathSeparator)
}

// ArmSensedPosition generates a state variable path for a sensed position of the joint at `jointIndex` of the part called `part_name`.
func ArmSensedPosition(partName string, jointIndex uint64) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, armPartNodeName)
	builder.addNodeWithIndex("sensed_position", jointIndex)
	return builder.build()
}

// ArmSensedVelocity generates a state variable path for the sensed velocity of the joint at jointIndex of the part called part_name.
func ArmSensedVelocity(partName string, jointIndex uint64) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, armPartNodeName)
	builder.addNodeWithIndex("sensed_velocity", jointIndex)
	return builder.build()
}

// ArmSensedAcceleration generates a state variable path for the sensed acceleration of the joint at jointIndex of the part called part_name.
func ArmSensedAcceleration(partName string, jointIndex uint64) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, armPartNodeName)
	builder.addNodeWithIndex("sensed_acceleration", jointIndex)
	return builder.build()
}

// ArmSensedTorque generates a state variable path for the sensed torque of the joint at jointIndex of the part called part_name.
func ArmSensedTorque(partName string, jointIndex uint64) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, armPartNodeName)
	builder.addNodeWithIndex("sensed_torque", jointIndex)
	return builder.build()
}

// TwistDimension is a type redefinition of uint64 for more convenient twist dimension specification.
type TwistDimension uint64

// Constants for twist dimensions. X,Y,Z are the linear velocities, RX,RY,RZ are the angular velocities.
const (
	TwistX  TwistDimension = 0
	TwistY  TwistDimension = 1
	TwistZ  TwistDimension = 2
	TwistRX TwistDimension = 3
	TwistRY TwistDimension = 4
	TwistRZ TwistDimension = 5
)

// ArmBaseTwistTipSensed generates a state variable path for the sensed twist at twistIndex of the part called part_name.
func ArmBaseTwistTipSensed(partName string, twistDimension TwistDimension) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, armPartNodeName)
	builder.addNodeWithIndex("base_twist_tip_sensed", uint64(twistDimension))
	return builder.build()
}

// ArmBaseLinearVelocityTipSensed generates a state variable path for the linear (translational) velocity of the arm at the tip.
func ArmBaseLinearVelocityTipSensed(partName string) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, armPartNodeName, "base_linear_velocity_tip_sensed")
	return builder.build()
}

// ArmBaseAngularVelocityTipSensed generates a state variable path for the linear (translational) velocity of the arm at the tip.
func ArmBaseAngularVelocityTipSensed(partName string) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, armPartNodeName, "base_angular_velocity_tip_sensed")
	return builder.build()
}

// ArmCurrentControlMode generates a state variable path for the current control mode of the arm.
func ArmCurrentControlMode(partName string) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, armPartNodeName, "current_control_mode")

	return builder.build()
}

// WrenchDimension is a type redefinition of uint64 for more convenient wrench dimension specification.
type WrenchDimension uint64

// Constants for wrench dimensions. X,Y,Z are the force, RX,RY,RZ are the torques.
const (
	WrenchX  WrenchDimension = 0
	WrenchY  WrenchDimension = 1
	WrenchZ  WrenchDimension = 2
	WrenchRX WrenchDimension = 3
	WrenchRY WrenchDimension = 4
	WrenchRZ WrenchDimension = 5
)

// FTWrenchAtTip generates a state variable path for the wrench as the tip sensed in the force torque sensor.
func FTWrenchAtTip(partName string, wrenchDimension WrenchDimension) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, ftPartNodeName)
	builder.addNodeWithIndex("wrench_at_tip", uint64(wrenchDimension))
	return builder.build()
}

// FTForceMagnitudeAtTip generates a state variable path for the magnitude of the force sensed at the force torque sensor in the frame of the arm tip.
func FTForceMagnitudeAtTip(partName string) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, ftPartNodeName, "force_magnitude_at_tip")
	return builder.build()
}

// FTTorqueMagnitudeAtTip generates a state variable path for the magnitude of the torque sensed at the force torque sensor in the frame of the arm tip.
func FTTorqueMagnitudeAtTip(partName string) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, ftPartNodeName, "torque_magnitude_at_tip")
	return builder.build()
}

// FTWrenchStabilityIndex generates a state variable path for the stability index of the sensed wrench at the force torque sensor.
func FTWrenchStabilityIndex(partName string) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, ftPartNodeName, "wrench_stability_index")
	return builder.build()
}

// GripperSensedState generates a state variable path for the current state of the gripper.
func GripperSensedState(partName string) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, gripperPartNodeName, "sensed_state")
	return builder.build()
}

// GripperOpeningWidth generates a state variable path for the opening width of the gripper.
func GripperOpeningWidth(partName string) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, gripperPartNodeName, "opening_width")
	return builder.build()
}

// RangefinderDistance generates a state variable path for the measured distance of the rangefinder.
func RangefinderDistance(partName string) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, rangefinderPartNodeName, "distance")
	return builder.build()
}

// ADIODigitalInput creates a state variable path for a digital input of the signal at `signalIndex` in block `blockName` of the part called `partName`.
func ADIODigitalInput(partName string, blockName string, signalIndex uint64) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, adioPartNodeName, "di")
	builder.addNodeWithIndex(blockName, signalIndex)
	return builder.build()
}

// ADIODigitalOutput creates a state variable path for a digital output of the signal at `signalIndex` in block `blockName` of the part called `partName`.
func ADIODigitalOutput(partName string, blockName string, signalIndex uint64) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, adioPartNodeName, "do")
	builder.addNodeWithIndex(blockName, signalIndex)
	return builder.build()
}

// ADIOAnalogInput creates a state variable path for an analog input of the signal at `signalIndex` in block `blockName` of the part called `partName`.
func ADIOAnalogInput(partName string, blockName string, signalIndex uint64) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, adioPartNodeName, "ai")
	builder.addNodeWithIndex(blockName, signalIndex)
	return builder.build()
}

// ADIOAnalogOutput creates a state variable path for an analog output of the signal at `signalIndex` in block `blockName` of the part called `partName`.
func ADIOAnalogOutput(partName string, blockName string, signalIndex uint64) string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(partName, adioPartNodeName, "ao")
	builder.addNodeWithIndex(blockName, signalIndex)
	return builder.build()
}

// SafetyEnableButtonStatus generates a state variable path for status of the enable button.
func SafetyEnableButtonStatus() string {
	builder := newStateVariablePathBuilder()
	builder.addNodes(safetyNodeName, "enable_button_status")
	return builder.build()
}
