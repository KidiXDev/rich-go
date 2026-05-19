# rich-go
An implementation of Discord's rich presence in Golang for Linux, macOS and Windows

## Project Status

This repository is a maintained continuation of the original [hugolgst/rich-go](https://github.com/hugolgst/rich-go) project.

## What's Updated In This Fork

- Kept IPC RPC flow compatible with current Discord desktop clients.
- Updated activity payload support for newer fields such as:
  - `type` (including `Listening` for music-player presence)
  - `status_display_type`
  - `details_url` and `state_url`
  - `large_url` and `small_url` in assets
  - `instance`
- Added activity type constants for easier usage:
  - `ActivityTypePlaying`
  - `ActivityTypeStreaming`
  - `ActivityTypeListening`
  - `ActivityTypeWatching`
  - `ActivityTypeCustom`
  - `ActivityTypeCompeting`

## Installation

Install `github.com/KidiXDev/rich-go`:

```
$ go get github.com/KidiXDev/rich-go
```

## Usage

First of all import rich-go
```golang
import "github.com/KidiXDev/rich-go/client"
```

then login by sending the first handshake
```golang
err := client.Login("DISCORD_APP_ID")
if err != nil {
	panic(err)
}
```

and you can set the Rich Presence activity:
```golang
now := time.Now()
activityType := client.ActivityTypeListening
statusDisplayType := client.StatusDisplayTypeDetails

err = client.SetActivity(client.Activity{
	Type:              activityType,
	StatusDisplayType: &statusDisplayType,
	Name:              "Never Gonna Give You Up",
	Details:           "Rick Astley",
	State:             "Whenever You Need Somebody",
	LargeImage:        "album_cover",
	LargeText:         "Never Gonna Give You Up",
	SmallImage:        "music",
	SmallText:         "Listening to...",
	Timestamps: &client.Timestamps{
		Start: &now,
	},
	Buttons: []*client.Button{
		{
			Label: "Open Player",
			Url:   "https://example.com/player",
		},
	},
})

if err != nil {
	panic(err)
}
```

`Activity.Type` supports Discord activity types via constants:
- `client.ActivityTypePlaying`
- `client.ActivityTypeStreaming`
- `client.ActivityTypeListening`
- `client.ActivityTypeWatching`
- `client.ActivityTypeCustom`
- `client.ActivityTypeCompeting`

More details in the [example](https://github.com/KidiXDev/rich-go/blob/master/example/main.go)

## Contributing

1. Fork it (https://github.com/KidiXDev/rich-go/fork)
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create a new Pull Request