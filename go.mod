module github.com/ONSdigital/dp-s3/v2

go 1.20

retract (
	v2.1.0-beta.2 // Contains retractions only
	[v2.1.0-beta, v2.1.0-beta.1] // Code never made into main and refers to non-released dependencies
)

require (
	github.com/ONSdigital/dp-healthcheck v1.6.1
	github.com/ONSdigital/log.go/v2 v2.4.3
	github.com/aws/aws-sdk-go v1.48.9
	github.com/smartystreets/goconvey v1.8.1
)

require (
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/smarty/assertions v1.15.1 // indirect
	go.opentelemetry.io/otel v1.21.0 // indirect
	go.opentelemetry.io/otel/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.21.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
)

require (
	github.com/ONSdigital/dp-api-clients-go/v2 v2.254.1 // indirect
	github.com/ONSdigital/dp-net/v2 v2.11.2 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect

)
