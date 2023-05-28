package main

import (
	"context"
	"fmt"
	"log"

	"robinhood/datastore"
	"robinhood/handlers"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	ctx := context.Background()
	cfg := readConfig()

	mongoDBClient := setupMongoDBClient(ctx, cfg)
	mongoInterviewsClient := datastore.NewInterviewsClient(mongoDBClient, cfg)
	mongoUsersClient := datastore.NewUsersClient(mongoDBClient, cfg)

	mongoInterviewsClient.InitInterviews(ctx)
	mongoUsersClient.InitUsers(ctx)

	router := gin.Default()

	if err := handlers.Setup(ctx, cfg, router, mongoInterviewsClient, mongoUsersClient).Run(); err != nil {
		return
	}
}

func readConfig() *viper.Viper {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}
	return viper.GetViper()
}

func setupMongoDBClient(ctx context.Context, cfg *viper.Viper) *mongo.Client {
	uri := fmt.Sprintf("mongodb://%s:%s@%s/test?authSource=admin",
		cfg.GetString("mongodb.dbuser"),
		cfg.GetString("mongodb.dbpassword"),
		cfg.GetString("mongodb.dbhost"))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(fmt.Errorf("could not connect to databse: %w", err))
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(fmt.Errorf("could not ping databse: %w", err))
	}

	log.Println("Connected to MongoDB")

	return client
}
