package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SeanLatimer/RecruitPool/handlers"
	"github.com/SeanLatimer/RecruitPool/modules"
	twitch "github.com/gempir/go-twitch-irc"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Verbose logging
var Verbose = false
var logger *logrus.Logger

const configName = ".recruitpool"

func init() {
	pflag.BoolVarP(&Verbose, "verbose", "V", false, "Verbose logging")
}

func main() {
	pflag.Parse()
	setupLogging()
	setupConfigDefaults()
	initConfigDir()

	if viper.GetString("AuthToken") == "CHANGEME" {
		firstRun()
	}

	client := twitch.NewClient(viper.GetString("Username"), viper.GetString("AuthToken"))

	modules.RecruitPool.SetLogger(logger)

	messageHandler := &handlers.MessageHandler{
		Logger: logger,
		Client: client,
	}
	client.OnNewMessage(messageHandler.Handle)

	client.Join(viper.GetString("Channel"))

	err := client.Connect()
	if err != nil {
		logger.Errorf("Error connecting", err)
	}
}

func firstRun() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Recruit Pool First Run")
	fmt.Println("---------------------")

	fmt.Println()
	fmt.Println("Please enter your username")
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	// convert CRLF to LF
	username = strings.Replace(username, "\r\n", "", -1)
	viper.Set("Username", username)

	// Get AuthToken
	fmt.Println()
	fmt.Println("Please enter your oauth token")
	fmt.Print("OAuthToken: ")
	authToken, _ := reader.ReadString('\n')
	// convert CRLF to LF
	authToken = strings.Replace(authToken, "\r\n", "", -1)
	if !strings.HasPrefix(authToken, "oauth:") {
		authToken = "oauth:" + authToken
	}
	viper.Set("AuthToken", authToken)

	// Get channel
	fmt.Println()
	fmt.Println("Please enter your channel")
	fmt.Print("Channel: ")
	channel, _ := reader.ReadString('\n')
	// convert CRLF to LF
	channel = strings.Replace(channel, "\r\n", "", -1)
	viper.Set("Channel", channel)

	home, err := homedir.Dir()
	if err != nil {
		logger.Fatalf("Error getting home directory", err)
	}

	configPath := fmt.Sprint(home, "/.recruitpool.yaml")
	err = viper.SafeWriteConfigAs(configPath)
	if err != nil {
		logger.Fatalf("Error saving config", err)
	}
	logger.Info("Config saved to: ", configPath)
}

func setupConfigDefaults() {
	viper.SetDefault("AuthToken", "CHANGEME")
	viper.SetDefault("Channel", "CHANGEME")
}

func initConfigDir() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		logger.Fatalf("Error getting home directory", err)
	}

	// Search config in home directory with name (without extension).
	viper.SetConfigType("yaml")
	viper.AddConfigPath(home)
	viper.SetConfigName(configName)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logger.Info("Using config file:", viper.ConfigFileUsed())
	}
}

func setupLogging() {
	if logger != nil {
		return
	}

	logger = logrus.New()
	if Verbose {
		logger.Level = logrus.InfoLevel
	} else {
		logger.Level = logrus.WarnLevel
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Fatal(err)
	}

	path := filepath.Join(dir, "recruitpool.log")
	rotatedLogsDir := filepath.Join(dir, "logs")
	if _, err := os.Stat(rotatedLogsDir); os.IsNotExist(err) {
		err = os.MkdirAll(rotatedLogsDir, 0755)
		if err != nil {
			logger.Fatalf("Failed to create rotated logs directory", err)
		}
	}

	pathWithTimeStamp := filepath.Join(rotatedLogsDir, "recruitpool-%Y%m%d%H%M.log")
	writer, err := rotatelogs.New(
		pathWithTimeStamp,
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
	)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Hooks.Add(lfshook.NewHook(
		writer,
		&logrus.TextFormatter{},
	))
}
