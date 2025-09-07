// Package config provides the configuration for the application.
package config

import "time"

// Usage is the configuration for the usage collector.
type Usage struct {
	Enabled bool          `yaml:"enabled"`
	Timeout time.Duration `yaml:"timeout"`
}

func (u *Usage) GetTimeout() time.Duration {
	if u.Timeout <= 0 {
		return 10 * time.Second
	}
	return u.Timeout
}

// SystemCollectorConfig is the configuration for the system collector.
type SystemCollectorConfig struct {
	CPUUsage     Usage `yaml:"cpu_usage"`
	MemoryUsage  Usage `yaml:"memory_usage"`
	DiskUsage    Usage `yaml:"disk_usage"`
	NetworkUsage Usage `yaml:"network_usage"`
	ProcessUsage Usage `yaml:"process_usage"`
	ThreadUsage  Usage `yaml:"thread_usage"`
	SocketUsage  Usage `yaml:"socket_usage"`
	FileUsage    Usage `yaml:"file_usage"`
}

type Config struct {
	Server                ServerConfig          `yaml:"server"`
	SystemCollectorConfig SystemCollectorConfig `yaml:"system_collector"`
}

// ServerConfig defines the HTTP server configuration
type ServerConfig struct {
	Address      string        `yaml:"address"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	TLS          TLSConfig     `yaml:"tls"`
}

// TLSConfig defines TLS configuration
type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}
