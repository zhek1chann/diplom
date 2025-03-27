package app

import (
	"context"
	"log"

	"diploma/internal/config"

	"diploma/pkg/client/db"
	"diploma/pkg/client/db/pg"
	"diploma/pkg/client/db/transaction"
	"diploma/pkg/closer"

	authApi "diploma/modules/auth/handler"
	authJWT "diploma/modules/auth/jwt"
	authMiddlware "diploma/modules/auth/middleware"
	userRepository "diploma/modules/auth/repository/user"
	authService "diploma/modules/auth/service/auth"

	productApi "diploma/modules/product/handler"
	productRepository "diploma/modules/product/repository/product"
	productService "diploma/modules/product/service"

	cartApi "diploma/modules/cart/handler"
	cartRepository "diploma/modules/cart/repository"
	cartService "diploma/modules/cart/service"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	jwtConfig     config.JWTConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig

	dbClient  db.Client
	txManager db.TxManager

	// auth
	authRepository authService.IAuthRepository
	authJWT        *authJWT.JSONWebToken
	authService    authApi.IAuthService
	authHanlder    *authApi.AuthHandler
	authMiddlware  *authMiddlware.AuthMiddleware

	// product

	productRepository productService.IProductRepository
	productService    *productService.ProductService
	productHanlder    *productApi.CatalogHandler

	// cart

	cartRepository cartService.ICartRepository
	cartService    cartApi.ICartService
	cartHanlder    *cartApi.CartHandler
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) JWTConfig() config.JWTConfig {
	if s.jwtConfig == nil {
		cfg, err := config.NewJWTConfig()
		if err != nil {
			log.Fatalf("failed to get jwt config: %s", err.Error())
		}

		s.jwtConfig = cfg
	}

	return s.jwtConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := config.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

// ========= authentication =========
func (s *serviceProvider) AuthRepository(ctx context.Context) authService.IAuthRepository {
	if s.authRepository == nil {
		s.authRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) JWT(ctx context.Context) *authJWT.JSONWebToken {
	if s.authJWT == nil {
		s.authJWT = authJWT.NewJSONWebToken(s.JWTConfig().GetSecretKey())
	}

	return s.authJWT
}

func (s *serviceProvider) AuthService(ctx context.Context) authApi.IAuthService {
	if s.authService == nil {
		s.authService = authService.NewService(s.AuthRepository(ctx), s.JWT(ctx), s.TxManager(ctx))
	}

	return s.authService
}

func (s *serviceProvider) AuthHandler(ctx context.Context) *authApi.AuthHandler {
	if s.authHanlder == nil {
		s.authHanlder = authApi.NewHandler(s.AuthService(ctx))
	}

	return s.authHanlder
}

func (s *serviceProvider) AuthMiddleware(ctx context.Context) *authMiddlware.AuthMiddleware {
	if s.authMiddlware == nil {
		s.authMiddlware = authMiddlware.NewAuthMiddleware(s.JWT(ctx))
	}

	return s.authMiddlware
}

// ========= product =========

func (s *serviceProvider) ProductRepository(ctx context.Context) productService.IProductRepository {
	if s.productRepository == nil {
		s.productRepository = productRepository.NewRepository(s.DBClient(ctx))
	}

	return s.productRepository
}

func (s *serviceProvider) ProductService(ctx context.Context) *productService.ProductService {
	if s.productService == nil {
		s.productService = productService.NewService(s.ProductRepository(ctx), s.TxManager(ctx))
	}

	return s.productService
}

func (s *serviceProvider) ProductHandler(ctx context.Context) *productApi.CatalogHandler {
	if s.productHanlder == nil {
		s.productHanlder = productApi.NewHandler(s.ProductService(ctx))
	}

	return s.productHanlder
}

// ========= cart =========

func (s *serviceProvider) CartRepo(ctx context.Context) cartService.ICartRepository {
	if s.cartRepository == nil {
		s.cartRepository = cartRepository.NewRepository(s.DBClient(ctx))
	}

	return s.cartRepository

}

func (s *serviceProvider) CartService(ctx context.Context) cartApi.ICartService {
	if s.cartService == nil {
		s.cartService = cartService.NewService(s.CartRepo(ctx), s.ProductService(ctx), s.TxManager(ctx))
	}

	return s.cartService
}

func (s *serviceProvider) CartHandler(ctx context.Context) *cartApi.CartHandler {
	if s.cartHanlder == nil {
		s.cartHanlder = cartApi.NewHandler(s.CartService(ctx))
	}

	return s.cartHanlder
}
