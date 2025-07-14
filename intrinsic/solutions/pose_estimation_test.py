# Copyright 2023 Intrinsic Innovation LLC

from unittest import mock

from absl.testing import absltest
from intrinsic.assets.data.proto.v1 import data_asset_pb2
from intrinsic.assets.proto import asset_type_pb2
from intrinsic.assets.proto import id_pb2
from intrinsic.assets.proto import installed_assets_pb2
from intrinsic.assets.proto import metadata_pb2
from intrinsic.perception.proto.v1 import perception_model_pb2
from intrinsic.resources.proto import resource_registry_pb2
from intrinsic.solutions import pose_estimation


class PoseEstimatorsTest(absltest.TestCase):

  def test_lists_pose_estimators(self):
    installed_assets_stub = mock.MagicMock()
    resource_registry_client = mock.MagicMock()
    resource_registry_client.list_all_resource_instances.return_value = [
        resource_registry_pb2.ResourceInstance(name="pose_estimator_1"),
        resource_registry_pb2.ResourceInstance(name="pose_estimator_2"),
    ]
    data_asset = data_asset_pb2.DataAsset()
    data_asset.data.Pack(perception_model_pb2.PerceptionModel())
    installed_assets_stub.ListInstalledAssets.return_value = installed_assets_pb2.ListInstalledAssetsResponse(
        installed_assets=[
            installed_assets_pb2.InstalledAsset(
                metadata=metadata_pb2.Metadata(
                    id_version=id_pb2.IdVersion(
                        id=id_pb2.Id(
                            package="ai.intrinsic",
                            name="pose_estimator_data_asset",
                        ),
                        version="0.0.1",
                    ),
                    display_name="Pose estimator data asset",
                    asset_type=asset_type_pb2.AssetType.ASSET_TYPE_DATA,
                ),
                deployment_data=installed_assets_pb2.InstalledAsset.DeploymentData(
                    data=installed_assets_pb2.InstalledAsset.DataDeploymentData(
                        data=data_asset,
                    )
                ),
                data_specific_metadata=installed_assets_pb2.InstalledAsset.DataMetadata(
                    proto_name=perception_model_pb2.PerceptionModel.DESCRIPTOR.full_name,
                ),
            )
        ]
    )
    pose_estimators = pose_estimation.PoseEstimators(
        resource_registry_client,
        installed_assets_stub,
    )
    self.assertLen(pose_estimators, 3)
    self.assertEqual(
        dir(pose_estimators),
        ["pose_estimator_1", "pose_estimator_2", "pose_estimator_data_asset"],
    )
    self.assertEqual(pose_estimators.pose_estimator_1.id, "pose_estimator_1")
    self.assertEqual(pose_estimators.pose_estimator_2.id, "pose_estimator_2")
    self.assertEqual(
        pose_estimators.pose_estimator_data_asset.id,
        "pose_estimator_data_asset",
    )
    resource_registry_client.list_all_resource_instances.assert_called_with(
        resource_family_id=pose_estimation._POSE_ESTIMATOR_RESOURCE_FAMILY_ID
    )


if __name__ == "__main__":
  absltest.main()
