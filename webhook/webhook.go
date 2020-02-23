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
	issueName := fmt.Sprintf("Issue [#%v %s](%s) ",
		payload.Issue.Number,
		payload.Issue.Title,
		payload.Repository.HTMLURL+"/issues/"+strconv.Itoa(payload.Issue.Number),
	)
	message := "### "

	switch payload.Action {
	case "opened":
		message += fmt.Sprintf(":git_issue_opened: %s Opened by `%s`\n", issueName, senderName)
	case "edited":
		message += fmt.Sprintf(":pencil: %s Edited by `%s`\n", issueName, senderName)
	case "assigned":
		message += fmt.Sprintf(":bust_in_silhouette: %s Assigned to `%s`\n", issueName, payload.Issue.Assignee.Username)
		message += fmt.Sprintf("By `%s`\n", senderName)
		message += fmt.Sprintf("Assignees: %s\n", getAssigneeNames(payload))
	case "unassigned":
		message += fmt.Sprintf(":bust_in_silhouette: %s Unassigned\n", issueName)
		message += fmt.Sprintf("By `%s`\n", senderName)
		message += fmt.Sprintf("Assignees: %s\n", getAssigneeNames(payload))
	case "label_updated":
		message += fmt.Sprintf(":label: %s Label Updated\n", issueName)
		message += fmt.Sprintf("By `%s`\n", senderName)
		message += fmt.Sprintf("Labels: %s\n", getLabelNames(payload))
	case "milestoned":
		message += fmt.Sprintf(":git_milestone: %s Milestone Set by `%s`\n", issueName, senderName)
		message += fmt.Sprintf("Milestone `%s` due by %s\n", payload.Issue.Milestone.Title, payload.Issue.Milestone.DueOn)
	case "demilestoned":
		message += fmt.Sprintf(":git_milestone: %s Milestone Removed by `%s`\n", issueName, senderName)
	case "closed":
		message += fmt.Sprintf(":git_issue_closed: %s Closed by `%s`\n", issueName, senderName)
	case "reopened":
		message += fmt.Sprintf(":git_issue_opened: %s Reopened by `%s`\n", issueName, senderName)
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
		message += ":comment: New Comment"
	case "edited":
		message += ":pencil: Comment Edited"
	case "deleted":
		message += ":pencil: Comment Deleted"
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
		message += fmt.Sprintf(":git_pull_request: %s Opened by `%s`\n", prName, senderName)
	case "edited":
		message += fmt.Sprintf(":pencil: %s Edited by `%s`\n", prName, senderName)
	case "synchronized":
		message += fmt.Sprintf(":git_push_repo: New Commit(s) to %s by `%s`\n", prName, senderName)
	case "assigned":
		message += fmt.Sprintf(":bust_in_silhouette: %s Assigned to `%s`\n", prName, payload.PullRequest.Assignee.Username)
		message += fmt.Sprintf("By `%s`\n", senderName)
		message += fmt.Sprintf("Assignees: %s\n", getAssigneeNames(payload))
	case "unassigned":
		message += fmt.Sprintf(":bust_in_silhouette: %s Unassigned\n", prName)
		message += fmt.Sprintf("By `%s`\n", senderName)
		message += fmt.Sprintf("Assignees: %s\n", getAssigneeNames(payload))
	case "milestoned":
		message += fmt.Sprintf(":git_milestone: %s Milestone Set by `%s`\n", prName, senderName)
		message += fmt.Sprintf("Milestone `%s` due to %s\n", payload.PullRequest.Milestone.Title, payload.PullRequest.Milestone.DueOn)
	case "demilestoned":
		message += fmt.Sprintf(":git_milestone: %s Milestone Removed by `%s`\n", prName, senderName)
	case "label_updated":
		message += fmt.Sprintf(":label: %s Label Updated\n", prName)
		message += fmt.Sprintf("By `%s`\n", senderName)
		message += fmt.Sprintf("Labels: %s\n", getLabelNames(payload))
	case "closed":
		switch payload.PullRequest.Merged {
		case true:
			message += fmt.Sprintf(":git_merged: %s Merged by `%s`\n", prName, senderName)
		case false:
			message += fmt.Sprintf(":git_pull_request_closed: %s Closed by `%s`\n", prName, senderName)
		}
	case "reopened":
		message += fmt.Sprintf(":git_pull_request: %s Reopened by `%s`\n", prName, senderName)
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
		message += fmt.Sprintf(":white_check_mark: %s Approved by `%s`", prName, senderName)
	case "comment":
		message += fmt.Sprintf(":comment: %s New Review Comment by `%s`", prName)
	case "rejected":
		message += fmt.Sprintf(":comment: %s Changes Requested by `%s`", prName, senderName)
	}

	return postMessage(c, message)
}
