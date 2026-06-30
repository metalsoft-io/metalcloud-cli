package fabric_template_config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestExamplesAreValidYAML(t *testing.T) {
	for name, doc := range map[string]string{"freeform": ExampleFreeformYAML(), "bgp": ExampleBgpYAML()} {
		var m map[string]interface{}
		if err := yaml.Unmarshal([]byte(doc), &m); err != nil {
			t.Errorf("%s example is not valid YAML: %v", name, err)
		}
		if _, ok := m[name]; !ok {
			t.Errorf("%s example is missing its %q section", name, name)
		}
	}
}
