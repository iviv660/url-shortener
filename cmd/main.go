// Package main URL Shortener API.
//
// Сервис сокращения URL и редирект по короткому коду.
//
// @title       URL Shortener API
// @version     1.0
// @description Мини-сервис для сокращения URL с хранением в Postgres и кэшем в Redis.
// @BasePath    /
package main

import (
	"app/internal/config"
	"app/internal/database"
	_ "app/internal/docs"
	"app/internal/handler"
	"app/internal/repository"
	"app/internal/usecase"
	"context"
	"log"
)

func main() {
	postgrsDB, err := database.ConnectPostgres(context.Background(), config.C.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer postgrsDB.Close()
	pDB := repository.NewRepoPostgres(postgrsDB)
	redisDb, err := database.ConnectRedis(context.Background(), config.C.RedisURL)
	if err != nil {
		log.Fatal(err)
	}
	defer redisDb.Close()
	rDB := repository.NewRepoRedis(redisDb)

	uc := usecase.New(pDB, rDB, config.C.CacheTTL, config.C.Secret)
	_, router := handler.NewHandler(uc)
	if err := router.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}
