package shopstore

import "github.com/golang-module/carbon/v2"

type ProductInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	CreatedAt() string
	CreatedAtCarbon() carbon.Carbon
	SetCreatedAt(createdAt string) ProductInterface
	DeletedAt() string
	DeletedAtCarbon() carbon.Carbon
	SetDeletedAt(deletedAt string) ProductInterface
	Description() string
	SetDescription(description string) ProductInterface
	ID() string
	SetID(id string) ProductInterface
	Memo() string
	SetMemo(memo string) ProductInterface
	Meta(name string) string
	SetMeta(name string, value string) error
	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	Price() string
	SetPrice(price string) ProductInterface
	PriceFloat() float64
	SetPriceFloat(price float64) ProductInterface
	Quantity() string
	SetQuantity(quantity string) ProductInterface
	QuantityInt() int64
	SetQuantityInt(quantity int64) ProductInterface
	Status() string
	SetStatus(status string) ProductInterface
	Title() string
	SetTitle(title string) ProductInterface
	UpdatedAt() string
	UpdatedAtCarbon() carbon.Carbon
	SetUpdatedAt(updatedAt string) ProductInterface
}
