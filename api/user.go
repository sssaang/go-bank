package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/sssaang/simplebank/db/sqlc"
	"github.com/sssaang/simplebank/db/util"
)


type createUserRequest struct {
	Username    string `json:"username" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, hashErr := util.HashPassword(req.Password)
	if hashErr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(hashErr))
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		HashedPassword: hashedPassword,
		FullName: req.FullName,
		Email: req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return 
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return 
	}

	res := createUserResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}

	ctx.JSON(http.StatusOK, res)
}

// type getUserRequest struct {
// 	Username int64 `uri:"id" binding:"required,min=1"`
// }

// func (server *Server) getUser(ctx *gin.Context) {
// 	var req getUserRequest
// 	if err := ctx.ShouldBindUri(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	account, err := server.store.GetAccount(ctx, req.ID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(err))
// 			return 
// 		}

// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return 
// 	}

// 	ctx.JSON(http.StatusOK, account)
// }