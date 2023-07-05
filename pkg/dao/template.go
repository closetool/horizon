// Copyright © 2023 Horizoncd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dao

import (
	"context"
	"fmt"

	"github.com/horizoncd/horizon/core/common"
	herrors "github.com/horizoncd/horizon/core/errors"
	"github.com/horizoncd/horizon/lib/q"
	dbsql "github.com/horizoncd/horizon/pkg/common"
	hctx "github.com/horizoncd/horizon/pkg/context"
	perror "github.com/horizoncd/horizon/pkg/errors"
	amodels "github.com/horizoncd/horizon/pkg/models"
	"gorm.io/gorm"
)

type TemplateDAO interface {
	Create(ctx context.Context, template *amodels.Template) (*amodels.Template, error)
	ListTemplate(ctx context.Context) ([]*amodels.Template, error)
	ListByGroupID(ctx context.Context, groupID uint) ([]*amodels.Template, error)
	DeleteByID(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*amodels.Template, error)
	GetByName(ctx context.Context, name string) (*amodels.Template, error)
	GetRefOfApplication(ctx context.Context, id uint) ([]*amodels.Application, uint, error)
	GetRefOfCluster(ctx context.Context, id uint) ([]*amodels.Cluster, uint, error)
	UpdateByID(ctx context.Context, id uint, template *amodels.Template) error
	ListByGroupIDs(ctx context.Context, ids []uint) ([]*amodels.Template, error)
	ListByIDs(ctx context.Context, ids []uint) ([]*amodels.Template, error)
	ListV2(ctx context.Context, query *q.Query, gorupIDs ...uint) ([]*amodels.Template, error)
}

// NewTemplateDAO returns an instance of the default TemplateDAO
func NewTemplateDAO(db *gorm.DB) TemplateDAO {
	return &templateDAO{db: db}
}

type templateDAO struct{ db *gorm.DB }

func (d templateDAO) Create(ctx context.Context, template *amodels.Template) (*amodels.Template, error) {
	result := d.db.WithContext(ctx).Create(template)
	return template, result.Error
}

func (d templateDAO) ListTemplate(ctx context.Context) ([]*amodels.Template, error) {
	var templates []*amodels.Template
	result := d.db.Raw(dbsql.TemplateList).Scan(&templates)
	if result.Error != nil {
		return nil, result.Error
	}
	return templates, nil
}

func (d templateDAO) ListByGroupID(ctx context.Context, groupID uint) ([]*amodels.Template, error) {
	var templates []*amodels.Template
	result := d.db.Raw(dbsql.TemplateListByGroup, groupID).Scan(&templates)
	if result.Error != nil {
		return nil, result.Error
	}
	return templates, nil
}

func (d templateDAO) DeleteByID(ctx context.Context, id uint) error {
	if res := d.db.Exec(dbsql.TemplateDelete, id); res.Error != nil {
		return perror.Wrap(herrors.NewErrDeleteFailed(herrors.TemplateInDB, res.Error.Error()),
			fmt.Sprintf("failed to delete template, id = %d", id))
	}
	return nil
}

func (d templateDAO) GetByID(ctx context.Context, id uint) (*amodels.Template, error) {
	var template amodels.Template
	res := d.db.Raw(dbsql.TemplateQueryByID, id).First(&template)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, perror.Wrap(herrors.NewErrNotFound(herrors.TemplateInDB, res.Error.Error()),
				fmt.Sprintf("failed to find template: id = %d", id))
		}
		return nil, perror.Wrap(herrors.NewErrGetFailed(herrors.TemplateInDB, res.Error.Error()),
			fmt.Sprintf("failed to get template: id = %d", id))
	}
	return &template, nil
}

func (d templateDAO) GetByName(ctx context.Context, name string) (*amodels.Template, error) {
	var template amodels.Template
	res := d.db.Raw(dbsql.TemplateQueryByName, name).First(&template)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return nil, perror.Wrap(herrors.NewErrNotFound(herrors.TemplateInDB, res.Error.Error()),
				fmt.Sprintf("failed to find template: name = %s", name))
		}
		return nil, perror.Wrap(herrors.NewErrGetFailed(herrors.TemplateInDB, res.Error.Error()),
			fmt.Sprintf("failed to get template: name = %s", name))
	}
	return &template, nil
}

func (d templateDAO) GetRefOfApplication(ctx context.Context, id uint) ([]*amodels.Application, uint, error) {
	onlyRefCount, ok := ctx.Value(hctx.TemplateOnlyRefCount).(bool)
	var (
		applications []*amodels.Application
		total        uint
	)
	res := d.db.Raw(dbsql.TemplateRefCountOfApplication, id).Scan(&total)
	if res.Error != nil {
		return nil, 0, perror.Wrap(herrors.NewErrGetFailed(herrors.TemplateInDB, res.Error.Error()),
			fmt.Sprintf("failed to get ref count of application: %s", res.Error.Error()))
	}

	if !ok || !onlyRefCount {
		res = d.db.Raw(dbsql.TemplateRefOfApplication, id).Scan(&applications)
		if res.Error != nil {
			return nil, 0, perror.Wrap(herrors.NewErrGetFailed(herrors.TemplateInDB, res.Error.Error()),
				fmt.Sprintf("failed to get ref of application: %s", res.Error.Error()))
		}
	}
	return applications, total, nil
}

func (d templateDAO) GetRefOfCluster(ctx context.Context, id uint) ([]*amodels.Cluster, uint, error) {
	onlyRefCount, ok := ctx.Value(hctx.TemplateOnlyRefCount).(bool)
	var (
		clusters []*amodels.Cluster
		total    uint
	)
	res := d.db.Raw(dbsql.TemplateRefCountOfCluster, id).Scan(&total)
	if res.Error != nil {
		return nil, 0, perror.Wrap(herrors.NewErrGetFailed(herrors.TemplateInDB, res.Error.Error()),
			fmt.Sprintf("failed to get ref count of cluster: %s", res.Error.Error()))
	}

	if !ok || !onlyRefCount {
		res = d.db.Raw(dbsql.TemplateRefOfCluster, id).Scan(&clusters)
		if res.Error != nil {
			return nil, 0, perror.Wrap(herrors.NewErrGetFailed(herrors.TemplateInDB, res.Error.Error()),
				fmt.Sprintf("failed to get ref of cluster: %s", res.Error.Error()))
		}
	}
	return clusters, total, nil
}

func (d templateDAO) UpdateByID(ctx context.Context, templateID uint, template *amodels.Template) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		var oldTemplate amodels.Template
		res := tx.Raw(dbsql.TemplateQueryByID, templateID).Scan(&oldTemplate)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return perror.Wrap(herrors.NewErrNotFound(herrors.TemplateInDB, res.Error.Error()),
				fmt.Sprintf("not found template with templateID = %d", templateID))
		}

		oldTemplate.UpdatedBy = template.UpdatedBy
		if template.Repository != "" {
			oldTemplate.Repository = template.Repository
		}
		if template.Description != "" {
			oldTemplate.Description = template.Description
		}
		if template.OnlyOwner != nil {
			oldTemplate.OnlyOwner = template.OnlyOwner
		}
		return tx.Model(&oldTemplate).Updates(oldTemplate).Error
	})
}

func (d templateDAO) ListByGroupIDs(ctx context.Context, ids []uint) ([]*amodels.Template, error) {
	templates := make([]*amodels.Template, 0)
	res := d.db.Where("group_id in ?", ids).Find(&templates)
	if res.Error != nil {
		return nil, perror.Wrapf(herrors.NewErrGetFailed(herrors.TemplateInDB, res.Error.Error()),
			"failed to get template:\n"+
				"template ids = %v\n err = %v", ids, res.Error)
	}
	return templates, nil
}

func (d templateDAO) ListByIDs(ctx context.Context, ids []uint) ([]*amodels.Template, error) {
	templates := make([]*amodels.Template, 0)
	res := d.db.Where("id in ?", ids).Find(&templates)
	if res.Error != nil {
		return nil, perror.Wrapf(herrors.NewErrGetFailed(herrors.TemplateInDB, res.Error.Error()),
			"failed to get template:\n"+
				"template ids = %v\n err = %v", ids, res.Error)
	}
	return templates, nil
}

func (d templateDAO) ListV2(ctx context.Context, query *q.Query, groupIDs ...uint) ([]*amodels.Template, error) {
	var templates []*amodels.Template

	statement := d.db.WithContext(ctx).
		Table("tb_template as t").
		Select("t.*")

	genSQL := func(statement *gorm.DB, query *q.Query) *gorm.DB {
		for k, v := range query.Keywords {
			switch k {
			case common.TemplateQueryWithoutCI:
				statement = statement.Where("t.without_ci = ?", v)
			case common.TemplateQueryByGroup:
				statement = statement.Where("t.group_id = ?", v)
			case common.TemplateQueryName:
				statement = statement.Where("t.name like ?", fmt.Sprintf("%%%v%%", v))
			case common.TemplateQueryByGroups:
				if _, ok := v.(uint); ok {
					statement = statement.Where("t.group_id = ?", v)
				} else if _, ok = v.([]uint); ok {
					statement = statement.Where("t.group_id in ?", v)
				}
			case common.TemplateQueryByUser:
				statement = statement.
					Joins("join tb_member as m on m.resource_id = t.id").
					Where("m.resource_type = ?", common.ResourceTemplate).
					Where("m.member_type = '0'").
					Where("m.deleted_ts = 0").
					Where("m.membername_id = ?", v)
			}
		}
		statement = statement.Where("t.deleted_ts = 0")
		return statement
	}

	if query != nil {
		statement = genSQL(statement, query)

		if len(groupIDs) > 0 &&
			query.Keywords != nil &&
			query.Keywords[common.TemplateQueryByUser] != nil {
			statementGroup := d.db.WithContext(ctx).
				Table("tb_template as t").
				Select("t.*")

			delete(query.Keywords, common.TemplateQueryByUser)
			statementGroup = genSQL(statementGroup, query)

			statementGroup = statementGroup.Where("group_id in ?", groupIDs)
			statement = d.db.Raw("? union ?", statement, statementGroup)
		}
	}

	statement = d.db.Raw("select distinct * from (?) as apps order by updated_at desc", statement)
	res := statement.Scan(&templates)
	if res.Error != nil {
		return nil, herrors.NewErrGetFailed(herrors.TemplateInDB, res.Error.Error())
	}

	return templates, nil
}
