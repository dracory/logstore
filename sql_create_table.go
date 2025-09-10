package logstore

import "github.com/dracory/sb"

// SqlCreateTable returns a SQL string for creating the setting table
func (store *storeImplementation) SqlCreateTable() string {
	sql := sb.NewBuilder(sb.DatabaseDriverName(store.db)).
		Table(store.logTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   COLUMN_LEVEL,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   COLUMN_MESSAGE,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 510,
		}).
		Column(sb.Column{
			Name: COLUMN_CONTEXT,
			Type: sb.COLUMN_TYPE_LONGTEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_TIME,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()

	return sql
}
