package nacoswrap

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// ---- Functional options for SelectInstances ----

// SelectInstanceOption configures the SelectInstancesParam for service queries.
type SelectInstanceOption func(*vo.SelectInstancesParam)

// WithHealthyOnly filters to only healthy instances (default true).
func WithHealthyOnly(v bool) SelectInstanceOption {
	return func(p *vo.SelectInstancesParam) { p.HealthyOnly = v }
}

// WithClusters filters instances by cluster names.
func WithClusters(clusters ...string) SelectInstanceOption {
	return func(p *vo.SelectInstancesParam) { p.Clusters = clusters }
}

// WithGroupName sets the group name filter for service queries.
func WithGroupName(group string) SelectInstanceOption {
	return func(p *vo.SelectInstancesParam) { p.GroupName = group }
}

// ---- Service Discovery ----

// GetInstances returns service instances matching the given service name.
// By default it returns only healthy instances; use options to override.
func (c *Client) GetInstances(serviceName string, opts ...SelectInstanceOption) ([]model.Instance, error) {
	if c.namingClient == nil {
		return nil, ErrNamingNotInit
	}

	param := vo.SelectInstancesParam{
		ServiceName: serviceName,
		HealthyOnly: true,
	}
	for _, opt := range opts {
		opt(&param)
	}

	instances, err := c.namingClient.SelectInstances(param)
	if err != nil {
		return nil, fmt.Errorf("nacoswrap: select instances %q: %w", serviceName, err)
	}
	return instances, nil
}

// ---- Service Subscription ----

// Watch subscribes to instance changes for a service. The onChange callback
// is invoked with the updated instance list whenever the server pushes a change.
func (c *Client) Watch(serviceName string, onChange func([]model.Instance)) error {
	if c.namingClient == nil {
		return ErrNamingNotInit
	}

	return c.namingClient.Subscribe(&vo.SubscribeParam{
		ServiceName: serviceName,
		SubscribeCallback: func(instances []model.Instance, err error) {
			if err != nil {
				c.logger.Error("nacoswrap: watch callback error",
					"service", serviceName, "error", err)
				return
			}
			c.logger.Info("nacoswrap: service instances changed",
				"service", serviceName, "count", len(instances))
			onChange(instances)
		},
	})
}

// Unwatch removes a previously registered subscription for the given service.
func (c *Client) Unwatch(serviceName string) error {
	if c.namingClient == nil {
		return ErrNamingNotInit
	}
	return c.namingClient.Unsubscribe(&vo.SubscribeParam{ServiceName: serviceName})
}

// ---- Service Registration ----

// RegisterInstance registers a service instance with Nacos.
func (c *Client) RegisterInstance(param vo.RegisterInstanceParam) (bool, error) {
	if c.namingClient == nil {
		return false, ErrNamingNotInit
	}
	return c.namingClient.RegisterInstance(param)
}

// DeregisterInstance removes a registered service instance from Nacos.
func (c *Client) DeregisterInstance(param vo.DeregisterInstanceParam) (bool, error) {
	if c.namingClient == nil {
		return false, ErrNamingNotInit
	}
	return c.namingClient.DeregisterInstance(param)
}
