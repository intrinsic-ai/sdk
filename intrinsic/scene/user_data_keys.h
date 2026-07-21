// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_SCENE_USER_DATA_KEYS_H_
#define INTRINSIC_SCENE_USER_DATA_KEYS_H_

// This file contain keys of user data World->SDF conversion might generate.
// The values format are documented for each key below. Additionally, these keys
// may be considered as restricted when populating user_data in SceneObject
// proto.
//
namespace intrinsic {
namespace sdf {

// The value is Gazebo plugins associated with a World Entity.
// The format of the value is <plugin .../><plugin .../>
// The <plugin> tags follow the same format as that of SDF.
inline constexpr char kGazeboPlugins[] = "gazebo_plugins";

// The value is the text representation of <surface> element under
// <link><collision> in SDF. It saves contact physics parameters for a World
// entity.
inline constexpr char kGazeboCollisionSurface[] = "gazebo_collision_surface";

// The value is the text representation of <physics> element under <joint> in
// SDF. It saves physics parameters for a joint entity.
// Note: the kGazeboJointPhysics key is deprecated as it is no longer used
// during the sdf -> scene object conversion process.
inline constexpr char kGazeboJointPhysics[] = "gazebo_joint_physics";

// User data associated with a Scene Object or Collection entity, to store
// joint simulation properties that are not otherwise represented in the World.
// The stored data is a map from joint name to the text representation of the
// `<joint>` element within the source `<model>` element in the source SDF file.
inline constexpr char kGazeboCustomJoint[] = "gazebo_custom_joint";

// The value is the text representation of <light> elements under a <world> or
// <link> element in SDF.
inline constexpr char kSdfLights[] = "sdf_lights";

// The value is an intrinsic_proto.icon.DigitalInputOutput proto that describes
// the DIO blocks for an entity
//
// Note that unlike the other keys in this file, this is a key for the newer
// user_data_protos field!
inline constexpr char kDioData[] = "digital_input_output";
}  // namespace sdf
}  // namespace intrinsic

#endif  // INTRINSIC_SCENE_USER_DATA_KEYS_H_
