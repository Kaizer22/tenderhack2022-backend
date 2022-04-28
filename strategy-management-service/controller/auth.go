package controller

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"main/auth"
	"main/logging"
	"main/model/entity"
	"main/model/response"
	"main/repository"
	"net/http"
	"os"
	"strconv"
)

func NewAuthController(c context.Context, r repository.AccountRepository) *AuthController {
	return &AuthController{
		tokenManager: auth.TokenManager{},
		ctx:          c,
		accountRepo:  r,
	}
}

type AuthController struct {
	tokenManager auth.TokenManager
	ctx          context.Context
	accountRepo  repository.AccountRepository
}

// Login godoc
// @Summary            Login
// @Description    Login request, use to get access_token and refresh_token
// @Tags                          auth
// @Accept                        json
// @Produce                       json
// @Param               info  body            entity.AccountData  true  "Account info"
// @Success             200            {object}  response.LoginResponse    "Login"
// @Failure        400  {object}  utils.HTTPError
// @Failure        404  {object}  utils.HTTPError
// @Failure        500  {object}  utils.HTTPError
// @Router                        /login [post]
func (a *AuthController) Login(c *gin.Context) {
	var acc entity.AccountData
	if err := c.ShouldBindJSON(&acc); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json")
		return
	}
	account, err := a.accountRepo.FindByLogin(a.ctx, acc.Login)
	if err != nil || acc.Login != account.Login || acc.Password != account.Password {
		logging.ErrorFormat("Error finding account %s by login", acc.Login)
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}

	ts, err := a.tokenManager.CreateToken(account.Id, string(account.Role))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	c.JSON(http.StatusOK, tokens)
}

// Register godoc
// @Summary            Register
// @Description    Add new account
// @Tags                          auth
// @Accept                        json
// @Produce                       json
// @Param               info  body            entity.AccountData  true  "Account info"
// @Success             201            {object}  response.AccountCreated    "Register"
// @Failure        400  {object}  utils.HTTPError
// @Failure        404  {object}  utils.HTTPError
// @Failure        500  {object}  utils.HTTPError
// @Router                       /register [post]
func (a *AuthController) Register(c *gin.Context) {
	var acc entity.AccountData
	if err := c.ShouldBindJSON(&acc); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json")
		return
	}
	account, err := a.accountRepo.FindByLogin(a.ctx, acc.Login)
	if err == nil || account != nil {
		logging.ErrorFormat("This login %s is already presented", acc.Login)
		c.JSON(http.StatusUnauthorized, "Please provide login, which is not in use")
		return
	}
	profileId, id, err := a.accountRepo.Insert(a.ctx, acc)
	if err != nil {
		logging.ErrorFormat("Cannot add new account %+v", acc)
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusCreated, response.AccountCreated{
		Msg:       "Account created",
		Id:        id,
		ProfileId: profileId,
	})

}

// Logout godoc
// @Summary            Logout
// @Description    Log out from account
// @Tags                          auth
// @Accept                        json
// @Produce                       json
// @Param        Authorization  header    string  true  "Authentication header"
// @Success        201            {object}  response.AccountCreated    "Register"
// @Failure        400  {object}  utils.HTTPError
// @Failure        404  {object}  utils.HTTPError
// @Failure        500  {object}  utils.HTTPError
// @Router                        /logout [post]
func (a *AuthController) Logout(c *gin.Context) {
	//If metadata is passed and the tokens valid, delete them from the redis store
	metadata, _ := a.tokenManager.ExtractTokenMetadata(c.Request)
	if metadata != nil {
		//deleteErr := servers.HttpServer.RD.DeleteTokens(metadata)
		//if deleteErr != nil {
		//	c.JSON(http.StatusBadRequest, deleteErr.Error())
		//	return
		//}
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

func (a *AuthController) Refresh(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	refreshToken := mapToken["refresh_token"]

	//verify the token
	token, err := jwt.Parse(refreshToken,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("REFRESH_SECRET")), nil
		})
	//if there is an error, the token must have expired
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Refresh token expired")
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		//refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		//if !ok {
		//	c.JSON(http.StatusUnprocessableEntity, err)
		//	return
		//}
		accountId, roleOk := claims["account_id"].(string)
		if roleOk == false {
			c.JSON(http.StatusUnprocessableEntity, "unauthorized")
			return
		}
		//Delete the previous Refresh Token
		//delErr := servers.HttpServer.RD.DeleteRefresh(refreshUuid)
		//if delErr != nil { //if any goes wrong
		//	c.JSON(http.StatusUnauthorized, "unauthorized")
		//	return
		//}
		//Create new pairs of refresh and access tokens

		accountID, err := strconv.Atoi(accountId)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "accountId invalid")
			return
		}
		account, err := a.accountRepo.FindById(a.ctx, int64(accountID))
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "User's not found ")
		}

		ts, createErr := a.tokenManager.CreateToken(int64(accountID),
			string(account.Role))
		if createErr != nil {
			c.JSON(http.StatusForbidden, createErr.Error())
			return
		}
		//save the tokens metadata to redis
		//saveErr := servers.HttpServer.RD.CreateAuth(accountId, ts)
		//if saveErr != nil {
		//	c.JSON(http.StatusForbidden, saveErr.Error())
		//	return
		//}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusCreated, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, "refresh expired")
	}
}
