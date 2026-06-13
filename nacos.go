// Package nacoswrap provides a clean, reusable wrapper around the nacos-sdk-go/v2.
// It simplifies Nacos service discovery (naming client) and configuration management
// (config client) into a single Client with sensible defaults and functional options.
package nacoswrap

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// Config holds all connection parameters needed to connect to a Nacos server.
type Config struct {
	// Host is the Nacos server address (required, e.g. "127.0.0.1").
	Host string

	// Port is the Nacos server port. Defaults to 8848.
	Port uint64

	// Namespace is the Nacos namespace ID. Defaults to "public".
	Namespace string

	// Username for Nacos authentication (optional).
	Username string

	// Password for Nacos authentication (optional).
	Password string

	// AppName is an optional application identifier.
	AppName string

	// LogDir is the directory for nacos SDK logs. When empty, the SDK uses
	// its own default; set to a specific path to control log output.
	LogDir string

	// LogLevel controls SDK log verbosity: "info", "debug", "warn", "error".
	// Defaults to "info".
	LogLevel string

	// TimeoutMs is the request timeout in milliseconds. Defaults to 5000.
	TimeoutMs uint64
}

// defaults applies safe defaults for zero-valued config fields.
func (c *Config) defaults() {
	if c.Port == 0 {
		c.Port = 8848
	}
	if c.Namespace == "" {
		c.Namespace = "public"
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	if c.TimeoutMs == 0 {
		c.TimeoutMs = 5000
	}
}

// Client is a high-level wrapper around the nacos-sdk-go v2.
// It manages both the naming client (service discovery) and the
// config client (configuration management).
type Client struct {
	namingClient naming_client.INamingClient
	configClient config_client.IConfigClient
	logger       *slog.Logger
}

// NewClient creates a Nacos Client from the given Config.
// Both naming and config client are initialized; if the config client
// fails, the Client is still returned with a warning log but config
// operations will return ErrConfigNotInit.
func NewClient(cfg Config) (*Client, error) {
	cfg.defaults()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      cfg.Host,
			Port:        cfg.Port,
			ContextPath: "/nacos",
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         cfg.Namespace,
		TimeoutMs:           cfg.TimeoutMs,
		NotLoadCacheAtStart: true,
		LogLevel:            cfg.LogLevel,
		Username:            cfg.Username,
		Password:            cfg.Password,
		AppName:             cfg.AppName,
	}

	if cfg.LogDir != "" {
		clientConfig.LogDir = cfg.LogDir
		clientConfig.CacheDir = cfg.LogDir + "/cache"
	}

	// Naming client — required
	nc, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		return nil, fmt.Errorf("nacoswrap: naming client: %w", err)
	}

	// Config client — optional; failure is non-fatal
	var cc config_client.IConfigClient
	cc, err = clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		logger.Warn("nacoswrap: config client creation failed (config operations disabled)",
			"error", err)
	}

	logger.Info("nacoswrap: client initialized",
		"host", cfg.Host,
		"port", cfg.Port,
		"namespace", cfg.Namespace,
	)

	return &Client{
		namingClient: nc,
		configClient: cc,
		logger:       logger,
	}, nil
}
