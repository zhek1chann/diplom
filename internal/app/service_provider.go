package app

import (
	"context"
	"log"

	"diploma/internal/config"

	"diploma/pkg/client/db"
	"diploma/pkg/client/db/pg"
	"diploma/pkg/client/db/transaction"
	"diploma/pkg/closer"

	"diploma/modules/auth/jwt"

	authApi "diploma/modules/auth/handler"
	userRepository "diploma/modules/auth/repository/user"
	authService "diploma/modules/auth/service/auth"

	productApi "diploma/modules/product/handler"
	productRepository "diploma/modules/product/repository/product"
	productService "diploma/modules/product/service"
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
	jwt            authService.IJWT
	authService    authApi.IAuthService
	authHanlder    *authApi.AuthHandler

	// product

	productRepository productService.IProductRepository
	productService    productApi.IProductService
	productHanlder    *productApi.CatalogHandler
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

func (s *serviceProvider) AuthRepository(ctx context.Context) authService.IAuthRepository {
	if s.authRepository == nil {
		s.authRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) JWT(ctx context.Context) authService.IJWT {
	if s.jwt == nil {
		s.jwt = jwt.NewJSONWebToken(s.JWTConfig().GetSecretKey())
	}

	return s.jwt
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

func (s *serviceProvider) ProductRepository(ctx context.Context) productService.IProductRepository {
	if s.productRepository == nil {
		s.productRepository = productRepository.NewRepository(s.DBClient(ctx))
	}

	return s.productRepository
}

func (s *serviceProvider) ProductService(ctx context.Context) productApi.IProductService {
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
