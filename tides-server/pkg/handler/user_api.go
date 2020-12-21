package handler

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/runtime/middleware"

	"tides-server/pkg/restapi/operations/user"

	"tides-server/pkg/config"
	"tides-server/pkg/models"
)

// RegisterUserHandler is API handler for /users/register POST
func RegisterUserHandler(params user.RegisterUserParams) middleware.Responder {
	body := params.ReqBody
	db := config.GetDB()
	var queryUser models.User
	db.Where("username = ?", body.Username).First(&queryUser)
	if queryUser.Username != "" {
		return user.NewRegisterUserBadRequest().WithPayload(&user.RegisterUserBadRequestBody{Message: "Username already used!"})
	}

	newUser := models.User{
		City:        body.City,
		CompanyName: body.CompanyName,
		Country:     body.Country,
		Email:       body.Email,
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		Password:    body.Password,
		Phone:       body.Phone,
		Position:    body.Position,
		Priority:    models.UserPriorityLow,
		Username:    body.Username,
	}

	err := db.Create(&newUser).Error
	if err != nil {
		return user.NewRegisterUserBadRequest()
	}

	res := &user.RegisterUserOKBodyUserInfo{
		City:        body.City,
		CompanyName: body.CompanyName,
		Country:     body.Country,
		Email:       body.Email,
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		Password:    body.Password,
		Phone:       body.Phone,
		Position:    body.Position,
		Priority:    models.UserPriorityLow,
		Username:    body.Username,
	}

	return user.NewRegisterUserOK().WithPayload(&user.RegisterUserOKBody{UserInfo: res})
}

// UserLoginHandler is API handler for /users/login POST
func UserLoginHandler(params user.UserLoginParams) middleware.Responder {
	body := params.ReqBody

	db := config.GetDB()
	var queryUser models.User
	db.Where("Username = ?", body.Username).First(&queryUser)
	if queryUser.Username == "" {
		return user.NewUserLoginUnauthorized()
	} else if queryUser.Password != body.Password {
		return user.NewUserLoginUnauthorized()
	}

	expirationTime := time.Now().Add(expireTime)
	claims := Claims{
		ID: queryUser.Model.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := config.GetConfig().SecretKey
	signedToken, _ := token.SignedString([]byte(secretKey))

	res := user.UserLoginOKBodyUserInfo{Priority: queryUser.Priority, Username: queryUser.Username}

	return user.NewUserLoginOK().WithPayload(&user.UserLoginOKBody{Token: signedToken, UserInfo: &res})
}

// GetUserProfileHandler is API handler for /users/profile GET
func GetUserProfileHandler(params user.GetUserProfileParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return user.NewGetUserProfileUnauthorized()
	}

	uid, _ := ParseUserIDFromToken(params.HTTPRequest)
	db := config.GetDB()
	var u models.User
	if db.Where("id = ?", uid).First(&u).RowsAffected == 0 {
		return user.NewGetUserProfileNotFound()
	}

	res := user.GetUserProfileOKBodyResults{
		City:        u.City,
		CompanyName: u.CompanyName,
		Country:     u.Country,
		Email:       u.Email,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Phone:       u.Phone,
		Position:    u.Position,
		Priority:    u.Priority,
		Username:    u.Username,
	}

	return user.NewGetUserProfileOK().WithPayload(&user.GetUserProfileOKBody{
		Message: "success",
		Results: &res,
	})
}

// UpdateUserProfileHandler is API handler for /users/profile PUT
func UpdateUserProfileHandler(params user.UpdateUserProfileParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return user.NewUpdateUserProfileUnauthorized()
	}

	uid, _ := ParseUserIDFromToken(params.HTTPRequest)
	body := params.ReqBody
	db := config.GetDB()
	var u models.User
	if db.Where("id = ?", uid).First(&u).RowsAffected == 0 {
		return user.NewUpdateUserProfileNotFound()
	}

	u.City = body.City
	u.CompanyName = body.CompanyName
	u.Country = body.Country
	u.Email = body.Email
	u.FirstName = body.FirstName
	u.LastName = body.LastName
	u.Phone = body.Phone
	u.Position = body.Position

	err := db.Save(&u).Error
	if err != nil {
		return user.NewUpdateUserProfileNotFound().WithPayload(&user.UpdateUserProfileNotFoundBody{
			Message: err.Error(),
		})
	}

	return user.NewUpdateUserProfileOK().WithPayload(&user.UpdateUserProfileOKBody{
		Message: "success",
	})
}
