package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	GlobalDataKey    = "AuthData"
	SessionUserIDKey = "AuthUserID"
)

// Login a user (stores in session)
func LoginUser(c *gin.Context, userID interface{}) error {
	session := sessions.Default(c)
	session.Set(SessionUserIDKey, userID)
	return session.Save()
}

// Get the stored user (from session)
func GetUserID(c *gin.Context) *interface{} {
	session := sessions.Default(c)
	if userID := session.Get(SessionUserIDKey); userID != nil {
		return &userID
	}
	return nil
}

// Logout a user, removing stored auth or clearing all session data
func LogoutUser(c *gin.Context, clearAll bool) error {
	session := sessions.Default(c)
	if clearAll {
		session.Clear()
	} else {
		session.Delete(SessionUserIDKey)
	}
	return session.Save()
}
