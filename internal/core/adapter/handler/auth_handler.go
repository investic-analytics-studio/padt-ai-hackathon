package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type AuthHandler struct {
	service    port.AuthService
	jwtService port.JwtService
}

func NewAuthHandler(service port.AuthService, jwtService port.JwtService) *AuthHandler {
	return &AuthHandler{
		service:    service,
		jwtService: jwtService,
	}
}

func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	var body model.SignUpBody
	if err := c.BodyParser(&body); err != nil {
		fmt.Println(err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	err := h.service.SignUp(body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Sign up successfully"})
}

// ExistEmail godoc
// @Summary      Check if email exists
// @Description  Verifies if a given email is already registered
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      model.EmailReq  true  "Email request"
// @Success      200   {object}  map[string]bool
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /auth/exist-email [post]
func (h *AuthHandler) ExistEmail(c *fiber.Ctx) error {
	var emailReq model.EmailReq
	if err := c.BodyParser(&emailReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	ois, err := h.service.ExistEmail(emailReq.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(ois)
}

func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	email := c.Locals("email").(string)
	twitterUID := c.Locals("twitter_uid").(string)
	twitterName := c.Locals("twitter_name").(string)
	return c.Status(http.StatusOK).JSON(fiber.Map{"uid": uid, "email": email, "twitter_uid": twitterUID, "twitter_name": twitterName})
}

func (h *AuthHandler) AutoValidateEmailInFirebase(c *fiber.Ctx) error {
	var autoValidate model.AutoValidateRequest
	if err := c.BodyParser(&autoValidate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	log.Println(autoValidate)
	err := h.service.AutoValidateEmailInFirebase(c.UserContext(), autoValidate.UID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Email is valid"})
}

func (h *AuthHandler) AllUserAutoValidateEmailInFirebase(c *fiber.Ctx) error {
	err := h.service.AllUserAutoValidateEmailInFirebase(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "All emails are valid"})
}

func (h *AuthHandler) GenerateToken(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	loginType := c.Locals("login_type").(string)
	token, err := h.jwtService.GenerateToken(c.UserContext(), uid, loginType)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"token": token})
}

func (h *AuthHandler) VerifyToken(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
	}

	claims, err := h.jwtService.VerifyToken(token)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(claims)
}

// func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
// 	token := c.Get("Authorization")
// 	if token == "" {
// 		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
// 	}

// 	newToken, err := h.jwtService.RefreshToken(c.UserContext(), token)
// 	if err != nil {
// 		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	return c.Status(http.StatusOK).JSON(fiber.Map{"token": newToken})
// }

// Login godoc
// @Summary      User login
// @Description  Authenticates user and returns JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      model.LoginBody  true  "Login credentials"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body model.LoginBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	token, err := h.service.Login(c.UserContext(), body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"token": token})
}

func (h *AuthHandler) ExistTwitterUID(c *fiber.Ctx) error {
	var twitterNameReq model.TwitterUIDReq
	if err := c.BodyParser(&twitterNameReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	ois, err := h.service.ExistTwitterUID(twitterNameReq.TwitterUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(ois)
}
