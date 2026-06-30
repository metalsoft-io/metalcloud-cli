package fabric_template_config

import (
	"fmt"
	"net"
	"os"

	"github.com/metalsoft-io/metalcloud-cli/pkg/utils"
)

var validModes = map[string]bool{"purel3": true, "l3evpn": true}
var validApplyModes = map[string]bool{"once": true, "always": true}

// Default labels / priorities (mirror configure_freeform.py / configure_bgp.py).
const (
	defaultFreeformLabel    = "spectrumx-freeform"
	defaultFreeformPriority = 50

	defaultUnderlayLabel    = "spectrumx-bgp-underlay"
	defaultUnderlayPriority = 60
	defaultOverlayLabel     = "spectrumx-bgp-overlay"
	defaultOverlayPriority  = 61
	defaultPfcLabel         = "spectrumx-qos-pfc"
	defaultPfcPriority      = 62
	defaultVrfLabel         = "switch-configure-vrf-create"

	defaultApplyMode = "once"
)

// vrfTemplateAnnotations bind the route-domain VRF template to the engine action
// (no per-device profile).
var vrfTemplateAnnotations = map[string]string{"action": "switch-configure-vrf-create", "position": "leaf"}

// templateSpec is a registered template: its find-or-create label, the profile
// priority, and the (decoded) template body read from its .j2 file.
type templateSpec struct {
	Label    string
	Priority float32
	Text     string
}

// FreeformConfig is the resolved `freeform:` section.
type FreeformConfig struct {
	Mode      string
	HgxPrefix string // override; "" => derive from the leaf-host pool
	ApplyMode string
	Template  templateSpec
}

// BgpConfig is the resolved `bgp:` section.
type BgpConfig struct {
	Mode      string
	ApplyMode string
	Underlay  templateSpec
	Overlay   templateSpec
	Pfc       templateSpec
	Vrf       templateSpec // priority unused (action-bound, no profile)
}

type rawTemplateDoc struct {
	Freeform *rawFreeform `json:"freeform" yaml:"freeform"`
	Bgp      *rawBgp      `json:"bgp" yaml:"bgp"`
}

type rawFreeform struct {
	Mode            string `json:"mode" yaml:"mode"`
	HgxPrefix       string `json:"hgxPrefix" yaml:"hgxPrefix"`
	TemplatePath    string `json:"templatePath" yaml:"templatePath"`
	TemplateLabel   string `json:"templateLabel" yaml:"templateLabel"`
	ProfilePriority *int   `json:"profilePriority" yaml:"profilePriority"`
	ApplyMode       string `json:"applyMode" yaml:"applyMode"`
}

type rawBgp struct {
	Mode      string `json:"mode" yaml:"mode"`
	ApplyMode string `json:"applyMode" yaml:"applyMode"`

	TemplatePath    string `json:"templatePath" yaml:"templatePath"`
	TemplateLabel   string `json:"templateLabel" yaml:"templateLabel"`
	ProfilePriority *int   `json:"profilePriority" yaml:"profilePriority"`

	OverlayTemplatePath    string `json:"overlayTemplatePath" yaml:"overlayTemplatePath"`
	OverlayTemplateLabel   string `json:"overlayTemplateLabel" yaml:"overlayTemplateLabel"`
	OverlayProfilePriority *int   `json:"overlayProfilePriority" yaml:"overlayProfilePriority"`

	PfcTemplatePath    string `json:"pfcTemplatePath" yaml:"pfcTemplatePath"`
	PfcTemplateLabel   string `json:"pfcTemplateLabel" yaml:"pfcTemplateLabel"`
	PfcProfilePriority *int   `json:"pfcProfilePriority" yaml:"pfcProfilePriority"`

	VrfTemplatePath  string `json:"vrfTemplatePath" yaml:"vrfTemplatePath"`
	VrfTemplateLabel string `json:"vrfTemplateLabel" yaml:"vrfTemplateLabel"`
}

func readTemplateFile(key, path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("%s is required (path to the .j2 template file)", key)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("%s: cannot read template file %q: %w", key, path, err)
	}
	return string(data), nil
}

func resolveLabel(value, def string) string {
	if value != "" {
		return value
	}
	return def
}

func resolvePriority(value *int, def int) (float32, error) {
	if value == nil {
		return float32(def), nil
	}
	if *value < 0 {
		return 0, fmt.Errorf("profile priority must be a non-negative integer")
	}
	return float32(*value), nil
}

func resolveApplyMode(value string) (string, error) {
	if value == "" {
		return defaultApplyMode, nil
	}
	if !validApplyModes[value] {
		return "", fmt.Errorf("applyMode must be one of [once always], got %q", value)
	}
	return value, nil
}

// LoadFreeformConfig parses the `freeform:` section from the config bytes.
func LoadFreeformConfig(data []byte) (*FreeformConfig, error) {
	var doc rawTemplateDoc
	if err := utils.UnmarshalContent(data, &doc); err != nil {
		return nil, err
	}
	if doc.Freeform == nil {
		return nil, fmt.Errorf("config has no 'freeform' section; nothing to do")
	}
	raw := doc.Freeform
	if !validModes[raw.Mode] {
		return nil, fmt.Errorf("freeform.mode must be one of [purel3 l3evpn], got %q", raw.Mode)
	}
	if doc.Bgp != nil && doc.Bgp.Mode != "" && doc.Bgp.Mode != raw.Mode {
		return nil, fmt.Errorf("freeform.mode %q differs from bgp.mode %q; the base and BGP layers should render the same fabric mode", raw.Mode, doc.Bgp.Mode)
	}
	if raw.HgxPrefix != "" {
		if _, _, err := net.ParseCIDR(raw.HgxPrefix); err != nil {
			return nil, fmt.Errorf("freeform.hgxPrefix: invalid network %q: %w", raw.HgxPrefix, err)
		}
	}
	text, err := readTemplateFile("freeform.templatePath", raw.TemplatePath)
	if err != nil {
		return nil, err
	}
	priority, err := resolvePriority(raw.ProfilePriority, defaultFreeformPriority)
	if err != nil {
		return nil, fmt.Errorf("freeform.%w", err)
	}
	applyMode, err := resolveApplyMode(raw.ApplyMode)
	if err != nil {
		return nil, fmt.Errorf("freeform.%w", err)
	}
	return &FreeformConfig{
		Mode:      raw.Mode,
		HgxPrefix: raw.HgxPrefix,
		ApplyMode: applyMode,
		Template:  templateSpec{Label: resolveLabel(raw.TemplateLabel, defaultFreeformLabel), Priority: priority, Text: text},
	}, nil
}

// LoadBgpConfig parses the `bgp:` section from the config bytes.
func LoadBgpConfig(data []byte) (*BgpConfig, error) {
	var doc rawTemplateDoc
	if err := utils.UnmarshalContent(data, &doc); err != nil {
		return nil, err
	}
	if doc.Bgp == nil {
		return nil, fmt.Errorf("config has no 'bgp' section; nothing to do")
	}
	raw := doc.Bgp
	if !validModes[raw.Mode] {
		return nil, fmt.Errorf("bgp.mode must be one of [purel3 l3evpn], got %q", raw.Mode)
	}
	applyMode, err := resolveApplyMode(raw.ApplyMode)
	if err != nil {
		return nil, fmt.Errorf("bgp.%w", err)
	}

	underlayText, err := readTemplateFile("bgp.templatePath", raw.TemplatePath)
	if err != nil {
		return nil, err
	}
	underlayPriority, err := resolvePriority(raw.ProfilePriority, defaultUnderlayPriority)
	if err != nil {
		return nil, fmt.Errorf("bgp.%w", err)
	}
	overlayText, err := readTemplateFile("bgp.overlayTemplatePath", raw.OverlayTemplatePath)
	if err != nil {
		return nil, err
	}
	overlayPriority, err := resolvePriority(raw.OverlayProfilePriority, defaultOverlayPriority)
	if err != nil {
		return nil, fmt.Errorf("bgp.%w", err)
	}
	pfcText, err := readTemplateFile("bgp.pfcTemplatePath", raw.PfcTemplatePath)
	if err != nil {
		return nil, err
	}
	pfcPriority, err := resolvePriority(raw.PfcProfilePriority, defaultPfcPriority)
	if err != nil {
		return nil, fmt.Errorf("bgp.%w", err)
	}
	vrfText, err := readTemplateFile("bgp.vrfTemplatePath", raw.VrfTemplatePath)
	if err != nil {
		return nil, err
	}

	return &BgpConfig{
		Mode:      raw.Mode,
		ApplyMode: applyMode,
		Underlay:  templateSpec{Label: resolveLabel(raw.TemplateLabel, defaultUnderlayLabel), Priority: underlayPriority, Text: underlayText},
		Overlay:   templateSpec{Label: resolveLabel(raw.OverlayTemplateLabel, defaultOverlayLabel), Priority: overlayPriority, Text: overlayText},
		Pfc:       templateSpec{Label: resolveLabel(raw.PfcTemplateLabel, defaultPfcLabel), Priority: pfcPriority, Text: pfcText},
		Vrf:       templateSpec{Label: resolveLabel(raw.VrfTemplateLabel, defaultVrfLabel), Text: vrfText},
	}, nil
}
