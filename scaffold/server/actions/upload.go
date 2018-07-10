package actions

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"strconv"

	flow "github.com/nerdynz/flow"
	"repo.nerdy.co.nz/thecollins/thecollins/server/models"
)

func FileUpload(ctx *flow.Context) {
	// 	fileType := ctx.URLParam("type")
	// 	if fileType == "" {
	// 		ctx.ErrorJSON(http.StatusBadRequest, "Unrecognised upload type", nil)
	// 		return
	// 	}

	// 	data := &fileUploadStruct{}
	// 	req := ctx.Req
	// 	req.ParseMultipartForm(32 << 20)

	// 	// get the request
	// 	var reqFile multipart.File
	// 	var handler *multipart.FileHeader
	// 	var err error

	// 	reqFile, handler, err = req.FormFile("file")
	// 	if err != nil {
	// 		ctx.JSON(http.StatusInternalServerError, err.Error())
	// 		return
	// 	}
	// 	defer reqFile.Close()

	// 	// proxy to image processsing server
	// 	filename := getValidFileName(ctx.Store.Settings.AttachmentsFolder, handler.Filename)
	// 	// lFile := strings.ToLower(filename)
	// 	f, err := os.OpenFile(ctx.Store.Settings.AttachmentsFolder+filename, os.O_RDWR|os.O_CREATE, 0666)
	// 	if err != nil {
	// 		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to create a file on the filesystem", err)
	// 		return
	// 	}
	// 	defer f.Close()

	// 	err = processedImage(f, reqFile, "jpg", 85)

	// 	data.URL = "/attachments/" + filename
	// 	data.FileName = filename
	// ctx.JSON(http.StatusOK, data)
	ctx.JSON(http.StatusOK, nil)
}

func CroppedFileUpload(ctx *flow.Context) {
	helper := models.ImageMetaHelper()
	meta, err := helper.NewFromRequest(ctx.Req)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to create image", err)
		return
	}

	// original
	if meta.IsExisting { // dont keep saving new copies of an original image, if they are just editing
		// clear the old image etc
	} else {
		f, err := os.OpenFile(ctx.Store.Settings.AttachmentsFolder+meta.Original, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			ctx.ErrorJSON(http.StatusInternalServerError, "Failed to create a file on the filesystem", err)
			return
		}
		defer f.Close()
		b, err := meta.OriginalBytes()
		if err != nil {
			ctx.ErrorJSON(http.StatusInternalServerError, "failed to get bytes from the original image", err)
			return
		}
		_, err = io.Copy(f, bytes.NewReader(b)) // no processing
		if err != nil {
			ctx.ErrorJSON(http.StatusInternalServerError, "failed to save the original image", err)
			return
		}
	}

	// new
	b, err := meta.Bytes()
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "failed to get bytes from the new image", err)
		return
	}

	f, err := os.OpenFile(ctx.Store.Settings.AttachmentsFolder+meta.Name, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to create a file on the filesystem", err)
		return
	}
	defer f.Close()
	err = processedImage(f, bytes.NewReader(b), meta.Ext, meta.Width, meta.Height, 85, meta.IsConvert())
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "failed to process image", err)
		return
	}
	helper.Save(meta)
	ctx.JSON(200, "image saved")
}

func getValidFileName(path string, filename string) string {
	return getValidFileNameWithDupIndex(path, filename, 0)
}

func getValidFileNameWithDupIndex(path string, filename string, duplicateIndex int) string {
	// filename = sanitize.Path(filename)
	dupStr := ""
	if duplicateIndex > 0 {
		dupStr = "" + strconv.Itoa(duplicateIndex) + "-"
	}
	fullpath := path + dupStr + filename

	// path doesn't exist so we can return this path
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		return dupStr + filename
	}

	//otherwise increase file index and
	duplicateIndex++
	return getValidFileNameWithDupIndex(path, filename, duplicateIndex)
}

type fileUploadStruct struct {
	FileName string
	URL      string `json:"link"`
}

type operation struct {
	Operation string                 `json:"operation"`
	Params    map[string]interface{} `json:"params"`
}

type operations struct {
	Ops []*operation
}

func (o *operation) addParam(key string, val interface{}) {
	if o.Params == nil {
		o.Params = make(map[string]interface{})
	}
	o.Params[key] = val
}

func (o *operations) add(op string) {
	if o.Ops == nil {
		o.Ops = make([]*operation, 0)
	}
	o.Ops = append(o.Ops, &operation{
		Operation: op,
	})
}

func (o *operations) last() *operation {
	return o.Ops[len(o.Ops)-1]
}

func processedImage(f *os.File, r io.Reader, imageType string, width int, height int, quality int, convert bool) error {
	ops := &operations{}

	originalImageType := "jpg"
	if convert {
		ops.add("convert")

		// converting
		if imageType == "jpg" {
			imageType = "jpeg"
			originalImageType = "png"
		} else if imageType == "png" {
			originalImageType = "jpg"
		}
		ops.last().addParam("type", imageType)
	}

	ops.add("fit")
	ops.last().addParam("width", width)    //absolute max
	ops.last().addParam("height", height)  // dont need its ratio based
	ops.last().addParam("stripmeta", true) // dont need its ratio based
	ops.last().addParam("quality", quality)
	// ops.last().addParam("compression", quality)
	bOps, err := json.Marshal(ops.Ops)
	if err != nil {
		return err
	}
	endpoint := "https://images.nerdy.co.nz/pipeline?operations=" + url.QueryEscape(string(bOps))
	// endpoint = "https://images.nerdy.co.nz/fit?width=200&height=200"

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("file", "filename_placeholder."+originalImageType)
	if err != nil {
		return err
		// ctx.ErrorJSON(http.StatusOK, "couldn't create form file ", err)
	}
	_, err = io.Copy(fw, r)
	if err != nil {
		// ctx.ErrorJSON(http.StatusOK, "failed to copy from reqFile", err)
		return err
	}
	err = w.Close()
	if err != nil {
		// ctx.ErrorJSON(http.StatusOK, "failed to copy from reqFile", err)
		return err
	}

	req, err := http.NewRequest("POST", endpoint, &b)
	if err != nil {
		// ctx.ErrorJSON(http.StatusOK, "failed to copy from reqFile", err)
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil || res.StatusCode != 200 {
		// ctx.ErrorJSON(http.StatusInternalServerError, "bad request", err)
		return err
	}
	defer res.Body.Close()

	// we read from tee reader as it hasn't already done its scan
	_, err = io.Copy(f, res.Body)
	if err != nil {
		// ctx.ErrorJSON(http.StatusInternalServerError, "Failed to create image", err)
		return err
	}
	return nil
}
