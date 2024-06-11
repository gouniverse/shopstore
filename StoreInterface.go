package shopstore

type StoreInterface interface {
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
