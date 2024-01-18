package main

import (
	"fmt"
	"store-procedures-mysql/internal/server"
)

func main() {

	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
	fmt.Println("Server started at port 8080")
}
