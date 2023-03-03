package warehouse

import (
	"errors"
	"fmt"
)

// ErrCategoryHiearchy thrown when there is an empty category and it's not the end of the branch
type ErrCategoryHiearchy struct {
	ErrorMsg string `json:"error"`
}

func (e ErrCategoryHiearchy) Error() string {
	return e.ErrorMsg
}

type Element struct {
	Item     bool                `json:"item,omitempty"`
	Children map[string]*Element `json:"children,omitempty"`
}

type ElementGetter interface {
	GetElements() (*Element, error)
}

type Warehouse struct {
	DB ElementGetter
}

func New(db ElementGetter) (*Warehouse, error) {

	if db == nil {
		return nil, errors.New("db is required")
	}

	return &Warehouse{DB: db}, nil
}

// CreateList creates a list of categories and items
// Returns
// * a list of categories and items
// Possible Errors
// * ErrCategoryHiearchy when a category is empty and the following is not empty
func (h *Warehouse) CreateList() (*Element, error) {
	t, err := h.DB.GetElements()
	if err != nil {
		return nil, fmt.Errorf("could not get elements: %w", err)
	}
	return t, err
}
