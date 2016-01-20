// Package ipamapi specifies the contract the IPAM service (built-in or remote) needs to satisfy.
package ipamapi

import (
	"net"

	"github.com/noironetworks/cilium-net/Godeps/_workspace/src/github.com/docker/libnetwork/types"
)

/********************
 * IPAM plugin types
 ********************/

const (
	// DefaultIPAM is the name of the built-in default ipam driver
	DefaultIPAM = "default"
	// PluginEndpointType represents the Endpoint Type used by Plugin system
	PluginEndpointType = "IpamDriver"
	// RequestAddressType represents the Address Type used when requesting an address
	RequestAddressType = "RequestAddressType"
)

// Callback provides a Callback interface for registering an IPAM instance into LibNetwork
type Callback interface {
	// RegisterIpamDriver provides a way for Remote drivers to dynamically register with libnetwork
	RegisterIpamDriver(name string, driver Ipam) error
	// RegisterIpamDriverWithCapabilities provides a way for Remote drivers to dynamically register with libnetwork and specify cpaabilities
	RegisterIpamDriverWithCapabilities(name string, driver Ipam, capability *Capability) error
}

/**************
 * IPAM Errors
 **************/

// Weel-known errors returned by IPAM
var (
	ErrIpamInternalError   = types.InternalErrorf("IPAM Internal Error")
	ErrInvalidAddressSpace = types.BadRequestErrorf("Invalid Address Space")
	ErrInvalidPool         = types.BadRequestErrorf("Invalid Address Pool")
	ErrInvalidSubPool      = types.BadRequestErrorf("Invalid Address SubPool")
	ErrInvalidRequest      = types.BadRequestErrorf("Invalid Request")
	ErrPoolNotFound        = types.BadRequestErrorf("Address Pool not found")
	ErrOverlapPool         = types.ForbiddenErrorf("Address pool overlaps with existing pool on this address space")
	ErrNoAvailablePool     = types.NoServiceErrorf("No available pool")
	ErrNoAvailableIPs      = types.NoServiceErrorf("No available addresses on this pool")
	ErrIPAlreadyAllocated  = types.ForbiddenErrorf("Address already in use")
	ErrIPOutOfRange        = types.BadRequestErrorf("Requested address is out of range")
	ErrPoolOverlap         = types.ForbiddenErrorf("Pool overlaps with other one on this address space")
	ErrBadPool             = types.BadRequestErrorf("Address space does not contain specified address pool")
	ErrNoIPReturned        = types.NoServiceErrorf("No address returned")
)

/*******************************
 * IPAM Service Interface
 *******************************/

// Ipam represents the interface the IPAM service plugins must implement
// in order to allow injection/modification of IPAM database.
type Ipam interface {
	// GetDefaultAddressSpaces returns the default local and global address spaces for this ipam
	GetDefaultAddressSpaces() (string, string, error)
	// RequestPool returns an address pool along with its unique id. Address space is a mandatory field
	// which denotes a set of non-overlapping pools. pool describes the pool of addresses in CIDR notation.
	// subpool indicates a smaller range of addresses from the pool, for now it is specified in CIDR notation.
	// Both pool and subpool are non mandatory fields. When they are not specified, Ipam driver may choose to
	// return a self chosen pool for this request. In such case the v6 flag needs to be set appropriately so
	// that the driver would return the expected ip version pool.
	RequestPool(addressSpace, pool, subPool string, options map[string]string, v6 bool) (string, *net.IPNet, map[string]string, error)
	// ReleasePool releases the address pool identified by the passed id
	ReleasePool(poolID string) error
	// Request address from the specified pool ID. Input options or preferred IP can be passed.
	RequestAddress(string, net.IP, map[string]string) (*net.IPNet, map[string]string, error)
	// Release the address from the specified pool ID
	ReleaseAddress(string, net.IP) error
}

// Capability represents the requirements and capabilities of the IPAM driver
type Capability struct {
	RequiresMACAddress bool
	SupportsAutoIPv6   bool
}