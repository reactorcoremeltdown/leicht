package main

import (
    "log"
    "net"
    "syscall"
    "strconv"
    "time"
    "os"
    "os/exec"
    "os/signal"
    "encoding/json"
    //"github.com/Syfaro/telegram-bot-api"
    "gopkg.in/go-telegram-bot-api/telegram-bot-api.v2"
    goopt "github.com/droundy/goopt"
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
            if CfgParams.Logging {
                logFilename := ''
                if update.Message.Chat.ID > 0 {
                    logFilename = CfgParams.LogDirectory +
                        "/user-" +
                        strconv.Itoa(update.Message.Chat.ID) +
                        ".log"
                } else {
                    logFilename = CfgParams.LogDirectory +
                        "/group" +
                        strconv.Itoa(update.Message.Chat.ID) +
                        ".log"
                }
                file, err := os.OpenFile(logFilename,
                                            os.O_WRONLY|os.O_APPEND|os.O_CREATE,
                                            0644)
                if err != nil {
                    log.Printf("Error at: %s\n", err.Error())
                }
                _, err = file.WriteString("[" +
                    time.Now().Format("2006-01-02T15:04:05-07:00") +
                    "] <" +
                    update.Message.From.UserName +
                    "> " +
                    update.Message.Text + "\n")
                if err != nil {
                    log.Printf("Error at: %s\n", err.Error())
                }
                file.Close()
            }

            if !CfgParams.WhitelistEnabled || usernameInWhitelist(update.Message.From.UserName, CfgParams.Whitelist) {
                cmd := exec.Command(CfgParams.Script,
                    update.Message.From.UserName,
                    strconv.Itoa(update.Message.Chat.ID),
                    strconv.Itoa(update.Message.MessageID),
                    update.Message.Text)
                err := cmd.Start()
                if err != nil {
                    log.Printf("Error at: %s\n", err.Error())
                }
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
            log.Fatalf("Error at: %s\n", err.Error())
        }
        for {
            conn, err := l.AcceptUnix()
            if err != nil {
                log.Fatalf("Error at: %s\n", err.Error())
            }
            var buf [1024]byte
            n, err := conn.Read(buf[:])
            if err != nil {
                log.Fatalf("Error at: %s\n", err.Error())
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
            log.Printf(err.Error())
        }
        var actionType string
        err = json.Unmarshal(*action["actionType"], &actionType)
        if err != nil {
            log.Printf(err.Error())
        }
        switch actionType {
            case "SendMessage":
                var settings tgbotapi.MessageConfig
                err = json.Unmarshal(*action["actionSettings"], &settings)
                if err != nil {
                    log.Printf(err.Error())
                }
                bot.Send(settings)
                if CfgParams.Logging {
                    logFilename := CfgParams.LogDirectory + "/" + strconv.Itoa(settings.ChatID) + ".log"
                    file, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
                    if err != nil {
                        log.Printf("Error at: %s\n", err.Error())
                    }
                    _, err = file.WriteString("[" + time.Now().Format("2006-01-02T15:04:05-07:00") + "] <" + bot.Self.UserName + "> " + settings.Text + "\n")
                    if err != nil {
                        log.Printf("Error at: %s\n", err.Error())
                    }
                    file.Close()
                }
        }
    }
}
