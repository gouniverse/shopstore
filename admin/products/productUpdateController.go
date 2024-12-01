package admin

// import (
// 	"net/http"
// 	"project/config"
// 	"project/controllers/admin/shop/shared"
// 	"slices"
// 	"strconv"
// 	"strings"

// 	"project/internal/helpers"
// 	"project/internal/layouts"
// 	"project/internal/links"

// 	"github.com/asaskevich/govalidator"
// 	"github.com/gouniverse/cdn"
// 	"github.com/gouniverse/form"
// 	"github.com/gouniverse/hb"
// 	"github.com/gouniverse/router"
// 	"github.com/gouniverse/shopstore"
// 	"github.com/gouniverse/utils"
// 	"github.com/mingrammer/cfmt"
// 	"github.com/samber/lo"
// 	"github.com/spf13/cast"
// )

// // == CONTROLLER ==============================================================

// type productUpdateController struct{}

// var _ router.HTMLControllerInterface = (*productUpdateController)(nil)

// // == CONSTRUCTOR =============================================================

// func NewProductUpdateController() *productUpdateController {
// 	return &productUpdateController{}
// }

// func (controller *productUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
// 	data, errorMessage := controller.prepareDataAndValidate(r)

// 	if errorMessage != "" {
// 		return helpers.ToFlashError(w, r, errorMessage, links.NewAdminLinks().Home(map[string]string{}), 10)
// 	}

// 	if r.Method == http.MethodPost && data.action == "update-details" {
// 		return controller.formDetails(data).ToHTML()
// 	}

// 	if r.Method == http.MethodPost && data.action == "update-metadata" {
// 		return controller.formMetadata(data).ToHTML()
// 	}

// 	return layouts.NewAdminLayout(r, layouts.Options{
// 		Title:   "Edit Product | Shop",
// 		Content: controller.page(data),
// 		ScriptURLs: []string{
// 			cdn.Htmx_1_9_4(),
// 			cdn.Sweetalert2_10(),
// 		},
// 		Styles: []string{},
// 	}).ToHTML()
// }

// func (controller *productUpdateController) page(data productUpdateControllerData) hb.TagInterface {
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
// 		{
// 			Name: "Edit Product",
// 			URL:  shared.NewLinks().ProductUpdate(map[string]string{}),
// 		},
// 	})

// 	buttonDetailsSave := hb.Button().
// 		Class("btn btn-primary ms-2 float-end").
// 		Child(hb.I().Class("bi bi-save").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
// 		HTML("Save").
// 		HxInclude("#FormProductUpdate").
// 		HxPost(shared.NewLinks().ProductUpdate(map[string]string{
// 			"productID": data.productID,
// 			"action":    "update-details",
// 		})).
// 		HxTarget("#FormProductUpdate")

// 	buttonCancel := hb.Hyperlink().
// 		Class("btn btn-secondary ms-2 float-end").
// 		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
// 		HTML("Back").
// 		Href(shared.NewLinks().Products(map[string]string{}))

// 	buttonMetadataSave := hb.Button().
// 		Class("btn btn-primary ms-2 float-end").
// 		Child(hb.I().Class("bi bi-save").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
// 		HTML("Save").
// 		HxInclude("#FormProductMetadataUpdate").
// 		HxPost(shared.NewLinks().ProductUpdate(map[string]string{
// 			"productID": data.productID,
// 			"action":    "update-metadata",
// 		})).
// 		HxTarget("#FormProductMetadataUpdate")

// 	heading := hb.Heading1().
// 		HTML("Shop. Product. Edit Product").
// 		// Child(buttonSave).
// 		Child(buttonCancel)

// 	cardProductDetails := hb.Div().
// 		Class("card").
// 		Child(
// 			hb.Div().
// 				Class("card-header").
// 				Style(`display:flex;justify-content:space-between;align-items:center;`).
// 				Child(hb.Heading4().
// 					HTML("Product Details").
// 					Style("margin-bottom:0;display:inline-block;")).
// 				Child(buttonDetailsSave),
// 		).
// 		Child(
// 			hb.Div().
// 				Class("card-body").
// 				Child(controller.formDetails(data)))

// 	cardProductMetadata := hb.Div().
// 		Class("card").
// 		Child(
// 			hb.Div().
// 				Class("card-header").
// 				Style(`display:flex;justify-content:space-between;align-items:center;`).
// 				Child(hb.Heading4().
// 					HTML("Product Metadata").
// 					Style("margin-bottom:0;display:inline-block;")).
// 				Child(buttonMetadataSave),
// 		).
// 		Child(
// 			hb.Div().
// 				Class("card-body").
// 				Child(controller.formMetadata(data)))

// 	productTitle := hb.Heading2().
// 		Class("mb-3").
// 		Text("Product: ").
// 		Text(data.product.Title())

// 	return hb.Div().
// 		Class("container").
// 		Child(breadcrumbs).
// 		Child(hb.HR()).
// 		Child(heading).
// 		Child(productTitle).
// 		Child(cardProductDetails).
// 		Child(hb.BR()).
// 		Child(cardProductMetadata)
// }

// func (controller *productUpdateController) formDetails(data productUpdateControllerData) hb.TagInterface {
// 	fieldsDetails := []form.FieldInterface{
// 		form.NewField(form.FieldOptions{
// 			Label: "Status",
// 			Name:  "product_status",
// 			Type:  form.FORM_FIELD_TYPE_SELECT,
// 			Value: data.formStatus,
// 			Help:  `The status of the product.`,
// 			Options: []form.FieldOption{
// 				{
// 					Value: "- not selected -",
// 					Key:   "",
// 				},
// 				{
// 					Value: "Active",
// 					Key:   shopstore.PRODUCT_STATUS_ACTIVE,
// 				},
// 				{
// 					Value: "Disabled",
// 					Key:   shopstore.PRODUCT_STATUS_DISABLED,
// 				},
// 				{
// 					Value: "Draft",
// 					Key:   shopstore.PRODUCT_STATUS_DRAFT,
// 				},
// 			},
// 		}),
// 		form.NewField(form.FieldOptions{
// 			Label: "Title",
// 			Name:  "product_title",
// 			Type:  form.FORM_FIELD_TYPE_STRING,
// 			Value: data.formTitle,
// 			Help:  `The title of the product.`,
// 		}),
// 		form.NewField(form.FieldOptions{
// 			Label: "Description",
// 			Name:  "product_description",
// 			Type:  form.FORM_FIELD_TYPE_TEXTAREA,
// 			Value: data.formDescription,
// 			Help:  `The description of the product.`,
// 		}),
// 		form.NewField(form.FieldOptions{
// 			Label: "Price",
// 			Name:  "product_price",
// 			Type:  form.FORM_FIELD_TYPE_NUMBER,
// 			Value: data.formPrice,
// 			Help:  `The price of the product.`,
// 		}),
// 		form.NewField(form.FieldOptions{
// 			Label: "Quantity",
// 			Name:  "product_quantity",
// 			Type:  form.FORM_FIELD_TYPE_NUMBER,
// 			Value: data.formQuantity,
// 			Help:  `The quantity of the product that is available to purchase.`,
// 		}),
// 		form.NewField(form.FieldOptions{
// 			Label: "Admin Notes",
// 			Name:  "product_memo",
// 			Type:  form.FORM_FIELD_TYPE_TEXTAREA,
// 			Value: data.formMemo,
// 			Help:  "Admin notes for this product. These notes will not be visible to the public.",
// 		}),
// 		form.NewField(form.FieldOptions{
// 			Label:    "Product ID",
// 			Name:     "product_id",
// 			Type:     form.FORM_FIELD_TYPE_STRING,
// 			Value:    data.productID,
// 			Readonly: true,
// 			Help:     "The reference number (ID) of the product.",
// 		}),
// 	}

// 	formUserUpdate := form.NewForm(form.FormOptions{
// 		ID: "FormProductUpdate",
// 	})

// 	formUserUpdate.SetFields(fieldsDetails)

// 	if data.formErrorMessage != "" {
// 		formUserUpdate.AddField(form.NewField(form.FieldOptions{
// 			Type:  form.FORM_FIELD_TYPE_RAW,
// 			Value: hb.Swal(hb.SwalOptions{Icon: "error", Text: data.formErrorMessage}).ToHTML(),
// 		}))
// 	}

// 	if data.formSuccessMessage != "" {
// 		formUserUpdate.AddField(form.NewField(form.FieldOptions{
// 			Type:  form.FORM_FIELD_TYPE_RAW,
// 			Value: hb.Swal(hb.SwalOptions{Icon: "success", Text: data.formSuccessMessage}).ToHTML(),
// 		}))
// 	}

// 	return formUserUpdate.Build()
// }

// func (controller *productUpdateController) formMetadata(data productUpdateControllerData) hb.TagInterface {
// 	fieldsDetails := []form.FieldInterface{
// 		form.NewField(form.FieldOptions{
// 			Label:    "Product ID",
// 			Name:     "product_id",
// 			Type:     form.FORM_FIELD_TYPE_HIDDEN,
// 			Value:    data.productID,
// 			Readonly: true,
// 			//Help:     "The reference number (ID) of the product.",
// 		}),
// 	}

// 	metas := data.formMetas

// 	index := 0
// 	keys := lo.Keys(metas)
// 	slices.Sort(keys)
// 	for _, key := range keys {
// 		value := metas[key]
// 		background := lo.Ternary(index%2 == 0, "bg-light", "bg-white")
// 		fieldsMeta := []form.FieldInterface{
// 			form.NewField(form.FieldOptions{
// 				Type:  form.FORM_FIELD_TYPE_RAW,
// 				Help:  `Opening row`,
// 				Value: `<div id="Row` + cast.ToString(index) + `" class="row ` + background + ` py-2">`,
// 			}),
// 			form.NewField(form.FieldOptions{
// 				Type:  form.FORM_FIELD_TYPE_RAW,
// 				Help:  `Opening column 1`,
// 				Value: `<div class="col-3">`,
// 			}),
// 			form.NewField(form.FieldOptions{
// 				Label: `Key`,
// 				Name:  `product_meta[` + cast.ToString(index) + `][key]`,
// 				Type:  form.FORM_FIELD_TYPE_STRING,
// 				Value: key,
// 				// Help:  "The metadata value.",
// 			}),
// 			form.NewField(form.FieldOptions{
// 				Type:  form.FORM_FIELD_TYPE_RAW,
// 				Help:  `Closing column 1`,
// 				Value: `</div>`,
// 			}),
// 			form.NewField(form.FieldOptions{
// 				Type:  form.FORM_FIELD_TYPE_RAW,
// 				Help:  `Opening column 2`,
// 				Value: `<div class="col-8">`,
// 			}),
// 			form.NewField(form.FieldOptions{
// 				Label: `Value`,
// 				Name:  `product_meta[` + cast.ToString(index) + `][value]`,
// 				Type:  form.FORM_FIELD_TYPE_TEXTAREA,
// 				Value: value,
// 				// Help:  "The metadata value.",
// 			}),
// 			form.NewField(form.FieldOptions{
// 				Type:  form.FORM_FIELD_TYPE_RAW,
// 				Help:  `Closing column 2`,
// 				Value: `</div>`,
// 			}),
// 			form.NewField(form.FieldOptions{
// 				Type:  form.FORM_FIELD_TYPE_RAW,
// 				Help:  `Opening column 3`,
// 				Value: `<div class="col-1">`,
// 			}),
// 			form.NewField(form.FieldOptions{
// 				Type:  form.FORM_FIELD_TYPE_RAW,
// 				Value: `<button onclick="document.getElementById('Row` + cast.ToString(index) + `').innerHTML='';" type="button" class="btn btn-sm btn-danger">x</button>`,
// 				Help:  "Delete...",
// 			}),
// 			form.NewField(form.FieldOptions{
// 				Type:  form.FORM_FIELD_TYPE_RAW,
// 				Help:  `Closing column 3`,
// 				Value: `</div>`,
// 			}),
// 			form.NewField(form.FieldOptions{
// 				Type:  form.FORM_FIELD_TYPE_RAW,
// 				Help:  `Closing the row.`,
// 				Value: `</div>`,
// 			}),
// 		}

// 		fieldsDetails = append(fieldsDetails, fieldsMeta...)

// 		index++
// 	}

// 	fieldsNewMeta := []form.FieldInterface{
// 		form.NewField(form.FieldOptions{
// 			Type:  form.FORM_FIELD_TYPE_RAW,
// 			Value: `<hr />`,
// 		}),
// 		form.NewField(form.FieldOptions{
// 			Type:  form.FORM_FIELD_TYPE_RAW,
// 			Value: `<div class="row bg-info py-2">`,
// 		}),
// 		form.NewField(form.FieldOptions{
// 			Type:  form.FORM_FIELD_TYPE_RAW,
// 			Value: `<h3>New Meta</h3>`,
// 		}),
// 		form.NewField(form.FieldOptions{
// 			Type:  form.FORM_FIELD_TYPE_RAW,
// 			Value: `<div class="col-6">`,
// 		}),

// 		form.NewField(form.FieldOptions{
// 			Label: `Key`,
// 			Name:  `product_meta[` + cast.ToString(index) + `][key]`,
// 			Type:  form.FORM_FIELD_TYPE_STRING,
// 			Value: "",
// 			// Help:  "The metadata value.",
// 		}),

// 		form.NewField(form.FieldOptions{
// 			Type:  form.FORM_FIELD_TYPE_RAW,
// 			Value: `</div>`,
// 		}),

// 		form.NewField(form.FieldOptions{
// 			Type:  form.FORM_FIELD_TYPE_RAW,
// 			Value: `<div class="col-6">`,
// 		}),

// 		form.NewField(form.FieldOptions{
// 			Label: `Value`,
// 			Name:  `product_meta[` + cast.ToString(index) + `][value]`,
// 			Type:  form.FORM_FIELD_TYPE_STRING,
// 			Value: "",
// 			// Help:  "The metadata value.",
// 		}),

// 		form.NewField(form.FieldOptions{
// 			Type:  form.FORM_FIELD_TYPE_RAW,
// 			Value: `</div>`,
// 		}),

// 		form.NewField(form.FieldOptions{
// 			Type:  form.FORM_FIELD_TYPE_RAW,
// 			Value: `</div>`,
// 		}),
// 	}

// 	fieldsDetails = append(fieldsDetails, fieldsNewMeta...)

// 	formMetadataUpdate := form.NewForm(form.FormOptions{
// 		ID:     "FormProductMetadataUpdate",
// 		Fields: fieldsDetails,
// 	})

// 	if data.formErrorMessage != "" {
// 		formMetadataUpdate.AddField(form.NewField(form.FieldOptions{
// 			Type:  form.FORM_FIELD_TYPE_RAW,
// 			Value: hb.Swal(hb.SwalOptions{Icon: "error", Text: data.formErrorMessage}).ToHTML(),
// 		}))
// 	}

// 	if data.formSuccessMessage != "" {
// 		formMetadataUpdate.AddField(form.NewField(form.FieldOptions{
// 			Type:  form.FORM_FIELD_TYPE_RAW,
// 			Value: hb.Swal(hb.SwalOptions{Icon: "success", Text: data.formSuccessMessage}).ToHTML(),
// 		}))
// 	}

// 	return formMetadataUpdate.Build()
// }

// func (controller *productUpdateController) saveProductDetails(r *http.Request, data productUpdateControllerData) (d productUpdateControllerData, errorMessage string) {
// 	data.formDescription = utils.Req(r, "product_description", "")
// 	data.formMemo = utils.Req(r, "product_memo", "")
// 	data.formPrice = utils.Req(r, "product_price", "")
// 	data.formQuantity = utils.Req(r, "product_quantity", "")
// 	data.formStatus = utils.Req(r, "product_status", "")
// 	data.formTitle = utils.Req(r, "product_title", "")

// 	if data.formStatus == "" {
// 		data.formErrorMessage = "Status is required"
// 		return data, ""
// 	}

// 	if data.formTitle == "" {
// 		data.formErrorMessage = "Title is required"
// 		return data, ""
// 	}

// 	if data.formPrice == "" {
// 		data.formErrorMessage = "Price is required"
// 		return data, ""
// 	}

// 	if data.formQuantity == "" {
// 		data.formErrorMessage = "Quantity is required"
// 		return data, ""
// 	}

// 	if !govalidator.IsFloat(data.formPrice) {
// 		data.formErrorMessage = "Price must be numeric"
// 		return data, ""
// 	}

// 	if !govalidator.IsInt(data.formQuantity) {
// 		data.formErrorMessage = "Quantity must be numeric"
// 		return data, ""
// 	}

// 	price, _ := strconv.ParseFloat(data.formPrice, 64)

// 	if price < 0 {
// 		data.formErrorMessage = "Price cannot be negative"
// 		return data, ""
// 	}

// 	quantity, _ := strconv.ParseInt(data.formQuantity, 10, 64)

// 	if quantity < 0 {
// 		data.formErrorMessage = "Quantity cannot be negative"
// 		return data, ""
// 	}

// 	data.product.SetDescription(data.formDescription)
// 	data.product.SetMemo(data.formMemo)
// 	data.product.SetQuantity(data.formQuantity)
// 	data.product.SetPrice(data.formPrice)
// 	data.product.SetStatus(data.formStatus)
// 	data.product.SetTitle(data.formTitle)

// 	err := config.ShopStore.ProductUpdate(data.product)

// 	if err != nil {
// 		config.LogStore.ErrorWithContext("At productUpdateController > prepareDataAndValidate", err.Error())
// 		data.formErrorMessage = "System error. Saving details failed"
// 		return data, ""
// 	}

// 	data.formSuccessMessage = "Product saved successfully"

// 	return data, ""
// }

// func ReqArrayOfMaps(r *http.Request, key string, defaultValue []map[string]string) []map[string]string {
// 	all := utils.ReqAll(r)

// 	reqArrayOfMaps := []map[string]string{}

// 	if all == nil {
// 		return reqArrayOfMaps
// 	}

// 	mapIndexMap := map[string]map[string]string{}

// 	// Iterate through all the parameters
// 	for k, v := range all {
// 		if !strings.HasPrefix(k, key+"[") {
// 			continue
// 		}
// 		if !strings.HasSuffix(k, "]") {
// 			continue
// 		}
// 		if !strings.Contains(k, "][") {
// 			continue
// 		}
// 		mapValue := v[0]

// 		str := strings.TrimSuffix(strings.TrimPrefix(k, key+"["), "]")
// 		split := strings.Split(str, "][")
// 		if len(split) != 2 {
// 			// Handle invalid format
// 			continue
// 		}

// 		index, key := split[0], split[1]

// 		if lo.HasKey(mapIndexMap, index) {
// 			mapIndexMap[index][key] = mapValue
// 		} else {
// 			mapIndexMap[index] = map[string]string{
// 				key: mapValue,
// 			}
// 		}
// 	}

// 	for _, v := range mapIndexMap {
// 		if v == nil {
// 			continue
// 		}
// 		reqArrayOfMaps = append(reqArrayOfMaps, v)
// 	}

// 	return reqArrayOfMaps
// }

// func (controller *productUpdateController) saveProductMetadata(r *http.Request, data productUpdateControllerData) (d productUpdateControllerData, errorMessage string) {
// 	// data.formMetas = utils.ReqMap(r, "product_meta")
// 	metas := ReqArrayOfMaps(r, "product_meta", []map[string]string{})

// 	cfmt.Infoln(metas)

// 	productMetas := map[string]string{}

// 	lo.ForEach(metas, func(meta map[string]string, index int) {
// 		metaKey := strings.TrimSpace(meta["key"])
// 		metaValue := strings.TrimSpace(meta["value"])
// 		if metaKey == "" {
// 			return
// 		}
// 		productMetas[metaKey] = metaValue
// 	})

// 	data.formMetas = productMetas

// 	cfmt.Infoln(data.formMetas)

// 	if data.formMetas == nil {
// 		data.formErrorMessage = "Metadata is required"
// 		return data, ""
// 	}

// 	data.product.SetMetas(data.formMetas)

// 	err := config.ShopStore.ProductUpdate(data.product)

// 	if err != nil {
// 		config.LogStore.ErrorWithContext("At productUpdateController > prepareDataAndValidate", err.Error())
// 		data.formErrorMessage = "System error. Saving metas failed"
// 		return data, ""
// 	}

// 	data.formSuccessMessage = "Metadata saved successfully"

// 	return data, ""
// }

// func (controller *productUpdateController) prepareDataAndValidate(r *http.Request) (data productUpdateControllerData, errorMessage string) {
// 	data.action = utils.Req(r, "action", "")
// 	data.productID = utils.Req(r, "product_id", "")

// 	if data.productID == "" {
// 		return data, "Product ID is required"
// 	}

// 	product, err := config.ShopStore.ProductFindByID(data.productID)

// 	if err != nil {
// 		config.LogStore.ErrorWithContext("At productUpdateController > prepareDataAndValidate", err.Error())
// 		return data, "Product not found"
// 	}

// 	if product == nil {
// 		return data, "Product not found"
// 	}

// 	data.product = product

// 	metas, err := product.Metas()

// 	if err != nil {
// 		config.LogStore.ErrorWithContext("At productUpdateController > prepareDataAndValidate", err.Error())
// 		return data, "Product metas not found"
// 	}

// 	data.formMemo = data.product.Memo()
// 	data.formStatus = data.product.Status()
// 	data.formTitle = data.product.Title()
// 	data.formDescription = data.product.Description()
// 	data.formPrice = data.product.Price()
// 	data.formQuantity = data.product.Quantity()
// 	data.formMetas = metas

// 	if r.Method != http.MethodPost {
// 		return data, ""
// 	}

// 	if data.action == "update-details" {
// 		return controller.saveProductDetails(r, data)
// 	}

// 	if data.action == "update-metadata" {
// 		return controller.saveProductMetadata(r, data)
// 	}

// 	return data, "action is required"
// }

// type productUpdateControllerData struct {
// 	action    string
// 	productID string
// 	product   shopstore.ProductInterface

// 	formErrorMessage   string
// 	formSuccessMessage string
// 	formDescription    string
// 	formMemo           string
// 	formMetas          map[string]string
// 	formQuantity       string
// 	formPrice          string
// 	formStatus         string
// 	formTitle          string
// }
