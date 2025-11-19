/*
 * MumbleDJ
 * By Matthieu Grieger
 * main.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package main

import (
	"context"
	"io/ioutil"
	"os"
	"strings"

	"github.com/leoverto/mumbledj/assets"
	"github.com/leoverto/mumbledj/bot"
	"github.com/leoverto/mumbledj/commands"
	"github.com/leoverto/mumbledj/services"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v3"
)

// DJ is a global variable that holds various details about the bot's state.
var DJ = bot.NewMumbleDJ()

// version is supplied by makefile
var version string

// Assets is global variable that allows access to config and sound assets
var Assets = assets.Assets

const assetConfigName = "config.yaml.example"

func init() {
	DJ.Commands = commands.Commands
	DJ.AvailableServices = services.Services

	// Injection into sub-packages.
	commands.DJ = DJ
	services.DJ = DJ
	bot.DJ = DJ

	if version != "" {
		DJ.Version = version
	} else {
		DJ.Version = "v0.0.0"
	}

	logrus.SetLevel(logrus.WarnLevel)
}

func main() {
	cmd := &cli.Command{
		Name:    "MumbleDJ",
		Usage:   "A Mumble bot that plays audio from various media sites.",
		Version: DJ.Version,
	}
	cmd.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   os.ExpandEnv("$HOME/.config/mumbledj/config.yaml"),
			Usage:   "location of MumbleDJ configuration file",
		},
		&cli.StringFlag{
			Name:    "server",
			Aliases: []string{"s"},
			Value:   "127.0.0.1",
			Usage:   "address of Mumble server to connect to",
		},
		&cli.StringFlag{
			Name:    "port",
			Aliases: []string{"o"},
			Value:   "64738",
			Usage:   "port of Mumble server to connect to",
		},
		&cli.StringFlag{
			Name:    "username",
			Aliases: []string{"u"},
			Value:   "MumbleDJ",
			Usage:   "username for the bot",
		},
		&cli.StringFlag{
			Name:    "password",
			Aliases: []string{"p"},
			Value:   "",
			Usage:   "password for the Mumble server",
		},
		&cli.StringFlag{
			Name:    "channel",
			Aliases: []string{"n"},
			Value:   "",
			Usage:   "channel the bot enters after connecting to the Mumble server",
		},
		&cli.StringFlag{
			Name:  "p12",
			Value: "",
			Usage: "path to user p12 file for authenticating as a registered user",
		},
		&cli.StringFlag{
			Name:    "cert",
			Aliases: []string{"e"},
			Value:   "",
			Usage:   "path to PEM certificate",
		},
		&cli.StringFlag{
			Name:    "key",
			Aliases: []string{"k"},
			Value:   "",
			Usage:   "path to PEM key",
		},
		&cli.StringFlag{
			Name:    "accesstokens",
			Aliases: []string{"a"},
			Value:   "",
			Usage:   "list of access tokens separated by spaces",
		},
		&cli.BoolFlag{
			Name:    "insecure",
			Aliases: []string{"i"},
			Usage:   "if present, the bot will not check Mumble certs for consistency",
		},
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "if present, all debug messages will be shown",
		},
	}

	hiddenFlags := make([]cli.Flag, len(viper.AllKeys()))
	for i, configValue := range viper.AllKeys() {
		hiddenFlags[i] = &cli.StringFlag{
			Name:   configValue,
			Hidden: true,
		}
	}
	cmd.Flags = append(cmd.Flags, hiddenFlags...)

	cmd.Action = func(ctx context.Context, c *cli.Command) error {
		if c.Bool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
			//Uncomment to show debug messages from Packr2
			//plog.Logger = logger.New(logger.DebugLevel)
			//Uncomment to show file and line number of log call
			//logrus.SetReportCaller(true)
		}

		for _, configValue := range viper.AllKeys() {
			if c.IsSet(configValue) {
				if strings.Contains(c.String(configValue), ",") {
					viper.Set(configValue, strings.Split(c.String(configValue), ","))
				} else {
					viper.Set(configValue, c.String(configValue))
				}
			}
		}

		viper.SetConfigFile(c.String("config"))
		if err := viper.ReadInConfig(); err != nil {
			logrus.WithFields(logrus.Fields{
				"file":  c.String("config"),
				"error": err.Error(),
			}).Warnln("An error occurred while reading the configuration file. Creating default configuration file...")
			if _, err := os.Stat(c.String("config")); os.IsNotExist(err) {
				createConfigWhenNotExists()
				// If we fail to re-read embedded config, true defaults will be used,
				// which are set in bot/config.go file. So we can safely ignore error here.
				_ = viper.ReadInConfig()
			}
		} else {
			if duplicateErr := bot.CheckForDuplicateAliases(); duplicateErr != nil {
				logrus.WithFields(logrus.Fields{
					"issue": duplicateErr.Error(),
				}).Fatalln("An issue was discoverd in your configuration.")
			}
			createNewConfigIfNeeded()
			viper.WatchConfig()
		}

		if c.IsSet("server") {
			viper.Set("connection.address", c.String("server"))
		}
		if c.IsSet("port") {
			viper.Set("connection.port", c.String("port"))
		}
		if c.IsSet("username") {
			viper.Set("connection.username", c.String("username"))
		}
		if c.IsSet("password") {
			viper.Set("connection.password", c.String("password"))
		}
		if c.IsSet("channel") {
			viper.Set("defaults.channel", c.String("channel"))
		}
		if c.IsSet("p12") {
			viper.Set("connection.user_p12", c.String("p12"))
		}
		if c.IsSet("cert") {
			viper.Set("connection.cert", c.String("cert"))
		}
		if c.IsSet("key") {
			viper.Set("connection.key", c.String("key"))
		}
		if c.IsSet("accesstokens") {
			viper.Set("connection.access_tokens", c.String("accesstokens"))
		}
		if c.IsSet("insecure") {
			viper.Set("connection.insecure", c.Bool("insecure"))
		}

		if err := DJ.Connect(); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatalln("An error occurred while connecting to the server.")
		}

		if viper.GetString("defaults.channel") != "" {
			defaultChannel := strings.Split(viper.GetString("defaults.channel"), "/")
			DJ.Client.Do(func() {
				DJ.Client.Self.Move(DJ.Client.Channels.Find(defaultChannel...))
			})
		}

		DJ.Client.Do(func() {
			DJ.Client.Self.SetComment(viper.GetString("defaults.comment"))
		})
		<-DJ.KeepAlive

		return nil
	}

	_ = cmd.Run(context.Background(), os.Args)
}

func createConfigWhenNotExists() {
	configFile, err := Assets.Find(assetConfigName)
	if err != nil {
		logrus.Warnln("An error occurred while accessing config binary data. A new config file will not be written.")
	} else {
		filePath := os.ExpandEnv("$HOME/.config/mumbledj/config.yaml")
		os.MkdirAll(os.ExpandEnv("$HOME/.config/mumbledj"), 0755)
		writeErr := ioutil.WriteFile(filePath, configFile, 0644)
		if writeErr == nil {
			logrus.WithFields(logrus.Fields{
				"file_path": filePath,
			}).Infoln("A default configuration file has been written.")
		} else {
			logrus.WithFields(logrus.Fields{
				"error": writeErr.Error(),
			}).Warnln("An error occurred while writing a new config file.")
		}
	}
}

func createNewConfigIfNeeded() {
	newConfigPath := os.ExpandEnv("$HOME/.config/mumbledj/config.yaml.new")

	// Check if we should write an updated config file to config.yaml.new.
	if asset, err := Assets.Find(assetConfigName); err == nil {

		assetF, _ := Assets.Open(assetConfigName)
		defer assetF.Close()
		assetInfo, _ := assetF.Stat()
		if configFile, err := os.Open(os.ExpandEnv("$HOME/.config/mumbledj/config.yaml")); err == nil {
			configInfo, _ := configFile.Stat()
			defer configFile.Close()
			if configNewFile, err := os.Open(newConfigPath); err == nil {
				defer configNewFile.Close()
				configNewInfo, _ := configNewFile.Stat()
				if assetInfo.ModTime().Unix() > configNewInfo.ModTime().Unix() {
					// The config asset is newer than the config.yaml.new file.
					// Write a new config.yaml.new file.
					ioutil.WriteFile(os.ExpandEnv(newConfigPath), asset, 0644)
					logrus.WithFields(logrus.Fields{
						"file_path": newConfigPath,
					}).Infoln("An updated default configuration file has been written.")
				}
			} else if assetInfo.ModTime().Unix() > configInfo.ModTime().Unix() {
				// The config asset is newer than the existing config file.
				// Write a config.yaml.new file.
				ioutil.WriteFile(os.ExpandEnv(newConfigPath), asset, 0644)
				logrus.WithFields(logrus.Fields{
					"file_path": newConfigPath,
				}).Infoln("An updated default configuration file has been written.")
			}
		}
	}
}
