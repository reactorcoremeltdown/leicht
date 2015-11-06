exports.sendMessage = function(chatID, replyToMessageID, text, disableWebPagePreview, socket) {
	var net = require('net');
	var action = {
		actionType: "SendMessage",
		actionSettings: {
			chatID: chatID,
			replyToMessageID: replyToMessageID,
			text: text,
			disableWebPagePreview: disableWebPagePreview,
			socket: socket
		}
	}
	
	var client = net.connect({path: socket},
		function() {
			client.write(JSON.stringify(action));
		});
}
exports.sendMessageToChannel = function(channelUsername, replyToMessageID, text, disableWebPagePreview, socket) {
	var net = require('net');
	var action = {
		actionType: "SendMessage",
		actionSettings: {
			channelUsername: channelUsername,
			replyToMessageID: replyToMessageID,
			text: text,
			disableWebPagePreview: disableWebPagePreview,
			socket: socket
		}
	}
	
	var client = net.connect({path: socket},
		function() {
			client.write(JSON.stringify(action));
		});
}
