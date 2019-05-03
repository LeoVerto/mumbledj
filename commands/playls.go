/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/aplayls.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"git.roshless.me/roshless/mumbledj/interfaces"
	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
)

// PlayLsCommand is a command that adds to queue every local storage song.
type PlayLsCommand struct{}

// Aliases returns the current aliases for the command.
func (c *PlayLsCommand) Aliases() []string {
	return viper.GetStringSlice("commands.playls.aliases")
}

// Description returns the description for the command.
func (c *PlayLsCommand) Description() string {
	return viper.GetString("commands.playls.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *PlayLsCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.playls.is_admin")
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
func (c *PlayLsCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {
	var (
		tag       string
		allTracks []interfaces.Track
		tracks    []interfaces.Track
		service   interfaces.Service
		err       error
	)

	if viper.GetBool("localstorage.enabled") == false {
		return "", true, errors.New(viper.GetString("common_messages.disabled_error"))
	}

	if len(args) != 0 {
		return "", true, errors.New(viper.GetString("commands.playls.messages.syntax_error"))
	}

	files, err := ioutil.ReadDir(viper.GetString("localstorage.directory"))
	if err != nil {
		return "", true, errors.New(viper.GetString("commands.playls.messages.directory_error"))
	}

	for i, f := range files {
		if i%2 == 1 {
			tag = f.Name() + ".ls"
			if service, err = DJ.GetService(tag); err == nil {
				tracks, err = service.GetTracks(tag, user)
				if err == nil {
					allTracks = append(allTracks, tracks...)
				}
			}
		}
	}

	// We dont use ShuffleTracks because it seeds once.
	rand.Seed(time.Now().UnixNano())
	for i := range allTracks {
		j := rand.Intn(i + 1)
		allTracks[i], allTracks[j] = allTracks[j], allTracks[i]
	}

	if len(allTracks) == 0 {
		return "", true, errors.New(viper.GetString("commands.playls.messages.directory_empty"))
	}

	numTooLong := 0
	numAdded := 0
	for _, track := range allTracks {
		if err = DJ.Queue.AppendTrack(track); err != nil {
			numTooLong++
		} else {
			numAdded++
		}
	}

	retString := fmt.Sprintf(viper.GetString("commands.playls.messages.many_tracks_added"), user.Name, numAdded)
	if numTooLong != 0 {
		retString += fmt.Sprintf(viper.GetString("commands.add.messages.num_tracks_too_long"), numTooLong)
	}
	return retString, false, nil
}
