package shopstore

import (
	"github.com/gouniverse/sb"
)

func (store *Store) sqlOrderTableCreate() string {
	sql := sb.NewBuilder(sb.DatabaseDriverName(store.db)).
		Table(store.orderTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   COLUMN_STATUS,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 20,
		}).
		Column(sb.Column{
			Name:   COLUMN_CUSTOMER_ID,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   COLUMN_QUANTITY,
			Type:   sb.COLUMN_TYPE_INTEGER,
			Length: 10,
		}).
		Column(sb.Column{
			Name:     COLUMN_PRICE,
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   10,
			Decimals: 2,
		}).
		Column(sb.Column{
			Name: COLUMN_METAS,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_MEMO,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_DELETED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()

	return sql
}
