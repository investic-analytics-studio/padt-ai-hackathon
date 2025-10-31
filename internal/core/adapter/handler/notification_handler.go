package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/config"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
)

type CryptoNotificationHandler struct {
	service port.NotificationService
	config  *config.Config
}

func NewCryptoNotificationHandler(service port.NotificationService, config *config.Config) *CryptoNotificationHandler {
	return &CryptoNotificationHandler{
		service: service,
		config:  config,
	}
}

type UpdateGroupNameRequest struct {
	GroupID      string `json:"group_id"`
	NewGroupName string `json:"new_group_name"`
}
type AddGroup struct {
	GroupName string `json:"group_name"`
}
type AddAuthor struct {
	GroupID    string `json:"group_id"`
	AuthorName string `json:"author_name"`
}
type RemoveAuthor struct {
	GroupAutherID string `json:"group_auther_id"`
}
type TelegramUpsert struct {
	ChatID string `json:"chat_id"`
	UserID string `json:"user_id"`
}
type ChatMemberResponse struct {
	Ok          bool       `json:"ok"`
	Result      ChatMember `json:"result"`
	Description string     `json:"description,omitempty"`
}
type DeleteGroup struct {
	GroupID string `json:"group_id"`
}
type ChatMember struct {
	Status string `json:"status"`
	User   struct {
		ID        int64  `json:"id"`
		FirstName string `json:"first_name"`
		Username  string `json:"username"`
	} `json:"user"`
}

// Get Notification group list by uid godoc
// @Summary      Get Notification group list by uid
// @Description  Get Notification group list by uid
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Success      200   {array}   map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notification/get-group [get]
// @Security     BearerAuth
func (h *CryptoNotificationHandler) GetNotificationGroupList(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)

	notiGroup, err := h.service.GetNotificationGroupList(c.Context(), uid)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})

	}

	return c.Status(fiber.StatusOK).JSON(notiGroup)
}

// Update group name by group id godoc
// @Summary      Update group name by group id
// @Description  Update group name by group id
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Param        body  body  UpdateGroupNameRequest  true  "Update Group Name"
// @Success      200   {array}   map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notification/update-group-name [post]
// @Security     BearerAuth
func (h *CryptoNotificationHandler) UpdateGroupName(c *fiber.Ctx) error {
	req := UpdateGroupNameRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	err := h.service.UpdateGroupName(c.Context(), req.GroupID, req.NewGroupName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(req)
}

// Add group godoc
// @Summary      Add new group godoc
// @Description  Add new group
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Param        body  body  AddGroup  true  "Add New Group"
// @Success      200   {array}   map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notification/add-group [post]
// @Security     BearerAuth
func (h *CryptoNotificationHandler) AddGroup(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	req := AddGroup{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	err := h.service.AddGroup(c.Context(), uid, req.GroupName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(req)
}

// Count user's group godoc
// @Summary      Count user's group godoc
// @Description  Count user's group
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Success      200   {array}   map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notification/count-group [get]
// @Security     BearerAuth
func (h *CryptoNotificationHandler) CountGroup(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	groupCount, err := h.service.CountGroup(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(groupCount)
}

// Add new author to group godoc
// @Summary      add new author to group godoc
// @Description  add new author to group
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Param        body  body  AddAuthor  true  "Add New Author"
// @Success      200   {array}   map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notification/add-author [post]
// @Security     BearerAuth
func (h *CryptoNotificationHandler) AddAuthor(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	req := AddAuthor{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	err := h.service.AddAuthor(c.Context(), uid, req.GroupID, req.AuthorName)
	if err != nil {
		if err.Error() == "group not found or does not belong to the user" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "group not found or does not belong to the user",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

// Remove author from group godoc
// @Summary      remove author from group godoc
// @Description  remove author from group
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Param        body  body  RemoveAuthor  true  "Remove Author"
// @Success      200   {array}   map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notification/remove-author [post]
// @Security     BearerAuth
func (h *CryptoNotificationHandler) RemoveAuthor(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	req := RemoveAuthor{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	err := h.service.RemoveAuthor(c.Context(), uid, req.GroupAutherID)
	if err != nil {
		if err.Error() == "author not found or does not belong to the user" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "author not found or does not belong to the user",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

// Update user telegram data godoc
// @Summary      Update user telegram data godoc
// @Description  Update user telegram data if telegram user id and chat id are valid
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Param        body  body  TelegramUpsert  true  "Update Telegram Data"
// @Success      200   {array}   map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notification/update-telegram [post]
// @Security     BearerAuth
func (h *CryptoNotificationHandler) UpdateTelegram(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	var req TelegramUpsert
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	botToken := config.GetConfig().Telegram.BotToken

	// Parse chatID and userID from string to int64
	chatID, err := strconv.ParseInt(req.ChatID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid chat_id or user_id"})
	}

	userID, err := strconv.ParseInt(req.UserID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid chat_id or user_id"})
	}

	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/getChatMember?chat_id=%d&user_id=%d",
		botToken, chatID, userID,
	)

	resp, err := http.Get(url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Telegram request failed"})
	}
	defer resp.Body.Close()

	var result ChatMemberResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode Telegram response"})
	}

	if !result.Ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid chat_id or user_id"})
	}
	err = h.service.UpdateTelegram(c.Context(), uid, req.ChatID, req.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id": req.UserID,
		"chat_id": req.ChatID,
		"status":  result.Result.Status,
	})

}

// get user telegram data godoc
// @Summary      Get user telegram data godoc
// @Description  Get user telegram data
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Success      200   {array}   map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notification/get-telegram [get]
// @Security     BearerAuth
func (h *CryptoNotificationHandler) GetTelegram(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	telegram, err := h.service.GetTelegram(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(telegram)
}

// Delete all author in group godoc
// @Summary      Delete all author in group godoc
// @Description  Delete all author in group
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Param        body  body  DeleteGroup  true  "Delete all author in group"
// @Success      200   {array}   map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notification/delete-group-author [post]
// @Security     BearerAuth
func (h *CryptoNotificationHandler) DeleteGroupAuthor(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	req := DeleteGroup{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	err := h.service.DeleteGroupAuthor(c.Context(), uid, req.GroupID)
	if err != nil {
		if err.Error() == "group not found or does not belong to the user" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "group not found or does not belong to the user",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

// Delete telegram data, group id and author godoc
// @Summary      Delete telegram data, group id and author godoc
// @Description  Delete telegram data, group id and author
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Success      200   {array}   map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notification/disconnect [get]
// @Security     BearerAuth
func (h *CryptoNotificationHandler) DisconnectNotification(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	err := h.service.DisconnectNotification(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
