package db

import warehouse "github.com/diegocmsantos/warehouse/core"

func CategoryHiearchy(line []string) (bool, error) {
	for i, l := range line {
		if l == "" && i != len(line)-2 {
			return false, warehouse.ErrCategoryHiearchy{
				ErrorMsg: "line cannot have categories after a blank space",
			}
		}
	}

	return true, nil
}
