package services

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/antonholmquist/jason"
	"github.com/layeh/gumble/gumble"
	// probably need to change that
	"github.com/matthieugrieger/mumbledj/bot"
	"github.com/matthieugrieger/mumbledj/interfaces"
	"github.com/spf13/viper"
)

// LocalStorage is music player from local storage. (WOW)
type LocalStorage struct {
	*GenericService
}

// NewLocalStorageService might or might not be needed
func NewLocalStorageService() *LocalStorage {
	return &LocalStorage{
		&GenericService{
			ReadableName: "LocalStorage",
			Format:       "opus", // ???
			TrackRegex: []*regexp.Regexp{
				//implement here checking file?
				regexp.MustCompile(`https?:\/\/www.youtube.com\/watch\?v=(?P<id>[\w-]+)(?P<timestamp>\&t=\d*m?\d*s?)?`),
				regexp.MustCompile(`https?:\/\/youtube.com\/watch\?v=(?P<id>[\w-]+)(?P<timestamp>\&t=\d*m?\d*s?)?`),
				regexp.MustCompile(`https?:\/\/youtu.be\/(?P<id>[\w-]+)(?P<timestamp>\?t=\d*m?\d*s?)?`),
				regexp.MustCompile(`https?:\/\/youtube.com\/v\/(?P<id>[\w-]+)(?P<timestamp>\?t=\d*m?\d*s?)?`),
				regexp.MustCompile(`https?:\/\/www.youtube.com\/v\/(?P<id>[\w-]+)(?P<timestamp>\?t=\d*m?\d*s?)?`),
			},
			PlaylistRegex: nil,
		},
	}
}

// CheckAPIKey performs a test API call with the API key
// provided in the configuration file to determine if the
// service should be enabled.
func (ls *LocalStorage) CheckAPIKey() error {
	// No key needed.
	return nil
}

// GetTracks uses the passed URL to find and return
// tracks associated with the URL. An error is returned
// if any error occurs during the API call.
func (ls *LocalStorage) GetTracks(url string, submitter *gumble.User) ([]interfaces.Track, error) {
	var (
		playlistURL      string
		playlistItemsURL string
		id               string
		err              error
		resp             *http.Response
		v                *jason.Object
		track            bot.Track
		tracks           []interfaces.Track
	)

	dummyOffset, _ := time.ParseDuration("0s")
	urlSplit := strings.Split(url, "?t=")

	playlistURL = "https://www.googleapis.com/youtube/v3/playlists?part=snippet&id=%s&key=%s"
	playlistItemsURL = "https://www.googleapis.com/youtube/v3/playlistItems?part=snippet,contentDetails&playlistId=%s&maxResults=%d&key=%s&pageToken=%s"
	id, err = yt.getID(urlSplit[0])
	if err != nil {
		return nil, err
	}

	if yt.isPlaylist(url) {
		resp, err = http.Get(fmt.Sprintf(playlistURL, id, viper.GetString("api_keys.youtube")))
		defer resp.Body.Close()
		if err != nil {
			return nil, err
		}

		v, err = jason.NewObjectFromReader(resp.Body)
		if err != nil {
			return nil, err
		}

		items, _ := v.GetObjectArray("items")
		item := items[0]

		title, _ := item.GetString("snippet", "title")

		playlist := &bot.Playlist{
			ID:        id,
			Title:     title,
			Submitter: submitter.Name,
			Service:   yt.ReadableName,
		}

		maxItems := math.MaxInt32
		if viper.GetInt("queue.max_tracks_per_playlist") > 0 {
			maxItems = viper.GetInt("queue.max_tracks_per_playlist")
		}

		// YouTube playlist searches return a max of 50 results per page
		maxResults := 50
		if maxResults > maxItems {
			maxResults = maxItems
		}

		pageToken := ""
		for len(tracks) < maxItems {
			curResp, curErr := http.Get(fmt.Sprintf(playlistItemsURL, id, maxResults, viper.GetString("api_keys.youtube"), pageToken))
			defer curResp.Body.Close()
			if curErr != nil {
				// An error occurred, simply skip this track.
				continue
			}

			v, err = jason.NewObjectFromReader(curResp.Body)
			if err != nil {
				// An error occurred, simply skip this track.
				continue
			}

			curTracks, _ := v.GetObjectArray("items")
			for _, track := range curTracks {
				videoID, _ := track.GetString("snippet", "resourceId", "videoId")

				// Unfortunately we have to execute another API call for each video as the YouTube API does not
				// return video durations from the playlistItems endpoint...
				newTrack, _ := yt.getTrack(videoID, submitter, dummyOffset)
				newTrack.Playlist = playlist
				tracks = append(tracks, newTrack)

				if len(tracks) >= maxItems {
					break
				}
			}

			pageToken, _ = v.GetString("nextPageToken")
			if pageToken == "" {
				break
			}
		}

		if len(tracks) == 0 {
			return nil, errors.New("Invalid playlist. No tracks were added")
		}
		return tracks, nil
	}

	// Submitter added a track!
	offset := dummyOffset
	if len(urlSplit) == 2 {
		offset, _ = time.ParseDuration(urlSplit[1])
	}

	track, err = yt.getTrack(id, submitter, offset)
	if err != nil {
		return nil, err
	}
	tracks = append(tracks, track)
	return tracks, nil
}

func (yt *YouTube) getTrack(id string, submitter *gumble.User, offset time.Duration) (bot.Track, error) {
	var (
		resp *http.Response
		err  error
		v    *jason.Object
	)

	// do a check for file here probably, not in regrex at top
	videoURL := "https://www.googleapis.com/youtube/v3/videos?part=snippet,contentDetails&id=%s&key=%s"
	resp, err = http.Get(fmt.Sprintf(videoURL, id, viper.GetString("api_keys.youtube")))
	defer resp.Body.Close()
	if err != nil {
		return bot.Track{}, err
	}

	v, err = jason.NewObjectFromReader(resp.Body)
	if err != nil {
		return bot.Track{}, err
	}
	items, _ := v.GetObjectArray("items")
	if len(items) == 0 {
		return bot.Track{}, errors.New("This YouTube video is private")
	}
	item := items[0]
	title, _ := item.GetString("snippet", "title")
	thumbnail, _ := item.GetString("snippet", "thumbnails", "high", "url")
	author, _ := item.GetString("snippet", "channelTitle")
	durationString, _ := item.GetString("contentDetails", "duration")
	durationConverted, _ := duration.FromString(durationString)
	duration := durationConverted.ToDuration()

	return bot.Track{
		ID:             id,
		URL:            "https://youtube.com/watch?v=" + id,
		Title:          title,
		Author:         author,
		Submitter:      submitter.Name,
		Service:        yt.ReadableName,
		Filename:       id + ".track",
		ThumbnailURL:   thumbnail,
		Duration:       duration,
		PlaybackOffset: offset,
		Playlist:       nil,
	}, nil
}
