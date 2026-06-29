package fabric

import "testing"

func TestDeepMergeStringMaps(t *testing.T) {
	base := map[string]interface{}{
		"driver":   "cumulus_linux",
		"username": "admin",
		"tagsMap":  map[string]interface{}{"environment": "production", "managed_by": "ra"},
	}
	override := map[string]interface{}{
		"identifierString": "leaf-01",
		"username":         "root", // per-switch wins
		"tagsMap":          map[string]interface{}{"nvidia/pod-id": "5", "environment": "lab"},
	}
	merged := deepMergeStringMaps(base, override)

	if merged["driver"] != "cumulus_linux" || merged["identifierString"] != "leaf-01" {
		t.Errorf("scalar merge wrong: %+v", merged)
	}
	if merged["username"] != "root" {
		t.Errorf("override should win: username=%v", merged["username"])
	}
	tags := merged["tagsMap"].(map[string]interface{})
	if tags["environment"] != "lab" || tags["managed_by"] != "ra" || tags["nvidia/pod-id"] != "5" {
		t.Errorf("nested tagsMap merge wrong: %+v", tags)
	}
	// base must be untouched.
	if base["username"] != "admin" {
		t.Errorf("base mutated: %+v", base)
	}
}

func TestCoerceTagsMap(t *testing.T) {
	sw := map[string]interface{}{
		"tagsMap": map[string]interface{}{
			"rack":    42,
			"managed": true,
			"off":     false,
			"name":    "keep",
		},
	}
	coerceTagsMap(sw)
	tags := sw["tagsMap"].(map[string]interface{})
	if tags["rack"] != "42" || tags["managed"] != "true" || tags["off"] != "false" || tags["name"] != "keep" {
		t.Errorf("coercion wrong: %+v", tags)
	}
}

func TestValidateSwitchMap(t *testing.T) {
	valid := map[string]interface{}{
		"driver": "cumulus_linux", "position": "leaf", "username": "admin",
		"managementPassword": "x", "managementAddress": "10.0.0.1",
		"managementPort": 22, "identifierString": "leaf-01",
	}
	if errs := validateSwitchMap(valid); len(errs) != 0 {
		t.Errorf("valid switch reported errors: %v", errs)
	}

	missing := map[string]interface{}{"driver": "cumulus_linux"}
	if errs := validateSwitchMap(missing); len(errs) == 0 {
		t.Error("missing required fields should be reported")
	}

	badTags := map[string]interface{}{
		"driver": "cumulus_linux", "position": "leaf", "username": "admin",
		"managementPassword": "x", "managementAddress": "10.0.0.1",
		"managementPort": 22, "identifierString": "leaf-01",
		"tagsMap": map[string]interface{}{"k": 5}, // non-string => invalid (pre-coercion)
	}
	if errs := validateSwitchMap(badTags); len(errs) == 0 {
		t.Error("non-string tag value should be reported")
	}
}

func TestFindExistingDeviceId(t *testing.T) {
	byMgmt := map[string]int64{"10.0.0.1": 11}
	byIdent := map[string]int64{"leaf-01": 22}
	bySerial := map[string]int64{"SN123": 33}

	cases := []struct {
		sw   map[string]interface{}
		want int64
		ok   bool
	}{
		{map[string]interface{}{"managementAddress": "10.0.0.1"}, 11, true},
		{map[string]interface{}{"identifierString": "leaf-01"}, 22, true},
		{map[string]interface{}{"serialNumber": "SN123"}, 33, true},
		{map[string]interface{}{"managementAddress": "10.0.0.9"}, 0, false},
		// mgmt takes precedence over ident
		{map[string]interface{}{"managementAddress": "10.0.0.1", "identifierString": "leaf-01"}, 11, true},
	}
	for i, c := range cases {
		id, ok := findExistingDeviceId(c.sw, byMgmt, byIdent, bySerial)
		if ok != c.ok || (ok && id != c.want) {
			t.Errorf("case %d: got (%d,%v), want (%d,%v)", i, id, ok, c.want, c.ok)
		}
	}
}

func TestBuildCreateNetworkDevice(t *testing.T) {
	sw := map[string]interface{}{
		"driver": "cumulus_linux", "position": "leaf", "username": "admin",
		"managementPassword": "x", "managementAddress": "10.0.0.1",
		"managementPort": 22, "identifierString": "leaf-01",
		"tagsMap": map[string]interface{}{"nvidia/pod-id": "5"},
	}
	dev, err := buildCreateNetworkDevice(sw, 11)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if dev.SiteId == nil || *dev.SiteId != 11 {
		t.Errorf("siteId should be injected, got %v", dev.SiteId)
	}
	if dev.Position != "leaf" {
		t.Errorf("position = %q, want leaf", dev.Position)
	}
	if dev.IdentifierString == nil || *dev.IdentifierString != "leaf-01" {
		t.Errorf("identifierString not mapped: %v", dev.IdentifierString)
	}
	if dev.TagsMap == nil || (*dev.TagsMap)["nvidia/pod-id"] != "5" {
		t.Errorf("tagsMap not mapped: %v", dev.TagsMap)
	}
}
