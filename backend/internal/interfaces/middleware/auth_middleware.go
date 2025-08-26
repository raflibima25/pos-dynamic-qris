package middleware

import (
	"strings"

	"qris-pos-backend/internal/domain/entities"
	"qris-pos-backend/pkg/auth"
	"qris-pos-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtService *auth.JWTService
}

func NewAuthMiddleware(jwtService *auth.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Check if header starts with Bearer
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := tokenParts[1]
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireRole(allowedRoles ...entities.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First check authentication
		m.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		userRole, exists := c.Get("user_role")
		if !exists {
			response.Forbidden(c, "User role not found")
			c.Abort()
			return
		}

		role, ok := userRole.(entities.UserRole)
		if !ok {
			response.Forbidden(c, "Invalid user role")
			c.Abort()
			return
		}

		// Check if user role is in allowed roles
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "Insufficient permissions")
		c.Abort()
	}
}

func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole(entities.RoleAdmin)
}

func (m *AuthMiddleware) RequireAdminOrCashier() gin.HandlerFunc {
	return m.RequireRole(entities.RoleAdmin, entities.RoleCashier)
}

// Optional auth - doesn't block request if no token
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		token := tokenParts[1]
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		// Set user info in context if valid
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

// Helper function to get current user from context
func GetCurrentUser(c *gin.Context) (*auth.Claims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}

	userClaims, ok := claims.(*auth.Claims)
	return userClaims, ok
}