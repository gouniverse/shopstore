package shopstore

import (
	"github.com/golang-module/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
	"github.com/gouniverse/utils"
)

// == CONSTANTS ==============================================================

const DISCOUNT_STATUS_DRAFT = "draft"
const DISCOUNT_STATUS_ACTIVE = "active"
const DISCOUNT_STATUS_INACTIVE = "inactive"

const DISCOUNT_TYPE_AMOUNT = "amount"
const DISCOUNT_TYPE_PERCENT = "percent"

const DISCOUNT_DURATION_FOREVER = "forever"
const DISCOUNT_DURATION_MONTHS = "months"
const DISCOUNT_DURATION_ONCE = "once"

// == CLASS ==================================================================

type Discount struct {
	dataobject.DataObject
}

// == INTERFACES =============================================================

var _ DiscountInterface = (*Discount)(nil)

// == CONSTRUCTORS ===========================================================

func NewDiscount() DiscountInterface {
	code := uid.Timestamp()

	d := (&Discount{}).
		SetID(uid.HumanUid()).
		SetStatus(DISCOUNT_STATUS_DRAFT).
		SetType(DISCOUNT_TYPE_PERCENT).
		SetTitle("").
		SetDescription("").
		SetAmount(0.00).
		SetCode(code).
		SetStartsAt(sb.NULL_DATETIME).
		SetEndsAt(sb.NULL_DATETIME).
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetDeletedAt(sb.MAX_DATETIME)

	return d
}

func NewDiscountFromExistingData(data map[string]string) DiscountInterface {
	o := &Discount{}
	o.Hydrate(data)
	return o
}

// == METHODS ================================================================

// == SETTERS AND GETTERS ====================================================

func (d *Discount) Amount() float64 {
	amountStr := d.Get(COLUMN_AMOUNT)
	amount, err := utils.ToFloat(amountStr)

	if err != nil {
		return 0
	}

	return amount
}

func (d *Discount) SetAmount(amount float64) DiscountInterface {
	amountStr := utils.ToString(amount)
	d.Set(COLUMN_AMOUNT, amountStr)
	return d
}

func (d *Discount) Code() string {
	return d.Get(COLUMN_CODE)
}

func (d *Discount) SetCode(code string) DiscountInterface {
	d.Set(COLUMN_CODE, code)
	return d
}

func (d *Discount) CreatedAt() string {
	return d.Get(COLUMN_CREATED_AT)
}

func (d *Discount) CreatedAtCarbon() carbon.Carbon {
	createdAt := d.CreatedAt()
	return carbon.Parse(createdAt)
}

func (d *Discount) SetCreatedAt(createdAt string) DiscountInterface {
	d.Set(COLUMN_CREATED_AT, createdAt)
	return d
}

func (d *Discount) DeletedAt() string {
	return d.Get(COLUMN_DELETED_AT)
}

func (d *Discount) SetDeletedAt(deletedAt string) DiscountInterface {
	d.Set(COLUMN_DELETED_AT, deletedAt)
	return d
}

func (d *Discount) Description() string {
	return d.Get(COLUMN_DESCRIPTION)
}

func (d *Discount) SetDescription(description string) DiscountInterface {
	d.Set(COLUMN_DESCRIPTION, description)
	return d
}

func (d *Discount) EndsAt() string {
	return d.Get(COLUMN_ENDS_AT)
}

func (d *Discount) EndsAtCarbon() carbon.Carbon {
	endsAt := d.EndsAt()
	return carbon.Parse(endsAt)
}

func (d *Discount) SetEndsAt(endsAt string) DiscountInterface {
	d.Set(COLUMN_ENDS_AT, endsAt)
	return d
}

// ID returns the ID of the exam
func (o *Discount) ID() string {
	return o.Get(COLUMN_ID)
}

// SetID sets the ID of the exam
func (o *Discount) SetID(id string) DiscountInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (d *Discount) StartsAt() string {
	return d.Get(COLUMN_STARTS_AT)
}

func (d *Discount) StartsAtCarbon() carbon.Carbon {
	startsAt := d.StartsAt()
	return carbon.Parse(startsAt)
}

func (d *Discount) SetStartsAt(startsAt string) DiscountInterface {
	d.Set(COLUMN_STARTS_AT, startsAt)
	return d
}

func (d *Discount) Status() string {
	return d.Get(COLUMN_STATUS)
}

func (d *Discount) SetStatus(status string) DiscountInterface {
	d.Set(COLUMN_STATUS, status)
	return d
}

func (d *Discount) Title() string {
	return d.Get(COLUMN_TITLE)
}

func (d *Discount) SetTitle(title string) DiscountInterface {
	d.Set(COLUMN_TITLE, title)
	return d
}

func (d *Discount) Type() string {
	return d.Get(COLUMN_TYPE)
}

func (d *Discount) SetType(type_ string) DiscountInterface {
	d.Set(COLUMN_TYPE, type_)
	return d
}

func (d *Discount) UpdatedAt() string {
	return d.Get(COLUMN_UPDATED_AT)
}

func (d *Discount) UpdatedAtCarbon() carbon.Carbon {
	updatedAt := d.UpdatedAt()
	return carbon.Parse(updatedAt)
}

func (d *Discount) SetUpdatedAt(updatedAt string) DiscountInterface {
	d.Set(COLUMN_UPDATED_AT, updatedAt)
	return d
}
