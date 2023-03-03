package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Get(t *testing.T) {
	t.Run("get successfully", func(t *testing.T) {
		c, err := New("../testdata/test_file.csv")
		require.NoError(t, err)
		elem, err := c.GetElements()
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

	t.Run("levels that contain empty strings should be interpreted as the end of that hierarchy branch", func(t *testing.T) {
		c, err := New("../testdata/specialcase.csv")
		require.NoError(t, err)
		elem, err := c.GetElements()
		require.NoError(t, err)

		require.Len(t, elem.Children, 2)

		category1 := elem.Children["category_1"]
		require.False(t, category1.Item)
		require.Len(t, category1.Children, 1)

		item1 := category1.Children["item_1"]
		require.True(t, item1.Item)
		require.Nil(t, item1.Children)

		category2 := elem.Children["category_2"]
		require.False(t, category2.Item)
		require.Len(t, category2.Children, 1)

		category3 := category2.Children["category_3"]
		require.False(t, category3.Item)
		require.Len(t, category3.Children, 1)

		item2 := category3.Children["item_2"]
		require.True(t, item2.Item)
		require.Nil(t, item2.Children)
	})
}
