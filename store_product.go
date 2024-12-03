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

func (store *Store) ProductCount(ctx context.Context, options ProductQueryOptions) (int64, error) {
	options.CountOnly = true
	q := store.productQuery(options)

	sqlStr, params, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, nil
	}

	store.logSql("count", sqlStr, params...)

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

func (store *Store) ProductCreate(ctx context.Context, product ProductInterface) error {
	if product == nil {
		return errors.New("product is nil")
	}

	product.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	product.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	product.SetDeletedAt(sb.MAX_DATETIME)

	data := product.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.productTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("insert", sqlStr, params...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return err
	}

	product.MarkAsNotDirty()

	return nil
}

func (store *Store) ProductDelete(ctx context.Context, product ProductInterface) error {
	if product == nil {
		return errors.New("product is nil")
	}

	return store.ProductDeleteByID(ctx, product.ID())
}

func (store *Store) ProductDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("product id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.productTableName).
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

func (store *Store) ProductSoftDelete(ctx context.Context, product ProductInterface) error {
	if product == nil {
		return errors.New("product is empty")
	}

	product.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.ProductUpdate(ctx, product)
}

func (store *Store) ProductSoftDeleteByID(ctx context.Context, id string) error {
	product, err := store.ProductFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.ProductSoftDelete(ctx, product)
}

func (store *Store) ProductFindByID(ctx context.Context, id string) (ProductInterface, error) {
	if id == "" {
		return nil, errors.New("product id is empty")
	}

	list, err := store.ProductList(ctx, ProductQueryOptions{
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

func (store *Store) ProductList(ctx context.Context, options ProductQueryOptions) ([]ProductInterface, error) {
	q := store.productQuery(options)

	sqlStr, _, errSql := q.Select().ToSQL()

	if errSql != nil {
		return []ProductInterface{}, nil
	}

	store.logSql("select", sqlStr)

	modelMaps, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr)
	if err != nil {
		return []ProductInterface{}, err
	}

	list := []ProductInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewProductFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) ProductUpdate(ctx context.Context, product ProductInterface) error {
	if product == nil {
		return errors.New("product is nil")
	}

	product.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := product.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable
	delete(dataChanged, "hash")    // Hash is not updateable
	delete(dataChanged, "data")    // Data is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.productTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(product.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("update", sqlStr, params...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	product.MarkAsNotDirty()

	return err
}

func (store *Store) productQuery(options ProductQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).From(store.productTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if len(options.IDIn) > 0 {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn))
	}

	if options.Title != "" {
		q = q.Where(goqu.C(COLUMN_TITLE).Eq(options.Title))
	}

	if options.Status != "" {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status))
	}

	if len(options.StatusIn) > 0 {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn))
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
