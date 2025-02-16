package system

type ContextKey string

const (
	ApiClientContextKey ContextKey = "apiClient"
)

const (
	ConfigPrefix   = "metalcloud"
	ConfigName     = "metalcloud"
	ConfigType     = "yaml"
	ConfigPath1    = "/etc/metalcloud/"
	ConfigPath2    = "$HOME/.metalcloud/"
	ConfigPath3    = "."
	ConfigEndpoint = "endpoint"
	ConfigApiKey   = "api_key"
	ConfigDebug    = "debug"
)
