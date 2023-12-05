module github.com/ONSdigital/dp-s3/v2

go 1.20

retract (
	v2.1.0-beta.2 // Contains retractions only
	[v2.1.0-beta, v2.1.0-beta.1] // Code never made into main and refers to non-released dependencies
)

require (
	github.com/ONSdigital/dp-healthcheck v1.6.1
	github.com/ONSdigital/log.go/v2 v2.4.1
	github.com/aws/aws-sdk-go v1.44.247
	github.com/smartystreets/goconvey v1.8.0
)

require (
	github.com/ONSdigital/dp-api-clients-go/v2 v2.252.1 // indirect
	github.com/ONSdigital/dp-net/v2 v2.9.1 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/smartystreets/assertions v1.13.1 // indirect
	golang.org/x/net v0.17.0 //indirect
	golang.org/x/sys v0.13.0 // indirect
)
