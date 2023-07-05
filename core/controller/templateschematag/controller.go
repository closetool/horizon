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

package templateschematag

import (
	"context"

	"github.com/horizoncd/horizon/core/common"
	clustermanager "github.com/horizoncd/horizon/pkg/manager"
	"github.com/horizoncd/horizon/pkg/param"

	"github.com/horizoncd/horizon/pkg/util/wlog"
)

type Controller interface {
	List(ctx context.Context, clusterID uint) (*ListResponse, error)
	Update(ctx context.Context, clusterID uint, r *UpdateRequest) error
}

type controller struct {
	clusterMgr          clustermanager.ClusterManager
	clusterSchemaTagMgr clustermanager.TemplateSchemaTagManager
}

func NewController(param *param.Param) Controller {
	return &controller{
		clusterMgr:          param.ClusterMgr,
		clusterSchemaTagMgr: param.ClusterSchemaTagMgr,
	}
}

func (c *controller) List(ctx context.Context, clusterID uint) (_ *ListResponse, err error) {
	const op = "cluster template scheme tag controller: list"
	defer wlog.Start(ctx, op).StopPrint()

	tags, err := c.clusterSchemaTagMgr.ListByClusterID(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	return ofClusterTemplateSchemaTags(tags), nil
}

func (c *controller) Update(ctx context.Context, clusterID uint, r *UpdateRequest) (err error) {
	const op = "cluster template scheme tag controller: update"
	defer wlog.Start(ctx, op).StopPrint()

	currentUser, err := common.UserFromContext(ctx)
	if err != nil {
		return err
	}

	clusterTemplateSchemaTags := r.toClusterTemplateSchemaTags(clusterID, currentUser)

	if err := clustermanager.ValidateUpsert(clusterTemplateSchemaTags); err != nil {
		return err
	}

	cluster, err := c.clusterMgr.GetByID(ctx, clusterID)
	if err != nil {
		return err
	}

	return c.clusterSchemaTagMgr.UpsertByClusterID(ctx, cluster.ID, clusterTemplateSchemaTags)
}
