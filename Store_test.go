package shopstore

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/gouniverse/sb"
	_ "modernc.org/sqlite"
)

func initDB(filepath string) (*sql.DB, error) {
	err := os.Remove(filepath) // remove database

	if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
		return nil, err
	}

	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestStoreDiscountCreate(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_create",
		OrderTableName:         "shop_order_create",
		OrderLineItemTableName: "shop_order_line_item_create",
		ProductTableName:       "shop_product_create",
		AutomigrateEnabled:     true,
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

func TestStoreDiscountDelete(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_delete",
		OrderTableName:         "shop_order_delete",
		OrderLineItemTableName: "shop_order_line_item_delete",
		ProductTableName:       "shop_product_delete",
		AutomigrateEnabled:     true,
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
		return
	}

	err = store.DiscountDelete(discount)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	discountFound, errFind := store.DiscountFindByID(discount.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if discountFound != nil {
		t.Fatal("Exam MUST be nil")
		return
	}
}

func TestStoreDiscountDeleteByID(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_delete_by_id",
		OrderTableName:         "shop_order_delete_by_id",
		OrderLineItemTableName: "shop_order_line_item_delete_by_id",
		ProductTableName:       "shop_product_delete_by_id",
		AutomigrateEnabled:     true,
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
		return
	}

	err = store.DiscountDeleteByID(discount.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	discountFound, errFind := store.DiscountFindByID(discount.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if discountFound != nil {
		t.Fatal("Exam MUST be nil")
		return
	}
}

func TestStoreDiscountFindByID(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_find_by_id",
		OrderTableName:         "shop_order_find_by_id",
		OrderLineItemTableName: "shop_order_line_item_find_by_id",
		ProductTableName:       "shop_product_find_by_id",
		AutomigrateEnabled:     true,
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
		t.Fatal("Discount MUST NOT be nil")
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

	if !strings.Contains(discountFound.DeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Exam MUST NOT be soft deleted", discountFound.DeletedAt())
		return
	}
}

func TestStoreDiscountSoftDelete(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_soft_delete",
		OrderTableName:         "shop_order_soft_delete",
		OrderLineItemTableName: "shop_order_line_item_soft_delete",
		ProductTableName:       "shop_product_soft_delete",
		AutomigrateEnabled:     true,
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
		return
	}

	err = store.DiscountSoftDelete(discount)
	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	discountFound, errFind := store.DiscountFindByID(discount.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if discountFound != nil {
		t.Fatal("Discount MUST be nil")
		return
	}

	discountList, errList := store.DiscountList(DiscountQueryOptions{
		ID:          discount.ID(),
		WithDeleted: true,
	})

	if errList != nil {
		t.Fatal("unexpected error:", errList)
		return
	}

	if len(discountList) != 1 {
		t.Fatal("Discount list MUST be 1")
		return
	}
}

func TestStoreDiscountUpdate(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_update",
		OrderTableName:         "shop_order_update",
		OrderLineItemTableName: "shop_order_line_item_update",
		ProductTableName:       "shop_product_update",
		AutomigrateEnabled:     true,
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

	discount.SetTitle("DISCOUNT_TITLE_UPDATED")

	err = store.DiscountUpdate(discount)
	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	discountFound, errFind := store.DiscountFindByID(discount.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if discountFound == nil {
		t.Fatal("Discount MUST NOT be nil")
	}

	if discountFound.Title() != "DISCOUNT_TITLE_UPDATED" {
		t.Fatal("Discount title MUST BE 'DISCOUNT_TITLE_UPDATED', found: ", discountFound.Title())
	}
}

func TestStoreOderCreate(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_create",
		OrderTableName:         "shop_order_create",
		OrderLineItemTableName: "shop_order_line_item_create",
		ProductTableName:       "shop_product_create",
		AutomigrateEnabled:     true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetCustomerID("CUSTOMER01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	err = order.SetMetas(map[string]string{
		"color": "green",
		"size":  "xxl",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderCreate(order)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreOderDelete(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_delete",
		OrderTableName:         "shop_order_delete",
		OrderLineItemTableName: "shop_order_line_item_delete",
		ProductTableName:       "shop_product_delete",
		AutomigrateEnabled:     true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetCustomerID("CUSTOMER01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	err = store.OrderCreate(order)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderDelete(order)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	orderFound, err := store.OrderFindByID(order.ID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if orderFound != nil {
		t.Fatal("expected nil order")
	}
}

func TestStoreOderDeleteByID(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_delete",
		OrderTableName:         "shop_order_delete",
		OrderLineItemTableName: "shop_order_line_item_delete",
		ProductTableName:       "shop_product_delete",
		AutomigrateEnabled:     true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetCustomerID("CUSTOMER01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	err = store.OrderCreate(order)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderDeleteByID(order.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	orderFound, err := store.OrderFindByID(order.ID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if orderFound != nil {
		t.Fatal("expected nil order")
	}
}

func TestStoreOrderFindByID(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_find_by_id",
		OrderTableName:         "shop_order_find_by_id",
		OrderLineItemTableName: "shop_order_line_item_find_by_id",
		ProductTableName:       "shop_product_find_by_id",
		AutomigrateEnabled:     true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetCustomerID("CUSTOMER01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99).
		SetMemo("test memo")

	err = order.SetMetas(map[string]string{
		"color": "green",
		"size":  "xxl",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

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

	if orderFound.CustomerID() != "CUSTOMER01_ID" {
		t.Fatal("Order user id MUST BE 'CUSTOMER01_ID', found: ", orderFound.CustomerID())
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

	if !strings.Contains(orderFound.DeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Order MUST NOT be soft deleted", orderFound.DeletedAt())
	}
}

func TestStoreOrderSoftDelete(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_soft_delete",
		OrderTableName:         "shop_order_soft_delete",
		OrderLineItemTableName: "shop_order_line_item_soft_delete",
		ProductTableName:       "shop_product_soft_delete",
		AutomigrateEnabled:     true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetCustomerID("USER01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	err = store.OrderCreate(order)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderSoftDeleteByID(order.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if order.DeletedAt() != sb.MAX_DATETIME {
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

func TestStoreOrderUpdate(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_update",
		OrderTableName:         "shop_order_update",
		OrderLineItemTableName: "shop_order_line_item_update",
		ProductTableName:       "shop_product_update",
		AutomigrateEnabled:     true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetCustomerID("CUSTOMER01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	err = store.OrderCreate(order)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	order.SetMemo("test memo")

	err = store.OrderUpdate(order)
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

	if orderFound.Memo() != "test memo" {
		t.Fatal("Order memo MUST BE 'test memo', found: ", orderFound.Memo())
	}
}

func TestStoreOrderLineItemCreate(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_create",
		OrderTableName:         "shop_order_create",
		OrderLineItemTableName: "shop_order_line_item_create",
		ProductTableName:       "shop_product_create",
		AutomigrateEnabled:     true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	orderLineItem := NewOrderLineItem().
		SetOrderID("ORDER01_ID").
		SetProductID("PRODUCT01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	err = store.OrderLineItemCreate(orderLineItem)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreOrderLineItemDelete(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_delete",
		OrderTableName:         "shop_order_delete",
		OrderLineItemTableName: "shop_order_line_item_delete",
		ProductTableName:       "shop_product_delete",
		AutomigrateEnabled:     true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	orderLineItem := NewOrderLineItem().
		SetOrderID("ORDER01_ID").
		SetProductID("PRODUCT01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	err = store.OrderLineItemCreate(orderLineItem)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderLineItemDelete(orderLineItem)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	orderLineItemFound, errFind := store.OrderLineItemFindByID(orderLineItem.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if orderLineItemFound != nil {
		t.Fatal("OrderLineItem MUST be nil")
	}
}

func TestStoreOrderLineItemFindByID(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_create",
		OrderTableName:         "shop_order_create",
		OrderLineItemTableName: "shop_order_line_item_create",
		ProductTableName:       "shop_product_create",
		AutomigrateEnabled:     true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	orderLineItem := NewOrderLineItem().
		SetOrderID("ORDER01_ID").
		SetProductID("PRODUCT01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	err = store.OrderLineItemCreate(orderLineItem)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	orderLineItemFound, errFind := store.OrderLineItemFindByID(orderLineItem.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if orderLineItemFound == nil {
		t.Fatal("OrderLineItem MUST NOT be nil")
	}
}

func TestStoreOrderLineItemList(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_create",
		OrderTableName:         "shop_order_create",
		OrderLineItemTableName: "shop_order_line_item_create",
		ProductTableName:       "shop_product_create",
		AutomigrateEnabled:     true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	orderLineItem := NewOrderLineItem().
		SetOrderID("ORDER01_ID").
		SetProductID("PRODUCT01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	err = store.OrderLineItemCreate(orderLineItem)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	orderLineItemsFound, errFind := store.OrderLineItemList(OrderLineItemQueryOptions{})

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if len(orderLineItemsFound) != 1 {
		t.Fatal("OrderLineItem MUST NOT be nil")
	}
}

func TestStoreOrderLineItemSoftDeleteByID(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_soft_delete",
		OrderTableName:         "shop_order_soft_delete",
		OrderLineItemTableName: "shop_order_line_item_soft_delete",
		ProductTableName:       "shop_product_soft_delete",
		AutomigrateEnabled:     true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	orderLineItem := NewOrderLineItem().
		SetOrderID("ORDER01_ID").
		SetProductID("PRODUCT01_ID").
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	err = store.OrderLineItemCreate(orderLineItem)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderLineItemSoftDeleteByID(orderLineItem.ID())
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	orderLineItemFound, errFind := store.OrderLineItemFindByID(orderLineItem.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if orderLineItemFound != nil {
		t.Fatal("OrderLineItem MUST be nil")
	}

	orderLineItems, errFind := store.OrderLineItemList(OrderLineItemQueryOptions{
		OrderID:     orderLineItem.OrderID(),
		WithDeleted: true,
	})

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if len(orderLineItems) != 1 {
		t.Fatal("OrderLineItem MUST be deleted")
	}
}

func TestStoreProductCreate(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_create",
		OrderTableName:         "shop_order_create",
		OrderLineItemTableName: "shop_order_line_item_create",
		ProductTableName:       "shop_product_create",
		AutomigrateEnabled:     true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	product := NewProduct().
		SetStatus(PRODUCT_STATUS_DRAFT).
		SetQuantityInt(1).
		SetPriceFloat(19.99)

	err = product.SetMetas(map[string]string{
		"color": "green",
		"size":  "xxl",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.ProductCreate(product)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreProductFindByID(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_find_by_id",
		OrderTableName:         "shop_order_find_by_id",
		OrderLineItemTableName: "shop_order_line_item_find_by_id",
		ProductTableName:       "shop_product_find_by_id",
		AutomigrateEnabled:     true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	product := NewProduct().
		SetStatus(PRODUCT_STATUS_DRAFT).
		SetTitle("Ruler").
		SetQuantityInt(1).
		SetPriceFloat(19.99).
		SetMemo("test ruler")

	err = product.SetMetas(map[string]string{
		"color": "green",
		"size":  "xxl",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.ProductCreate(product)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	productFound, errFind := store.ProductFindByID(product.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if productFound == nil {
		t.Fatal("Product MUST NOT be nil")
	}

	if productFound.Title() != "Ruler" {
		t.Fatal("Product title MUST BE 'Ruler', found: ", productFound.Title())
	}

	if productFound.Status() != PRODUCT_STATUS_DRAFT {
		t.Fatal("Product status MUST BE 'draft', found: ", productFound.Status())
	}

	if productFound.Quantity() != "1" {
		t.Fatal("Product quantity MUST BE '1', found: ", productFound.Quantity())
	}

	if productFound.Price() != "19.9900" {
		t.Fatal("Product price MUST BE '19.9900', found: ", productFound.Price())
	}

	if productFound.Memo() != "test ruler" {
		t.Fatal("Product memo MUST BE 'test ruler', found: ", productFound.Memo())
	}

	if productFound.Meta("color") != "green" {
		t.Fatal("Product color meta MUST BE 'green', found: ", productFound.Meta("color"))
	}

	if productFound.Meta("size") != "xxl" {
		t.Fatal("Product size meta MUST BE 'xxl', found: ", productFound.Meta("xxl"))
	}

	if !strings.Contains(productFound.DeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Product MUST NOT be soft deleted", productFound.DeletedAt())
	}
}

func TestStoreProductSoftDelete(t *testing.T) {
	db, err := initDB(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	store, err := NewStore(NewStoreOptions{
		DB:                     db,
		DiscountTableName:      "shop_discount_soft_delete",
		OrderTableName:         "shop_order_soft_delete",
		OrderLineItemTableName: "shop_order_line_item_soft_delete",
		ProductTableName:       "shop_product_soft_delete",
		AutomigrateEnabled:     true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	product := NewProduct().
		SetStatus(PRODUCT_STATUS_DRAFT).
		SetTitle("Ruler").
		SetQuantityInt(1).
		SetPriceFloat(19.99).
		SetMemo("test ruler")

	err = store.ProductCreate(product)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if product.DeletedAt() != sb.MAX_DATETIME {
		t.Fatal("Product MUST NOT be soft deleted")
	}

	err = store.ProductSoftDelete(product)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	productFound, errFind := store.ProductFindByID(product.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if productFound != nil {
		t.Fatal("Product MUST be nil")
	}

	productFindWithDeleted, errFind := store.ProductList(ProductQueryOptions{
		ID:          product.ID(),
		Limit:       1,
		WithDeleted: true,
	})

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if len(productFindWithDeleted) < 1 {
		t.Fatal("Product list MUST NOT be empty")
		return
	}

	if strings.Contains(productFindWithDeleted[0].DeletedAt(), sb.NULL_DATETIME) {
		t.Fatal("Product MUST be soft deleted", productFound.DeletedAt())
	}

}
