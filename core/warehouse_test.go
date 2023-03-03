package warehouse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type elementGetterMock struct {
	getElements func() (*Element, error)
}

func (m *elementGetterMock) GetElements() (*Element, error) {
	return m.getElements()
}

func Test_CreateList(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		getElements := func() (*Element, error) {
			return &Element{
				Children: map[string]*Element{
					"category_1": {
						Children: map[string]*Element{
							"category_2": {
								Children: map[string]*Element{
									"item_1": {
										Item: true,
									},
								},
							},
							"category_3": {
								Children: map[string]*Element{
									"item_2": {
										Item: true,
									},
								},
							},
						},
					},
				},
			}, nil
		}
		elemGetter := elementGetterMock{getElements: getElements}
		wh, err := New(&elemGetter)
		require.NoError(t, err)
		elem, err := wh.CreateList()
		require.NoError(t, err)

		require.Len(t, elem.Children, 1)

		category1 := elem.Children["category_1"]
		require.False(t, category1.Item)
		require.Len(t, category1.Children, 2)

		category2 := category1.Children["category_2"]
		require.False(t, category2.Item)
		require.Len(t, category2.Children, 1)

		item1 := category2.Children["item_1"]
		require.True(t, item1.Item)
		require.Nil(t, item1.Children)

		category3 := category1.Children["category_3"]
		require.False(t, category3.Item)
		require.Len(t, category3.Children, 1)

		item2 := category3.Children["item_2"]
		require.True(t, item2.Item)
		require.Nil(t, item2.Children)
	})

	t.Run("hierarchy error", func(t *testing.T) {
		getElements := func() (*Element, error) {
			return nil, ErrCategoryHiearchy{ErrorMsg: "error message"}
		}
		elemGetter := elementGetterMock{getElements: getElements}
		wh, err := New(&elemGetter)
		require.NoError(t, err)
		_, err = wh.CreateList()
		require.Error(t, err)
		require.EqualError(t, err, "could not get elements: error message")
	})
}
