package fabric_template_config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	fsc "github.com/metalsoft-io/metalcloud-cli/internal/fabric_switch_config"
)

// fakeTemplateClient is an in-memory TemplateClient that records and applies
// writes so a re-run observes the new state.
type fakeTemplateClient struct {
	siteId    int64
	devices   []*deviceRecord
	nextId    int64
	templates map[int64]*templateRecord
	profiles  map[string]*profileRecord // key: templateId|deviceId

	templatesCreated int
	templatesUpdated int
	profilesCreated  int
	profilesUpdated  int
	customVarsSet    int
	renders          int
}

func (f *fakeTemplateClient) newId() int64 { f.nextId++; return f.nextId }

func (f *fakeTemplateClient) GetFabric(int64) (*int64, string, error) {
	s := f.siteId
	return &s, "Test Fabric", nil
}
func (f *fakeTemplateClient) ListFabricDevices(int64) ([]*deviceRecord, error) { return f.devices, nil }
func (f *fakeTemplateClient) ListDevicesBySite(int64) ([]*deviceRecord, error) { return f.devices, nil }

func (f *fakeTemplateClient) ListTemplates() ([]*templateRecord, error) {
	out := []*templateRecord{}
	for _, t := range f.templates {
		out = append(out, t)
	}
	return out, nil
}
func (f *fakeTemplateClient) GetTemplateContent(id int64) (string, string, error) {
	return f.templates[id].TemplateB64, f.templates[id].Revision, nil
}
func (f *fakeTemplateClient) CreateTemplate(t templateCreate) (int64, error) {
	f.templatesCreated++
	id := f.newId()
	f.templates[id] = &templateRecord{Id: id, Label: t.Label, TemplateB64: t.ContentB64, HasContent: true, Annotations: t.Annotations, Revision: "1"}
	return id, nil
}
func (f *fakeTemplateClient) UpdateTemplate(id int64, contentB64 *string, annotations *map[string]string, revision string) error {
	f.templatesUpdated++
	if contentB64 != nil {
		f.templates[id].TemplateB64 = *contentB64
	}
	if annotations != nil {
		f.templates[id].Annotations = *annotations
	}
	return nil
}

func (f *fakeTemplateClient) ListProfiles() ([]*profileRecord, error) {
	out := []*profileRecord{}
	for _, p := range f.profiles {
		out = append(out, p)
	}
	return out, nil
}
func (f *fakeTemplateClient) CreateProfile(p profileCreate) error {
	f.profilesCreated++
	id := f.newId()
	prio := p.Priority
	enabled := p.IsEnabled
	f.profiles[fmt.Sprintf("%d|%d", p.TemplateId, p.DeviceId)] = &profileRecord{
		Id: fmt.Sprintf("%d", id), TemplateId: p.TemplateId, DeviceId: p.DeviceId,
		Variables: p.Variables, Priority: &prio, ApplyMode: p.ApplyMode, IsEnabled: &enabled, Revision: "1",
	}
	return nil
}
func (f *fakeTemplateClient) UpdateProfile(id int64, p profileUpdate, revision string) error {
	f.profilesUpdated++
	return nil
}
func (f *fakeTemplateClient) RenderTemplate(string, map[string]interface{}) (string, error) {
	f.renders++
	return "rendered", nil
}
func (f *fakeTemplateClient) GetDeviceCustomVariables(int64) (map[string]interface{}, string, error) {
	return map[string]interface{}{}, "1", nil
}
func (f *fakeTemplateClient) UpdateDeviceCustomVariables(int64, map[string]interface{}, string) error {
	f.customVarsSet++
	return nil
}

func fixtureDeviceRecords() []*deviceRecord {
	dev := func(id int64, position, mgmt string, tags map[string]string) *deviceRecord {
		return &deviceRecord{Device: fsc.Device{Id: id, Position: position, ManagementAddress: mgmt,
			IdentifierString: "old-" + position, Driver: "cumulus_linux", TagsMap: tags}}
	}
	return []*deviceRecord{
		dev(1, "spine", "10.0.0.22", map[string]string{"nvidia/pod-id": "5", "nvidia/rail-group-id": "1", "nvidia/spine-index": "2"}),
		dev(2, "spine", "10.0.0.21", map[string]string{"nvidia/pod-id": "5", "nvidia/rail-group-id": "1", "nvidia/spine-index": "1"}),
		dev(3, "leaf", "10.0.0.12", map[string]string{"nvidia/pod-id": "5", "nvidia/scalability-unit-id": "2", "nvidia/rail-group-id": "1"}),
		dev(4, "leaf", "10.0.0.11", map[string]string{"nvidia/pod-id": "5", "nvidia/scalability-unit-id": "1", "nvidia/rail-group-id": "1"}),
		dev(5, "leaf", "10.0.0.13", map[string]string{"nvidia/pod-id": "5", "nvidia/scalability-unit-id": "3", "nvidia/rail-group-id": "1"}),
		dev(6, "super_spine", "10.0.0.31", map[string]string{"nvidia/ssp-group-id": "1"}),
		dev(7, "super_spine", "10.0.0.32", map[string]string{"nvidia/ssp-group-id": "2"}),
		dev(8, "super_spine", "10.0.0.33", map[string]string{"nvidia/ssp-group-id": "1"}),
	}
}

func writeTemplate(t *testing.T, dir, name, body string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	return path
}

func newFakeClient() *fakeTemplateClient {
	return &fakeTemplateClient{siteId: 11, devices: fixtureDeviceRecords(), nextId: 1000,
		templates: map[int64]*templateRecord{}, profiles: map[string]*profileRecord{}}
}

func TestRunFreeformWriteIdempotencyDryRun(t *testing.T) {
	dir := t.TempDir()
	tmpl := writeTemplate(t, dir, "freeform.j2", "hostname {{ identifierString }}\nmode {{ mode }}\n")
	config := []byte(fmt.Sprintf(`
ordering: managementAddress
loopback:
  subnet: 10.253.128.0/18
topology:
  leafSpine:
    linksPerPair: auto
  spineSuperSpine:
    linksPerPair: 4
  leafHost:
    nodeCount: 2
p2p:
  mtu: 9216
freeform:
  mode: l3evpn
  templatePath: %s
`, tmpl))

	f := newFakeClient()
	res, err := RunFreeform(f, config, 5, false, false)
	if err != nil {
		t.Fatalf("RunFreeform: %v", err)
	}
	if res.Failures != 0 {
		t.Fatalf("failures = %d", res.Failures)
	}
	if f.templatesCreated != 1 {
		t.Errorf("templates created = %d, want 1", f.templatesCreated)
	}
	if f.profilesCreated != 8 {
		t.Errorf("profiles created = %d, want 8", f.profilesCreated)
	}

	// Idempotency: second run sees the created template + profiles -> no writes.
	f.templatesCreated, f.profilesCreated, f.templatesUpdated, f.profilesUpdated = 0, 0, 0, 0
	if _, err := RunFreeform(f, config, 5, false, false); err != nil {
		t.Fatalf("second RunFreeform: %v", err)
	}
	if f.templatesCreated != 0 || f.profilesCreated != 0 || f.templatesUpdated != 0 || f.profilesUpdated != 0 {
		t.Errorf("second run made writes: tCreate=%d pCreate=%d tUpd=%d pUpd=%d",
			f.templatesCreated, f.profilesCreated, f.templatesUpdated, f.profilesUpdated)
	}

	// Dry-run on a fresh fake -> no creates recorded.
	fdry := newFakeClient()
	if _, err := RunFreeform(fdry, config, 5, true, false); err != nil {
		t.Fatalf("dry-run RunFreeform: %v", err)
	}
	if fdry.templatesCreated != 0 || fdry.profilesCreated != 0 {
		t.Errorf("dry-run made writes")
	}
}

func TestRunBgpL3evpn(t *testing.T) {
	dir := t.TempDir()
	body := "{{ mode }} {{ position }}\n"
	config := []byte(fmt.Sprintf(`
ordering: managementAddress
loopback:
  subnet: 10.253.128.0/18
asn: {}
topology:
  leafSpine:
    linksPerPair: auto
  spineSuperSpine:
    linksPerPair: 4
  leafHost:
    nodeCount: 2
p2p:
  mtu: 9216
bgp:
  mode: l3evpn
  templatePath: %s
  overlayTemplatePath: %s
  pfcTemplatePath: %s
  vrfTemplatePath: %s
`,
		writeTemplate(t, dir, "underlay.j2", body),
		writeTemplate(t, dir, "overlay.j2", body),
		writeTemplate(t, dir, "pfc.j2", body),
		writeTemplate(t, dir, "vrf.j2", body)))

	f := newFakeClient()
	// asn must be present on records (configure-switches first); set it.
	for _, d := range f.devices {
		asn := int64(4200000000 + d.Id)
		d.Asn = &asn
	}
	res, err := RunBgp(f, config, 5, false, false)
	if err != nil {
		t.Fatalf("RunBgp: %v", err)
	}
	if res.Failures != 0 {
		t.Fatalf("failures = %d", res.Failures)
	}
	// l3evpn registers 4 templates (underlay, overlay, pfc, vrf).
	if f.templatesCreated != 4 {
		t.Errorf("templates created = %d, want 4", f.templatesCreated)
	}
	// underlay: a profile per switch (8); overlay: leaves(3)+RRs(2)=5; pfc: 8.
	if f.profilesCreated != 8+5+8 {
		t.Errorf("profiles created = %d, want %d", f.profilesCreated, 8+5+8)
	}
	// customVariables set on all 8 switches (l3evpn).
	if f.customVarsSet != 8 {
		t.Errorf("customVariables set = %d, want 8", f.customVarsSet)
	}
}

func TestRunBgpRequiresAsnLoopback(t *testing.T) {
	dir := t.TempDir()
	body := "x\n"
	config := []byte(fmt.Sprintf(`
loopback:
  subnet: 10.253.128.0/18
topology:
  leafSpine:
    linksPerPair: auto
  spineSuperSpine:
    linksPerPair: 4
  leafHost:
    nodeCount: 2
p2p:
  mtu: 9216
bgp:
  mode: purel3
  templatePath: %s
  overlayTemplatePath: %s
  pfcTemplatePath: %s
  vrfTemplatePath: %s
`,
		writeTemplate(t, dir, "u.j2", body), writeTemplate(t, dir, "o.j2", body),
		writeTemplate(t, dir, "p.j2", body), writeTemplate(t, dir, "v.j2", body)))

	f := newFakeClient() // records have no asn -> should error
	if _, err := RunBgp(f, config, 5, false, false); err == nil {
		t.Error("expected error when devices lack asn/loopback")
	}
}
