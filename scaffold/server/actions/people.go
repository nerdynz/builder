package actions

import (
	"net/http"

	"github.com/nerdynz/builder/scaffold/server/models"

	"github.com/nerdynz/datastore"
	flow "github.com/nerdynz/flow"
)

func NewPerson(w http.ResponseWriter, req *http.Request, ctx *flow.Context, store *datastore.Datastore) {
	personHelper := models.PersonHelper()
	person := personHelper.New()
	ctx.JSON(http.StatusOK, person)
}

func CreatePerson(w http.ResponseWriter, req *http.Request, ctx *flow.Context, store *datastore.Datastore) {
	personHelper := models.PersonHelper()
	person, err := personHelper.FromRequest(ctx.Req)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to create person record", nil)
		return
	}
	err = personHelper.Save(person)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to create person record", err)
		return
	}
	ctx.JSON(http.StatusOK, person)
}

func RetrievePerson(w http.ResponseWriter, req *http.Request, ctx *flow.Context, store *datastore.Datastore) {
	if ctx.URLParam("personID") == "" {
		RetrievePeople(w, req, ctx, store)
		return
	}

	//get the personID from the request
	personID := ctx.URLIntParamWithDefault("personID", -1)
	if personID == -1 {
		ctx.ErrorJSON(http.StatusInternalServerError, "Invalid personID", nil)
		return
	}

	personHelper := models.PersonHelper()
	person, err := personHelper.Load(personID)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to retrieve Person record", err)
		return
	}

	ctx.JSON(http.StatusOK, person)
}

func RetrievePeople(w http.ResponseWriter, req *http.Request, ctx *flow.Context, store *datastore.Datastore) {
	personHelper := models.PersonHelper()
	people, err := personHelper.All()

	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to retrieve Person records", err)
		return
	}

	ctx.JSON(http.StatusOK, people)
}

func PagedPeople(w http.ResponseWriter, req *http.Request, ctx *flow.Context, store *datastore.Datastore) {
	personHelper := models.PersonHelper()
	pageNum := ctx.URLIntParamWithDefault("pagenum", 1)
	limit := ctx.URLIntParamWithDefault("limit", 10)
	sort := ctx.URLParam("sort")
	direction := ctx.URLParam("direction")

	data, err := personHelper.PagedBy(pageNum, limit, sort, direction)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Unabled to get paged Person data", err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func UpdatePerson(w http.ResponseWriter, req *http.Request, ctx *flow.Context, store *datastore.Datastore) {
	personHelper := models.PersonHelper()
	person, err := personHelper.FromRequest(ctx.Req)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to load Person record for update", err)
		return
	}

	// save and validate
	err = personHelper.Save(person)
	// other type of error
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to save updated Person record", err)
		return
	}

	ctx.JSON(http.StatusOK, person)
}

func DeletePerson(w http.ResponseWriter, req *http.Request, ctx *flow.Context, store *datastore.Datastore) {
	personHelper := models.PersonHelper()
	//get the personID from the request
	personID := ctx.URLIntParamWithDefault("personID", -1)
	if personID == -1 {
		ctx.JSON(http.StatusInternalServerError, "Invalid PersonID for remove")
		return
	}

	isDeleted, err := personHelper.Delete(personID)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to remove the Person record", err)
		return
	}
	ctx.JSON(http.StatusOK, isDeleted)
}
