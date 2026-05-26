package util

import "go-pure-admin/server/internal/model"

func BuildMenuTree(items []model.SysMenu, parentID *int64) []model.SysMenu {
	var roots []model.SysMenu
	for i := range items {
		if sameParent(items[i].ParentID, parentID) {
			node := items[i]
			node.Children = BuildMenuTree(items, &node.ID)
			roots = append(roots, node)
		}
	}
	return roots
}

func BuildDeptTree(items []model.SysDept, parentID *int64) []model.SysDept {
	var roots []model.SysDept
	for i := range items {
		if sameParent(items[i].ParentID, parentID) {
			node := items[i]
			node.Children = BuildDeptTree(items, &node.ID)
			roots = append(roots, node)
		}
	}
	return roots
}

func sameParent(a, b *int64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// FilterMenusForRoute 返回可注册动态路由的目录/页面（含 hidden；排除按钮与已禁用）。
func FilterMenusForRoute(items []model.SysMenu) []model.SysMenu {
	out := make([]model.SysMenu, 0, len(items))
	for _, m := range items {
		if m.MenuType == 2 || !m.Enabled {
			continue
		}
		out = append(out, m)
	}
	return out
}

func CollectPermCodes(items []model.SysMenu) []string {
	var codes []string
	for _, m := range items {
		if m.MenuType == 2 && m.PermCode != "" {
			codes = append(codes, m.PermCode)
		}
	}
	return codes
}
