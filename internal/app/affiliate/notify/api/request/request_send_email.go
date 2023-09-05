package request

import (
	"fmt"
	"net/mail"
)

type BaseAPISendEmailRequest struct {
	// Recipients is the list of recipients of the email.
	Recipients []string `json:"Recipients" binding:"required"`

	// CC is the list of CC-recipients of the email.
	CC []string `json:"CC,omitempty"`

	// BCC is the list of BCC-recipients of the email
	BCC []string `json:"BCC,omitempty"`
}

func (r *BaseAPISendEmailRequest) ValidateEmail(email string) (bool, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false, err
	}

	return true, nil
}

// ValidateEmails checks if the given list of emails are valid.
func (r *BaseAPISendEmailRequest) ValidateEmails(emails []string) (bool, error) {
	if len(emails) == 0 {
		return true, nil
	}

	for _, email := range emails {
		if _, err := r.ValidateEmail(email); err != nil {
			return false, fmt.Errorf("email `%v` is invalid: %v", email, err)
		}
	}

	return true, nil
}

func (r *BaseAPISendEmailRequest) IsValid() (bool, error) {
	if len(r.Recipients) == 0 {
		return false, fmt.Errorf("empty recipients")
	}
	if _, err := r.ValidateEmails(r.Recipients); err != nil {
		return false, fmt.Errorf("invalid recipients: %v", err)
	}

	if _, err := r.ValidateEmails(r.CC); err != nil {
		return false, fmt.Errorf("invalid CC: %v", err)
	}

	if _, err := r.ValidateEmails(r.BCC); err != nil {
		return false, fmt.Errorf("invalid BCC: %v", err)
	}

	return true, nil
}

// APISendPlainEmailRequest is a request for a
type APISendPlainEmailRequest struct {
	BaseAPISendEmailRequest

	// Subject is the subject of the email.
	Subject string `json:"Subject" binding:"required"`

	// PlainMessage is the text message part of the email.
	PlainMessage string `json:"PlainMessage,omitempty" binding:"required_without=HtmlMessage,omitempty"`

	// HtmlMessage is the HTML part of the message.
	HtmlMessage string `json:"HtmlMessage,omitempty" biding:"required_without=PlainMessage,omitempty"`
}

// IsValid checks if the request is valid.
func (r *APISendPlainEmailRequest) IsValid() (bool, error) {
	if r.PlainMessage+r.HtmlMessage == "" {
		return false, fmt.Errorf("required either PlainMessage or HtmlMessage to be not empty")
	}

	if _, err := r.BaseAPISendEmailRequest.IsValid(); err != nil {
		return false, err
	}

	return true, nil
}

// APISendEmailWithTemplateRequest is a request for a
type APISendEmailWithTemplateRequest struct {
	BaseAPISendEmailRequest

	// TemplateName is the name of the template.
	TemplateName string `json:"TemplateName" binding:"required"`

	// Args consists of arguments for the given template. Should this field not be empty, TemplateName must not be empty.
	Args map[string]interface{} `json:"Args,omitempty"`
}

func (r *APISendEmailWithTemplateRequest) IsValid() (bool, error) {
	if r.TemplateName == "" {
		return false, fmt.Errorf("no template name specified")
	}

	if _, err := r.BaseAPISendEmailRequest.IsValid(); err != nil {
		return false, err
	}

	return true, nil
}
