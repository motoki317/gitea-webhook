package webhook

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/motoki317/gitea-webhook/model"
	"log"
	"net/http"
	"os"
)

var (
	TraqWebhookId     = os.Getenv("TRAQ_WEBHOOK_ID")
	TraqWebhookSecret = os.Getenv("TRAQ_WEBHOOK_SECRET")
)

func MakeWebhookHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		event := c.Request().Header.Get("X-Gitea-Event")
		log.Printf("Received event %s\n", event)

		switch event {
		case "issues":
			return handleIssuesEvent(c)
		case "issue_comment":
			return handleIssueCommentEvent(c)
		case "pull_request":
			return handlePullRequestEvent(c)
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func handleIssuesEvent(c echo.Context) error {
	payload := model.IssueEvent{}
	if err := c.Bind(&payload); err != nil {
		log.Printf("Error occured while binding payload: %s\n", err)
		return err
	}

	log.Printf("Issue event action: %s\n", payload.Action)

	senderName := payload.Sender.Username
	message := fmt.Sprintf("### Issue [#%v %s](%s) ", payload.Issue.Number, payload.Issue.Title, payload.Issue.URL)

	switch payload.Action {
	case "opened":
		message += fmt.Sprintf("Opened by `%s`\n", senderName)
	case "edited":
		message += fmt.Sprintf("Edited by `%s`\n", senderName)
	case "assigned":
		message += fmt.Sprintf("Assigned to `%s`\n", payload.Issue.Assignee.Username)
		message += fmt.Sprintf("By `%s`\n", senderName)
		message += fmt.Sprintf("Assignees: %s\n", getAssigneeNames(payload))
	case "unassigned":
		message += fmt.Sprintf("Unassigned\n")
		message += fmt.Sprintf("By `%s`\n", senderName)
		message += fmt.Sprintf("Assignees: %s\n", getAssigneeNames(payload))
	case "label_updated":
		message += fmt.Sprintf("Label Updated\n")
		message += fmt.Sprintf("By `%s`\n", senderName)
		message += fmt.Sprintf("Labels: %s\n", getLabelNames(payload))
	case "milestoned":
		message += fmt.Sprintf("Milestone Set by `%s`\n", senderName)
		message += fmt.Sprintf("Milestone `%s` due to %s\n", payload.Issue.Milestone.Title, payload.Issue.Milestone.DueOn)
	case "demilestoned":
		message += fmt.Sprintf("Milestone Removed by `%s`\n", senderName)
	case "closed":
		message += fmt.Sprintf("Closed by `%s`\n", senderName)
	case "reopened":
		message += fmt.Sprintf("Reopened by `%s`\n", senderName)
	}

	message += fmt.Sprintf("\n---\n")
	message += fmt.Sprintf("%s", payload.Issue.Body)

	return postMessage(c, message)
}

func handleIssueCommentEvent(c echo.Context) error {
	payload := model.IssueCommentEvent{}
	if err := c.Bind(&payload); err != nil {
		log.Printf("Error occured while binding payload: %s\n", err)
		return err
	}

	senderName := payload.Sender.Username
	issueName := fmt.Sprintf("[#%v %s](%s)", payload.Issue.Number, payload.Issue.Title, payload.Issue.URL)
	message := "### "

	switch payload.Action {
	case "created":
		message += "New Comment"
	case "edited":
		message += "Comment Edited"
	case "deleted":
		message += "Comment Deleted"
	}

	message += fmt.Sprintf(" by `%s`\n", senderName)
	message += fmt.Sprintf("Issue or Pull Request: %s\n", issueName)
	message += fmt.Sprintf("\n---\n")
	message += fmt.Sprintf("%s", payload.Comment.Body)

	return postMessage(c, message)
}

func handlePullRequestEvent(c echo.Context) error {
	payload := model.PullRequestEvent{}
	if err := c.Bind(&payload); err != nil {
		log.Printf("Error occured while binding payload: %s\n", err)
		return err
	}

	senderName := payload.Sender.Username
	message := fmt.Sprintf("### Pull Request [#%v %s](%s) ", payload.PullRequest.Number, payload.PullRequest.Title, payload.PullRequest.URL)

	switch payload.Action {
	case "opened":
		message += fmt.Sprintf("Opened by `%s`\n", senderName)
	case "synchronized":
		message += fmt.Sprintf("Changed by `%s`\n", senderName)
	case "assigned":
		message += fmt.Sprintf("Assigned to `%s`\n", payload.PullRequest.Assignee.Username)
		message += fmt.Sprintf("By `%s`\n", senderName)
		message += fmt.Sprintf("Assignees: %s\n", getAssigneeNames(payload))
	case "unassigned":
		message += fmt.Sprintf("Unassigned\n")
		message += fmt.Sprintf("By `%s`\n", senderName)
		message += fmt.Sprintf("Assignees: %s\n", getAssigneeNames(payload))
	case "milestoned":
		message += fmt.Sprintf("Milestone Set by `%s`\n", senderName)
		message += fmt.Sprintf("Milestone `%s` due to %s\n", payload.PullRequest.Milestone.Title, payload.PullRequest.Milestone.DueOn)
	case "demilestoned":
		message += fmt.Sprintf("Milestone Removed by `%s`\n", senderName)
	case "label_updated":
		message += fmt.Sprintf("Label Updated\n")
		message += fmt.Sprintf("By `%s`\n", senderName)
		message += fmt.Sprintf("Labels: %s\n", getLabelNames(payload))
	case "closed":
		switch payload.PullRequest.Merged {
		case true:
			message += fmt.Sprintf("Merged by `%s`\n", senderName)
		case false:
			message += fmt.Sprintf("Closed by `%s`\n", senderName)
		}
	case "reopened":
		message += fmt.Sprintf("Reopened by `%s`\n", senderName)
	}

	message += fmt.Sprintf("\n---\n")
	message += fmt.Sprintf("%s", payload.PullRequest.Body)

	return postMessage(c, message)
}
