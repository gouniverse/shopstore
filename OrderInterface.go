package shopstore

import "github.com/golang-module/carbon/v2"

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
