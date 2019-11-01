package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type RestAPI struct {
	engine *gin.Engine
	db     *gorm.DB
}

var (
	totalRequestsMetric = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "user_service",
		Name:      "total_requests",
		Help:      "Total number of requests",
	})

	internalServerErrorsResponsesMetric = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "user_service",
		Name:      "total_internal_server_errors_reponses",
		Help:      "Total number of internal server error reponses",
	})

	successfulResponsesMetric = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "user_service",
		Name:      "total_success_responses",
		Help:      "Total number of successful reponses",
	})
)

func NewRestAPI(cfg Config) (*RestAPI, error) {
	api := &RestAPI{}

	engine := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	api.engine = engine
	api.engine.Use(api.monitorMiddleware())
	api.setRoutes()

	api.setMonitorRoute()

	s := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Addr, cfg.DBName)

	db, err := gorm.Open("mysql", s)
	if err != nil {
		return nil, err
	}
	// db.LogMode(true)
	db.AutoMigrate(&Auth{}, &Profile{})
	api.db = db

	return api, err
}

func (a *RestAPI) setRoutes() {
	a.engine.GET("/profile/:userid", func(ctx *gin.Context) {
		GetProfile(a.db, ctx)
	})

	a.engine.POST("/login", func(ctx *gin.Context) {
		ValidLogin(a.db, ctx)
	})

	a.engine.POST("/account", func(ctx *gin.Context) {
		CreateAccount(a.db, ctx)
	})

	a.engine.POST("/profile/update", func(ctx *gin.Context) {
		UpdateProfile(a.db, ctx)
	})

	a.engine.POST("/password", func(ctx *gin.Context) {
		ChangePassword(a.db, ctx)
	})
}

func (a *RestAPI) setMonitorRoute() {
	a.engine.Any("/metrics", gin.WrapH(promhttp.Handler()))
}

func (a *RestAPI) Run(addr string) error {
	log.WithFields(appFields).Info("starting service on addr: " + addr)
	return a.engine.Run(addr)
}

func (a *RestAPI) monitorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		totalRequestsMetric.Inc()

		c.Next()

		if c.Writer.Status() == http.StatusInternalServerError {
			internalServerErrorsResponsesMetric.Inc()
			return
		}

		if c.Writer.Status() == http.StatusOK {
			successfulResponsesMetric.Inc()
			return
		}
	}
}
