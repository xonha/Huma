package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/xonha/huma/databases"
	"github.com/xonha/huma/views"
)

func main() {
	databases.Init()
	views.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Printf("🚀 Server listening on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, views.Router); err != nil {
		panic(err)
	}
}
