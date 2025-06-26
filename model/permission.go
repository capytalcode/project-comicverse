package model

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Permissions int64

var (
	_ sql.Scanner  = (*Permissions)(nil)
	_ driver.Value = Permissions(0)
	_ fmt.Stringer = Permissions(0)
)

func (p Permissions) Has(perm ...Permissions) bool {
	// Bitwise AND to compare if p has a permission
	//
	// If for example, p is 0x0010 ("edit.accessibility") and perm is
	// 0x0001 ("read"): 0x0010 AND 0x0001 = 0x0000, which is not equal
	// to 0x0001, return false.
	//
	// If p is 0x0011 ("edit.accessibility" and "read") and perm is
	// 0x0001 ("read"): 0x0011 AND 0x0001 results in 0x0001, which
	// is equal to 0x0001 ("read").
	if len(perm) == 0 {
		return false
	}
	if len(perm) == 1 {
		return p&perm[0] == perm[0]
	}
	for _, pe := range perm {
		if p&pe != pe {
			return false
		}
	}
	return true
}

func (p *Permissions) Add(perm ...Permissions) {
	if p == nil {
		t := Permissions(0)
		p = &t
	}
	// Bitwise OR to add permissions.
	//
	// If p is 0x0001 ("read") and pe is 0x0010 ("edit.accessibility"):
	// 0x0001 OR 0x0010 results in 0x0011, which means we added the "edit.accessibility" bit.
	for _, pe := range perm {
		*p = *p | pe
	}
}

func (p *Permissions) Remove(perm ...Permissions) {
	if p == nil {
		return
	}
	// Bitwise NOT AND
	//
	// If p is 0x0011 ("read" + "edit.accessibility"), and perm is 0x0010 ("edit.accessibility"):
	// we first convert perm to a bit-mask using NOT, so it becomes 0x1101; then we use AND to
	// remove the "edit.accessibility", since 0x0011 AND 0x1101 results in 0x0001 ("read").
	for _, pe := range perm {
		*p = *p & (^pe)
	}
}

func (p *Permissions) Scan(src any) error {
	switch src := src.(type) {
	case nil:
		return nil
	case int64:
		*p = Permissions(src)
	case string:
		if strings.HasPrefix(src, "0x") {
			i, err := strconv.ParseInt(strings.TrimPrefix(src, "0x"), 2, 64)
			if err != nil {
				return errors.Join(errors.New("Scan: unable to scan binary Permissions"), err)
			}
			return p.Scan(i)
		}
		i, err := strconv.ParseInt(src, 10, 64)
		if err != nil {
			return errors.Join(errors.New("Scan: unable to scan base10 Permissions"), err)
		}
		return p.Scan(i)
	case []byte:
		return p.Scan(string(src))
	default:
		return fmt.Errorf("Scan: unable to scan type %T into Permissions", src)
	}

	return nil
}

func (p Permissions) Value() (driver.Value, error) {
	return int64(p), nil
}

func (p Permissions) String() string {
	if p.Has(PermissionAuthor) {
		return "author"
	}

	labels := []string{}
	for perm, l := range PermissionLabels {
		if p.Has(perm) {
			labels = append(labels, l)
		}
	}

	return strings.Join(labels, ",")
}

const (
	PermissionAuthor            Permissions = 0x1111111111111111 // "author"
	PermissionAdminDelete       Permissions = 0x1000000000000000 // "admin.delete" -----
	PermissionAdminAll          Permissions = 0x0111110000000001 // "admin.all"
	PermissionAdminProject      Permissions = 0x0100000000000000 // "admin.project"
	PermissionAdminMembers      Permissions = 0x0010000000000000 // "admin.members"
	PermissionEditAll           Permissions = 0x0000001111111111 // "edit.all" ---------
	PermissionEditPages         Permissions = 0x0000000100000000 // "edit.pages"
	PermissionEditInteractions  Permissions = 0x0000000010000000 // "edit.interactions"
	PermissionEditDialogs       Permissions = 0x0000000000001000 // "edit.dialogs"
	PermissionEditTranslations  Permissions = 0x0000000000000100 // "edit.translations"
	PermissionEditAccessibility Permissions = 0x0000000000000010 // "edit.accessibility"
	PermissionRead              Permissions = 0x0000000000000001 // "read"
)

var PermissionLabels = map[Permissions]string{
	PermissionAuthor:            "author",
	PermissionAdminDelete:       "admin.delete",
	PermissionAdminProject:      "admin.project",
	PermissionAdminMembers:      "admin.members",
	PermissionEditPages:         "edit.pages",
	PermissionEditInteractions:  "edit.interactions",
	PermissionEditDialogs:       "edit.dialogs",
	PermissionEditTranslations:  "edit.translations",
	PermissionEditAccessibility: "edit.accessibility",
	PermissionRead:              "read",
}
