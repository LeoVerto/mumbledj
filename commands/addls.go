/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/addls.go
 * Copyright (c) 2018 Roshless (MIT License)
 */

package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/antonholmquist/jason"
	"github.com/layeh/gumble/gumble"
	"github.com/spf13/viper"
	"github.com/vincent-petithory/dataurl"
)

// AddLsCommand is a command that adds an audio track to a local storage.
// TODO: change name to AddLocalStorage
type AddLsCommand struct{}

// Aliases returns the current aliases for the command.
func (c *AddLsCommand) Aliases() []string {
	return viper.GetStringSlice("commands.addls.aliases")
}

// Description returns the description for the command.
func (c *AddLsCommand) Description() string {
	return viper.GetString("commands.addls.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *AddLsCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.addls.is_admin")
}

// Execute executes the command with the given user and arguments.
// Return value descriptions:
//    string: A message to be returned to the user upon successful execution.
//    bool:   Whether the message should be private or not. true = private,
//            false = public (sent to whole channel).
//    error:  An error message to be returned upon unsuccessful execution.
//            If no error has occurred, pass nil instead.
// Example return statement:
//    return "This is a private message!", true, nil
func (c *AddLsCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	var (
		resp *http.Response
		err  error
		v    *jason.Object
		cmd  *exec.Cmd
	)

	if viper.GetBool("localstorage.enabled") == false {
		return "", true, errors.New(viper.GetString("common_messages.disabled_error"))
	}

	if len(args) != 2 {
		return "", true, errors.New(viper.GetString("commands.addls.messages.syntax_error"))
	}

	// TODO: replace with proper code
	jsonInput :=
		`{
		    "URL": "%s",
		    "title": "%s",
		    "thumbnail": "%s",
		    "artist": "%s",
			"duration": "%s"
		}`

	URL := args[0]
	tag := args[1]

	idUnformatted := strings.Split(URL, "watch?v=")
	if len(idUnformatted) != 2 {
		return "", true, errors.New(viper.GetString("commands.addls.messages.syntax_error"))
	}
	id := idUnformatted[1]

	videoURL := "https://www.googleapis.com/youtube/v3/videos?part=snippet,contentDetails&id=%s&key=%s"
	// Get info about video from youtube
	resp, err = http.Get(fmt.Sprintf(videoURL, id, viper.GetString("api_keys.youtube")))
	if err != nil {
		return "", true, errors.New(viper.GetString("commands.addls.messages.api_error"))
	}
	defer resp.Body.Close()

	v, err = jason.NewObjectFromReader(resp.Body)
	if err != nil {
		return "", true, errors.New(viper.GetString("commands.addls.messages.parsing_error"))
	}
	items, _ := v.GetObjectArray("items")
	if len(items) == 0 {
		return "", true, errors.New(viper.GetString("commands.addls.messages.api_error"))
	}
	// Get desired info from json
	item := items[0]
	title, _ := item.GetString("snippet", "title")
	thumbnail, _ := item.GetString("snippet", "thumbnails", "high", "url")
	author, _ := item.GetString("snippet", "channelTitle")
	duration, _ := item.GetString("contentDetails", "duration")

	player := "--prefer-ffmpeg"
	if viper.GetString("defaults.player_command") == "avconv" {
		player = "--prefer-avconv"
	}

	// Download file to local storage folder
	filePath := viper.GetString("localstorage.directory") + "/" + tag + ".track"
	cmd = exec.Command("youtube-dl", "--verbose", "--no-mtime", "--output", filePath, "--format", "bestaudio", player, URL)
	output, err := cmd.CombinedOutput()
	if err != nil {
		args := ""
		for s := range cmd.Args {
			args += cmd.Args[s] + " "
		}
		logrus.Warnf("%s\n%s\nyoutube-dl: %s", args, string(output), err.Error())
		return "", true, errors.New("Track download failed")
	}

	// Create json file with info
	jsonFilePath := viper.GetString("localstorage.directory") + "/" + tag + ".json"
	fo, err := os.Create(jsonFilePath)
	if err != nil {
		return "", true, errors.New("json file creation failed")
	}
	defer fo.Close()

	// Mumble no longer displays links to images so we have to convert it into data uri
	resp, err = http.Get(thumbnail)
	if err != nil {
		return "", true, errors.New("Couldn't get thumbnail image")
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(make([]byte, 0, resp.ContentLength))
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", true, errors.New("Response into buffer conversion failed")
	}
	img := buf.Bytes()

	thumbnail = dataurl.EncodeBytes(img)

	_, err = io.Copy(fo, strings.NewReader(fmt.Sprintf(jsonInput, URL, title, thumbnail, author, duration)))
	if err != nil {
		return "", true, errors.New("Copying info to json file failed")
	}

	retString := fmt.Sprintf(viper.GetString("commands.addls.messages.track_added"), tag, URL)
	return retString, true, nil
}
