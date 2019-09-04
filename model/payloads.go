package model

import "time"

// Common structs

type User struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Language  string `json:"language"`
	Username  string `json:"username"`
}

type Repository struct {
	ID              int         `json:"id"`
	Owner           User        `json:"owner"`
	Name            string      `json:"name"`
	FullName        string      `json:"full_name"`
	Description     string      `json:"description"`
	Empty           bool        `json:"empty"`
	Private         bool        `json:"private"`
	Fork            bool        `json:"fork"`
	Parent          interface{} `json:"parent"`
	Mirror          bool        `json:"mirror"`
	Size            int         `json:"size"`
	HTMLURL         string      `json:"html_url"`
	SSHURL          string      `json:"ssh_url"`
	CloneURL        string      `json:"clone_url"`
	Website         string      `json:"website"`
	StarsCount      int         `json:"stars_count"`
	ForksCount      int         `json:"forks_count"`
	WatchersCount   int         `json:"watchers_count"`
	OpenIssuesCount int         `json:"open_issues_count"`
	DefaultBranch   string      `json:"default_branch"`
	Archived        bool        `json:"archived"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	Permissions     struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permissions"`
}

type Issue struct {
	ID          int         `json:"id"`
	URL         string      `json:"url"`
	Number      int         `json:"number"`
	User        User        `json:"user"`
	Title       string      `json:"title"`
	Body        string      `json:"body"`
	Labels      []Label     `json:"labels"`
	Milestone   Milestone   `json:"milestone"`
	Assignee    User        `json:"assignee"`
	Assignees   []User      `json:"assignees"`
	State       string      `json:"state"`
	Comments    int         `json:"comments"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	ClosedAt    time.Time   `json:"closed_at"`
	DueDate     interface{} `json:"due_date"`
	PullRequest interface{} `json:"pull_request"`
}

type Label struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	URL   string `json:"url"`
}

type Milestone struct {
	ID           int         `json:"id"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	State        string      `json:"state"`
	OpenIssues   int         `json:"open_issues"`
	ClosedIssues int         `json:"closed_issues"`
	ClosedAt     interface{} `json:"closed_at"`
	DueOn        time.Time   `json:"due_on"`
}

type Comment struct {
	ID             int       `json:"id"`
	HTMLURL        string    `json:"html_url"`
	PullRequestURL string    `json:"pull_request_url"`
	IssueURL       string    `json:"issue_url"`
	User           User      `json:"user"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PullRequest struct {
	ID             int         `json:"id"`
	URL            string      `json:"url"`
	Number         int         `json:"number"`
	User           User        `json:"user"`
	Title          string      `json:"title"`
	Body           string      `json:"body"`
	Labels         []Label     `json:"labels"`
	Milestone      Milestone   `json:"milestone"`
	Assignee       User        `json:"assignee"`
	Assignees      []User      `json:"assignees"`
	State          string      `json:"state"`
	Comments       int         `json:"comments"`
	HTMLURL        string      `json:"html_url"`
	DiffURL        string      `json:"diff_url"`
	PatchURL       string      `json:"patch_url"`
	Mergeable      bool        `json:"mergeable"`
	Merged         bool        `json:"merged"`
	MergedAt       interface{} `json:"merged_at"`
	MergeCommitSha interface{} `json:"merge_commit_sha"`
	MergedBy       interface{} `json:"merged_by"`
	Base           struct {
		Label  string     `json:"label"`
		Ref    string     `json:"ref"`
		Sha    string     `json:"sha"`
		RepoID int        `json:"repo_id"`
		Repo   Repository `json:"repo"`
	} `json:"base"`
	Head struct {
		Label  string     `json:"label"`
		Ref    string     `json:"ref"`
		Sha    string     `json:"sha"`
		RepoID int        `json:"repo_id"`
		Repo   Repository `json:"repo"`
	} `json:"head"`
	MergeBase string      `json:"merge_base"`
	DueDate   interface{} `json:"due_date"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	ClosedAt  interface{} `json:"closed_at"`
}

// ---------- Events ----------

type IssueEvent struct {
	Secret     string     `json:"secret"`
	Action     string     `json:"action"`
	Number     int        `json:"number"`
	Issue      Issue      `json:"issue"`
	Repository Repository `json:"repository"`
	Sender     User       `json:"sender"`
}

type IssueCommentEvent struct {
	Secret     string     `json:"secret"`
	Action     string     `json:"action"`
	Issue      Issue      `json:"issue"`
	Comment    Comment    `json:"comment"`
	Repository Repository `json:"repository"`
	Sender     User       `json:"sender"`
}

type PullRequestEvent struct {
	Secret      string      `json:"secret"`
	Action      string      `json:"action"`
	Number      int         `json:"number"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  Repository  `json:"repository"`
	Sender      User        `json:"sender"`
}
