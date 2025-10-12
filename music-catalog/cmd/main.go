package main

import (
	"log"
	"music-catalog/internal/configs"
	membershipsHdl "music-catalog/internal/handler/memberships"
	"music-catalog/internal/models/memberships"
	membershipsRepo "music-catalog/internal/repository/memberships"
	membershipsSvc "music-catalog/internal/service/memberships"
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

	db, err := internalsql.Connect(cfg.Database.DataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect database : %v", err)
	}

	db.AutoMigrate(&memberships.User{})
	r := gin.Default()

	membershipRepo := membershipsRepo.NewRepository(db)
	membershipSvc := membershipsSvc.NewService(cfg, membershipRepo)
	membershipHdl := membershipsHdl.NewHandler(r, membershipSvc)
	membershipHdl.RegisterRoute()

	r.Run(cfg.Service.Port)
}
