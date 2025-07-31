package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// errorResponse error response struct
// @Description error message
type errorResponse struct {
	Error any `json:"error"`
}

func (h *Handler) logError(c *gin.Context, err error) {
	method := c.Request.Method
	uri := c.Request.RequestURI

	h.log.Error(err.Error(), "method", method, "uri", uri)
}

func (h *Handler) errorResponse(c *gin.Context, status int, message any) {
	env := errorResponse{message}

	c.AbortWithStatusJSON(status, env)
}

func (h *Handler) serverErrorResponse(c *gin.Context, err error) {
	h.logError(c, err)

	message := "the server encountered a problem and could not process your request"
	h.errorResponse(c, http.StatusInternalServerError, message)
}

func (h *Handler) notFoundResponse(c *gin.Context) {
	message := "the requested resource could not be found"
	h.errorResponse(c, http.StatusNotFound, message)
}

func (h *Handler) methodNotAllowedResponse(c *gin.Context) {
	message := fmt.Sprintf("the requested method %s is not supported for this resourse", c.Request.Method)
	h.errorResponse(c, http.StatusMethodNotAllowed, message)
}

func (h *Handler) badRequestResponse(c *gin.Context, err error) {
	h.errorResponse(c, http.StatusBadRequest, err.Error())
}

func (h *Handler) editConflictResponse(c *gin.Context) {
	message := "unable to update the record due to an edit conflict, please try again"
	h.errorResponse(c, http.StatusConflict, message)
}

func (h *Handler) failedValidationResponse(c *gin.Context, errs map[string]string) {
	h.errorResponse(c, http.StatusUnprocessableEntity, errs)
}
