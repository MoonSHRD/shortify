package app

import (
	"context"
	"fmt"
	"github.com/MoonSHRD/logger"
	"net/http"

	"github.com/MoonSHRD/shortify/config"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	Config   *config.Config
	DBClient *mongo.Client
	DB       *mongo.Database

	server *http.Server
}

func NewApp(config *config.Config) (*App, error) {
	app := &App{}
	app.Config = config
	err := app.createMongoDBConnection(&config.MongoDB)
	if err != nil {
		return nil, err
	}

	port := app.Config.HTTP.Port
	addr := app.Config.HTTP.Address
	server := &http.Server{Addr: fmt.Sprintf("%s:%d", addr, port)}
	app.server = server

	return app, nil
}

func (app *App) createMongoDBConnection(config *config.MongoDB) error {
	var mongoURI string
	if config.User == "" && config.Password == "" {
		mongoURI = fmt.Sprintf("mongodb://%s:%d", config.Host, config.Port)
	} else {
		mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%d", config.User, config.Password, config.Host, config.Port)
	}
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	app.DBClient = client
	app.DB = client.Database(app.Config.MongoDB.DatabaseName)
	return nil
}

func (app *App) Run(r *mux.Router) {
	port := app.Config.HTTP.Port
	addr := app.Config.HTTP.Address
	app.server.Handler = r

	logger.Infof("HTTP server starts listening at %s:%d", addr, port)
	go func() {
		if err := app.server.ListenAndServe(); err != nil {
			logger.Info(err)
		}
	}()
}

func (app *App) Destroy() {
	ctx := context.Background()
	_ = app.DBClient.Disconnect(ctx)
	_ = app.server.Shutdown(ctx)
}
