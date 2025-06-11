package category

import (
	categoryHandler "diploma/modules/category/handler"
	categoryRepo "diploma/modules/category/repository"
	categorySvc "diploma/modules/category/service"
	"diploma/pkg/client/db"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes sets up all category routes
func RegisterRoutes(r chi.Router, dbClient db.Client) {
	repo := categoryRepo.NewRepository(dbClient)
	svc := categorySvc.NewService(repo)
	h := categoryHandler.NewHandler(svc)

	h.RegisterRoutes(r)
}
