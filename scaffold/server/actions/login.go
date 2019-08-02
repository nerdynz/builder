package actions

import (
	"net/http"

	"github.com/nerdynz/builder/scaffold/server/models"
	flow "github.com/nerdynz/flow"
)

func Login(ctx *flow.Context) {
	helper := models.PersonHelper()

	// create a blank person. dont load from request because we need to check their creds are valid
	person, err := helper.FromRequest(ctx.Req)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to read login details", err)
		return
	}

	sessionInfo, err := ctx.Padlock.LoginReturningInfo(person.Email, person.Password, "person")
	if err != nil {
		ctx.ErrorJSON(http.StatusUnauthorized, "Failed to login. Incorrect username or password", err)
		return
	}

	ctx.JSON(http.StatusOK, sessionInfo)
}

func RetrieveLoginUsers(ctx *flow.Context) {
	var people models.People
	err := ctx.Store.DB.
		Select("person_id", "email", "name", "picture").
		From("person").
		QueryStructs(&people)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, people)
}

func Logout(ctx *flow.Context) {
	ctx.Padlock.Logout()
	ctx.Redirect("/", http.StatusSeeOther)
}

func UserDetails(ctx *flow.Context) {
	// create a blank person. dont load from request because we need to check their creds are valid
	user, err := ctx.Padlock.LoggedInUser()
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Invalid user details", err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}
