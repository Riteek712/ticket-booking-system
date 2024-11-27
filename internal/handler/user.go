package handler

import (
	"log"
	"ticketing/internal/database"
	"ticketing/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UserHandler represents the handler for user-related operations.
type UserHandler struct {
	db database.Service
}

// NewUserHandler initializes a new UserHandler with the given database service.
func NewUserHandler(db database.Service) *UserHandler {
	return &UserHandler{db: db}
}

// LoginUser godoc
// @Summary User login
// @Description This endpoint allows users to log in with their email and password.
// @Tags Users
// @Accept json
// @Produce json
// @Param body body database.LoginDTO true "User login details"
// @Success 200 {object} map[string]string "Login successful with token"
// @Failure 400 {object} map[string]string "User not found or invalid request"
// @Failure 401 {object} map[string]string "Incorrect password"
// @Failure 500 {object} map[string]string "Error generating token"
// @Router /users/login [post]
func (h *UserHandler) RegisterUser(c *fiber.Ctx) error {
	var req database.SignUpDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Hash password and create user
	user := database.User{
		UserID:       uuid.New().String(),
		Email:        req.Email,
		PasswordHash: utils.GeneratePassword(req.Password),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
	}

	if err := h.db.CreateUser(&user); err != nil {
		log.Printf("Error registering user: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Could not register, try again later!",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// LoginUser handles user login.
func (h *UserHandler) LoginUser(c *fiber.Ctx) error {
	var req database.LoginDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Retrieve user by email
	user, err := h.db.GetUserByEmail(req.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User not found!",
		})
	}

	// Compare hashed password
	if !utils.ComparePassword(user.PasswordHash, req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Incorrect password!",
		})
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate token!",
		})
	}

	return c.JSON(fiber.Map{"token": token})
}
