package app

import (
	"context"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/alexedwards/scs/goredisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type App struct {
	DB       *db.Queries
	Pool     *pgxpool.Pool
	Redis    *redis.Client
	Scs      *scs.SessionManager
	Validate *validator.Validate
}

func NewApp() *App {
	// Set up database connection pools
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln("Can't create pgxpool config: " + err.Error())
	}

	config.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		registerType(ctx, c, "UserStatus")
		registerType(ctx, c, "Role")
		registerType(ctx, c, "TransactionType")
		registerType(ctx, c, "TransactionStatus")
		registerType(ctx, c, "EventType")
		registerType(ctx, c, "EventResult")
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalln("Can't create connection pool: " + err.Error())
	}

	// Redis
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalln("Can't connect to redis: " + err.Error())
	}
	client := redis.NewClient(opt)

	// Sessions
	sessionManager := scs.New()
	sessionManager.Cookie.Domain = os.Getenv("DOMAIN")
	sessionManager.Store = goredisstore.New(client)
	sessionManager.Lifetime = time.Duration(30 * 24 * time.Hour)
	if os.Getenv("ENVIRONMENT") != "development" {
		sessionManager.Cookie.Secure = true
	}

	// Validator
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		regex := regexp.MustCompile("^[a-zA-Z0-9_]+$")
		return regex.MatchString(fl.Field().String())
	})

	return &App{
		DB:       db.New(pool),
		Pool:     pool,
		Redis:    client,
		Scs:      sessionManager,
		Validate: validate,
	}
}

func registerType(ctx context.Context, c *pgx.Conn, name string) {
	t, err := c.LoadType(ctx, "\""+name+"\"")
	if err != nil {
		log.Fatal(err)
	}
	c.TypeMap().RegisterType(t)

	t, err = c.LoadType(ctx, "\"_"+name+"\"")
	if err != nil {
		log.Fatal(err)
	}
	c.TypeMap().RegisterType(t)
}
