package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/akrylysov/algnhsa"
	"github.com/xonha/todos/databases"
	"github.com/xonha/todos/views"
)

func main() {
	databases.Init()
	views.Init()

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		algnhsa.ListenAndServe(views.Router, nil)
	} else {

		port := os.Getenv("PORT")
		if port == "" {
			port = "3000"
		}
		fmt.Printf("🚀 Server listening on port %s...\n", port)
		if err := http.ListenAndServe(":"+port, views.Router); err != nil {
			panic(err)
		}
	}
}
