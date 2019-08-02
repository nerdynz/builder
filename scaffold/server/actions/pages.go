package actions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	flow "github.com/nerdynz/flow"
	"github.com/nerdynz/builder/scaffold/server/models"
)

func loadPageExtras(ctx *flow.Context) {
	navItems, err := models.PageHelper().LoadTopNav()
	if err != nil {
		ctx.ErrorHTML(http.StatusNotFound, "Something went wrong", err)
		return
	}
	ctx.Add("TopNav", navItems)
	navItems, err = models.PageHelper().LoadSideNav()
	if err != nil {
		ctx.ErrorHTML(http.StatusNotFound, "Something went wrong", err)
		return
	}
	ctx.Add("SideNav", navItems)
	navItems, err = models.PageHelper().LoadFooterNav()
	if err != nil {
		ctx.ErrorHTML(http.StatusNotFound, "Something went wrong", err)
		return
	}
	ctx.Add("FooterNav", navItems)
	// if ctx.Padlock.IsLoggedIn() {
	// 	user, _ := ctx.Padlock.LoggedInUser()
	// 	member, _ := models.MemberHelper().Load(user.ID)
	// 	if member != nil {
	// 		ctx.Add("Member", member)
	// 	}
	// }
	settings, err := models.SettingHelper().Load(1)
	if err != nil {
		ctx.ErrorHTML(http.StatusNotFound, "Something went wrong", err)
		return
	}
	ctx.Add("Settings", settings)
}

func RedirectHome(ctx *flow.Context) {
	ctx.Redirect("/", http.StatusMovedPermanently)
}

func Analytics(ctx *flow.Context) {
	helper := models.AnalyticsHelper()
	a := helper.New()

	if ctx.URLParam("t") == "pageview" {
		// screenres
		sr := ctx.URLParam("sr")
		if sr != "" {
			res := strings.Split(sr, "x")
			if w, err := strconv.Atoi(res[0]); err == nil {
				a.ScreenWidth = w
			}
			if h, err := strconv.Atoi(res[1]); err == nil {
				a.ScreenHeight = h
			}
		}

		// viewport
		vp := ctx.URLParam("vp")
		if vp != "" {
			res := strings.Split(vp, "x")
			if w, err := strconv.Atoi(res[0]); err == nil {
				a.VpWidth = w
			}
			if h, err := strconv.Atoi(res[1]); err == nil {
				a.VpHeight = h
			}
		}

		// page
		a.Page = strings.ToLower(ctx.URLParam("dl"))
		ua := ctx.Req.Header.Get("User-Agent")
		a.UserAgent = ua
		// unique id0
		a.SetUaInfo()
		a.SetUniqueID(ctx.URLParam("cid"))

		err := helper.Save(a) // ignore errors
		if err != nil {
			ctx.ErrorHTML(500, "failed", err)
			return
		}
	}
	ctx.Text(200, "done")
}

func Fix(ctx *flow.Context) {
	as, _ := models.AnalyticsHelper().All()
	for _, a := range as {
		a.SetUaInfo()
	}
	models.AnalyticsHelper().SaveMany(as)
}

func Home(ctx *flow.Context) {
	loadPageExtras(ctx)
	page, err := models.PageHelper().LoadBySpecialPage("home")
	renderPage(ctx, page, err)
}

func KitchenSink(ctx *flow.Context) {
	page, err := models.PageHelper().LoadKitchenSink()
	renderPage(ctx, page, err)
}

func RenderPageBySlug(ctx *flow.Context) {
	loadPageExtras(ctx)
	pageSlug := ctx.URLParam("slug")
	if pageSlug == "favicon.ico" {
		return
	}
	page, err := models.PageHelper().LoadBySlug(pageSlug)

	renderPage(ctx, page, err)
}

func ContactUs(ctx *flow.Context) {
	loadPageExtras(ctx)
	form, err := models.FormHelper().LoadFullForm("Contact")
	if err != nil {
		ctx.ErrorHTML(http.StatusInternalServerError, "Failed to Load contact Form Details", err)
		return
	}
	ctx.Add("Form", form)
	// ctx.JSON(200, ctx.Bucket)
	// return
	page, err := models.PageHelper().LoadBySpecialPage("contact")
	renderPage(ctx, page, err)
}

func renderPage(ctx *flow.Context, page *models.Page, err error) {
	if err != nil {
		ctx.ErrorHTML(http.StatusNotFound, "We couldn't find the page you were looking for.", err)
		return
	}
	ctx.Add("Page", page)

	kind := page.Kind
	specialFor := page.SpecialPageFor
	if specialFor != "" {
		if strings.Contains(specialFor, ":") {
			s := strings.Split(specialFor, ":")
			kind = s[0]
			specialFor = s[1]
		} else {
			ctx.HTML(page.SpecialPageFor, http.StatusOK)
			return
		}
	}

	// // check for testimonials
	// hasTestimonial := false
	// for _, block := range page.Blocks {
	// 	if block.Type == "Testimonials" {
	// 		hasTestimonial = true
	// 		break
	// 	}
	// }
	// if hasTestimonial {
	// 	testimonials, err := models.TestimonialHelper().Random(3)
	// 	if err != nil {
	// 		ctx.ErrorHTML(http.StatusInternalServerError, "Testimonials Can't be loaded", err)
	// 		return
	// 	}
	// 	ctx.Add("Testimonials", testimonials)
	// }

	if kind == "Hero Image" && page.HasPictures() {
		ctx.HTML("hero-image", http.StatusOK)
		return
	}
	ctx.HTML("default", http.StatusOK)
}

// func EditPage(ctx *flow.Context) {
// 	pageID, err := ctx.URLIntParam("pageID")
// 	if err != nil || pageID == 0 {
// 		ctx.ErrorJSON(http.StatusInternalServerError, "invalid pageID")
// 		return
// 	}

// 	pageHelper := models.PageHelper()
// 	page, err := pageHelper.Load(pageID)
// 	if err != nil {
// 		ctx.ErrorJSON(http.StatusInternalServerError, "", err)
// 		return
// 	}
// 	page.IsBeingEdited = true

// 	content, err := template.Content(ctx, page)
// 	if err != nil {
// 		ctx.ErrorJSON(http.StatusInternalServerError, "", err)
// 		return
// 	}
// 	ctx.RenderByPageKind(content)
// }

// RESTFUL METHODS
func NewPage(ctx *flow.Context) {
	pageHelper := models.PageHelper()
	page := pageHelper.New()
	ctx.JSON(http.StatusOK, page)
}

func CreatePage(ctx *flow.Context) {
	pageHelper := models.PageHelper()
	page, err := pageHelper.NewFromRequest(ctx.Req)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to load page changes", err)
		return
	}
	err = pageHelper.Save(page)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to save page", err)
		return
	}
	ctx.JSON(http.StatusOK, page)
}

func RetrievePage(ctx *flow.Context) {
	if ctx.URLParam("pageID") == "" {
		RetrievePages(ctx)
		return
	}

	//get the pageID from the request
	pageID, err := ctx.URLIntParam("pageID")
	if err != nil || pageID == 0 {
		ctx.ErrorJSON(http.StatusInternalServerError, "Invalid Page ID", nil)
		return
	}

	pageHelper := models.PageHelper()
	page, err := pageHelper.Load(pageID)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "", err)
		return
	}

	err = page.LoadBlocks()
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "", err)
		return
	}

	ctx.JSON(http.StatusOK, page)
}

func RetrievePageBySlug(ctx *flow.Context) {
	slug := ctx.URLParam("slug")
	if slug == "" {
		ctx.ErrorJSON(http.StatusInternalServerError, "Invalid Page ID", nil)
		return
	}

	pageHelper := models.PageHelper()
	page, err := pageHelper.LoadBySlug(slug)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "", err)
		return
	}

	err = page.LoadBlocks()
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "", err)
		return
	}

	ctx.JSON(http.StatusOK, page)
}

func RetrievePages(ctx *flow.Context) {
	pageHelper := models.PageHelper()
	pages, err := pageHelper.All()

	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "", err)
		return
	}

	ctx.JSON(http.StatusOK, pages)
}

func UpdatePage(ctx *flow.Context) {
	pageHelper := models.PageHelper()
	page, err := pageHelper.LoadAndUpdateFromRequest(ctx.Req)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "", err)
		return
	}

	// save and validate
	err = pageHelper.Save(page)
	// other type of error
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "", err)
		return
	}

	ctx.JSON(http.StatusOK, page)
}

// type pageSort struct {
// 	PageID       int `json:"PageID"`
// 	SortPosition int `json:"SortPosition"`
// }

func ChangePageSort(ctx *flow.Context) {
	var sort models.Pages
	decoder := json.NewDecoder(ctx.Req.Body)
	err := decoder.Decode(&sort)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "", err)
		return
	}

	err = models.PageHelper().SaveMany(sort)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to update sort Position", err)
		return
	}

	ctx.JSON(http.StatusOK, sort)
}

func DeletePage(ctx *flow.Context) {
	pageHelper := models.PageHelper()
	//get the pageID from the request
	pageID, err := ctx.URLIntParam("pageID")
	if err != nil || pageID == 0 {
		ctx.ErrorJSON(http.StatusInternalServerError, "Invalid Page ID", nil)
		return
	}

	isDeleted, err := pageHelper.Delete(pageID)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "", err)
		return
	}
	ctx.JSON(http.StatusOK, isDeleted)
}

func TestMessage(ctx *flow.Context) {

}

func SPA(ctx *flow.Context) {
	ctx.W.Header().Add("content-type", "text/html")
	file, err := ioutil.ReadFile("admin/dist/index.html")
	if err != nil {
		ctx.ErrorHTML(500, "Failed to load SPA", err)
		return
	}
	ctx.Renderer.Data(ctx.W, 200, file)
}

func FirebaseMessaging(ctx *flow.Context) {
	ctx.W.Header().Add("content-type", "application/javascript")
	file, err := ioutil.ReadFile("public/js/firebase-messaging-sw.js")
	if err != nil {
		ctx.ErrorHTML(500, "Failed to load SPA", err)
		return
	}
	ctx.Renderer.Data(ctx.W, 200, file)
}
