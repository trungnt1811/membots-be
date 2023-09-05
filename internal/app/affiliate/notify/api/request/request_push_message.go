package request

type PushNotification struct {
	// Topic is the topic of the message.
	Topic string `json:"Topic,omitempty"`

	// Data is the message data.
	Data map[string]string `json:"Data,omitempty"`

	// CustomID is another ID of the request (optional), for client-side tracking purposes. If provided, it must be a valid 64-character long hex string.
	CustomID string `json:"CustomID,omitempty"`

	// Notification is the notification message.
	Notification struct {
		//Title is the title of the message.
		Title string `json:"Title,omitempty"`

		// Body is the message body.
		Body string `json:"Body,omitempty"`

		// ImageURL is the URL of the notification image.
		ImageURL string `json:"Image,omitempty"`
	} `json:"Notification,omitempty"`
}

// APISendPushRequest specifies the necessary fields for sending a push notification.
type APISendPushRequest struct {
	PushNotification

	// Recipients is the list of userIDs to send the message to.
	Recipients []uint `json:"Recipients"`

	// Webhook is the URl for receiving the status of the request.
	Webhook string `json:"Webhook"`
}
