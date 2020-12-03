module github.com/bigstepinc/metalcloud-cli

go 1.12

require (
	github.com/bigstepinc/metal-cloud-sdk-go v1.6.0
	github.com/golang/mock v1.4.4
	github.com/metalsoft-io/tableformatter v1.0.4
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/onsi/gomega v1.10.3
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/yaml.v2 v2.3.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)

//replace github.com/bigstepinc/metal-cloud-sdk-go => /Users/alex/code/go/src/github.com/bigstepinc/metal-cloud-sdk-go
