package shopstore

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/base/database"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
)

var _ StoreInterface = (*Store)(nil) // verify it extends the interface

type Store struct {
	categoryTableName      string
	discountTableName      string
	mediaTableName         string
	orderTableName         string
	orderLineItemTableName string
	productTableName       string
	db                     *sql.DB
	dbDriverName           string
	timeoutSeconds         int64
	automigrateEnabled     bool
	debugEnabled           bool
	sqlLogger              *slog.Logger
}

// logSql logs sql to the sql logger
func (store *Store) logSql(sqlOperationType string, sql string, params ...interface{}) {
	if !store.debugEnabled {
		return
	}

	if store.sqlLogger != nil {
		store.sqlLogger.Debug("sql: "+sqlOperationType, slog.String("sql", sql), slog.Any("params", params))
	}
}

// AutoMigrate auto migrate
func (store *Store) AutoMigrate() error {
	sqls := []string{
		store.sqlCategoryTableCreate(),
		store.sqlDiscountTableCreate(),
		store.sqlMediaTableCreate(),
		store.sqlOrderTableCreate(),
		store.sqlOrderLineItemTableCreate(),
		store.sqlProductTableCreate(),
	}

	for _, sql := range sqls {
		_, err := store.db.Exec(sql)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (store *Store) DB() *sql.DB {
	return store.db
}

// EnableDebug - enables the debug option
func (store *Store) EnableDebug(debug bool, sqlLogger ...*slog.Logger) {
	store.debugEnabled = debug
	if store.debugEnabled {
		if len(sqlLogger) < 1 {
			store.sqlLogger = slog.Default()
			return
		}
		store.sqlLogger = sqlLogger[0]
	} else {
		store.sqlLogger = nil
	}
}

func (store *Store) CategoryCount(ctx context.Context, options CategoryQueryInterface) (int64, error) {
	options.SetCountOnly(true)

	q, _, err := store.categoryQuery(options)

	if err != nil {
		return -1, err
	}

	sqlStr, params, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, nil
	}

	store.logSql("select", sqlStr, params...)

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

func (store *Store) CategoryCreate(ctx context.Context, category CategoryInterface) error {
	if category == nil {
		return errors.New("category is nil")
	}

	category.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	category.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	category.SetSoftDeletedAt(sb.MAX_DATETIME)

	data := category.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.categoryTableName).
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

	return nil
}

func (store *Store) CategoryDelete(ctx context.Context, category CategoryInterface) error {
	if category == nil {
		return errors.New("category is nil")
	}

	return store.CategoryDeleteByID(ctx, category.ID())
}

func (store *Store) CategoryDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.categoryTableName).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("delete", sqlStr, params...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return err
	}

	return nil
}

func (store *Store) CategoryFindByID(ctx context.Context, id string) (CategoryInterface, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}

	q := NewCategoryQuery().SetID(id).SetLimit(1)

	list, err := store.CategoryList(ctx, q)

	if err != nil {
		return nil, err
	}

	if len(list) < 1 {
		return nil, nil
	}

	return list[0], nil
}

func (store *Store) CategoryList(ctx context.Context, options CategoryQueryInterface) ([]CategoryInterface, error) {
	err := options.Validate()

	if err != nil {
		return nil, err
	}

	q, columns, err := store.categoryQuery(options)

	if err != nil {
		return nil, err
	}

	sqlStr, params, errSql := q.Prepared(true).
		Select(columns...).
		ToSQL()

	if errSql != nil {
		return nil, errSql
	}

	store.logSql("select", sqlStr, params...)

	modelMaps, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return nil, err
	}

	list := []CategoryInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewCategoryFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) CategorySoftDelete(ctx context.Context, category CategoryInterface) error {
	if category == nil {
		return errors.New("category is nil")
	}

	category.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.CategoryUpdate(ctx, category)
}

func (store *Store) CategorySoftDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is empty")
	}

	category, err := store.CategoryFindByID(ctx, id)

	if err != nil {
		return err
	}

	if category == nil {
		return nil
	}

	return store.CategorySoftDelete(ctx, category)
}

func (store *Store) CategoryUpdate(ctx context.Context, category CategoryInterface) (err error) {
	if category == nil {
		return errors.New("category is nil")
	}

	category.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := category.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable
	delete(dataChanged, "hash")    // Hash is not updateable
	delete(dataChanged, "data")    // Data is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.categoryTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C("id").Eq(category.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err = database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return err
	}

	category.MarkAsNotDirty()

	return nil

}

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

func (store *Store) categoryQuery(options CategoryQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, nil, errors.New("category options is nil")
	}

	if err := options.Validate(); err != nil {
		return nil, nil, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.categoryTableName)

	if options.HasID() {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if options.HasIDIn() {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn()))
	}

	if options.HasParentID() {
		q = q.Where(goqu.C(COLUMN_PARENT_ID).Eq(options.ParentID()))
	}

	if options.HasStatus() {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status()))
	}

	if options.HasTitleLike() {
		q = q.Where(goqu.C(COLUMN_TITLE).ILike(options.TitleLike()))
	}

	columns = []any{}

	for _, column := range options.Columns() {
		columns = append(columns, column)
	}

	if options.SoftDeletedIncluded() {
		return q, columns, nil // soft deleted blocks requested specifically
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), columns, nil
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

func (store *Store) MediaCount(ctx context.Context, options MediaQueryInterface) (int64, error) {
	options.SetCountOnly(true)

	q, _, err := store.mediaQuery(options)

	if err != nil {
		return -1, err
	}

	sqlStr, params, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, nil
	}

	store.logSql("select", sqlStr, params...)

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

func (store *Store) MediaCreate(ctx context.Context, media MediaInterface) error {
	if media == nil {
		return errors.New("media is nil")
	}

	media.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	media.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	media.SetSoftDeletedAt(sb.MAX_DATETIME)

	data := media.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.mediaTableName).
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

	return nil
}

func (store *Store) MediaDelete(ctx context.Context, media MediaInterface) error {
	if media == nil {
		return errors.New("media is nil")
	}

	return store.MediaDeleteByID(ctx, media.ID())
}

func (store *Store) MediaDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.mediaTableName).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("delete", sqlStr, params...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return err
	}

	return nil
}

func (store *Store) MediaFindByID(ctx context.Context, id string) (MediaInterface, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}

	q := NewMediaQuery().SetID(id).SetLimit(1)

	list, err := store.MediaList(ctx, q)

	if err != nil {
		return nil, err
	}

	if len(list) < 1 {
		return nil, nil
	}

	return list[0], nil
}

func (store *Store) MediaList(ctx context.Context, options MediaQueryInterface) ([]MediaInterface, error) {
	err := options.Validate()

	if err != nil {
		return nil, err
	}

	q, columns, err := store.mediaQuery(options)

	if err != nil {
		return nil, err
	}

	sqlStr, params, errSql := q.Prepared(true).
		Select(columns...).
		ToSQL()

	if errSql != nil {
		return nil, errSql
	}

	store.logSql("select", sqlStr, params...)

	modelMaps, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return nil, err
	}

	list := []MediaInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewMediaFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) MediaSoftDelete(ctx context.Context, media MediaInterface) error {
	if media == nil {
		return errors.New("media is nil")
	}

	media.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.MediaUpdate(ctx, media)
}

func (store *Store) MediaSoftDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is empty")
	}

	media, err := store.MediaFindByID(ctx, id)

	if err != nil {
		return err
	}

	if media == nil {
		return nil
	}

	return store.MediaSoftDelete(ctx, media)
}

func (store *Store) MediaUpdate(ctx context.Context, media MediaInterface) (err error) {
	if media == nil {
		return errors.New("media is nil")
	}

	media.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := media.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable
	delete(dataChanged, "hash")    // Hash is not updateable
	delete(dataChanged, "data")    // Data is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.mediaTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(media.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("update", sqlStr, params...)

	_, err = database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return err
	}

	media.MarkAsNotDirty()

	return nil
}

func (store *Store) mediaQuery(options MediaQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, nil, errors.New("category options is nil")
	}

	if err := options.Validate(); err != nil {
		return nil, nil, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.mediaTableName)

	if options.HasID() {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if options.HasIDIn() {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn()))
	}

	if options.HasEntityID() {
		q = q.Where(goqu.C(COLUMN_ENTITY_ID).Eq(options.EntityID()))
	}

	if options.HasStatus() {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status()))
	}

	if options.HasTitleLike() {
		q = q.Where(goqu.C(COLUMN_TITLE).ILike(options.TitleLike()))
	}

	if options.HasType() {
		q = q.Where(goqu.C(COLUMN_TYPE).Eq(options.Type()))
	}

	columns = []any{}

	for _, column := range options.Columns() {
		columns = append(columns, column)
	}

	if options.SoftDeletedIncluded() {
		return q, columns, nil // soft deleted blocks requested specifically
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), columns, nil
}

func (store *Store) OrderCount(ctx context.Context, options OrderQueryOptions) (int64, error) {
	options.CountOnly = true
	q := store.orderQuery(options)

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

func (store *Store) OrderCreate(ctx context.Context, order OrderInterface) error {
	order.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	order.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	order.SetDeletedAt(sb.MAX_DATETIME)

	data := order.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.orderTableName).
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

	order.MarkAsNotDirty()

	return nil
}

func (store *Store) OrderDelete(ctx context.Context, order OrderInterface) error {
	if order == nil {
		return errors.New("order is nil")
	}

	return store.OrderDeleteByID(ctx, order.ID())
}

func (store *Store) OrderDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("order id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.orderTableName).
		Prepared(true).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("delete", sqlStr, params...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	return err
}

func (store *Store) OrderSoftDelete(ctx context.Context, order OrderInterface) error {
	if order == nil {
		return errors.New("order is empty")
	}

	order.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.OrderUpdate(ctx, order)
}

func (store *Store) OrderSoftDeleteByID(ctx context.Context, id string) error {
	order, err := store.OrderFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.OrderSoftDelete(ctx, order)
}

func (store *Store) OrderFindByID(ctx context.Context, id string) (OrderInterface, error) {
	if id == "" {
		return nil, errors.New("order id is empty")
	}

	list, err := store.OrderList(ctx, OrderQueryOptions{
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

func (store *Store) OrderList(ctx context.Context, options OrderQueryOptions) ([]OrderInterface, error) {
	q := store.orderQuery(options)

	sqlStr, _, errSql := q.Select().ToSQL()

	if errSql != nil {
		return []OrderInterface{}, nil
	}

	store.logSql("select", sqlStr)

	modelMaps, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr)

	if err != nil {
		return []OrderInterface{}, err
	}

	list := []OrderInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewOrderFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) OrderUpdate(ctx context.Context, order OrderInterface) error {
	if order == nil {
		return errors.New("order is nil")
	}

	order.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := order.DataChanged()

	delete(dataChanged, "id")   // ID is not updateable
	delete(dataChanged, "hash") // Hash is not updateable
	delete(dataChanged, "data") // Data is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.orderTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C("id").Eq(order.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("update", sqlStr, params...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	order.MarkAsNotDirty()

	return err
}

func (store *Store) orderQuery(options OrderQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).From(store.orderTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if len(options.IDIn) > 0 {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn))
	}

	if options.CustomerID != "" {
		q = q.Where(goqu.C(COLUMN_CUSTOMER_ID).Eq(options.CustomerID))
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

func (store *Store) OrderLineItemCount(ctx context.Context, options OrderLineItemQueryOptions) (int64, error) {
	options.CountOnly = true
	q := store.orderLineItemQuery(options)

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

func (store *Store) OrderLineItemCreate(ctx context.Context, orderLineItem OrderLineItemInterface) error {
	if orderLineItem == nil {
		return errors.New("orderLineItem is nil")
	}

	orderLineItem.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	orderLineItem.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	orderLineItem.SetDeletedAt(sb.MAX_DATETIME)

	data := orderLineItem.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.orderLineItemTableName).
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

	orderLineItem.MarkAsNotDirty()

	return nil
}

func (store *Store) OrderLineItemDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("order line id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.orderLineItemTableName).
		Prepared(true).
		Where(goqu.C("id").Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("delete", sqlStr, params...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	return err
}

func (store *Store) OrderLineItemDelete(ctx context.Context, orderLineItem OrderLineItemInterface) error {
	return store.OrderLineItemDeleteByID(ctx, orderLineItem.ID())
}

func (store *Store) OrderLineItemFindByID(ctx context.Context, id string) (OrderLineItemInterface, error) {
	if id == "" {
		return nil, errors.New("order line id is empty")
	}

	list, err := store.OrderLineItemList(ctx, OrderLineItemQueryOptions{
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

func (store *Store) OrderLineItemList(ctx context.Context, options OrderLineItemQueryOptions) ([]OrderLineItemInterface, error) {
	q := store.orderLineItemQuery(options)

	sqlStr, params, errSql := q.Prepared(true).ToSQL()

	if errSql != nil {
		return nil, errSql
	}

	store.logSql("select", sqlStr, params...)

	modelMaps, err := database.SelectToMapString(store.toQuerableContext(ctx), sqlStr, params...)

	if err != nil {
		return []OrderLineItemInterface{}, err
	}

	list := []OrderLineItemInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewOrderLineItemFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *Store) OrderLineItemSoftDelete(ctx context.Context, orderLineItem OrderLineItemInterface) error {
	if orderLineItem == nil {
		return errors.New("order line is empty")
	}

	orderLineItem.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.OrderLineItemUpdate(ctx, orderLineItem)
}

func (store *Store) OrderLineItemSoftDeleteByID(ctx context.Context, id string) error {
	item, err := store.OrderLineItemFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.OrderLineItemSoftDelete(ctx, item)
}

func (store *Store) OrderLineItemUpdate(ctx context.Context, orderLineItem OrderLineItemInterface) error {
	if orderLineItem == nil {
		return errors.New("orderLineItem is nil")
	}

	orderLineItem.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	dataChanged := orderLineItem.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable
	delete(dataChanged, "hash")    // Hash is not updateable
	delete(dataChanged, "data")    // Data is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.orderLineItemTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C(COLUMN_ID).Eq(orderLineItem.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("update", sqlStr, params...)

	_, err := database.Execute(store.toQuerableContext(ctx), sqlStr, params...)

	orderLineItem.MarkAsNotDirty()

	return err
}

func (store *Store) orderLineItemQuery(options OrderLineItemQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).From(store.orderLineItemTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
	}

	if len(options.IDIn) > 0 {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn))
	}

	if options.OrderID != "" {
		q = q.Where(goqu.C(COLUMN_ORDER_ID).Eq(options.OrderID))
	}

	if options.ProductID != "" {
		q = q.Where(goqu.C(COLUMN_PRODUCT_ID).Eq(options.ProductID))
	}

	if options.Status != "" {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(options.Status))
	}

	if len(options.StatusIn) > 0 {
		q = q.Where(goqu.C(COLUMN_STATUS).In(options.StatusIn))
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

type DiscountQueryOptions struct {
	ID           string
	IDIn         []string
	Status       string
	StatusIn     []string
	Code         string
	CreatedAtGte string
	CreatedAtLte string
	Offset       int
	Limit        int
	SortOrder    string
	OrderBy      string
	CountOnly    bool
	WithDeleted  bool
}

type MediaQueryOptions struct {
	ID          string
	IDIn        []string
	EntityID    string
	Status      string
	StatusIn    []string
	Offset      int
	Limit       int
	SortOrder   string
	OrderBy     string
	CountOnly   bool
	WithDeleted bool
}

type OrderQueryOptions struct {
	ID           string
	IDIn         string
	CustomerID   string
	Status       string
	StatusIn     []string
	CreatedAtGte string
	CreatedAtLte string
	Offset       int
	Limit        int
	SortOrder    string
	OrderBy      string
	CountOnly    bool
	WithDeleted  bool
}

type OrderLineItemQueryOptions struct {
	ID          string
	IDIn        string
	OrderID     string
	ProductID   string
	Status      string
	StatusIn    []string
	Offset      int
	Limit       int
	SortOrder   string
	OrderBy     string
	CountOnly   bool
	WithDeleted bool
}

type ProductQueryOptions struct {
	ID           string
	IDIn         []string
	Status       string
	StatusIn     []string
	Title        string
	CreatedAtGte string
	CreatedAtLte string
	Offset       int
	Limit        int
	SortOrder    string
	OrderBy      string
	CountOnly    bool
	WithDeleted  bool
}

func (store *Store) toQuerableContext(context context.Context) database.QueryableContext {
	if database.IsQueryableContext(context) {
		return context.(database.QueryableContext)
	}

	return database.Context(context, store.db)
}
