package actions

import (
	"net/http"

	flow "github.com/nerdynz/flow"
	"repo.nerdy.co.nz/displayworks/displayworks-signs/server/models"
)

func NewImageMeta(ctx *flow.Context) {
	imageMetaHelper := models.ImageMetaHelper()
	imageMeta := imageMetaHelper.New()
	ctx.JSON(http.StatusOK, imageMeta)
}

func CreateImageMeta(ctx *flow.Context) {
	imageMetaHelper := models.ImageMetaHelper()
	imageMeta, err := imageMetaHelper.NewFromRequest(ctx.Req)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to create imageMeta record", err)
		return
	}
	err = imageMetaHelper.Save(imageMeta)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to create imageMeta record", err)
		return
	}
	ctx.JSON(http.StatusOK, imageMeta)
}

func RetrieveImageMeta(ctx *flow.Context) {
	uniqueID := ctx.URLUnique()
	if uniqueID == "" {
		ctx.ErrorJSON(http.StatusInternalServerError, "Invalid uniqueID", nil)
		return
	}

	imageMetaHelper := models.ImageMetaHelper()
	imageMeta, err := imageMetaHelper.One("unique_id = $1", uniqueID)
	if err != nil {
		ctx.JSON(500, nil)
		return
	}

	ctx.JSON(http.StatusOK, imageMeta)
}

func PagedImageMeta(ctx *flow.Context) {
	imageMetaHelper := models.ImageMetaHelper()
	pageNum := ctx.URLIntParamWithDefault("pagenum", 1)
	limit := ctx.URLIntParamWithDefault("limit", 10)
	sort := ctx.URLParam("sort")
	direction := ctx.URLParam("direction")

	data, err := imageMetaHelper.PagedBy(pageNum, limit, sort, direction)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Unabled to get paged ImageMeta data", err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func UpdateImageMeta(ctx *flow.Context) {
	imageMetaHelper := models.ImageMetaHelper()
	imageMeta, err := imageMetaHelper.LoadAndUpdateFromRequest(ctx.Req)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to load ImageMeta record for update", err)
		return
	}

	// save and validate
	err = imageMetaHelper.Save(imageMeta)
	// other type of error
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to save updated ImageMeta record", err)
		return
	}

	ctx.JSON(http.StatusOK, imageMeta)
}

func DeleteImageMeta(ctx *flow.Context) {
	imageMetaHelper := models.ImageMetaHelper()
	//get the imageMetaID from the request
	imageMetaID := ctx.URLIntParamWithDefault("imageMetaID", -1)
	if imageMetaID == -1 {
		ctx.JSON(http.StatusInternalServerError, "Invalid ImageMetaID for remove")
		return
	}

	isDeleted, err := imageMetaHelper.Delete(imageMetaID)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to remove the ImageMeta record", err)
		return
	}
	ctx.JSON(http.StatusOK, isDeleted)
}
