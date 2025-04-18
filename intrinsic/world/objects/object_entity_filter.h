// Copyright 2023 Intrinsic Innovation LLC

#ifndef INTRINSIC_WORLD_OBJECTS_OBJECT_ENTITY_FILTER_H_
#define INTRINSIC_WORLD_OBJECTS_OBJECT_ENTITY_FILTER_H_

#include <set>
#include <string>
#include <tuple>

#include "absl/strings/string_view.h"
#include "absl/types/span.h"
#include "intrinsic/world/objects/object_world_ids.h"
#include "intrinsic/world/proto/object_world_refs.pb.h"

namespace intrinsic {
namespace world {

// Selects one or more entities of a WorldObject, which are the parts of an
// object.
//
// Not all objects have more than one entity. For single-entity objects all of
// the provided options will behave identically and select the one and only
// existing entity.
class ObjectEntityFilter {
 public:
  // Changes the filter to also include the base entity of an object.
  ObjectEntityFilter& IncludeBaseEntity();

  // Return True if filter includes the base entity of an object.
  bool IncludesBaseEntity() const { return include_base_entity_; }

  // Changes the filter to also include the final entity of an object. Using the
  // resulting filter on an object that does not have a unique final entity will
  // result in an error.
  ObjectEntityFilter& IncludeFinalEntity();

  // Return True if filter includes the final entity of an object.
  bool IncludesFinalEntity() const { return include_final_entity_; }

  // Changes the filter to include all entities of an object.
  ObjectEntityFilter& IncludeAllEntities();

  // Return True if filter includes all entities of an object.
  bool IncludesAllEntities() const { return include_all_entities_; }

  // Changes the filter to include the given entity of an object.
  ObjectEntityFilter& IncludeEntityId(ObjectWorldResourceId entity_id);

  // Removes any previously added explicit entity ids.
  ObjectEntityFilter& ClearExplicitEntityIds();

  // Return Entity IDs referenced by this filter.
  const std::set<ObjectWorldResourceId>& EntityIds() const {
    return entity_ids_;
  }

  // Changes the filter to include the given entity local name.
  ObjectEntityFilter& IncludeEntityName(absl::string_view entity_name);

  // Removes any previously added explicit entity local names.
  ObjectEntityFilter& ClearExplicitEntityNames();

  // Return Entity names referenced by this filter.
  const std::set<std::string>& EntityNames() const { return entity_names_; }

  // Returns the proto representation of this instance.
  intrinsic_proto::world::ObjectEntityFilter ToProto() const;

  // Parses the given proto and returns an instance of ObjectEntityFilter.
  static ObjectEntityFilter FromProto(
      const ::intrinsic_proto::world::ObjectEntityFilter& entity_filter);

  // Returns an ObjectEntityFilter that includes all the named entities.
  static ObjectEntityFilter FromEntityNames(
      absl::Span<const absl::string_view> entity_names);

  // Returns an ObjectEntityFilter that includes all the ID'd entities.
  static ObjectEntityFilter FromEntityIds(
      absl::Span<const ObjectWorldResourceId> entity_ids);

  // Returns a statically constructed ObjectEntityFilter that includes the base
  // entity. This is a helper to minimize creation of temporary objects.
  static const ObjectEntityFilter& BaseEntity();

  // Returns a statically constructed ObjectEntityFilter that includes the final
  // entity. This is a helper to minimize creation of temporary objects.
  static const ObjectEntityFilter& FinalEntity();

  // Returns a statically constructed ObjectEntityFilter that includes all the
  // entities. This is a helper to minimize creation of temporary objects.
  static const ObjectEntityFilter& AllEntities();

  friend bool operator==(const ObjectEntityFilter& lhs,
                         const ObjectEntityFilter& rhs);

 private:
  bool include_base_entity_ = false;
  bool include_final_entity_ = false;
  bool include_all_entities_ = false;
  std::set<ObjectWorldResourceId> entity_ids_;
  std::set<std::string> entity_names_;
};

// Equality operator
inline bool operator==(const ObjectEntityFilter& lhs,
                       const ObjectEntityFilter& rhs) {
  return std::tie(lhs.include_base_entity_, lhs.include_final_entity_,
                  lhs.include_all_entities_, lhs.entity_ids_,
                  lhs.entity_names_) ==
         std::tie(rhs.include_base_entity_, rhs.include_final_entity_,
                  rhs.include_all_entities_, rhs.entity_ids_,
                  rhs.entity_names_);
}

// Inequality operator
inline bool operator!=(const ObjectEntityFilter& lhs,
                       const ObjectEntityFilter& rhs) {
  return !(lhs == rhs);
}

}  // namespace world
}  // namespace intrinsic

#endif  // INTRINSIC_WORLD_OBJECTS_OBJECT_ENTITY_FILTER_H_
