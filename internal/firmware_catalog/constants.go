package firmware_catalog

// Firmware catalog vendor, update type and binary update severity values.
//
// These used to be typed enums exported by the SDK (ServerFirmwareCatalogVendor,
// CatalogUpdateType, FirmwareBinaryUpdateSeverity). The SDK now models the
// corresponding fields as plain strings, so the allowed values are defined here.
const (
	VendorDell       = "dell"
	VendorLenovo     = "lenovo"
	VendorHp         = "hp"
	VendorSupermicro = "supermicro"

	UpdateTypeOnline  = "online"
	UpdateTypeOffline = "offline"

	UpdateSeverityCritical    = "critical"
	UpdateSeverityRecommended = "recommended"
	UpdateSeverityOptional    = "optional"
	UpdateSeverityUnknown     = "unknown"
)

// ValidVendors lists the supported firmware catalog vendors.
var ValidVendors = []string{VendorDell, VendorLenovo, VendorHp, VendorSupermicro}

// ValidUpdateTypes lists the supported catalog update types.
var ValidUpdateTypes = []string{UpdateTypeOnline, UpdateTypeOffline}
