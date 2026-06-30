package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func TestBuildFreeformConfigFromFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "configure-freeform"}
	registerConfigureFreeformFlags(cmd)
	for k, v := range map[string]string{
		"mode":                          "l3evpn",
		"template-path":                 "/tmp/freeform.j2",
		"hgx-prefix":                    "172.0.0.0/8",
		"topology-leaf-spine":           "true",
		"topology-leaf-host-node-count": "32",
		"p2p-pool-leaf-host":            "172.16.0.0/12",
		"p2p-mtu":                       "9216",
	} {
		if err := cmd.Flags().Set(k, v); err != nil {
			t.Fatalf("set %s: %v", k, err)
		}
	}
	data, err := buildFreeformConfigFromFlags(cmd)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	var doc map[string]interface{}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("unmarshal: %v\n%s", err, data)
	}
	ff := doc["freeform"].(map[string]interface{})
	if ff["mode"] != "l3evpn" || ff["templatePath"] != "/tmp/freeform.j2" || ff["hgxPrefix"] != "172.0.0.0/8" {
		t.Errorf("freeform section wrong: %v", ff)
	}
	topo := doc["topology"].(map[string]interface{})
	if _, ok := topo["leafSpine"]; !ok {
		t.Errorf("topology.leafSpine missing: %v", topo)
	}
	if lh := topo["leafHost"].(map[string]interface{}); lh["nodeCount"] != 32 {
		t.Errorf("leafHost nodeCount wrong: %v", lh)
	}
	p2p := doc["p2p"].(map[string]interface{})
	if p2p["mtu"] != 9216 {
		t.Errorf("p2p mtu wrong: %v", p2p)
	}
}

func TestBuildBgpConfigFromFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "configure-bgp"}
	registerConfigureBgpFlags(cmd)
	for k, v := range map[string]string{
		"mode":                                      "l3evpn",
		"template-path":                             "/tmp/underlay.j2",
		"overlay-template-path":                     "/tmp/overlay.j2",
		"pfc-template-path":                         "/tmp/pfc.j2",
		"vrf-template-path":                         "/tmp/vrf.j2",
		"topology-leaf-spine-links-per-pair":        "auto",
		"topology-spine-super-spine-links-per-pair": "4",
		"p2p-pool-leaf-spine":                       "10.254.0.0/16",
		"p2p-mtu":                                   "9216",
	} {
		if err := cmd.Flags().Set(k, v); err != nil {
			t.Fatalf("set %s: %v", k, err)
		}
	}
	data, err := buildBgpConfigFromFlags(cmd)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	var doc map[string]interface{}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("unmarshal: %v\n%s", err, data)
	}
	bgp := doc["bgp"].(map[string]interface{})
	if bgp["mode"] != "l3evpn" || bgp["templatePath"] != "/tmp/underlay.j2" ||
		bgp["overlayTemplatePath"] != "/tmp/overlay.j2" || bgp["vrfTemplatePath"] != "/tmp/vrf.j2" {
		t.Errorf("bgp section wrong: %v", bgp)
	}
	topo := doc["topology"].(map[string]interface{})
	ls := topo["leafSpine"].(map[string]interface{})
	if ls["linksPerPair"] != "auto" {
		t.Errorf("leafSpine linksPerPair should be 'auto', got %v", ls)
	}
	ssp := topo["spineSuperSpine"].(map[string]interface{})
	if ssp["linksPerPair"] != 4 {
		t.Errorf("spineSuperSpine linksPerPair should be 4, got %v", ssp)
	}
	p2p := doc["p2p"].(map[string]interface{})
	pools := p2p["pools"].(map[string]interface{})
	if pools["leafSpine"] != "10.254.0.0/16" {
		t.Errorf("p2p pool leafSpine wrong: %v", pools)
	}
}

func TestBuildTemplateConfigFromFlagsEmpty(t *testing.T) {
	ff := &cobra.Command{Use: "configure-freeform"}
	registerConfigureFreeformFlags(ff)
	if _, err := buildFreeformConfigFromFlags(ff); err == nil {
		t.Error("freeform: expected error with no flags set")
	}
	bf := &cobra.Command{Use: "configure-bgp"}
	registerConfigureBgpFlags(bf)
	if _, err := buildBgpConfigFromFlags(bf); err == nil {
		t.Error("bgp: expected error with no flags set")
	}
}
