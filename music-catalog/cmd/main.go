package main

import (
	"log"
	"music-catalog/internal/configs"
	"music-catalog/pkg/internalsql"

	"github.com/gin-gonic/gin"
)

func main() {
	var (
		cfg *configs.Config
	)

	err := configs.Init(
		configs.WithConfigFolder(
			[]string{"./internal/configs"},
		),
		configs.WithConfigFile("config"),
		configs.WithConfigType("yaml"),
	)

	if err != nil {
		log.Fatal("Failed to initialized config")
	}

	cfg = configs.Get()

	_, err = internalsql.Connect(cfg.Database.DataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect database : %v", err)
	}

	r := gin.Default()
	r.Run(cfg.Service.Port)
}
