package webhook

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"github.com/labstack/echo"
	"github.com/motoki317/gitea-webhook/model"
	"log"
	"net/http"
	"strings"
)

// postMessage Webhookにメッセージを投稿します
func postMessage(c echo.Context, message string) error {
	url := "https://q.trap.jp/api/1.0/webhooks/" + TraqWebhookId
	req, err := http.NewRequest("POST",
		url,
		strings.NewReader(message))
	if err != nil {
		log.Printf("Error occured while creating a new request: %s\n", err)
		return err
	}

	req.Header.Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
	req.Header.Set("X-TRAQ-Signature", generateSignature(message))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	response := make([]byte, 512)
	_, err = resp.Body.Read(response)
	if err != nil {
		log.Printf("Error occured while reading response from traq webhook: %s\n", err)
	}

	log.Printf("Message sent to %s, message: %s, response: %s\n", url, message, response)

	return c.NoContent(http.StatusNoContent)
}

func generateSignature(message string) string {
	mac := hmac.New(sha1.New, []byte(TraqWebhookSecret))
	_, _ = mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func getAssigneeNames(payload interface{}) (ret string) {
	var assignees []model.User
	switch payload.(type) {
	case model.IssueEvent:
		assignees = payload.(model.IssueEvent).Issue.Assignees
	case model.PullRequestEvent:
		assignees = payload.(model.PullRequestEvent).PullRequest.Assignees
	default:
		return
	}

	if assignees == nil {
		return
	}

	for i, v := range assignees {
		ret += "`" + v.Username + "`"
		if i != len(assignees)-1 {
			ret += ", "
		}
	}
	return
}

func getLabelNames(payload interface{}) (ret string) {
	var labels []model.Label
	switch payload.(type) {
	case model.IssueEvent:
		labels = payload.(model.IssueEvent).Issue.Labels
	case model.PullRequestEvent:
		labels = payload.(model.PullRequestEvent).PullRequest.Labels
	default:
		return
	}

	if labels == nil {
		return
	}

	for i, v := range labels {
		ret += ":0x" + v.Color + ": `" + v.Name + "`"
		if i != len(labels)-1 {
			ret += ", "
		}
	}
	return ret
}
