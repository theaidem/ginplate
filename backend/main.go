package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/host"
)

var mode, buildstamp string

func init() {
	gin.SetMode(mode)
}

func main() {

	router := gin.Default()

	if gin.Mode() == gin.DebugMode {
		// For Dev Mode
		router.Static("/js", "frontend/app/js")
		router.Static("/css", "frontend/app/css")
		router.Static("/bower_components", "frontend/app/bower_components")
		router.LoadHTMLGlob("frontend/app/index.html")
	} else if gin.Mode() == gin.ReleaseMode {
		// For Prod Mode
		router.Use(static.Serve("/index", BinaryFileSystem("frontend/dist/index.html")))
		router.Use(static.Serve("/css", BinaryFileSystem("frontend/dist/css")))
		router.Use(static.Serve("/js", BinaryFileSystem("frontend/dist/js")))
	}

	// For SPA Router
	router.NoRoute(index)

	api := router.Group("/api")

	xhr := api.Group("/xhr")

	xhr.Group("/admin").
		GET("/host", hostInfo)

	sse := api.Group("/sse")
	sse.GET("/stream", stream)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
		//Comment timeouts if you use SSE streams
		//ReadTimeout: 10 * time.Second,
		//WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()
}

func index(ctx *gin.Context) {

	if gin.Mode() == gin.DebugMode {

		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "appName",
		})

	} else if gin.Mode() == gin.ReleaseMode {

		fmt.Println(buildstamp, mode)
		templateString, err := Asset("frontend/dist/index.html")
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		t, err := template.New("index").Parse(string(templateString))
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var index bytes.Buffer
		err = t.Execute(&index, gin.H{
			"title": "appName",
		})

		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.Data(http.StatusOK, "text/html; charset=utf-8", index.Bytes())

	}

}

func hostInfo(ctx *gin.Context) {
	// For static
	info, err := host.HostInfo()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, info)
}

func stream(ctx *gin.Context) {

	ticker := time.NewTicker(500 * time.Millisecond)
	defer func() {
		ticker.Stop()
	}()

	ctx.Stream(func(w io.Writer) bool {
		select {
		case tm := <-ticker.C:

			ctx.SSEvent("", tm)

		}
		return true
	})
}
