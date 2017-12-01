## Sanctify

Convert JSON data into a Go struct with idiomatic field names
  
## Synopsis

`go get -u github.com/as/sanctify/...`

`sanctify < data.json`
  
## Description

Much of the work in creating a marshaller is in the tweaks that occur after the specification is
copied into Go source code. Because Go allows structs to have tags, you can keep the bad names in
an already-established API specification out of the struct field names. This package is an executable
that reads JSON data and converts it into Go, but also takes the labor and guess work out of such decisions. 

It uses the following process:
	
- Marshal JSON into a Go interface{}
- Recursively descend into arrays, amalagating fields of underlying JSON objects into a set
- Parse the amalagate tree, generating basic Go source in a main package
- Vet the package on the fly with golint, capturing naming suggestions in a buffer
- Remove underscores in variable names
- Capitalize letter occupying position of deleted underscores
- Compile rules to correct improperly-punctuated acronyms in struct field names
- Reparse the amalagate, applying corrections during the recursive descent step
- Prefix compare child fields to parent fields, remove stuttercase naming in Go field names
- Run gofmt -s to simplify code

Then you:

- Check the auto-generated code for consistency
- Paste it into your project

## Options
	
	-p    Package name to generate (default: main)
	-t    JSON root element type (default: X)
	-o    Add omitempty to all fields
  

## Why

- Generates Go without adding third party dependencies.
- Save time correcting huge datasets by hand

## Example
 
`echo {"msg":{"msg_string":"hi","msg_num": 3}} | sanctify`
   
```
package main
type X struct {
  Msg struct {
    String string `json:"msg_string"`
    Num    int    `json:"msg_num"`
  } `json:"msg"`
}
```

## Example 2

This hideous thing is straight from the json.org website. 

```
{"menu": {
    "header": "SVG Viewer",
    "items": [
        {"id": "Open"},
        {"id": "OpenNew", "label": "Open New"},
        null,
        {"id": "ZoomIn", "label": "Zoom In"},
        {"id": "ZoomOut", "label": "Zoom Out"},
        {"id": "OriginalView", "label": "Original View"},
        null,
        {"id": "Quality"},
        {"id": "Pause"},
        {"id": "Mute"},
        null,
        {"id": "Find", "label": "Find..."},
        {"id": "FindAgain", "label": "Find Again"},
        {"id": "Copy"},
        {"id": "CopyAgain", "label": "Copy Again"},
        {"id": "CopySVG", "label": "Copy SVG"},
        {"id": "ViewSVG", "label": "View SVG"},
        {"id": "ViewSource", "label": "View Source"},
        {"id": "SaveAs", "label": "Save As"},
        null,
        {"id": "Help"},
        {"id": "About", "label": "About Adobe CVG Viewer..."}
    ]
}}
```

``` cat horridmonkey.json | sanctify```

```
package main

type X struct {
        Menu struct {
                Header string `json:"header"`
                Items  []struct {
                        Id    string `json:"id"`
                        Label string `json:"label"`
                } `json:"items"`
        } `json:"menu"`
}
```

# Example 3

```hget https://api.github.com/repos/as/sanctify | sanctify```

```
type X struct {
        IssueEventsURL   string      `json:"issue_events_url"`
        BlobsURL         string      `json:"blobs_url"`
        CloneURL         string      `json:"clone_url"`
        Description      string      `json:"description"`
        CollaboratorsURL string      `json:"collaborators_url"`
        TeamsURL         string      `json:"teams_url"`
        ArchiveURL       string      `json:"archive_url"`
        PushedAt         string      `json:"pushed_at"`
        SubscribersCount int         `json:"subscribers_count"`
        Watchers         int         `json:"watchers"`
        Private          bool        `json:"private"`
        HooksURL         string      `json:"hooks_url"`
        LanguagesURL     string      `json:"languages_url"`
        SubscribersURL   string      `json:"subscribers_url"`
        HasPages         bool        `json:"has_pages"`
        MirrorURL        interface{} `json:"mirror_url"`
        Archived         bool        `json:"archived"`
        Id               int         `json:"id"`
        PullsURL         string      `json:"pulls_url"`
        NotificationsURL string      `json:"notifications_url"`
        SvnURL           string      `json:"svn_url"`
        NetworkCount     int         `json:"network_count"`
        IssuesURL        string      `json:"issues_url"`
        TreesURL         string      `json:"trees_url"`
        IssueCommentURL  string      `json:"issue_comment_url"`
        SshURL           string      `json:"ssh_url"`
        GitCommitsURL    string      `json:"git_commits_url"`
        CommentsURL      string      `json:"comments_url"`
        Size             int         `json:"size"`
        WatchersCount    int         `json:"watchers_count"`
        Url              string      `json:"url"`
        KeysURL          string      `json:"keys_url"`
        EventsURL        string      `json:"events_url"`
        TagsURL          string      `json:"tags_url"`
        Language         string      `json:"language"`
        Owner            struct {
                AvatarURL         string `json:"avatar_url"`
                GravatarID        string `json:"gravatar_id"`
                StarredURL        string `json:"starred_url"`
                ReposURL          string `json:"repos_url"`
                FollowersURL      string `json:"followers_url"`
                SubscriptionsURL  string `json:"subscriptions_url"`
                Type              string `json:"type"`
                SiteAdmin         bool   `json:"site_admin"`
                Id                int    `json:"id"`
                GistsURL          string `json:"gists_url"`
                EventsURL         string `json:"events_url"`
                Login             string `json:"login"`
                Url               string `json:"url"`
                HtmlURL           string `json:"html_url"`
                FollowingURL      string `json:"following_url"`
                OrganizationsURL  string `json:"organizations_url"`
                ReceivedEventsURL string `json:"received_events_url"`
        } `json:"owner"`
        GitTagsURL      string `json:"git_tags_url"`
        StargazersURL   string `json:"stargazers_url"`
        ForksCount      int    `json:"forks_count"`
        GitRefsURL      string `json:"git_refs_url"`
        CommitsURL      string `json:"commits_url"`
        ContentsURL     string `json:"contents_url"`
        CompareURL      string `json:"compare_url"`
        DownloadsURL    string `json:"downloads_url"`
        DeploymentsURL  string `json:"deployments_url"`
        HasIssues       bool   `json:"has_issues"`
        FullName        string `json:"full_name"`
        UpdatedAt       string `json:"updated_at"`
        HasProjects     bool   `json:"has_projects"`
        SubscriptionURL string `json:"subscription_url"`
        ReleasesURL     string `json:"releases_url"`
        HasDownloads    bool   `json:"has_downloads"`
        HtmlURL         string `json:"html_url"`
        Fork            bool   `json:"fork"`
        CreatedAt       string `json:"created_at"`
        HasWiki         bool   `json:"has_wiki"`
        License         struct {
                Key    string `json:"key"`
                Name   string `json:"name"`
                SpdxID string `json:"spdx_id"`
                Url    string `json:"url"`
        } `json:"license"`
        Name            string      `json:"name"`
        ForksURL        string      `json:"forks_url"`
        AssigneesURL    string      `json:"assignees_url"`
        BranchesURL     string      `json:"branches_url"`
        GitURL          string      `json:"git_url"`
        StatusesURL     string      `json:"statuses_url"`
        ContributorsURL string      `json:"contributors_url"`
        LabelsURL       string      `json:"labels_url"`
        Homepage        interface{} `json:"homepage"`
        OpenIssuesCount int         `json:"open_issues_count"`
        Forks           int         `json:"forks"`
        OpenIssues      int         `json:"open_issues"`
        MergesURL       string      `json:"merges_url"`
        MilestonesURL   string      `json:"milestones_url"`
        StargazersCount int         `json:"stargazers_count"`
        DefaultBranch   string      `json:"default_branch"`
}
```

