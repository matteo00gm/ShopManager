package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"shopsweb.com/gateway-service/utils"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}

	parsedToken, err := utils.VerifyToken(token)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}

	tokenParams, err := utils.GetParamsFromToken(parsedToken, "userId")

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}

	//using 0 as i'm retrieving userId only
	userId := int64(tokenParams[0].(float64))

	//setting the userId in the context for ownership checks in update and delete
	context.Set("userId", userId)
	context.Next()
}
