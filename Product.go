package shopstore

import (
	"github.com/golang-module/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/maputils"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
	"github.com/gouniverse/utils"
)

// == CLASS ====================================================================

type Product struct {
	dataobject.DataObject
}

// == INTERFACES ===============================================================

var _ ProductInterface = (*Product)(nil)

// == CONSTRUCTORS =============================================================

func NewProduct() ProductInterface {
	o := (&Product{}).
		SetID(uid.HumanUid()).
		SetStatus(PRODUCT_STATUS_DRAFT).
		SetTitle("").
		SetDescription("").
		SetQuantityInt(0). // By default 1
		SetPriceFloat(0).  // Free. By default
		SetMemo("").
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetDeletedAt(sb.MAX_DATETIME)

	o.SetMetas(map[string]string{})

	return o
}

func NewProductFromExistingData(data map[string]string) ProductInterface {
	o := &Product{}
	o.Hydrate(data)
	return o
}

// == METHODS ==================================================================

func (product *Product) IsActive() bool {
	return product.Status() == PRODUCT_STATUS_ACTIVE
}

func (product *Product) IsDisabled() bool {
	return product.Status() == PRODUCT_STATUS_DISABLED
}

func (product *Product) IsDraft() bool {
	return product.Status() == PRODUCT_STATUS_DRAFT
}

// == GETTERS & SETTERS ========================================================

func (product *Product) CreatedAt() string {
	return product.Get(COLUMN_CREATED_AT)
}

func (product *Product) CreatedAtCarbon() carbon.Carbon {
	return carbon.NewCarbon().Parse(product.CreatedAt(), carbon.UTC)
}

func (product *Product) SetCreatedAt(createdAt string) ProductInterface {
	product.Set(COLUMN_CREATED_AT, createdAt)
	return product
}

func (product *Product) DeletedAt() string {
	return product.Get(COLUMN_DELETED_AT)
}

func (product *Product) DeletedAtCarbon() carbon.Carbon {
	return carbon.NewCarbon().Parse(product.DeletedAt(), carbon.UTC)
}

func (product *Product) SetDeletedAt(deletedAt string) ProductInterface {
	product.Set(COLUMN_DELETED_AT, deletedAt)
	return product
}

func (product *Product) ID() string {
	return product.Get(COLUMN_ID)
}

func (product *Product) SetID(id string) ProductInterface {
	product.Set(COLUMN_ID, id)
	return product
}

func (product *Product) Memo() string {
	return product.Get(COLUMN_MEMO)
}

func (product *Product) SetMemo(memo string) ProductInterface {
	product.Set(COLUMN_MEMO, memo)
	return product
}

func (product *Product) Metas() (map[string]string, error) {
	metasStr := product.Get(COLUMN_METAS)

	if metasStr == "" {
		metasStr = "{}"
	}

	metasJson, errJson := utils.FromJSON(metasStr, map[string]string{})
	if errJson != nil {
		return map[string]string{}, errJson
	}

	return maputils.MapStringAnyToMapStringString(metasJson.(map[string]any)), nil
}

func (product *Product) Meta(name string) string {
	metas, err := product.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

func (product *Product) SetMeta(name string, value string) error {
	return product.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (product *Product) SetMetas(metas map[string]string) error {
	mapString, err := utils.ToJSON(metas)
	if err != nil {
		return err
	}
	product.Set(COLUMN_METAS, mapString)
	return nil
}

func (product *Product) UpsertMetas(metas map[string]string) error {
	currentMetas, err := product.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return product.SetMetas(currentMetas)
}

func (product *Product) Status() string {
	return product.Get(COLUMN_STATUS)
}

func (product *Product) SetStatus(status string) ProductInterface {
	product.Set(COLUMN_STATUS, status)
	return product
}

func (product *Product) Price() string {
	return product.Get(COLUMN_PRICE)
}

func (product *Product) SetPrice(price string) ProductInterface {
	product.Set(COLUMN_PRICE, price)
	return product
}

func (product *Product) PriceFloat() float64 {
	price, _ := utils.ToFloat(product.Get(COLUMN_PRICE))
	return price
}

func (product *Product) SetPriceFloat(price float64) ProductInterface {
	product.SetPrice(utils.ToString(price))
	return product
}

func (product *Product) Quantity() string {
	return product.Get(COLUMN_QUANTITY)
}

func (product *Product) SetQuantity(quantity string) ProductInterface {
	product.Set(COLUMN_QUANTITY, quantity)
	return product
}

func (product *Product) QuantityInt() int64 {
	quantity, _ := utils.ToInt(product.Quantity())
	return quantity
}

func (product *Product) SetQuantityInt(quantity int64) ProductInterface {
	product.SetQuantity(utils.ToString(quantity))
	return product
}

func (product *Product) Title() string {
	return product.Get(COLUMN_TITLE)
}

func (product *Product) SetTitle(title string) ProductInterface {
	product.Set(COLUMN_TITLE, title)
	return product
}

func (product *Product) Description() string {
	return product.Get(COLUMN_DESCRIPTION)
}

func (product *Product) SetDescription(description string) ProductInterface {
	product.Set(COLUMN_DESCRIPTION, description)
	return product
}

func (product *Product) UpdatedAt() string {
	return product.Get(COLUMN_UPDATED_AT)
}

func (product *Product) UpdatedAtCarbon() carbon.Carbon {
	return carbon.NewCarbon().Parse(product.UpdatedAt(), carbon.UTC)
}

func (product *Product) SetUpdatedAt(updatedAt string) ProductInterface {
	product.Set(COLUMN_UPDATED_AT, updatedAt)
	return product
}
