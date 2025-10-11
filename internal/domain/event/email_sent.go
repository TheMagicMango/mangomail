package event

import (
	"time"
)

type EmailSent struct {
	Name    string
	Payload interface{}
}

func NewEmailSent() *EmailSent {
	return &EmailSent{
		Name: "email.sent",
	}
}

func (e *EmailSent) GetName() string {
	return e.Name
}

func (e *EmailSent) GetPayload() interface{} {
	return e.Payload
}

func (e *EmailSent) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *EmailSent) GetDateTime() time.Time {
	return time.Now()
}
