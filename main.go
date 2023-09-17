package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Movie struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	MovieID int                `json:"movie_id" bson:"movie_id"`
	Name    string             `json:"Name" bson:"Name"`
	Genre   string             `json:"Genre" bson:"Genre"`
	Quality string             `json:"Quality" bson:"Quality"`
	Rating  float64            `json:"Rating" bson:"Rating"`
	Year    int                `json:"year" bson:"year"`
}

var client *mongo.Client

func main() {
	r := gin.Default()

	// Welcome page with image
	r.GET("/", func(c *gin.Context) {
		htmlContent := `
        <!DOCTYPE html>
        <html>
        <head>
            <title>Welcome to Kubernetflix!</title>
        </head>
        <body style="text-align: center; padding-top: 100px;">
            <h1 style="color: #FF5733; font-family: Arial, sans-serif;">Welcome to Kubernetflix!</h1>
            <img src="https://scontent.fhfa2-2.fna.fbcdn.net/v/t39.30808-6/373003940_315302124329871_2467345526071715973_n.jpg?stp=cp6_dst-jpg&_nc_cat=110&ccb=1-7&_nc_sid=8bfeb9&_nc_ohc=_oPxNNy6jYsAX_1epgu&_nc_ht=scontent.fhfa2-2.fna&oh=00_AfCSaB9qh-cjKxwp86szSSw0M8RF308kf7IOwdlEsJQ8fQ&oe=64F45B10" alt="Netflix Logo" style="width: 300px; height: auto;">
            <p style="font-size: 18px;">Enjoy your stay in the world of Kubernetflix!</p>
        </body>
        </html>
        `
		c.Data(200, "text/html; charset=utf-8", []byte(htmlContent))
	})

	r.GET("/gmovies", getAllMovies)
	r.GET("/gmovies/:movie_id", getMovieByID)
	r.POST("/cmovies", createMovie)
	r.PUT("/umovies/:movie_id", updateMovie)
	r.DELETE("/dmovies/:movie_id", deleteMovie)

	// Initialize MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)

	if err := r.Run(":8010"); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func getAllMovies(c *gin.Context) {
	movies := []Movie{}
	collection := client.Database("movies_database").Collection("movies_collection")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var movie Movie
		cursor.Decode(&movie)
		movies = append(movies, movie)
	}
	c.JSON(http.StatusOK, movies)
}

func getMovieByID(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("movie_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	collection := client.Database("movies_database").Collection("movies_collection")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var movie Movie
	err = collection.FindOne(ctx, bson.M{"movie_id": movieID}).Decode(&movie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, movie)
}

func createMovie(c *gin.Context) {
	var newMovie Movie
	if err := c.BindJSON(&newMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := client.Database("movies_database").Collection("movies_collection")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	result, err := collection.InsertOne(ctx, newMovie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result.InsertedID)
}

func updateMovie(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("movie_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	var updatedMovie Movie
	if err := c.BindJSON(&updatedMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := client.Database("movies_database").Collection("movies_collection")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	result, err := collection.UpdateOne(ctx, bson.M{"movie_id": movieID}, bson.M{"$set": updatedMovie})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie updated"})
}

func deleteMovie(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("movie_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	collection := client.Database("movies_database").Collection("movies_collection")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	result, err := collection.DeleteOne(ctx, bson.M{"movie_id": movieID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted"})
}
