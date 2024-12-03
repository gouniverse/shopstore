package shopstore

type ProductQueryOptions struct {
	ID           string
	IDIn         []string
	Status       string
	StatusIn     []string
	Title        string
	CreatedAtGte string
	CreatedAtLte string
	Offset       int
	Limit        int
	SortOrder    string
	OrderBy      string
	CountOnly    bool
	WithDeleted  bool
}
