package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"diploma/modules/category/model"
	"diploma/modules/category/service"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
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

// CreateCategory godoc
// @Summary      Create a new category
// @Description  Create a new product category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        category  body      model.CreateCategoryRequest  true  "Category data"
// @Success      201       {object}  model.Category
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /api/categories [post]
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req model.CreateCategoryRequest
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
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// UpdateCategory godoc
// @Summary      Update a category
// @Description  Update an existing category by ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        categoryID  path      int                        true  "Category ID"
// @Param        category    body      model.UpdateCategoryRequest  true  "Updated category data"
// @Success      200         {object}  model.Category
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /api/categories/{categoryID} [put]
func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.Atoi(chi.URLParam(r, "categoryID"))
	if err != nil {
		http.Error(w, "invalid category ID", http.StatusBadRequest)
		return
	}

	var req model.UpdateCategoryRequest
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

// DeleteCategory godoc
// @Summary      Delete a category
// @Description  Delete a category by ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        categoryID  path  int  true  "Category ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/categories/{categoryID} [delete]
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

// GetCategory godoc
// @Summary      Get a category
// @Description  Get a category by ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        categoryID  path      int  true  "Category ID"
// @Success      200         {object}  model.Category
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /api/categories/{categoryID} [get]
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

// ListCategories godoc
// @Summary      List all categories
// @Description  Get a list of all categories
// @Tags         categories
// @Accept       json
// @Produce      json
// @Success      200  {array}   model.Category
// @Failure      500  {object}  map[string]string
// @Router       /api/categories [get]
func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.ListCategories(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// CreateSubcategory godoc
// @Summary      Create a new subcategory
// @Description  Create a new subcategory for a specific category
// @Tags         subcategories
// @Accept       json
// @Produce      json
// @Param        categoryID    path      int                           true  "Category ID"
// @Param        subcategory   body      model.CreateSubcategoryRequest  true  "Subcategory data"
// @Success      201           {object}  model.Subcategory
// @Failure      400           {object}  map[string]string
// @Failure      404           {object}  map[string]string
// @Failure      500           {object}  map[string]string
// @Router       /api/categories/{categoryID}/subcategories [post]
func (h *Handler) CreateSubcategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.Atoi(chi.URLParam(r, "categoryID"))
	if err != nil {
		http.Error(w, "invalid category ID", http.StatusBadRequest)
		return
	}

	var req model.CreateSubcategoryRequest
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
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subcategory)
}

// UpdateSubcategory godoc
// @Summary      Update a subcategory
// @Description  Update an existing subcategory by ID
// @Tags         subcategories
// @Accept       json
// @Produce      json
// @Param        subcategoryID  path      int                           true  "Subcategory ID"
// @Param        subcategory    body      model.UpdateSubcategoryRequest  true  "Updated subcategory data"
// @Success      200            {object}  model.Subcategory
// @Failure      400            {object}  map[string]string
// @Failure      404            {object}  map[string]string
// @Failure      500            {object}  map[string]string
// @Router       /api/subcategories/{subcategoryID} [put]
func (h *Handler) UpdateSubcategory(w http.ResponseWriter, r *http.Request) {
	subcategoryID, err := strconv.Atoi(chi.URLParam(r, "subcategoryID"))
	if err != nil {
		http.Error(w, "invalid subcategory ID", http.StatusBadRequest)
		return
	}

	var req model.UpdateSubcategoryRequest
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

// DeleteSubcategory godoc
// @Summary      Delete a subcategory
// @Description  Delete a subcategory by ID
// @Tags         subcategories
// @Accept       json
// @Produce      json
// @Param        subcategoryID  path  int  true  "Subcategory ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/subcategories/{subcategoryID} [delete]
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

// GetSubcategory godoc
// @Summary      Get a subcategory
// @Description  Get a subcategory by ID
// @Tags         subcategories
// @Accept       json
// @Produce      json
// @Param        subcategoryID  path      int  true  "Subcategory ID"
// @Success      200            {object}  model.Subcategory
// @Failure      400            {object}  map[string]string
// @Failure      404            {object}  map[string]string
// @Failure      500            {object}  map[string]string
// @Router       /api/subcategories/{subcategoryID} [get]
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

// ListSubcategories godoc
// @Summary      List subcategories
// @Description  Get a list of all subcategories for a specific category
// @Tags         subcategories
// @Accept       json
// @Produce      json
// @Param        categoryID  path      int  true  "Category ID"
// @Success      200         {array}   model.Subcategory
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /api/categories/{categoryID}/subcategories [get]
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
