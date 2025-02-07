package delivery

import (
	"backend/domain"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"

	"backend/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type HTTPHandler struct {
	authService    *service.AuthService
	backendService *service.BackendService
}

func NewHTTPHandler(authService *service.AuthService, backendService *service.BackendService) *HTTPHandler {
	return &HTTPHandler{
		authService:    authService,
		backendService: backendService,
	}
}

func (h *HTTPHandler) Register(c echo.Context) error {
	ctx := c.Request().Context()
	var req domain.Account
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.authService.RegisterAccount(ctx, req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to register account"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Account created successfully"})
}

func (h *HTTPHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	var req domain.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	token, err := h.authService.Login(ctx, req.Login, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid login or password"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *HTTPHandler) GetContainers(c echo.Context) error {
	ctx := c.Request().Context()
	containers, err := h.backendService.GetAllContainers(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch containers"})
	}
	return c.JSON(http.StatusOK, containers)
}

func (h *HTTPHandler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing authorization header"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid authorization header"})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("mysecretkey")), nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		return next(c)
	}
}
