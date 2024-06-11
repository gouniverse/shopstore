package shopstore

import (
	"github.com/golang-module/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
	"github.com/gouniverse/utils"
)

const ORDER_STATUS_PENDING = "pending"
const ORDER_STATUS_PAID = "paid"
const ORDER_STATUS_CANCELLED = "cancelled"

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

	return o
}

func NewOrderFromExistingData(data map[string]string) *Order {
	o := &Order{}
	o.Hydrate(data)
	return o
}

// == METHODS ==================================================================

// func (order *Order) Update() error {
// 	return NewUserService().OrderUpdate(order)
// }

func (order *Order) IsPaid() bool {
	return order.Get("status") == ORDER_STATUS_PAID
}

func (order *Order) IsPending() bool {
	return order.Get("status") == ORDER_STATUS_PENDING
}

func (order *Order) IsCancelled() bool {
	return order.Get("status") == ORDER_STATUS_CANCELLED
}

// == GETTERS & SETTERS ========================================================

func (order *Order) CreatedAt() string {
	return order.Get("created_at")
}

func (order *Order) SetCreatedAt(createdAt string) *Order {
	order.Set("created_at", createdAt)
	return order
}

func (order *Order) DeletedAt() string {
	return order.Get("deleted_at")
}

func (order *Order) SetDeletedAt(deletedAt string) *Order {
	order.Set("deleted_at", deletedAt)
	return order
}

func (order *Order) ExamID() string {
	return order.Get("exam_id")
}

func (order *Order) SetExamID(id string) *Order {
	order.Set("exam_id", id)
	return order
}

func (order *Order) ID() string {
	return order.Get("id")
}

func (order *Order) SetID(id string) *Order {
	order.Set("id", id)
	return order
}

func (order *Order) Memo() string {
	return order.Get("memo")
}

func (order *Order) SetMemo(memo string) *Order {
	order.Set("memo", memo)
	return order
}

func (order *Order) Status() string {
	return order.Get("status")
}

func (order *Order) SetStatus(status string) *Order {
	order.Set("status", status)
	return order
}

func (order *Order) Price() string {
	return order.Get("price")
}

func (order *Order) PriceFloat() float64 {
	price, _ := utils.ToFloat(order.Get("price"))
	return price
}

func (order *Order) SetPrice(price float64) *Order {
	order.Set("price", utils.ToString(price))
	return order
}

func (order *Order) Quantity() string {
	return order.Get("quantity")
}

func (order *Order) QuantityInt() int64 {
	quantity, _ := utils.ToInt(order.Get("quantity"))
	return quantity
}

func (order *Order) SetQuantity(quantity int) *Order {
	order.Set("quantity", utils.ToString(quantity))
	return order
}

func (order *Order) UpdatedAt() string {
	return order.Get("updated_at")
}

func (order *Order) SetUpdatedAt(updatedAt string) *Order {
	order.Set("updated_at", updatedAt)
	return order
}

func (order *Order) UserID() string {
	return order.Get("user_id")
}

func (order *Order) SetUserID(id string) *Order {
	order.Set("user_id", id)
	return order
}
