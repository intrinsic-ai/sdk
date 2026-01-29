// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_WORLD_ENTITY_ID_H_
#define INTRINSIC_WORLD_ENTITY_ID_H_

// This file describes an Entity Id and typed entity id concept for use with the
// World. EntityId is a strongly typed integral value that uniquely identifies
// an entity within a World. The TypedEntityId types are wrappers around this
// strongly typed id with some extra information and type safety.
//
// TypedEntityId is the base that provides base functionality and safety. It
// has a set of template parameters that describe which components this entity
// id has. You can implicitly convert one TypedEntityId to another as long
// as the conversion is a subset operation of the entity component types.
//
// There exist a few premade typed entity id types for convenience
//   * one for each component type (e.g. AttachmentEntityId, CollisionEntityId)
//   * one for each logical type (e.g. JointId, LinkId)
//
// The intention of the TypedEntityId type is to be constructed by the World
// and during construction we can do a runtime check to make sure the entity has
// the correct types. World is also the place that we will consume these types
// and perform actions with them. At that time we can also do runtime validation
// around the type and components the entity contains.

#include <cstdint>
#include <ostream>
#include <type_traits>
#include <utility>

#include "intrinsic/production/external/intops/strong_int.h"
#include "intrinsic/util/aggregate_type.h"

namespace intrinsic {

// The strongly typed integral value that holds the entity id.
DEFINE_STRONG_INT_TYPE(EntityId, uint32_t);

// EntityId == 0 is invalid.
constexpr EntityId kInvalidEntityId = EntityId(0);

// The first available entity id in a new World object.
extern const EntityId kFirstEntityId;

// Default prefix used for entity id generation
constexpr uint16_t kDefaultEntityIdPrefix = 0;

inline bool operator==(const EntityId& lhs, const EntityId& rhs) {
  return lhs.value() == rhs.value();
}

inline bool operator!=(const EntityId& lhs, const EntityId& rhs) {
  return lhs.value() != rhs.value();
}

inline bool operator>(const EntityId& lhs, const EntityId& rhs) {
  return lhs.value() > rhs.value();
}

inline bool operator>=(const EntityId& lhs, const EntityId& rhs) {
  return lhs.value() >= rhs.value();
}

inline bool operator<(const EntityId& lhs, const EntityId& rhs) {
  return lhs.value() < rhs.value();
}

inline bool operator<=(const EntityId& lhs, const EntityId& rhs) {
  return lhs.value() <= rhs.value();
}

// Base struct for all TypedEntityId that will hold the only data needed within
// a TypedEntityId besides its type which is the id of the entity.
struct EntityBase {
  explicit EntityBase(EntityId _id) : id(_id) {}
  explicit EntityBase(EntityId::ValueType _id) : id(_id) {}

  bool operator==(const EntityBase& other) const { return id == other.id; }

  EntityId id;
};

// These structs are intended as strong types to identify the corresponding
// components but should not contain any actual component data themselves.
struct AttachmentComponentType {
  bool operator==(const AttachmentComponentType& other) const { return true; }
};
struct CollectionsComponentType {
  bool operator==(const CollectionsComponentType& other) const { return true; }
};
struct CollectionsMemberComponentType {
  bool operator==(const CollectionsMemberComponentType& other) const {
    return true;
  }
};
struct CollisionComponentType {
  bool operator==(const CollisionComponentType& other) const { return true; }
};
struct EquipmentComponentType {
  bool operator==(const EquipmentComponentType& other) const { return true; }
};
struct GeometryComponentType {
  bool operator==(const GeometryComponentType& other) const { return true; }
};
struct KinematicsComponentType {
  bool operator==(const KinematicsComponentType& other) const { return true; }
};
struct GripperComponentType {
  bool operator==(const GripperComponentType& other) const { return true; }
};
struct OutfeedComponentType {
  bool operator==(const OutfeedComponentType& other) const { return true; }
};
struct PhysicsComponentType {
  bool operator==(const PhysicsComponentType& other) const { return true; }
};
struct PPRComponentType {
  bool operator==(const PPRComponentType& other) const { return true; }
};
struct ProjectorComponentType {
  bool operator==(const ProjectorComponentType& other) const { return true; }
};
struct RegionsComponentType {
  bool operator==(const RegionsComponentType& other) const { return true; }
};
struct RobotComponentType {
  bool operator==(const RobotComponentType& other) const { return true; }
};
struct UserDataComponentType {
  bool operator==(const UserDataComponentType& other) const { return true; }
};
struct SensorComponentType {
  bool operator==(const SensorComponentType& other) const { return true; }
};
struct SimulationComponentType {
  bool operator==(const SimulationComponentType& other) const { return true; }
};
struct SpawnerComponentType {
  bool operator==(const SpawnerComponentType& other) const { return true; }
};

// Base class for all of the typed Entity id c++ types, always specifies
// EntityBase.
template <typename... ComponentTypes>
struct TypedEntityId : public AggregateType<EntityBase, ComponentTypes...> {
  static_assert(std::conjunction<std::is_empty<ComponentTypes>...>::value,
                "Must have basic types within TypedEntityId. Did you use "
                "Component when you should have used ComponentType?");

  explicit TypedEntityId()
      : AggregateType<EntityBase, ComponentTypes...>(kInvalidEntityId,
                                                     ComponentTypes{}...) {}

  template <typename... InputComponentTypes>
  TypedEntityId(  // NOLINT
      const TypedEntityId<InputComponentTypes...>& other)
      : AggregateType<EntityBase, ComponentTypes...>(other) {}
  explicit TypedEntityId(EntityId _id)
      : AggregateType<EntityBase, ComponentTypes...>(_id, ComponentTypes{}...) {
  }

  explicit TypedEntityId(EntityId::ValueType _id)
      : AggregateType<EntityBase, ComponentTypes...>(_id, ComponentTypes{}...) {
  }

  // Returns the raw value of the entity id held by the typed id.
  EntityId::ValueType value() const { return EntityBase::id.value(); }

  // The typed id is implicitly convertible to the EntityId type so that it can
  // be used to to pass into methods that do not require a typed id but instead
  // only care about the raw id itself.
  operator EntityId() const { return EntityBase::id; }  // NOLINT

  template <typename H>
  friend H AbslHashValue(H h, const TypedEntityId<ComponentTypes...>& e) {
    return H::combine(std::move(h), e.id);
  }

  using AggregateType<EntityBase, ComponentTypes...>::operator=;

  bool operator<(const TypedEntityId<ComponentTypes...>& other) const {
    return EntityBase::id < other.id;
  }
};

template <typename... ComponentTypes>
inline std::ostream& operator<<(std::ostream& os,
                                const TypedEntityId<ComponentTypes...>& e) {
  os << e.id;
  return os;
}

// Predefined common EntityId types, one for each component.
using EntityHandle = TypedEntityId<>;
using AttachmentEntityId = TypedEntityId<AttachmentComponentType>;
using CollectionsEntityId = TypedEntityId<CollectionsComponentType>;
using CollectionsMemberEntityId = TypedEntityId<CollectionsMemberComponentType>;
using CollisionEntityId = TypedEntityId<CollisionComponentType>;
using EquipmentEntityId = TypedEntityId<EquipmentComponentType>;
using GeometryEntityId = TypedEntityId<GeometryComponentType>;
using KinematicsEntityId = TypedEntityId<KinematicsComponentType>;
using GripperEntityId = TypedEntityId<GripperComponentType>;
using OutfeedEntityId = TypedEntityId<OutfeedComponentType>;
using PhysicsEntityId = TypedEntityId<PhysicsComponentType>;
using RegionsEntityId = TypedEntityId<RegionsComponentType>;
using RobotEntityId = TypedEntityId<RobotComponentType>;
using UserDataEntityId = TypedEntityId<UserDataComponentType>;
using SensorEntityId = TypedEntityId<SensorComponentType>;
using PPREntityId = TypedEntityId<PPRComponentType>;
using ProjectorEntityId = TypedEntityId<ProjectorComponentType>;
using SimulationEntityId = TypedEntityId<SimulationComponentType>;
using SpawnerEntityId = TypedEntityId<SpawnerComponentType>;

// Typed entity id representing a logical Joint Entity.
using JointEntityId =
    TypedEntityId<AttachmentComponentType, KinematicsComponentType,
                  CollectionsMemberComponentType>;

// Typed entity id representing a logical Link Entity.
using LinkEntityId =
    TypedEntityId<AttachmentComponentType, CollisionComponentType,
                  PhysicsComponentType, GeometryComponentType,
                  CollectionsMemberComponentType>;

// Typed entity id representing a robot collections Entity.
using RobotCollectionsEntityId =
    TypedEntityId<CollectionsComponentType, RobotComponentType>;

// Typed entity id representing a robot-associated coordinate frame.
using RobotCoordinateFrameEntityId =
    TypedEntityId<AttachmentComponentType, CollectionsMemberComponentType>;

// Typed entity id representing a gripper collections Entity.
using GripperCollectionsEntityId =
    TypedEntityId<CollectionsComponentType, GripperComponentType>;

// Typed entity id representing a logical Object Entity.
using PhysicalEntityId =
    TypedEntityId<AttachmentComponentType, CollisionComponentType,
                  GeometryComponentType>;

// The World has a single root coordinate frame entity that is always there and
// never changes.
extern const AttachmentEntityId kRootEntityId;

}  // namespace intrinsic
#endif  // INTRINSIC_WORLD_ENTITY_ID_H_
