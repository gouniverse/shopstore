package admin

import (
	"net/http"

	"github.com/gouniverse/cdn"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/shopstore/admin/shared"
)

// == CONSTRUCTOR ==============================================================

func home(ui ui) shared.PageInterface {
	return &homeController{
		ui: ui,
	}
}

// == CONTROLLER ===============================================================

type homeController struct {
	ui ui
}

type homeControllerData struct{}

func (c *homeController) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, errorMessage := c.prepareData()

	c.ui.layout.SetTitle("Dashboard | Kalleidoscope")

	if errorMessage != "" {
		c.ui.layout.SetBody(hb.Div().
			Class("alert alert-danger").
			Text(errorMessage).ToHTML())

		return hb.Raw(c.ui.layout.Render(w, r))
	}

	htmxScript := `setTimeout(() => async function() {
		if (!window.htmx) {
			let script = document.createElement('script');
			document.head.appendChild(script);
			script.type = 'text/javascript';
			script.src = '` + cdn.Htmx_2_0_0() + `';
			await script.onload
		}
	}, 1000);`

	swalScript := `setTimeout(() => async function() {
		if (!window.Swal) {
			let script = document.createElement('script');
			document.head.appendChild(script);
			script.type = 'text/javascript';
			script.src = '` + cdn.Sweetalert2_11() + `';
			await script.onload
		}
	}, 1000);`

	// cdn.Jquery_3_7_1(),
	// // `https://cdnjs.cloudflare.com/ajax/libs/Chart.js/1.0.2/Chart.min.js`,
	// `https://cdn.jsdelivr.net/npm/chart.js`,

	c.ui.layout.SetBody(c.page(data).ToHTML())
	c.ui.layout.SetScripts([]string{htmxScript, swalScript})

	return hb.Raw(c.ui.layout.Render(w, r))
}

func (c *homeController) ToHTML() string {
	return c.ToTag(nil, nil).ToHTML()
}

// == PRIVATE METHODS ==========================================================

func (c *homeController) prepareData() (data homeControllerData, errorMessage string) {
	return homeControllerData{}, ""
}

func (c *homeController) page(data homeControllerData) hb.TagInterface {
	breadcrumbs := breadcrumbs(c.ui.request, []Breadcrumb{
		{
			Name: "Home",
			URL:  url(c.ui.request, c.ui.homeURL, nil),
		},
		{
			Name: "Kalleidoscope",
			URL:  url(c.ui.request, pathHome, nil),
		},
	})

	title := hb.Heading1().
		HTML("Bazaar. Home")

	options :=
		hb.Section().
			Class("mb-3 mt-3").
			Style("background-color: #f8f9fa;").
			Child(
				hb.UL().
					Class("list-group").
					Child(hb.LI().
						Class("list-group-item").
						Child(hb.A().
							Href(url(c.ui.request, pathProducts, nil)).
							Text("Products"))).
					Child(hb.LI().
						Class("list-group-item").
						Child(hb.A().
							Href(url(c.ui.request, pathDiscounts, nil)).
							Text("Discounts"))))

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(adminHeader(c.ui)).
		Child(hb.HR()).
		Child(title).
		Child(options)
}

// == PRIVATE METHODS ==========================================================
