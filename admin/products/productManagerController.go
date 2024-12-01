package admin

// import (
// 	"net/http"
// 	"project/config"
// 	"project/controllers/admin/shop/shared"

// 	"project/internal/helpers"
// 	"project/internal/layouts"
// 	"project/internal/links"
// 	"strings"

// 	"github.com/gouniverse/bs"
// 	"github.com/gouniverse/cdn"
// 	"github.com/gouniverse/form"
// 	"github.com/gouniverse/hb"
// 	"github.com/gouniverse/router"
// 	"github.com/gouniverse/sb"
// 	"github.com/gouniverse/shopstore"
// 	"github.com/gouniverse/utils"
// 	"github.com/samber/lo"
// 	"github.com/spf13/cast"
// )

// const ActionModalProductFilterShow = "modal_product_filter_show"

// // == CONTROLLER ==============================================================

// type productManagerController struct{}

// var _ router.HTMLControllerInterface = (*productManagerController)(nil)

// // == CONSTRUCTOR =============================================================

// func NewProductManagerController() *productManagerController {
// 	return &productManagerController{}
// }

// func (controller *productManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
// 	data, errorMessage := controller.prepareData(r)

// 	if errorMessage != "" {
// 		return helpers.ToFlashError(w, r, errorMessage, links.NewAdminLinks().Home(map[string]string{}), 10)
// 	}

// 	if data.action == ActionModalProductFilterShow {
// 		return controller.onModalProductFilterShow(data).ToHTML()
// 	}

// 	return layouts.NewAdminLayout(r, layouts.Options{
// 		Title:   "Products | Shop",
// 		Content: controller.page(data),
// 		ScriptURLs: []string{
// 			cdn.Htmx_1_9_4(),
// 			cdn.Sweetalert2_10(),
// 		},
// 		Styles: []string{},
// 	}).ToHTML()
// }

// func (controller *productManagerController) onModalProductFilterShow(data productManagerControllerData) *hb.Tag {
// 	modalCloseScript := `document.getElementById('ModalMessage').remove();document.getElementById('ModalBackdrop').remove();`

// 	title := hb.Heading5().
// 		Text("Filters").
// 		Style(`margin:0px;padding:0px;`)

// 	buttonModalClose := hb.Button().Type("button").
// 		Class("btn-close").
// 		Data("bs-dismiss", "modal").
// 		OnClick(modalCloseScript)

// 	buttonCancel := hb.Button().
// 		Child(hb.I().Class("bi bi-chevron-left me-2")).
// 		HTML("Cancel").
// 		Class("btn btn-secondary float-start").
// 		OnClick(modalCloseScript)

// 	buttonOk := hb.Button().
// 		Child(hb.I().Class("bi bi-check me-2")).
// 		HTML("Apply").
// 		Class("btn btn-primary float-end").
// 		OnClick(`FormFilters.submit();` + modalCloseScript)

// 	filterForm := form.NewForm(form.FormOptions{
// 		ID:     "FormFilters",
// 		Method: http.MethodGet,
// 		Fields: []form.FieldInterface{
// 			form.NewField(form.FieldOptions{
// 				Label: "Status",
// 				Name:  "filter_status",
// 				Type:  form.FORM_FIELD_TYPE_SELECT,
// 				Help:  `The status of the product.`,
// 				Value: data.formStatus,
// 				Options: []form.FieldOption{
// 					{
// 						Value: "",
// 						Key:   "",
// 					},
// 					{
// 						Value: "Active",
// 						Key:   shopstore.PRODUCT_STATUS_ACTIVE,
// 					},
// 					{
// 						Value: "Disabled",
// 						Key:   shopstore.PRODUCT_STATUS_DISABLED,
// 					},
// 					{
// 						Value: "Draft",
// 						Key:   shopstore.PRODUCT_STATUS_DRAFT,
// 					},
// 				},
// 			}),
// 			form.NewField(form.FieldOptions{
// 				Label: "Title",
// 				Name:  "filter_title",
// 				Type:  form.FORM_FIELD_TYPE_STRING,
// 				Value: data.formTitle,
// 				Help:  `Filter by title.`,
// 			}),
// 			form.NewField(form.FieldOptions{
// 				Label: "Created From",
// 				Name:  "filter_created_from",
// 				Type:  form.FORM_FIELD_TYPE_DATE,
// 				Value: data.formCreatedFrom,
// 				Help:  `Filter by creation date.`,
// 			}),
// 			form.NewField(form.FieldOptions{
// 				Label: "Created To",
// 				Name:  "filter_created_to",
// 				Type:  form.FORM_FIELD_TYPE_DATE,
// 				Value: data.formCreatedTo,
// 				Help:  `Filter by creation date.`,
// 			}),
// 			form.NewField(form.FieldOptions{
// 				Label: "Updated From",
// 				Name:  "filter_updated_from",
// 				Type:  form.FORM_FIELD_TYPE_DATE,
// 				Value: data.formUpdatedFrom,
// 				Help:  `Filter by update date.`,
// 			}),
// 			form.NewField(form.FieldOptions{
// 				Label: "Updated To",
// 				Name:  "filter_updated_to",
// 				Type:  form.FORM_FIELD_TYPE_DATE,
// 				Value: data.formUpdatedTo,
// 				Help:  `Filter by update date.`,
// 			}),
// 			form.NewField(form.FieldOptions{
// 				Label: "Product ID",
// 				Name:  "filter_product_id",
// 				Type:  form.FORM_FIELD_TYPE_STRING,
// 				Value: data.formProductID,
// 				Help:  `Find product by reference number (ID).`,
// 			}),
// 		},
// 	}).Build()

// 	modal := bs.Modal().
// 		ID("ModalMessage").
// 		Class("fade show").
// 		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
// 		Children([]hb.TagInterface{
// 			bs.ModalDialog().Children([]hb.TagInterface{
// 				bs.ModalContent().Children([]hb.TagInterface{
// 					bs.ModalHeader().Children([]hb.TagInterface{
// 						title,
// 						buttonModalClose,
// 					}),

// 					bs.ModalBody().
// 						Child(filterForm),

// 					bs.ModalFooter().
// 						Style(`display:flex;justify-content:space-between;`).
// 						Child(buttonCancel).
// 						Child(buttonOk),
// 				}),
// 			}),
// 		})

// 	backdrop := hb.Div().
// 		ID("ModalBackdrop").
// 		Class("modal-backdrop fade show").
// 		Style("display:block;")

// 	return hb.Wrap().Children([]hb.TagInterface{
// 		modal,
// 		backdrop,
// 	})

// }

// func (controller *productManagerController) page(data productManagerControllerData) hb.TagInterface {
// 	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
// 		{
// 			Name: "Home",
// 			URL:  links.NewAdminLinks().Home(map[string]string{}),
// 		},
// 		{
// 			Name: "Shop",
// 			URL:  shared.NewLinks().Home(map[string]string{}),
// 		},
// 		{
// 			Name: "Product Manager",
// 			URL:  shared.NewLinks().Products(map[string]string{}),
// 		},
// 	})

// 	buttonProductNew := hb.Button().
// 		Class("btn btn-primary float-end").
// 		Child(hb.I().Class("bi bi-plus-circle").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
// 		HTML("New Product").
// 		HxGet(shared.NewLinks().ProductCreate(map[string]string{})).
// 		HxTarget("body").
// 		HxSwap("beforeend")

// 	title := hb.Heading1().
// 		HTML("Shop. Product Manager").
// 		Child(buttonProductNew)

// 	return hb.Div().
// 		Class("container").
// 		Child(breadcrumbs).
// 		Child(hb.HR()).
// 		Child(title).
// 		Child(controller.tableProducts(data))
// }

// func (controller *productManagerController) tableProducts(data productManagerControllerData) hb.TagInterface {
// 	type columnStruct struct {
// 		Children    []columnStruct
// 		Label       string
// 		Sortable    bool
// 		SortableKey string
// 	}
// 	columns := []columnStruct{
// 		{
// 			Children: []columnStruct{
// 				{
// 					Label:       "Title",
// 					Sortable:    true,
// 					SortableKey: "title",
// 				},
// 				{
// 					Label:       "Reference",
// 					Sortable:    true,
// 					SortableKey: "id",
// 				},
// 			},
// 		},
// 		{
// 			Label:       "Status",
// 			Sortable:    true,
// 			SortableKey: "status",
// 		},
// 		{
// 			Label:       "Created",
// 			Sortable:    true,
// 			SortableKey: "created_at",
// 		},
// 		{
// 			Label:       "Modified",
// 			Sortable:    true,
// 			SortableKey: "updated_at",
// 		},
// 		{
// 			Label:    "Actions",
// 			Sortable: false,
// 		},
// 	}
// 	table := hb.Table().
// 		Class("table table-striped table-hover table-bordered").
// 		Children([]hb.TagInterface{
// 			hb.Thead().
// 				Child(hb.TR().Children(lo.Map(columns, func(column columnStruct, _ int) hb.TagInterface {
// 					children := []columnStruct{}
// 					if len(column.Children) == 0 {
// 						children = append(children, column)
// 					} else {
// 						children = append(children, column.Children...)
// 					}

// 					links := lo.Map(children, func(column columnStruct, _ int) string {
// 						return hb.Span().
// 							ChildIf(column.Sortable, controller.sortableColumnLabel(data, column.Label, column.SortableKey)).
// 							StyleIf(column.Sortable, "cursor: pointer;").
// 							TextIf(!column.Sortable, column.Label).
// 							ToHTML()
// 					})

// 					return hb.TH().
// 						HTML(strings.Join(links, ", ")).
// 						Style("width: 1px;")
// 				}))),
// 			hb.Tbody().Children(lo.Map(data.productList, func(product shopstore.ProductInterface, _ int) hb.TagInterface {
// 				productLink := hb.Hyperlink().
// 					Text(product.Title()).
// 					Href(links.NewAdminLinks().Tasks(map[string]string{"product_id": product.ID()}))

// 				status := hb.Span().
// 					Style(`font-weight: bold;`).
// 					StyleIf(product.Status() == shopstore.PRODUCT_STATUS_ACTIVE, `color:green;`).
// 					StyleIf(product.Status() == shopstore.PRODUCT_STATUS_DISABLED, `color:silver;`).
// 					StyleIf(product.Status() == shopstore.PRODUCT_STATUS_DRAFT, `color:blue;`).
// 					HTML(product.Status())

// 				buttonEdit := hb.Hyperlink().
// 					Class("btn btn-primary me-2").
// 					Child(hb.I().Class("bi bi-pencil-square")).
// 					Title("Edit").
// 					Href(shared.NewLinks().ProductUpdate(map[string]string{"product_id": product.ID()}))

// 				buttonDelete := hb.Hyperlink().
// 					Class("btn btn-danger").
// 					Child(hb.I().Class("bi bi-trash")).
// 					Title("Delete").
// 					HxGet(shared.NewLinks().ProductDelete(map[string]string{"product_id": product.ID()})).
// 					HxTarget("body").
// 					HxSwap("beforeend")

// 				return hb.TR().Children([]hb.TagInterface{
// 					hb.TD().
// 						Child(hb.Div().Child(productLink)).
// 						Child(hb.Div().
// 							Style("font-size: 11px;").
// 							HTML("Ref: ").
// 							HTML(product.ID())),
// 					hb.TD().
// 						Child(status),
// 					hb.TD().
// 						Child(hb.Div().
// 							Style("font-size: 13px;white-space: nowrap;").
// 							HTML(product.CreatedAtCarbon().Format("d M Y"))),
// 					hb.TD().
// 						Child(hb.Div().
// 							Style("font-size: 13px;white-space: nowrap;").
// 							HTML(product.UpdatedAtCarbon().Format("d M Y"))),
// 					hb.TD().
// 						Child(buttonEdit).
// 						Child(buttonDelete),
// 				})
// 			})),
// 		})

// 	// cfmt.Successln("Table: ", table)

// 	return hb.Wrap().Children([]hb.TagInterface{
// 		controller.tableFilter(data),
// 		table,
// 		controller.tablePagination(data, int(data.productCount), data.pageInt, data.perPage),
// 	})
// }

// func (controller *productManagerController) sortableColumnLabel(data productManagerControllerData, columnLabel string, columnSortKey string) hb.TagInterface {
// 	isSelected := strings.EqualFold(data.sortBy, columnSortKey)

// 	changeDirection := sb.ASC

// 	if isSelected {
// 		changeDirection = lo.If(data.sortOrder == sb.ASC, sb.DESC).Else(sb.ASC)
// 	}

// 	link := shared.NewLinks().Products(map[string]string{
// 		"page":       "0",
// 		"by":         columnSortKey,
// 		"sort":       changeDirection,
// 		"date_from":  data.formCreatedFrom,
// 		"date_to":    data.formCreatedTo,
// 		"status":     data.formStatus,
// 		"product_id": data.formProductID,
// 	})
// 	return hb.Hyperlink().
// 		HTML(columnLabel).
// 		Child(controller.sortingIndicator(columnSortKey, data.sortBy, changeDirection)).
// 		Href(link)
// }

// func (controller *productManagerController) sortingIndicator(columnSortKey string, sortByColumnKey string, sortOrder string) hb.TagInterface {
// 	isSelected := strings.EqualFold(sortByColumnKey, columnSortKey)

// 	direction := lo.If(isSelected && sortOrder == "asc", "up").
// 		ElseIf(isSelected && sortOrder == "desc", "down").
// 		Else("none")

// 	sortingIndicator := hb.Span().
// 		Class("sorting").
// 		HTMLIf(direction == "up", "&#8595;").
// 		HTMLIf(direction == "down", "&#8593;").
// 		HTMLIf(direction != "down" && direction != "up", "")

// 	return sortingIndicator
// }

// func (controller *productManagerController) tableFilter(data productManagerControllerData) hb.TagInterface {
// 	buttonFilter := hb.Button().
// 		Class("btn btn-sm btn-info me-2").
// 		Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
// 		Child(hb.I().Class("bi bi-filter me-2")).
// 		Text("Filters").
// 		HxPost(shared.NewLinks().Products(map[string]string{
// 			"action":              ActionModalProductFilterShow,
// 			"filter_title":        data.formTitle,
// 			"filter_status":       data.formStatus,
// 			"filter_product_id":   data.formProductID,
// 			"filter_created_from": data.formCreatedFrom,
// 			"filter_created_to":   data.formCreatedTo,
// 			"filter_updated_from": data.formUpdatedFrom,
// 			"filter_updated_to":   data.formUpdatedTo,
// 		})).
// 		HxTarget("body").
// 		HxSwap("beforeend")

// 	description := []string{
// 		hb.Span().HTML("Showing products").Text(" ").ToHTML(),
// 	}

// 	if data.formStatus != "" {
// 		description = append(description, hb.Span().Text("with status: "+data.formStatus).ToHTML())
// 	} else {
// 		description = append(description, hb.Span().Text("with status: any").ToHTML())
// 	}

// 	if data.formTitle != "" {
// 		description = append(description, hb.Span().Text("and email: "+data.formTitle).ToHTML())
// 	}

// 	if data.formProductID != "" {
// 		description = append(description, hb.Span().Text("and ID: "+data.formProductID).ToHTML())
// 	}

// 	if data.formCreatedFrom != "" && data.formCreatedTo != "" {
// 		description = append(description, hb.Span().Text("and created between: "+data.formCreatedFrom+" and "+data.formCreatedTo).ToHTML())
// 	} else if data.formCreatedFrom != "" {
// 		description = append(description, hb.Span().Text("and created after: "+data.formCreatedFrom).ToHTML())
// 	} else if data.formCreatedTo != "" {
// 		description = append(description, hb.Span().Text("and created before: "+data.formCreatedTo).ToHTML())
// 	}

// 	return hb.Div().
// 		Class("card bg-light mb-3").
// 		Style("").
// 		Children([]hb.TagInterface{
// 			hb.Div().Class("card-body").
// 				Child(buttonFilter).
// 				Child(hb.Span().
// 					HTML(strings.Join(description, " "))),
// 		})
// }

// func (controller *productManagerController) tablePagination(data productManagerControllerData, count int, page int, perPage int) hb.TagInterface {
// 	url := shared.NewLinks().Products(map[string]string{
// 		"filter_status":       data.formStatus,
// 		"filter_created_from": data.formCreatedFrom,
// 		"filter_created_to":   data.formCreatedTo,
// 		"filter_product_id":   data.formProductID,
// 		"filter_updated_from": data.formUpdatedFrom,
// 		"filter_updated_to":   data.formUpdatedTo,
// 		"by":                  data.sortBy,
// 		"order":               data.sortOrder,
// 	})

// 	url = lo.Ternary(strings.Contains(url, "?"), url+"&page=", url+"?page=") // page must be last

// 	pagination := bs.Pagination(bs.PaginationOptions{
// 		NumberItems:       count,
// 		CurrentPageNumber: page,
// 		PagesToShow:       5,
// 		PerPage:           perPage,
// 		URL:               url,
// 	})

// 	return hb.Div().
// 		Class(`d-flex justify-content-left mt-5 pagination-primary-soft rounded mb-0`).
// 		HTML(pagination)
// }

// func (controller *productManagerController) prepareData(r *http.Request) (data productManagerControllerData, errorMessage string) {
// 	var err error
// 	data.request = r
// 	data.action = utils.Req(r, "action", "")
// 	data.page = utils.Req(r, "page", "0")
// 	data.pageInt = cast.ToInt(data.page)
// 	data.perPage = cast.ToInt(utils.Req(r, "per_page", "10"))
// 	data.sortOrder = utils.Req(r, "sort", sb.DESC)
// 	data.sortBy = utils.Req(r, "by", shopstore.COLUMN_CREATED_AT)
// 	data.formProductID = utils.Req(r, "filter_product_id", "")
// 	data.formTitle = utils.Req(r, "filter_title", "")
// 	data.formStatus = utils.Req(r, "filter_status", "")
// 	data.formCreatedFrom = utils.Req(r, "filter_created_from", "")
// 	data.formCreatedTo = utils.Req(r, "filter_created_to", "")
// 	data.formUpdatedFrom = utils.Req(r, "filter_updated_from", "")
// 	data.formUpdatedTo = utils.Req(r, "filter_updated_to", "")

// 	productList, productCount, err := controller.fetchProductList(data)

// 	if err != nil {
// 		config.LogStore.ErrorWithContext("At productManagerController > prepareData", err.Error())
// 		return data, "error retrieving products"
// 	}

// 	data.productList = productList
// 	data.productCount = productCount

// 	return data, ""
// }

// func (controller *productManagerController) fetchProductList(data productManagerControllerData) ([]shopstore.ProductInterface, int64, error) {
// 	productIDs := []string{}

// 	if !lo.Contains([]string{sb.DESC, sb.ASC}, data.sortOrder) {
// 		data.sortOrder = sb.ASC
// 	}

// 	if !lo.Contains([]string{
// 		shopstore.COLUMN_CREATED_AT,
// 		shopstore.COLUMN_TITLE,
// 		shopstore.COLUMN_ID,
// 		shopstore.COLUMN_STATUS,
// 		shopstore.COLUMN_UPDATED_AT,
// 	}, data.sortBy) {
// 		data.sortBy = shopstore.COLUMN_CREATED_AT
// 	}

// 	query := shopstore.ProductQueryOptions{
// 		IDIn:      productIDs,
// 		Offset:    data.pageInt * data.perPage,
// 		Limit:     data.perPage,
// 		Status:    data.formStatus,
// 		SortOrder: data.sortOrder,
// 		OrderBy:   data.sortBy,
// 	}

// 	if data.formCreatedFrom != "" {
// 		query.CreatedAtGte = data.formCreatedFrom + " 00:00:00"
// 	}

// 	if data.formCreatedTo != "" {
// 		query.CreatedAtLte = data.formCreatedTo + " 23:59:59"
// 	}

// 	productList, err := config.ShopStore.ProductList(query)

// 	if err != nil {
// 		config.LogStore.ErrorWithContext("At productManagerController > prepareData", err.Error())
// 		return []shopstore.ProductInterface{}, 0, err
// 	}

// 	productCount, err := config.ShopStore.ProductCount(query)

// 	if err != nil {
// 		config.LogStore.ErrorWithContext("At productManagerController > prepareData", err.Error())
// 		return []shopstore.ProductInterface{}, 0, err
// 	}

// 	return productList, productCount, nil
// }

// type productManagerControllerData struct {
// 	request         *http.Request
// 	action          string
// 	page            string
// 	pageInt         int
// 	perPage         int
// 	sortOrder       string
// 	sortBy          string
// 	formStatus      string
// 	formTitle       string
// 	formCreatedFrom string
// 	formCreatedTo   string
// 	formUpdatedFrom string
// 	formUpdatedTo   string
// 	formProductID   string
// 	productList     []shopstore.ProductInterface
// 	productCount    int64
// }
