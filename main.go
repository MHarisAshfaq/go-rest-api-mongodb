package main

import (
	"fmt"
	"net/http"

	"github.com/MHarisAshfaq/go-rest-api-mongodb/db"
	"github.com/MHarisAshfaq/go-rest-api-mongodb/handlers"
	"github.com/gorilla/mux"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	router := mux.NewRouter()

	collection := db.ConnectDB()
	// Print the pointer value
	fmt.Printf("Pointer address: %p\n", collection)

	// Print the actual value (dereferencing the pointer)
	fmt.Printf("Collection details: %+v\n", *collection)

	// Optional: print a field to verify its details
	fmt.Printf("Collection name: %s\n", collection.Name())

	handler := &handlers.Handler{Collection: collection}
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	router.HandleFunc("/books", handler.GetBooks).Methods("GET")
	router.HandleFunc("/books/{id}", handler.GetBook).Methods("GET")
	router.HandleFunc("/books", handler.CreateBook).Methods("POST")

	http.ListenAndServe(":8080", router)
	// person := Person{Name: "John", Age: 30}
	// fmt.Println(person)

	// jsonData, err := json.Marshal(person)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// fmt.Println(jsonData)

	// var person2 Person

	// newErr := json.Unmarshal(jsonData, &person2)
	// if newErr != nil {
	// 	log.Fatal(newErr)
	// 	return
	// }
	// fmt.Println(person2)
}
