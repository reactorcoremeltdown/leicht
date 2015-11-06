function leicht_send_message() {
    chat_id=$1
    reply_id=$2
    message_text=`echo "$3" | jq -R -s '.'`
    disable_web_preview=$4
    leicht_socket=$5

    echo '{ "actionType": "SendMessage", "actionSettings": {"chatID": '$chat_id', "replyToMessageID": '$reply_id', "text": '$message_text', "disableWebPagePreview": '$disable_web_preview'}}' | socat stdio unix-connect:$5
}

function leicht_send_message_to_channel() {
    channel_username=$1
    reply_id=$2
    message_text=`echo "$3" | jq -R -s '.'`
    disable_web_preview=$4
    leicht_socket=$5

    echo '{ "actionType": "SendMessage", "actionSettings": {"channelUsername": "'$channel_username'", "replyToMessageID": '$reply_id', "text": '$message_text', "disableWebPagePreview": '$disable_web_preview'}}' | socat stdio unix-connect:$5
}

