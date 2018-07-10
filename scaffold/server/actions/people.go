package actions

import (
	"net/http"

	"repo.nerdy.co.nz/displayworks/displayworks-signs/server/models"

	flow "github.com/nerdynz/flow"
)

func NewPerson(ctx *flow.Context) {
	personHelper := models.PersonHelper()
	person := personHelper.New()
	ctx.JSON(http.StatusOK, person)
}

func CreatePerson(ctx *flow.Context) {
	personHelper := models.PersonHelper()
	person, err := personHelper.NewFromRequest(ctx.Req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	err = personHelper.Save(person)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, person)
}

func RetrievePerson(ctx *flow.Context) {
	if ctx.URLParam("personID") == "" {
		RetrievePeople(ctx)
		return
	}

	//get the personID from the request
	personID, err := ctx.URLIntParam("personID")
	if err != nil || personID == 0 {
		ctx.JSON(http.StatusInternalServerError, "invalid personID")
		return
	}

	personHelper := models.PersonHelper()
	person, err := personHelper.Load(personID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, person)
}

func RetrievePeople(ctx *flow.Context) {
	personHelper := models.PersonHelper()
	people, err := personHelper.All()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, people)
}

func UpdatePerson(ctx *flow.Context) {
	personHelper := models.PersonHelper()
	person, err := personHelper.LoadAndUpdateFromRequest(ctx.Req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// save and validate
	err = personHelper.Save(person)
	// other type of error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, person)
}

func DeletePerson(ctx *flow.Context) {
	personHelper := models.PersonHelper()
	//get the personID from the request
	personID, err := ctx.URLIntParam("personID")
	if err != nil || personID == 0 {
		ctx.JSON(http.StatusInternalServerError, "invalid personID")
		return
	}

	isDeleted, err := personHelper.Delete(personID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, isDeleted)
}
