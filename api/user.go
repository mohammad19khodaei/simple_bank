package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/mohammad19khodaei/simple_bank/db/sqlc"
	"github.com/mohammad19khodaei/simple_bank/utils"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	Username          string             `json:"username"`
	FullName          string             `json:"full_name"`
	Email             string             `json:"email"`
	PasswordChangedAt pgtype.Timestamptz `json:"password_changed_at"`
	CreatedAt         pgtype.Timestamptz `json:"created_at"`
}

func (s *server) createUserHandler(ctx *gin.Context) {
	var request createUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, s.errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, s.errorResponse(err))
		return
	}

	user, err := s.store.CreateUser(ctx, db.CreateUserParams{
		Username:       request.Username,
		HashedPassword: hashedPassword,
		FullName:       request.FullName,
		Email:          request.Email,
	})

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				ctx.JSON(http.StatusConflict, s.errorResponse(err))
			}
			return
		}
		ctx.JSON(http.StatusInternalServerError, s.errorResponse(err))
		return
	}

	resp := createUserResponse(user)
	ctx.JSON(http.StatusCreated, resp)
}

type loginRequest struct {
	Username string `json:"username" biding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

func (s *server) login(ctx *gin.Context) {
	var request loginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, s.errorResponse(err))
		return
	}

	user, err := s.store.GetUser(ctx, request.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "username or password is incorrect",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, s.errorResponse(err))
		return
	}

	if !utils.IsHashPasswordValid(user.HashedPassword, request.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "username or password is incorrect",
		})
		return
	}

	token, err := s.tokenMaker.GenerateToken(user.Username, s.config.TokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, s.errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{
		AccessToken: token,
		User:        createUserResponse(user),
	})
}

func createUserResponse(user db.User) UserResponse {
	return UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}
