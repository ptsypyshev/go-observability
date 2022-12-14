package main

import (
	"github.com/ptsypyshev/go-observability/lesson02/handler"
	"github.com/ptsypyshev/go-observability/lesson02/l"
	"github.com/ptsypyshev/go-observability/lesson02/s"
	"github.com/ptsypyshev/go-observability/lesson02/store"

	"github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
)

// Переписать не на Martini
func main() {
	//Sentry error handler
	s.NewSentryLogger()

	//Initialize Stores
	articleStore, err := store.NewArticleStore()
	parseErr(err)

	//Initialize Handlers
	articleHandler := handler.NewArticleHandler(articleStore)
	panicHandler := handler.NewPanicHandler()

	//Initialize Router and add Middleware
	router := gin.Default()
	router.Use(nice.Recovery(handler.RecoveryHandler))
	router.LoadHTMLFiles("template/error.tpl")

	//Routes
	router.GET("/article/id/:id", articleHandler.Id)
	router.POST("/article/add", articleHandler.Add)
	router.POST("/article/search", articleHandler.Search)
	router.GET("/panic", panicHandler.Handle)
	router.POST("/log/add", panicHandler.Log)

	// Start serving the application
	router.Run()
}

func parseErr(err error) {
	if err != nil {
		l.F(err)
	}
	l.Log.Log("Application started")
}
