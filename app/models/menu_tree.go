package models

import (
	"sort"
)

func MenuTree(list []Menu) menuSlice {
	data := buildMenuData(list)
	result := makeMenuTreeCore(0, data)
	sort.Stable(result)

	// body, err := json.Marshal(result)
	return result
	// body, err := json.Marshal(result)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// return string(body)
}

type menuSlice []Menu

func (s menuSlice) Len() int           { return len(s) }
func (s menuSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s menuSlice) Less(i, j int) bool { return s[i].Position < s[j].Position }

func buildMenuData(list menuSlice) map[uint64]map[uint64]Menu {
	var data map[uint64]map[uint64]Menu = make(map[uint64]map[uint64]Menu)
	for _, v := range list {
		id := v.ID
		fid := v.ParentID
		if _, ok := data[fid]; !ok {
			data[fid] = make(map[uint64]Menu)
		}
		data[fid][id] = v
	}
	return data
}

func makeMenuTreeCore(index uint64, data map[uint64]map[uint64]Menu) menuSlice {

	tmp := make(menuSlice, 0)
	for id, item := range data[index] {

		if data[id] != nil {
			item.Children = makeMenuTreeCore(id, data)
		}

		tmp = append(tmp, item)
	}
	sort.Stable(tmp)
	return tmp
}
