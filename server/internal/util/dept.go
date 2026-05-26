package util

import (
	"fmt"
	"strconv"
	"strings"

	"go-pure-admin/server/internal/model"
)

func BuildAncestors(parent *model.SysDept) string {
	if parent == nil {
		return "0"
	}
	if parent.Ancestors == "" {
		return fmt.Sprintf("0,%d", parent.ID)
	}
	return fmt.Sprintf("%s,%d", parent.Ancestors, parent.ID)
}

func ParseAncestorIDs(ancestors string) []int64 {
	parts := strings.Split(ancestors, ",")
	ids := make([]int64, 0, len(parts))
	for _, p := range parts {
		if p == "" || p == "0" {
			continue
		}
		id, err := strconv.ParseInt(p, 10, 64)
		if err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}
