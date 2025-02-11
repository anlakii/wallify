package main

import (
	"os/exec"
	"time"

	"github.com/anlakii/wallify/config"
	"github.com/anlakii/wallify/os/darwin"
	"github.com/anlakii/wallify/process"
	"github.com/anlakii/wallify/providers/spotify"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	conf, err := config.Load()
	if err != nil {
		panic(err)
	} 

	client := spotify.New(&conf)

	processor := process.ImageProcessor{
		Config: &conf,
	}

	for {
		time.Sleep(time.Duration(conf.Interval) * time.Millisecond)

		updated, err := client.Update()
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		if !updated {
			continue
		}

		if err := processor.Process(); err != nil {
			log.Fatal().Err(err).Send()
		}

		if conf.Posthook != "" {
			cmd := exec.Command(conf.Posthook, conf.SavePath)
			go cmd.Run()
		}

		darwin.SetWallpaper(conf.SavePath)
	}
}
