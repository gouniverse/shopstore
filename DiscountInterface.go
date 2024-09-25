package shopstore

import "github.com/golang-module/carbon/v2"

type DiscountInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	Amount() float64
	SetAmount(amount float64) DiscountInterface
	Code() string
	SetCode(code string) DiscountInterface
	CreatedAt() string
	CreatedAtCarbon() carbon.Carbon
	SetCreatedAt(createdAt string) DiscountInterface
	DeletedAt() string
	DeletedAtCarbon() carbon.Carbon
	SetDeletedAt(deletedAt string) DiscountInterface
	Description() string
	SetDescription(description string) DiscountInterface
	EndsAt() string
	EndsAtCarbon() carbon.Carbon
	SetEndsAt(endsAt string) DiscountInterface
	ID() string
	SetID(id string) DiscountInterface
	StartsAt() string
	StartsAtCarbon() carbon.Carbon
	SetStartsAt(startsAt string) DiscountInterface
	Status() string
	SetStatus(status string) DiscountInterface
	Title() string
	SetTitle(title string) DiscountInterface
	Type() string
	SetType(type_ string) DiscountInterface
	UpdatedAt() string
	UpdatedAtCarbon() carbon.Carbon
	SetUpdatedAt(updatedAt string) DiscountInterface
}
