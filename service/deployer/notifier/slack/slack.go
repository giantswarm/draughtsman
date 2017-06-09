package slack

import (
	"fmt"
	"time"

	"github.com/nlopes/slack"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	eventerspec "github.com/giantswarm/draughtsman/service/deployer/eventer/spec"
	"github.com/giantswarm/draughtsman/service/deployer/notifier/spec"
)

const (
	// goodColour is the colour to use for success Slack messages.
	goodColour = "good"
	// dangerColour is the colour to use for failure Slack messages.
	dangerColour = "danger"

	// titleFormat is the format for titles for Slack messages.
	titleFormat = "%v - %v"
	// successMessage is the message for success Slack messages.
	successMessage = "Successfully deployed"
	// failedMessageFormat is the format for failure Slack messages.
	failedMessageFormat = "Encountered an error ```%v```"
	// footerFormat is the format for footers for Slack messages.
	footerFormat = "%v (%v)"
)

// SlackNotifierType is an Notifier that uses Slack.
var SlackNotifierType spec.NotifierType = "SlackNotifier"

// Config represents the configuration used to create a Slack Notifier..
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Channel     string
	Emoji       string
	Environment string
	SlackToken  string
	Username    string
}

// DefaultConfig provides a default configuration to create a new Slack
// Notifier by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,
	}
}

// New creates a new configured Slack Notifier.
func New(config Config) (*SlackNotifier, error) {
	if config.Channel == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "channel must not be empty")
	}
	if config.Emoji == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "emoji must not be empty")
	}
	if config.Environment == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "environment must not be empty")
	}
	if config.SlackToken == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "slack token must not be empty")
	}
	if config.Username == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "username must not be empty")
	}

	client := slack.New(config.SlackToken)

	if _, err := client.AuthTest(); err != nil {
		return nil, microerror.MaskAnyf(err, "could not authenticate with slack")
	}

	notifier := &SlackNotifier{
		// Dependencies.
		client: client,
		logger: config.Logger,

		// Settings.
		channel:     config.Channel,
		emoji:       config.Emoji,
		environment: config.Environment,
		slackToken:  config.SlackToken,
		username:    config.Username,
	}

	return notifier, nil
}

// SlackNotifier is an implementation of the Notifier interface,
// that uses Slack.
type SlackNotifier struct {
	// Dependencies.
	client *slack.Client
	logger micrologger.Logger

	// Settings.
	channel     string
	emoji       string
	environment string
	slackToken  string
	username    string
}

// postSlackMessage takes a DeploymentEvent and a possible error message,
// and posts a helpful message to the configured Slack channel.
func (n *SlackNotifier) postSlackMessage(event eventerspec.DeploymentEvent, errorMessage string) error {
	startTime := time.Now()
	defer updateSlackMetrics(startTime)

	success := true
	if len(errorMessage) > 0 {
		success = false
	}

	attachment := slack.Attachment{}

	if success {
		attachment.Color = goodColour
	} else {
		attachment.Color = dangerColour
	}

	attachment.MarkdownIn = []string{"text"}

	attachment.Title = fmt.Sprintf(titleFormat, event.Name, event.Sha)
	if success {
		attachment.Text = successMessage
	} else {
		attachment.Text = fmt.Sprintf(failedMessageFormat, errorMessage)
	}
	attachment.Footer = fmt.Sprintf(footerFormat, n.environment, event.ID)

	params := slack.PostMessageParameters{}

	params.Username = n.username
	params.IconEmoji = n.emoji
	params.Attachments = []slack.Attachment{attachment}

	_, _, err := n.client.PostMessage(n.channel, "", params)
	if err != nil {
		return microerror.MaskAny(err)
	}

	return nil
}

func (n *SlackNotifier) Success(event eventerspec.DeploymentEvent) error {
	n.logger.Log("debug", "sending success message to slack")

	return n.postSlackMessage(event, "")
}

func (n *SlackNotifier) Failed(event eventerspec.DeploymentEvent, errorMessage string) error {
	n.logger.Log("debug", "sending failed message to slack")

	return n.postSlackMessage(event, errorMessage)
}
