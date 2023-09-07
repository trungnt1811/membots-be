package request

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
)

// APISendSMSRequest specifies the necessary fields for sending sms.
type APISendSMSRequest struct {
	// Message is the message to be sent.
	Message string `json:"Message"`

	// Recipient is the receiver of the message.
	Recipient string `json:"Recipient"`
}

// ToValidE164PhoneNumber returns the valid E164-confront phone number from the given phone number.
// If the given number is a common VN phone number (ten digits with leading `0`, or beginning with `84`), it will be converted to +84 format.
func ToValidE164PhoneNumber(phoneNumber string) (string, error) {
	if len(phoneNumber) == 10 && phoneNumber[0] == '0' {
		phoneNumber = fmt.Sprintf("+84%v", phoneNumber[1:])
	}
	if len(phoneNumber) == 11 && phoneNumber[:2] == "84" {
		phoneNumber = fmt.Sprintf("+%v", phoneNumber)
	}

	type phoneE164 struct {
		PhoneNumber string `json:"PhoneNumber" validate:"required,e164"`
	}

	v := validator.New()
	err := v.Struct(phoneE164{PhoneNumber: phoneNumber})
	if err != nil {
		return "", err
	}

	return phoneNumber, nil
}

func (req *APISendSMSRequest) UnmarshalJSON(data []byte) error {
	type structHolder struct {
		// Message is the message to be sent.
		Message string `json:"Message" binding:"required" validate:"required"`

		// Recipient is the receiver of the message.
		Recipient string `json:"Recipient" binding:"required" validate:"required"`
	}
	var tmpReq structHolder
	err := json.Unmarshal(data, &tmpReq)
	if err != nil {
		return err
	}

	req.Message = tmpReq.Message
	req.Recipient, err = ToValidE164PhoneNumber(tmpReq.Recipient)
	if err != nil {
		return err
	}

	return nil
}

func (req *APISendSMSRequest) IsValid() (bool, error) {
	v := validator.New()
	type SMSRequest struct {
		Message   string `json:"Message" binding:"required" validate:"required"`
		Recipient string `json:"Recipient" binding:"required" validate:"required,e164"`
	}

	holder := SMSRequest{
		Message:   req.Message,
		Recipient: req.Recipient,
	}
	err := v.Struct(holder)
	if err != nil {
		return false, err
	}

	return true, nil
}
