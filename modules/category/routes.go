package category

import (
	categoryHandler "diploma/modules/category/handler"
	categoryRepo "diploma/modules/category/repository"
	categorySvc "diploma/modules/category/service"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

// RegisterRoutes sets up all category routes
func RegisterRoutes(r chi.Router, db *sqlx.DB) {
	repo := categoryRepo.NewRepository(db)
	svc := categorySvc.NewService(repo)
	h := categoryHandler.NewHandler(svc)

	h.RegisterRoutes(r)
}
