package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type CRMHandler struct {
	service port.CRMService
	crmRepo port.CRMRepo
}
type ErrorResponse struct {
	Error string `json:"error"`
}
type KolReferRequest struct {
	Page int `json:"page"`
}
type UpdateDisplayCodeRequest struct {
	CryptoUserID string `json:"crypto_user_id"`
	DisplayCode  string `json:"display_code"`
}

// UpdateCryptoUserApproveRequest is the request body for updating copytrade approval status
type UpdateCryptoUserApproveRequest struct {
	UserUUID string `json:"uuid"`
	Approve  *bool  `json:"approve"`
}

func NewCryptoCRMHandler(service port.CRMService) *CRMHandler {
	return &CRMHandler{service: service}
}

// Login CRM
// @Summary      CRM User login
// @Description  Authenticates CRM user and returns JWT token
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Param        body  body      model.CRMLoginBody  true  "Login credentials"
// @Success      200   {object}  ErrorResponse   "Successful login, returns JWT token. Example: {\"token\": \"<jwt>\"}"
// @Failure      400   {object}  ErrorResponse   "Bad Request - Invalid request body"
// @Failure      401   {object}  ErrorResponse   "Unauthorized - Invalid credentials"
// @Router       /crm/login [post]
func (h *CRMHandler) Login(c *fiber.Ctx) error {
	var body model.CRMLoginBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	token, err := h.service.CRMLogin(c.UserContext(), body)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

// New Kol CRM
// @Summary      Create a new KOL user
// @Description  Creates a new KOL user and associates a KOL code. Requires valid CRM user and unique display code.
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Param        body  body      model.KOLUser  true  "KOL User Payload"
// @Success      200   "KOL user created successfully", e.g., {\"kolUserID\": \"abc123\"}"
// @Failure      400   {object}  ErrorResponse "Bad Request - Invalid input, validation failed, or duplicate KOL user"
// @Failure      404   {object}  ErrorResponse "Not Found - CRM user does not exist"
// @Failure      500   {object}  ErrorResponse "Internal Server Error - Validation or DB operation failed"
// @Router       /crm/new-kol [post]
// @Security     BearerAuth
func (h *CRMHandler) NewKolUser(c *fiber.Ctx) error {
	var body model.KOLUser
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	exit, err := h.service.CheckCRMUserIsExit(c.UserContext(), string(body.ID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot validate user",
		})
	}
	if !exit {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User is not exit",
		})
	}
	kolexit, err := h.service.CheckKOLUserIsExit(c.UserContext(), string(body.ID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot validate kol user",
		})
	}
	if kolexit {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Kol user already exit",
		})
	}
	validate, err := h.service.ValidateDisplaycode(c.UserContext(), string(body.DisplayCode))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if !validate {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot validate kol code",
		})
	}
	kolUserID, err := h.service.InsertKolUser(c.UserContext(), body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot insert new KOL user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"kolUserID": kolUserID,
	})
}

// KolReferDetail godoc
// @Summary      Get KOL referral details
// @Description  Fetches KOL referral details with pagination
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Success      200   {array}   map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /crm/kol_refer_detail [get]
// @Security     BearerAuth
func (h *CRMHandler) KolReferDetail(c *fiber.Ctx) error {
	// var requestBody map[string]int
	// if err := c.BodyParser(&requestBody); err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "Invalid request body",
	// 	})
	// }

	// page, exists := requestBody["page"]
	// if !exists || page < 1 {
	// 	page = 1
	// }

	results, err := h.service.GetKolReferDetail(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch KOL referral details",
		})
	}

	return c.Status(fiber.StatusOK).JSON(results)
}

// GetAllUsers godoc
// @Summary      Get all users
// @Description  Returns all users with uuid, email, and twitter_uid
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Success      200   {array}   map[string]interface{}
// @Failure      500   {object}  map[string]string
// @Router       /crm/get_users [get]
// @Security     BearerAuth
func (h *CRMHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}
	resp := make([]map[string]interface{}, len(users))
	for i, u := range users {
		resp[i] = map[string]interface{}{
			"uuid":         u.UUID,
			"email":        u.Email,
			"twitter_name": u.TwitterName,
		}
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

// UpdateDisplayCode godoc
// @Summary      Update display_code for a user
// @Description  Updates the display_code for a user if not duplicated
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Param        body  body  UpdateDisplayCodeRequest  true  "Update display_code"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /crm/update_display_code [post]
// @Security     BearerAuth
func (h *CRMHandler) UpdateDisplayCode(c *fiber.Ctx) error {
	var req struct {
		CryptoUserID string `json:"crypto_user_id"`
		DisplayCode  string `json:"display_code"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	if req.CryptoUserID == "" || req.DisplayCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "crypto_user_id and display_code are required"})
	}
	err := h.service.UpdateDisplayCode(c.UserContext(), req.CryptoUserID, req.DisplayCode)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

// GetCryptoUser godoc
// @Summary      Get crypto user data
// @Description  Get crypto user data
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Param        page  query     int   false  "Page number (starting from 1)"  default(1)
// @Param        order query     string false  "Sort by waiting_list_timestamp (asc or desc)" Enums(asc,desc) default(desc)
// @Param        search query    string false  "Search by twitter_name or email"
// @Param        is_copytrade_approved query bool   false  "Filter by copytrade approval (true/false)"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /crm/get-crypto-user [get]
// @Security     BearerAuth
func (h *CRMHandler) GetCryptoUser(c *fiber.Ctx) error {
	// parse page query param, default to 1
	pageStr := c.Query("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// parse order (asc/desc), default desc
	order := c.Query("order", "desc")
	if order != "asc" && order != "ASC" && order != "desc" && order != "DESC" {
		order = "desc"
	}

	// parse search
	search := c.Query("search", "")

	// parse optional is_copytrade_approved
	var approvedPtr *bool
	if v := c.Query("is_copytrade_approved"); v != "" {
		// Fiber doesn't have direct bool query parse, use strconv
		b, err := strconv.ParseBool(v)
		if err == nil {
			approvedPtr = &b
		}
	}

	users, err := h.service.GetCryptoUser(c.UserContext(), page, order, search, approvedPtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch crypto users",
		})
	}

	// get total count for pagination
	total, err := h.service.CountCryptoUser(c.UserContext(), search, approvedPtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count crypto users",
		})
	}

	pageSize := 200
	totalPages := (total + pageSize - 1) / pageSize

	return c.JSON(fiber.Map{
		"page":       page,
		"pageSize":   pageSize,
		"total":      total,
		"totalPages": totalPages,
		"order":      order,
		"items":      users,
	})
}

// GetRefferalScore godoc
// @Summary      Get refferal score
// @Description  Get refferal score
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /crm/get_refferal_score [get]
// @Security     BearerAuth
func (h *CRMHandler) GetRefferalScore(c *fiber.Ctx) error {
	users, err := h.service.GetRefferalScore(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch crypto users",
		})
	}
	return c.JSON(users)
}

// GetUserReferral godoc
// @Summary      Get User Referral
// @Description  Get User Referral
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /crm/get_user_referral [get]
// @Security     BearerAuth
func (h *CRMHandler) GetUserReferral(c *fiber.Ctx) error {
	users, err := h.service.GetUserReferral(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch crypto users",
		})
	}
	return c.JSON(users)
}

// UpdateCryptoUserApprove godoc
// @Summary      Update copytrade approval status
// @Description  Update a user's copytrade approval (waiting list) status
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Param        body  body      UpdateCryptoUserApproveRequest  true  "Approve payload"
// @Success      200   {object}  map[string]bool
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /crm/update-crypto-user-approve [post]
// @Security     BearerAuth
func (h *CRMHandler) UpdateCryptoUserApprove(c *fiber.Ctx) error {
	var req UpdateCryptoUserApproveRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	if req.UserUUID == "" || req.Approve == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_uuid and approve are required"})
	}
	err := h.service.UpdateUserApprove(c.UserContext(), req.UserUUID, *req.Approve)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed Approve User",
		})
	}
	return c.JSON(fiber.Map{"success": true})
}

// GetPrivyUserOverview godoc
// @Summary      Search privy copytrade user overview by user id
// @Description  Returns user profile, wallets (with subscriptions), and trade logs. Filter logs by status and order by execution date.
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Param        user_id query     string true  "User ID (crypto_user.id)"
// @Param        status query     string false "Filter trade_logs.status (e.g. success, fail)"
// @Param        order  query     string false "Order by executed_at (asc or desc)" Enums(asc,desc) default(desc)
// @Param        limit  query     int    false "Max logs to return" default(200)
// @Param        offset query     int    false "Logs offset for pagination" default(0)
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /crm/privy-user-overview [get]
// @Security     BearerAuth
func (h *CRMHandler) GetPrivyUserOverview(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		userID = c.Query("uuid")
	}
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id is required"})
	}

	// optional filters
	var statusPtr *string
	if v := c.Query("status"); v != "" {
		statusPtr = &v
	}
	order := c.Query("order", "desc")
	if order != "asc" && order != "ASC" && order != "desc" && order != "DESC" {
		order = "desc"
	}
	limit := 200
	if s := c.Query("limit"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			limit = v
		}
	}
	offset := 0
	if s := c.Query("offset"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			offset = v
		}
	}

	overview, err := h.service.GetPrivyUserOverview(c.UserContext(), userID, statusPtr, order, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(overview)
}
