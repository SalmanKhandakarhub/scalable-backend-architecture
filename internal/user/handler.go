package user

import (
	"net/http"
	"strconv"

	"github.com/SalmanKhandakarhub/scalable-backend-architecture/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{
		service: service,
	}
}

// GetAll - GET(/api/v1/users)
func (h *Handler) GetAll(ctx *gin.Context) {
	var query PaginationQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.Error(
			ctx,
			http.StatusBadRequest,
			"Invalid quarey parameters",
			err.Error(),
		)
		return
	}

	if query.Page > 0 && query.PageSize > 0 {
		users, total, err := h.service.GetWithPagination(
			ctx.Request.Context(),
			query.Page,
			query.PageSize,
		)
		if err != nil {
			response.Error(
				ctx,
				http.StatusInternalServerError,
				"Failed to fetch users",
				err.Error(),
			)
			return
		}

		response.SuccessWithPagination(
			ctx,
			http.StatusOK,
			"Users fetched successfully",
			ToUserResponseList(users),
			query.Page,
			query.PageSize,
			int(total),
		)
		return
	}

	// Get all users without pagination
	users, err := h.service.GetAll(ctx.Request.Context())
	if err != nil {
		response.Error(
			ctx,
			http.StatusInternalServerError,
			"Failed to fetch users",
			err.Error(),
		)
		return
	}
	response.Success(
		ctx,
		http.StatusOK,
		"Users fetched successfully",
		ToUserResponseList(users),
	)
}

// GetById - GET(/api/v1/users/:id)
func (h *Handler) GetById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(
			ctx,
			http.StatusBadRequest,
			"Invalid user id",
			err.Error(),
		)
		return
	}

	user, err := h.service.GetById(ctx.Request.Context(), uint(id))
	if err != nil {
		response.Error(
			ctx,
			http.StatusNotFound,
			"User not found",
			err.Error(),
		)
		return
	}
	response.Success(
		ctx,
		http.StatusOK,
		"Users fetched successfully",
		ToUserResponse(user),
	)
}

// Create - POST /api/v1/users
func (h *Handler) Create(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(
			ctx,
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		return
	}

	user, err := h.service.Create(ctx.Request.Context(), req)
	if err != nil {
		response.Error(
			ctx,
			http.StatusBadRequest,
			"Failed to create user",
			err.Error(),
		)
		return
	}

	response.Success(
		ctx,
		http.StatusCreated,
		"User created successfully",
		ToUserResponse(user),
	)
}

// Update - PUT /api/v1/users/:id
func (h *Handler) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(
			ctx,
			http.StatusBadRequest,
			"Invalid user ID",
			err.Error(),
		)
		return
	}

	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(
			ctx,
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		return
	}

	user, err := h.service.Update(ctx.Request.Context(), uint(id), req)
	if err != nil {
		response.Error(
			ctx,
			http.StatusBadRequest,
			"Failed to update user",
			err.Error(),
		)
		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"User updated successfully",
		ToUserResponse(user),
	)
}

// Delete - DELETE /api/v1/users/:id
func (h *Handler) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(
			ctx,
			http.StatusBadRequest,
			"Invalid user ID",
			err.Error(),
		)
		return
	}

	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.Error(
			ctx,
			http.StatusNotFound,
			"Failed to delete user",
			err.Error(),
		)
		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"User deleted successfully",
		nil,
	)
}

// Login - POST /api/v1/auth/login
func (h *Handler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(
			ctx,
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		return
	}

	user, token, err := h.service.Login(ctx.Request.Context(), req)
	if err != nil {
		response.Error(
			ctx,
			http.StatusUnauthorized,
			"Login failed",
			err.Error())
		return
	}

	loginResponse := LoginResponse{
		Token: token,
		User:  ToUserResponse(user),
	}

	response.Success(ctx, http.StatusOK, "Login successful", loginResponse)
}

// ChangePassword - PUT /api/v1/users/:id/password
func (h *Handler) ChangePassword(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.Error(
			ctx,
			http.StatusBadRequest,
			"Invalid user ID",
			err.Error(),
		)
		return
	}

	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(
			ctx,
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		return
	}

	if err := h.service.ChangePassword(ctx.Request.Context(), uint(id), req); err != nil {
		response.Error(
			ctx,
			http.StatusBadRequest,
			"Failed to change password",
			err.Error(),
		)
		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"Password changed successfully",
		nil,
	)
}

// Search - GET /api/v1/users/search
func (h *Handler) Search(ctx *gin.Context) {
	var query SearchQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.Error(
			ctx,
			http.StatusBadRequest,
			"Invalid query parameters",
			err.Error(),
		)
		return
	}

	if query.Keyword == "" {
		response.Error(
			ctx,
			http.StatusBadRequest,
			"Keyword is required",
			nil,
		)
		return
	}

	// With pagination
	if query.Page > 0 && query.PageSize > 0 {
		users, total, err := h.service.SearchWithPagination(
			ctx.Request.Context(),
			query.Keyword,
			query.Page,
			query.PageSize,
		)
		if err != nil {
			response.Error(
				ctx,
				http.StatusInternalServerError,
				"Failed to search users",
				err.Error(),
			)
			return
		}

		response.SuccessWithPagination(
			ctx,
			http.StatusOK,
			"Search results",
			ToUserResponseList(users),
			query.Page,
			query.PageSize,
			int(total),
		)
		return
	}

	// Without pagination
	users, err := h.service.Search(ctx.Request.Context(), query.Keyword)
	if err != nil {
		response.Error(
			ctx,
			http.StatusInternalServerError,
			"Failed to search users",
			err.Error(),
		)
		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"Search results",
		ToUserResponseList(users),
	)
}

// GetStats - GET /api/v1/users/stats
func (h *Handler) GetStats(ctx *gin.Context) {
	stats, err := h.service.GetStats(ctx.Request.Context())
	if err != nil {
		response.Error(
			ctx,
			http.StatusInternalServerError,
			"Failed to fetch stats",
			err.Error(),
		)
		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"Stats fetched successfully",
		stats,
	)
}
