package main

import (
	"encoding/json"
	goopt "github.com/droundy/goopt"
	"gopkg.in/go-telegram-bot-api/telegram-bot-api.v4"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var param_cfgpath = goopt.String([]string{"-c", "--config"}, "/etc/leicht/default.json", "set config file path")

func usernameInWhitelist(username string, whitelist []string) bool {
	present := false
	for _, item := range whitelist {
		if username == item {
			present = true
		}
	}
	return present
}

func writeMessageToLog(chatID int64, logDirectory, botUsername, messageText string)  {
	logFilename := ""
	if chatID > 0 {
		logFilename = logDirectory +
			"/user-" +
			strconv.FormatInt(chatID, 10) +
			".log"
	} else {
		logFilename = logDirectory +
			"/group" +
			strconv.FormatInt(chatID, 10) +
			".log"
	}

	file, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Could not open log file: %s\n", err.Error())
	}
	_, err = file.WriteString("[" + time.Now().Format("2006-01-02T15:04:05-07:00") + "] <" + botUsername + "> " + messageText + "\n")
	if err != nil {
		log.Printf("Could not write log file entry: %s\n", err.Error())
	}
	file.Close()
}

func main() {

	goopt.Description = func() string {
		return "Leicht - universal telegram bot"
	}

	goopt.Version = "0.1"
	goopt.Summary = "leicht -c [config]"
	goopt.Parse(nil)

	CfgParams, _ := LoadConfig(*param_cfgpath)

	msgbus := make(chan string)

	bot, err := tgbotapi.NewBotAPI(CfgParams.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = CfgParams.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}


	go func() {
		for update := range updates {
			var MessageID int
			var MessageText string
			var UserID string
			var ChatID int64

			if update.Message != nil {
				MessageID = update.Message.MessageID
				MessageText = update.Message.Text
				UserID = update.Message.From.UserName
				ChatID = update.Message.Chat.ID
			} else if update.ChannelPost != nil {
				MessageID = update.ChannelPost.MessageID
				MessageText = update.ChannelPost.Text
				UserID = update.ChannelPost.Chat.Title
				ChatID = update.ChannelPost.Chat.ID
			}

			if CfgParams.Logging {
				if CfgParams.WhitelistEnabled {
					if !CfgParams.DoNotLogBlacklisted || usernameInWhitelist(UserID, CfgParams.Whitelist) {
						writeMessageToLog(ChatID,
															CfgParams.LogDirectory,
															UserID,
															MessageText)
					}
				} else {
					writeMessageToLog(ChatID,
														CfgParams.LogDirectory,
														UserID,
														MessageText)
				}
			}

			if !CfgParams.WhitelistEnabled || usernameInWhitelist(UserID, CfgParams.Whitelist) {
				var args []string
				args = append(args,
					UserID,
					strconv.FormatInt(ChatID, 10),
					strconv.Itoa(MessageID))
				if update.Message.Voice != nil {
					url, err := bot.GetFileDirectURL(update.Message.Voice.FileID)
					if err != nil {
						log.Printf("Could not get direct URL of a voice file: %s\n", err.Error())
					}
					args = append(args, "handle_voice "+url)
				} else {
					args = append(args, MessageText)
				}
				cmd := exec.Command(CfgParams.Script, args...)
				err := cmd.Start()
				if err != nil {
					log.Printf("Could not execute command: %s\n", err.Error())
				}
				logstring := make([]interface{}, len(args)+1)
				logstring[0] = CfgParams.Script
				for i := range args {
					logstring[i+1] = args[i]
				}
				log.Printf("Executing script: %s %s %s %s \"%s\"\n", logstring...)
			}
		}
	}()

	go func() {
		if _, err := os.Stat(CfgParams.Socket); err == nil {
			log.Printf("Socket %s exists! Removing...\n", CfgParams.Socket)
			os.Remove(CfgParams.Socket)
		}
		l, err := net.ListenUnix("unix", &net.UnixAddr{CfgParams.Socket, "unix"})
		if err != nil {
			log.Fatalf("Could not open UNIX socket: %s\n", err.Error())
		}
		for {
			conn, err := l.AcceptUnix()
			if err != nil {
				log.Fatalf("Could not accept data from UNIX socket: %s\n", err.Error())
			}
			var buf [1024]byte
			n, err := conn.Read(buf[:])
			if err != nil {
				log.Fatalf("Could not read from UNIX socket: %s\n", err.Error())
			}
			msgbus <- string(buf[:n])
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		for sig := range c {
			os.Remove(CfgParams.Socket)
			log.Printf("Captured %v, Exiting\n", sig)
			os.Exit(0)
		}
	}()

	for {
		msg := <-msgbus
		var action map[string]*json.RawMessage
		data := []byte(msg)
		err := json.Unmarshal(data, &action)
		if err != nil {
			log.Printf("Could not unmarshal incoming JSON: %s\n", err.Error())
		} else {
			var actionType string
			err = json.Unmarshal(*action["actionType"], &actionType)
			if err != nil {
				log.Printf("Could not detect action type: %s\n", err.Error())
			} else {
				switch actionType {
				case "SendMessage":
					var settings tgbotapi.MessageConfig
					err = json.Unmarshal(*action["actionSettings"], &settings)
					if err != nil {
						log.Printf("Could not get action settings: %s\n", err.Error())
					} else {
						bot.Send(settings)
						if CfgParams.Logging {
							writeMessageToLog(settings.ChatID,
																CfgParams.LogDirectory,
																bot.Self.UserName,
																settings.Text)
						}
					}
				}
			}
		}
	}
}
