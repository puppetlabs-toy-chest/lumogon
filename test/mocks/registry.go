package mocks

import (
	"fmt"

	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/types"
	"github.com/puppetlabs/lumogon/utils"
)

// MockRegistry TODO
type MockRegistry struct {
	AttachedCapabilitiesFn  func() []types.AttachedCapability
	DockerAPICapabilitiesFn func() []dockeradapter.DockerAPICapability
	CountFn                 func() int
	TypesCountFn            func() int
	DescribeCapabilityFn    func(capabilityID string) (string, error)
}

// AttachedCapabilities TODO
func (r MockRegistry) AttachedCapabilities() []types.AttachedCapability {
	if r.AttachedCapabilitiesFn != nil {
		fmt.Println("[MockRegistry] In ", utils.CurrentFunctionName())
		return r.AttachedCapabilitiesFn()
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// DockerAPICapabilities TODO
func (r MockRegistry) DockerAPICapabilities() []dockeradapter.DockerAPICapability {
	if r.DockerAPICapabilitiesFn != nil {
		fmt.Println("[MockRegistry] In ", utils.CurrentFunctionName())
		return r.DockerAPICapabilitiesFn()
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// Count TODO
func (r MockRegistry) Count() int {
	if r.CountFn != nil {
		fmt.Println("[MockRegistry] In ", utils.CurrentFunctionName())
		return r.CountFn()
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// TypesCount TODO
func (r MockRegistry) TypesCount() int {
	if r.TypesCountFn != nil {
		fmt.Println("[MockRegistry] In ", utils.CurrentFunctionName())
		return r.TypesCountFn()
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}

// DescribeCapability TODO
func (r MockRegistry) DescribeCapability(c string) (string, error) {
	if r.DescribeCapabilityFn != nil {
		fmt.Println("[MockRegistry] In ", utils.CurrentFunctionName())
		return r.DescribeCapabilityFn(c)
	}
	panic(fmt.Sprintf("No function defined for: %s", utils.CurrentFunctionName()))
}
