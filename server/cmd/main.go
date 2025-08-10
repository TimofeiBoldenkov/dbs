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

	conn, err := pgx.Connect(context.Background(), os.Getenv("DEFAULT_DATABASE_URL"))
	utils.ExitOnErr(err, "unable to connect to default database: %v", err)
	defer conn.Close(context.Background())

	var dbsExists string
	var dbsDatabaseName = os.Getenv("DBS_DATABASE_NAME")
	err = conn.QueryRow(context.Background(),
		"SELECT 'true' FROM pg_database WHERE datname = $1",
		dbsDatabaseName).Scan(&dbsExists)
	if !errors.Is(err, pgx.ErrNoRows) {
		utils.ExitOnErr(err, "unable to find out whether '%v' database exists or not: %v", dbsDatabaseName, err)
	}
	if dbsExists != "true" {
		err = conn.QueryRow(context.Background(), "CREATE DATABASE $1", dbsDatabaseName).Scan()
		utils.ExitOnErr(err, "unable to create '%v' database: %v", dbsDatabaseName, err)
	}

	dbsConn, err := pgx.Connect(context.Background(), os.Getenv("DBS_DATABASE_URL"))
	utils.ExitOnErr(err, "unable to connect to '%v' database: %v", dbsDatabaseName, err)
	defer dbsConn.Close(context.Background())

	tableName := os.Getenv("DBS_TABLE_NAME")
	createTableQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %v (
			id				SERIAL PRIMARY KEY,
			provider_name	TEXT NOT NULL,
			info			JSONB NOT NULL
		)
	`, tableName)
	_, err = dbsConn.Exec(context.Background(), createTableQuery)
	utils.ExitOnErr(err, "unable to create '%v' table in '%v' database: %v", tableName, dbsDatabaseName, err)

	app := fiber.New()

	app.Post("/api/new-info/:providername", addprovidersinfo.AddProvidersInfo)

	err = app.Listen(":" + os.Getenv("PORT"))
	utils.ExitOnErr(err, "server exits on error: %v", err)
}
