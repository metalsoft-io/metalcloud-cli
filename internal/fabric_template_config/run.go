package fabric_template_config

import (
	"fmt"
	"strings"

	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
)

const (
	freeformDescription = "Spectrum-X base freeform device configuration (hostname, RoCE/QoS, AR, telemetry, BFD; l3evpn EVPN/VXLAN data plane)"
	underlayDescription = "Spectrum-X BGP underlay"
	overlayDescription  = "Spectrum-X EVPN overlay (RR mesh)"
	pfcDescription      = "Spectrum-X QoS PFC defaults"
	vrfDescription      = "Spectrum-X route-domain tenant VRF (EVPN VNI/RD + leaf /26 aggregate-routes from device customVariables)"
)

// RunFreeform registers the base freeform template + one profile per switch.
// verify pushes each device's render through the engine first.
func RunFreeform(client TemplateClient, data []byte, fabricId int64, dryRun, verify bool) (*Result, error) {
	freeform, err := LoadFreeformConfig(data)
	if err != nil {
		return nil, err
	}
	plan, planConfig, warnings, err := buildPlan(client, data, fabricId)
	if err != nil {
		return nil, err
	}
	isThreeTier := threeTier(plan.groups)
	hgx := hgxPrefix(planConfig, isThreeTier, freeform.HgxPrefix)
	variables, err := computeFreeformVariables(plan.groups, plan.state, plan.records, freeform.Mode, hgx)
	if err != nil {
		return nil, err
	}

	logger.Get().Debug().Msgf("freeform: mode=%s, hgx_prefix=%s, %d switch(es), template=%q (priority %g)",
		freeform.Mode, hgx, len(plan.devices), freeform.Template.Label, freeform.Template.Priority)
	for _, dev := range plan.devices {
		logger.Get().Debug().Msgf("[%s] freeform vars: %s", dev.Label(), varSummary(variables[dev.Id]))
	}

	r := &runner{client: client, fabricId: fabricId, dryRun: dryRun, apply: freeform.ApplyMode,
		result: &Result{Counters: map[string]int{}, Warnings: warnings}}

	if verify {
		m := r.verifyRender(freeform.Template, plan.devices, func(rec *deviceRecord) map[string]interface{} {
			return renderContextFreeform(&rec.Device, variables[rec.Id], plan.state, plan.records)
		})
		if m > 0 {
			summarize(r.result, dryRun)
			return r.result, fmt.Errorf("render verification failed for %d device(s); not writing", m)
		}
		logger.Get().Info().Msgf("render verification passed for all %d device(s)", len(plan.devices))
	}

	id, usable := r.ensureTemplate(freeform.Template, freeformDescription, nil)
	r.ensureProfiles(id, usable, plan.devices, variables, freeform.Template.Priority)

	summarize(r.result, dryRun)
	if r.result.Failures > 0 {
		return r.result, fmt.Errorf("freeform registration completed with %d failure(s)", r.result.Failures)
	}
	return r.result, nil
}

// RunBgp registers the BGP underlay (+ l3evpn overlay/PFC/VRF) templates and
// per-switch profiles, and reconciles device customVariables.
func RunBgp(client TemplateClient, data []byte, fabricId int64, dryRun, verify bool) (*Result, error) {
	bgp, err := LoadBgpConfig(data)
	if err != nil {
		return nil, err
	}
	plan, planConfig, warnings, err := buildPlan(client, data, fabricId)
	if err != nil {
		return nil, err
	}
	if planConfig.Topology == nil || planConfig.Topology.LeafSpine == nil {
		return nil, fmt.Errorf("'bgp' requires 'topology.leafSpine' (the neighbor set is the link plan)")
	}
	if planConfig.P2p == nil {
		return nil, fmt.Errorf("'bgp' requires 'p2p' (neighbor IPs are the link /31s)")
	}

	variables, err := computeBgpVariables(plan.groups, plan.state, plan.records, bgp.Mode)
	if err != nil {
		return nil, err
	}

	// asn/loopback must already be on the device records (configure-switches first).
	var unconfigured []string
	for _, dev := range plan.devices {
		if dev.Asn == nil || loopbackOf(&dev.Device, plan.state, plan.records) == "" {
			unconfigured = append(unconfigured, dev.Label())
		}
	}
	if len(unconfigured) > 0 {
		return nil, fmt.Errorf("device(s) missing asn/loopbackAddress (run 'fabric configure-switches' first): %s", strings.Join(unconfigured, ", "))
	}

	overlay, err := computeOverlayVariables(plan.groups, plan.state, plan.records, bgp.Mode)
	if err != nil {
		return nil, err
	}
	pfc := computePfcVariables(plan.groups, bgp.Mode)

	var overlayTargets, pfcTargets []*deviceRecord
	for _, dev := range plan.devices {
		if overlayApplies(&dev.Device, overlay[dev.Id]) {
			overlayTargets = append(overlayTargets, dev)
		}
		if pfcApplies(pfc[dev.Id]) {
			pfcTargets = append(pfcTargets, dev)
		}
	}

	logger.Get().Debug().Msgf("bgp: mode=%s, %d switch(es); overlay targets=%d, pfc targets=%d",
		bgp.Mode, len(plan.devices), len(overlayTargets), len(pfcTargets))
	for _, dev := range plan.devices {
		logger.Get().Debug().Msgf("[%s] underlay: {%s}; overlay: {%s}", dev.Label(),
			varSummary(variables[dev.Id]), varSummary(overlay[dev.Id]))
	}

	r := &runner{client: client, fabricId: fabricId, dryRun: dryRun, apply: bgp.ApplyMode,
		result: &Result{Counters: map[string]int{}, Warnings: warnings}}

	if verify {
		m := r.verifyRender(bgp.Underlay, plan.devices, bgpCtx(variables, plan))
		m += r.verifyRender(bgp.Overlay, overlayTargets, bgpCtx(overlay, plan))
		m += r.verifyRender(bgp.Pfc, pfcTargets, bgpCtx(pfc, plan))
		if m > 0 {
			summarize(r.result, dryRun)
			return r.result, fmt.Errorf("render verification failed for %d device(s); not writing", m)
		}
		logger.Get().Info().Msgf("render verification passed (underlay %d, overlay %d, pfc %d)",
			len(plan.devices), len(overlayTargets), len(pfcTargets))
	}

	uid, uok := r.ensureTemplate(bgp.Underlay, underlayDescription, nil)
	r.ensureProfiles(uid, uok, plan.devices, variables, bgp.Underlay.Priority)

	if bgp.Mode == "l3evpn" {
		oid, ook := r.ensureTemplate(bgp.Overlay, overlayDescription, nil)
		r.ensureProfiles(oid, ook, overlayTargets, overlay, bgp.Overlay.Priority)
		pid, pok := r.ensureTemplate(bgp.Pfc, pfcDescription, nil)
		r.ensureProfiles(pid, pok, pfcTargets, pfc, bgp.Pfc.Priority)
		// Action-bound route-domain VRF template: registered with annotations, no profile.
		r.ensureTemplate(bgp.Vrf, vrfDescription, vrfTemplateAnnotations)
		r.reconcileCustomVariables(plan.devices, variables, overlay)
	} else {
		logger.Get().Info().Msgf("overlay + pfc skipped (mode=%s is underlay-only)", bgp.Mode)
	}

	summarize(r.result, dryRun)
	if r.result.Failures > 0 {
		return r.result, fmt.Errorf("bgp registration completed with %d failure(s)", r.result.Failures)
	}
	return r.result, nil
}

func bgpCtx(variables map[int64]map[string]interface{}, plan *planContext) func(*deviceRecord) map[string]interface{} {
	return func(rec *deviceRecord) map[string]interface{} {
		return renderContextBgp(&rec.Device, variables[rec.Id], plan.records)
	}
}

// reconcileCustomVariables writes each device's customVariables = {aggregates,
// is_evpn_rr} (strings), which the engine's route-domain tenant-VRF render reads.
func (r *runner) reconcileCustomVariables(devices []*deviceRecord, variables, overlay map[int64]map[string]interface{}) {
	for _, dev := range devices {
		var aggregates []string
		if a, ok := variables[dev.Id]["aggregates"].([]string); ok {
			aggregates = a
		}
		isRR := false
		if v, ok := overlay[dev.Id]["is_evpn_rr"].(bool); ok {
			isRR = v
		}
		desired := map[string]interface{}{
			"aggregates": strings.Join(aggregates, ","),
			"is_evpn_rr": boolStr(isRR),
		}

		current, driftStatus, revision, err := r.client.GetDeviceCustomVariables(dev.Id)
		if err != nil {
			r.fail("[%s] device GET failed: %s", dev.Label(), err.Error())
			continue
		}
		merged := map[string]interface{}{}
		for k, v := range current {
			merged[k] = v
		}
		for k, v := range desired {
			merged[k] = v
		}
		if canonicalJSON(merged) == canonicalJSON(current) {
			r.count("device custom variables unchanged")
			continue
		}
		if r.dryRun {
			r.count("device custom variables set")
			continue
		}
		if err := r.client.UpdateDeviceCustomVariables(dev.Id, merged, driftStatus, revision); err != nil {
			r.fail("[%s] customVariables PATCH failed: %s", dev.Label(), err.Error())
			continue
		}
		r.count("device custom variables set")
	}
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
