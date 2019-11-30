package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const tpl = `<html>
	<head>
		<meta name="go-import"
		      content="{{.Domain}}
                   git https://{{.CodePath}}/{{.RepoPath}}">
		<meta name="go-source"
		      content="{{.Domain}}
                   https://{{.CodePath}}/{{.RepoPath}}
                   https://{{.CodePath}}/{{.RepoPath}}/tree/master{/dir}
                   https://{{.CodePath}}/{{.RepoPath}}/blob/master{/dir}/{file}#L{line}">
		<meta http-equiv="refresh" content="1; url=https://godoc.org/{{.Domain}}{{.Path}}/">
	</head>
	<body>
		Nothing to see here; <a href="https://godoc.org/{{.Domain}}{{.Path}}/">see the package on godoc</a>.
	</body>
</html>`

type goHTMLData struct {
	Domain   string
	CodePath string
	Path     string
	RepoPath string
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	pathsplit := strings.Split(request.Path, "/")
	data := goHTMLData{
		Domain:   os.Getenv("DOMAIN"),
		CodePath: os.Getenv("CODEPATH"),
		Path:     request.Path,
		RepoPath: pathsplit[1],
	}

	t, err := template.New("index").Parse(tpl)
	if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("Couldn't construct template")
	}

	var resp bytes.Buffer
	err = t.Execute(&resp, data)
	if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("Couldn't construct gopath %s", err.Error())
	}

	return events.APIGatewayProxyResponse{
		Body:       resp.String(),
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":  "text/html",
			"Cache-Control": "public, max-age=86400",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
