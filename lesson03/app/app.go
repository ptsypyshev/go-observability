package app

import (
	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github/ptsypyshev/go-observability/lesson03/handler"
	"github/ptsypyshev/go-observability/lesson03/l"
	"github/ptsypyshev/go-observability/lesson03/s"
	"github/ptsypyshev/go-observability/lesson03/store"
	"go.uber.org/zap"
	"io"
)

type App struct {
	logger *zap.Logger
	tracer opentracing.Tracer
}

func (a *App) Init() (io.Closer, error) {
	//ctx := context.Background()
	// Предустановленный конфиг. Можно выбрать
	// NewProduction/NewDevelopment/NewExample или создать свой
	// Production - уровень логгирования InfoLevel, формат вывода: json
	// Development - уровень логгирования DebugLevel, формат вывода: console
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	defer func() { _ = logger.Sync() }()
	// Трейсер
	// Можно "захардкодить" при инициализации
	//tracer, closer := l.InitJaeger("App", "jaeger:6831", logger)
	// Или использовать переменные окружения
	tracer, closer := l.InitJaeger(logger)
	// можно установить глобальный логгер (но лучше не надо: используйте внедрение
	// зависимостей где это возможно)
	// undo := zap.ReplaceGlobals(logger)
	// defer undo()
	// zap.L().Info("replaced zap's global loggers")

	a.logger = logger
	a.tracer = tracer

	return closer, nil
}

//func (a *App) JaegerMiddleware() gin.HandlerFunc {
//	JaegerOpen := true
//	return func(c *gin.Context) {
//		var serverSpan opentracing.Span
//		if JaegerOpen == true {
//			operationName := c.Request.Method + " " + c.Request.URL.Path
//			wireContext, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
//			serverSpan = opentracing.StartSpan(
//				operationName,
//				ext.RPCServerOption(wireContext),
//				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"})
//			defer serverSpan.Finish()
//			c.Set("Tracer", opentracing.GlobalTracer())
//			c.Set("SpanHttpContext", opentracing.ContextWithSpan(context.Background(), serverSpan))
//			serverSpan.SetTag("http.status_code", c.Writer.Status())
//
//			spanContext := serverSpan.Context()
//			if spanContext, ok := spanContext.(jaegerClient.SpanContext); ok {
//				c.Set("trace_id", spanContext.TraceID().String())
//				c.Set("span_id", spanContext.SpanID().String())
//			}
//		}
//		c.Next()
//	}
//}

func (a *App) Serve() error {
	//Sentry error handler
	s.NewSentryLogger()

	//Initialize Stores
	articleStore, err := store.NewArticleStore(a.logger, a.tracer)
	parseErr(err)

	//Initialize Handlers
	articleHandler := handler.NewArticleHandler(articleStore, a.logger, a.tracer)
	panicHandler := handler.NewPanicHandler(a.logger, a.tracer)

	//Initialize Router and add Middleware
	router := gin.Default()
	router.Use(nice.Recovery(panicHandler.RecoveryHandler))
	router.LoadHTMLFiles("template/error.tpl")

	//Routes
	router.GET("/article/id/:id", articleHandler.Id)
	router.POST("/article/add", articleHandler.Add)
	router.POST("/article/search", articleHandler.Search)
	router.GET("/panic", panicHandler.Panic)
	router.POST("/log/add", panicHandler.Log)

	// Start serving the application
	return router.Run()
}

func parseErr(err error) {
	if err != nil {
		l.F(err)
	}
	l.L("Application started")
}
