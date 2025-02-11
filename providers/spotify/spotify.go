package spotify

import (
	"context"
	"net/http"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"github.com/anlakii/wallify/config"
	"github.com/anlakii/wallify/providers"
	"github.com/rs/zerolog/log"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	_ "github.com/reujab/wallpaper"
	"golang.org/x/oauth2"
)

var tokenChan chan *oauth2.Token
var auth *spotifyauth.Authenticator

type SpotifyClient struct {
	lastImg       string
	spotifyClient *spotify.Client

	config      *config.Config
	redirectURL string
	csrfToken   string
}

func (sc *SpotifyClient) Authenticate() error {
	tokenChan = make(chan *oauth2.Token)

	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(sc.redirectURL),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadCurrentlyPlaying,
			spotifyauth.ScopeUserReadPlaybackState,
		),
		spotifyauth.WithClientID(sc.config.Spotify.ClientID),
		spotifyauth.WithClientSecret(sc.config.Spotify.ClientSecret),
	)

	url := auth.AuthURL(sc.csrfToken)

	err := exec.Command("open", url).Run()

	if err != nil {
		return err
	}

	http.HandleFunc("/callback", sc.callback)

	go http.ListenAndServe(":8080", nil)

	token := <-tokenChan

	sc.spotifyClient = spotify.New(auth.Client(context.TODO(), token))

	sc.config.Spotify.AccessToken = token.AccessToken
	sc.config.Spotify.RefreshToken = token.RefreshToken
	sc.config.Spotify.TokenType = token.TokenType
	sc.config.Spotify.Expiry = token.Expiry

	if err := sc.config.Save(); err != nil {
		return err
	}

	return nil
}

func New(config *config.Config) providers.Provider {
	var client SpotifyClient
	tokenChan = make(chan *oauth2.Token)

	client.redirectURL = "http://localhost:8080/callback"
	client.csrfToken = uuid.New().String()

	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(client.redirectURL),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadCurrentlyPlaying,
			spotifyauth.ScopeUserReadPlaybackState,
		),
		spotifyauth.WithClientID(config.Spotify.ClientID),
		spotifyauth.WithClientSecret(config.Spotify.ClientSecret),
	)

	if config.Spotify.AccessToken != "" {
		// var conf = oauth2.Config{
		// 	ClientID:     config.Spotify.ClientID,
		// 	ClientSecret: config.Spotify.ClientSecret,
		// 	Scopes:       []string{spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopeUserReadPlaybackState},
		// 	RedirectURL:  redirectURL,
		// 	Endpoint: oauth2.Endpoint{
		// 		AuthURL:  spotifyauth.AuthURL,
		// 		TokenURL: spotifyauth.TokenURL,
		// 	},
		// }

		token := &oauth2.Token{

			AccessToken:  config.Spotify.AccessToken,
			RefreshToken: config.Spotify.RefreshToken,
			Expiry:       config.Spotify.Expiry,
			TokenType:    config.Spotify.TokenType,
		}

		client.spotifyClient = spotify.New(auth.Client(context.TODO(), token))
		client.config = config

		return &client
	}

	if err := client.Authenticate(); err != nil {
		panic(err)
	}

	return &client
}

func (sc *SpotifyClient) callback(w http.ResponseWriter, r *http.Request) {
	token, err := auth.Token(r.Context(), sc.csrfToken, r)
	if err != nil {
		panic(err)
	}

	tokenChan <- token
}

func (sc *SpotifyClient) Update() (bool, error) {
	playerState, err := sc.spotifyClient.PlayerState(context.Background())
	if err != nil {
		return false, sc.Authenticate()
	}
	if playerState == nil {
		return false, nil
	}

	currentlyPlaying := playerState.Item
	if currentlyPlaying == nil {
		return false, nil
	}

	firstImg := currentlyPlaying.Album.Images[0]

	if sc.lastImg == firstImg.URL {
		return false, nil
	}

	log.Info().Msgf("[NEW ALBUM] %s", currentlyPlaying.Album.Name)

	out, err := os.Create(sc.config.CoverPath)
	if err != nil {
		return false, err
	}

	if err := firstImg.Download(out); err != nil {
		return false, err
	}

	out.Close()

	sc.lastImg = firstImg.URL

	return true, nil
}
