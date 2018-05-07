package modules

import "github.com/sirupsen/logrus"

type IModule interface {
	HandleCommands()
	SetLogger(*logrus.Logger)
	GetBaseCommand() string
}
