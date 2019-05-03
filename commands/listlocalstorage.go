/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/listlocalstorage.go
 * Copyright (c) 2018 Roshless (MIT License)
 */

package commands

import (
	"errors"
	"fmt"
	"io/ioutil"

	"git.roshless.me/roshless/mumbledj/interfaces"
	"layeh.com/gumble/gumble"
	"github.com/spf13/viper"
)

// ListLocalStorageCommand is a command that lists the tracks that were
// added to Local Storage service.
type ListLocalStorageCommand struct{}

// Aliases returns the current aliases for the command.
func (c *ListLocalStorageCommand) Aliases() []string {
	return viper.GetStringSlice("commands.listlocalstorage.aliases")
}

// Description returns the description for the command.
func (c *ListLocalStorageCommand) Description() string {
	return viper.GetString("commands.listlocalstorage.description")
}

// IsAdminCommand returns true if the command is only for admin use, and
// returns false otherwise.
func (c *ListLocalStorageCommand) IsAdminCommand() bool {
	return viper.GetBool("commands.listlocalstorage.is_admin")
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
func (c *ListLocalStorageCommand) Execute(user *gumble.User, args ...string) (string, bool, error) {

	var (
		message   string
		tag       string
		allTracks []interfaces.Track
		tracks    []interfaces.Track
		service   interfaces.Service
		err       error
	)

	files, err := ioutil.ReadDir(viper.GetString("localstorage.directory"))
	if err != nil {
		return "", true, errors.New(viper.GetString("commands.localstorage.messages.directory_error"))
	}

	if len(files) == 0 {
		return "", true, errors.New(viper.GetString("commands.localstorage.common_messages.no_songs_error"))
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

	for _, track := range allTracks {
		message += fmt.Sprintf("%v.ls - %v<br>", track.GetID(), track.GetTitle())
	}

	// IDEA: return array of string?
	return message, false, nil
}
