package main

import (
	"html/template"
	"net/http"

	"github.com/nerdynz/builder/scaffold/server"
	"github.com/nerdynz/builder/scaffold/server/models"
	"github.com/bkono/datastore"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

func main() {
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	store := datastore.New()
	// defer store.Cleanup()
	renderer := render.New(render.Options{
		Layout:     "application",
		Extensions: []string{".html"},
		Funcs: []template.FuncMap{
			server.HelperFuncs,
		},
		// prevent having to rebuild for every template reload... This is an important setting for development speed
		IsDevelopment:               store.Settings.ServerIsDEV,
		RequirePartials:             store.Settings.ServerIsDEV,
		RequireBlocks:               store.Settings.ServerIsDEV,
		RenderPartialsWithoutPrefix: true,
	})
	store.Renderer = renderer

	models.Init(store.DB, store.Cache)
	attachments := negroni.NewStatic(http.Dir(store.Settings.AttachmentsFolder))
	attachments.Prefix = "/attachments"
	n.Use(attachments)
	adminNuxt := negroni.NewStatic(http.Dir("admin/dist/_nuxt"))
	adminNuxt.Prefix = "/admin/_nuxt"
	n.Use(adminNuxt)
	admin := negroni.NewStatic(http.Dir("./admin/dist/"))
	admin.Prefix = "/admin"
	n.Use(admin)
	public := negroni.NewStatic(http.Dir("./public/"))
	n.Use(public)
	//n.Use(cachecontrol.New(http.Dir("public")))
	n.UseHandler(server.Routes(store))
	n.Run(store.Settings.ServerPort)
}
