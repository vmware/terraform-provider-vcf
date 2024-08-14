package credentials

const (
	AccountTypeSystem  = "SYSTEM"
	AccountTypeService = "SERVICE"
	AccountTypeUser    = "USER"
)

const (
	ResourceTypeEsxi          = "ESXI"
	ResourceTypeVcenter       = "VCENTER"
	ResourceTypePsc           = "PSC"
	ResourceTypeNsxManager    = "NSX_MANAGER"
	ResourceTypeNsxController = "NSX_CONTROLLER"
	ResourceTypeNsxEdge       = "NSXT_EDGE"
	ResourceTypeNsxtManager   = "NSXT_MANAGER"
	ResourceTypeVrli          = "VRLI"
	ResourceTypeVra           = "VRA"
	ResourceTypeWsa           = "WSA"
	ResourceTypeVrslcm        = "VRSLCM"
	ResourceTypeVxrailManager = "VXRAIL_MANAGER"
	ResourceTypeNsxAlb        = "NSX_ALB"
	ResourceTypeBackup        = "BACKUP"
)

const (
	ConfigAutoRotate = "UPDATE_AUTO_ROTATE_POLICY"
	Rotate           = "ROTATE"
	Update           = "UPDATE"
)

const (
	AutorotateDays30  = 30
	AutorotateDays60  = 60
	AutorotateDays90  = 90
	AutorotateDaysMax = AutorotateDays90
	AutoRotateDaysMin = 1
)

func AllAccountTypes() []string {
	return []string{AccountTypeUser, AccountTypeService, AccountTypeSystem}
}

func AllCredentialTypes() []string {
	return []string{"SSO", "SSH", "API", "FTP", "AUDIT"}
}

func AllResourceTypes() []string {
	return []string{
		ResourceTypeBackup,
		ResourceTypeEsxi,
		ResourceTypeNsxAlb,
		ResourceTypeNsxEdge,
		ResourceTypeNsxController,
		ResourceTypeNsxManager,
		ResourceTypeNsxtManager,
		ResourceTypeVcenter,
		ResourceTypePsc,
		ResourceTypeVrli,
		ResourceTypeVra,
		ResourceTypeWsa,
		ResourceTypeVrslcm,
		ResourceTypeVxrailManager,
	}
}
