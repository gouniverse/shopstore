package admin

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/shopstore"
	"github.com/gouniverse/shopstore/admin/shared"
	"github.com/gouniverse/utils"
)

func UI(options shared.UIOptions) (hb.TagInterface, error) {
	if options.ResponseWriter == nil {
		return nil, errors.New("options.ResponseWriter is required")
	}

	if options.Request == nil {
		return nil, errors.New("options.Request is required")
	}

	if options.Store == nil {
		return nil, errors.New("options.Store is required")
	}

	if options.Logger == nil {
		return nil, errors.New("options.Logger is required")
	}

	if options.Layout == nil {
		return nil, errors.New("options.Layout is required")
	}

	ui := &ui{
		response:   options.ResponseWriter,
		request:    options.Request,
		store:      options.Store,
		logger:     *options.Logger,
		layout:     options.Layout,
		homeURL:    options.HomeURL,
		websiteUrl: options.WebsiteUrl,
	}

	return ui.handler(), nil
}

type ui struct {
	response   http.ResponseWriter
	request    *http.Request
	store      shopstore.StoreInterface
	logger     slog.Logger
	layout     shared.Layout
	homeURL    string
	websiteUrl string
}

func (ui *ui) handler() hb.TagInterface {
	controller := utils.Req(ui.request, "controller", "")

	if controller == "" {
		controller = pathHome
	}

	if controller == pathHome {
		return home(*ui)
	}

	if controller == pathDiscounts {
		// return visitorActivity(*ui)
	}

	if controller == pathProducts {
		// return visitorPaths(*ui)
	}

	ui.layout.SetBody(hb.H1().HTML(controller).ToHTML())
	return hb.Raw(ui.layout.Render(ui.response, ui.request))
	// redirect(a.response, a.request, url(a.request, pathQueueManager, map[string]string{}))
	// return nil
}
