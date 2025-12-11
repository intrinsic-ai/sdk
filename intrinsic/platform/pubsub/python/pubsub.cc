// Copyright 2023 Intrinsic Innovation LLC

#include "intrinsic/platform/pubsub/pubsub.h"

#include <pybind11/functional.h>
#include <pybind11/pybind11.h>
#include <pybind11/pytypes.h>
#include <pybind11/stl.h>

#include <memory>
#include <optional>
#include <string>
#include <string_view>
#include <utility>
#include <vector>

#include "absl/container/flat_hash_map.h"
#include "absl/log/log.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"
#include "absl/time/time.h"
#include "google/protobuf/any.pb.h"
#include "google/protobuf/message.h"
#include "intrinsic/platform/pubsub/kvstore.h"
#include "intrinsic/platform/pubsub/publisher.h"
#include "intrinsic/platform/pubsub/subscription.h"
#include "pybind11/cast.h"
#include "pybind11_abseil/absl_casters.h"
#include "pybind11_abseil/no_throw_status.h"
#include "pybind11_abseil/status_casters.h"
#include "pybind11_protobuf/native_proto_caster.h"

namespace intrinsic {
namespace pubsub {

namespace {

absl::StatusOr<Subscription> CreateSubscriptionWithConfig(
    PubSub* self, absl::string_view topic, const TopicConfig& config,
    const google::protobuf::Message& exemplar, pybind11::object msg_callback,
    pybind11::object err_callback) {
  // The callback passed to the adapter must be able to be copied in a
  // separate thread without copying the msg_callback.
  // This allows the message callback to capture variables which
  // are not possible (or safe) to copy in a separate thread. This is the
  // case when the callback captures a python function, since those cannot
  // be copied without holding the GIL, and the adapter thread executing the
  // callback does not know to acquire the GIL. Using a shared pointer to
  // own the adapter callback satisfies these requirements.

  SubscriptionOkCallback<google::protobuf::Message> message_callback = {};
  SubscriptionErrorCallback error_callback = {};

  if (msg_callback && !msg_callback.is_none()) {
    message_callback = [py_msg_cb = std::move(msg_callback)](
                           const google::protobuf::Message& msg) {
      pybind11::gil_scoped_acquire gil;
      try {
        // This will create a copy in the py proto caster
        py_msg_cb(msg);
      } catch (const pybind11::error_already_set& e) {
        LOG(ERROR) << "Exception in message callback: " << e.what();
      }
    };
  }

  if (err_callback && !err_callback.is_none()) {
    error_callback = [py_err_cb = std::move(err_callback)](
                         absl::string_view packet, absl::Status error) {
      pybind11::gil_scoped_acquire gil;
      try {
        py_err_cb(packet, pybind11::google::DoNotThrowStatus(error));
      } catch (const pybind11::error_already_set& e) {
        LOG(ERROR) << "Exception in error callback: " << e.what();
      }
    };
  }

  return self->CreateSubscription(topic, config, exemplar,
                                  std::move(message_callback),
                                  std::move(error_callback));
}

absl::StatusOr<Subscription> CreateSubscription(
    PubSub* self, absl::string_view topic,
    const google::protobuf::Message& exemplar, pybind11::object msg_callback,
    pybind11::object err_callback) {
  return CreateSubscriptionWithConfig(self, topic, TopicConfig{}, exemplar,
                                      std::move(msg_callback),
                                      std::move(err_callback));
}

absl::StatusOr<Subscription> CreateRawSubscription(
    PubSub* self, absl::string_view topic, const TopicConfig& config,
    pybind11::object msg_callback) {
  // The callback passed to the adapter must be able to be copied in a
  // separate thread without copying the msg_callback.
  // This allows the message callback to capture variables which
  // are not possible (or safe) to copy in a separate thread. This is the
  // case when the callback captures a python function, since those cannot
  // be copied without holding the GIL, and the adapter thread executing the
  // callback does not know to acquire the GIL. Using a shared pointer to
  // own the adapter callback satisfies these requirements.

  SubscriptionOkCallback<intrinsic_proto::pubsub::PubSubPacket>
      message_callback = {};

  if (msg_callback && !msg_callback.is_none()) {
    message_callback = [py_msg_cb = std::move(msg_callback)](
                           const intrinsic_proto::pubsub::PubSubPacket& msg) {
      pybind11::gil_scoped_acquire gil;
      try {
        // This will create a copy in the py proto caster
        py_msg_cb(msg.payload());
      } catch (const pybind11::error_already_set& e) {
        LOG(ERROR) << "Exception in message callback: " << e.what();
      }
    };
  }

  return self->CreateSubscription(topic, config, std::move(message_callback));
}

absl::StatusOr<KeyValueStore> CreateKeyValueStore(
    PubSub* self, std::optional<std::string> prefix_override) {
  return self->KeyValueStore(prefix_override);
}

absl::StatusOr<KeyValueStore> CreateReplicationKVStore(PubSub* self) {
  return self->KeyValueStore(std::string(intrinsic::kReplicationPrefix));
}

absl::StatusOr<KVQuery> GetAll(KeyValueStore* self, const std::string& key,
                               KeyValueCallback callback,
                               OnDoneCallback on_done) {
  return self->GetAll(key, callback, on_done);
}

absl::StatusOr<google::protobuf::Any> Get(KeyValueStore* self,
                                          const std::string& key, int timeout) {
  return self->Get<google::protobuf::Any>(key, absl::Seconds(timeout));
}

absl::StatusOr<std::vector<std::string>> ListAllKeys(KeyValueStore* self,
                                                     int timeout) {
  return self->ListAllKeys(absl::Seconds(timeout));
}

absl::StatusOr<std::vector<std::string>> ListAllGlobalKeys(KeyValueStore* self,
                                                           int timeout) {
  return self->ListAllGlobalKeys(absl::Seconds(timeout));
}

absl::StatusOr<std::vector<std::string>> ListAllOnpremKeys(
    KeyValueStore* self, const std::string& workcell_name, int timeout) {
  return self->ListAllOnpremKeys(workcell_name, absl::Seconds(timeout));
}

absl::Status AdminCloudCopy(KeyValueStore* self, const std::string& source_key,
                            const std::string& target_key, int timeout) {
  return self->AdminCloudCopy(source_key, target_key, absl::Seconds(timeout));
}

absl::StatusOr<absl::flat_hash_map<std::string, google::protobuf::Any>>
GetAllSynchronous(KeyValueStore* self, const std::string& keyexpr,
                  int timeout) {
  return self->GetAllSynchronous(keyexpr, absl::Seconds(timeout));
}

struct PySubscriptionDeleter {
  void operator()(Subscription* s) {
    // To avoid deadlock, the call to Zenoh.imw_destroy_subscription() needs
    // to happen with the GIL released. Otherwise, the GIL and the internal
    // callback mutex are potentially locked in opposite order by this thread
    // and the Zenoh callback thread pool, which can deadlock, especially on
    // high-frequency topics.
    {
      pybind11::gil_scoped_release release_gil;
      s->Unsubscribe();
    }

    // The Python GIL will be re-acquired now that the previous scoped_release
    // has disappeared. With the re-acquired GIL, we can safely delete the
    // subscription_data_ struct in Subscription, which contains the Python
    // callback object. A deadlock can no longer occur, because a message
    // callback will no longer occur because the remainder of the destruction
    // call chain is holding the GIL.
    delete s;
  }
};

}  // namespace

PYBIND11_MODULE(pubsub, m) {
  pybind11::google::ImportStatusModule();
  pybind11_protobuf::ImportNativeProtoCasters();

  pybind11::enum_<TopicConfig::TopicQoS>(m, "TopicQoS")
      .value("HighReliability", TopicConfig::TopicQoS::HighReliability)
      .value("Sensor", TopicConfig::TopicQoS::Sensor)
      .export_values();

  pybind11::class_<TopicConfig>(m, "TopicConfig")
      .def(pybind11::init<>())
      .def_readwrite("topic_qos", &TopicConfig::topic_qos);

  pybind11::class_<PubSub>(m, "PubSub")
      .def(pybind11::init<>())
      .def(pybind11::init<std::string_view>(),
           pybind11::arg("participant_name"))
      .def(pybind11::init<std::string_view, std::string_view>(),
           pybind11::arg("participant_name"), pybind11::arg("config"))
      // Cast required for overloaded methods:
      // https://pybind11.readthedocs.io/en/stable/classes.html#overloaded-methods
      .def("CreatePublisher", &PubSub::CreatePublisher, pybind11::arg("topic"),
           pybind11::arg("config") = TopicConfig{})
      .def("CreateSubscription", &CreateRawSubscription, pybind11::arg("topic"),
           pybind11::arg("config"), pybind11::arg("msg_callback") = nullptr)
      .def("CreateSubscription", &CreateSubscriptionWithConfig,
           pybind11::arg("topic"), pybind11::arg("config"),
           pybind11::arg("exemplar"), pybind11::arg("msg_callback") = nullptr,
           pybind11::arg("error_callback") = nullptr)
      .def("CreateSubscription", &CreateSubscription, pybind11::arg("topic"),
           pybind11::arg("exemplar"), pybind11::arg("msg_callback") = nullptr,
           pybind11::arg("error_callback") = nullptr)
      .def("KeyValueStore", &CreateKeyValueStore,
           pybind11::arg("prefix_override") = std::nullopt)
      .def("ReplicationKeyValueStore", &CreateReplicationKVStore);

  pybind11::class_<Publisher>(m, "Publisher")
      .def("Publish",
           static_cast<absl::Status (Publisher::*)(
               const google::protobuf::Message&) const>(&Publisher::Publish),
           pybind11::arg("message"))
      .def("TopicName", &Publisher::TopicName)
      .def("HasMatchingSubscribers", &Publisher::HasMatchingSubscribers);

  pybind11::class_<KVQuery>(m, "KVQuery");

  pybind11::class_<KeyValueStore>(m, "KeyValueStore")
      .def("Set", &KeyValueStore::Set<google::protobuf::Message>,
           pybind11::arg("key"), pybind11::arg("value"),
           pybind11::arg("high_consistency") = false)
      .def("Get", &Get, pybind11::arg("key"), pybind11::arg("timeout") = 10)
      .def("GetAll", &GetAll)
      .def("GetAllSynchronous", &GetAllSynchronous, pybind11::arg("keyexpr"),
           pybind11::arg("timeout") = 10)
      .def("ListAllKeys", &ListAllKeys, pybind11::arg("timeout") = 10)
      .def("ListAllGlobalKeys", &ListAllGlobalKeys,
           pybind11::arg("timeout") = 10)
      .def("ListAllOnpremKeys", &ListAllOnpremKeys,
           pybind11::arg("workcell_name"), pybind11::arg("timeout") = 10)
      .def("Delete", &KeyValueStore::Delete, pybind11::arg("key"))
      .def("AdminCloudCopy", &AdminCloudCopy, pybind11::arg("source_key"),
           pybind11::arg("target_key"), pybind11::arg("timeout") = 10);

  // The python GIL does not need to be locked during the entire destructor
  // of this class. Instead, the custom deleter provided during its
  // construction will acquire the GIL only during the deletion of the
  // SubscriptionData object, which holds the Python callback.
  pybind11::class_<Subscription,
                   std::unique_ptr<Subscription, PySubscriptionDeleter>>(
      m, "Subscription")
      .def("TopicName", &Subscription::TopicName);
}

}  // namespace pubsub
}  // namespace intrinsic
