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

// == CONSTRUCTORS ===========================================================

func NewDiscount() *Discount {
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
		SetDeletedAt(sb.NULL_DATETIME)

	return d
}

func NewDiscountFromExistingData(data map[string]string) *Discount {
	o := &Discount{}
	o.Hydrate(data)
	return o
}

// == METHODS ================================================================

// == SETTERS AND GETTERS ====================================================

func (d *Discount) Amount() float64 {
	amountStr := d.Get("amount")
	amount, err := utils.ToFloat(amountStr)

	if err != nil {
		return 0
	}

	return amount
}

func (d *Discount) SetAmount(amount float64) *Discount {
	amountStr := utils.ToString(amount)
	d.Set("amount", amountStr)
	return d
}

func (d *Discount) Code() string {
	return d.Get("code")
}

func (d *Discount) SetCode(code string) *Discount {
	d.Set("code", code)
	return d
}

func (d *Discount) CreatedAt() string {
	return d.Get("created_at")
}

func (d *Discount) CreatedAtCarbon() carbon.Carbon {
	createdAt := d.CreatedAt()
	return carbon.Parse(createdAt)
}

func (d *Discount) SetCreatedAt(createdAt string) *Discount {
	d.Set("created_at", createdAt)
	return d
}

func (d *Discount) DeletedAt() string {
	return d.Get("deleted_at")
}

func (d *Discount) SetDeletedAt(deletedAt string) *Discount {
	d.Set("deleted_at", deletedAt)
	return d
}

func (d *Discount) Description() string {
	return d.Get("description")
}

func (d *Discount) SetDescription(description string) *Discount {
	d.Set("description", description)
	return d
}

func (d *Discount) EndsAt() string {
	return d.Get("ends_at")
}

func (d *Discount) EndsAtCarbon() carbon.Carbon {
	endsAt := d.EndsAt()
	return carbon.Parse(endsAt)
}

func (d *Discount) SetEndsAt(endsAt string) *Discount {
	d.Set("ends_at", endsAt)
	return d
}

// ID returns the ID of the exam
func (o *Discount) ID() string {
	return o.Get("id")
}

// SetID sets the ID of the exam
func (o *Discount) SetID(id string) *Discount {
	o.Set("id", id)
	return o
}

func (d *Discount) StartsAt() string {
	return d.Get("starts_at")
}

func (d *Discount) StartsAtCarbon() carbon.Carbon {
	startsAt := d.StartsAt()
	return carbon.Parse(startsAt)
}

func (d *Discount) SetStartsAt(startsAt string) *Discount {
	d.Set("starts_at", startsAt)
	return d
}

func (d *Discount) Status() string {
	return d.Get("status")
}

func (d *Discount) SetStatus(status string) *Discount {
	d.Set("status", status)
	return d
}

func (d *Discount) Title() string {
	return d.Get("title")
}

func (d *Discount) SetTitle(title string) *Discount {
	d.Set("title", title)
	return d
}

func (d *Discount) Type() string {
	return d.Get("type")
}

func (d *Discount) SetType(type_ string) *Discount {
	d.Set("type", type_)
	return d
}

func (d *Discount) UpdatedAt() string {
	return d.Get("updated_at")
}

func (d *Discount) UpdatedAtCarbon() carbon.Carbon {
	updatedAt := d.UpdatedAt()
	return carbon.Parse(updatedAt)
}

func (d *Discount) SetUpdatedAt(updatedAt string) *Discount {
	d.Set("updated_at", updatedAt)
	return d
}
