package interfaces

import (
	metalcloud "github.com/bigstepinc/metal-cloud-sdk-go"
)

//MetalCloudClient interface describing the functions used by this cli in interacting with the metal cloud client, used for testing
type MetalCloudClient interface {
	//InfrastructureCreate creates an infrastructure
	InfrastructureCreate(infrastructure metalcloud.Infrastructure) (*metalcloud.Infrastructure, error)
	//InfrastructureEdit alters an infrastructure
	InfrastructureEdit(infrastructureID metalcloud.ID, infrastructureOperation metalcloud.InfrastructureOperation) (*metalcloud.Infrastructure, error)
	//InfrastructureDelete deletes an infrastructure and all associated elements. Requires deploy
	InfrastructureDelete(infrastructureID metalcloud.ID) error
	//InfrastructureOperationCancel reverts (undos) alterations done before deploy
	InfrastructureOperationCancel(infrastructureID metalcloud.ID) error
	//InfrastructureDeploy initiates a deploy operation that will apply all registered changes for the respective infrastructure
	InfrastructureDeploy(infrastructureID metalcloud.ID, shutdownOptions metalcloud.ShutdownOptions, allowDataLoss bool, skipAnsible bool) error

	//Infrastructures returns a list of infrastructures
	Infrastructures() (*map[string]metalcloud.Infrastructure, error)
	//InfrastructureGet returns a specific infrastructure
	InfrastructureGet(infrastructureID metalcloud.ID) (*metalcloud.Infrastructure, error)
	//InfrastructureUserLimits returns user metadata
	InfrastructureUserLimits(infrastructureID metalcloud.ID) (*map[string]interface{}, error)

	//InstanceArrayGet returns an InstanceArray with specified id
	InstanceArrayGet(instanceArrayID metalcloud.ID) (*metalcloud.InstanceArray, error)
	//InstanceArrays returns list of instance arrays of specified infrastructure
	InstanceArrays(infrastructureID metalcloud.ID) (*map[string]metalcloud.InstanceArray, error)
	//InstanceArrayCreate creates an instance array (colletion of identical instances). Requires Deploy.
	InstanceArrayCreate(infrastructureID metalcloud.ID, instanceArray metalcloud.InstanceArray) (*metalcloud.InstanceArray, error)
	//InstanceArrayEdit alterns a deployed instance array. Requires deploy.
	InstanceArrayEdit(instanceArrayID metalcloud.ID, instanceArrayOperation metalcloud.InstanceArrayOperation, bSwapExistingInstancesHardware *bool, bKeepDetachingDrives *bool, objServerTypeMatches *[]metalcloud.ServerType, arrInstancesToBeDeleted *[]int) (*metalcloud.InstanceArray, error)
	//InstanceArrayDelete deletes an instance array. Requires deploy.
	InstanceArrayDelete(instanceArrayID metalcloud.ID) error
	//InstanceArrayInterfaceAttachNetwork attaches an InstanceArrayInterface to a Network
	InstanceArrayInterfaceAttachNetwork(instanceArrayID int, instanceArrayInterfaceIndex int, networkID int) (*metalcloud.InstanceArray, error)
	//InstanceArrayInterfaceDetach detaches an InstanceArrayInterface from any Network element that is attached to.
	InstanceArrayInterfaceDetach(instanceArrayID int, instanceArrayInterfaceIndex int) (*metalcloud.InstanceArray, error)

	//NetworkGet retrieves a network object
	NetworkGet(networkID metalcloud.ID) (*metalcloud.Network, error)
	//Networks returns a list of all network objects of an infrastructure
	Networks(infrastructureID metalcloud.ID) (*map[string]metalcloud.Network, error)
	//NetworkCreate creates a network
	NetworkCreate(infrastructureID metalcloud.ID, network metalcloud.Network) (*metalcloud.Network, error)
	//NetworkEdit applies a change to an existing network
	NetworkEdit(networkID metalcloud.ID, networkOperation metalcloud.NetworkOperation) (*metalcloud.Network, error)
	//NetworkDelete deletes a network.
	NetworkDelete(networkID metalcloud.ID) error
	//NetworkJoin merges two specified Network objects.
	NetworkJoin(networkID metalcloud.ID, networkToBeDeletedID metalcloud.ID) error

	//ServerTypeGet returns the respective server type
	ServerTypeGet(serverTypeID metalcloud.ID) (*metalcloud.ServerType, error)
	//ServerTypesMatches matches available servers with a certain Instance's configuration, using the properties specified in the objHardwareConfiguration object, and returns the number of compatible servers for each server_type_id.
	ServerTypesMatches(infrastructureID metalcloud.ID, hardwareConfiguration metalcloud.HardwareConfiguration, instanceArrayID *int, bAllowServerSwap bool) (*map[string]metalcloud.ServerType, error)
	//ServerTypesMatchHardwareConfiguration Retrieves a list of server types that match the provided hardware configuration. The function does not check for availability, only compatibility, so physical servers associated with the returned server types might be unavailable.
	ServerTypesMatchHardwareConfiguration(datacenterName string, hardwareConfiguration metalcloud.HardwareConfiguration) (*map[int]metalcloud.ServerType, error)
	//ServerTypeDatacenter retrieves all the server type IDs for servers found in a specified Datacenter
	ServerTypeDatacenter(datacenterName string) (*[]int, error)
	//ServerTypes retrieves all ServerType objects from the database.
	ServerTypes(datacenterName string, bOnlyAvailable bool) (*map[int]metalcloud.ServerType, error)

	//VolumeTemplates retrives the list of available templates
	VolumeTemplates() (*map[string]metalcloud.VolumeTemplate, error)
	//VolumeTemplateGet returns the specified volume template
	VolumeTemplateGet(volumeTemplateID metalcloud.ID) (*metalcloud.VolumeTemplate, error)

	//DriveArrays retrieves the list of drives arrays of an infrastructure
	DriveArrays(infrastructureID metalcloud.ID) (*map[string]metalcloud.DriveArray, error)
	//DriveArrayGet retrieves a DriveArray object with specified ids
	DriveArrayGet(driveArrayID metalcloud.ID) (*metalcloud.DriveArray, error)
	//DriveArrayCreate creates a drive array. Requires deploy.
	DriveArrayCreate(infrastructureID metalcloud.ID, driveArray metalcloud.DriveArray) (*metalcloud.DriveArray, error)
	//DriveArrayEdit alters a deployed drive array. Requires deploy.
	DriveArrayEdit(driveArrayID metalcloud.ID, driveArrayOperation metalcloud.DriveArrayOperation) (*metalcloud.DriveArray, error)
	//DriveArrayDelete deletes a Drive Array with specified id
	DriveArrayDelete(driveArrayID metalcloud.ID) error
}
