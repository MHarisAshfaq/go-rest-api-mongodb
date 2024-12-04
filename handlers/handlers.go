package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MHarisAshfaq/go-rest-api-mongodb/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Handler struct {
	Collection *mongo.Collection
}

// get all books
func (h *Handler) GetBooks(w http.ResponseWriter, r *http.Request) {
	// get all books from the database
	w.Header().Set("Content-Type", "application/json")

	var books []models.Book
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	cursor, err := h.Collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer cursor.Close(ctx)

	fmt.Println(cursor)
	for cursor.Next(ctx) {
		var book models.Book
		cursor.Decode(&book)
		books = append(books, book)
	}
	if err := cursor.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(books)
}

// get a book

func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book models.Book
	fmt.Println(mux.Vars(r))
	id, _ := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	fmt.Println(id)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := h.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(book)

}

// create a books
func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var book models.Book
	json.NewDecoder(r.Body).Decode(&book)
	fmt.Println(book)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	book.ID = primitive.NewObjectID()
	result, err := h.Collection.InsertOne(ctx, book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(result)
}

// update a book
func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book models.Book
	json.NewDecoder(r.Body).Decode(&book)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	id, _ := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	result, err := h.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.D{{Key: "$set", Value: book}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(result)
}

// delete a book
func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	id, _ := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	result, err := h.Collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(result)
}
