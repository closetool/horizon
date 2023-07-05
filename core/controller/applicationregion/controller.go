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

package applicationregion

import (
	"context"

	perror "github.com/horizoncd/horizon/pkg/errors"
	"github.com/horizoncd/horizon/pkg/manager"
	"github.com/horizoncd/horizon/pkg/models"
	"github.com/horizoncd/horizon/pkg/param"
)

type Controller interface {
	List(ctx context.Context, applicationID uint) (ApplicationRegion, error)
	Update(ctx context.Context, applicationID uint, regions ApplicationRegion) error
}

type controller struct {
	mgr                  manager.ApplicationRegionManager
	regionMgr            manager.RegionManager
	environmentMgr       manager.EnvironmentManager
	environmentRegionMgr manager.EnvironmentRegionManager
	groupMgr             manager.GroupManager
	applicationMgr       manager.ApplicationManager
}

var _ Controller = (*controller)(nil)

func NewController(param *param.Param) Controller {
	return &controller{
		mgr:                  param.ApplicationRegionManager,
		regionMgr:            param.RegionMgr,
		environmentMgr:       param.EnvMgr,
		environmentRegionMgr: param.EnvironmentRegionMgr,
		groupMgr:             param.GroupManager,
		applicationMgr:       param.ApplicationManager,
	}
}

func (c *controller) List(ctx context.Context, applicationID uint) (ApplicationRegion, error) {
	applicationRegions, err := c.mgr.ListByApplicationID(ctx, applicationID)
	if err != nil {
		return nil, perror.WithMessage(err, "failed to list application regions")
	}

	environments, err := c.environmentMgr.ListAllEnvironment(ctx)
	if err != nil {
		return nil, perror.WithMessage(err, "failed to list environment")
	}

	application, err := c.applicationMgr.GetByID(ctx, applicationID)
	if err != nil {
		return nil, err
	}
	environmentRegions, err := c.groupMgr.GetDefaultRegions(ctx, application.GroupID)
	if err != nil {
		return nil, perror.WithMessage(err, "failed to list environmentRegions")
	}

	return ofApplicationRegion(applicationRegions, environments, environmentRegions), nil
}

func (c *controller) Update(ctx context.Context, applicationID uint, regions ApplicationRegion) error {
	applicationRegions := make([]*models.ApplicationRegion, 0)

	for _, r := range regions {
		if r.Environment != "" && r.Region != "" {
			_, err := c.environmentRegionMgr.GetByEnvironmentAndRegion(ctx, r.Environment, r.Region)
			if err != nil {
				return perror.WithMessagef(err,
					"environment/region %s/%s is not exists", r.Environment, r.Region)
			}
			applicationRegions = append(applicationRegions, &models.ApplicationRegion{
				ApplicationID:   applicationID,
				EnvironmentName: r.Environment,
				RegionName:      r.Region,
			})
		}
	}

	return c.mgr.UpsertByApplicationID(ctx, applicationID, applicationRegions)
}
