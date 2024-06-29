// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package models

import (
	"database/sql/driver"
	"fmt"
)

type Alignment string

const (
	AlignmentGOOD    Alignment = "GOOD"
	AlignmentNEUTRAL Alignment = "NEUTRAL"
	AlignmentEVIL    Alignment = "EVIL"
)

func (e *Alignment) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Alignment(s)
	case string:
		*e = Alignment(s)
	default:
		return fmt.Errorf("unsupported scan type for Alignment: %T", src)
	}
	return nil
}

type NullAlignment struct {
	Alignment Alignment `json:"alignment"`
	Valid     bool      `json:"valid"` // Valid is true if Alignment is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAlignment) Scan(value interface{}) error {
	if value == nil {
		ns.Alignment, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Alignment.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAlignment) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Alignment), nil
}

type Rarity string

const (
	RarityCOMMON       Rarity = "COMMON"
	RarityUNCOMMON     Rarity = "UNCOMMON"
	RarityRARE         Rarity = "RARE"
	RarityEPIC         Rarity = "EPIC"
	RarityLEGENDARY    Rarity = "LEGENDARY"
	RarityMYTHICAL     Rarity = "MYTHICAL"
	RarityROLESPECIFIC Rarity = "ROLE_SPECIFIC"
)

func (e *Rarity) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Rarity(s)
	case string:
		*e = Rarity(s)
	default:
		return fmt.Errorf("unsupported scan type for Rarity: %T", src)
	}
	return nil
}

type NullRarity struct {
	Rarity Rarity `json:"rarity"`
	Valid  bool   `json:"valid"` // Valid is true if Rarity is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRarity) Scan(value interface{}) error {
	if value == nil {
		ns.Rarity, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Rarity.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRarity) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Rarity), nil
}

type AbilityCategory struct {
	AbilityID  int32 `json:"ability_id"`
	CategoryID int32 `json:"category_id"`
}

type AbilityInfo struct {
	ID             int32  `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	DefaultCharges int32  `json:"default_charges"`
	AnyAbility     bool   `json:"any_ability"`
	Rarity         Rarity `json:"rarity"`
}

type Category struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type PerkInfo struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Role struct {
	ID          int32     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Alignment   Alignment `json:"alignment"`
}

type RoleAbility struct {
	RoleID    int32 `json:"role_id"`
	AbilityID int32 `json:"ability_id"`
}

type RolePerk struct {
	RoleID int32 `json:"role_id"`
	PerkID int32 `json:"perk_id"`
}
