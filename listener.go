package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/davecgh/go-spew/spew"
	// "os"
)

func setupListener() {

	httpURI := cfg.Section("http").Key("http_uri").String()
	httpPort := cfg.Section("http").Key("http_port").String()

	mux := http.NewServeMux()
	mux.HandleFunc(httpURI, PostHandler)

	log.Printf("listening on port %s", httpPort)
	log.Fatal(http.ListenAndServe(":"+httpPort, mux))

}

// PostHandler converts post request body to string
func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}

		var postData GiteaPostData

		//TODO check secret!
		//TODO check repo

		json.Unmarshal(body, &postData)

		// fmt.Println(string(body))

		fmt.Printf("%+v\n", postData)
		fmt.Println("=================================================")
		spew.Dump(postData)
		fmt.Println("=================================================")

		fmt.Println("=================================================")
		fmt.Println(generateMessage(postData, r.Header.Get("X-Gitea-Event")))
		fmt.Println("=================================================")
		fmt.Fprint(w, "POST done")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

//generateMessage generates the message string for a given event
func generateMessage(data GiteaPostData, eventHeader string) string {

	templ := template.New("notification")
	var tpl bytes.Buffer

	fmt.Println("======================")

	fmt.Println(eventHeader)
	fmt.Println(data.Action)

	fmt.Println("======================")

	switch eventHeader {

	case "push":
		templ.Parse("{{.Pusher.FullName}} pushed " + strconv.Itoa(len(data.Commits)) + " commit(s) to {{.Repository.Name}}")

	case "issues":
		switch data.Action {
		case "assigned":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} assigned issue #{{.Issue.Number}} {{.Issue.Title}} to  TODO ")
		case "closed":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} closed issue #{{.Issue.Number}} {{.Issue.Title}}")
		case "demilestoned":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} removed milestone TODO from issue #{{.Issue.Number}} {{.Issue.Title}}")
		case "edited":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} edited issue #{{.Issue.Number}} {{.Issue.Title}}")
		case "label_cleared":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} cleared labels from issue #{{.Issue.Number}} {{.Issue.Title}}")
		case "label_updated":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} updated labels of issue #{{.Issue.Number}} {{.Issue.Title}}")
		case "milestoned":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} added issue #{{.Issue.Number}} {{.Issue.Title}} to milestone TODO")
		case "opened":
			templ.Parse("{{.Repository.FullName}}: {{.Issue.User.FullName}} opened issue #{{.Issue.Number}} {{.Issue.Title}}")
		case "reopened":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} re-opened issue #{{.Issue.Number}} {{.Issue.Title}}")
		case "synchronized":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} synchronized issue #{{.Issue.Number}} {{.Issue.Title}}")
		case "unassigned":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} removed assignee from issue #{{.Issue.Number}} {{.Issue.Title}}")
		}

	case "fork":
		// switch data.Action {
		// case "created":
		templ.Parse("{{.Sender.FullName}} forked repository {{.Repository.Parent.FullName}} to {{.Repository.FullName}}")
		// case "edited":
		// templ.Parse("{{.Sender.FullName}} edited repository fork {{}}")
		// case "deleted":
		// templ.Parse("{{.Sender.FullName}} deleted repository fork {{}}")
		// }

	case "pull_request":
		switch data.Action {
		case "assigned":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} assigned pull-request {{}} {{}} to {{}}")
		case "closed":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} closed pull-request {{}} {{}}")
		case "demilestoned":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} removed milestone TODO from pull-request {{}} {{}}")
		case "edited":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} edited pull-request {{}} {{}}")
		case "label_cleared":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} removed labels from pull-request {{}} {{}}")
		case "label_updated":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} updated labels from pull-request {{}} {{}}")
		case "milestoned":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} added pull-request {{}} {{}} to milestone TODO")
		case "opened":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} opened pull-request {{}} {{}}")
		case "reopened":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} re-opened pull-request {{}} {{}}")
		case "synchronized":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} synchronized pull-request {{}} {{}}")
		case "unassigned":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} removed assinee {{}} from pull-request {{}} {{}}")
		}

	case "issue_comment":
		switch data.Action {
		case "created":
			templ.Parse("{{.Sender.FullName}} commented issue #{{.Issue.Number}} {{.Issue.Title}}: {{.Comment.Body}}")
		case "deleted":
			templ.Parse("{{.Sender.FullName}} deleted commented on issue #{{.Issue.Number}} {{.Issue.Title}}")
		}

	case "repository":
		switch data.Action {
		case "created":
			templ.Parse("{{.Sender.FullName}} created repository {{}}")
		case "deleted":
			templ.Parse("{{.Sender.FullName}} deleted repository {{}}")
		}

	case "release":
		switch data.Action {
		case "published":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} published release {{}}")
		case "updated":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} updated release {{}}")
		case "deleted":
			templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} deleted release {{}}")
		}

	case "pull_request_approved":
		templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} approved pull request {{}}{{}}")
	case "pull_request_rejected":
		templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} rejected pull request {{}}{{}}")
	case "pull_request_comment":
		templ.Parse("{{.Repository.FullName}}: {{.Sender.FullName}} commented on pull request {{}}{{}}: {{.Comment.Body}}")
	default:

		fmt.Println("Unknown action: " + eventHeader + " " + data.Action)
		templ.Parse("Gitea did something unexpected, seriously wtf was that?! Event: " + eventHeader + " Action: " + data.Action)

	}

	if err := templ.Execute(&tpl, data); err != nil {
		panic(err)
	}

	return tpl.String()

	// tmplIssue, err := template.New("test").Parse("{{.User}} {{.Action}} Issue {{.IssueID}} in repository {{.Repo}}")

}
