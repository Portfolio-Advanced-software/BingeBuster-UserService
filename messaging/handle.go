package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Portfolio-Advanced-software/BingeBuster-UserService/globals"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Message struct {
	UserId           string `json:"user_id"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	DateOfBirth      string `json:"date_of_birth"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	CreditCardNumber int32  `json:"creditcard_number"`
	ExpirationDate   string `json:"expiration_date"`
	CVC              int32  `json:"cvc"`
	Action           string `json:"action"`
}

func HandleMessage(body []byte) error {
	jsonStr := string(body)
	var msg Message
	err := json.Unmarshal([]byte(jsonStr), &msg)
	if err != nil {
		log.Println("Failed to unmarshal JSON:", err)
		return err
	}

	switch msg.Action {
	case "saveRecord":
		// Insert the data into the database, result contains the newly generated Object ID for the new document
		_, err := globals.UserDb.InsertOne(globals.MongoCtx, msg)
		// check for potential errors
		if err != nil {
			// return internal gRPC error to be handled later
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Internal error: %v", err),
			)
		}
	default:
		fmt.Println("Unknown action:", msg.Action)
	}

	return nil
}
