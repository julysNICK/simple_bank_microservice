package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)
type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	
}

type renewAccessTokenResponse struct {
	AccessToken string       `json:"acess_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`

}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}


	refreshToken, err:=server.tokenMaker.VerifyToken(req.RefreshToken)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
	}

	session, err := server.store.GetSession(ctx, refreshToken.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	if session.IsBlocked {
		err := fmt.Errorf("refresh token is blocked")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != refreshToken.Username {
		err := fmt.Errorf("refresh token was issued to another user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("refresh token doesn't match the records")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("refresh token was expired")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}




	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(refreshToken.Username, server.config.AccessTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := renewAccessTokenResponse{
		
		AccessToken: accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpiredAt,
	
	}

	ctx.JSON(http.StatusOK, rsp)

}
