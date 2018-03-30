/*
 * MumbleDJ
 * By Matthieu Grieger
 * commands/listlocalstorage.go
 * Copyright (c) 2018 Roshless (MIT License)
 */

 package commands

 import (
	 "bytes"
	 "errors"
	 "fmt"
	 "strconv"
 
	 "git.roshless.me/roshless/mumbledj/interfaces"
	 "github.com/layeh/gumble/gumble"
	 "github.com/spf13/viper"
 )
 
 // ListLocalStorageCommand is a command that lists the tracks that were
 // added to Local Storage service.
 type ListLocalStorageCommand struct{}
 
 // Aliases returns the current aliases for the command.
 func (c *ListLocalStorageCommand) Aliases() []string {
	 return viper.GetStringSlice("commands.listtracks.aliases")
 }
 
 // Description returns the description for the command.
 func (c *ListLocalStorageCommand) Description() string {
	 return viper.GetString("commands.listtracks.description")
 }
 
 // IsAdminCommand returns true if the command is only for admin use, and
 // returns false otherwise.
 func (c *ListLocalStorageCommand) IsAdminCommand() bool {
	 return viper.GetBool("commands.listtracks.is_admin")
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
	
	//if no songs in localstorage return

	//list songs with names from json etc

	//either return 1 long ass message and chop it up while writing
	// or just do it here before.

	 return buffer.String(), true, nil
 }
 