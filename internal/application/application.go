package application

import (
	"context"
	"log/slog"

	"github.com/Grsaiago/gonebook_server/internal/database"
	"github.com/Grsaiago/gonebook_server/internal/services"
	"github.com/jackc/pgx/v5"
)

type Application struct {
	DbCtx          context.Context
	DbConn         *pgx.Conn
	ContactService *services.ContactService
}

func New() (Application, error) {
	// connect to database and setup pgx
	dbCtx := context.Background()
	dbConn, err := pgx.Connect(dbCtx, "postgres://postgres:123456@localhost:8080/postgres") // expose this as an envfile
	// conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL")) // should be this
	if err != nil {
		return Application{}, err
	}
	slog.Info("Database connection established!")

	// get queries from sqlc
	queries := database.New(dbConn)

	// create newApp
	newApp := Application{
		DbCtx:  dbCtx,
		DbConn: dbConn,
	}

	// initialize services
	newApp.ContactService = &services.ContactService{
		Repo:       queries,
		AppContext: &newApp.DbCtx,
	}

	return newApp, nil
}
