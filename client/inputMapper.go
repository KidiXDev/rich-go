package client

import (
	"time"
)

type ActivityType int

const (
	ActivityTypePlaying   ActivityType = 0
	ActivityTypeStreaming ActivityType = 1
	ActivityTypeListening ActivityType = 2
	ActivityTypeWatching  ActivityType = 3
	ActivityTypeCustom    ActivityType = 4
	ActivityTypeCompeting ActivityType = 5
)

type StatusDisplayType int

const (
	StatusDisplayTypeName    StatusDisplayType = 0
	StatusDisplayTypeState   StatusDisplayType = 1
	StatusDisplayTypeDetails StatusDisplayType = 2
)

// Activity holds the data for discord rich presence
type Activity struct {
	// Name used in activity text, e.g. "Listening to {name}"
	Name string
	// Type of activity, e.g. Playing (0), Listening (2), Watching (3)
	Type ActivityType
	// Which field is highlighted in the activity header
	StatusDisplayType *StatusDisplayType
	// What the player is currently doing
	Details string
	// URL opened when clicking the details text
	DetailsURL string
	// The user's current party status
	State string
	// URL opened when clicking the state text
	StateURL string
	// The id for a large asset of the activity, usually a snowflake
	LargeImage string
	// Text displayed when hovering over the large image of the activity
	LargeText string
	// URL opened when clicking the large image
	LargeURL string
	// The id for a small asset of the activity, usually a snowflake
	SmallImage string
	// Text displayed when hovering over the small image of the activity
	SmallText string
	// URL opened when clicking the small image
	SmallURL string
	// Information for the current party of the player
	Party *Party
	// Unix timestamps for start and/or end of the game
	Timestamps *Timestamps
	// Secrets for Rich Presence joining and spectating
	Secrets *Secrets
	// Clickable buttons that open a URL in the browser
	Buttons []*Button
	// Whether this activity is an instanced game session
	Instance *bool
}

// Button holds a label and the corresponding URL that is opened on press
type Button struct {
	// The label of the button
	Label string
	// The URL of the button
	Url string
}

// Party holds information for the current party of the player
type Party struct {
	// The ID of the party
	ID string
	// Used to show the party's current size
	Players int
	// Used to show the party's maximum size
	MaxPlayers int
}

// Timestamps holds unix timestamps for start and/or end of the game
type Timestamps struct {
	// unix time (in milliseconds) of when the activity started
	Start *time.Time
	// unix time (in milliseconds) of when the activity ends
	End *time.Time
}

// NewListeningActivity builds a Discord "Listening to ..." activity payload
// that music players can reuse with minimal boilerplate.
func NewListeningActivity(track, artist, album string, startedAt, endAt *time.Time) Activity {
	statusDisplayType := StatusDisplayTypeDetails

	activity := Activity{
		Type:              ActivityTypeListening,
		StatusDisplayType: &statusDisplayType,
		Name:              track,
		Details:           artist,
		State:             album,
	}

	if startedAt != nil || endAt != nil {
		activity.Timestamps = &Timestamps{
			Start: startedAt,
			End:   endAt,
		}
	}

	return activity
}

// Secrets holds secrets for Rich Presence joining and spectating
type Secrets struct {
	// The secret for a specific instanced match
	Match string
	// The secret for joining a party
	Join string
	// The secret for spectating a game
	Spectate string
}

func mapActivity(activity *Activity) *PayloadActivity {
	activityType := int(activity.Type)
	final := &PayloadActivity{
		Name:       activity.Name,
		Type:       &activityType,
		Details:    activity.Details,
		DetailsURL: activity.DetailsURL,
		State:      activity.State,
		StateURL:   activity.StateURL,
		Assets: PayloadAssets{
			LargeImage: activity.LargeImage,
			LargeText:  activity.LargeText,
			LargeURL:   activity.LargeURL,
			SmallImage: activity.SmallImage,
			SmallText:  activity.SmallText,
			SmallURL:   activity.SmallURL,
		},
		Instance: activity.Instance,
	}

	if activity.StatusDisplayType != nil {
		statusDisplayType := int(*activity.StatusDisplayType)
		final.StatusDisplayType = &statusDisplayType
	}

	if activity.Timestamps != nil {
		timestamps := &PayloadTimestamps{}
		if activity.Timestamps.Start != nil {
			start := uint64(activity.Timestamps.Start.UnixNano() / 1e6)
			timestamps.Start = &start
		}
		if activity.Timestamps.End != nil {
			end := uint64(activity.Timestamps.End.UnixNano() / 1e6)
			timestamps.End = &end
		}
		if timestamps.Start != nil || timestamps.End != nil {
			final.Timestamps = timestamps
		}
	}

	if activity.Party != nil {
		final.Party = &PayloadParty{
			ID:   activity.Party.ID,
			Size: [2]int{activity.Party.Players, activity.Party.MaxPlayers},
		}
	}

	if activity.Secrets != nil {
		final.Secrets = &PayloadSecrets{
			Join:     activity.Secrets.Join,
			Match:    activity.Secrets.Match,
			Spectate: activity.Secrets.Spectate,
		}
	}

	if len(activity.Buttons) > 0 {
		for _, btn := range activity.Buttons {
			final.Buttons = append(final.Buttons, &PayloadButton{
				Label: btn.Label,
				Url:   btn.Url,
			})
		}
	}

	return final
}
