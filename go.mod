module api-gateway/v2

go 1.21.0

toolchain go1.21.6

require (
	contrib.go.opencensus.io/exporter/aws v0.0.0-20230502192102-15967c811cec
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/prometheus v0.0.0-20190424224027-f02a6e68f94d
	contrib.go.opencensus.io/exporter/stackdriver v0.13.14
	contrib.go.opencensus.io/exporter/zipkin v0.0.0-20190424224031-c96617f51dc6
	github.com/DataDog/opencensus-go-exporter-datadog v0.0.0-20191210083620-6965a1cfed68
	github.com/Masterminds/sprig/v3 v3.2.3
	github.com/alecthomas/chroma v0.6.3
	github.com/aws/aws-sdk-go v1.47.13
	github.com/catalinc/hashcash v0.0.0-20161205220751-e6bc29ff4de9
	github.com/clbanning/mxj v1.8.4
	github.com/davron112/lura/v2 v2.8.0
	github.com/elliotchance/phpserialize v1.4.0
	github.com/fatih/color v1.16.0
	github.com/gin-gonic/gin v1.9.1
	github.com/go-contrib/uuid v1.2.0
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/golang/protobuf v1.5.3
	github.com/google/cel-go v0.11.4
	github.com/google/martian v2.1.1-0.20190517191504-25dcb96d9e51+incompatible
	github.com/google/uuid v1.6.0
	github.com/hashicorp/golang-lru v1.0.2
	github.com/influxdata/influxdb v1.9.7
	github.com/jhump/protoreflect v1.15.6
	github.com/juju/ratelimit v1.0.1
	github.com/kpacha/opencensus-influxdb v0.0.0-20180520162117-1b490a38de4c
	github.com/krakend/go-auth0 v1.0.0
	github.com/mattn/go-isatty v0.0.20
	github.com/mmcdole/gofeed v1.2.1
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/openzipkin/zipkin-go v0.2.5
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.16.0
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475
	github.com/rs/cors v1.8.2
	github.com/rs/cors/wrapper/gin v0.0.0-20221003140808-fcebdb403f4d
	github.com/santhosh-tekuri/jsonschema/v5 v5.0.0
	github.com/sony/gobreaker v0.5.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/streadway/amqp v1.0.0
	github.com/tmthrgd/go-bitset v0.0.0-20190904054048-394d9a556c05
	github.com/unrolled/secure v1.13.0
	github.com/valyala/fastrand v1.1.0
	github.com/xeipuuv/gojsonschema v1.2.1-0.20200424115421-065759f9c3d7
	github.com/yuin/gopher-lua v0.0.0-20190206043414-8bfc7677f583
	go.opencensus.io v0.24.0
	gocloud.dev v0.34.0
	gocloud.dev/pubsub/kafkapubsub v0.25.0
	gocloud.dev/pubsub/natspubsub v0.25.0
	gocloud.dev/pubsub/rabbitpubsub v0.25.0
	gocloud.dev/secrets/hashivault v0.34.0
	golang.org/x/mod v0.17.0
	golang.org/x/net v0.27.0
	golang.org/x/oauth2 v0.16.0
	golang.org/x/sync v0.7.0
	golang.org/x/text v0.16.0
	google.golang.org/grpc v1.62.1
	gopkg.in/Graylog2/go-gelf.v2 v2.0.0-20191017102106-1550ee647df0
	gopkg.in/square/go-jose.v2 v2.6.0
)

require (
	cloud.google.com/go/compute v1.23.3 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v1.1.5 // indirect
	cloud.google.com/go/kms v1.15.5 // indirect
	cloud.google.com/go/monitoring v1.17.0 // indirect
	cloud.google.com/go/pubsub v1.34.0 // indirect
	cloud.google.com/go/trace v1.10.4 // indirect
	github.com/Azure/azure-amqp-common-go/v3 v3.2.3 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.9.2 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.5.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.5.2 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys v0.10.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal v0.7.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus v1.5.0 // indirect
	github.com/Azure/go-amqp v1.0.4 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.2.2 // indirect
	github.com/DataDog/datadog-go v3.4.1+incompatible // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/PuerkitoBio/goquery v1.8.1 // indirect
	github.com/Shopify/sarama v1.34.1 // indirect
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/antlr/antlr4/runtime/Go/antlr v0.0.0-20220418222510-f25a4f6275ed // indirect
	github.com/aws/aws-sdk-go-v2 v1.23.0 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.25.2 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.16.1 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.14.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.7.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.10.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.10.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/kms v1.26.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sns v1.21.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sqs v1.24.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.17.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.19.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.25.2 // indirect
	github.com/aws/smithy-go v1.17.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bufbuild/protocompile v0.8.0 // indirect
	github.com/bytedance/sonic v1.11.6 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cenkalti/backoff/v3 v3.2.2 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/danwakefield/fnmatch v0.0.0-20160403171240-cbb64ac3d964 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dlclark/regexp2 v1.1.6 // indirect
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-jose/go-jose/v3 v3.0.1 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/wire v0.5.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.16.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.5 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.8 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.6 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/hcl v1.0.1-vault-5 // indirect
	github.com/hashicorp/vault/api v1.10.0 // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.2 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/klauspost/cpuid/v2 v2.2.6 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/mmcdole/goxpp v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nats-io/nats-server/v2 v2.9.23 // indirect
	github.com/nats-io/nats.go v1.28.0 // indirect
	github.com/nats-io/nkeys v0.4.6 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pelletier/go-toml v1.9.3 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.14 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.0 // indirect
	github.com/prometheus/prometheus v0.46.0 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/spf13/afero v1.9.2 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/jwalterweatherman v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/tinylib/msgp v1.1.2 // indirect
	github.com/tmthrgd/atomics v0.0.0-20190904060638-dc7a5fcc7e0d // indirect
	github.com/tmthrgd/go-bitwise v0.0.0-20190904053232-1430ee983fca // indirect
	github.com/tmthrgd/go-byte-test v0.0.0-20190904060354-2794345b9929 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/tmthrgd/go-memset v0.0.0-20190904060434-6fb7a21f88f1 // indirect
	github.com/tmthrgd/go-popcount v0.0.0-20190904054823-afb1ace8b04f // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/uber/jaeger-client-go v2.28.0+incompatible // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.46.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.46.1 // indirect
	go.opentelemetry.io/otel v1.21.0 // indirect
	go.opentelemetry.io/otel/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.21.0 // indirect
	golang.org/x/arch v0.6.0 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	google.golang.org/api v0.155.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20240123012728-ef4313101c80 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240123012728-ef4313101c80 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240123012728-ef4313101c80 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
	gopkg.in/DataDog/dd-trace-go.v1 v1.22.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
