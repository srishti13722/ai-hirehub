package handlers

import (
	"context"
	"time"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/srishti13722/ai-hirehub/auth-service/config"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *fiber.Ctx) error{
	type SignUpRequest struct{
		FirstName string `json:"firstname"`
		LastName string `json:"lastname"`
		Email string `json:"email"`
		Password string `json:"password"`
		Role string `json:"role"`
	}

	req := new(SignUpRequest)
	if err := c.BodyParser(&req); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Invalid Request"})
	}

	//validate role
	validRoles := map[string]bool{"job_seeker" : true, "recruiter" : true, "admin":true}
	if !validRoles[req.Role]{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Invalid Role"})
	}

	//hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	}

	//insert into dp

	query := "INSERT INTO users (firstname, lastname, email, password, role) values ($1, $2, $3, $4, $5) RETURNING id"
	var UserID string
	err = config.DB.QueryRow(context.Background(), query, req.FirstName, req.LastName, req.Email, string(hashedPassword), req.Role).Scan(&UserID)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "DataBse error" + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message":"User Created SuccessFully!!", "id": UserID})
}

func Login( c *fiber.Ctx) error{
	type LoginRequest struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	loginReq := new(LoginRequest)
	if err := c.BodyParser(&loginReq); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Invalid Request"})
	}

	//fetch user from db

	var storedPassword, userID, role string 
	query := "SELECT id, password, role FROM users WHERE email = $1"
	err := config.DB.QueryRow(context.Background(), query, loginReq.Email).Scan(&storedPassword, &userID,&role)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Database error"+ err.Error()})
	}

	//compare password
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(loginReq.Password)); err != nil{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"erroer":"Invalid email or password"+ err.Error()})
	}

	//Generate JWT token with role

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id" : userID,
		"role" : role,
		"exp" : time.Now().Add(time.Hour * 3).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
	}

	return c.JSON(fiber.Map{"token": tokenString, "role": role})
}