# LEICHT 1 0.1

## NAME

leicht - Telegram bot framework

## SYNOPSIS

**leicht** [OPTION]...

## DESCRIPTION

leicht is a simple multipurpose Telegram bot that can handle and log direct messages or multi-user chat.

**-c, --config FILE**
       Path to JSON config file. Default file path is /etc/kurz/default.json

**--help**
       show usage

## CONFIGURATION

leicht aims to be a framework for building Telegram-oriented services. You can run leicht as a service providing different config files. Here is the description of JSON config fields:

**token**
       Telegram token of the bot

**debug**
       Log all ongoing operations to stderr for debugging. Takes a boolean value(*true* or *false*)

**socket**
       Path to AF_UNIX socket of the bot. You can send a message, send a sticker or any other media to any chat by sending a special JSON-formatted action to the socket. Please refer to the *ACTIONS* section for more information

**script**
       Path to a script that should handle incoming messages. leicht invokes the script with four arguments: user name of Telegram contact, chat ID, message ID and the message text. The script should return processing results via the socket. Please refer to the *LIBRARIES* section for more information

**logging**
       Enable or disable plaintext chat logging. Takes a boolean value(*true* or *false*)

**logDirectory**
       Path to the directory where the plaintext log files will be stored

**whitelistEnabled**
       Enable or disable whitelist. Takes a boolean value(*true* or *false*)

**whitelist**
       A list of Telegram user names who are allowed to invoke the script

## ACTIONS

There are a number of Telegram actions supported by leicht:

**SendMessage**

Sends a message to a chat

```
{
    "actionType": "SendMessage",
    "actionSettings": {
        "chatID": 0,
        "replyToMessageID": 0,
        "text": "foo\nbar",
        "disableWebPagePreview": false
    }
}
```

**SendMessageToChannel**

It is basically the same as *SendMessage* action, but uses *channelUsername* instead of *chatID*

```
{
    "actionType": "SendMessage",
    "actionSettings": {
        "channelUsername": "@examplechannel",
        "replyToMessageID": 0,
        "text": "foo\nbar",
        "disableWebPagePreview": false
    }
}
```

## LIBRARIES

There are some inter-process communication libraries which already ship with leicht. They are intended to help communicating with leicht via its socket. You can easily send messages to recipients using them:

**bash**

```
source /usr/lib/leicht/leicht.sh
leicht_send_message 1234567 0 "Hello, buddy" false /tmp/example.socket
leicht_send_message_to_channel "@examplechannel" 0 "Hello, buddy" false /tmp/example.socket
```

**node.js**

```
var leicht = require('/usr/lib/leicht/leicht.js');
leicht.sendMessage(1234567, 0, "Hello, buddy", false, "/tmp/example.socket");
leicht.sendMessageToChannel("@examplechannel", 0, "Hello, buddy", false, "/tmp/example.socket");
```

**Go**

```
package main

import (
    "log"
    "github.com/Like-all/leicht/src/lib/golang"
)

func main () {
    err := leicht.SendMessage( 1234567, 0, "Hello, buddy", false, "/tmp/example.socket" )
    if err != nil {
        log.Fatal(err)
    }
    err = leicht.SendMessageToChannel( "@examplechannel", 0, "Hello, buddy", false, "/tmp/example.socket" )
    if err != nil {
        log.Fatal(err)
    }
}
```

## TIPS AND TRICKS

There is a special systemd target unit that ships with leicht. If you want to restart your services when leicht updates, you may place Requires-style dependencies from your services on the target. Here is an example:

```
[Unit]
Description=Example bot
After=leicht.target
Requires=leicht.target

[Service]
User=somebody
ExecStart=/usr/bin/leicht -c /etc/leicht/example.json
```

leicht builds for Debian and Ubuntu already ship with postinstall script which restart *leicht.target* during updates.

## BUGS

Please send your bugreports here: https://github.com/Like-all/leicht/issues

## AUTHOR

Written by Azer Abdullaev aka Like-all

## SEE ALSO

systemd.target(5)
