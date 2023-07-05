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

	"github.com/horizoncd/horizon/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/horizoncd/horizon/core/common"
	herrors "github.com/horizoncd/horizon/core/errors"
	"github.com/horizoncd/horizon/lib/q"
	perror "github.com/horizoncd/horizon/pkg/errors"
)

type EventDAO interface {
	CreateEvent(ctx context.Context, event ...*models.Event) ([]*models.Event, error)
	List(ctx context.Context, query *q.Query) ([]*models.Event, error)
	CreateOrUpdateCursor(ctx context.Context,
		eventIndex *models.EventCursor) (*models.EventCursor, error)
	GetCursor(ctx context.Context) (*models.EventCursor, error)
	GetEvent(ctx context.Context, id uint) (*models.Event, error)
	DeleteEvents(ctx context.Context, id ...uint) (int64, error)
}

type eventDAO struct{ db *gorm.DB }

// NewEventDAO returns an instance of the default EventDAO
func NewEventDAO(db *gorm.DB) EventDAO {
	return &eventDAO{db: db}
}

func (d *eventDAO) CreateEvent(ctx context.Context, events ...*models.Event) ([]*models.Event, error) {
	if result := d.db.WithContext(ctx).Create(events); result.Error != nil {
		return nil, herrors.NewErrInsertFailed(herrors.EventInDB, result.Error.Error())
	}
	return events, nil
}

func (d *eventDAO) List(ctx context.Context, query *q.Query) ([]*models.Event, error) {
	var events []*models.Event
	statement := d.db.WithContext(ctx).Debug().Order("id asc")
	for k, v := range query.Keywords {
		switch k {
		case common.Offset:
			offset, ok := v.(int)
			if !ok {
				return nil, perror.Wrap(herrors.ErrParamInvalid, fmt.Sprintf("invalid offset %v", v))
			}
			statement = statement.Offset(offset)
		case common.Limit:
			limit, ok := v.(int)
			if !ok {
				return nil, perror.Wrap(herrors.ErrParamInvalid, fmt.Sprintf("invalid limit %v", v))
			}
			statement = statement.Limit(limit)
		case common.StartID:
			statement = statement.Where("id > ?", v)
		case common.EndID:
			statement = statement.Where("id <= ?", v)
		case common.ReqID:
			statement = statement.Where("req_id = ?", v)
		}
	}

	if result := statement.Find(&events); result.Error != nil {
		return nil, herrors.NewErrInsertFailed(herrors.EventInDB, result.Error.Error())
	}
	return events, nil
}

func (d *eventDAO) GetEvent(ctx context.Context, id uint) (*models.Event, error) {
	var event *models.Event
	if result := d.db.WithContext(ctx).Where("id = ?", id).First(&event); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, herrors.NewErrNotFound(herrors.EventInDB,
				result.Error.Error())
		}
		return nil, herrors.NewErrGetFailed(herrors.EventInDB, result.Error.Error())
	}
	return event, nil
}

func (d *eventDAO) CreateOrUpdateCursor(ctx context.Context,
	eventCursor *models.EventCursor) (*models.EventCursor, error) {
	if result := d.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "id",
			},
		},
		DoUpdates: clause.AssignmentColumns([]string{"position"}),
	}).Create(eventCursor); result.Error != nil {
		return nil, herrors.NewErrInsertFailed(herrors.EventInDB, result.Error.Error())
	}
	return eventCursor, nil
}

func (d *eventDAO) GetCursor(ctx context.Context) (*models.EventCursor, error) {
	var eventIndex models.EventCursor
	statement := d.db.WithContext(ctx)
	if result := statement.First(&eventIndex); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, herrors.NewErrNotFound(herrors.EventCursorInDB,
				result.Error.Error())
		}
		return nil, herrors.NewErrGetFailed(herrors.EventCursorInDB, result.Error.Error())
	}
	return &eventIndex, nil
}

func (d *eventDAO) DeleteEvents(ctx context.Context, ids ...uint) (int64, error) {
	var events []*models.Event
	tx := d.db.WithContext(ctx).Begin()

	result := tx.Where("id in (?)", ids).Delete(&events)
	if result.Error != nil {
		tx.Rollback()
		return 0, herrors.NewErrDeleteFailed(herrors.EventInDB, result.Error.Error())
	}

	if result := tx.Where("event_id in (?)", ids).Delete(&models.WebhookLog{}); result.Error != nil {
		tx.Rollback()
		return 0, herrors.NewErrDeleteFailed(herrors.WebhookLogInDB, result.Error.Error())
	}

	tx.Commit()
	return result.RowsAffected, nil
}

// TODO: must add gc
