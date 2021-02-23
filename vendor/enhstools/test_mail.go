package enhstools

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mailgun/mailgun-go"
)

var yourDomain string = "#"
var privateAPIKey string = "#"

func main() {
	mg := mailgun.NewMailgun(yourDomain, privateAPIKey)

	sender := "#"
	subject := "test"
	body := ""
	recipient := "#"

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)
	message.SetTemplate("newmessage-2021-02-21.181649")
	message.AddTemplateVariable("Location_id", "test this")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
