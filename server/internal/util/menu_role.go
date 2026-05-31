package util

import (
	"pure-go-admin/server/internal/model"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

// ExpandRoleMenuIDs 补全祖先目录，并在 autoIncludeButtons 为 true 时为已选页面补全其下按钮节点。
func ExpandRoleMenuIDs(db *gorm.DB, ids []int64, autoIncludeButtons bool) ([]int64, error) {
	if len(ids) == 0 {
		return ids, nil
	}

	var all []model.SysMenu
	if err := db.Select("id", "parent_id", "menu_type").Find(&all).Error; err != nil {
		return nil, err
	}

	byID := make(map[int64]model.SysMenu, len(all))
	children := make(map[int64][]int64)
	for _, m := range all {
		byID[m.ID] = m
		if m.ParentID != nil {
			children[*m.ParentID] = append(children[*m.ParentID], m.ID)
		}
	}

	set := make(map[int64]struct{}, len(ids)*2)
	add := func(id int64) {
		if id > 0 {
			set[id] = struct{}{}
		}
	}

	for _, id := range ids {
		add(id)
		cur := id
		for {
			m, ok := byID[cur]
			if !ok || m.ParentID == nil {
				break
			}
			add(*m.ParentID)
			cur = *m.ParentID
		}
	}

	if autoIncludeButtons {
		// 对已选页面补全按钮；已单独勾选的按钮保持不变
		snapshot := lo.Keys(set)
		for _, id := range snapshot {
			m, ok := byID[id]
			if !ok || m.MenuType != 1 {
				continue
			}
			for _, cid := range children[id] {
				if child, ok := byID[cid]; ok && child.MenuType == 2 {
					add(cid)
				}
			}
		}
	}

	return lo.Keys(set), nil
}
