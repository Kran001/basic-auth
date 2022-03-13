package config

import (
	"fmt"
	"os"

	"github.com/go-yaml/yaml"
)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func (d *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s",
		"postgres",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Database)
}

type HttpConfig struct {
	Name      string `yaml:"name,attr"`
	HostName  string `yaml:"host"`
	PortHttps string `yaml:"porthttps"`
	CertLoc   string `yaml:"cert-loc"`
	KeyLoc    string `yaml:"key-loc"`
}

func (c *Config) Host() string {
	return c.HttpConfig.HostName
}

func (c *Config) PortHTTPS() string {
	return c.HttpConfig.PortHttps
}

func (c *Config) HTTPSConnectionString() string {
	return fmt.Sprintf("%s:%s", c.HttpConfig.HostName, c.HttpConfig.PortHttps)
}
func (h *HttpConfig) CertLocation() string {
	return h.CertLoc
}

func (h *HttpConfig) KeyLocation() string {
	return h.KeyLoc
}

type Config struct {
	DatabaseConfig *DatabaseConfig `yaml:"database"`
	HttpConfig     *HttpConfig     `yaml:"http"`
}

func NewConfig(pathToConfig string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(pathToConfig)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	if err = yaml.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Database() *DatabaseConfig {
	return c.DatabaseConfig
}

func (c *Config) HTTPSettings() *HttpConfig {
	return c.HttpConfig
}
