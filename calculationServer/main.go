package main

import (
	_ "calculationServer/docs"
	"calculationServer/internal/Api"
)

//	@title			Calculation Server API
//	@version		1.0
//	@description	This is a calculation server.

// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	a := Api.NewApi(STORAGE_URL, SECRET_KEY)
	a.Start(a.SetupRouter())
}
