package main

import (
	"os"
	"os/exec"
	"time"

	"github.com/anlakii/wallify/config"
	wos "github.com/anlakii/wallify/os"
	"github.com/anlakii/wallify/process"
	"github.com/anlakii/wallify/providers"
	"github.com/anlakii/wallify/providers/lastfm"
	"github.com/anlakii/wallify/providers/spotify"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	wm := wos.WallpaperManager{}

	conf, err := config.Load(wm)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	var client providers.Provider

	log.Info().Msgf("Selected provider: %s", conf.Provider)

	switch conf.Provider {
	case "spotify":
		client = spotify.New(&conf)
		log.Info().Msg("Initialized Spotify provider.")
	case "lastfm":
		client = lastfm.New(&conf)
		log.Info().Msg("Initialized Last.fm provider.")
	default:
		log.Fatal().Msgf("Invalid provider '%s' specified in configuration", conf.Provider)
	}

	processor := process.ImageProcessor{
		Config: &conf,
	}

	log.Info().Msgf("Starting wallpaper update loop (Interval: %dms)", conf.Interval)

	performUpdate(client, &processor, &conf, wm)

	ticker := time.NewTicker(time.Duration(conf.Interval) * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		performUpdate(client, &processor, &conf, wm)
	}
}

func performUpdate(client providers.Provider, processor *process.ImageProcessor, conf *config.Config, wm wos.WallpaperManager) {
	log.Debug().Msg("Checking for updates...")
	updated, err := client.Update()
	if err != nil {
		log.Error().Err(err).Msg("Provider update failed")
		return
	}

	if !updated {
		log.Debug().Msg("No update from provider.")
		return
	}

	log.Info().Msg("Provider reported update, processing image...")
	if err := processor.Process(); err != nil {
		log.Error().Err(err).Msg("Image processing failed")
		return
	}

	log.Info().Msgf("Setting wallpaper to: %s", conf.SavePath)
	if err := wm.SetWallpaper(conf.SavePath); err != nil {
		log.Error().Err(err).Msg("Failed to set wallpaper")
		return
	}

	if conf.Posthook != "" {
		log.Info().Msgf("Running post-hook command: %s %s", conf.Posthook, conf.SavePath)
		cmd := exec.Command(conf.Posthook, conf.SavePath)
		go func() {
			if err := cmd.Run(); err != nil {
				log.Error().Err(err).Msg("Post-hook command failed")
			}
		}()
	}

	log.Info().Msg("Wallpaper update cycle complete.")
}
