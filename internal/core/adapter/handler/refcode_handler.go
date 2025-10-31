package handler

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type GenerateRefcodeRequest struct {
	CryptoUserID string `json:"crypto_user_id" validate:"required"`
}
type GenerateRefcodeResponse struct {
	Refcodes []string `json:"refcodes"`
}
type CheckUserIDExistsResponse struct {
	Result bool `json:"result"`
}
type RefcodeHandler struct {
	service  port.CryptoUserRefcodeService
	authRepo port.AuthRepo
}
type CheckRefcodeRequest struct {
	// UserID  string `json:"user_id" validate:"required"`
	Refcode string `json:"ref_code" validate:"required"`
}

type CheckRefcodeResponse struct {
	Result bool `json:"result"`
}
type GetCryptoRefUserRequest struct {
	CryptoUserID string `json:"crypto_user_id" validate:"required"`
}
type OffsetDaysRequest struct {
	OffsetDays int `json:"offset_days" validate:"required"`
}
type CheckXUser struct {
	OffsetDays int `json:"offset_days" validate:"required"`
}
type GetCryptoRefUserResponse map[string]string

func NewCryptoUserRefcodeHandler(service port.CryptoUserRefcodeService, authRepo port.AuthRepo) *RefcodeHandler {
	return &RefcodeHandler{service: service, authRepo: authRepo}
}

// GenerateRefcodeRequest godoc
// @Summary      Generate a new referral code
// @Description  Generates and saves a referral code for the provided Crypto User ID
// @Tags         refcode
// @Accept       json
// @Produce      json
// @Param        body  body      GenerateRefcodeRequest  true  "Crypto User ID to generate refcode"
// @Success      200   {object}  GenerateRefcodeResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /generate-refcode [post]
// @Security     BearerAuth
func (h *RefcodeHandler) GenerateRefcodeRequest(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	refcodes, err := h.service.GenerateAndSaveRefcode(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(GenerateRefcodeResponse{
		Refcodes: refcodes,
	})

}

// @Summary Generate referral codes by number
// @Description Generates a specific number of referral codes for a given user ID
// @Tags refcode
// @Accept json
// @Produce json
// @Param uid query string true "User ID"
// @Param genNum query int true "Number of referral codes to generate"
// @Success 200 {object} map[string][]string "refcodes"
// @Failure 400 {object} map[string]string "error"
// @Failure 404 {object} map[string]string "error"
// @Router /generate-refcode-bynum [post]
// @Security     BearerAuth
func (h *RefcodeHandler) GenerateRefcodeBynumRequest(c *fiber.Ctx) error {
	uid := c.Query("uid")
	genNumStr := c.Query("genNum")

	if uid == "" || genNumStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required query parameters: uid or genNum",
		})
	}
	genNum, err := strconv.Atoi(genNumStr)
	if err != nil || genNum <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid genNum value",
		})
	}
	refcodes, err := h.service.GenerateRefcodeBynumRequest(c.Context(), uid, genNum)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(GenerateRefcodeResponse{
		Refcodes: refcodes,
	})
}

// CheckUserIDExists godoc
// @Summary      Check if User ID exists
// @Description  Verifies if the given User ID is already registered in the system
// @Tags         refcode
// @Accept       json
// @Produce      json
// @Success      200 {object} CheckUserIDExistsResponse
// @Failure      500 {object} map[string]string
// @Router       /check-user-id [get]
// @Security     BearerAuth
func (h *RefcodeHandler) CheckUserIDExists(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	exists, err := h.service.CheckUserIDExists(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to check user ID",
		})
	}

	return c.JSON(CheckUserIDExistsResponse{
		Result: !exists, // Return true if user doesn't exist, false if exists
	})
}

// CheckAndUpdateRefcode godoc
// @Summary      Check and update a referral code
// @Description  Checks if a referral code is valid and not yet used, then updates it with the user's UUID if available
// @Tags         refcode
// @Accept       json
// @Produce      json
// @Param        body  body      CheckRefcodeRequest  true  "Referral code request with code and user ID"
// @Success      200   {object}  CheckRefcodeResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /check-refcode [post]
// @Security     BearerAuth
func (h *RefcodeHandler) CheckAndUpdateRefcode(c *fiber.Ctx) error {
	var req CheckRefcodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	uid := c.Locals("uid").(string)

	isAvailable, err := h.service.CheckAndUpdateRefcode(c.Context(), req.Refcode, uid)
	if err != nil {
		// c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 	"error": err.Error(),
		// })
		return c.Next()
	}
	return c.JSON(CheckRefcodeResponse{
		Result: isAvailable,
	})

}
func (h *RefcodeHandler) CheckAndInsertKolcode(c *fiber.Ctx) error {
	var req CheckRefcodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	uid := c.Locals("uid").(string)
	KolCodeEntry, err := h.service.CheckKolcode(c.Context(), req.Refcode, uid)
	if err != nil {
		// handle error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if KolCodeEntry.CryptoUserID != "" && KolCodeEntry.DisplayCode != "" && KolCodeEntry.Refcode != "" {
		err := h.service.InsertKolcode(c.Context(), KolCodeEntry.CryptoUserID, KolCodeEntry.Refcode, uid)
		if err != nil {
			// handle error
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}
	return c.JSON(CheckRefcodeResponse{
		Result: true,
	})
}

// GetCryptoRefUser godoc
// @Summary      Get Crypto Referral User details
// @Description  Fetches refcodes and their corresponding user's email who used the refcode
// @Tags         refcode
// @Accept       json
// @Produce      json
// @Success      200   {object}  GetCryptoRefUserResponse
// @Failure      500   {object}  map[string]string
// @Router       /get-crypto-ref-user [get]
// @Security     BearerAuth
func (h *RefcodeHandler) GetCryptoRefUser(c *fiber.Ctx) error {

	uid := c.Locals("uid").(string)

	refcodes, err := h.service.GetByCryptoUserID(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to get crypto ref user",
		})
	}

	// Create a map to store all refcodes and their corresponding emails
	response := GetCryptoRefUserResponse{}

	// Add all refcodes to the map
	for _, refcode := range refcodes {
		email := ""
		twitterID := ""
		if refcode.CryptoRefUserID != nil {
			// Get user info to fetch email
			// fmt.Println("CryptoRefUserID:", *refcode.CryptoRefUserID)
			userInfo, err := h.authRepo.GetUserInfo(*refcode.CryptoRefUserID)
			// fmt.Println(userInfo)
			if err == nil {
				email = userInfo.Email
				twitterID = userInfo.TwitterName
			}
		}

		if email == "" {
			response[refcode.Refcode] = twitterID
		} else {
			response[refcode.Refcode] = email
		}
	}

	return c.JSON(response)
}

// GetCryptoRefUser godoc
// @Summary      Get KOL code and number of uses.
// @Description  Fetches KOL code and number of uses."
// @Tags         refcode
// @Accept       json
// @Produce      json
// @Success      200   {object}  GetCryptoRefUserResponse
// @Failure      500   {object}  map[string]string
// @Router       /get-kolcode [get]
// @Security     BearerAuth
func (h *RefcodeHandler) GetCryptoKolCode(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	kolUsed, err := h.service.GetCryptoKolCode(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to get kol used data",
		})
	}
	return c.JSON(kolUsed)
}

// GetRefferalScore godoc
// @Summary      Get users refferal score.
// @Description  Get users refferal score."
// @Tags         refcode
// @Accept       json
// @Produce      json
// @Success      200   {object}  GetCryptoRefUserResponse
// @Failure      500   {object}  map[string]string
// @Router       /get-refferal-score [get]
func (h *RefcodeHandler) GetRefferalScore(c *fiber.Ctx) error {
	refScore, err := h.service.GetRefferalScore(c.Context())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to get refferal score",
		})
	}
	return c.JSON(refScore)
}

// GetRefferalScoreRanking godoc
// @Summary      Get referral score rankings with rank changes
// @Description  Returns referral score rankings with rank changes compared to N days ago.
//
// Response Codes:
// 200 OK: Successfully retrieved ranking data
// 400 Bad Request: Invalid request parameters
// 500 Internal Server Error: Server-side error occurred
//
// Example Request:
//
//	{
//	  "offset_days": 7
//	}
//
// @Tags         refcode
// @Accept       json
// @Produce      json
// @Param        body  body  OffsetDaysRequest  true  "Number of days to look back for rank comparison"  example(7)
// @Success      200   {array}   model.RefferalScoreRanking  "Successfully retrieved ranking data"
// @Failure      400   {object}  map[string]string           "Invalid request parameters"
// @Failure      500   {object}  map[string]string           "Server-side error occurred"
// @Router       /get-refferal-score-ranking [post]
func (h *RefcodeHandler) GetRefferalScoreRanking(c *fiber.Ctx) error {
	var req OffsetDaysRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	offsetDays := req.OffsetDays
	if offsetDays < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "offset_days must be a positive integer",
		})
	}

	rankings, err := h.service.GetRefferalScoreRanking(c.Context(), offsetDays)
	if err != nil {
		// Log the error (replace with your logger if you have one)
		fmt.Printf("GetRefferalScoreRanking error: %+v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get referral score rankings: %v", err),
		})
	}

	return c.JSON(rankings)
}

// CheckXUserIsExit godoc
// @Summary      Check X user is exit
// @Description  Check X user and telegram alert are exit
// @Tags         refcode
// @Accept       json
// @Produce      json
// @Param        xUser body model.XUser true "XUser"
// @example xUser {"twitter_name": "NoMoonNoBuy"}
// @Success      200 {object} model.CheckXUser
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /check-xuser [post]
func (h *RefcodeHandler) CheckXUserIsExit(c *fiber.Ctx) error {
	var req model.XUser
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	CheckXUser, err := h.service.CheckXUserIsExit(c.Context(), req.TwitterName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check x user",
		})
	}
	return c.JSON(CheckXUser)
}
