package server

import (
	"fmt"
	"log"
	"net/http"

	"nikwallet/handlers"
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
	userService := services.NewUserService(database)
	authService := services.NewAuthService(database)
	walletService := services.NewWalletService(database)

	userHandlers := handlers.NewUserHandlers(userService, authService)
	walletHandlers := handlers.NewWalletHandlers(walletService, authService, userService)

	router := routers.NewRouter(userHandlers, walletHandlers)

	fmt.Println("Server listening on port 8080...")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
