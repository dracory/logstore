package logstore

import (
	"errors"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/dracory/sb"
)

// logQuery implements the LogQueryInterface
type logQuery struct {
	isIDSet bool
	id      string

	isIDInSet bool
	idIn      []string

	isLevelSet bool
	level      string

	isLevelInSet bool
	levelIn      []string

	isTimeGteSet bool
	timeGte      string

	isTimeLteSet bool
	timeLte      string

	isLimitSet bool
	limit      int

	isOffsetSet bool
	offset      int

	isOrderDirectionSet bool
	orderDirection      string

	isOrderBySet bool
	orderBy      string
}

var _ LogQueryInterface = (*logQuery)(nil)

// LogQuery creates a new log query
func LogQuery() LogQueryInterface {
	return &logQuery{}
}

// Validate validates the query parameters
func (q *logQuery) Validate() error {
	if q.IsIDSet() && q.GetID() == "" {
		return errors.New("log query: id cannot be empty")
	}

	if q.IsIDInSet() && len(q.GetIDIn()) < 1 {
		return errors.New("log query: id_in cannot be empty array")
	}

	if q.IsLevelSet() && q.GetLevel() == "" {
		return errors.New("log query: level cannot be empty")
	}

	if q.IsLevelInSet() && len(q.GetLevelIn()) < 1 {
		return errors.New("log query: level_in cannot be empty array")
	}

	if q.IsLimitSet() && q.GetLimit() < 0 {
		return errors.New("log query: limit cannot be negative")
	}

	if q.IsOffsetSet() && q.GetOffset() < 0 {
		return errors.New("log query: offset cannot be negative")
	}

	return nil
}

func (q *logQuery) ToSelectDataset(st StoreInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if st == nil {
		return nil, []any{}, errors.New("store cannot be nil")
	}

	if err := q.Validate(); err != nil {
		return nil, []any{}, err
	}

	sql := goqu.Dialect(st.GetDriverName()).From(st.GetLogTableName())

	// ID filter
	if q.IsIDSet() {
		sql = sql.Where(goqu.C(COLUMN_ID).Eq(q.GetID()))
	}

	// ID IN filter
	if q.IsIDInSet() {
		sql = sql.Where(goqu.C(COLUMN_ID).In(q.GetIDIn()))
	}

	// Level filter
	if q.IsLevelSet() {
		sql = sql.Where(goqu.C(COLUMN_LEVEL).Eq(q.GetLevel()))
	}

	// Level IN filter
	if q.IsLevelInSet() {
		sql = sql.Where(goqu.C(COLUMN_LEVEL).In(q.GetLevelIn()))
	}

	// Time filters
	if q.IsTimeGteSet() {
		sql = sql.Where(goqu.C(COLUMN_TIME).Gte(q.GetTimeGte()))
	}

	if q.IsTimeLteSet() {
		sql = sql.Where(goqu.C(COLUMN_TIME).Lte(q.GetTimeLte()))
	}

	// Limit and offset
	if q.IsLimitSet() {
		sql = sql.Limit(uint(q.GetLimit()))
	}

	if q.IsOffsetSet() {
		sql = sql.Offset(uint(q.GetOffset()))
	}

	// Sort order
	orderDirection := sb.DESC
	if q.IsOrderDirectionSet() {
		orderDirection = q.GetOrderDirection()
	}

	if q.IsOrderBySet() {
		if strings.EqualFold(orderDirection, sb.ASC) {
			sql = sql.Order(goqu.I(q.GetOrderBy()).Asc())
		} else {
			sql = sql.Order(goqu.I(q.GetOrderBy()).Desc())
		}
	}

	return sql, []any{}, nil
}

// ============================================================================
// == Getters and Setters
// ============================================================================

func (q *logQuery) IsIDSet() bool {
	return q.isIDSet
}

func (q *logQuery) GetID() string {
	if q.IsIDSet() {
		return q.id
	}

	return ""
}

func (q *logQuery) SetID(id string) LogQueryInterface {
	q.isIDSet = true
	q.id = id
	return q
}

func (q *logQuery) IsIDInSet() bool {
	return q.isIDInSet
}

func (q *logQuery) GetIDIn() []string {
	if q.IsIDInSet() {
		return q.idIn
	}

	return []string{}
}

func (q *logQuery) SetIDIn(idIn []string) LogQueryInterface {
	q.isIDInSet = true
	q.idIn = idIn
	return q
}

func (q *logQuery) IsLevelSet() bool {
	return q.isLevelSet
}

func (q *logQuery) GetLevel() string {
	if q.IsLevelSet() {
		return q.level
	}

	return ""
}

func (q *logQuery) SetLevel(level string) LogQueryInterface {
	q.isLevelSet = true
	q.level = level
	return q
}

func (q *logQuery) IsLevelInSet() bool {
	return q.isLevelInSet
}

func (q *logQuery) GetLevelIn() []string {
	if q.IsLevelInSet() {
		return q.levelIn
	}

	return []string{}
}

func (q *logQuery) SetLevelIn(levelIn []string) LogQueryInterface {
	q.isLevelInSet = true
	q.levelIn = levelIn
	return q
}

func (q *logQuery) IsTimeGteSet() bool {
	return q.isTimeGteSet
}

func (q *logQuery) GetTimeGte() string {
	if q.IsTimeGteSet() {
		return q.timeGte
	}

	return ""
}

func (q *logQuery) SetTimeGte(timeGte string) LogQueryInterface {
	q.isTimeGteSet = true
	q.timeGte = timeGte
	return q
}

func (q *logQuery) IsTimeLteSet() bool {
	return q.isTimeLteSet
}

func (q *logQuery) GetTimeLte() string {
	if q.IsTimeLteSet() {
		return q.timeLte
	}

	return ""
}

func (q *logQuery) SetTimeLte(timeLte string) LogQueryInterface {
	q.isTimeLteSet = true
	q.timeLte = timeLte
	return q
}

func (q *logQuery) IsLimitSet() bool {
	return q.isLimitSet
}

func (q *logQuery) GetLimit() int {
	if q.IsLimitSet() {
		return q.limit
	}

	return 0
}

func (q *logQuery) SetLimit(limit int) LogQueryInterface {
	q.isLimitSet = true
	q.limit = limit
	return q
}

func (q *logQuery) IsOffsetSet() bool {
	return q.isOffsetSet
}

func (q *logQuery) GetOffset() int {
	if q.IsOffsetSet() {
		return q.offset
	}

	return 0
}

func (q *logQuery) SetOffset(offset int) LogQueryInterface {
	q.isOffsetSet = true
	q.offset = offset
	return q
}

func (q *logQuery) IsOrderDirectionSet() bool {
	return q.isOrderDirectionSet
}

func (q *logQuery) GetOrderDirection() string {
	if q.IsOrderDirectionSet() {
		return q.orderDirection
	}

	return ""
}

func (q *logQuery) SetOrderDirection(orderDirection string) LogQueryInterface {
	q.isOrderDirectionSet = true
	q.orderDirection = orderDirection
	return q
}

func (q *logQuery) IsOrderBySet() bool {
	return q.isOrderBySet
}

func (q *logQuery) GetOrderBy() string {
	if q.IsOrderBySet() {
		return q.orderBy
	}

	return ""
}

func (q *logQuery) SetOrderBy(orderBy string) LogQueryInterface {
	q.isOrderBySet = true
	q.orderBy = orderBy
	return q
}
