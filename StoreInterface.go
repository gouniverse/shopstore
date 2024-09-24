package shopstore

type StoreInterface interface {
	OrderCreate(order OrderInterface) error
	OrderDelete(order OrderInterface) error
	OrderDeleteByID(id string) error
	OrderFindByID(id string) (OrderInterface, error)
	OrderList(options OrderQueryOptions) ([]OrderInterface, error)
	OrderSoftDelete(order OrderInterface) error
	OrderSoftDeleteByID(id string) error

	DiscountCreate(discount DiscountInterface) error
	DiscountDelete(discount DiscountInterface) error
	DiscountDeleteByID(discountID string) error
	DiscountFindByID(discountID string) (DiscountInterface, error)
	DiscountFindByCode(code string) (DiscountInterface, error)
	DiscountList(options DiscountQueryOptions) ([]DiscountInterface, error)
	DiscountSoftDelete(discount DiscountInterface) error
	DiscountSoftDeleteByID(discountID string) error
	DiscountUpdate(discount DiscountInterface) error
}
