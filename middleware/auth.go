package middleware

import (
	"a21hc3NpZ25tZW50/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		tokenValue, err := ctx.Cookie("session_token")
		if err != nil {
			if ctx.GetHeader("Content-Type") == "application/json" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.NewErrorResponse("unauthorized"))
			}
			ctx.Redirect(http.StatusSeeOther, "login")
			return
		}

		claimsToken := &model.Claims{}
		token, err := jwt.ParseWithClaims(tokenValue, claimsToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return model.JwtKey, nil
		})

		if err != nil || !token.Valid {
			if err == jwt.ErrSignatureInvalid {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				ctx.Abort()
				return
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad req"})
			ctx.Abort()
			return
		}

		ctx.Set("email", claimsToken.Email)
		ctx.Next()
	})
}
