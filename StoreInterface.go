package shopstore

type StoreInterface interface {
	OrderCreate(order OrderInterface) error
	// OrderDelete(order OrderInterface) error
	// OrderDeleteByID(id string) error
	OrderFindByID(id string) (OrderInterface, error)
	OrderList(options OrderQueryOptions) ([]OrderInterface, error)
	OrderSoftDelete(order OrderInterface) error
	OrderSoftDeleteByID(id string) error
	DiscountCreate(block *Discount) error
	DiscountDelete(block *Discount) error
	DiscountDeleteByID(blockID string) error
	DiscountFindByID(id string) (*Discount, error)
	DiscountFindByCode(code string) (*Discount, error)
	DiscountList(options DiscountQueryOptions) ([]Discount, error)
	DiscountSoftDelete(discount *Discount) error
	DiscountSoftDeleteByID(discountID string) error
	DiscountUpdate(block *Discount) error
}
