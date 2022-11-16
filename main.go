package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"gitlab.com/techpotion/leadershack2022/api/di"
	"go.uber.org/zap"
)

// @title          Swagger for TechPotion's leadershack2022 solution
// @version        1.0
// @description    TechPotion ЛЦТ2022 Swagger
// @termsOfService http://swagger.io/terms/

// @contact.name  API Support
// @contact.url   http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1
// @Schemes  http
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	ndi, err := di.NewDI(ctx)
	if err != nil {
		log.Fatalf("di init error: %v", err)
	}

	defer cancel()

	z := zap.S().With("context", "main")

	z.Info("Starting service")

	ctx = ndi.Start(ctx)

	<-ctx.Done()

	z.Infow("Stopping service", "reason", ctx.Err())
	ndi.Stop()

	z.Info("Service stopped")
}
