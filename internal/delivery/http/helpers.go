package http

import (
	"eff-subscriptions/internal/domain/models"
	"eff-subscriptions/internal/validator"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

func readIDParam(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func readString(c *gin.Context, key string, defaultValue string) string {
	s := c.Query(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func readDate(c *gin.Context, key string, defaultValue models.CustomDate, v *validator.Validator) models.CustomDate {
	s := c.Query(key)

	if s == "" {
		return defaultValue
	}

	var date models.CustomDate
	s = strings.Trim(s, "\"")
	t, err := time.Parse("01-2006", s)
	if err != nil {
		v.AddError(key, "must be a valid date")
		return defaultValue
	}
	date = models.CustomDate(t)

	return date
}

func readInt(c *gin.Context, key string, defaultValue int, v *validator.Validator) int {
	s := c.Query(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer")
		return defaultValue
	}

	return i
}

func readUUID(c *gin.Context, key string, defaultValue uuid.UUID, v *validator.Validator) uuid.UUID {
	s := c.Query(key)

	if s == "" {
		return defaultValue
	}

	u, err := uuid.Parse(s)
	if err != nil {
		v.AddError(key, "must be a valid UUID")
		return defaultValue
	}

	return u
}
