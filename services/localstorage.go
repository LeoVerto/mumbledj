/*
 * MumbleDJ
 * By Matthieu Grieger
 * services/localstorage.go
 * Copyright (c) 2018 Roshless (MIT License)
 */

package services

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"git.roshless.me/roshless/mumbledj/bot"
	"git.roshless.me/roshless/mumbledj/interfaces"
	"github.com/ChannelMeter/iso8601duration"
	"github.com/antonholmquist/jason"
	"github.com/layeh/gumble/gumble"
	"github.com/spf13/viper"
)

// LocalStorage is simple local file player.
type LocalStorage struct {
	*GenericService
}

// NewLocalStorageService returns an initialized LocalStorage service object
func NewLocalStorageService() *LocalStorage {
	return &LocalStorage{
		&GenericService{
			ReadableName: "LocalStorage",
			Format:       "",
			TrackRegex: []*regexp.Regexp{
				regexp.MustCompile(`+\.ls`),
			},
			PlaylistRegex: nil,
		},
	}
}

// CheckAPIKey performs a test API call with the API key
// provided in the configuration file to determine if the
// service should be enabled.
func (ls *LocalStorage) CheckAPIKey() error {
	// We check for
	if viper.GetBool("localstorage.enabled") == false {
		return errors.New("LocalStorage disabled")
	}
	return nil
}

// GetTracks uses the passed tag to find and return
// tracks associated with the tag (file name without extension).
// An error is returned if any error occurs during the API call.
func (ls *LocalStorage) GetTracks(tag string, submitter *gumble.User) ([]interfaces.Track, error) {
	var (
		id        string
		err       error
		directory string
		filePath  string
		jsonPath  string
		v         *jason.Object
		tracks    []interfaces.Track
	)

	fileExtension := ".track"
	dummyOffset, _ := time.ParseDuration("0s")
	tagSplit := strings.Split(tag, ".")

	// getID needs to have some kind of magic <id> in regrex
	/*
		id, err = ls.getID(tagSplit[0])
		if err != nil {
			return nil, err
		}
	*/
	id = tagSplit[0]

	offset := dummyOffset
	/*
		// I don't plan to use offset feature but I'll leave it in for now.
		if len(tagSplit) == 2 {
			// If ?t=XmXs add that time to offset
			offset, _ = time.ParseDuration(tagSplit[1])
		}
	*/

	directory = viper.GetString("localstorage.directory") + "/"
	filePath = directory + id + fileExtension
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, err
	}

	jsonPath = directory + id + ".json"
	jsonFile, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}

	v, err = jason.NewObjectFromBytes(jsonFile)
	if err != nil {
		return nil, err
	}

	// If we got this far I assume json file
	// has all needed stuff.
	title, _ := v.GetString("title")
	thumbnail, _ := v.GetString("thumbnail")
	artist, _ := v.GetString("artist")
	durationString, _ := v.GetString("duration")
	durationConverted, _ := duration.FromString(durationString)
	duration := durationConverted.ToDuration()

	track := bot.Track{
		ID:             id,
		URL:            "localstorage/" + id,
		Title:          title,
		Author:         artist,
		Submitter:      submitter.Name,
		Service:        ls.ReadableName,
		Filename:       id + fileExtension,
		ThumbnailURL:   thumbnail,
		Duration:       duration,
		PlaybackOffset: offset,
		Playlist:       nil,
	}

	tracks = append(tracks, track)
	return tracks, nil
}
