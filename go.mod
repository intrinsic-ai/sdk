// Copyright 2023 Intrinsic Innovation LLC

module intrinsic

go 1.23.0

toolchain go1.23.1

require (
	cloud.google.com/go/bigquery v1.64.0
	cloud.google.com/go/compute/metadata v0.6.0
	cloud.google.com/go/firestore v1.17.0
	cloud.google.com/go/logging v1.12.0
	cloud.google.com/go/longrunning v0.6.2
	cloud.google.com/go/secretmanager v1.14.2
	cloud.google.com/go/spanner v1.73.0
	cloud.google.com/go/storage v1.43.0
	contrib.go.opencensus.io/exporter/ocagent v0.7.0
	contrib.go.opencensus.io/exporter/prometheus v0.4.2
	contrib.go.opencensus.io/exporter/stackdriver v0.13.14
	dario.cat/mergo v1.0.0
	firebase.google.com/go v3.13.0+incompatible
	github.com/authzed/authzed-go v1.3.0
	github.com/bazelbuild/buildtools v0.0.0-20240918101019-be1c24cc9a44
	github.com/bazelbuild/remote-apis-sdks v0.0.0-20230919142202-aa1c266ae342
	github.com/bazelbuild/rules_go v0.51.0
	github.com/bits-and-blooms/bitset v1.20.0
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/containerd/containerd v1.7.20
	github.com/dustin/go-humanize v1.0.1
	github.com/fsnotify/fsnotify v1.6.0
	github.com/fsouza/fake-gcs-server v1.47.4
	github.com/gobwas/glob v0.2.3
	github.com/golang-jwt/jwt/v4 v4.4.2
	github.com/golang/glog v1.2.3
	github.com/golang/protobuf v1.5.4
	github.com/google/brotli/go/cbrotli v1.1.0
	github.com/google/go-cmp v0.6.0
	github.com/google/go-containerregistry v0.20.3
	github.com/google/safehtml v0.1.0
	github.com/google/safetext v0.0.0-20221026122733-23539d61753f
	github.com/google/subcommands v1.2.0
	github.com/google/uuid v1.6.0
	github.com/googlecloudrobotics/core/src v0.0.0-20230426093931-c9725477ada9
	github.com/gorilla/websocket v1.5.0
	github.com/jedib0t/go-pretty/v6 v6.6.7
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kylelemons/godebug v1.1.0
	github.com/lestrrat-go/libxml2 v0.0.0-20240905100032-c934e3fcb9d3
	github.com/minio/highwayhash v1.0.2
	github.com/opencontainers/go-digest v1.0.0
	github.com/opencontainers/image-spec v1.1.0
	github.com/pborman/uuid v1.2.1
	github.com/pion/ice/v4 v4.0.10
	github.com/pion/rtcp v1.2.15
	github.com/pion/transport/v3 v3.0.7
	github.com/pion/turn/v4 v4.0.2
	github.com/pion/webrtc/v4 v4.1.1
	github.com/pkg/errors v0.9.1
	github.com/protocolbuffers/txtpbfmt v0.0.0-20230830074515-d9d085e6be90
	github.com/robfig/cron/v3 v3.0.1
	github.com/rs/xid v1.5.0
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/spf13/cobra v1.8.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.16.0
	github.com/stoewer/go-strcase v1.3.0
	github.com/stripe/stripe-go/v79 v79.12.0
	github.com/tdewolff/parse v2.3.4+incompatible
	github.com/traviswt/gke-auth-plugin v0.0.0-20230623230742-7b40450b6d49
	go.etcd.io/bbolt v1.3.10
	go.opencensus.io v0.24.0
	go.uber.org/atomic v1.11.0
	golang.org/x/crypto v0.33.0
	golang.org/x/exp v0.0.0-20240909161429-701f63a606c0
	golang.org/x/net v0.35.0
	golang.org/x/oauth2 v0.25.0
	golang.org/x/sync v0.11.0
	golang.org/x/sys v0.30.0
	golang.org/x/time v0.8.0
	gonum.org/v1/gonum v0.14.0
	google.golang.org/api v0.214.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.4
	gopkg.in/yaml.v3 v3.0.1
	k8s.io/api v0.30.3
	k8s.io/apimachinery v0.30.3
	k8s.io/client-go v0.30.3
	k8s.io/kubelet v0.30.3
	k8s.io/metrics v0.33.1
	sigs.k8s.io/yaml v1.3.0
)

require (
	github.com/pion/datachannel v1.5.10 // indirect
	github.com/pion/dtls/v3 v3.0.6 // indirect
	github.com/pion/interceptor v0.1.37 // indirect
	github.com/pion/logging v0.2.3 // indirect
	github.com/pion/mdns/v2 v2.0.7 // indirect
	github.com/pion/randutil v0.1.0 // indirect
	github.com/pion/rtp v1.8.15 // indirect
	github.com/pion/sctp v1.8.39 // indirect
	github.com/pion/sdp/v3 v3.0.11 // indirect
	github.com/pion/srtp/v3 v3.0.4 // indirect
	github.com/pion/stun/v3 v3.0.0 // indirect
	github.com/wlynxg/anet v0.0.5 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	cel.dev/expr v0.19.0 // indirect
	cloud.google.com/go v0.116.0 // indirect
	cloud.google.com/go/auth v0.13.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.6 // indirect
	cloud.google.com/go/iam v1.2.2 // indirect
	cloud.google.com/go/monitoring v1.21.2 // indirect
	cloud.google.com/go/pubsub v1.45.1
	cloud.google.com/go/trace v1.11.2 // indirect
	github.com/AdaLogics/go-fuzz-headers v0.0.0-20230811130428-ced1acdcaa24 // indirect
	github.com/AdamKorcz/go-118-fuzz-build v0.0.0-20230306123547-8075edf89bb0 // indirect
	github.com/GoogleCloudPlatform/grpc-gcp-go/grpcgcp v1.5.2 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.25.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/Microsoft/hcsshim v0.11.7 // indirect
	github.com/apache/arrow/go/v15 v15.0.2 // indirect
	github.com/aws/aws-sdk-go v1.44.129 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cncf/xds/go v0.0.0-20240905190251-b4127c9b8d78 // indirect
	github.com/containerd/cgroups v1.1.0 // indirect
	github.com/containerd/containerd/api v1.7.19 // indirect
	github.com/containerd/continuity v0.4.2 // indirect
	github.com/containerd/errdefs v0.1.0 // indirect
	github.com/containerd/fifo v1.1.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/containerd/platforms v0.2.1 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.16.3 // indirect
	github.com/containerd/ttrpc v1.2.5 // indirect
	github.com/containerd/typeurl/v2 v2.1.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.6 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/cli v27.5.0+incompatible // indirect
	github.com/docker/distribution v2.8.3+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.8.2 // indirect
	github.com/docker/go-events v0.0.0-20190806004212-e31b211e4f1c // indirect
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/envoyproxy/go-control-plane v0.13.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/flatbuffers v23.5.26+incompatible // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/renameio/v2 v2.0.0 // indirect
	github.com/google/s2a-go v0.1.8 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.4 // indirect
	github.com/googleapis/gax-go/v2 v2.14.0 // indirect
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.0
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/jzelinskie/stringz v0.0.3 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moby/locker v1.0.1 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/moby/sys/sequential v0.5.0 // indirect
	github.com/moby/sys/signal v0.7.0 // indirect
	github.com/moby/sys/user v0.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/opencontainers/runtime-spec v1.1.0 // indirect
	github.com/opencontainers/selinux v1.11.0 // indirect
	github.com/paul-mannino/go-fuzzywuzzy v0.0.0-20241117160931-a1769aeb6b21
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pierrec/lz4/v4 v4.1.18 // indirect
	github.com/pkg/xattr v0.4.9 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/prometheus/client_golang v1.16.0 // indirect
	github.com/prometheus/client_model v0.6.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.10.1 // indirect
	github.com/prometheus/prometheus v0.39.1 // indirect
	github.com/prometheus/statsd_exporter v0.22.8 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/samber/lo v1.49.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/tdewolff/test v1.0.9 // indirect
	github.com/tklauser/go-sysconf v0.3.14 // indirect
	github.com/tklauser/numcpus v0.8.0 // indirect
	github.com/vbatts/tar-split v0.11.6 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/detectors/gcp v1.32.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.54.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.58.0 // indirect
	go.opentelemetry.io/otel v1.33.0 // indirect
	go.opentelemetry.io/otel/metric v1.33.0 // indirect
	go.opentelemetry.io/otel/sdk v1.33.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.32.0 // indirect
	go.opentelemetry.io/otel/trace v1.33.0 // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/term v0.29.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	golang.org/x/tools v0.29.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20241118233622-e639e219e697 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250115164207-1a7da9e5054f
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.5.0
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	k8s.io/klog/v2 v2.120.1 // indirect
	k8s.io/kube-openapi v0.0.0-20240228011516-70dd3763d340 // indirect
	k8s.io/utils v0.0.0-20230726121419-3b25d923346b // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
)
