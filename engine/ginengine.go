package engine

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
)

type GinEngine struct {
	engine *gin.Engine
	server *http.Server
}

func NewGinEngine(port string, routersToInstall ...func(*gin.Engine)) *GinEngine {
	engine := gin.Default()
	installHealthCheck(engine)
	installMiddlewares(engine)
	for _, router := range routersToInstall {
		router(engine)
	}
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           engine,
		ReadHeaderTimeout: time.Minute * 5,
	}
	return &GinEngine{engine: engine, server: server}
}

func (s *GinEngine) RunHttpServer() {
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}

func (s *GinEngine) Shutdown() error {
	// The context is used to inform the server it has 5 seconds to finish
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(c)
}

func installHealthCheck(engine *gin.Engine) {
	engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "up and running"})
	})
}

func installMiddlewares(engine *gin.Engine) {
	engine.Use(recoverConfig())
	engine.RedirectTrailingSlash = false
	engine.RemoveExtraSlash = true
}

func recoverConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("ginEngine:server recovered ,error:", r)
				c.JSON(http.StatusInternalServerError, gin.H{"message": "server panic"})
			}
		}()
		c.Next()
	}
}
