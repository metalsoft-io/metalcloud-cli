package fabric_template_config

import (
	"encoding/json"
	"fmt"

	fsc "github.com/metalsoft-io/metalcloud-cli/internal/fabric_switch_config"
	"github.com/metalsoft-io/metalcloud-cli/pkg/logger"
)

const profileLifecycleStage = "configuration"

// Result summarizes a freeform/BGP registration run.
type Result struct {
	Counters map[string]int
	Failures int
	Warnings []string
}

type runner struct {
	client   TemplateClient
	fabricId int64
	dryRun   bool
	apply    string
	result   *Result
}

func (r *runner) count(key string) { r.result.Counters[key]++ }
func (r *runner) fail(format string, args ...any) {
	r.result.Failures++
	logger.Get().Error().Msgf(format, args...)
}

func canonicalJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

// planContext bundles the computed plan + device records shared by the runs.
type planContext struct {
	groups  map[string][]*fsc.Device
	state   *fsc.DesiredState
	records map[int64]*deviceRecord
	devices []*deviceRecord // in (leaf, spine, super_spine) order
}

func buildPlan(client TemplateClient, data []byte, fabricId int64) (*planContext, *fsc.Config, []string, error) {
	planConfig, err := fsc.LoadConfig(data)
	if err != nil {
		return nil, nil, nil, err
	}
	siteId, name, err := client.GetFabric(fabricId)
	if err != nil {
		return nil, nil, nil, err
	}
	logger.Get().Info().Msgf("Target fabric: %q (id=%d)", name, fabricId)

	records, err := client.ListFabricDevices(fabricId)
	if err != nil {
		return nil, nil, nil, err
	}
	if len(records) == 0 {
		return nil, nil, nil, fmt.Errorf("fabric %d has no network devices attached", fabricId)
	}

	var warnings []string
	if needTagHydration(records) && siteId != nil {
		siteDevices, err := client.ListDevicesBySite(*siteId)
		if err != nil {
			return nil, nil, nil, err
		}
		tagsById := map[int64]map[string]string{}
		for _, d := range siteDevices {
			if len(d.TagsMap) > 0 {
				tagsById[d.Id] = d.TagsMap
			}
		}
		n := 0
		for _, d := range records {
			if len(d.TagsMap) == 0 {
				if tags, ok := tagsById[d.Id]; ok {
					d.TagsMap = tags
					n++
				}
			}
		}
		if n > 0 {
			warnings = append(warnings, fmt.Sprintf("backfilled tagsMap for %d device(s) from the siteId=%d listing", n, *siteId))
		}
	}

	engineDevices := make([]*fsc.Device, len(records))
	recByID := map[int64]*deviceRecord{}
	for i, d := range records {
		engineDevices[i] = &d.Device
		recByID[d.Id] = d
	}
	groups, err := fsc.GroupAndOrder(engineDevices, planConfig.Ordering)
	if err != nil {
		return nil, nil, nil, err
	}
	state, err := fsc.ComputeDesired(planConfig, groups)
	if err != nil {
		return nil, nil, nil, err
	}

	ordered := make([]*deviceRecord, 0, len(records))
	for _, position := range switchPositions {
		for _, dev := range groups[position] {
			ordered = append(ordered, recByID[dev.Id])
		}
	}
	return &planContext{groups: groups, state: state, records: recByID, devices: ordered}, planConfig, warnings, nil
}

func needTagHydration(records []*deviceRecord) bool {
	for _, d := range records {
		if len(d.TagsMap) == 0 {
			return true
		}
	}
	return false
}

// ensureTemplate finds-or-creates the template by label; updates only on drift.
// Returns the template id and whether it is usable for profile registration.
func (r *runner) ensureTemplate(spec templateSpec, description string, annotations map[string]string) (int64, bool) {
	contentB64 := base64Encode(spec.Text)
	templates, err := r.client.ListTemplates()
	if err != nil {
		r.fail("listing templates failed: %s", err.Error())
		return 0, false
	}
	var existing *templateRecord
	for _, t := range templates {
		if t.Label == spec.Label {
			existing = t
			break
		}
	}

	if existing == nil {
		if r.dryRun {
			r.count("templates created")
			return 0, false
		}
		id, err := r.client.CreateTemplate(templateCreate{Label: spec.Label, Description: description, ContentB64: contentB64, Annotations: annotations})
		if err != nil {
			r.fail("template %q POST failed: %s", spec.Label, err.Error())
			return 0, false
		}
		r.count("templates created")
		return id, true
	}

	remoteB64 := existing.TemplateB64
	revision := existing.Revision
	if !existing.HasContent {
		content, rev, err := r.client.GetTemplateContent(existing.Id)
		if err != nil {
			r.fail("template %q GET failed: %s", spec.Label, err.Error())
			return existing.Id, true
		}
		remoteB64, revision = content, rev
	}
	contentDrift := base64Decode(remoteB64) != spec.Text
	annDrift := annotations != nil && canonicalJSON(existing.Annotations) != canonicalJSON(annotations)
	if !contentDrift && !annDrift {
		r.count("templates unchanged")
		return existing.Id, true
	}
	if r.dryRun {
		r.count("templates updated")
		return existing.Id, true
	}
	var contentArg *string
	if contentDrift {
		contentArg = &contentB64
	}
	var annArg *map[string]string
	if annDrift {
		annArg = &annotations
	}
	if err := r.client.UpdateTemplate(existing.Id, contentArg, annArg, revision); err != nil {
		r.fail("template %q PATCH failed: %s", spec.Label, err.Error())
		return existing.Id, true
	}
	r.count("templates updated")
	return existing.Id, true
}

// ensureProfiles reconciles one profile per targeted device (POST missing, PATCH
// drifted, leave matches). Devices absent from variables are skipped.
func (r *runner) ensureProfiles(templateId int64, usable bool, devices []*deviceRecord, variables map[int64]map[string]interface{}, priority float32) {
	if !usable && !r.dryRun {
		return
	}
	existingByDevice := map[int64]*profileRecord{}
	if usable {
		profiles, err := r.client.ListProfiles()
		if err != nil {
			r.fail("listing profiles failed: %s", err.Error())
			return
		}
		for _, p := range profiles {
			if p.TemplateId == templateId {
				existingByDevice[p.DeviceId] = p
			}
		}
	}

	for _, dev := range devices {
		vars, ok := variables[dev.Id]
		if !ok {
			continue
		}
		existing := existingByDevice[dev.Id]
		if existing == nil {
			if r.dryRun {
				r.count("profiles created")
				continue
			}
			err := r.client.CreateProfile(profileCreate{
				TemplateId: templateId, DeviceId: dev.Id, FabricId: r.fabricId,
				LifecycleStage: profileLifecycleStage, Variables: vars,
				IsEnabled: true, Priority: priority, ApplyMode: r.apply,
			})
			if err != nil {
				r.fail("[%s] profile POST failed: %s", dev.Label(), err.Error())
				continue
			}
			r.count("profiles created")
			continue
		}

		varsDrift := canonicalJSON(existing.Variables) != canonicalJSON(vars)
		priorityDrift := existing.Priority == nil || *existing.Priority != priority
		applyDrift := existing.ApplyMode != r.apply
		enabledDrift := existing.IsEnabled == nil || !*existing.IsEnabled
		if !varsDrift && !priorityDrift && !applyDrift && !enabledDrift {
			r.count("profiles unchanged")
			continue
		}
		if r.dryRun {
			r.count("profiles updated")
			continue
		}
		id, _ := parseInt64(existing.Id)
		err := r.client.UpdateProfile(id, profileUpdate{Variables: vars, IsEnabled: true, Priority: priority, ApplyMode: r.apply}, existing.Revision)
		if err != nil {
			r.fail("[%s] profile PATCH failed: %s", dev.Label(), err.Error())
			continue
		}
		r.count("profiles updated")
	}
}

// verifyRender pushes each device's render context through the engine's
// stateless render endpoint; a render error counts as a verification failure.
// (Unlike the Python gate, there is no local Jinja2 render to compare against.)
func (r *runner) verifyRender(spec templateSpec, devices []*deviceRecord, contextFor func(*deviceRecord) map[string]interface{}) int {
	contentB64 := base64Encode(spec.Text)
	mismatches := 0
	for _, dev := range devices {
		if _, err := r.client.RenderTemplate(contentB64, contextFor(dev)); err != nil {
			mismatches++
			r.fail("[%s] engine render failed: %s", dev.Label(), err.Error())
		}
	}
	return mismatches
}

func parseInt64(s string) (int64, error) {
	var n int64
	_, err := fmt.Sscan(s, &n)
	return n, err
}

func summarize(result *Result, dryRun bool) {
	for _, w := range result.Warnings {
		logger.Get().Warn().Msg(w)
	}
	parts := ""
	for k, v := range result.Counters {
		if parts != "" {
			parts += ", "
		}
		parts += fmt.Sprintf("%s=%d", k, v)
	}
	if parts == "" {
		parts = "nothing to do"
	}
	suffix := ""
	if dryRun {
		suffix = " (dry-run, no changes made)"
	}
	logger.Get().Info().Msgf("Summary: %s, failures=%d%s", parts, result.Failures, suffix)
}
