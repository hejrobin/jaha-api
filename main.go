package main

import (
	// Local packages
	"jaha-api/routers"
)

func main() {
	defaultRouter := routers.GetDefaultRouter()

	defaultRouter.Run(":" + env.GetPort())
}
