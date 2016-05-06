package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/KSubedi/ip-geolocation-api-server/handlers"
	"github.com/KSubedi/ip-geolocation-api-server/utils"
)

func main() {
	// Initialize the IP Database
	utils.NewDB()

	// Close IP database on program exit
	defer utils.DB.Close()

	router := mux.NewRouter()

	// Handle the /V1/ API path
	router.HandleFunc("/v1/", handlers.V1)
	router.HandleFunc("/v1/{ip}", handlers.V1)

	// Always return 404 for Favicon
	router.Handle("/favicon.ico", http.NotFoundHandler())

	// Run SSL Server if possible
	go runSSLServer(router)

	// Run web server
	fmt.Println("Starting Server On Port 3000")
	http.ListenAndServe(":3000", router)

}

// Runs SSL server if certificated are available
func runSSLServer(router *mux.Router) {
	_, err1 := os.Stat("ssl.cert")
	_, err2 := os.Stat("ssl.key")

	if err1 == nil && err2 == nil {
		fmt.Println("Starting SSL Server On Port 3001")
		http.ListenAndServeTLS(":3001", "ssl.cert", "ssl.key", router)
	}
}
