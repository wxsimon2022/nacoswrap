# nacoswrap

> A clean, reusable Go wrapper around the `nacos-sdk-go/v2`, simplifying service discovery and configuration management into a single `Client`.

[![Go Reference](https://pkg.go.dev/badge/github.com/wxsimon2022/nacoswrap.svg)](https://pkg.go.dev/github.com/wxsimon2022/nacoswrap)
[![Go Report Card](https://goreportcard.com/badge/github.com/wxsimon2022/nacoswrap)](https://goreportcard.com/report/github.com/wxsimon2022/nacoswrap)

## Features

- **Service Discovery** — query, watch, register, and deregister service instances
- **Configuration Management** — get, publish, delete, and listen for config changes
- **Authentication** — supports Nacos username/password auth
- **Functional Options** — clean, composable option pattern for queries
- **Zero Dependencies** beyond the official nacos-sdk-go
- **`slog` Logging** — uses Go's standard `log/slog` logger

## Installation

```bash
go get github.com/wxsimon2022/nacoswrap
```

## Quick Start

### 1. Connect to Nacos

```go
package main

import (
	"log"

	"github.com/wxsimon2022/nacoswrap"
)

func main() {
	client, err := nacoswrap.NewClient(nacoswrap.Config{
		Host:      "127.0.0.1",
		Port:      8848,
		Namespace: "public",
		Username:  "nacos",
		Password:  "nacos",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// ... use client
}
```

### 2. Service Discovery — Query Instances

```go
instances, err := client.GetInstances("my-service")
if err != nil {
	log.Fatal(err)
}
for _, inst := range instances {
	fmt.Printf("  %s:%d (healthy=%v)\n", inst.Ip, inst.Port, inst.Healthy)
}
```

With filtering options:

```go
instances, err := client.GetInstances("my-service",
	nacoswrap.WithClusters("DEFAULT"),
	nacoswrap.WithGroupName("DEFAULT_GROUP"),
	nacoswrap.WithHealthyOnly(true),
)
```

### 3. Service Discovery — Watch Changes

```go
err := client.Watch("my-service", func(instances []model.Instance) {
	fmt.Printf("instances changed, now %d\n", len(instances))
})
```

### 4. Configuration — Get & Listen

```go
// Get a config
val, err := client.GetConfig("app.yml", nacoswrap.WithGroup("APP_GROUP"))

// Listen for changes
err := client.ListenConfig("app.yml", func(newValue string) {
	fmt.Println("config updated:", newValue)
}, nacoswrap.WithGroup("APP_GROUP"))
```

## API Overview

### Naming (Service Discovery)

| Method | Description |
|--------|-------------|
| `GetInstances(serviceName, ...opts)` | Query service instances |
| `Watch(serviceName, onChange)` | Subscribe to instance changes |
| `Unwatch(serviceName)` | Unsubscribe |
| `RegisterInstance(param)` | Register a service instance |
| `DeregisterInstance(param)` | Deregister a service instance |

### Config (Configuration Management)

| Method | Description |
|--------|-------------|
| `GetConfig(dataId, ...opts)` | Get a config value |
| `PublishConfig(dataId, content, ...opts)` | Create or update a config |
| `DeleteConfig(dataId, ...opts)` | Delete a config |
| `ListenConfig(dataId, onChange, ...opts)` | Listen for config changes |
| `CancelListenConfig(dataId, ...opts)` | Stop listening |

## License

MIT
