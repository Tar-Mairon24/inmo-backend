package cmd

import (
	"inmo-backend/internal/infrastructure/db"
	"inmo-backend/internal/interface/api"
)

func main() {

	db.Init()
	r := api.SetupRouter()
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
	r.Run(":8080")
}