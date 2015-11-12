package leicht

import (
    "net"
    "encoding/json"
    "github.com/Syfaro/telegram-bot-api"
)

type action struct {
    ActionType string
    ActionSettings interface{}
}

func SendMessage(chatID, replyToMessageID int, text string, disableWebPagePreview bool, socket string) (err error) {
    var act action
    msg := tgbotapi.MessageConfig{
        ChatID: chatID,
        ReplyToMessageID: replyToMessageID,
        Text: text,
        DisableWebPagePreview: disableWebPagePreview,
    }

    act.ActionType = "SendMessage"
    act.ActionSettings = msg

    b, err := json.Marshal(act)
    if err != nil {
        return err
    }

    conn, err := net.DialUnix("unix", nil, &net.UnixAddr{socket, "unix"})
    if err != nil {
        return err
    }

    _, err = conn.Write(b)
    if err != nil {
        return err
    }

    conn.Close()
    return nil
}
