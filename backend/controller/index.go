package controller

import (
	"fmt"
	"net/http"

	"team/model"
	"team/web"
)

const page = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<title>%s</title>
		<link rel="shortcut icon" href="/www/favicon.ico" />
	</head>
	<body>
		<div id="%s"></div>
		<script src="/www/app.js"></script>
	</body>
</html>
`

// Index handler.
func Index(c *web.Context) {
	container := "install"
	if model.Environment.Installed {
		container = "app"
	}

	c.HTML(http.StatusOK, fmt.Sprintf(page, model.Environment.AppName, container))
}
