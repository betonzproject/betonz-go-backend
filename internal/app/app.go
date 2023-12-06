package app

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/alexedwards/scs/goredisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type App struct {
	DB       *db.Queries
	Redis    *redis.Client
	Scs      *scs.SessionManager
	Validate *validator.Validate
}

func NewApp() *App {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln("Can't create connection pool: " + err.Error())
	}

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalln("Can't connect to redis: " + err.Error())
	}
	client := redis.NewClient(opt)

	sessionManager := scs.New()
	sessionManager.Cookie.Domain = os.Getenv("DOMAIN")
	sessionManager.Store = goredisstore.New(client)
	sessionManager.Lifetime = time.Duration(30 * 24 * time.Hour)
	if os.Getenv("ENVIRONMENT") != "development" {
		sessionManager.Cookie.Secure = true
	}

	validator := validator.New(validator.WithRequiredStructEnabled())

	return &App{
		DB:       db.New(pool),
		Redis:    client,
		Scs:      sessionManager,
		Validate: validator,
	}
}
