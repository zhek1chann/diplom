package handler

import (
	"net/http"
	"strconv"

	"diploma/modules/category/model"
	"diploma/modules/category/service"

	"github.com/gin-gonic/gin"
)

type GinHandler struct {
	service service.Service
}

func NewGinHandler(service service.Service) *GinHandler {
	return &GinHandler{service: service}
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
func (h *GinHandler) CreateCategory(c *gin.Context) {
	var req model.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.service.CreateCategory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
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
func (h *GinHandler) UpdateCategory(c *gin.Context) {
	categoryIDStr := c.Param("categoryID")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	var req model.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.service.UpdateCategory(c.Request.Context(), categoryID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
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
func (h *GinHandler) DeleteCategory(c *gin.Context) {
	categoryIDStr := c.Param("categoryID")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	if err := h.service.DeleteCategory(c.Request.Context(), categoryID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
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
func (h *GinHandler) GetCategory(c *gin.Context) {
	categoryIDStr := c.Param("categoryID")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	category, err := h.service.GetCategory(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
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
func (h *GinHandler) ListCategories(c *gin.Context) {
	categories, err := h.service.ListCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoriesTree godoc
// @Summary      Get categories tree
// @Description  Get a hierarchical tree of all categories with their subcategories
// @Tags         categories
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.CategoriesTreeResponse
// @Failure      500  {object}  map[string]string
// @Router       /api/categories/tree [get]
func (h *GinHandler) GetCategoriesTree(c *gin.Context) {
	categoriesTree, err := h.service.GetCategoriesTree(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categoriesTree)
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
func (h *GinHandler) CreateSubcategory(c *gin.Context) {
	categoryIDStr := c.Param("categoryID")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	var req model.CreateSubcategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.CategoryID = categoryID

	subcategory, err := h.service.CreateSubcategory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subcategory)
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
func (h *GinHandler) UpdateSubcategory(c *gin.Context) {
	subcategoryIDStr := c.Param("subcategoryID")
	subcategoryID, err := strconv.Atoi(subcategoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subcategory ID"})
		return
	}

	var req model.UpdateSubcategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subcategory, err := h.service.UpdateSubcategory(c.Request.Context(), subcategoryID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subcategory)
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
func (h *GinHandler) DeleteSubcategory(c *gin.Context) {
	subcategoryIDStr := c.Param("subcategoryID")
	subcategoryID, err := strconv.Atoi(subcategoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subcategory ID"})
		return
	}

	if err := h.service.DeleteSubcategory(c.Request.Context(), subcategoryID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
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
func (h *GinHandler) GetSubcategory(c *gin.Context) {
	subcategoryIDStr := c.Param("subcategoryID")
	subcategoryID, err := strconv.Atoi(subcategoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subcategory ID"})
		return
	}

	subcategory, err := h.service.GetSubcategory(c.Request.Context(), subcategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subcategory)
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
func (h *GinHandler) ListSubcategories(c *gin.Context) {
	categoryIDStr := c.Param("categoryID")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	subcategories, err := h.service.ListSubcategories(c.Request.Context(), categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subcategories)
}
