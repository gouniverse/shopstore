package admin

// import (
// 	"net/http"
// 	"project/config"
// 	"project/controllers/admin/shop/shared"
// 	"project/internal/helpers"

// 	"github.com/gouniverse/bs"
// 	"github.com/gouniverse/hb"
// 	"github.com/gouniverse/router"
// 	"github.com/gouniverse/shopstore"
// 	"github.com/gouniverse/utils"
// )

// type productDeleteController struct{}

// var _ router.HTMLControllerInterface = (*productDeleteController)(nil)

// type productDeleteControllerData struct {
// 	productID      string
// 	product        shopstore.ProductInterface
// 	successMessage string
// 	//errorMessage   string
// }

// func NewProductDeleteController() *productDeleteController {
// 	return &productDeleteController{}
// }

// func (controller productDeleteController) Handler(w http.ResponseWriter, r *http.Request) string {
// 	data, errorMessage := controller.prepareDataAndValidate(r)

// 	if errorMessage != "" {
// 		return hb.Swal(hb.SwalOptions{
// 			Icon: "error",
// 			Text: errorMessage,
// 		}).ToHTML()
// 	}

// 	if data.successMessage != "" {
// 		return hb.Wrap().
// 			Child(hb.Swal(hb.SwalOptions{
// 				Icon: "success",
// 				Text: data.successMessage,
// 			})).
// 			Child(hb.Script("setTimeout(() => {window.location.href = window.location.href}, 2000)")).
// 			ToHTML()
// 	}

// 	return controller.
// 		modal(data).
// 		ToHTML()
// }

// func (controller *productDeleteController) modal(data productDeleteControllerData) hb.TagInterface {
// 	submitUrl := shared.NewLinks().ProductDelete(map[string]string{
// 		"product_id": data.productID,
// 	})

// 	modalID := "ModalProductDelete"
// 	modalBackdropClass := "ModalBackdrop"

// 	formGroupProductId := hb.Input().
// 		Type(hb.TYPE_HIDDEN).
// 		Name("product_id").
// 		Value(data.productID)

// 	buttonDelete := hb.Button().
// 		HTML("Delete").
// 		Class("btn btn-primary float-end").
// 		HxInclude("#Modal" + modalID).
// 		HxPost(submitUrl).
// 		HxSelectOob("#ModalProductDelete").
// 		HxTarget("body").
// 		HxSwap("beforeend")

// 	modalCloseScript := `closeModal` + modalID + `();`

// 	modalHeading := hb.Heading5().HTML("Delete Product").Style(`margin:0px;`)

// 	modalClose := hb.Button().Type("button").
// 		Class("btn-close").
// 		Data("bs-dismiss", "modal").
// 		OnClick(modalCloseScript)

// 	jsCloseFn := `function closeModal` + modalID + `() {document.getElementById('ModalProductDelete').remove();[...document.getElementsByClassName('` + modalBackdropClass + `')].forEach(el => el.remove());}`

// 	modal := bs.Modal().
// 		ID(modalID).
// 		Class("fade show").
// 		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
// 		Child(hb.Script(jsCloseFn)).
// 		Child(bs.ModalDialog().
// 			Child(bs.ModalContent().
// 				Child(
// 					bs.ModalHeader().
// 						Child(modalHeading).
// 						Child(modalClose)).
// 				Child(
// 					bs.ModalBody().
// 						Child(hb.Paragraph().Text("Are you sure you want to delete this product?").Style(`margin-bottom:20px;color:red;`)).
// 						Child(hb.Paragraph().Text("This action cannot be undone.")).
// 						Child(formGroupProductId)).
// 				Child(bs.ModalFooter().
// 					Style(`display:flex;justify-content:space-between;`).
// 					Child(
// 						hb.Button().HTML("Close").
// 							Class("btn btn-secondary float-start").
// 							Data("bs-dismiss", "modal").
// 							OnClick(modalCloseScript)).
// 					Child(buttonDelete)),
// 			))

// 	backdrop := hb.Div().Class(modalBackdropClass).
// 		Class("modal-backdrop fade show").
// 		Style("display:block;z-index:1000;")

// 	return hb.Wrap().
// 		Children([]hb.TagInterface{
// 			modal,
// 			backdrop,
// 		})
// }

// func (controller *productDeleteController) prepareDataAndValidate(r *http.Request) (data productDeleteControllerData, errorMessage string) {
// 	authUser := helpers.GetAuthUser(r)
// 	data.productID = utils.Req(r, "product_id", "")

// 	if authUser == nil {
// 		return data, "You are not logged in. Please login to continue."
// 	}

// 	if data.productID == "" {
// 		return data, "product id is required"
// 	}

// 	product, err := config.ShopStore.ProductFindByID(data.productID)

// 	if err != nil {
// 		config.LogStore.ErrorWithContext("Error. At productDeleteController > prepareDataAndValidate", err.Error())
// 		return data, "Product not found"
// 	}

// 	if product == nil {
// 		return data, "Product not found"
// 	}

// 	data.product = product

// 	if r.Method != "POST" {
// 		return data, ""
// 	}

// 	err = config.ShopStore.ProductSoftDelete(product)

// 	if err != nil {
// 		config.LogStore.ErrorWithContext("Error. At productDeleteController > prepareDataAndValidate", err.Error())
// 		return data, "Deleting product failed. Please contact an administrator."
// 	}

// 	data.successMessage = "product deleted successfully."

// 	return data, ""

// }
