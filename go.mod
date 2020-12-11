module github.com/bigstepinc/metalcloud-cli

go 1.12

require (
	github.com/bigstepinc/metal-cloud-sdk-go/v2 v2.0.1
	// github.com/bigstepinc/metal-cloud-sdk-go v1.5.2
	// github.com/bigstepinc/metal-cloud-sdk-go v2.0.0
	github.com/golang/mock v1.4.4
	github.com/kr/text v0.2.0 // indirect
	github.com/metalsoft-io/tableformatter v1.0.4
	github.com/onsi/gomega v1.10.4
	github.com/savaki/jq v0.0.0-20161209013833-0e6baecebbf8
	golang.org/x/crypto v0.0.0-20201208171446-5f87f3452ae9
	golang.org/x/net v0.0.0-20201209123823-ac852fbbde11 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)

//replace github.com/bigstepinc/metal-cloud-sdk-go => /Users/alex/code/go/src/github.com/bigstepinc/metal-cloud-sdk-go
