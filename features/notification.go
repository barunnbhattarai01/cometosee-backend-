package features

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/messaging"

	firebase "firebase.google.com/go/v4"
)

var FirebaseApp *firebase.App

func SendCallNotification(fcmToken string, callerName string, channelName string) error {

	ctx := context.Background()

	client, err := FirebaseApp.Messaging(ctx)
	if err != nil {
		return err
	}

	message := &messaging.Message{
		Token: fcmToken,

		Notification: &messaging.Notification{
			Title: "Incoming Call ",
			Body:  fmt.Sprintf("%s is calling you", callerName),
		},

		Data: map[string]string{
			"type":         "call",
			"caller_name":  callerName,
			"channel_name": channelName,
		},
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		return err
	}

	fmt.Println("FCM sent:", response)
	return nil
}
