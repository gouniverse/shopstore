package shopstore

type DiscountInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	Amount() float64
	SetAmount(amount float64) DiscountInterface
	Code() string
	SetCode(code string) DiscountInterface
	CreatedAt() string
	SetCreatedAt(createdAt string) DiscountInterface
	DeletedAt() string
	SetDeletedAt(deletedAt string) DiscountInterface
	Description() string
	SetDescription(description string) DiscountInterface
	EndsAt() string
	SetEndsAt(endsAt string) DiscountInterface
	ID() string
	SetID(id string) DiscountInterface
	StartsAt() string
	SetStartsAt(startsAt string) DiscountInterface
	Status() string
	SetStatus(status string) DiscountInterface
	Title() string
	SetTitle(title string) DiscountInterface
	Type() string
	SetType(type_ string) DiscountInterface
	UpdatedAt() string
	SetUpdatedAt(updatedAt string) DiscountInterface
}
