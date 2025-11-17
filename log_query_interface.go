package logstore

import "github.com/doug-martin/goqu/v9"

// LogQueryInterface defines the interface for querying logs
type LogQueryInterface interface {
	// Validation method
	Validate() error

	// Dataset conversion methods
	ToSelectDataset(store StoreInterface) (selectDataset *goqu.SelectDataset, columns []any, err error)

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
}
