# FM Shockwave

<img src="shockwave.jpg" width="122" height="150" align="right" />

SOON\_ FM Volume Management Service. This service subscribes to a Redis
Pub/Sub service listening for volume change events. Once an event has been
recieved it changes the volume and publishes a volume chnaged event back
to the channel.

Note: Currently this can only run on **Linux**.

## Dependencies

Ensure that `libasound2-dev` and `g++` is installed.

## Install

Ensure Go is installed your `$GOPATH` set as you desire.

```
go get github.com/thisissoon/FM-Shockwave/...
go install github.com/thisissoon/FM-Shockwave/...
```

This `shockwave` binary will be installed to `$GOPATH/bin/shockwave`.

## Usage

The application has the following usage options:

```
Usage:
  shockwave [flags]
Flags:
  -c, --channel="": Redis Channel Name
  -d, --device="default": Audio Device Name
  -h, --help=false: help for shockwave
      --max_volume=100: Max Volume Level
      --min_volume=0: Min Volume Level
  -m, --mixer="PCM": Audio Mixer Name
  -r, --redis="127.0.0.1:6379": Redis Server Address
```
