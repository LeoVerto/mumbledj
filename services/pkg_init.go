/*
 * MumbleDJ
 * By Matthieu Grieger
 * services/pkg_init.go
 * Copyright (c) 2016 Matthieu Grieger (MIT License)
 */

package services

import (
	"git.roshless.me/roshless/mumbledj/bot"
	"git.roshless.me/roshless/mumbledj/interfaces"
)

// DJ is an injected MumbleDJ struct.
var DJ *bot.MumbleDJ

// Services is a slice of enabled MumbleDJ services.
var Services []interfaces.Service

func init() {
	Services = []interfaces.Service{
		NewLocalStorageService(),
		NewMixcloudService(),
		NewSoundCloudService(),
		NewYouTubeService(),
	}
}
