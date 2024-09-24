package shopstore

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/golang-module/carbon/v2"
	"github.com/gouniverse/sb"
	"github.com/mingrammer/cfmt"
	"github.com/samber/lo"
)

// const DISCOUNT_TABLE_NAME = "shop_discount"

var _ StoreInterface = (*Store)(nil) // verify it extends the interface

type Store struct {
	discountTableName  string
	orderTableName     string
	productTableName   string
	db                 *sql.DB
	dbDriverName       string
	timeoutSeconds     int64
	automigrateEnabled bool
	debugEnabled       bool
}

// AutoMigrate auto migrate
func (store *Store) AutoMigrate() error {
	sql := store.sqlDiscountTableCreate()

	_, err := store.db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	sql = store.sqlOrderTableCreate()

	_, err = store.db.Exec(sql)

	if err != nil {
		log.Println(err)
		return err
	}

	sql = store.sqlProductTableCreate()

	_, err = store.db.Exec(sql)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// EnableDebug - enables the debug option
func (st *Store) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

func (store *Store) DiscountCreate(discount DiscountInterface) error {
	discount.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	discount.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

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

	_, err := store.db.Exec(sqlStr, params...)

	if err != nil {
		return err
	}

	discount.MarkAsNotDirty()

	return nil
}

func (store *Store) DiscountDelete(discount DiscountInterface) error {
	if discount == nil {
		return errors.New("discount is nil")
	}

	return store.DiscountDeleteByID(discount.ID())
}

func (store *Store) DiscountDeleteByID(id string) error {
	if id == "" {
		return errors.New("discount id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.discountTableName).
		Prepared(true).
		Where(goqu.C("id").Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, params...)

	return err
}

func (store *Store) DiscountFindByID(id string) (DiscountInterface, error) {
	if id == "" {
		return nil, errors.New("discount id is empty")
	}

	list, err := store.DiscountList(DiscountQueryOptions{
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

func (store *Store) DiscountFindByCode(code string) (DiscountInterface, error) {
	if code == "" {
		return nil, errors.New("discount code is empty")
	}

	list, err := store.DiscountList(DiscountQueryOptions{
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

func (store *Store) DiscountList(options DiscountQueryOptions) ([]DiscountInterface, error) {
	q := store.discountQuery(options)

	sqlStr, _, errSql := q.Select().ToSQL()

	if errSql != nil {
		return []DiscountInterface{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)
	modelMaps, err := db.SelectToMapString(sqlStr)
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

func (store *Store) DiscountSoftDelete(discount DiscountInterface) error {
	if discount == nil {
		return errors.New("discount is nil")
	}

	discount.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.DiscountUpdate(discount)
}

func (store *Store) DiscountSoftDeleteByID(id string) error {
	discount, err := store.DiscountFindByID(id)

	if err != nil {
		return err
	}

	return store.DiscountSoftDelete(discount)
}

func (store *Store) DiscountUpdate(discount DiscountInterface) error {
	if discount == nil {
		return errors.New("order is nil")
	}

	discount.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := discount.DataChanged()

	delete(dataChanged, "id")   // ID is not updateable
	delete(dataChanged, "hash") // Hash is not updateable
	delete(dataChanged, "data") // Data is not updateable

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

	_, err := store.db.Exec(sqlStr, params...)

	discount.MarkAsNotDirty()

	return err
}

func (store *Store) discountQuery(options DiscountQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).From(store.discountTableName)

	if options.ID != "" {
		q = q.Where(goqu.C("id").Eq(options.ID))
	}

	if options.Status != "" {
		q = q.Where(goqu.C("status").Eq(options.Status))
	}

	if len(options.StatusIn) > 0 {
		q = q.Where(goqu.C("status").In(options.StatusIn))
	}

	if options.Code != "" {
		q = q.Where(goqu.C("code").Eq(options.Code))
	}

	if !options.CountOnly {
		if options.Limit > 0 {
			q = q.Limit(uint(options.Limit))
		}

		if options.Offset > 0 {
			q = q.Offset(uint(options.Offset))
		}
	}

	sortOrder := "desc"
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

	if !options.WithDeleted {
		q = q.Where(goqu.C("deleted_at").Eq(sb.NULL_DATETIME))
	}

	return q
}

func (store *Store) OrderCreate(order OrderInterface) error {
	order.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	order.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	order.SetDeletedAt(sb.NULL_DATETIME)

	data := order.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.orderTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		cfmt.Infoln(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, params...)

	if err != nil {
		return err
	}

	order.MarkAsNotDirty()

	return nil
}

func (store *Store) OrderDelete(order OrderInterface) error {
	if order == nil {
		return errors.New("order is nil")
	}

	return store.OrderDeleteByID(order.ID())
}

func (store *Store) OrderDeleteByID(id string) error {
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

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, params...)

	return err
}

func (store *Store) OrderSoftDelete(order OrderInterface) error {
	if order == nil {
		return errors.New("order is empty")
	}

	order.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.OrderUpdate(order)
}

func (store *Store) OrderSoftDeleteByID(id string) error {
	order, err := store.OrderFindByID(id)

	if err != nil {
		return err
	}

	return store.OrderSoftDelete(order)
}

func (store *Store) OrderFindByID(id string) (OrderInterface, error) {
	if id == "" {
		return nil, errors.New("order id is empty")
	}

	list, err := store.OrderList(OrderQueryOptions{
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

func (store *Store) OrderList(options OrderQueryOptions) ([]OrderInterface, error) {
	q := store.orderQuery(options)

	sqlStr, _, errSql := q.Select().ToSQL()

	if errSql != nil {
		return []OrderInterface{}, nil
	}

	if store.debugEnabled {
		cfmt.Infoln(sqlStr)
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)
	modelMaps, err := db.SelectToMapString(sqlStr)
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

func (store *Store) OrderUpdate(order OrderInterface) error {
	if order == nil {
		return errors.New("order is nil")
	}

	order.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := order.DataChanged()

	delete(dataChanged, "id") // ID is not updateable

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

	if store.debugEnabled {
		cfmt.Infoln(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, params...)

	order.MarkAsNotDirty()

	return err
}

func (store *Store) orderQuery(options OrderQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).From(store.orderTableName)

	if options.ID != "" {
		q = q.Where(goqu.C("id").Eq(options.ID))
	}

	if options.UserID != "" {
		q = q.Where(goqu.C("user_id").Eq(options.UserID))
	}

	if options.ExamID != "" {
		q = q.Where(goqu.C("exam_id").Eq(options.ExamID))
	}

	if options.Status != "" {
		q = q.Where(goqu.C("status").Eq(options.Status))
	}

	if len(options.StatusIn) > 0 {
		q = q.Where(goqu.C("status").In(options.StatusIn))
	}

	if !options.CountOnly {
		if options.Limit > 0 {
			q = q.Limit(uint(options.Limit))
		}

		if options.Offset > 0 {
			q = q.Offset(uint(options.Offset))
		}
	}

	sortOrder := "desc"
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

	if !options.WithDeleted {
		q = q.Where(goqu.C("deleted_at").Eq(sb.NULL_DATETIME))
	}

	return q
}

func (store *Store) ProductCreate(product ProductInterface) error {
	product.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	product.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	product.SetDeletedAt(sb.NULL_DATETIME)

	data := product.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.productTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		cfmt.Infoln(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, params...)

	if err != nil {
		return err
	}

	product.MarkAsNotDirty()

	return nil
}

func (store *Store) ProductSoftDelete(product ProductInterface) error {
	if product == nil {
		return errors.New("product is empty")
	}

	product.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.ProductUpdate(product)
}

func (store *Store) ProductSoftDeleteByID(id string) error {
	product, err := store.ProductFindByID(id)

	if err != nil {
		return err
	}

	return store.ProductSoftDelete(product)
}

func (store *Store) ProductFindByID(id string) (ProductInterface, error) {
	if id == "" {
		return nil, errors.New("product id is empty")
	}

	list, err := store.ProductList(ProductQueryOptions{
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

func (store *Store) ProductList(options ProductQueryOptions) ([]ProductInterface, error) {
	q := store.productQuery(options)

	sqlStr, _, errSql := q.Select().ToSQL()

	if errSql != nil {
		return []ProductInterface{}, nil
	}

	if store.debugEnabled {
		cfmt.Infoln(sqlStr)
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)
	modelMaps, err := db.SelectToMapString(sqlStr)
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

func (store *Store) ProductUpdate(product ProductInterface) error {
	if product == nil {
		return errors.New("product is nil")
	}

	product.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := product.DataChanged()

	delete(dataChanged, "id") // ID is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.productTableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C("id").Eq(product.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		cfmt.Infoln(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, params...)

	product.MarkAsNotDirty()

	return err
}

func (store *Store) productQuery(options ProductQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).From(store.productTableName)

	if options.ID != "" {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID))
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

	if !options.WithDeleted {
		q = q.Where(goqu.C("deleted_at").Eq(sb.NULL_DATETIME))
	}

	return q
}

type DiscountQueryOptions struct {
	ID          string
	IDIn        []string
	Status      string
	StatusIn    []string
	Code        string
	Offset      int
	Limit       int
	SortOrder   string
	OrderBy     string
	CountOnly   bool
	WithDeleted bool
}

type OrderQueryOptions struct {
	ID          string
	IDIn        string
	UserID      string
	ExamID      string
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
	ID          string
	IDIn        string
	Status      string
	StatusIn    []string
	Title       string
	Offset      int
	Limit       int
	SortOrder   string
	OrderBy     string
	CountOnly   bool
	WithDeleted bool
}
