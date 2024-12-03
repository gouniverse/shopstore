package shopstore

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/base/database"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
)

func (store *Store) DiscountCount(ctx context.Context, options DiscountQueryOptions) (int64, error) {
	options.CountOnly = true
	q := store.discountQuery(options)

	sqlStr, params, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	mapped, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return -1, err
	}

	if len(mapped) < 1 {
		return -1, nil
	}

	countStr := mapped[0]["count"]

	i, err := strconv.ParseInt(countStr, 10, 64)

	if err != nil {
		return -1, err

	}

	return i, nil
}

func (store *Store) DiscountCreate(ctx context.Context, discount DiscountInterface) error {
	discount.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	discount.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	discount.SetDeletedAt(sb.MAX_DATETIME)

	data := discount.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.discountTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return err
	}

	discount.MarkAsNotDirty()

	return nil
}

func (store *Store) DiscountDelete(ctx context.Context, discount DiscountInterface) error {
	if discount == nil {
		return errors.New("discount is nil")
	}

	return store.DiscountDeleteByID(ctx, discount.ID())
}

func (store *Store) DiscountDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("discount id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.discountTableName).
		Prepared(true).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	return err
}

func (store *Store) DiscountFindByID(ctx context.Context, id string) (DiscountInterface, error) {
	if id == "" {
		return nil, errors.New("discount id is empty")
	}

	list, err := store.DiscountList(ctx, DiscountQueryOptions{
		ID:    id,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *Store) DiscountFindByCode(ctx context.Context, code string) (DiscountInterface, error) {
	if code == "" {
		return nil, errors.New("discount code is empty")
	}

	list, err := store.DiscountList(ctx, DiscountQueryOptions{
		Status: DISCOUNT_STATUS_ACTIVE,
		Code:   code,
		Limit:  1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *Store) DiscountList(ctx context.Context, options DiscountQueryOptions) ([]DiscountInterface, error) {
	q := store.discountQuery(options)

	sqlStr, _, errSql := q.Select().ToSQL()

	if errSql != nil {
		return []DiscountInterface{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	modelMaps, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr)

	if err != nil {
		return []DiscountInterface{}, err
	}

	list := []DiscountInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewDiscountFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) DiscountSoftDelete(ctx context.Context, discount DiscountInterface) error {
	if discount == nil {
		return errors.New("discount is nil")
	}

	discount.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.DiscountUpdate(ctx, discount)
}

func (store *Store) DiscountSoftDeleteByID(ctx context.Context, id string) error {
	discount, err := store.DiscountFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.DiscountSoftDelete(ctx, discount)
}

func (store *Store) DiscountUpdate(ctx context.Context, discount DiscountInterface) error {
	if discount == nil {
		return errors.New("discount is nil")
	}

	discount.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := discount.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable
	delete(dataChanged, "hash")    // Hash is not updateable
	delete(dataChanged, "data")    // Data is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.discountTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C("id").Eq(discount.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	discount.MarkAsNotDirty()

	return err
}

func (store *Store) discountQuery(options DiscountQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).From(store.discountTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if len(options.IDIn) > 0 {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn))
	}

	if options.Status != "" {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status))
	}

	if len(options.StatusIn) > 0 {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn))
	}

	if options.Code != "" {
		q = q.Where(goqu.C(COLUMN_CODE).Eq(options.Code))
	}

	if options.CreatedAtGte != "" && options.CreatedAtLte != "" {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Between(exp.NewRangeVal(options.CreatedAtGte, options.CreatedAtLte)))
	} else if options.CreatedAtGte != "" {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Gte(options.CreatedAtGte))
	} else if options.CreatedAtLte != "" {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Lte(options.CreatedAtLte))
	}

	if !options.CountOnly {
		if options.Limit > 0 {
			q = q.Limit(uint(options.Limit))
		}

		if options.Offset > 0 {
			q = q.Offset(uint(options.Offset))
		}
	}

	sortOrder := sb.DESC
	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.OrderBy != "" {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy).Desc())
		}
	}

	if options.WithDeleted {
		return q
	}

	softDeleted := goqu.C(COLUMN_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted)
}
