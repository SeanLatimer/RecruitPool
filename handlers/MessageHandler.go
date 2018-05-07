package handlers

import (
	"strings"

	"github.com/SeanLatimer/RecruitPool/modules"

	twitch "github.com/gempir/go-twitch-irc"
	"github.com/sirupsen/logrus"
)

type MessageHandler struct {
	Logger *logrus.Logger
	Client *twitch.Client
}

// Handle Handles RoomState Messages
func (h *MessageHandler) Handle(channel string, user twitch.User, message twitch.Message) {
	// h.Logger.WithFields(logrus.Fields{
	// 	"channel": channel,
	// 	"user":    user,
	// 	"message": message,
	// }).Info("Message")
	if strings.HasPrefix(message.Text, "!") {
		command := strings.TrimPrefix(message.Text, "!")
		switch {
		case strings.HasPrefix(command, modules.RecruitPool.GetBaseCommand()):
			args := strings.TrimPrefix(command, modules.RecruitPool.GetBaseCommand())
			args = strings.TrimSpace(args)
			resp := modules.RecruitPool.HandleCommands(args, user, message.Tags, (user.Username == channel))
			h.Client.Say(channel, resp)
		}
	}
}
