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

type Order struct {
	dataobject.DataObject
}

// == CONSTRUCTORS =============================================================

func NewOrder() *Order {
	o := (&Order{}).
		SetID(uid.HumanUid()).
		SetStatus(ORDER_STATUS_PENDING).
		SetQuantity(1). // By default 1
		SetPrice(0).    // Free. By default
		SetMemo("").
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetDeletedAt(sb.NULL_DATETIME)

	o.SetMetas(map[string]string{})

	return o
}

func NewOrderFromExistingData(data map[string]string) *Order {
	o := &Order{}
	o.Hydrate(data)
	return o
}

// == METHODS ==================================================================

func (order *Order) IsAwaitingFulfillment() bool {
	return order.Status() == ORDER_STATUS_AWAITING_FULFILLMENT
}

func (order *Order) IsAwaitingPayment() bool {
	return order.Status() == ORDER_STATUS_AWAITING_PAYMENT
}

func (order *Order) IsAwaitingPickup() bool {
	return order.Status() == ORDER_STATUS_AWAITING_PICKUP
}

func (order *Order) IsAwaitingShipment() bool {
	return order.Status() == ORDER_STATUS_AWAITING_SHIPMENT
}

func (order *Order) IsCancelled() bool {
	return order.Status() == ORDER_STATUS_CANCELLED
}

func (order *Order) IsCompleted() bool {
	return order.Status() == ORDER_STATUS_COMPLETED
}

func (order *Order) IsDeclined() bool {
	return order.Status() == ORDER_STATUS_DECLINED
}

func (order *Order) IsDisputed() bool {
	return order.Status() == ORDER_STATUS_DISPUTED
}

func (order *Order) IsManualVerificationRequired() bool {
	return order.Status() == ORDER_STATUS_MANUAL_VERIFICATION_REQUIRED
}

func (order *Order) IsPending() bool {
	return order.Status() == ORDER_STATUS_PENDING
}

func (order *Order) IsRefunded() bool {
	return order.Status() == ORDER_STATUS_REFUNDED
}

func (order *Order) IsShipped() bool {
	return order.Status() == ORDER_STATUS_SHIPPED
}

// == GETTERS & SETTERS ========================================================

func (order *Order) CreatedAt() string {
	return order.Get(COLUMN_CREATED_AT)
}

func (order *Order) CreatedAtCarbon() carbon.Carbon {
	return carbon.NewCarbon().Parse(order.CreatedAt(), carbon.UTC)
}

func (order *Order) SetCreatedAt(createdAt string) *Order {
	order.Set(COLUMN_CREATED_AT, createdAt)
	return order
}

func (order *Order) DeletedAt() string {
	return order.Get(COLUMN_DELETED_AT)
}

func (order *Order) DeletedAtCarbon() carbon.Carbon {
	return carbon.NewCarbon().Parse(order.DeletedAt(), carbon.UTC)
}

func (order *Order) SetDeletedAt(deletedAt string) *Order {
	order.Set(COLUMN_DELETED_AT, deletedAt)
	return order
}

func (order *Order) ID() string {
	return order.Get(COLUMN_ID)
}

func (order *Order) SetID(id string) *Order {
	order.Set(COLUMN_ID, id)
	return order
}

func (order *Order) Memo() string {
	return order.Get(COLUMN_MEMO)
}

func (order *Order) SetMemo(memo string) *Order {
	order.Set(COLUMN_MEMO, memo)
	return order
}

func (order *Order) Metas() (map[string]string, error) {
	metasStr := order.Get(COLUMN_METAS)

	if metasStr == "" {
		metasStr = "{}"
	}

	metasJson, errJson := utils.FromJSON(metasStr, map[string]string{})
	if errJson != nil {
		return map[string]string{}, errJson
	}

	return maputils.MapStringAnyToMapStringString(metasJson.(map[string]any)), nil
}

func (order *Order) Meta(name string) string {
	metas, err := order.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

func (order *Order) SetMeta(name string, value string) error {
	return order.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (order *Order) SetMetas(metas map[string]string) error {
	mapString, err := utils.ToJSON(metas)
	if err != nil {
		return err
	}
	order.Set(COLUMN_METAS, mapString)
	return nil
}

func (order *Order) UpsertMetas(metas map[string]string) error {
	currentMetas, err := order.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return order.SetMetas(currentMetas)
}

func (order *Order) Status() string {
	return order.Get(COLUMN_STATUS)
}

func (order *Order) SetStatus(status string) *Order {
	order.Set(COLUMN_STATUS, status)
	return order
}

func (order *Order) Price() string {
	return order.Get(COLUMN_PRICE)
}

func (order *Order) PriceFloat() float64 {
	price, _ := utils.ToFloat(order.Get(COLUMN_PRICE))
	return price
}

func (order *Order) SetPrice(price float64) *Order {
	order.Set(COLUMN_PRICE, utils.ToString(price))
	return order
}

func (order *Order) Quantity() string {
	return order.Get(COLUMN_QUANTITY)
}

func (order *Order) QuantityInt() int64 {
	quantity, _ := utils.ToInt(order.Quantity())
	return quantity
}

func (order *Order) SetQuantity(quantity int) *Order {
	order.Set(COLUMN_QUANTITY, utils.ToString(quantity))
	return order
}

func (order *Order) UpdatedAt() string {
	return order.Get(COLUMN_UPDATED_AT)
}

func (order *Order) UpdatedAtCarbon() carbon.Carbon {
	return carbon.NewCarbon().Parse(order.UpdatedAt(), carbon.UTC)
}

func (order *Order) SetUpdatedAt(updatedAt string) *Order {
	order.Set(COLUMN_UPDATED_AT, updatedAt)
	return order
}

func (order *Order) UserID() string {
	return order.Get(COLUMN_USER_ID)
}

func (order *Order) SetUserID(id string) *Order {
	order.Set(COLUMN_USER_ID, id)
	return order
}
