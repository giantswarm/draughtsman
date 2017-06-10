package slack

import (
	"github.com/nlopes/slack"
)

// Client is an interface for Slack clients to implement.
type Client interface {
	AuthTest() (response *slack.AuthTestResponse, error error)
	PostMessage(channel, text string, params slack.PostMessageParameters) (string, string, error)
}
