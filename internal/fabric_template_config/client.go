package fabric_template_config

// TemplateClient is the narrow set of operations the freeform/BGP runner needs,
// abstracted from the SDK so the runner can be unit-tested with a fake.
type TemplateClient interface {
	GetFabric(fabricId int64) (siteId *int64, name string, err error)
	ListFabricDevices(fabricId int64) ([]*deviceRecord, error)
	ListDevicesBySite(siteId int64) ([]*deviceRecord, error)

	ListTemplates() ([]*templateRecord, error)
	GetTemplateContent(id int64) (contentB64 string, revision string, err error)
	CreateTemplate(t templateCreate) (int64, error)
	UpdateTemplate(id int64, contentB64 *string, annotations *map[string]string, revision string) error

	ListProfiles() ([]*profileRecord, error)
	CreateProfile(p profileCreate) error
	UpdateProfile(id int64, p profileUpdate, revision string) error

	RenderTemplate(contentB64 string, variables map[string]interface{}) (rendered string, err error)

	GetDeviceCustomVariables(deviceId int64) (current map[string]interface{}, driftStatus string, revision string, err error)
	UpdateDeviceCustomVariables(deviceId int64, customVariables map[string]interface{}, driftStatus string, revision string) error
}

type templateRecord struct {
	Id          int64
	Label       string
	TemplateB64 string // may be empty if the list omits it
	HasContent  bool
	Annotations map[string]string
	Revision    string
}

type templateCreate struct {
	Label       string
	Description string
	ContentB64  string
	Annotations map[string]string // nil for profile-bound templates
}

type profileRecord struct {
	Id         string
	TemplateId int64
	DeviceId   int64
	Variables  map[string]interface{}
	Priority   *float32
	ApplyMode  string
	IsEnabled  *bool
	Revision   string
}

type profileCreate struct {
	TemplateId     int64
	DeviceId       int64
	FabricId       int64
	LifecycleStage string
	Variables      map[string]interface{}
	IsEnabled      bool
	Priority       float32
	ApplyMode      string
}

type profileUpdate struct {
	Variables map[string]interface{}
	IsEnabled bool
	Priority  float32
	ApplyMode string
}
