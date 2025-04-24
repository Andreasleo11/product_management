package main

import (
	"backend_prodman/db"
	"backend_prodman/routes"
)

func main() {
	db.InitDB()

	router := routes.SetupRouter()
	router.Run(":8081")
}
