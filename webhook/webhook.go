package webhook

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/motoki317/gitea-webhook/model"
	"log"
	"net/http"
	"os"
	"strconv"
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
		case "pull_request_approved":
			return handlePullRequestReviewEvent(c, "approved")
		case "pull_request_comment":
			return handlePullRequestReviewEvent(c, "comment")
		case "pull_request_rejected":
			return handlePullRequestReviewEvent(c, "rejected")
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
	message := fmt.Sprintf("### Issue [#%v %s](%s) ",
		payload.Issue.Number,
		payload.Issue.Title,
		payload.Repository.HTMLURL+"/issues/"+strconv.Itoa(payload.Issue.Number),
	)

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
		message += fmt.Sprintf("Milestone `%s` due by %s\n", payload.Issue.Milestone.Title, payload.Issue.Milestone.DueOn)
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
	issueName := fmt.Sprintf("[#%v %s](%s)",
		payload.Issue.Number,
		payload.Issue.Title,
		payload.Repository.HTMLURL+"/issues/"+strconv.Itoa(payload.Issue.Number),
	)
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
	message += fmt.Sprintf("%s\n", issueName)
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
	message := "### "
	prName := fmt.Sprintf("Pull Request [#%v %s](%s)", payload.PullRequest.Number, payload.PullRequest.Title, payload.PullRequest.HTMLURL)

	switch payload.Action {
	case "opened":
		message += fmt.Sprintf("%s Opened by `%s`\n", prName, senderName)
	case "edited":
		message += fmt.Sprintf("%s Edited by `%s`\n", prName, senderName)
	case "synchronized":
		message += fmt.Sprintf("New Commit(s) to %s by `%s`\n", prName, senderName)
	case "assigned":
		message += fmt.Sprintf("%s Assigned to `%s`\n", prName, payload.PullRequest.Assignee.Username)
		message += fmt.Sprintf("By `%s`\n", senderName)
		message += fmt.Sprintf("Assignees: %s\n", getAssigneeNames(payload))
	case "unassigned":
		message += fmt.Sprintf("%s Unassigned\n", prName)
		message += fmt.Sprintf("By `%s`\n", senderName)
		message += fmt.Sprintf("Assignees: %s\n", getAssigneeNames(payload))
	case "milestoned":
		message += fmt.Sprintf("%s Milestone Set by `%s`\n", prName, senderName)
		message += fmt.Sprintf("Milestone `%s` due to %s\n", payload.PullRequest.Milestone.Title, payload.PullRequest.Milestone.DueOn)
	case "demilestoned":
		message += fmt.Sprintf("%s Milestone Removed by `%s`\n", prName, senderName)
	case "label_updated":
		message += fmt.Sprintf("%s Label Updated\n", prName)
		message += fmt.Sprintf("By `%s`\n", senderName)
		message += fmt.Sprintf("Labels: %s\n", getLabelNames(payload))
	case "closed":
		switch payload.PullRequest.Merged {
		case true:
			message += fmt.Sprintf("%s Merged by `%s`\n", prName, senderName)
		case false:
			message += fmt.Sprintf("%s Closed by `%s`\n", prName, senderName)
		}
	case "reopened":
		message += fmt.Sprintf("%s Reopened by `%s`\n", prName, senderName)
	}

	message += fmt.Sprintf("\n---\n")
	message += fmt.Sprintf("%s", payload.PullRequest.Body)

	return postMessage(c, message)
}

func handlePullRequestReviewEvent(c echo.Context, status string) error {
	payload := model.PullRequestEvent{}
	if err := c.Bind(&payload); err != nil {
		log.Printf("Error occured while binding payload: %s\n", err)
		return err
	}

	senderName := payload.Sender.Username
	message := "### "
	prName := fmt.Sprintf("Pull Request [#%v %s](%s)", payload.PullRequest.Number, payload.PullRequest.Title, payload.PullRequest.HTMLURL)
	switch status {
	case "approved":
		message += fmt.Sprintf("%s Approved by `%s`", prName, senderName)
	case "comment":
		message += fmt.Sprintf("New Review Comment to %s", prName)
	case "rejected":
		message += fmt.Sprintf("%s Rejected by `%s`", prName, senderName)
	}

	return postMessage(c, message)
}
