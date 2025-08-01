package config

import (
	"errors"
	stdos "os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/anlakii/wallify/os"
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

type Lastfm struct {
	APIKey   string `yaml:"api_key"`
	Username string `yaml:"username"`
}

type Config struct {
	Provider  string `yaml:"provider"`
	SavePath  string `yaml:"save_path"`
	CoverPath string `yaml:"cover_path"`
	Posthook  string `yaml:"post_hook"`
	Height    uint   `yaml:"height"`
	Width     uint   `yaml:"width"`
	Interval  uint   `yaml:"interval"`

	Spotify Spotify `yaml:"spotify"`
	Lastfm  Lastfm  `yaml:"lastfm"`

	configPath string
}

func (c *Config) Save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	configDir := filepath.Dir(c.configPath)
	if _, err := stdos.Stat(configDir); stdos.IsNotExist(err) {
		if err := stdos.MkdirAll(configDir, 0700); err != nil {
			return err
		}
	}

	return stdos.WriteFile(c.configPath, data, 0600)
}

func Load(wm os.WallpaperManager) (Config, error) {
	var conf Config

	usr, err := user.Current()
	if err != nil {
		return conf, err
	}
	home := usr.HomeDir

	conf.configPath = filepath.Join(home, ".config", "wallify", "config.yaml")

	data, err := stdos.ReadFile(conf.configPath)
	if err != nil {
		if stdos.IsNotExist(err) {
			conf.Provider = "spotify"
			conf.Interval = 1000

			if conf.SavePath == "" {
				conf.SavePath = filepath.Join(home, ".wallify.png")
			}
			if conf.CoverPath == "" {
				conf.CoverPath = filepath.Join(stdos.TempDir(), "cover.jpg")
			}

			resolution, err := wm.Resolution()
			if err != nil {
				return conf, err
			}
			conf.Width = resolution.Width
			conf.Height = resolution.Height

			if err := conf.Save(); err != nil {
				return conf, err
			}
			return conf, errors.New("configuration file created at " + conf.configPath + "; please fill in your API credentials")
		}
		return conf, err
	}

	if err := yaml.Unmarshal(data, &conf); err != nil {
		return conf, err
	}

	conf.Provider = strings.ToLower(conf.Provider)
	switch conf.Provider {
	case "spotify":
		if conf.Spotify.ClientID == "" {
			return conf, errors.New("provider is 'spotify' but spotify.client_id is empty")
		}
		if conf.Spotify.ClientSecret == "" {
			return conf, errors.New("provider is 'spotify' but spotify.client_secret is empty")
		}
	case "lastfm":
		if conf.Lastfm.APIKey == "" {
			return conf, errors.New("provider is 'lastfm' but lastfm.api_key is empty")
		}
		if conf.Lastfm.Username == "" {
			return conf, errors.New("provider is 'lastfm' but lastfm.username is empty")
		}
	case "":
		return conf, errors.New("config error: 'provider' field cannot be empty (must be 'spotify' or 'lastfm')")
	default:
		return conf, errors.New("config error: unknown provider '" + conf.Provider + "' (must be 'spotify' or 'lastfm')")
	}

	if conf.SavePath == "" {
		conf.SavePath = filepath.Join(home, ".wallify.png")
	}

	if conf.CoverPath == "" {
		conf.CoverPath = filepath.Join(stdos.TempDir(), "cover.jpg")
	}

	if conf.Interval == 0 {
		conf.Interval = 1000
	}

	if conf.Width == 0 || conf.Height == 0 {
		resolution, err := wm.Resolution()
		if err != nil {
			// ignore error
		} else {
			conf.Width = resolution.Width
			conf.Height = resolution.Height
		}
	}

	return conf, nil
}
