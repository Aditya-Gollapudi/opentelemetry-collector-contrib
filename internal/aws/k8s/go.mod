module github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/k8s

go 1.15

require (
	github.com/aws/aws-sdk-go v1.38.52
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/stretchr/testify v1.7.0
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/zap v1.16.0
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.20.4
)

replace go.opentelemetry.io/collector => /Users/adgollap/Documents/GitHub/opentelemetry-collector-contrib/../opentelemetry-collector
