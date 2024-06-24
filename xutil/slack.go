package xutil

import (
	"bytes"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

type Slack struct {
	SlackClient *slack.Client
	ChannelID   string
}

func (s *Slack) SendSlackMessage(message string) error {
	ChannelID, timestamp, err := s.SlackClient.PostMessage(s.ChannelID, slack.MsgOptionText(message, false))
	if err != nil {
		return err
	}
	log.Info().Str("message", message).Msgf("Message sent successfully to %s channel at %s", ChannelID, timestamp)
	return nil
}

// https://hooks.slack.com/services/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
	
func PostSlackBetaSignup(req *http.Request, msg string) string {
	ctx := appengine.NewContext(req);
	request := urlfetch.Client(ctx);
	data := []byte("{'text': '" + msg + "'}");
	body := bytes.NewReader(data);
	resp, err := request.Post("https://hooks.slack.com/services/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", "application/json", body);
	if err != nil {
		return err.Error();
	} else {
		return resp.Status;
	}
}