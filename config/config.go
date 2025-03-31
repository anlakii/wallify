package config

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/anlakii/wallify/os/darwin"
	"gopkg.in/yaml.v3"
)

type Spotify struct {
	ClientID     string    `yaml:"client_id"`
	ClientSecret string    `yaml:"client_secret"`
	AccessToken  string    `yaml:"access_token"`
	RefreshToken string    `yaml:"refresh_token"`
	TokenType    string    `yaml:"token_type"`
	Expiry       time.Time `yaml:"expiry"`
}

type Config struct {
	SavePath  string `yaml:"save_path"`
	CoverPath string `yaml:"cover_path"`
	Posthook  string `yaml:"post_hook"`
	Height    uint   `yaml:"height"`
	Width     uint   `yaml:"width"`
	Interval  uint   `yaml:"interval"`
	// CacheSize uint   `yaml:"cache_size"`

	Spotify Spotify `yaml:"spotify"`

	configPath string
}

func (c *Config) Save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(c.configPath, data, 0700)
}

func Load() (Config, error) {
	var conf Config

	home, err := os.UserHomeDir()
	if err != nil {
		return conf, err
	}

	conf.configPath = filepath.Join(home, ".config", "wallify", "config.yaml")

	data, err := os.ReadFile(conf.configPath)
	if err != nil {
		return conf, err
	}

	if err := yaml.Unmarshal(data, &conf); err != nil {
		return conf, err
	}

	if conf.Spotify.ClientID == "" {
		return conf, errors.New("client_id empty")
	}

	if conf.Spotify.ClientSecret == "" {
		return conf, errors.New("client_secret empty")
	}

	if conf.SavePath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return conf, err
		}

		conf.SavePath = filepath.Join(home, ".wallify.png")
	}

	if conf.CoverPath == "" {
		conf.CoverPath = filepath.Join(os.TempDir(), "cover.jpg")
	}

	if conf.Interval == 0 {
		conf.Interval = 1000
	}

	if conf.Width == 0 || conf.Height == 0 {
		resolution := darwin.GetResolution()
		conf.Width = resolution.Width
		conf.Height = resolution.Height
	}

	return conf, nil

}
