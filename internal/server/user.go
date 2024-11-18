package server

import (
	"log"
	"ticketing/internal/database"
	"ticketing/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RegisterUser godoc
// @Summary Register a new user
// @Description This endpoint allows users to register by providing their details.
// @Tags Users
// @Accept json
// @Produce json
// @Param body body database.SignUpDTO true "User registration details"
// @Success 201 {object} database.User "User registered successfully"
// @Failure 400 {object} map[string]string "Error while processing the request"
// @Router /users/register [post]
func (s *FiberServer) RegisterUser(c *fiber.Ctx) error {
	var req database.SignUpDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user := database.User{
		Email:        req.Email,
		PasswordHash: utils.GeneratePassword(req.Password),
		UserID:       uuid.New().String(),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
	}
	if err := s.db.CreateUser(&user); err != nil {
		log.Printf("Error registering user: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": "Could not register, try again later!"})
	}

	return c.Status(fiber.StatusCreated).JSON(user)

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
func (s *FiberServer) LoginUser(c *fiber.Ctx) error {
	var req database.LoginDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user, err := s.db.GetUserByEmail(req.Email)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"message": "user not found!"})
	}

	if !utils.ComparePassword(user.PasswordHash, req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "incorrect password!"})
	}

	token, err := utils.GenerateToken(user.UserID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(fiber.Map{"token": token})

}
