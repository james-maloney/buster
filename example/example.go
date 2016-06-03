package main

import (
	"html/template"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/james-maloney/buster"
)

var tmpl = template.Must(template.New("base").Parse(indexTmpl))

func main() {
	e := gin.New()

	fs := buster.NewFileServer("./files", "/assets/")
	e.Use(fs.GinFunc())

	e.GET("/", func(ctx *gin.Context) {
		tmpl.Execute(ctx.Writer, map[string]interface{}{
			"Style": fs.BuildURL("/assets/style.css"),
			"JS":    fs.BuildURL("/assets/page.js"),
		})
	})

	log.Fatal(e.Run(":8080"))
}

var indexTmpl = `
<!DOCTYPE html>
<html>
	<head>
		<link href="{{ .Style }}" rel="stylesheet" type="text/css" />
		<script src="{{ .JS }}"></script>
	</head>
	<body>
		<h1>Hello, World!</h2>
	</body>
</html>
`
