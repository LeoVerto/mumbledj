
<h1 align="center">MumbleDJ</h1>
<p align="center"><b>A Mumble bot that plays audio fetched from various media websites.</b></p>

Fork of [matthieugrieger's mumbledj](https://github.com/matthieugrieger/mumbledj)

## Table of Contents

* [Features](#features)
* [Installation](#installation)
  * [Requirements](#requirements)
    * [YouTube API Key](#youtube-api-key)
    * [SoundCloud API Key](#soundcloud-api-key)
  * [Via `go get`](#via-go-get-recommended)
  * [Pre-compiled Binaries](#pre-compiled-binaries-easiest)
  * [From Source](#from-source)
  * [Docker](#docker)
* [Usage](#usage)
* [Commands](#commands)
* [Contributing](#contributing)
* [Author](#author)
* [License](#license)
* [Thanks](#thanks)

## Features
* Plays audio from many media websites, including YouTube, SoundCloud, and Mixcloud.
* Supports playlists and individual videos/tracks.
* Displays metadata in the text chat whenever a new track starts playing.
* Incredibly customizable. Nearly everything is able to be tweaked via configuration files (by default located at `$HOME/.config/mumbledj/config.yaml`).
* A large array of [commands](#commands) that perform a wide variety of functions.
* Built-in vote-skipping.
* Built-in caching system (disabled by default).
* Built-in play/pause/volume control.

## Installation
**IMPORTANT NOTE:** MumbleDJ is only tested and developed for Linux systems. Support will not be given for non-Linux systems if problems are encountered.

### Requirements
**All MumbleDJ installations must also have the following installed:**
* [`youtube-dl`](https://rg3.github.io/youtube-dl/download.html)
* [`ffmpeg`](https://ffmpeg.org) OR [`avconv`](https://libav.org)
* [`aria2`](https://aria2.github.io/) if you plan on using services that throttle download speeds (like Mixcloud)

**If installing via `go install` or from source, the following must be installed:**
* [Go 1.5+](https://golang.org)
  * __NOTE__: Extra installation steps are required for a working Go installation. Once Go is installed, type `go help gopath` for more information.
  * If the repositories for your distro contain a version of Go older than 1.5, try using [`gvm`](https://github.com/moovweb/gvm) to install Go 1.5 or newer.

#### YouTube API Key
A YouTube API key must be present in your configuration file in order to use the YouTube service within the bot. Below is a guide for retrieving an API key:

**1)** Navigate to the [Google Developers Console](https://console.developers.google.com) and sign in with your Google account, or create one if you haven't already.

**2)** Click the "Create Project" button and give your project a name. It doesn't matter what you set your project name to. Once you have a name click the "Create" button. You should be redirected to your new project once it's ready.

**3)** Click on "APIs & auth" on the sidebar, and then click APIs. Under the "YouTube APIs" header, click "YouTube Data API". Click on the "Enable API" button.

**4)** Click on the "Credentials" option underneath "APIs & auth" on the sidebar. Underneath "Public API access" click on "Create New Key". Choose the "Server key" option.

**5)** Add the IP address of the machine MumbleDJ will run on in the box that appears (this is optional, but improves security). Click "Create".

**6)** You should now see that an API key has been generated. Copy/paste this API key into the configuration file located at `$HOME/.config/mumbledj/mumbledj.yaml`.

#### SoundCloud API Key
A SoundCloud client ID must be present in your configuration file in order to use the SoundCloud service within the bot. Below is a guide for retrieving a client ID:

**1)** Login/sign up for a SoundCloud account on https://soundcloud.com.

**2)** Create a new app: https://soundcloud.com/you/apps/new.

**3)** You should now see that a client ID has been generated. Copy/paste this ID (NOT the client secret) into the configuration file located at `$HOME/.config/mumbledj/mumbledj.yaml`.


### Via `go get` (might not work)
After verifying that the [requirements](#requirements) are installed, simply issue the following command:
```
go get -u git.roshless.me/Roshless/MumbleDJ
```

This should place a binary in `$GOPATH/bin` that can be used to start the bot.

### From Source
First, clone the MumbleDJ repository to your machine:
```
git clone https://git.roshless.me/roshless/mumbledj.git
```

Install the required software as described in the [requirements section](#requirements), and execute the following:
```
make
```

This will place a compiled `mumbledj` binary in the cloned directory if successful. If you would like to make the binary more accessible by adding it to `/usr/local/bin`, simply execute the following:
```
sudo make install
```

## Usage
```
NAME:
   MumbleDJ - A Mumble bot that plays audio from various media sites.

USAGE:
   mumbledj [global options] command [command options] [arguments...]

VERSION:
   v3.1.0

COMMANDS:
GLOBAL OPTIONS:
   --config value, -c value		location of MumbleDJ configuration file (default: "/home/$USER/.config/mumbledj/config.yaml")
   --server value, -s value		address of Mumble server to connect to (default: "127.0.0.1")
   --port value, -o value		port of Mumble server to connect to (default: "64738")
   --username value, -u value		username for the bot (default: "MumbleDJ")
   --password value, -p value		password for the Mumble server
   --channel value, -n value		channel the bot enters after connecting to the Mumble server
   --p12 value				path to user p12 file for authenticating as a registered user
   --cert value, -e value		path to PEM certificate
   --key value, -k value		path to PEM key
   --accesstokens value, -a value	list of access tokens separated by spaces
   --insecure, -i			if present, the bot will not check Mumble certs for consistency
   --debug, -d				if present, all debug messages will be shown
   --help, -h				show help
   --version, -v			print the version

```

__NOTE__: You can also override all settings found within `config.yaml` directly from the commandline. Here's an example:

```
mumbledj --admins.names="SuperUser,Matt" --volume.default="0.5" --volume.lowest="0.2" --queue.automatic_shuffle_on="true"
```

Keep in mind that values that contain commas (such as `"SuperUser,Matt"`) will be interpreted as string slices, or arrays if you are not familiar with Go. If you want your value to be interpreted as a normal string, it is best to avoid commas for now.

## Commands

You can check that in config.yaml
