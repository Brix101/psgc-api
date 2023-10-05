package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/Brix101/psgc-api/internal/cmd"
)

//	@title		Philippine Standard Geographic Code (PSGC) API
//	@version	1.0
// description This is a sample server Petstore server.
// termsOfService http://swagger.io/terms/

// contact.name API Support
// contact.url http://www.swagger.io/support
// contact.email support@swagger.io

//	@BasePath	/api
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	ret := cmd.Execute(ctx)
	os.Exit(ret)
}
