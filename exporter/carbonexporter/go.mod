module github.com/open-telemetry/opentelemetry-collector-contrib/exporter/carbonexporter

go 1.16

require (
	github.com/armon/go-metrics v0.3.3 // indirect
	github.com/census-instrumentation/opencensus-proto v0.3.0
	github.com/hashicorp/go-immutable-radix v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/collector v0.28.1-0.20210616151306-cdc163427b8e
	go.uber.org/zap v1.17.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/ini.v1 v1.57.0 // indirect
)

replace go.opentelemetry.io/collector => /Users/adgollap/Documents/GitHub/opentelemetry-collector-contrib/../opentelemetry-collector
