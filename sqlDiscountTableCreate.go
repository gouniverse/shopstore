package shopstore

import "github.com/gouniverse/sb"

// sqlDiscountTableCreate returns a SQL string for creating the discount table
func (st *Store) sqlDiscountTableCreate() string {
	sql := sb.NewBuilder(st.dbDriverName).
		Table(st.discountTableName).
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
			Name:   "title",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name: "description",
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name:   COLUMN_TYPE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 20,
		}).
		Column(sb.Column{
			Name:     COLUMN_AMOUNT,
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   10,
			Decimals: 2,
		}).
		Column(sb.Column{
			Name:   "code",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name: COLUMN_STARTS_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_ENDS_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
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
