package dto

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type FilterOperator string

const (
	OpEqual              FilterOperator = "eq"
	OpNotEqual           FilterOperator = "neq"
	OpGreaterThan        FilterOperator = "gt"
	OpGreaterThanOrEqual FilterOperator = "gte"
	OpLessThan           FilterOperator = "lt"
	OpLessThanOrEqual    FilterOperator = "lte"
	OpLike               FilterOperator = "like"
	OpILike              FilterOperator = "ilike"
	OpIn                 FilterOperator = "in"
	OpNotIn              FilterOperator = "nin"
	OpIsNull             FilterOperator = "null"
	OpIsNotNull          FilterOperator = "notnull"
)

type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

type Filter struct {
	Field    string      `json:"field" binding:"required"`
	Operator string      `json:"operator" binding:"required"`
	Value    interface{} `json:"value"`
}

type Sort struct {
	Field     string `json:"field" binding:"required"`
	Direction string `json:"direction" binding:"required"`
}

func (f *Filter) Validate() error {
	if f.Field == "" {
		return errors.New("filter field is required")
	}
	if f.Operator == "" {
		return errors.New("filter operator is required")
	}

	validOps := []string{
		string(OpEqual), string(OpNotEqual),
		string(OpGreaterThan), string(OpGreaterThanOrEqual),
		string(OpLessThan), string(OpLessThanOrEqual),
		string(OpLike), string(OpILike),
		string(OpIn), string(OpNotIn),
		string(OpIsNull), string(OpIsNotNull),
	}

	valid := false
	for _, op := range validOps {
		if f.Operator == op {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("invalid operator: %s", f.Operator)
	}

	if f.Operator != string(OpIsNull) && f.Operator != string(OpIsNotNull) {
		if f.Value == nil {
			return errors.New("filter value is required for this operator")
		}
	}

	return nil
}

func (s *Sort) Validate() error {
	if s.Field == "" {
		return errors.New("sort field is required")
	}

	dir := strings.ToLower(s.Direction)
	if dir != string(SortAsc) && dir != string(SortDesc) {
		return fmt.Errorf("invalid sort direction: %s (must be 'asc' or 'desc')", s.Direction)
	}

	return nil
}

func (f *Filter) GetSQLOperator() string {
	switch FilterOperator(f.Operator) {
	case OpEqual:
		return "="
	case OpNotEqual:
		return "!="
	case OpGreaterThan:
		return ">"
	case OpGreaterThanOrEqual:
		return ">="
	case OpLessThan:
		return "<"
	case OpLessThanOrEqual:
		return "<="
	case OpLike:
		return "LIKE"
	case OpILike:
		return "ILIKE"
	case OpIn:
		return "IN"
	case OpNotIn:
		return "NOT IN"
	case OpIsNull:
		return "IS NULL"
	case OpIsNotNull:
		return "IS NOT NULL"
	default:
		return "="
	}
}

type FilteredRequestInterface interface {
	DecodeFilters() error
	DecodeSort() error
	GetFilters() []Filter
	GetSort() []Sort
	AddFilters(filters ...*Filter)
	AddSort(sorts ...*Sort)
}

type ListRequest struct {
	Limit  uint64 `form:"limit" binding:"required,numeric,gte=1,lte=100"`
	Page   uint64 `form:"page" binding:"required,numeric,gte=1"`
	UserID uuid.UUID
	FilteredRequest
}

type FilteredRequest struct {
	FiltersStringEncoded string `form:"filters"`
	Filters              []Filter
	SortStringEncoded    string `form:"sort"`
	Sort                 []Sort
}

func (r *FilteredRequest) DecodeFilters() error {
	filters, err := parseFilters(r.FiltersStringEncoded)
	if err != nil {
		return err
	}

	for i, f := range filters {
		if err := f.Validate(); err != nil {
			return fmt.Errorf("filter[%d] validation error: %w", i, err)
		}
	}

	r.Filters = filters
	return nil
}

func (r *FilteredRequest) DecodeSort() error {
	sorts, err := parseSort(r.SortStringEncoded)
	if err != nil {
		return err
	}

	for i := range sorts {
		sorts[i].Direction = strings.ToLower(sorts[i].Direction)
		if err := sorts[i].Validate(); err != nil {
			return fmt.Errorf("sort[%d] validation error: %w", i, err)
		}
	}

	r.Sort = sorts
	return nil
}

func (r *FilteredRequest) GetFilters() []Filter {
	if r.Filters == nil {
		return []Filter{}
	}
	return r.Filters
}

func (r *FilteredRequest) GetSort() []Sort {
	if r.Sort == nil {
		return []Sort{}
	}
	return r.Sort
}

func (r *FilteredRequest) AddFilters(filters ...*Filter) {
	for _, f := range filters {
		if f == nil {
			continue
		}
		if err := f.Validate(); err != nil {
			continue // Skip invalid filters
		}
		r.Filters = append(r.Filters, *f)
	}
}

func (r *FilteredRequest) AddSort(sorts ...*Sort) {
	for _, s := range sorts {
		if s == nil {
			continue
		}
		s.Direction = strings.ToLower(s.Direction)
		if err := s.Validate(); err != nil {
			continue // Skip invalid sorts
		}
		r.Sort = append(r.Sort, *s)
	}
}

func (r *FilteredRequest) HasFilters() bool {
	return len(r.Filters) > 0
}

func (r *FilteredRequest) HasSort() bool {
	return len(r.Sort) > 0
}

func parseFilters(b64 string) ([]Filter, error) {
	var filters []Filter

	if b64 == "" {
		return filters, nil
	}

	decoded, err := base64.URLEncoding.DecodeString(b64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode base64 filters")
	}

	if err := json.Unmarshal(decoded, &filters); err != nil {
		return nil, errors.Wrap(err, "invalid filter JSON format")
	}

	return filters, nil
}

func parseSort(b64 string) ([]Sort, error) {
	var sorts []Sort

	if b64 == "" {
		return sorts, nil
	}

	decoded, err := base64.URLEncoding.DecodeString(b64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode base64 sort")
	}

	if err := json.Unmarshal(decoded, &sorts); err != nil {
		return nil, errors.Wrap(err, "invalid sort JSON format")
	}

	return sorts, nil
}
