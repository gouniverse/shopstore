package shopstore

import (
	"log/slog"

	"github.com/dromara/carbon/v2"
)

type DiscountInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	Amount() float64
	SetAmount(amount float64) DiscountInterface
	Code() string
	SetCode(code string) DiscountInterface
	CreatedAt() string
	CreatedAtCarbon() carbon.Carbon
	SetCreatedAt(createdAt string) DiscountInterface
	DeletedAt() string
	DeletedAtCarbon() carbon.Carbon
	SetDeletedAt(deletedAt string) DiscountInterface
	Description() string
	SetDescription(description string) DiscountInterface
	EndsAt() string
	EndsAtCarbon() carbon.Carbon
	SetEndsAt(endsAt string) DiscountInterface
	ID() string
	SetID(id string) DiscountInterface
	StartsAt() string
	StartsAtCarbon() carbon.Carbon
	SetStartsAt(startsAt string) DiscountInterface
	Status() string
	SetStatus(status string) DiscountInterface
	Title() string
	SetTitle(title string) DiscountInterface
	Type() string
	SetType(type_ string) DiscountInterface
	UpdatedAt() string
	UpdatedAtCarbon() carbon.Carbon
	SetUpdatedAt(updatedAt string) DiscountInterface
}

type OrderInterface interface {
	// Inherited from DataObject
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	// Methods
	IsAwaitingFulfillment() bool
	IsAwaitingPayment() bool
	IsAwaitingPickup() bool
	IsAwaitingShipment() bool
	IsCancelled() bool
	IsCompleted() bool
	IsDeclined() bool
	IsDisputed() bool
	IsManualVerificationRequired() bool
	IsPending() bool
	IsRefunded() bool
	IsShipped() bool

	// Setters and Getters
	CreatedAt() string
	CreatedAtCarbon() carbon.Carbon
	SetCreatedAt(createdAt string) OrderInterface

	CustomerID() string
	SetCustomerID(customerID string) OrderInterface

	DeletedAt() string
	DeletedAtCarbon() carbon.Carbon
	SetDeletedAt(deletedAt string) OrderInterface

	ID() string
	SetID(id string) OrderInterface

	Memo() string
	SetMemo(memo string) OrderInterface

	Meta(name string) string
	SetMeta(name string, value string) error
	Metas() (map[string]string, error)
	SetMetas(metas map[string]string) error
	UpsertMetas(metas map[string]string) error

	Price() string
	SetPrice(price string) OrderInterface
	PriceFloat() float64
	SetPriceFloat(price float64) OrderInterface

	Quantity() string
	SetQuantity(quantity string) OrderInterface
	QuantityInt() int64
	SetQuantityInt(quantity int64) OrderInterface

	Status() string
	SetStatus(status string) OrderInterface

	UpdatedAt() string
	UpdatedAtCarbon() carbon.Carbon
	SetUpdatedAt(updatedAt string) OrderInterface
}

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
	UpsertMetas(metas map[string]string) error

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

type StoreInterface interface {
	AutoMigrate() error
	EnableDebug(debug bool, sqlLogger ...*slog.Logger)

	DiscountCount(options DiscountQueryOptions) (int64, error)
	DiscountCreate(discount DiscountInterface) error
	DiscountDelete(discount DiscountInterface) error
	DiscountDeleteByID(discountID string) error
	DiscountFindByID(discountID string) (DiscountInterface, error)
	DiscountFindByCode(code string) (DiscountInterface, error)
	DiscountList(options DiscountQueryOptions) ([]DiscountInterface, error)
	DiscountSoftDelete(discount DiscountInterface) error
	DiscountSoftDeleteByID(discountID string) error
	DiscountUpdate(discount DiscountInterface) error

	OrderCount(options OrderQueryOptions) (int64, error)
	OrderCreate(order OrderInterface) error
	OrderDelete(order OrderInterface) error
	OrderDeleteByID(id string) error
	OrderFindByID(id string) (OrderInterface, error)
	OrderList(options OrderQueryOptions) ([]OrderInterface, error)
	OrderSoftDelete(order OrderInterface) error
	OrderSoftDeleteByID(id string) error
	OrderUpdate(order OrderInterface) error

	OrderLineItemCount(options OrderLineItemQueryOptions) (int64, error)
	OrderLineItemCreate(orderLineItem OrderLineItemInterface) error
	OrderLineItemDelete(orderLineItem OrderLineItemInterface) error
	OrderLineItemDeleteByID(id string) error
	OrderLineItemFindByID(id string) (OrderLineItemInterface, error)
	OrderLineItemList(options OrderLineItemQueryOptions) ([]OrderLineItemInterface, error)
	OrderLineItemSoftDelete(orderLineItem OrderLineItemInterface) error
	OrderLineItemSoftDeleteByID(id string) error
	OrderLineItemUpdate(orderLineItem OrderLineItemInterface) error

	ProductCount(options ProductQueryOptions) (int64, error)
	ProductCreate(product ProductInterface) error
	ProductDelete(product ProductInterface) error
	ProductDeleteByID(productID string) error
	ProductFindByID(productID string) (ProductInterface, error)
	ProductList(options ProductQueryOptions) ([]ProductInterface, error)
	ProductSoftDelete(product ProductInterface) error
	ProductSoftDeleteByID(productID string) error
	ProductUpdate(product ProductInterface) error
}
