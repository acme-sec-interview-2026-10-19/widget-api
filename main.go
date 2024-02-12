package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const defaultConfigPath = "config/settings.yaml"

func main() {
	path := os.Getenv("WIDGET_API_CONFIG")
	if path == "" {
		path = defaultConfigPath
	}

	cfg, err := loadConfig(path)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	if cfg.Auth.JWTSigningKey == "" {
		log.Fatal("auth.jwt_signing_key is empty")
	}
	if cfg.AWS.AccessKeyID == "" || cfg.AWS.SecretAccessKey == "" {
		log.Printf("AWS credentials not set; inventory sync disabled")
	}

	addr := fmt.Sprintf(":%d", cfg.Service.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(newRouter(cfg), &http2.Server{}),
	}

	log.Printf("%s listening on %s (region %s)", cfg.Service.Name, addr, cfg.AWS.Region)
	log.Fatal(srv.ListenAndServe())
}

func newRouter(cfg *Config) *gin.Engine {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello, world")
	})

	admin := r.Group("/admin", requireJWT([]byte(cfg.Auth.JWTSigningKey), cfg.Service.Name))
	admin.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": cfg.Service.Name,
			"region":  cfg.AWS.Region,
			"status":  "ok",
			"caller":  c.GetString("subject"),
		})
	})

	return r
}
