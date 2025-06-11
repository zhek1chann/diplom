package category

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/categories", func(r chi.Router) {
		r.Post("/", h.CreateCategory)
		r.Get("/", h.ListCategories)
		r.Route("/{categoryID}", func(r chi.Router) {
			r.Get("/", h.GetCategory)
			r.Put("/", h.UpdateCategory)
			r.Delete("/", h.DeleteCategory)

			// Subcategory routes
			r.Get("/subcategories", h.ListSubcategories)
			r.Post("/subcategories", h.CreateSubcategory)
		})
	})

	r.Route("/subcategories/{subcategoryID}", func(r chi.Router) {
		r.Get("/", h.GetSubcategory)
		r.Put("/", h.UpdateSubcategory)
		r.Delete("/", h.DeleteSubcategory)
	})
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category, err := h.service.CreateCategory(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.Atoi(chi.URLParam(r, "categoryID"))
	if err != nil {
		http.Error(w, "invalid category ID", http.StatusBadRequest)
		return
	}

	var req UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category, err := h.service.UpdateCategory(r.Context(), categoryID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.Atoi(chi.URLParam(r, "categoryID"))
	if err != nil {
		http.Error(w, "invalid category ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteCategory(r.Context(), categoryID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.Atoi(chi.URLParam(r, "categoryID"))
	if err != nil {
		http.Error(w, "invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := h.service.GetCategory(r.Context(), categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.ListCategories(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *Handler) CreateSubcategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.Atoi(chi.URLParam(r, "categoryID"))
	if err != nil {
		http.Error(w, "invalid category ID", http.StatusBadRequest)
		return
	}

	var req CreateSubcategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.CategoryID = categoryID

	subcategory, err := h.service.CreateSubcategory(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subcategory)
}

func (h *Handler) UpdateSubcategory(w http.ResponseWriter, r *http.Request) {
	subcategoryID, err := strconv.Atoi(chi.URLParam(r, "subcategoryID"))
	if err != nil {
		http.Error(w, "invalid subcategory ID", http.StatusBadRequest)
		return
	}

	var req UpdateSubcategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subcategory, err := h.service.UpdateSubcategory(r.Context(), subcategoryID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subcategory)
}

func (h *Handler) DeleteSubcategory(w http.ResponseWriter, r *http.Request) {
	subcategoryID, err := strconv.Atoi(chi.URLParam(r, "subcategoryID"))
	if err != nil {
		http.Error(w, "invalid subcategory ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteSubcategory(r.Context(), subcategoryID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetSubcategory(w http.ResponseWriter, r *http.Request) {
	subcategoryID, err := strconv.Atoi(chi.URLParam(r, "subcategoryID"))
	if err != nil {
		http.Error(w, "invalid subcategory ID", http.StatusBadRequest)
		return
	}

	subcategory, err := h.service.GetSubcategory(r.Context(), subcategoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subcategory)
}

func (h *Handler) ListSubcategories(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.Atoi(chi.URLParam(r, "categoryID"))
	if err != nil {
		http.Error(w, "invalid category ID", http.StatusBadRequest)
		return
	}

	subcategories, err := h.service.ListSubcategories(r.Context(), categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subcategories)
}
