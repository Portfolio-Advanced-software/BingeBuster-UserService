package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Email            string             `bson:"email,omitempty"`
	Phone            string             `bson:"phone,omitempty"`
	DateOfBirth      string             `bson:"dateofbirth,omitempty"`
	FirstName        string             `bson:"firstname,omitempty"`
	LastName         string             `bson:"lastname,omitempty"`
	CreditCardNumber int32              `bson:"creditcardnumber,omitempty"`
	ExpirationDate   string             `bson:"expirationdate,omitempty"`
	CVC              int32              `bson:"cvc,omitempty"`
}
