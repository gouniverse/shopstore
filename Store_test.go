package shopstore

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/gouniverse/sb"
	_ "modernc.org/sqlite"
)

func initDB(filepath string) *sql.DB {
	os.Remove(filepath) // remove database
	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func TestStoreDiscountCreate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		DiscountTableName:  "shop_discount_create",
		OrderTableName:     "shop_order_create",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	discount := NewDiscount().
		SetStatus(DISCOUNT_STATUS_DRAFT).
		SetTitle("DISCOUNT_TITLE")

	err = store.DiscountCreate(discount)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreDiscountFindByID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		DiscountTableName:  "shop_discount_find_by_id",
		OrderTableName:     "shop_order_find_by_id",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	discount := NewDiscount().
		SetStatus(DISCOUNT_STATUS_DRAFT).
		SetTitle("DISCOUNT_TITLE").
		SetDescription("DISCOUNT_DESCRIPTION").
		SetType(DISCOUNT_TYPE_AMOUNT).
		SetAmount(19.99).
		SetStartsAt(`2022-01-01 00:00:00`).
		SetEndsAt(`2022-01-01 23:59:59`)

	err = store.DiscountCreate(discount)

	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	discountFound, errFind := store.DiscountFindByID(discount.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if discountFound == nil {
		t.Fatal("Exam MUST NOT be nil")
		return
	}

	if discountFound.Title() != "DISCOUNT_TITLE" {
		t.Fatal("Exam title MUST BE 'DISCOUNT_TITLE', found: ", discountFound.Title())
		return
	}

	if discountFound.Description() != "DISCOUNT_DESCRIPTION" {
		t.Fatal("Exam description MUST BE 'DISCOUNT_DESCRIPTION', found: ", discountFound.Description())
	}

	if discountFound.Status() != DISCOUNT_STATUS_DRAFT {
		t.Fatal("Exam status MUST BE 'draft', found: ", discountFound.Status())
		return
	}

	if discountFound.Type() != DISCOUNT_TYPE_AMOUNT {
		t.Fatal("Exam type MUST BE 'amount', found: ", discountFound.Type())
	}

	if discountFound.Type() != DISCOUNT_TYPE_AMOUNT {
		t.Fatal("Exam type MUST BE 'amount', found: ", discountFound.Type())
	}

	if discountFound.Amount() != 19.9900 {
		t.Fatal("Exam price MUST BE '19.9900', found: ", discountFound.Amount())
		return
	}

	if discountFound.StartsAt() != "2022-01-01 00:00:00 +0000 UTC" {
		t.Fatal("Exam start date MUST BE '2022-01-01 00:00:00', found: ", discountFound.StartsAt())
	}

	if discountFound.EndsAt() != "2022-01-01 23:59:59 +0000 UTC" {
		t.Fatal("Exam end date MUST BE '2022-01-01 23:59:59', found: ", discountFound.EndsAt())
	}

	// if examFound.Memo() != "test memo" {
	// 	t.Fatal("Exam memo MUST BE 'test memo', found: ", examFound.Memo())
	// }

	if !strings.Contains(discountFound.DeletedAt(), sb.NULL_DATETIME) {
		t.Fatal("Exam MUST NOT be soft deleted", discountFound.DeletedAt())
		return
	}
}

func TestStoreOderCreate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		DiscountTableName:  "shop_discount_create",
		OrderTableName:     "shop_order_create",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetUserID("USER01_ID").
		SetQuantity(1).
		SetPrice(19.99)

	order.SetMetas(map[string]string{
		"color": "green",
		"size":  "xxl",
	})

	err = store.OrderCreate(order)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreOrderFindByID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		DiscountTableName:  "shop_discount_create",
		OrderTableName:     "shop_order_create",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetUserID("USER01_ID").
		SetQuantity(1).
		SetPrice(19.99).
		SetMemo("test memo")

	order.SetMetas(map[string]string{
		"color": "green",
		"size":  "xxl",
	})

	err = store.OrderCreate(order)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	orderFound, errFind := store.OrderFindByID(order.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if orderFound == nil {
		t.Fatal("Order MUST NOT be nil")
	}

	if orderFound.UserID() != "USER01_ID" {
		t.Fatal("Order user id MUST BE 'USER01_ID', found: ", orderFound.UserID())
	}

	if orderFound.Status() != ORDER_STATUS_PENDING {
		t.Fatal("Order status MUST BE 'pending', found: ", orderFound.Status())
	}

	if orderFound.Quantity() != "1" {
		t.Fatal("Order quantity MUST BE '1', found: ", orderFound.Quantity())
	}

	if orderFound.Price() != "19.9900" {
		t.Fatal("Order price MUST BE '19.9900', found: ", orderFound.Price())
	}

	if orderFound.Memo() != "test memo" {
		t.Fatal("Order memo MUST BE 'test memo', found: ", orderFound.Memo())
	}

	if orderFound.Meta("color") != "green" {
		t.Fatal("Order color meta MUST BE 'green', found: ", orderFound.Meta("color"))
	}

	if orderFound.Meta("size") != "xxl" {
		t.Fatal("Order size meta MUST BE 'xxl', found: ", orderFound.Meta("xxl"))
	}

	if !strings.Contains(orderFound.DeletedAt(), sb.NULL_DATETIME) {
		t.Fatal("Order MUST NOT be soft deleted", orderFound.DeletedAt())
	}
}

func TestStoreOrderSoftDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		DiscountTableName:  "shop_discount_create",
		OrderTableName:     "shop_order_create",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetUserID("USER01_ID").
		SetQuantity(1).
		SetPrice(19.99)

	err = store.OrderCreate(order)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderSoftDeleteByID(order.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if order.DeletedAt() != sb.NULL_DATETIME {
		t.Fatal("Order MUST NOT be soft deleted")
	}

	orderFound, errFind := store.OrderFindByID(order.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if orderFound != nil {
		t.Fatal("Order MUST be nil")
	}

	orderFindWithDeleted, errFind := store.OrderList(OrderQueryOptions{
		ID:          order.ID(),
		Limit:       1,
		WithDeleted: true,
	})

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if len(orderFindWithDeleted) < 1 {
		t.Fatal("Order list MUST NOT be empty")
		return
	}

	if strings.Contains(orderFindWithDeleted[0].DeletedAt(), sb.NULL_DATETIME) {
		t.Fatal("Order MUST be soft deleted", orderFound.DeletedAt())
	}

}
