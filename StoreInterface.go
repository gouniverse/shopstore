package shopstore

type StoreInterface interface {
	AutoMigrate() error
	EnableDebug(debug bool)

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
