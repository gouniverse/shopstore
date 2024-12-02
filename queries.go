package shopstore

import "errors"

type CategoryQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) CategoryQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) CategoryQueryInterface

	HasID() bool
	ID() string
	SetID(id string) CategoryQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) CategoryQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) CategoryQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) CategoryQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) CategoryQueryInterface

	HasSortDirection() bool
	SortDirection() string
	SetSortDirection(sortDirection string) CategoryQueryInterface

	HasParentID() bool
	ParentID() string
	SetParentID(parentID string) CategoryQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(softDeletedIncluded bool) CategoryQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) CategoryQueryInterface

	HasTitleLike() bool
	TitleLike() string
	SetTitleLike(titleLike string) CategoryQueryInterface

	hasProperty(name string) bool
}

func NewCategoryQuery() CategoryQueryInterface {
	return &categoryQueryImplementation{
		properties: make(map[string]any),
	}
}

type categoryQueryImplementation struct {
	properties map[string]any
}

func (c *categoryQueryImplementation) Validate() error {
	if c.HasID() && c.ID() == "" {
		return errors.New("category query. id cannot be empty")
	}

	if c.HasIDIn() && len(c.IDIn()) == 0 {
		return errors.New("category query. id_in cannot be empty")
	}

	if c.HasParentID() && c.ParentID() == "" {
		return errors.New("category query. parent_id cannot be empty")
	}

	if c.HasStatus() && c.Status() == "" {
		return errors.New("category query. status cannot be empty")
	}

	if c.HasTitleLike() && c.TitleLike() == "" {
		return errors.New("category query. title_like cannot be empty")
	}

	if c.HasOrderBy() && c.OrderBy() == "" {
		return errors.New("category query. order_by cannot be empty")
	}

	if c.HasSortDirection() && c.SortDirection() == "" {
		return errors.New("category query. sort_direction cannot be empty")
	}

	if c.HasLimit() && c.Limit() <= 0 {
		return errors.New("category query. limit must be greater than 0")
	}

	if c.HasOffset() && c.Offset() < 0 {
		return errors.New("category query. offset must be greater than or equal to 0")
	}

	return nil
}

func (c *categoryQueryImplementation) Columns() []string {
	if !c.hasProperty("columns") {
		return []string{}
	}

	return c.properties["columns"].([]string)
}

func (c *categoryQueryImplementation) SetColumns(columns []string) CategoryQueryInterface {
	c.properties["columns"] = columns

	return c
}

func (c *categoryQueryImplementation) HasCountOnly() bool {
	return c.hasProperty("count_only")
}

func (c *categoryQueryImplementation) IsCountOnly() bool {
	if !c.HasCountOnly() {
		return false
	}

	return c.properties["count_only"].(bool)
}

func (c *categoryQueryImplementation) SetCountOnly(countOnly bool) CategoryQueryInterface {
	c.properties["count_only"] = countOnly

	return c
}

func (c *categoryQueryImplementation) HasID() bool {
	return c.hasProperty("id")
}

func (c *categoryQueryImplementation) ID() string {
	if !c.HasID() {
		return ""
	}

	return c.properties["id"].(string)
}

func (c *categoryQueryImplementation) SetID(id string) CategoryQueryInterface {
	c.properties["id"] = id

	return c
}

func (c *categoryQueryImplementation) HasIDIn() bool {
	return c.hasProperty("id_in")
}

func (c *categoryQueryImplementation) IDIn() []string {
	if !c.HasIDIn() {
		return []string{}
	}

	return c.properties["id_in"].([]string)
}

func (c *categoryQueryImplementation) SetIDIn(idIn []string) CategoryQueryInterface {
	c.properties["id_in"] = idIn

	return c
}

func (c *categoryQueryImplementation) HasLimit() bool {
	return c.hasProperty("limit")
}

func (c *categoryQueryImplementation) Limit() int {
	if !c.HasLimit() {
		return 0
	}

	return c.properties["limit"].(int)
}

func (c *categoryQueryImplementation) SetLimit(limit int) CategoryQueryInterface {
	c.properties["limit"] = limit

	return c
}

func (c *categoryQueryImplementation) HasOffset() bool {
	return c.hasProperty("offset")
}

func (c *categoryQueryImplementation) Offset() int {
	if !c.HasOffset() {
		return 0
	}

	return c.properties["offset"].(int)
}

func (c *categoryQueryImplementation) SetOffset(offset int) CategoryQueryInterface {
	c.properties["offset"] = offset

	return c
}

func (c *categoryQueryImplementation) HasOrderBy() bool {
	return c.hasProperty("order_by")
}

func (c *categoryQueryImplementation) OrderBy() string {
	if !c.HasOrderBy() {
		return ""
	}

	return c.properties["order_by"].(string)
}

func (c *categoryQueryImplementation) SetOrderBy(orderBy string) CategoryQueryInterface {
	c.properties["order_by"] = orderBy

	return c
}

func (c *categoryQueryImplementation) HasSortDirection() bool {
	return c.hasProperty("sort_direction")
}

func (c *categoryQueryImplementation) SortDirection() string {
	if !c.HasSortDirection() {
		return ""
	}

	return c.properties["sort_direction"].(string)
}

func (c *categoryQueryImplementation) SetSortDirection(sortDirection string) CategoryQueryInterface {
	c.properties["sort_direction"] = sortDirection

	return c
}

func (c *categoryQueryImplementation) HasParentID() bool {
	return c.hasProperty("parent_id")
}

func (c *categoryQueryImplementation) ParentID() string {
	if !c.HasParentID() {
		return ""
	}

	return c.properties["parent_id"].(string)
}

func (c *categoryQueryImplementation) SetParentID(parentID string) CategoryQueryInterface {
	c.properties["parent_id"] = parentID

	return c
}

func (c *categoryQueryImplementation) HasSoftDeletedIncluded() bool {
	return c.hasProperty("soft_deleted_included")
}

func (c *categoryQueryImplementation) SoftDeletedIncluded() bool {
	if !c.HasSoftDeletedIncluded() {
		return false
	}

	return c.properties["soft_deleted_included"].(bool)
}

func (c *categoryQueryImplementation) SetSoftDeletedIncluded(softDeletedIncluded bool) CategoryQueryInterface {
	c.properties["soft_deleted_included"] = softDeletedIncluded

	return c
}

func (c *categoryQueryImplementation) HasStatus() bool {
	return c.hasProperty("status")
}

func (c *categoryQueryImplementation) Status() string {
	if !c.HasStatus() {
		return ""
	}

	return c.properties["status"].(string)
}

func (c *categoryQueryImplementation) SetStatus(status string) CategoryQueryInterface {
	c.properties["status"] = status

	return c
}

func (c *categoryQueryImplementation) HasTitleLike() bool {
	return c.hasProperty("title_like")
}

func (c *categoryQueryImplementation) TitleLike() string {
	if !c.HasTitleLike() {
		return ""
	}

	return c.properties["title_like"].(string)
}

func (c *categoryQueryImplementation) SetTitleLike(titleLike string) CategoryQueryInterface {
	c.properties["title_like"] = titleLike

	return c
}

func (c *categoryQueryImplementation) hasProperty(name string) bool {
	_, ok := c.properties[name]
	return ok
}

type MediaQueryInterface interface {
	Validate() error

	Columns() []string
	SetColumns(columns []string) MediaQueryInterface

	HasCountOnly() bool
	IsCountOnly() bool
	SetCountOnly(countOnly bool) MediaQueryInterface

	HasEntityID() bool
	EntityID() string
	SetEntityID(entityID string) MediaQueryInterface

	HasID() bool
	ID() string
	SetID(id string) MediaQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) MediaQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) MediaQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) MediaQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) MediaQueryInterface

	HasSortDirection() bool
	SortDirection() string
	SetSortDirection(sortDirection string) MediaQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(softDeletedIncluded bool) MediaQueryInterface

	HasStatus() bool
	Status() string
	SetStatus(status string) MediaQueryInterface

	HasTitleLike() bool
	TitleLike() string
	SetTitleLike(titleLike string) MediaQueryInterface

	HasType() bool
	Type() string
	SetType(mediaType string) MediaQueryInterface

	hasProperty(name string) bool
}

func NewMediaQuery() MediaQueryInterface {
	return &mediaQueryImplementation{
		properties: make(map[string]any),
	}
}

type mediaQueryImplementation struct {
	properties map[string]any
}

func (c *mediaQueryImplementation) Validate() error {
	if c.HasID() && c.ID() == "" {
		return errors.New("media query. id cannot be empty")
	}

	if c.HasIDIn() && len(c.IDIn()) == 0 {
		return errors.New("media query. id_in cannot be empty")
	}

	if c.HasEntityID() && c.EntityID() == "" {
		return errors.New("media query. entity_id cannot be empty")
	}

	if c.HasStatus() && c.Status() == "" {
		return errors.New("media query. status cannot be empty")
	}

	if c.HasTitleLike() && c.TitleLike() == "" {
		return errors.New("media query. title_like cannot be empty")
	}

	if c.HasOrderBy() && c.OrderBy() == "" {
		return errors.New("media query. order_by cannot be empty")
	}

	if c.HasSortDirection() && c.SortDirection() == "" {
		return errors.New("media query. sort_direction cannot be empty")
	}

	if c.HasOffset() && c.Offset() < 0 {
		return errors.New("media query. offset cannot be negative")
	}

	if c.HasLimit() && c.Limit() < 0 {
		return errors.New("media query. limit cannot be negative")
	}

	return nil
}

func (c *mediaQueryImplementation) Columns() []string {
	if !c.hasProperty("columns") {
		return []string{}
	}

	return c.properties["columns"].([]string)
}

func (c *mediaQueryImplementation) SetColumns(columns []string) MediaQueryInterface {
	c.properties["columns"] = columns

	return c
}

func (c *mediaQueryImplementation) HasCountOnly() bool {
	return c.hasProperty("count_only")
}

func (c *mediaQueryImplementation) IsCountOnly() bool {
	if !c.HasCountOnly() {
		return false
	}

	return c.properties["count_only"].(bool)
}

func (c *mediaQueryImplementation) SetCountOnly(countOnly bool) MediaQueryInterface {
	c.properties["count_only"] = countOnly

	return c
}

func (c *mediaQueryImplementation) HasEntityID() bool {
	return c.hasProperty("entity_id")
}

func (c *mediaQueryImplementation) EntityID() string {
	if !c.HasEntityID() {
		return ""
	}

	return c.properties["entity_id"].(string)
}

func (c *mediaQueryImplementation) SetEntityID(entityID string) MediaQueryInterface {
	c.properties["entity_id"] = entityID

	return c
}

func (c *mediaQueryImplementation) HasID() bool {
	return c.hasProperty("id")
}

func (c *mediaQueryImplementation) ID() string {
	if !c.HasID() {
		return ""
	}

	return c.properties["id"].(string)
}

func (c *mediaQueryImplementation) SetID(id string) MediaQueryInterface {
	c.properties["id"] = id

	return c
}

func (c *mediaQueryImplementation) HasIDIn() bool {
	return c.hasProperty("id_in")
}

func (c *mediaQueryImplementation) IDIn() []string {
	if !c.HasIDIn() {
		return []string{}
	}

	return c.properties["id_in"].([]string)
}

func (c *mediaQueryImplementation) SetIDIn(idIn []string) MediaQueryInterface {
	c.properties["id_in"] = idIn

	return c
}

func (c *mediaQueryImplementation) HasLimit() bool {
	return c.hasProperty("limit")
}

func (c *mediaQueryImplementation) Limit() int {
	if !c.HasLimit() {
		return 0
	}

	return c.properties["limit"].(int)
}

func (c *mediaQueryImplementation) SetLimit(limit int) MediaQueryInterface {
	c.properties["limit"] = limit

	return c
}

func (c *mediaQueryImplementation) HasOffset() bool {
	return c.hasProperty("offset")
}

func (c *mediaQueryImplementation) Offset() int {
	if !c.HasOffset() {
		return 0
	}

	return c.properties["offset"].(int)
}

func (c *mediaQueryImplementation) SetOffset(offset int) MediaQueryInterface {
	c.properties["offset"] = offset

	return c
}

func (c *mediaQueryImplementation) HasOrderBy() bool {
	return c.hasProperty("order_by")
}

func (c *mediaQueryImplementation) OrderBy() string {
	if !c.HasOrderBy() {
		return ""
	}

	return c.properties["order_by"].(string)
}

func (c *mediaQueryImplementation) SetOrderBy(orderBy string) MediaQueryInterface {
	c.properties["order_by"] = orderBy

	return c
}

func (c *mediaQueryImplementation) HasSortDirection() bool {
	return c.hasProperty("sort_direction")
}

func (c *mediaQueryImplementation) SortDirection() string {
	if !c.HasSortDirection() {
		return ""
	}

	return c.properties["sort_direction"].(string)
}

func (c *mediaQueryImplementation) SetSortDirection(sortDirection string) MediaQueryInterface {
	c.properties["sort_direction"] = sortDirection

	return c
}

func (c *mediaQueryImplementation) SetParentID(parentID string) MediaQueryInterface {
	c.properties["parent_id"] = parentID

	return c
}

func (c *mediaQueryImplementation) HasSoftDeletedIncluded() bool {
	return c.hasProperty("soft_deleted_included")
}

func (c *mediaQueryImplementation) SoftDeletedIncluded() bool {
	if !c.HasSoftDeletedIncluded() {
		return false
	}

	return c.properties["soft_deleted_included"].(bool)
}

func (c *mediaQueryImplementation) SetSoftDeletedIncluded(softDeletedIncluded bool) MediaQueryInterface {
	c.properties["soft_deleted_included"] = softDeletedIncluded

	return c
}

func (c *mediaQueryImplementation) HasStatus() bool {
	return c.hasProperty("status")
}

func (c *mediaQueryImplementation) Status() string {
	if !c.HasStatus() {
		return ""
	}

	return c.properties["status"].(string)
}

func (c *mediaQueryImplementation) SetStatus(status string) MediaQueryInterface {
	c.properties["status"] = status

	return c
}

func (c *mediaQueryImplementation) HasTitleLike() bool {
	return c.hasProperty("title_like")
}

func (c *mediaQueryImplementation) TitleLike() string {
	if !c.HasTitleLike() {
		return ""
	}

	return c.properties["title_like"].(string)
}

func (c *mediaQueryImplementation) SetTitleLike(titleLike string) MediaQueryInterface {
	c.properties["title_like"] = titleLike

	return c
}

func (c *mediaQueryImplementation) HasType() bool {
	return c.hasProperty("type")
}

func (c *mediaQueryImplementation) Type() string {
	if !c.HasType() {
		return ""
	}

	return c.properties["type"].(string)
}

func (c *mediaQueryImplementation) SetType(type_ string) MediaQueryInterface {
	c.properties["type"] = type_

	return c
}

func (c *mediaQueryImplementation) hasProperty(name string) bool {
	_, ok := c.properties[name]
	return ok
}
