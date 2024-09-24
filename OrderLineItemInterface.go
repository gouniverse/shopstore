package shopstore

import "github.com/golang-module/carbon/v2"

type OrderLineItemInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	CreatedAt() string
	CreatedAtCarbon() carbon.Carbon
	SetCreatedAt(createdAt string) OrderLineItemInterface

	DeletedAt() string
	DeletedAtCarbon() carbon.Carbon
	SetDeletedAt(deletedAt string) OrderLineItemInterface

	ID() string
	SetID(id string) OrderLineItemInterface

	Memo() string
	SetMemo(memo string) OrderLineItemInterface

	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	Meta(name string) string
	SetMeta(name string, value string) error
	UpsertMetas(metas map[string]string) error

	OrderID() string
	SetOrderID(orderID string) OrderLineItemInterface

	Price() string
	SetPrice(price string) OrderLineItemInterface

	PriceFloat() float64
	SetPriceFloat(price float64) OrderLineItemInterface

	ProductID() string
	SetProductID(productID string) OrderLineItemInterface

	Quantity() string
	SetQuantity(quantity string) OrderLineItemInterface

	QuantityInt() int64
	SetQuantityInt(quantity int64) OrderLineItemInterface

	Status() string
	SetStatus(status string) OrderLineItemInterface

	Title() string
	SetTitle(title string) OrderLineItemInterface

	UpdatedAt() string
	UpdatedAtCarbon() carbon.Carbon
	SetUpdatedAt(updatedAt string) OrderLineItemInterface
}
