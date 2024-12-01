package shared

import (
	"log/slog"
	"net/http"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/shopstore"
)

type Layout interface {
	SetTitle(title string)
	SetScriptURLs(scripts []string)
	SetScripts(scripts []string)
	SetStyleURLs(styles []string)
	SetStyles(styles []string)
	SetBody(string)
	Render(w http.ResponseWriter, r *http.Request) string
}

type UIOptions struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Logger         *slog.Logger
	Store          shopstore.StoreInterface
	Layout         Layout
	HomeURL        string
	WebsiteUrl     string
}

type PageInterface interface {
	hb.TagInterface
	ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface
}
