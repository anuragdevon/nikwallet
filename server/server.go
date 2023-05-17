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
	err := db.ConnectToDB("testdb")
	if err != nil {
		log.Panic("failed to connect to database:", err)
	}
	defer db.DB.Close()

	userService := services.NewUserService(db.DB)
	authService := services.NewAuthService(db.DB)
	walletService := services.NewWalletService(db.DB)

	userHandlers := handlers.NewUserHandlers(userService, authService)
	walletHandlers := handlers.NewWalletHandlers(walletService, authService, userService)

	router := routers.NewRouter(userHandlers, walletHandlers)

	fmt.Println("Server listening on port 8080...")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
