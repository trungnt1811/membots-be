package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astraprotocol/affiliate-system/internal/app/notify/api/request"
	"github.com/astraprotocol/affiliate-system/internal/app/notify/api/response"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type APISendSMSRequest struct {
	// Message is the message to be sent.
	Message string `json:"Message"`

	// Recipient is the receiver of the message.
	Recipient string `json:"Recipient"`

	// CustomID is another ID of the request (optional), for client-side tracking purposes. If provided, it must be a valid 64-character long hex string.
	CustomID string `json:"CustomID,omitempty"`

	// Webhook is the URl for receiving the status of the request.
	Webhook string `json:"Webhook"`
}

type BaseAPISendEmailRequest struct {
	// Recipients is the list of recipients of the email.
	Recipients []string `json:"Recipients" binding:"required"`

	// CC is the list of CC-recipients of the email.
	CC []string `json:"CC,omitempty"`

	// BCC is the list of BCC-recipients of the email
	BCC []string `json:"BCC,omitempty"`

	// Webhook is the URl for receiving the status of the request.
	Webhook string `json:"Webhook"`

	// CustomID is another ID of the request (optional), for client-side tracking purposes. If provided, it must be a valid 64-character long hex string.
	CustomID string `json:"CustomID,omitempty"`
}

// APISendEmailWithTemplateRequest is a request for a
type APISendEmailWithTemplateRequest struct {
	BaseAPISendEmailRequest

	// TemplateName is the name of the template.
	TemplateName string `json:"TemplateName" binding:"required"`

	// Args consists of arguments for the given template. Should this field not be empty, TemplateName must not be empty.
	Args map[string]interface{} `json:"Args,omitempty"`
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

//type APISendEmailRequest struct {
//	// Subject is the subject of the email.
//	Subject string `json:"Subject,omitempty"`
//
//	// TemplateName is the name of the template.
//	TemplateName string `json:"TemplateName" binding:"required_without=PlainMessage,omitempty"`
//
//	// Args consists of arguments for the given template. Should this field not be empty, TemplateName must not be empty.
//	Args map[string]interface{} `json:"Args" binding:"required_with=TemplateName"`
//
//	// PlainMessage is the plain message of the email.
//	PlainMessage string `json:"PlainMessage" binding:"required_without=TemplateName,omitempty"`
//
//	// Recipients is the list of recipients of the email.
//	Recipients []string `json:"Recipients" binding:"required"`
//
//	// CC is the list of CC-recipients of the email.
//	CC []string `json:"CC"`
//
//	// BCC is the list of BCC-recipients of the email
//	BCC []string `json:"BCC"`
//}

type HTTPClientConfig struct {
	Endpoint    string
	AccessToken string
}

type HTTPClient struct {
	*http.Client
	cfg *HTTPClientConfig
}

func NewHTTPClient(cfg *HTTPClientConfig) *HTTPClient {
	return &HTTPClient{
		Client: &http.Client{},
		cfg:    cfg,
	}
}

func NewNotifyClient(endpoint string, AccessToken string) *HTTPClient {
	cfg := &HTTPClientConfig{endpoint, AccessToken} //"127.0.0.1:12345"}
	client := NewHTTPClient(cfg)

	return client
}

func (c *HTTPClient) RequestSendEmail(message APISendPlainEmailRequest) error {

	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	msgBuffer := bytes.NewBuffer(msgBytes)

	req, err := http.NewRequest("POST", c.parseUrl("msg/send/email"), msgBuffer)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.cfg.AccessToken)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		body, err := c.parseResponse(resp)
		if err != nil {
			return fmt.Errorf("request failed: %v", resp)
		}
		log.LG.Infof("resp: %v, body: %v", resp.Status, string(body))
		return fmt.Errorf("request failed with status: %v", resp.Status)
	}

	_, err = c.parseResponse(resp)
	if err != nil {
		return fmt.Errorf("parseReponse error: %v", err)
	}

	return nil
}

func (c *HTTPClient) RequestSendEmailWithTemplate(message APISendEmailWithTemplateRequest) error {

	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	msgBuffer := bytes.NewBuffer(msgBytes)

	req, err := http.NewRequest("POST", c.parseUrl("msg/send/email/template"), msgBuffer)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.cfg.AccessToken)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		body, err := c.parseResponse(resp)
		if err != nil {
			return fmt.Errorf("request failed: %v", resp)
		}
		log.LG.Infof("resp: %v, body: %v", resp.Status, string(body))
		return fmt.Errorf("request failed with status: %v", resp.Status)
	}

	_, err = c.parseResponse(resp)
	if err != nil {
		return fmt.Errorf("parseReponse error: %v", err)
	}

	return nil
}

func (c *HTTPClient) buildAndDoRequest(message interface{}, method, path string, ret interface{}) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	msgBuffer := bytes.NewBuffer(msgBytes)

	req, err := http.NewRequest(method, c.parseUrl(path), msgBuffer)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.cfg.AccessToken)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		body, err := c.parseResponse(resp)
		if err != nil {
			return fmt.Errorf("request failed: %v", resp)
		}
		log.LG.Infof("resp: %v, body: %v", resp.Status, string(body))
		return fmt.Errorf("request failed with status: %v", resp.Status)
	}

	err = c.parseResponseToRet(resp, ret)
	if err != nil {
		return fmt.Errorf("parseResponseToRet error: %v", err)
	}

	return nil
}

func (c *HTTPClient) parseResponseToRet(resp *http.Response, ret interface{}) error {
	body, err := c.parseResponse(resp)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return err
	}

	return nil
}

func (c *HTTPClient) SendEmailSMSTemplate(message interface{}, sendtype string) (string, error) {
	var ret struct {
		Data response.APISendResponse `json:"data"`
	}

	var api string
	if sendtype == "EMAIL" {
		api = "msg/send/email/template"
	} else if sendtype == "SMS" {
		api = "msg/send/sms"
	}

	err := c.buildAndDoRequest(message, http.MethodPost, api, &ret)
	if err != nil {
		return "", errors.Errorf("message: %v, err: %v", message, err)
	}
	log.LG.Infof("New Notify Request: %v", ret.Data.ID)

	return ret.Data.ID, nil
}

func (c *HTTPClient) SendEmailSMSTemplateAndCheckStatus(message interface{}, sendtype string) error {
	var ret struct {
		Data response.APISendResponse `json:"data"`
	}

	var api string
	if sendtype == "EMAIL" {
		api = "msg/send/email/template"
	} else if sendtype == "SMS" {
		api = "msg/send/sms"
	}

	err := c.buildAndDoRequest(message, http.MethodPost, api, &ret)
	if err != nil {
		return errors.Errorf("message: %v, err: %v", message, err)
	}
	log.LG.Infof("New Notify Request: %v", ret.Data.ID)

	var tmpRet struct {
		Data response.APISendStatusResponse `json:"data"`
	}

	msg := request.APISendStatusRequest{ID: ret.Data.ID}

	ticker := time.NewTicker(10 * time.Second)
	maxFailed := 10
	numFailed := 0

	for range ticker.C {
		err = c.buildAndDoRequest(msg, http.MethodPost, "msg/status", &tmpRet)

		if err == nil {
			switch tmpRet.Data.Status {
			case "DELIVERED":
				log.LG.Infof("Req: %v %v", ret.Data.ID, tmpRet.Data.Status)
				return nil
			case "NOT FOUND", "FAILED", "REJECTED":
				return errors.Errorf("request %v FAILED", ret.Data.ID)
			default:
				//QUEUED
			}
		}

		numFailed += 1
		if numFailed >= maxFailed {
			return errors.Errorf("request %v FAILED: %v", ret.Data.ID)
		}
	}

	return nil
}

func (c *HTTPClient) parseResponse(resp *http.Response) ([]byte, error) {
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			errors.Errorf("Close body error: %v", err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func (c *HTTPClient) parseUrl(shortPath string) string {
	return fmt.Sprintf("%v/%v", c.cfg.Endpoint, shortPath)
}

func (c *HTTPClient) RequestSendSMS(message APISendSMSRequest) error {

	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	msgBuffer := bytes.NewBuffer(msgBytes)

	req, err := http.NewRequest("POST", c.parseUrl("msg/send/sms"), msgBuffer)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.cfg.AccessToken)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		body, err := c.parseResponse(resp)
		if err != nil {
			return fmt.Errorf("request failed: %v", resp)
		}
		log.LG.Infof("resp: %v, body: %v", resp.Status, string(body))
		return fmt.Errorf("request failed with status: %v", resp.Status)
	}

	_, err = c.parseResponse(resp)
	if err != nil {
		return fmt.Errorf("parseReponse error: %v", err)
	}

	return nil
}

func (c *HTTPClient) PushNotiMessageToApp(message *request.APISendPushRequest) error {

	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	msgBuffer := bytes.NewBuffer(msgBytes)

	req, err := http.NewRequest("POST", c.parseUrl("msg/send/push"), msgBuffer)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.cfg.AccessToken)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		body, err := c.parseResponse(resp)
		if err != nil {
			return fmt.Errorf("request failed: %v", resp)
		}
		log.LG.Infof("resp: %v, body: %v", resp.Status, string(body))
		return fmt.Errorf("request failed with status: %v", resp.Status)
	}

	body, err := c.parseResponse(resp)
	if err != nil {
		return fmt.Errorf("parseReponse error: %v", err)
	}

	log.LG.Infof("PushAppMsg : %v, body: %v", resp.Status, string(body))

	return nil
}
