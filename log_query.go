package logstore

import "errors"

// LogQueryInterface defines the interface for querying logs
type LogQueryInterface interface {
	// Validation method
	Validate() error

	// Field query methods

	IsIDSet() bool
	GetID() string
	SetID(id string) LogQueryInterface

	IsIDInSet() bool
	GetIDIn() []string
	SetIDIn(ids []string) LogQueryInterface

	IsLevelSet() bool
	GetLevel() string
	SetLevel(level string) LogQueryInterface

	IsLevelInSet() bool
	GetLevelIn() []string
	SetLevelIn(levels []string) LogQueryInterface

	IsMessageContainsSet() bool
	GetMessageContains() string
	SetMessageContains(term string) LogQueryInterface

	IsMessageNotContainsSet() bool
	GetMessageNotContains() string
	SetMessageNotContains(term string) LogQueryInterface

	IsContextContainsSet() bool
	GetContextContains() string
	SetContextContains(term string) LogQueryInterface

	IsContextNotContainsSet() bool
	GetContextNotContains() string
	SetContextNotContains(term string) LogQueryInterface

	IsTimeGteSet() bool
	GetTimeGte() string
	SetTimeGte(time string) LogQueryInterface

	IsTimeLteSet() bool
	GetTimeLte() string
	SetTimeLte(time string) LogQueryInterface

	IsLimitSet() bool
	GetLimit() int
	SetLimit(limit int) LogQueryInterface

	IsOffsetSet() bool
	GetOffset() int
	SetOffset(offset int) LogQueryInterface

	IsOrderBySet() bool
	GetOrderBy() string
	SetOrderBy(orderBy string) LogQueryInterface

	IsOrderDirectionSet() bool
	GetOrderDirection() string
	SetOrderDirection(orderDirection string) LogQueryInterface

	IsColumnsSet() bool
	GetColumns() []string
	SetColumns(columns []string) LogQueryInterface
}

// logQueryImplementation implements the LogQueryInterface
type logQueryImplementation struct {
	isIDSet bool
	id      string

	isIDInSet bool
	idIn      []string

	isLevelSet bool
	level      string

	isLevelInSet bool
	levelIn      []string

	isMessageContainsSet    bool
	messageContains         string
	isMessageNotContainsSet bool
	messageNotContains      string

	isContextContainsSet    bool
	contextContains         string
	isContextNotContainsSet bool
	contextNotContains      string

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

	isColumnsSet bool
	columns      []string
}

var _ LogQueryInterface = (*logQueryImplementation)(nil)

// LogQuery creates a new log query
func LogQuery() LogQueryInterface {
	return &logQueryImplementation{}
}

// Validate validates the query parameters
func (q *logQueryImplementation) Validate() error {
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

	if q.IsMessageContainsSet() && q.GetMessageContains() == "" {
		return errors.New("log query: message_contains cannot be empty")
	}

	if q.IsMessageNotContainsSet() && q.GetMessageNotContains() == "" {
		return errors.New("log query: message_not_contains cannot be empty")
	}

	if q.IsContextContainsSet() && q.GetContextContains() == "" {
		return errors.New("log query: context_contains cannot be empty")
	}

	if q.IsContextNotContainsSet() && q.GetContextNotContains() == "" {
		return errors.New("log query: context_not_contains cannot be empty")
	}

	if q.IsLimitSet() && q.GetLimit() < 0 {
		return errors.New("log query: limit cannot be negative")
	}

	if q.IsOffsetSet() && q.GetOffset() < 0 {
		return errors.New("log query: offset cannot be negative")
	}

	return nil
}

// ============================================================================
// == Getters and Setters
// ============================================================================

func (q *logQueryImplementation) IsIDSet() bool {
	return q.isIDSet
}

func (q *logQueryImplementation) GetID() string {
	if q.IsIDSet() {
		return q.id
	}
	return ""
}

func (q *logQueryImplementation) SetID(id string) LogQueryInterface {
	q.isIDSet = true
	q.id = id
	return q
}

func (q *logQueryImplementation) IsIDInSet() bool {
	return q.isIDInSet
}

func (q *logQueryImplementation) GetIDIn() []string {
	if q.IsIDInSet() {
		return q.idIn
	}
	return []string{}
}

func (q *logQueryImplementation) SetIDIn(idIn []string) LogQueryInterface {
	q.isIDInSet = true
	q.idIn = idIn
	return q
}

func (q *logQueryImplementation) IsLevelSet() bool {
	return q.isLevelSet
}

func (q *logQueryImplementation) GetLevel() string {
	if q.IsLevelSet() {
		return q.level
	}
	return ""
}

func (q *logQueryImplementation) SetLevel(level string) LogQueryInterface {
	q.isLevelSet = true
	q.level = level
	return q
}

func (q *logQueryImplementation) IsLevelInSet() bool {
	return q.isLevelInSet
}

func (q *logQueryImplementation) GetLevelIn() []string {
	if q.IsLevelInSet() {
		return q.levelIn
	}
	return []string{}
}

func (q *logQueryImplementation) SetLevelIn(levelIn []string) LogQueryInterface {
	q.isLevelInSet = true
	q.levelIn = levelIn
	return q
}

func (q *logQueryImplementation) IsMessageContainsSet() bool {
	return q.isMessageContainsSet
}

func (q *logQueryImplementation) GetMessageContains() string {
	if q.IsMessageContainsSet() {
		return q.messageContains
	}
	return ""
}

func (q *logQueryImplementation) SetMessageContains(term string) LogQueryInterface {
	q.isMessageContainsSet = true
	q.messageContains = term
	return q
}

func (q *logQueryImplementation) IsMessageNotContainsSet() bool {
	return q.isMessageNotContainsSet
}

func (q *logQueryImplementation) GetMessageNotContains() string {
	if q.IsMessageNotContainsSet() {
		return q.messageNotContains
	}
	return ""
}

func (q *logQueryImplementation) SetMessageNotContains(term string) LogQueryInterface {
	q.isMessageNotContainsSet = true
	q.messageNotContains = term
	return q
}

func (q *logQueryImplementation) IsContextContainsSet() bool {
	return q.isContextContainsSet
}

func (q *logQueryImplementation) GetContextContains() string {
	if q.IsContextContainsSet() {
		return q.contextContains
	}
	return ""
}

func (q *logQueryImplementation) SetContextContains(term string) LogQueryInterface {
	q.isContextContainsSet = true
	q.contextContains = term
	return q
}

func (q *logQueryImplementation) IsContextNotContainsSet() bool {
	return q.isContextNotContainsSet
}

func (q *logQueryImplementation) GetContextNotContains() string {
	if q.IsContextNotContainsSet() {
		return q.contextNotContains
	}
	return ""
}

func (q *logQueryImplementation) SetContextNotContains(term string) LogQueryInterface {
	q.isContextNotContainsSet = true
	q.contextNotContains = term
	return q
}

func (q *logQueryImplementation) IsTimeGteSet() bool {
	return q.isTimeGteSet
}

func (q *logQueryImplementation) GetTimeGte() string {
	if q.IsTimeGteSet() {
		return q.timeGte
	}
	return ""
}

func (q *logQueryImplementation) SetTimeGte(time string) LogQueryInterface {
	q.isTimeGteSet = true
	q.timeGte = time
	return q
}

func (q *logQueryImplementation) IsTimeLteSet() bool {
	return q.isTimeLteSet
}

func (q *logQueryImplementation) GetTimeLte() string {
	if q.IsTimeLteSet() {
		return q.timeLte
	}
	return ""
}

func (q *logQueryImplementation) SetTimeLte(time string) LogQueryInterface {
	q.isTimeLteSet = true
	q.timeLte = time
	return q
}

func (q *logQueryImplementation) IsLimitSet() bool {
	return q.isLimitSet
}

func (q *logQueryImplementation) GetLimit() int {
	if q.IsLimitSet() {
		return q.limit
	}
	return 0
}

func (q *logQueryImplementation) SetLimit(limit int) LogQueryInterface {
	q.isLimitSet = true
	q.limit = limit
	return q
}

func (q *logQueryImplementation) IsOffsetSet() bool {
	return q.isOffsetSet
}

func (q *logQueryImplementation) GetOffset() int {
	if q.IsOffsetSet() {
		return q.offset
	}
	return 0
}

func (q *logQueryImplementation) SetOffset(offset int) LogQueryInterface {
	q.isOffsetSet = true
	q.offset = offset
	return q
}

func (q *logQueryImplementation) IsOrderDirectionSet() bool {
	return q.isOrderDirectionSet
}

func (q *logQueryImplementation) GetOrderDirection() string {
	if q.IsOrderDirectionSet() {
		return q.orderDirection
	}
	return ""
}

func (q *logQueryImplementation) SetOrderDirection(orderDirection string) LogQueryInterface {
	q.isOrderDirectionSet = true
	q.orderDirection = orderDirection
	return q
}

func (q *logQueryImplementation) IsOrderBySet() bool {
	return q.isOrderBySet
}

func (q *logQueryImplementation) GetOrderBy() string {
	if q.IsOrderBySet() {
		return q.orderBy
	}
	return ""
}

func (q *logQueryImplementation) SetOrderBy(orderBy string) LogQueryInterface {
	q.isOrderBySet = true
	q.orderBy = orderBy
	return q
}

func (q *logQueryImplementation) IsColumnsSet() bool {
	return q.isColumnsSet
}

func (q *logQueryImplementation) GetColumns() []string {
	if q.IsColumnsSet() {
		return q.columns
	}
	return []string{}
}

func (q *logQueryImplementation) SetColumns(columns []string) LogQueryInterface {
	q.isColumnsSet = true
	q.columns = columns
	return q
}
