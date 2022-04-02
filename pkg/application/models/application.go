package models

import (
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type Priority string

const (
	P0 Priority = "P0"
	P1 Priority = "P1"
	P2 Priority = "P2"
	P3 Priority = "P3"
)

type Application struct {
	gorm.Model
	GroupID         uint
	Name            string
	Description     string
	Priority        Priority
	GitURL          string
	GitSubfolder    string
	GitBranch       string
	Template        string
	TemplateRelease string
	CreatedBy       uint
	UpdatedBy       uint
	DeletedTs       soft_delete.DeletedAt
}
