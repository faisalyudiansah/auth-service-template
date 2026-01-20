package queryBuilder

import (
	"fmt"
	"strings"

	dtoPkg "github.com/faisalyudiansah/auth-service-template/pkg/dto"
)

// QueryBuilder helps build SQL queries with filters and sort
type QueryBuilder struct {
	baseQuery       string
	whereConditions []string
	args            []interface{}
	argCounter      int
	sortClauses     []string
	allowedFilters  map[string]string // map[dtoField]dbColumn
	allowedSorts    map[string]string // map[dtoField]dbColumn
}

// NewQueryBuilder creates a new query builder
func NewQueryBuilder(baseQuery string) *QueryBuilder {
	return &QueryBuilder{
		baseQuery:       baseQuery,
		whereConditions: []string{},
		args:            []interface{}{},
		argCounter:      1,
		sortClauses:     []string{},
		allowedFilters:  make(map[string]string),
		allowedSorts:    make(map[string]string),
	}
}

// SetAllowedFilters sets allowed filter fields mapping
func (qb *QueryBuilder) SetAllowedFilters(mapping map[string]string) *QueryBuilder {
	qb.allowedFilters = mapping
	return qb
}

// SetAllowedSorts sets allowed sort fields mapping
func (qb *QueryBuilder) SetAllowedSorts(mapping map[string]string) *QueryBuilder {
	qb.allowedSorts = mapping
	return qb
}

// AddWhere adds a WHERE condition
func (qb *QueryBuilder) AddWhere(condition string, args ...interface{}) *QueryBuilder {
	qb.whereConditions = append(qb.whereConditions, condition)
	qb.args = append(qb.args, args...)
	qb.argCounter += len(args)
	return qb
}

// ApplyFilters applies filters from request
func (qb *QueryBuilder) ApplyFilters(filters []dtoPkg.Filter) error {
	for _, filter := range filters {
		// Check if field is allowed
		dbColumn, allowed := qb.allowedFilters[filter.Field]
		if !allowed {
			return fmt.Errorf("filter on field '%s' is not allowed", filter.Field)
		}

		// Build condition based on operator
		condition, args, err := qb.buildFilterCondition(dbColumn, filter)
		if err != nil {
			return err
		}

		qb.whereConditions = append(qb.whereConditions, condition)
		qb.args = append(qb.args, args...)
	}

	return nil
}

// ApplySort applies sort from request
func (qb *QueryBuilder) ApplySort(sorts []dtoPkg.Sort) error {
	for _, sort := range sorts {
		// Check if field is allowed
		dbColumn, allowed := qb.allowedSorts[sort.Field]
		if !allowed {
			return fmt.Errorf("sort on field '%s' is not allowed", sort.Field)
		}

		direction := strings.ToUpper(sort.Direction)
		qb.sortClauses = append(qb.sortClauses, fmt.Sprintf("%s %s", dbColumn, direction))
	}

	return nil
}

// AddDefaultSort adds default sort if no sort is specified
func (qb *QueryBuilder) AddDefaultSort(field string, direction string) *QueryBuilder {
	if len(qb.sortClauses) == 0 {
		qb.sortClauses = append(qb.sortClauses, fmt.Sprintf("%s %s", field, strings.ToUpper(direction)))
	}
	return qb
}

// Build builds the final query
func (qb *QueryBuilder) Build() (string, []interface{}) {
	query := qb.baseQuery

	// Add WHERE clause
	if len(qb.whereConditions) > 0 {
		query += " WHERE " + strings.Join(qb.whereConditions, " AND ")
	}

	// Add ORDER BY clause
	if len(qb.sortClauses) > 0 {
		query += " ORDER BY " + strings.Join(qb.sortClauses, ", ")
	}

	return query, qb.args
}

// BuildWithPagination builds query with pagination
func (qb *QueryBuilder) BuildWithPagination(limit, offset uint64) (string, []interface{}) {
	query, args := qb.Build()

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", qb.argCounter, qb.argCounter+1)
	args = append(args, limit, offset)

	return query, args
}

// buildFilterCondition builds SQL condition for a filter
func (qb *QueryBuilder) buildFilterCondition(dbColumn string, filter dtoPkg.Filter) (string, []interface{}, error) {
	var condition string
	var args []interface{}

	switch dtoPkg.FilterOperator(filter.Operator) {
	case dtoPkg.OpIsNull:
		condition = fmt.Sprintf("%s IS NULL", dbColumn)

	case dtoPkg.OpIsNotNull:
		condition = fmt.Sprintf("%s IS NOT NULL", dbColumn)

	case dtoPkg.OpIn, dtoPkg.OpNotIn:
		// Handle IN and NOT IN operators
		values, ok := filter.Value.([]interface{})
		if !ok {
			return "", nil, fmt.Errorf("value for IN/NOT IN operator must be an array")
		}

		if len(values) == 0 {
			return "", nil, fmt.Errorf("value array for IN/NOT IN cannot be empty")
		}

		placeholders := make([]string, len(values))
		for i, val := range values {
			placeholders[i] = fmt.Sprintf("$%d", qb.argCounter+i)
			args = append(args, val)
		}

		operator := "IN"
		if filter.Operator == string(dtoPkg.OpNotIn) {
			operator = "NOT IN"
		}

		condition = fmt.Sprintf("%s %s (%s)", dbColumn, operator, strings.Join(placeholders, ", "))
		qb.argCounter += len(values)

	case dtoPkg.OpLike, dtoPkg.OpILike:
		// Add wildcards for LIKE/ILIKE
		value := fmt.Sprintf("%%%v%%", filter.Value)
		condition = fmt.Sprintf("%s %s $%d", dbColumn, filter.GetSQLOperator(), qb.argCounter)
		args = append(args, value)
		qb.argCounter++

	default:
		// Standard comparison operators
		condition = fmt.Sprintf("%s %s $%d", dbColumn, filter.GetSQLOperator(), qb.argCounter)
		args = append(args, filter.Value)
		qb.argCounter++
	}

	return condition, args, nil
}

// GetArgCounter returns current argument counter
func (qb *QueryBuilder) GetArgCounter() int {
	return qb.argCounter
}

// GetArgs returns all accumulated arguments
func (qb *QueryBuilder) GetArgs() []interface{} {
	return qb.args
}
