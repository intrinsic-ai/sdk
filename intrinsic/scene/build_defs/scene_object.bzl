# Copyright 2023 Intrinsic Innovation LLC

"""Custom rules for creating SceneObject from SDF."""

SceneObjectInfo = provider(
    "provided by the scene_object() rule",
    fields = {
        "fdsets": "The transitive descriptor sets needed to parse the user provided data in pbtxt",
        "gzf": "The gzf file, containing the extracted scene object and associated geometry",
        "pbtxt": "A pbtxt file representing a intrinsic_proto::scene_object::v1::SceneObject message",
    },
)
