package lastfm

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/shkh/lastfm-go/lastfm"

	"github.com/anlakii/wallify/config"
	"github.com/anlakii/wallify/providers"
)

type LastfmClient struct {
	api      *lastfm.Api
	config   *config.Config
	lastImgURL string
}

func New(config *config.Config) providers.Provider {
	api := lastfm.New(config.Lastfm.APIKey, "")
	return &LastfmClient{
		api:    api,
		config: config,
	}
}

func (lc *LastfmClient) Update() (bool, error) {
	args := map[string]interface{}{
		"user":  lc.config.Lastfm.Username,
		"limit": 1,
	}
	recentTracks, err := lc.api.User.GetRecentTracks(args)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch recent tracks from Last.fm")
		return false, err
	}

	if len(recentTracks.Tracks) == 0 {
		log.Info().Msg("No recent tracks found for user on Last.fm.")
		return false, nil
	}

	track := recentTracks.Tracks[0]

	if track.Artist.Name == "" || track.Album.Name == "" {
		log.Warn().Msgf("Last.fm track missing artist or album info: %+v", track)
		return false, nil
	}

	imageURL := track.Images[0].Url

	if !strings.Contains(imageURL, "34s") {
		return false, errors.New("failed to get high resolution image")
	}

	imageURL = strings.ReplaceAll(imageURL, "34s", "1000x1000")

	if imageURL == "" {
		log.Warn().Msg("No suitable image found for the track on Last.fm.")
		return false, nil
	}

	if lc.lastImgURL == imageURL {
		return false, nil
	}

	log.Info().Msgf("[NEW ALBUM ART] Downloading from Last.fm: %s", imageURL)

	resp, err := http.Get(imageURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start download from Last.fm")
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("failed to download image from Last.fm, status: %s", resp.Status))
		log.Error().Err(err).Send()
		return false, err
	}

	out, err := os.Create(lc.config.CoverPath)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create cover file: %s", lc.config.CoverPath)
		return false, err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save downloaded image from Last.fm")
		os.Remove(lc.config.CoverPath)
		return false, err
	}

	lc.lastImgURL = imageURL
	return true, nil
}
