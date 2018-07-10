package actions

import (
	"net/http"

	"github.com/nerdynz/flow"
	"repo.nerdy.co.nz/displayworks/displayworks-signs/server/models"
)

func NewSettings(ctx *flow.Context) {
	settingsHelper := models.SettingsHelper()
	settings := settingsHelper.New()
	ctx.JSON(http.StatusOK, settings)
}

func CreateSettings(ctx *flow.Context) {
	settingsHelper := models.SettingsHelper()
	settings, err := settingsHelper.NewFromRequest(ctx.Req)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to create settings record", nil)
		return
	}
	err = settingsHelper.Save(settings)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to create settings record", err)
		return
	}
	ctx.JSON(http.StatusOK, settings)
}

func RetrieveSettings(ctx *flow.Context) {
	//get the settingsID from the request
	settingsID := ctx.URLIntParamWithDefault("settingsID", -1)
	if settingsID == -1 {
		ctx.ErrorJSON(http.StatusInternalServerError, "Invalid settingsID", nil)
		return
	}

	settingsHelper := models.SettingsHelper()
	settings, err := settingsHelper.Load(settingsID)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to retrieve settings record", err)
		return
	}

	ctx.JSON(http.StatusOK, settings)
}

func UpdateSettings(ctx *flow.Context) {
	settingsHelper := models.SettingsHelper()
	settings, err := settingsHelper.LoadAndUpdateFromRequest(ctx.Req)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to load Settings record for update", err)
		return
	}

	// save and validate
	err = settingsHelper.Save(settings)
	// other type of error
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to save updated Settings record", err)
		return
	}

	ctx.JSON(http.StatusOK, settings)
}

func DeleteSettings(ctx *flow.Context) {
	settingsHelper := models.SettingsHelper()
	//get the settingsID from the request
	settingsID := ctx.URLIntParamWithDefault("settingsID", -1)
	if settingsID == -1 {
		ctx.JSON(http.StatusInternalServerError, "Invalid SettingsID for remove")
		return
	}

	isDeleted, err := settingsHelper.Delete(settingsID)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to remove the Settings record", err)
		return
	}
	ctx.JSON(http.StatusOK, isDeleted)
}
