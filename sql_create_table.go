package logstore

import "github.com/dracory/sb"

// sqlCreateTable returns a SQL string for creating the setting table
func (store *storeImplementation) sqlCreateTable() (string, error) {
	sql, err := sb.NewBuilder(sb.DatabaseDriverName(store.db)).
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
			Length: 20,
		}).
		Column(sb.Column{
			Name: COLUMN_MESSAGE,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_CONTEXT,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_TIME,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()

	return sql, err
}

// sqlDropTable returns a SQL string for dropping the log table
func (store *storeImplementation) sqlDropTable() (string, error) {
	sql, err := sb.NewBuilder(sb.DatabaseDriverName(store.db)).
		Table(store.logTableName).
		Drop()
	return sql, err
}
