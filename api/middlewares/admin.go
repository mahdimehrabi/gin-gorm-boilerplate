package middlewares

import (
	"boilerplate/api/repositories"
	"boilerplate/api/responses"
	"boilerplate/api/services"
	"boilerplate/infrastructure"
	"net/http"

	"github.com/gin-gonic/gin"
)

//AdminMiddleware -> struct for transaction
type AdminMiddleware struct {
	logger         infrastructure.Logger
	authService    services.AuthService
	env            infrastructure.Env
	userRepository repositories.UserRepository
}

//NewAdminMiddleware -> new instance of transaction
func NewAdminMiddleware(
	logger infrastructure.Logger,
	authService services.AuthService,
	env infrastructure.Env,
	userRepository repositories.UserRepository,
) AdminMiddleware {
	return AdminMiddleware{
		authService:    authService,
		logger:         logger,
		env:            env,
		userRepository: userRepository,
	}
}

func (m AdminMiddleware) AdminHandle() gin.HandlerFunc {
	m.logger.Zap.Info("setting up admin middleware")

	return func(c *gin.Context) {
		user, err := m.userRepository.GetAuthenticatedUser(c)
		if err != nil {
			m.logger.Zap.Error("Failed to get user in admin middleware", err.Error())
			responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured 😢")
			c.Abort()
			return
		}
		if !user.IsAdmin {
			responses.ErrorJSON(c, http.StatusForbidden, gin.H{}, "You don't have access to this page 😥")
			c.Abort()
			return
		}
		c.Next()
	}
}
