package handler

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type TwitterCryptoHandler struct {
	twitterService    port.TwitterCryptoService
	authorTierService port.AuthorTierService
}

func NewTwitterCryptoHandler(twitterService port.TwitterCryptoService, authorTierService port.AuthorTierService) *TwitterCryptoHandler {
	return &TwitterCryptoHandler{twitterService: twitterService, authorTierService: authorTierService}
}

func (h *TwitterCryptoHandler) GetAllSentiments(c *fiber.Ctx) error {
	sentiments, err := h.twitterService.GetAllSentiments(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(sentiments)
}

func (h *TwitterCryptoHandler) GetAllTweets(c *fiber.Ctx) error {
	tweets, err := h.twitterService.GetAllTweets(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(tweets)
}
func (h *TwitterCryptoHandler) GetAuthorWinrate(c *fiber.Ctx) error {
	selectedWinratePeriod := c.Query("selectedWinratePeriod", "overall")
	profiles, err := h.twitterService.GetAuthorWinrate(c.Context(), selectedWinratePeriod)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(profiles)
}
func (h *TwitterCryptoHandler) GetAuthorProfiles(c *fiber.Ctx) error {
	profiles, err := h.twitterService.GetAuthorProfiles(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(profiles)
}

func (h *TwitterCryptoHandler) GetPaginatedTweets(c *fiber.Ctx) error {
	start := c.QueryInt("start", 0)
	limit := c.QueryInt("limit", 100)
	sortBy := c.Query("sort_by", "tweet_created_at")
	sortOrder := c.Query("sort_order", "desc")

	tweets, total, err := h.twitterService.GetPaginatedTweets(c.Context(), start, limit, sortBy, sortOrder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":       tweets,
		"total":      total,
		"start":      start,
		"limit":      limit,
		"sort_by":    sortBy,
		"sort_order": sortOrder,
	})
}

func (h *TwitterCryptoHandler) GetPaginatedSentiments(c *fiber.Ctx) error {
	start := c.QueryInt("start", 0)
	limit := c.QueryInt("limit", 100)
	sortBy := c.Query("sort_by", "created_at")
	sortOrder := c.Query("sort_order", "desc")

	sentiments, total, err := h.twitterService.GetPaginatedSentiments(c.Context(), start, limit, sortBy, sortOrder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":       sentiments,
		"total":      total,
		"start":      start,
		"limit":      limit,
		"sort_by":    sortBy,
		"sort_order": sortOrder,
	})
}

func (h *TwitterCryptoHandler) GetTweetsWithSentiments(c *fiber.Ctx) error {
	start := c.QueryInt("start", 0)
	limit := c.QueryInt("limit", 100)
	sortBy := c.Query("sort_by", "tweet_created_at")
	sortOrder := c.Query("sort_order", "desc")

	tweets, total, err := h.twitterService.GetTweetsWithSentiments(c.Context(), start, limit, sortBy, sortOrder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":       tweets,
		"total":      total,
		"start":      start,
		"limit":      limit,
		"sort_by":    sortBy,
		"sort_order": sortOrder,
	})
}

func (h *TwitterCryptoHandler) GetTweetsWithSentimentAuthorSignal(c *fiber.Ctx) error {
	start := c.QueryInt("start", 0)
	limit := c.QueryInt("limit", 100)
	sortBy := c.Query("sort_by", "tweet_created_at")
	sortOrder := c.Query("sort_order", "desc")

	// Get authors from query parameter
	authorsParam := c.Query("authors", "")
	var authors []string
	if authorsParam != "" {
		authors = strings.Split(authorsParam, ",")
		for i, author := range authors {
			authors[i] = strings.TrimSpace(author)
		}
	}

	// Parse date parameters
	var fromDate, toDate *time.Time
	fromDateStr := c.Query("from_date", "")
	toDateStr := c.Query("to_date", "")

	if fromDateStr != "" {
		parsedFromDate, err := time.Parse("2006-01-02", fromDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid from_date format. Use YYYY-MM-DD",
			})
		}
		fromDate = &parsedFromDate
	}

	if toDateStr != "" {
		parsedToDate, err := time.Parse("2006-01-02", toDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid to_date format. Use YYYY-MM-DD",
			})
		}
		// Set to end of day
		parsedToDate = parsedToDate.Add(24*time.Hour - time.Second)
		toDate = &parsedToDate
	}

	tweets, total, err := h.twitterService.GetTweetsWithSentimentAuthorSignal(
		c.Context(),
		start,
		limit,
		sortBy,
		sortOrder,
		authors,
		fromDate,
		toDate,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":       tweets,
		"total":      total,
		"start":      start,
		"limit":      limit,
		"sort_by":    sortBy,
		"sort_order": sortOrder,
		"authors":    authors,
		"from_date":  fromDate,
		"to_date":    toDate,
	})
}

func (h *TwitterCryptoHandler) GetTweetsWithSentimentsAndAuthor(c *fiber.Ctx) error {
	start := c.QueryInt("start", 0)
	limit := c.QueryInt("limit", 100)
	sortBy := c.Query("sort_by", "tweet_created_at")
	sortOrder := c.Query("sort_order", "desc")
	searchTokenSymbolValue := c.Query("search", "")

	// Get authors from query parameter
	authorsParam := c.Query("authors", "")
	var authors []string
	if authorsParam != "" {
		authors = strings.Split(authorsParam, ",")
		for i, author := range authors {
			authors[i] = strings.TrimSpace(author)
		}
	}

	// Parse date parameters
	var fromDate, toDate *time.Time
	fromDateStr := c.Query("from_date", "")
	toDateStr := c.Query("to_date", "")

	if fromDateStr != "" {
		parsedFromDate, err := time.Parse("2006-01-02", fromDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid from_date format. Use YYYY-MM-DD",
			})
		}
		fromDate = &parsedFromDate
	}

	if toDateStr != "" {
		parsedToDate, err := time.Parse("2006-01-02", toDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid to_date format. Use YYYY-MM-DD",
			})
		}
		// Set to end of day
		parsedToDate = parsedToDate.Add(24*time.Hour - time.Second)
		toDate = &parsedToDate
	}

	tweets, total, err := h.twitterService.GetTweetsWithSentimentsAndAuthor(
		c.Context(),
		start,
		limit,
		sortBy,
		sortOrder,
		authors,
		fromDate,
		toDate,
		searchTokenSymbolValue,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":       tweets,
		"total":      total,
		"start":      start,
		"limit":      limit,
		"sort_by":    sortBy,
		"sort_order": sortOrder,
		"authors":    authors,
		"from_date":  fromDate,
		"to_date":    toDate,
		"search":     searchTokenSymbolValue,
	})
}

func (h *TwitterCryptoHandler) GetTweetsWithSentimentsAndTier(c *fiber.Ctx) error {
	start := c.QueryInt("start", 0)
	limit := c.QueryInt("limit", 100)
	sortBy := c.Query("sort_by", "tweet_created_at")
	sortOrder := c.Query("sort_order", "desc")

	// Get authors from query parameter
	tierParam := c.Query("tier", "")
	var authors []string

	if tierParam != "" {
		authorList, err := h.authorTierService.GetAuthorsByTier(tierParam)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		authors = make([]string, len(authorList))
		for i, author := range authorList {
			authors[i] = strings.TrimSpace(author)
		}
	}

	// Parse date parameters
	var fromDate, toDate *time.Time
	fromDateStr := c.Query("from_date", "")
	toDateStr := c.Query("to_date", "")

	if fromDateStr != "" {
		parsedFromDate, err := time.Parse("2006-01-02", fromDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid from_date format. Use YYYY-MM-DD",
			})
		}
		fromDate = &parsedFromDate
	}

	if toDateStr != "" {
		parsedToDate, err := time.Parse("2006-01-02", toDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid to_date format. Use YYYY-MM-DD",
			})
		}
		// Set to end of day
		parsedToDate = parsedToDate.Add(24*time.Hour - time.Second)
		toDate = &parsedToDate
	}

	tweets, total, err := h.twitterService.GetTweetsWithSentimentsAndTier(
		c.Context(),
		start,
		limit,
		sortBy,
		sortOrder,
		authors,
		fromDate,
		toDate,
	)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":       tweets,
		"total":      total,
		"start":      start,
		"limit":      limit,
		"sort_by":    sortBy,
		"sort_order": sortOrder,
		"from_date":  fromDate,
		"to_date":    toDate,
	})
}

func (h *TwitterCryptoHandler) GetSummaries(c *fiber.Ctx) error {
	start := c.QueryInt("start", 0)
	limit := c.QueryInt("limit", 100)
	sortBy := c.Query("sort_by", "date_time")
	sortOrder := c.Query("sort_order", "desc")

	// Parse date parameters (optional)
	var fromDate, toDate *time.Time
	fromDateStr := c.Query("from_date", "")
	toDateStr := c.Query("to_date", "")

	if fromDateStr != "" {
		parsedFromDate, err := time.Parse("2006-01-02", fromDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid from_date format. Use YYYY-MM-DD",
			})
		}
		fromDate = &parsedFromDate
	}

	if toDateStr != "" {
		parsedToDate, err := time.Parse("2006-01-02", toDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid to_date format. Use YYYY-MM-DD",
			})
		}
		// Set to end of day
		parsedToDate = parsedToDate.Add(24*time.Hour - time.Second)
		toDate = &parsedToDate
	}

	summaries, total, err := h.twitterService.GetSummaries(c.Context(), start, limit, sortBy, sortOrder, fromDate, toDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":       summaries,
		"total":      total,
		"start":      start,
		"limit":      limit,
		"sort_by":    sortBy,
		"sort_order": sortOrder,
		"from_date":  fromDate,
		"to_date":    toDate,
	})
}

func (h *TwitterCryptoHandler) GetBubbleSentiment(c *fiber.Ctx) error {
	data, err := h.twitterService.GetBubbleSentiment(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"bubble_sentiment": data,
	})
}

func (h *TwitterCryptoHandler) SearchTokenMentionSymbolsByAuthors(c *fiber.Ctx) error {
	symbol := c.Query("symbol", "")
	authors := c.Query("authors", "")
	tier := c.Query("tier", "")
	limit := c.QueryInt("limit")
	createdAt := c.Query("created_at", "")
	id := c.Query("id", "")
	selectedTimeRange := c.Query("selected_time_range", "")

	var authorsList []string

	if tier != "" {
		tierAuthorList, err := h.authorTierService.GetAuthorsByTier(tier)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		authorsList = tierAuthorList
	} else {
		if authors != "" {
			authorsList = strings.Split(authors, ",")
			for i, author := range authorsList {
				authorsList[i] = strings.TrimSpace(author)
			}
		}
	}

	tokens, nextCreatedAt, nextId, totalCount, err := h.twitterService.SearchTokenMentionSymbolsByAuthors(
		c.Context(),
		symbol,
		createdAt,
		id,
		limit,
		selectedTimeRange,
		authorsList,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := model.TokenMentionSymbolWithAuthorResponse{
		Tokens: tokens,
		NextCursor: model.NextCursor{
			CreatedAt: nextCreatedAt,
			ID:        nextId,
		},
		TotalCount: totalCount,
	}

	return c.JSON(response)
}

func (h *TwitterCryptoHandler) GetAllTiers(c *fiber.Ctx) error {
	tiers, err := h.authorTierService.GetAllTiers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(tiers)
}
