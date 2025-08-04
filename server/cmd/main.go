package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	addprovidersinfo "github.com/TimofeiBoldenkov/dbs/server/handlers/add_providers_info"
	"github.com/TimofeiBoldenkov/dbs/server/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	utils.ExitOnErr(err, "unable to load .env file: %v", err)
	var dbsDatabaseName = os.Getenv("DBS_DATABASE_NAME")

	conn, err := pgx.Connect(context.Background(), os.Getenv("DEFAULT_DATABASE_URL"))
	utils.ExitOnErr(err, "unable to connect to database: %v", err)
	defer conn.Close(context.Background())

	var dbsExists string
	err = conn.QueryRow(context.Background(),
		"SELECT 'true' FROM pg_database WHERE datname = $1",
		dbsDatabaseName).Scan(&dbsExists)
	if !errors.Is(err, pgx.ErrNoRows) {
		utils.ExitOnErr(err, "unable to execute QueryRow: %v", err)
	}
	if dbsExists != "true" {
		err = conn.QueryRow(context.Background(), "CREATE DATABASE $1", dbsDatabaseName).Scan()
		utils.ExitOnErr(err, "unable to execute QueryRow: %v", err)
	}

	dbsConn, err := pgx.Connect(context.Background(), os.Getenv("DBS_DATABASE_URL"))
	utils.ExitOnErr(err, "unable to connect to database: %v", err)
	defer dbsConn.Close(context.Background())

	createTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %v (
			id				SERIAL PRIMARY KEY,
			provider_name	TEXT NOT NULL,
			info			JSONB NOT NULL
		)
	`, os.Getenv("DBS_TABLE_NAME"))
	_, err = dbsConn.Exec(context.Background(), createTableQuery)
	utils.ExitOnErr(err, "unable to execute QueryRow: %v", err)

	app := fiber.New()

	app.Post("/api/new-info/:providername", addprovidersinfo.AddProvidersInfo)

	app.Listen(":" + os.Getenv("PORT"))
}
