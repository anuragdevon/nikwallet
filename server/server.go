package server

import (
	"fmt"
	"log"
	"net/http"

	"nikwallet/pkg/db"
	"nikwallet/routers"
	"nikwallet/services"
)

func StartServer() {
	database, err := db.ConnectToDB("testdb")
	if err != nil {
		log.Panic("failed to connect to database:")
	}
	defer database.Close()
	// Initialize services
	userService := services.NewUserService(database)
	authService := services.NewAuthService(database)
	walletService := services.NewWalletService(database)

	// Initialize handlers with the services
	// userHandlers := handlers.NewUserHandlers(userService, authService)
	// walletHandlers := handlers.NewWalletHandlers(*walletService, *authService, *userService)

	// Initialize the router
	router := routers.NewRouter(userService, authService, walletService)

	// Start the server
	fmt.Println("Server listening on port 8080...")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
