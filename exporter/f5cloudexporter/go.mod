module github.com/open-telemetry/opentelemetry-collector-contrib/exporter/f5cloudexporter

go 1.16

require (
	github.com/armon/go-metrics v0.3.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/collector v0.28.1-0.20210616151306-cdc163427b8e
	go.uber.org/zap v1.17.0
	golang.org/x/oauth2 v0.0.0-20210514164344-f6687ab2804c
	google.golang.org/api v0.48.0
	gopkg.in/ini.v1 v1.57.0 // indirect
)

replace go.opentelemetry.io/collector => /Users/adgollap/Documents/GitHub/opentelemetry-collector-contrib/../opentelemetry-collector
