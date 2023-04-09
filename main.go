package main

import (
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

	log.Fatal(http.ListenAndServe(":6000", router))
}
