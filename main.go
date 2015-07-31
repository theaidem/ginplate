package main

import (
	"theaidem/ginplate/backend/handlers"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.HTMLRender = createMyRender()

	router.Static("/static", "./static/")
	router.Static("/bower", "./static/bower")

	router.GET("/", handlers.Home)

	endless.ListenAndServe(":8080", router)
}

func createMyRender() multitemplate.Render {
	render := multitemplate.New()

	root := "templates/"
	base := root + "base/"
	partials := root + "partials/"

	layout := []string{
		root + "layout.html",
		partials + "header.html",
		root + "styles.html",
		partials + "footer.html",
		root + "scripts.html",
	}

	render.AddFromFiles("home", append(layout, base+"home.html")...)

	return render
}
