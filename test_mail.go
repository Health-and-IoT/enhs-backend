package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mailgun/mailgun-go"
)

var yourDomain string = "sandboxcedc9b051efc44eba1486189ca99064b.mailgun.org"
var privateAPIKey string = "a2ce223b0e6cd87b3d75099d78cc96ae-d32d817f-a5693b86"

func main() {
	mg := mailgun.NewMailgun(yourDomain, privateAPIKey)

	sender := "sandboxcedc9b051efc44eba1486189ca99064b.mailgun.org"
	subject := "Fancy subject!"
	body := "Hello from Mailgun Go!"
	recipient := "40405884@mailinator.com"

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
