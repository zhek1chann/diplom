package handler

import (
	"log"
	"net/http"

	"diploma/modules/auth/model"
	"diploma/pkg/validator"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) Register(c *gin.Context) {
	var form struct {
		Name            string `json:"name"`
		PhoneNumber     string `json:"phone_number"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
		validator.Validator
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "Name cannot be blank")
	form.CheckField(validator.NotBlank(form.PhoneNumber), "phone_number", "Phone number cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be blank")
	form.CheckField(form.Password == form.ConfirmPassword, "confirm_password", "Passwords do not match")

	if !form.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"errors": form.Errors})
		return
	}

	authUser := &model.AuthUser{
		Info: &model.UserInfo{
			Name:        form.Name,
			PhoneNumber: form.PhoneNumber,
			Role:        0, // По умолчанию роль пользователя
		},
		Password: form.Password,
	}

	id, err := h.service.Register(c.Request.Context(), authUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Registered user with id: %d", id)

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var form struct {
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
		validator.Validator
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form.CheckField(validator.NotBlank(form.PhoneNumber), "phone_number", "Phone number cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be blank")

	if !form.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"errors": form.Errors})
		return
	}

	accessToken, refreshToken, err := h.service.Login(c.Request.Context(), form.PhoneNumber, form.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
