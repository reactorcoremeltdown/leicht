package leicht

import (
    "net"
    "encoding/json"
    "github.com/Syfaro/telegram-bot-api"
)

type action struct {
    ActionType string           `json:"actionType"`
    ActionSettings interface{}  `json:"actionSettings"`
}

func sendBytesToSocket(bytes []byte, socket string) (err error) {
    conn, err := net.DialUnix("unix", nil, &net.UnixAddr{socket, "unix"})
    if err != nil {
        return err
    }

    _, err = conn.Write(bytes)
    if err != nil {
        return err
    }

    conn.Close()

    return nil
}

func SendMessage(chatID, replyToMessageID int, text string, disableWebPagePreview bool, socket string) (err error) {
    msg := tgbotapi.MessageConfig{
        ChatID: chatID,
        ReplyToMessageID: replyToMessageID,
        Text: text,
        DisableWebPagePreview: disableWebPagePreview,
    }

    act := action{
        ActionType: "SendMessage",
        ActionSettings: msg,
    }

    b, err := json.Marshal(act)
    if err != nil {
        return err
    }

    err = sendBytesToSocket(b, socket)
    if err != nil {
        return err
    }

    return nil
}


func SendMessageToChannel(channelUsername string, replyToMessageID int, text string, disableWebPagePreview bool, socket string) (err error) {
    msg := tgbotapi.MessageConfig{
        ChannelUsername: channelUsername,
        ReplyToMessageID: replyToMessageID,
        Text: text,
        DisableWebPagePreview: disableWebPagePreview,
    }

    act := action{
        ActionType: "SendMessage",
        ActionSettings: msg,
    }

    b, err := json.Marshal(act)
    if err != nil {
        return err
    }

    err = sendBytesToSocket(b, socket)
    if err != nil {
        return err
    }

    return nil
}
