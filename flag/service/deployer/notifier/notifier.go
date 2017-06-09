package notifier

import (
	"github.com/giantswarm/draughtsman/flag/service/deployer/notifier/slack"
)

type Notifier struct {
	Slack slack.Slack
	Type  string
}
