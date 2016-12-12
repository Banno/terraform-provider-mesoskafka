package terraform

import (
	"fmt"

	"github.com/hashicorp/terraform/config"
)

// NodeAbstractProvider represents a provider that has no associated operations.
// It registers all the common interfaces across operations for providers.
type NodeAbstractProvider struct {
	NameValue string
	PathValue []string

	// The fields below will be automatically set using the Attach
	// interfaces if you're running those transforms, but also be explicitly
	// set if you already have that information.

	Config *config.ProviderConfig
}

func (n *NodeAbstractProvider) Name() string {
	result := fmt.Sprintf("provider.%s", n.NameValue)
	if len(n.PathValue) > 1 {
		result = fmt.Sprintf("%s.%s", modulePrefixStr(n.PathValue), result)
	}

	return result
}

// GraphNodeSubPath
func (n *NodeAbstractProvider) Path() []string {
	return n.PathValue
}

// GraphNodeReferencer
func (n *NodeAbstractProvider) References() []string {
	if n.Config == nil {
		return nil
	}

	return ReferencesFromConfig(n.Config.RawConfig)
}

// GraphNodeProvider
func (n *NodeAbstractProvider) ProviderName() string {
	return n.NameValue
}

// GraphNodeProvider
func (n *NodeAbstractProvider) ProviderConfig() *config.RawConfig {
	if n.Config == nil {
		return nil
	}

	return n.Config.RawConfig
}

// GraphNodeAttachProvider
func (n *NodeAbstractProvider) AttachProvider(c *config.ProviderConfig) {
	n.Config = c
}
