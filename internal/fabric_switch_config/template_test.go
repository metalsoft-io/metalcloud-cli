package fabric_switch_config

import "testing"

func TestExpandTemplateForms(t *testing.T) {
	tags := map[string]string{"nvidia/pod-id": "5", "nvidia/rail-group-id": "1"}
	ordinals := map[string]int{"nvidia/ssp-group-id": 2}
	ordinalBy := func(key string) (int, error) { return ordinals[key], nil }

	cases := []struct {
		name     string
		template string
		want     string
	}{
		{"tag zero-padded", "pod{tag:nvidia/pod-id:02d}", "pod05"},
		{"tag raw", "r{tag:nvidia/rail-group-id}", "r1"},
		{"ordinalBy 1-based", "s{ordinalBy:nvidia/ssp-group-id}", "s2"},
		{"ordinalBy0 padded", "s{ordinalBy0:nvidia/ssp-group-id:02d}", "s01"},
		{"plain value padded", "h{node:02d}", "h07"},
		{"plain value raw", "{nic}", "enp26s0f0np0"},
		{"hex", "{node:x}", "7"},
		{"mixed", "to_hgx-pod{tag:nvidia/pod-id:02d}-h{node:02d}_{nic}", "to_hgx-pod05-h07_enp26s0f0np0"},
	}
	for _, c := range cases {
		got, err := expandTemplate(c.template, tags, ordinalBy, map[string]any{"node": 7, "nic": "enp26s0f0np0"})
		if err != nil {
			t.Errorf("%s: unexpected error %v", c.name, err)
			continue
		}
		if got != c.want {
			t.Errorf("%s: %q -> %q, want %q", c.name, c.template, got, c.want)
		}
	}
}

func TestExpandTemplateErrors(t *testing.T) {
	if _, err := expandTemplate("{tag:missing}", map[string]string{}, nil, nil); err == nil {
		t.Error("missing tag should error")
	}
	if _, err := expandTemplate("{unknown}", nil, nil, map[string]any{}); err == nil {
		t.Error("unknown placeholder should error")
	}
	if _, err := expandTemplate("{ordinalBy:x}", nil, nil, nil); err == nil {
		t.Error("ordinalBy without callback should error")
	}
}

func TestLogicalToSwp(t *testing.T) {
	cases := map[int]string{1: "swp1s0", 2: "swp1s1", 65: "swp33s0", 128: "swp64s1"}
	for logical, want := range cases {
		if got := logicalToSwp(logical); got != want {
			t.Errorf("logicalToSwp(%d) = %q, want %q", logical, got, want)
		}
	}
}
