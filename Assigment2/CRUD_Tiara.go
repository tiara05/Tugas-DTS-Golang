package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Order represents the structure of the order request body
type Order struct {
	OrderedAt     string `json:"orderedAt"`
	CustomerName  string `json:"customerName"`
	Items         []Item `json:"items"`
}

// Item represents the structure of each item in the order
type Item struct {
	ItemCode    string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

// createOrderHandler handles the POST request to create a new order
func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Here you would normally insert the order into a database.
	// For demonstration purposes, we'll just log the order to the console.
	fmt.Printf("Order received: %+v\n", order)

	// Respond with the created order
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// getOrderHandler handles the GET request to retrieve an order by ID
func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	// Here you would retrieve the order from the database based on the ID.
	// For demonstration purposes, let's assume we have a sample order.
	order := Order{
		OrderedAt:    "2024-03-28T10:00:00Z",
		CustomerName: "John Doe",
		Items: []Item{
			{ItemCode: "item001", Description: "Product A", Quantity: 2},
			{ItemCode: "item002", Description: "Product B", Quantity: 1},
		},
	}

	// Respond with the retrieved order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// updateOrderHandler handles the PUT request to update an existing order
func updateOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Here you would update the order in the database based on the received data.
	// For demonstration purposes, we'll just log the updated order to the console.
	fmt.Printf("Order updated: %+v\n", order)

	// Respond with the updated order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// deleteOrderHandler handles the DELETE request to delete an existing order
func deleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	// Here you would delete the order from the database based on the ID.
	// For demonstration purposes, let's assume we have successfully deleted the order.
	fmt.Println("Order deleted successfully.")

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Order deleted successfully.")
}

func main() {
	http.HandleFunc("/orders", createOrderHandler)
	http.HandleFunc("/orders", getOrderHandler).Methods("GET")
	http.HandleFunc("/orders", updateOrderHandler).Methods("PUT")
	http.HandleFunc("/orders", deleteOrderHandler).Methods("DELETE")

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
