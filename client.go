package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/MHarisAshfaq/go-rest-api-mongodb/models"
)

func client() {
	baseURL := "http://localhost:8080"
	sampleBook := models.Book{
		Title:  "The Alchemist",
		Author: "Paulo Coelho",
		ISBN:   "978-0062315007",
	}
	// Create a new book
	// createBook(baseURL, sampleBook)
	// Get all books
	books, err := getBooks(baseURL)

	if err != nil {
		// handle error
	}
	fmt.Println("Books=len", len(books))

	for _, book := range books {
		fmt.Printf("Book: %+v\n", book)
	}

	fmt.Println("Books=", books)
	fmt.Println("First Book id=", books[0].ID.Hex())

	// Update a book (replace {id} with actual book ID)
	sampleBook.Title = "Updated Book"
	updateBook(baseURL, books[0].ID.Hex(), sampleBook)

	// Get a single book (replace {id} with actual book ID)
	getBook(baseURL, books[0].ID.Hex())

	// Delete a book (replace {id} with actual book ID)
	deleteBook(baseURL, books[0].ID.Hex())
}

func createBook(baseURL string, book models.Book) {
	// Create a new book
	jsonBook, _ := json.Marshal(book)
	res, err := http.Post(baseURL+"/books", "application/json", bytes.NewBuffer(jsonBook))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}

func getBooks(baseURL string) ([]models.Book, error) {
	// Get all books
	res, err := http.Get(baseURL + "/books")
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	// body, _ := ioutil.ReadAll(res.Body)
	// fmt.Println(string(body))
	var books []models.Book
	jsonErr := json.NewDecoder(res.Body).Decode(&books)
	if err != nil {
		fmt.Println("Error decoding JSON:", jsonErr)
		return nil, jsonErr
	}

	return books, nil
}

func getBook(baseURL string, id string) {
	// Get a book
	res, err := http.Get(baseURL + "/books/" + id)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}

func updateBook(baseURL string, id string, book models.Book) {
	// Update a book
	fmt.Println("Updating book with ID:", id)
	jsonBook, _ := json.Marshal(book)
	req, _ := http.NewRequest(http.MethodPut, baseURL+"/books/"+id, bytes.NewBuffer(jsonBook))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}

func deleteBook(baseURL string, id string) {
	// Delete a book
	fmt.Println("Deleting book with ID:", id)
	req, _ := http.NewRequest(http.MethodDelete, baseURL+"/books/"+id, nil)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}
