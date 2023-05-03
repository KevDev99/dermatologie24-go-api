package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/KevDev99/dermatologie24-go-api/configs"
	"github.com/KevDev99/dermatologie24-go-api/routes"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	configs.ConnectDB()

	routes.UserRoute(router)
	routes.BookingRoute(router)
	routes.AuthRoute(router)
	routes.StatsRoute(router)

	fmt.Println("Server is ready to handle requests at", configs.EnvPort())
	err := http.ListenAndServe(":"+configs.EnvPort(), router)

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Fatal("Server stopped.")

}
