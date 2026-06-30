package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// templatePlanFlags are the switch-config plan sections (ordering / topology /
// p2p) that configure-freeform and configure-bgp need to compute per-device
// variables. They mirror the configure-switches flag names/semantics.
type templatePlanFlags struct {
	ordering string

	topoLeafSpine          bool
	topoLeafSpineLPP       string
	topoSpineSuperSpine    bool
	topoSpineSuperSpineLPP string

	topoLeafHost            bool
	topoLeafHostNodeCount   int
	topoLeafHostNodes       []int
	topoLeafHostPortPattern string
	topoLeafHostNicNames    []string
	topoLeafHostDescription string

	p2pPoolLeafSpine       string
	p2pPoolSpineSuperSpine string
	p2pPoolLeafHost        string
	p2pMtu                 int32
}

var planFlagNames = []string{
	"ordering",
	"topology-leaf-spine", "topology-leaf-spine-links-per-pair",
	"topology-spine-super-spine", "topology-spine-super-spine-links-per-pair",
	"topology-leaf-host", "topology-leaf-host-node-count", "topology-leaf-host-nodes",
	"topology-leaf-host-port-pattern", "topology-leaf-host-nic-names",
	"topology-leaf-host-description-template",
	"p2p-pool-leaf-spine", "p2p-pool-spine-super-spine", "p2p-pool-leaf-host", "p2p-mtu",
}

func registerPlanFlags(cmd *cobra.Command, pf *templatePlanFlags) {
	f := cmd.Flags()
	f.StringVar(&pf.ordering, "ordering", "managementAddress", "Device ordering: managementAddress | identifierString | id.")
	f.BoolVar(&pf.topoLeafSpine, "topology-leaf-spine", false, "Enable leaf<->spine pairing.")
	f.StringVar(&pf.topoLeafSpineLPP, "topology-leaf-spine-links-per-pair", "", "Leaf<->spine links per pair: 'auto' or an integer.")
	f.BoolVar(&pf.topoSpineSuperSpine, "topology-spine-super-spine", false, "Enable spine<->superspine pairing (3-tier only).")
	f.StringVar(&pf.topoSpineSuperSpineLPP, "topology-spine-super-spine-links-per-pair", "", "Spine<->superspine links per pair: 'auto' or an integer.")
	f.BoolVar(&pf.topoLeafHost, "topology-leaf-host", false, "Enable leaf->host downlinks (for the /26 aggregates).")
	f.IntVar(&pf.topoLeafHostNodeCount, "topology-leaf-host-node-count", 0, "Number of host port-pairs per leaf.")
	f.IntSliceVar(&pf.topoLeafHostNodes, "topology-leaf-host-nodes", nil, "Exact 0-based node indices (mutually exclusive with node-count).")
	f.StringVar(&pf.topoLeafHostPortPattern, "topology-leaf-host-port-pattern", "", "Leaf host port pattern, e.g. swp{port}s{sub}.")
	f.StringSliceVar(&pf.topoLeafHostNicNames, "topology-leaf-host-nic-names", nil, "Remote host NIC names (even count).")
	f.StringVar(&pf.topoLeafHostDescription, "topology-leaf-host-description-template", "", "Leaf->host description template.")
	f.StringVar(&pf.p2pPoolLeafSpine, "p2p-pool-leaf-spine", "", "Leaf<->spine /31 pool.")
	f.StringVar(&pf.p2pPoolSpineSuperSpine, "p2p-pool-spine-super-spine", "", "Spine<->superspine /31 pool.")
	f.StringVar(&pf.p2pPoolLeafHost, "p2p-pool-leaf-host", "", "Leaf->host /31 pool.")
	f.Int32Var(&pf.p2pMtu, "p2p-mtu", 0, "MTU applied to created links.")
}

// addPlanSections adds the topology / p2p / ordering sections built from the
// plan flags that were set into doc.
func addPlanSections(cmd *cobra.Command, pf *templatePlanFlags, doc map[string]interface{}) error {
	f := cmd.Flags()
	if f.Changed("ordering") {
		doc["ordering"] = pf.ordering
	}

	topology := map[string]interface{}{}
	if leafSpine, present, err := buildLayerFlags(f, "topology-leaf-spine", pf.topoLeafSpine, "topology-leaf-spine-links-per-pair", pf.topoLeafSpineLPP); err != nil {
		return err
	} else if present {
		topology["leafSpine"] = leafSpine
	}
	if spineSsp, present, err := buildLayerFlags(f, "topology-spine-super-spine", pf.topoSpineSuperSpine, "topology-spine-super-spine-links-per-pair", pf.topoSpineSuperSpineLPP); err != nil {
		return err
	} else if present {
		topology["spineSuperSpine"] = spineSsp
	}
	leafHost := map[string]interface{}{}
	leafHostPresent := f.Changed("topology-leaf-host") && pf.topoLeafHost
	if f.Changed("topology-leaf-host-node-count") {
		leafHost["nodeCount"] = pf.topoLeafHostNodeCount
		leafHostPresent = true
	}
	if f.Changed("topology-leaf-host-nodes") {
		leafHost["nodes"] = pf.topoLeafHostNodes
		leafHostPresent = true
	}
	if f.Changed("topology-leaf-host-port-pattern") {
		leafHost["portPattern"] = pf.topoLeafHostPortPattern
		leafHostPresent = true
	}
	if f.Changed("topology-leaf-host-nic-names") {
		leafHost["nicNames"] = pf.topoLeafHostNicNames
		leafHostPresent = true
	}
	if f.Changed("topology-leaf-host-description-template") {
		leafHost["descriptionTemplate"] = pf.topoLeafHostDescription
		leafHostPresent = true
	}
	if leafHostPresent {
		topology["leafHost"] = leafHost
	}
	if len(topology) > 0 {
		doc["topology"] = topology
	}

	p2p := map[string]interface{}{}
	pools := map[string]interface{}{}
	if f.Changed("p2p-pool-leaf-spine") {
		pools["leafSpine"] = pf.p2pPoolLeafSpine
	}
	if f.Changed("p2p-pool-spine-super-spine") {
		pools["spineSuperSpine"] = pf.p2pPoolSpineSuperSpine
	}
	if f.Changed("p2p-pool-leaf-host") {
		pools["leafHost"] = pf.p2pPoolLeafHost
	}
	if len(pools) > 0 {
		p2p["pools"] = pools
	}
	if f.Changed("p2p-mtu") {
		p2p["mtu"] = pf.p2pMtu
	}
	if len(p2p) > 0 {
		doc["p2p"] = p2p
	}
	return nil
}

// ---- freeform ----------------------------------------------------------------

var configureFreeformFlags = struct {
	mode            string
	templatePath    string
	templateLabel   string
	profilePriority int
	applyMode       string
	hgxPrefix       string
	templatePlanFlags
}{}

var freeformDetailFlags = append([]string{
	"mode", "template-path", "template-label", "profile-priority", "apply-mode", "hgx-prefix",
}, planFlagNames...)

func registerConfigureFreeformFlags(cmd *cobra.Command) {
	ff := &configureFreeformFlags
	f := cmd.Flags()
	f.StringVar(&ff.mode, "mode", "", "Fabric mode: purel3 | l3evpn (must match bgp.mode).")
	f.StringVar(&ff.templatePath, "template-path", "", "Path to the base freeform .j2 template body.")
	f.StringVar(&ff.templateLabel, "template-label", "", "Find-or-create template label (default spectrumx-freeform).")
	f.IntVar(&ff.profilePriority, "profile-priority", 0, "Profile priority (default 50).")
	f.StringVar(&ff.applyMode, "apply-mode", "", "Profile apply mode: once | always (default once).")
	f.StringVar(&ff.hgxPrefix, "hgx-prefix", "", "Tenant HGX supernet prefix (default: derived from the leafHost pool).")
	registerPlanFlags(cmd, &ff.templatePlanFlags)
}

func buildFreeformConfigFromFlags(cmd *cobra.Command) ([]byte, error) {
	f := cmd.Flags()
	ff := &configureFreeformFlags
	doc := map[string]interface{}{}
	section := map[string]interface{}{}
	if f.Changed("mode") {
		section["mode"] = ff.mode
	}
	if f.Changed("template-path") {
		section["templatePath"] = ff.templatePath
	}
	if f.Changed("template-label") {
		section["templateLabel"] = ff.templateLabel
	}
	if f.Changed("profile-priority") {
		section["profilePriority"] = ff.profilePriority
	}
	if f.Changed("apply-mode") {
		section["applyMode"] = ff.applyMode
	}
	if f.Changed("hgx-prefix") {
		section["hgxPrefix"] = ff.hgxPrefix
	}
	if len(section) == 0 {
		return nil, fmt.Errorf("specify --config-source or at least the freeform flags (--mode, --template-path)")
	}
	doc["freeform"] = section
	if err := addPlanSections(cmd, &ff.templatePlanFlags, doc); err != nil {
		return nil, err
	}
	return yaml.Marshal(doc)
}

// ---- bgp ----------------------------------------------------------------------

var configureBgpFlags = struct {
	mode      string
	applyMode string

	templatePath    string
	templateLabel   string
	profilePriority int

	overlayTemplatePath    string
	overlayTemplateLabel   string
	overlayProfilePriority int

	pfcTemplatePath    string
	pfcTemplateLabel   string
	pfcProfilePriority int

	vrfTemplatePath  string
	vrfTemplateLabel string
	templatePlanFlags
}{}

var bgpDetailFlags = append([]string{
	"mode", "apply-mode",
	"template-path", "template-label", "profile-priority",
	"overlay-template-path", "overlay-template-label", "overlay-profile-priority",
	"pfc-template-path", "pfc-template-label", "pfc-profile-priority",
	"vrf-template-path", "vrf-template-label",
}, planFlagNames...)

func registerConfigureBgpFlags(cmd *cobra.Command) {
	bf := &configureBgpFlags
	f := cmd.Flags()
	f.StringVar(&bf.mode, "mode", "", "Fabric mode: purel3 | l3evpn (must match freeform.mode).")
	f.StringVar(&bf.applyMode, "apply-mode", "", "Profile apply mode: once | always (default once).")
	f.StringVar(&bf.templatePath, "template-path", "", "Path to the BGP underlay .j2 template body.")
	f.StringVar(&bf.templateLabel, "template-label", "", "Underlay template label (default spectrumx-bgp-underlay).")
	f.IntVar(&bf.profilePriority, "profile-priority", 0, "Underlay profile priority (default 60).")
	f.StringVar(&bf.overlayTemplatePath, "overlay-template-path", "", "Path to the EVPN overlay .j2 template body.")
	f.StringVar(&bf.overlayTemplateLabel, "overlay-template-label", "", "Overlay template label (default spectrumx-bgp-overlay).")
	f.IntVar(&bf.overlayProfilePriority, "overlay-profile-priority", 0, "Overlay profile priority (default 61).")
	f.StringVar(&bf.pfcTemplatePath, "pfc-template-path", "", "Path to the QoS PFC .j2 template body.")
	f.StringVar(&bf.pfcTemplateLabel, "pfc-template-label", "", "PFC template label (default spectrumx-qos-pfc).")
	f.IntVar(&bf.pfcProfilePriority, "pfc-profile-priority", 0, "PFC profile priority (default 62).")
	f.StringVar(&bf.vrfTemplatePath, "vrf-template-path", "", "Path to the route-domain VRF .j2 template body.")
	f.StringVar(&bf.vrfTemplateLabel, "vrf-template-label", "", "VRF template label (default switch-configure-vrf-create).")
	registerPlanFlags(cmd, &bf.templatePlanFlags)
}

func buildBgpConfigFromFlags(cmd *cobra.Command) ([]byte, error) {
	f := cmd.Flags()
	bf := &configureBgpFlags
	doc := map[string]interface{}{}
	section := map[string]interface{}{}
	set := func(flag, key string, value interface{}) {
		if f.Changed(flag) {
			section[key] = value
		}
	}
	set("mode", "mode", bf.mode)
	set("apply-mode", "applyMode", bf.applyMode)
	set("template-path", "templatePath", bf.templatePath)
	set("template-label", "templateLabel", bf.templateLabel)
	set("profile-priority", "profilePriority", bf.profilePriority)
	set("overlay-template-path", "overlayTemplatePath", bf.overlayTemplatePath)
	set("overlay-template-label", "overlayTemplateLabel", bf.overlayTemplateLabel)
	set("overlay-profile-priority", "overlayProfilePriority", bf.overlayProfilePriority)
	set("pfc-template-path", "pfcTemplatePath", bf.pfcTemplatePath)
	set("pfc-template-label", "pfcTemplateLabel", bf.pfcTemplateLabel)
	set("pfc-profile-priority", "pfcProfilePriority", bf.pfcProfilePriority)
	set("vrf-template-path", "vrfTemplatePath", bf.vrfTemplatePath)
	set("vrf-template-label", "vrfTemplateLabel", bf.vrfTemplateLabel)
	if len(section) == 0 {
		return nil, fmt.Errorf("specify --config-source or at least the bgp flags (--mode, --template-path, ...)")
	}
	doc["bgp"] = section
	if err := addPlanSections(cmd, &bf.templatePlanFlags, doc); err != nil {
		return nil, err
	}
	return yaml.Marshal(doc)
}
