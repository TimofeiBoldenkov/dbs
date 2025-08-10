package addprovidersinfo

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func AddProvidersInfo(ctx *fiber.Ctx) error {
	contentType := ctx.Get("Content-Type")
	if contentType != "application/json" {
		return fmt.Errorf("invalid Content-Type: %v (should be 'application/json')", contentType)
	}

	providerName := ctx.Params("providername")
	info := string(ctx.BodyRaw())

	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("can't load .env file: %v", err)
	}
	dbsDatabaseUrl := os.Getenv("DBS_DATABASE_URL")

	conn, err := pgx.Connect(context.Background(), dbsDatabaseUrl)
	if err != nil {
		return fmt.Errorf("can't connect to %v: %v", dbsDatabaseUrl, err)
	}
	defer conn.Close(context.Background())

	query := fmt.Sprintf("INSERT INTO %v (provider_name, info) VALUES ($1, $2)", 
		os.Getenv("DBS_TABLE_NAME"))
	_, err = conn.Exec(context.Background(), query, providerName, info)
	if err != nil {
		return fmt.Errorf("can't insert into dbs database: %v", err)
	}

	return ctx.SendStatus(201)
}
