package logstore

// IsColumnsSet returns true if columns are set
func (q *logQueryImplementation) IsColumnsSet() bool {
	return q.isColumnsSet
}

// GetColumns returns the columns to select
func (q *logQueryImplementation) GetColumns() []string {
	if q.IsColumnsSet() {
		return q.columns
	}
	return []string{}
}

// SetColumns sets the columns to select
func (q *logQueryImplementation) SetColumns(columns []string) LogQueryInterface {
	q.isColumnsSet = true
	q.columns = columns
	return q
}
