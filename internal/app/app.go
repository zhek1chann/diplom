package app

import (
	"context"
	"diploma/docs"
	"diploma/internal/config"
	"diploma/modules/auth"
	"diploma/modules/cart"
	"diploma/modules/product"

	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	serviceProvider *serviceProvider
	httpServer      *gin.Engine
}

func NewApp(ctx context.Context) (*App, error) {
	app := &App{}

	err := app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Run() error {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := a.httpServerRun(); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	wg.Wait()
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	router := gin.Default()

	docs.SwaggerInfo.BasePath = ""
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiGroup := router.Group("/api")

	authHandler := a.serviceProvider.AuthHandler(ctx)
	auth.RegisterRoutes(apiGroup, authHandler)
	authMiddleware := a.serviceProvider.AuthMiddleware(ctx)

	productHandler := a.serviceProvider.ProductHandler(ctx)
	product.RegisterRoutes(apiGroup, productHandler)

	cartHanlder := a.serviceProvider.CartHandler(ctx)
	cartGroup := router.Group("/api")
	cartGroup.Use(authMiddleware.AuthMiddleware())
	cart.RegisterRoutes(cartGroup, cartHanlder)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = router
	a.httpServer.Use(func(c *gin.Context) {
		corsMiddleware.HandlerFunc(c.Writer, c.Request)
		c.Next()
	})
	return nil
}

func (a *App) httpServerRun() error {
	address := a.serviceProvider.HTTPConfig().Address()
	log.Printf("HTTP server is running on %s", address)
	return a.httpServer.Run(address)
}
