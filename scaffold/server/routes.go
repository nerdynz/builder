package server

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/nerdynz/datastore"
	"github.com/nerdynz/flow"
	"github.com/nerdynz/router"
	"github.com/nerdynz/security"
	"github.com/snabb/sitemap"
	"repo.nerdy.co.nz/thecollins/thecollins/server/actions"
	"repo.nerdy.co.nz/thecollins/thecollins/server/models"
)

var store *datastore.Datastore

func Routes(ds *datastore.Datastore) *bone.Mux {
	store = ds
	r := router.New(ds)
	// r.Mux.Handle("/admin/", http.FileServer(http.Dir("./admin/dist/")))
	if !store.Settings.ServerIsDEV {
		r.GET("/admin/:a", actions.SPA, security.NoAuth)
		r.GET("/admin/:a/:a", actions.SPA, security.NoAuth)
		r.GET("/admin/:a/:a/:a", actions.SPA, security.NoAuth)
		r.GET("/admin/:a/:a/:a/:a", actions.SPA, security.NoAuth)
		r.GET("/admin/:a/:a/:a/:a/:a", actions.SPA, security.NoAuth)
		r.GET("/admin/:a/:a/:a/:a/:a/:a", actions.SPA, security.NoAuth)
		r.GET("/admin/:a/:a/:a/:a/:a/:a/:a", actions.SPA, security.NoAuth)
		r.GET("/admin/:a/:a/:a/:a/:a/:a/:a/:a", actions.SPA, security.NoAuth)
		r.GET("/admin/:a/:a/:a/:a/:a/:a/:a/:a/:a", actions.SPA, security.NoAuth)
		r.GET("/admin/:a/:a/:a/:a/:a/:a/:a/:a/:a/:a", actions.SPA, security.NoAuth)
		r.GET("/admin/:a/:a/:a/:a/:a/:a/:a/:a/:a/:a/:a", actions.SPA, security.NoAuth)
		r.GET("/admin/:a/:a/:a/:a/:a/:a/:a/:a/:a/:a/:a/:a", actions.SPA, security.NoAuth)
		r.GET("/admin/:a/:a/:a/:a/:a/:a/:a/:a/:a/:a/:a/:a/:a", actions.SPA, security.NoAuth)
	}

	r.GET("/__ua", actions.Analytics, security.NoAuth)

	r.GET("/", actions.Home, security.NoAuth)
	r.GET("/home", actions.RedirectHome, security.NoAuth)
	r.GET("/contact", actions.ContactUs, security.NoAuth)

	// r.GET("/fix", actions.Fix, security.NoAuth)
	r.GET("/kitchen-sink", actions.KitchenSink, security.NoAuth)

	// r.GET("/register", actions.Register, security.NoAuth)
	// Scaffold routes
	r.GET("/api/v1/views/:month/:year", actions.Views, security.NoAuth)
	r.GET("/api/v1/sitesettings", siteSettings, security.NoAuth)
	r.GET("/api/v1/schema", Schema, security.NoAuth)
	r.PST("/api/v1/login", actions.Login, security.NoAuth)
	r.GET("/api/v1/user", actions.UserDetails, security.Disallow)

	r.GET("/api/v1/people/new", actions.NewPerson, security.NoAuth)
	r.PST("/api/v1/people/create", actions.CreatePerson, security.Disallow)
	r.GET("/api/v1/people/retrieve", actions.RetrievePeople, security.NoAuth)
	r.GET("/api/v1/people/retrieve/:personID", actions.RetrievePerson, security.NoAuth)
	r.PUT("/api/v1/people/update/:personID", actions.UpdatePerson, security.Disallow)

	// r.GET("/edit/:slug", actions.EditPage, security.NoAuth)

	r.GET("/api/v1/page/new", actions.NewPage, security.NoAuth)
	r.PST("/api/v1/page/create", actions.CreatePage, security.Disallow)
	r.GET("/api/v1/page/retrieve", actions.RetrievePages, security.NoAuth)
	r.GET("/api/v1/page/retrieve/:pageID", actions.RetrievePage, security.NoAuth)
	r.GET("/api/v1/page/retrieve/byslug/:slug", actions.RetrievePageBySlug, security.NoAuth)
	r.PUT("/api/v1/page/update/:pageID", actions.UpdatePage, security.Disallow)
	r.DEL("/api/v1/page/delete/:pageID", actions.DeletePage, security.Disallow)
	r.PUT("/api/v1/page/sort", actions.ChangePageSort, security.Disallow)

	r.GET("/:api/v1/work/new", actions.NewWork, security.Disallow)
	r.PST("/:api/v1/work/create", actions.CreateWork, security.Disallow)
	r.GET("/:api/v1/work/retrieve", actions.RetrieveWorks, security.Disallow)
	r.GET("/:api/v1/work/retrieve/:workID", actions.RetrieveWork, security.Disallow)
	r.GET("/:api/v1/work/paged/:sort/:direction/limit/:limit/pagenum/:pagenum", actions.PagedWorks, security.Disallow)
	r.PUT("/:api/v1/work/update/:workID", actions.UpdateWork, security.Disallow)
	r.DEL("/:api/v1/work/delete/:workID", actions.DeleteWork, security.Disallow)
	r.PUT("/:api/v1/work/sort", actions.SortWork, security.Disallow)

	r.GET("/:api/v1/person/new", actions.NewPerson, security.Disallow)
	r.PST("/:api/v1/person/create", actions.CreatePerson, security.Disallow)
	r.GET("/:api/v1/person/retrieve", actions.RetrievePeople, security.Disallow)
	r.GET("/:api/v1/person/retrieve/:personID", actions.RetrievePerson, security.Disallow)
	r.GET("/:api/v1/person/paged/:sort/:direction/limit/:limit/pagenum/:pagenum", actions.PagedPeople, security.Disallow)
	r.PUT("/:api/v1/person/update/:personID", actions.UpdatePerson, security.Disallow)
	r.DEL("/:api/v1/person/delete/:personID", actions.DeletePerson, security.Disallow)

	r.GET("/:api/v1/blog/new", actions.NewBlog, security.NoAuth)
	r.PST("/:api/v1/blog/create", actions.CreateBlog, security.NoAuth)
	r.GET("/:api/v1/blog/retrieve", actions.RetrieveBlogs, security.NoAuth)
	r.GET("/:api/v1/blog/retrieve/:blogID", actions.RetrieveBlog, security.NoAuth)
	r.GET("/:api/v1/blog/paged/:sort/:direction/limit/:limit/pagenum/:pagenum", actions.PagedBlogs, security.NoAuth)
	r.PUT("/:api/v1/blog/update/:blogID", actions.UpdateBlog, security.Disallow)
	r.DEL("/:api/v1/blog/delete/:blogID", actions.DeleteBlog, security.Disallow)

	r.POST("/api/v1/upload/crop", actions.CroppedFileUpload, security.Disallow)
	r.POST("/api/v1/upload/:quality/:type", actions.FileUpload, security.NoAuth)
	r.POST("/api/v1/upload/:type", actions.FileUpload, security.NoAuth)

	// r.GET("/:api/v1/imagemeta/new", actions.NewImageMeta, security.Disallow)
	// r.PST("/:api/v1/imagemeta/create", actions.CreateImageMeta, security.Disallow)
	r.GET("/:api/v1/imagemeta/retrieve/:uniqueid", actions.RetrieveImageMeta, security.Disallow)
	// r.GET("/:api/v1/imagemeta/paged/:sort/:direction/limit/:limit/pagenum/:pagenum", actions.PagedImageMeta, security.Disallow)
	// r.PUT("/:api/v1/imagemeta/update/:uniqueid", actions.UpdateImageMeta, security.Disallow)
	// r.DEL("/:api/v1/imagemeta/delete/:uniqueid", actions.DeleteImageMeta, security.Disallow)

	// r.GET("/testxx", actions.Test, security.NoAuth)

	r.GET("/sitemap.xml", websitemap, security.NoAuth)
	r.GET("/robots.txt", robots, security.NoAuth)

	// GOES LAST FOR GOOD MEASURE
	r.GET("/:slug", actions.RenderPageBySlug, security.NoAuth)
	return r.Mux
}

func Schema(ctx *flow.Context) {
	data := struct {
		Page             *models.Page
		Block            *models.Block
		Work             *models.Work
		Blog             *models.Blog
		IsSocketsEnabled bool
	}{
		Page:             models.PageHelper().New(),
		Block:            models.BlockHelper().New(),
		Work:             models.WorkHelper().New(),
		Blog:             models.BlogHelper().New(),
		IsSocketsEnabled: (ctx.Settings.Get("SocketsEnabled") == "true"),
	}
	ctx.JSON(http.StatusOK, data)
}

func robots(ctx *flow.Context) {
	robotsTxt := `User-agent: Teoma
Disallow: /
User-agent: twiceler
Disallow: /
User-agent: Gigabot
Disallow: /
User-agent: Scrubby
Disallow: /
User-agent: Nutch
Disallow: /
User-agent: baiduspider
Disallow: /
User-agent: naverbot
Disallow: /
User-agent: yeti
Disallow: /
User-agent: psbot
Disallow: /
User-agent: asterias
Disallow: /
User-agent: yahoo-blogs
Disallow: /
User-agent: YandexBot
Disallow: /
User-agent: Sosospider
Disallow: /
User-agent: *
Disallow: /admin
User-agent: *
Disallow: /df9249a6-0d56-11e8-ba89-0ed5f89f718b
User-agent: *
Disallow: /kitchen-sink
Sitemap: ` + ctx.Settings.WebsiteBaseURL + `sitemap.xml`
	ctx.Renderer.Text(ctx.W, 200, robotsTxt)
}

func websitemap(ctx *flow.Context) {
	sm := sitemap.New()
	pages, _ := models.PageHelper().All()

	for _, page := range pages {
		if page.ShowInNav == "Placeholder" {
			continue // skip placeholder
		}
		sm.Add(&sitemap.URL{
			Loc:        ctx.Settings.WebsiteBaseURL + page.Slug + "/",
			LastMod:    &page.DateModified,
			ChangeFreq: sitemap.Weekly,
		})
	}

	ctx.W.Header().Set("Content-Type", "text/xml")
	sm.WriteTo(ctx.W)
}

func siteSettings(ctx *flow.Context) {
	topNav, err := models.PageHelper().LoadTopNav()
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "Failed to load top nav", err)
		return
	}
	data := struct {
		TopNav models.NavItems
	}{
		topNav,
	}
	ctx.JSON(http.StatusOK, data)
}
