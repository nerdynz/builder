package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"repo.nerdy.co.nz/thecollins/thecollins/server/models"

	"github.com/nerdynz/helpers"

	"github.com/nerdynz/flow"

	dat "github.com/nerdynz/dat"
)

var HelperFuncs = template.FuncMap{
	"javascript":          javascriptTag,
	"stylesheet":          stylesheetTag,
	"javascriptAsync":     javascriptTagAsync,
	"stylesheetAsync":     stylesheetTagAsync,
	"image":               imageTag,
	"imagepath":           imagePath,
	"content":             content,
	"htmlblock":           htmlblock,
	"heroImagePageSlider": heroImagePageSlider,
	"slugify":             slugify,
	"JSON": func(v interface{}) template.JS {
		a, _ := json.Marshal(v)
		return template.JS(a)
	},
	"link":           link,
	"title":          title,
	"year":           year,
	"getPartialName": getPartialName,
	"hasValue":       hasValue,
	"isBlank":        isBlank,
	"isNotBlank":     isNotBlank,
	"formatDate":     formatDate,
	"htmlsafe":       htmlSafe,
	"gt":             greaterThan,
	"replace":        replace,
	"pictureBox":     pictureBox,
	"icon":           icon,
}

func pictureBox(str string, index int, w int, h int) template.HTML {
	if str == "" {
		return template.HTML("")
	}
	index = index - 1

	html := `	<div class="column is-4">
		<a href="javascript:void(0)" onclick="gallery.open(` + strconv.Itoa(index) + `);return false;" class="image" data-galleryitem='{"src": "` + imagePath(str) + `", "w": ` + strconv.Itoa(w) + `, "h": ` + strconv.Itoa(h) + `}'>
			<img src="` + imagePath(str) + `" />
		</a>
	</div>`
	return template.HTML(html)
}

func icon(str string, size string) template.HTML {
	html := `<span class="icon ` + size + `"><i class="fa fa-` + str + `"></i></span>`
	return template.HTML(html)
}

func replace(str string, old string, new string) string {
	return strings.Replace(str, old, new, -1)
}

func getPartialName(block *models.Block) string {
	return helpers.KebabCase(block.Type)
}

func year() string {
	return strconv.Itoa(time.Now().Year())
}

func greaterThan(num int, amt int) bool {
	return num > amt
}

// func content(contents ...string) template.HTML {
// 	var str string
// 	for _, content := range contents {
// 		str += "<div class='standard'>" + content + "</standard>"
// 	}
// 	return template.HTML(str)
// }

func javascriptTag(names ...string) template.HTML {
	var str string
	for _, name := range names {
		if strings.HasPrefix(name, "http") {
			str += "<script src='" + name + ".js' type='text/javascript'></script>"
		} else {
			str += "<script src='/js/" + name + ".js' type='text/javascript'></script>"
		}
	}
	return template.HTML(str)
}

func javascriptTagAsync(names ...string) template.HTML {
	var str string
	for _, name := range names {
		href := ""
		if strings.HasPrefix(name, "http") {
			href = name + ".js"
		} else {
			href = "/js/" + name + ".js"
		}
		str += `<script type="text/javascript">(function () { head.load('` + href + `');})() </script>`
	}
	return template.HTML(str)
}

func stylesheetTagAsync(names ...string) template.HTML {
	var str string
	for _, name := range names {
		href := "/css/" + name + ".css"
		//str += `<script type="text/javascript">(function () { var rl = document.createElement('link'); rl.rel = 'stylesheet';rl.href = '` + href + `';var rh = document.getElementsByTagName('head')[0]; rh.parentNode.insertBefore(rl, rh);})();</script>`
		str += `<script type="text/javascript">(function () { head.load('` + href + `');})() </script>`
	}
	return template.HTML(str)
}

func stylesheetTag(names ...string) template.HTML {
	var str string
	for _, name := range names {
		str += "<link rel='stylesheet' href='/css/" + name + ".css' type='text/css' media='screen'  />\n"
	}
	return template.HTML(str)
}

func imagePath(n interface{}) string {
	name := n.(string)
	// if store.Settings.ServerIsDEV && strings.HasPrefix(name, "/attachments") {
	// 	return strings.Replace(name, "/attachments/", "https://nerdy.co.nz/attachments/", 1)
	// 	// return "https://cdn.nerdy.co.nz/attachments/" + store.Settings.Sitename + "/" + name
	// }
	if strings.HasPrefix(name, "/attachments") {
		return name
	}
	if strings.HasPrefix(name, "/images") {
		return name
	}
	if strings.HasPrefix(name, "/") {
		return name
	}
	return "/images/" + name
}

func imageTag(name interface{}, alt interface{}, class string) template.HTML {
	if name.(string) != "" {
		return template.HTML("<img src='" + imagePath(name) + "' alt='" + alt.(string) + "' class='" + class + "' />")
	}
	return template.HTML("")
}

func htmlSafe(str string) template.HTML {
	return template.HTML(str)
}

func content(str string) template.HTML {
	return template.HTML("<div class='content'>" + str + "</div>")
}

func htmlblock(page *Page, code string) template.HTML {
	html := "<div class='textblock editable' "
	html += " data-textblock='page-" + strconv.FormatInt(page.PageID, 10) + "-" + code + "'"
	html += " data-placeholder='#{placeholder}'> "
	html += getHTMLFromTextblock(page, code)
	html += "</div>"
	return template.HTML(html)
}

func link(bucket flow.Bucket, text string, link string, cssClass string) template.HTML {
	currentLink := strings.ToLower(bucket["currentURL"].(string))
	class := "link link-" + strings.ToLower(link)
	disabled := ""
	if strings.HasSuffix(link, "home") {
		link = ""
	}
	if "/"+strings.ToLower(link) == currentLink {
		class += " active"
		disabled = "disabled"
	}
	class += " " + cssClass
	fullLink := bucket["websiteBaseURL"].(string) + link
	return template.HTML(fmt.Sprintf(`<a class="%v" href="%v" `+disabled+`><span class="inner-link">%v</span></a>`, class, fullLink, text))
}

func title(text string) string {
	return strings.Title(text)
}

func heroImagePageSlider(page *models.Page) template.HTML {
	html := ""
	if page.Picture != "" {
		html += "<img src='" + page.Picture + "' />"
	}
	if page.PictureTwo != "" {
		html += "<img src='" + page.PictureTwo + "' />"
	}
	if page.PictureThree != "" {
		html += "<img src='" + page.PictureThree + "' />"
	}
	if page.PictureFour != "" {
		html += "<img src='" + page.PictureFour + "' />"
	}
	return template.HTML(html)
}

// func navigation(viewBag *ViewBucket) template.HTML {
// 	html := ""
// 	if viewBag.Data["NavItems"] != nil {
// 		navItems := viewBag.Data["NavItems"].([]*NavItem)
// 		html = "<nav class='main-nav closed'>"
// 		for _, navItem := range navItems {
// 			html += "<a href='/" + navItem.Slug + "'>" + navItem.Title + "</a>"
// 		}
// 		html += "</nav>"
// 	}
// 	return template.HTML(html)
// }
func hasValue(val interface{}) bool {
	return val != nil && val != ""
}

func isBlank(str string) bool {
	return str == ""
}

func isNotBlank(str string) bool {
	return !isBlank(str)
}

type Page struct {
	PageID     int64        `db:"page_id"`
	Title      string       `db:"title"`
	Body       string       `db:"body"`
	Slug       string       `db:"slug"`
	Template   string       `db:"template"`
	CreatedAt  dat.NullTime `db:"created_at"`
	UpdatedAt  dat.NullTime `db:"updated_at"`
	Textblocks []*Textblock
}

type NavItem struct {
	Title string `db:"title"`
	Slug  string `db:"slug"`
}

func (navItem *NavItem) getURL() string {
	return ""
}

type Textblock struct {
	TextblockID int64        `db:"textblock_id"`
	Code        string       `db:"code"`
	Body        string       `db:"body"`
	CreatedAt   dat.NullTime `db:"created_at"`
	UpdatedAt   dat.NullTime `db:"updated_at"`
	PageID      int64        `db:"page_id"`
}

func getHTMLFromTextblock(page *Page, code string) string {
	var body string
	for _, tb := range page.Textblocks {
		if tb.Code == code {
			body = tb.Body
		}
	}
	return body
}

func formatDate(time time.Time, layout string) string {
	return time.Format(layout)
}

func slugify(str string) string {
	return slug.Make(str)
}
