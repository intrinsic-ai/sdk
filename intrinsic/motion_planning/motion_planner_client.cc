// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/motion_planning/motion_planner_client.h"

#include <cstdint>
#include <memory>
#include <optional>
#include <string>
#include <utility>
#include <vector>

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "google/protobuf/duration.pb.h"
#include "google/protobuf/empty.pb.h"
#include "grpcpp/client_context.h"
#include "intrinsic/eigenmath/types.h"
#include "intrinsic/icon/proto/joint_space.pb.h"
#include "intrinsic/logging/proto/context.pb.h"
#include "intrinsic/math/pose3.h"
#include "intrinsic/math/proto_conversion.h"
#include "intrinsic/motion_planning/conversions.h"
#include "intrinsic/motion_planning/proto/motion_target.pb.h"
#include "intrinsic/motion_planning/proto/v1/motion_planner_config.pb.h"
#include "intrinsic/motion_planning/proto/v1/motion_planner_service.grpc.pb.h"
#include "intrinsic/motion_planning/proto/v1/motion_planner_service.pb.h"
#include "intrinsic/util/eigen.h"
#include "intrinsic/util/status/status_conversion_grpc.h"
#include "intrinsic/util/status/status_macros.h"
#include "intrinsic/world/objects/kinematic_object.h"
#include "intrinsic/world/objects/transform_node.h"
#include "intrinsic/world/proto/object_world_refs.pb.h"

namespace intrinsic {
namespace motion_planning {

const MotionPlannerClient::MotionPlanningOptions&
MotionPlannerClient::MotionPlanningOptions::Defaults() {
  static const auto* defaults = new MotionPlannerClient::MotionPlanningOptions({
      .path_planning_time_out = 30,
      .path_planning_step_size = std::nullopt,
      .compute_swept_volume = false,
      .lock_motion_configuration = std::nullopt,
      .skip_fuzzy_cache_check = false,
  });

  return *defaults;
}

MotionPlannerClient::MotionPlannerClient(
    absl::string_view world_id,
    std::shared_ptr<intrinsic_proto::motion_planning::v1::MotionPlannerService::
                        StubInterface>
        motion_planner_service)
    : world_id_(world_id),
      motion_planner_service_(std::move(motion_planner_service)) {}

absl::StatusOr<MotionPlannerClient::PlanTrajectoryResult>
MotionPlannerClient::PlanTrajectory(
    const intrinsic_proto::motion_planning::v1::RobotSpecification&
        robot_specification,
    const intrinsic_proto::motion_planning::v1::MotionSpecification&
        motion_specification,
    const MotionPlanningOptions& options, const std::string& caller_id,
    const intrinsic_proto::data_logger::Context& context) {
  intrinsic_proto::motion_planning::v1::MotionPlanningRequest request;
  *request.mutable_robot_specification() = robot_specification;
  *request.mutable_motion_specification() = motion_specification;
  request.set_world_id(world_id_);
  request.set_compute_swept_volume(options.compute_swept_volume);
  request.mutable_motion_planner_config()->mutable_timeout_sec()->set_seconds(
      options.path_planning_time_out);
  const int64_t s = request.motion_planner_config().timeout_sec().seconds();
  request.mutable_motion_planner_config()->mutable_timeout_sec()->set_nanos(
      (options.path_planning_time_out - s) * 1e9);
  if (options.path_planning_step_size.has_value()) {
    request.mutable_motion_planner_config()->set_path_planning_step_size(
        *options.path_planning_step_size);
  }
  if (options.lock_motion_configuration.has_value()) {
    *request.mutable_motion_planner_config()
         ->mutable_lock_motion_configuration() =
        options.lock_motion_configuration.value();
  }
  request.set_caller_id(caller_id);
  *request.mutable_context() = context;

  intrinsic_proto::motion_planning::v1::TrajectoryPlanningResponse response;
  grpc::ClientContext ctx;
  INTR_RETURN_IF_ERROR(ToAbslStatus(
      motion_planner_service_->PlanTrajectory(&ctx, request, &response)));

  MotionPlannerClient::PlanTrajectoryResult result;
  result.trajectory = response.discretized();
  result.swept_volume.insert(result.swept_volume.begin(),
                             response.swept_volume().begin(),
                             response.swept_volume().end());
  result.lock_motion_id = response.has_lock_motion_id()
                              ? std::optional(response.lock_motion_id())
                              : std::nullopt;
  result.logging_id = response.logging_id();

  return result;
}

absl::StatusOr<std::vector<eigenmath::VectorXd>> MotionPlannerClient::ComputeIk(
    const world::KinematicObject& robot,
    const intrinsic_proto::motion_planning::CartesianMotionTarget&
        cartesian_target,
    const IkOptions& options) {
  // Convert CartesianMotionTarget to GeometricConstraint::PoseEquality.
  intrinsic_proto::motion_planning::v1::GeometricConstraint geometric_target;
  *geometric_target.mutable_cartesian_pose()->mutable_target_frame() =
      cartesian_target.frame();
  *geometric_target.mutable_cartesian_pose()->mutable_moving_frame() =
      cartesian_target.tool();
  if (cartesian_target.has_offset()) {
    *geometric_target.mutable_cartesian_pose()->mutable_target_frame_offset() =
        cartesian_target.offset();
  }

  return ComputeIk(robot, geometric_target, options);
}

absl::StatusOr<std::vector<eigenmath::VectorXd>> MotionPlannerClient::ComputeIk(
    const world::KinematicObject& robot,
    const intrinsic_proto::motion_planning::v1::GeometricConstraint&
        geometric_target,
    const IkOptions& options) {
  intrinsic_proto::motion_planning::v1::IkRequest request;
  request.set_world_id(world_id_);
  request.mutable_robot_reference()->mutable_object_id()->set_id(
      robot.Id().value());

  *request.mutable_target() = geometric_target;

  if (options.starting_joints.size() > 0) {
    VectorXdToRepeatedDouble(
        options.starting_joints,
        request.mutable_starting_joints()->mutable_joints());
  }

  if (options.max_num_solutions.has_value()) {
    request.set_max_num_solutions(*options.max_num_solutions);
  }

  if (options.collision_settings.has_value()) {
    *request.mutable_collision_settings() = *options.collision_settings;
  }

  request.set_ensure_same_branch(options.ensure_same_branch);

  request.set_prefer_same_branch(options.prefer_same_branch);

  intrinsic_proto::motion_planning::v1::IkResponse response;
  grpc::ClientContext ctx;
  INTR_RETURN_IF_ERROR(ToAbslStatus(
      motion_planner_service_->ComputeIk(&ctx, request, &response)));

  return ToVectorXds(response.solutions());
}

namespace {

// Common adaptor to handle different specifications of reference, target.
absl::StatusOr<Pose3d> ComputeFkInternal(
    const world::KinematicObject& robot,
    const eigenmath::VectorXd& joint_values,
    const intrinsic_proto::world::TransformNodeReference& reference,
    const intrinsic_proto::world::TransformNodeReference& target,
    const std::string& world_id,
    intrinsic_proto::motion_planning::v1::MotionPlannerService::StubInterface&
        motion_planner_service) {
  intrinsic_proto::motion_planning::v1::FkRequest request;
  request.set_world_id(world_id);
  request.mutable_robot_reference()->mutable_object_id()->set_id(
      robot.Id().value());
  VectorXdToRepeatedDouble(joint_values,
                           request.mutable_joints()->mutable_joints());

  *request.mutable_reference() = reference;
  *request.mutable_target() = target;

  intrinsic_proto::motion_planning::v1::FkResponse response;
  grpc::ClientContext ctx;
  INTR_RETURN_IF_ERROR(
      ToAbslStatus(motion_planner_service.ComputeFk(&ctx, request, &response)));

  return FromProto(response.reference_t_target());
}

}  // namespace

absl::StatusOr<Pose3d> MotionPlannerClient::ComputeFk(
    const world::KinematicObject& robot,
    const eigenmath::VectorXd& joint_values,
    const intrinsic_proto::world::TransformNodeReferenceByName& reference,
    const intrinsic_proto::world::TransformNodeReferenceByName& target) {
  intrinsic_proto::world::TransformNodeReference reference_proto;
  *reference_proto.mutable_by_name() = reference;
  intrinsic_proto::world::TransformNodeReference target_proto;
  *target_proto.mutable_by_name() = target;
  return ComputeFkInternal(robot, joint_values, reference_proto, target_proto,
                           world_id_, *motion_planner_service_);
}

absl::StatusOr<Pose3d> MotionPlannerClient::ComputeFk(
    const world::KinematicObject& robot,
    const eigenmath::VectorXd& joint_values,
    const world::TransformNode& reference, const world::TransformNode& target) {
  intrinsic_proto::world::TransformNodeReference reference_proto;
  reference_proto.set_id(reference.Id().value());
  intrinsic_proto::world::TransformNodeReference target_proto;
  target_proto.set_id(target.Id().value());
  return ComputeFkInternal(robot, joint_values, reference_proto, target_proto,
                           world_id_, *motion_planner_service_);
}

absl::StatusOr<intrinsic_proto::motion_planning::v1::CheckCollisionsResponse>
MotionPlannerClient::CheckCollisions(
    const world::KinematicObject& robot,
    const std::vector<eigenmath::VectorXd>& waypoints,
    const CheckCollisionsOptions& options) {
  intrinsic_proto::motion_planning::v1::CheckCollisionsRequest request;
  request.set_world_id(world_id_);
  request.mutable_robot_reference()->mutable_object_id()->set_id(
      robot.Id().value());
  ToJointVecs(waypoints, request.mutable_waypoint());

  if (options.collision_settings.has_value()) {
    *request.mutable_collision_settings() = *options.collision_settings;
  }

  intrinsic_proto::motion_planning::v1::CheckCollisionsResponse response;
  grpc::ClientContext ctx;
  INTR_RETURN_IF_ERROR(ToAbslStatus(
      motion_planner_service_->CheckCollisions(&ctx, request, &response)));
  return response;
}

absl::StatusOr<google::protobuf::Empty> MotionPlannerClient::ClearCache() {
  google::protobuf::Empty response;
  grpc::ClientContext ctx;
  INTR_RETURN_IF_ERROR(ToAbslStatus(motion_planner_service_->ClearCache(
      &ctx, google::protobuf::Empty(), &response)));
  return response;
}

}  // namespace motion_planning
}  // namespace intrinsic
