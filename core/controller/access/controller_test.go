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

package access

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/horizoncd/horizon/core/common"
	"github.com/horizoncd/horizon/core/middleware"
	"github.com/horizoncd/horizon/lib/orm"
	userauth "github.com/horizoncd/horizon/pkg/authentication/user"
	applicationmodels "github.com/horizoncd/horizon/pkg/models"
	"github.com/horizoncd/horizon/pkg/param/managerparam"
	"github.com/horizoncd/horizon/pkg/rbac"
	roleservice "github.com/horizoncd/horizon/pkg/rbac/role"
	memberservice "github.com/horizoncd/horizon/pkg/service"
	"github.com/stretchr/testify/assert"
)

var (
	ctx         context.Context
	c           Controller
	group       *applicationmodels.Group
	application *applicationmodels.Application
	cluster     *applicationmodels.Cluster
	manager     *managerparam.Manager
)

// nolint
func TestMain(m *testing.M) {
	db, _ := orm.NewSqliteDB("")
	manager = managerparam.InitManager(db)
	if err := db.AutoMigrate(&applicationmodels.Member{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&applicationmodels.Group{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&applicationmodels.Application{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&applicationmodels.Cluster{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&applicationmodels.User{}); err != nil {
		panic(err)
	}

	roleService, err := roleservice.NewFileRole(context.Background(), strings.NewReader(roleConfig))
	if err != nil {
		panic(err)
	}
	memberService := memberservice.NewMemberService(roleService, nil, manager)
	if err != nil {
		panic(err)
	}
	ctx = context.WithValue(context.Background(), common.UserContextKey(), &userauth.DefaultInfo{
		Name: "Tony",
		ID:   uint(110),
	})

	rbacAuthorizer := rbac.NewAuthorizer(roleService, memberService)
	skippers := middleware.MethodAndPathSkipper("*",
		regexp.MustCompile("(^/apis/front/.*)|(^/health)|(^/metrics)|(^/apis/login)|"+
			"(^/apis/core/v1/roles)|(^/apis/internal/.*)"))
	c = NewController(rbacAuthorizer, skippers)

	group, err = manager.GroupManager.Create(ctx, &applicationmodels.Group{
		Name:            "group",
		Path:            "/group",
		VisibilityLevel: "private",
	})
	if err != nil {
		panic(err)
	}
	application, _ = manager.ApplicationManager.Create(ctx, &applicationmodels.Application{
		Name:    "application",
		GroupID: group.ID,
	}, nil)

	cluster, _ = manager.ClusterMgr.Create(ctx, &applicationmodels.Cluster{
		Name:          "cluster",
		ApplicationID: application.ID,
	}, nil, nil)

	os.Exit(m.Run())
}

// nolint
func TestController_GetAccesses_Guest(t *testing.T) {
	guest, err := manager.UserManager.Create(ctx, &applicationmodels.User{
		Name: "guest",
	})

	nonMemberCtx := context.WithValue(ctx, common.UserContextKey(), &userauth.DefaultInfo{
		ID: 2,
	})
	guestCtx := context.WithValue(ctx, common.UserContextKey(), &userauth.DefaultInfo{
		ID: guest.ID,
	})

	deniedAPIs := []API{
		{
			URL:    fmt.Sprintf("/apis/core/v1/clusters/%d/shell", cluster.ID),
			Method: "GET",
		},
		{
			URL:    fmt.Sprintf("/apis/core/v1/applications/%d/clusters?scope=dev/hz", application.ID),
			Method: "POST",
		},
	}

	allowAPIs := []API{
		{
			URL:    fmt.Sprintf("/apis/core/v1/clusters/%d/status", cluster.ID),
			Method: "GET",
		},
		{
			URL:    fmt.Sprintf("/apis/core/v1/applications/%d/clusters?scope=test/hz", application.ID),
			Method: "POST",
		},
	}

	apis := append(deniedAPIs, allowAPIs...)

	reviewResults, err := c.Review(nonMemberCtx, apis)
	assert.Nil(t, err)
	for _, api := range deniedAPIs {
		assert.Equal(t, false, reviewResults[api.URL][api.Method].Allowed)
	}
	for _, api := range allowAPIs {
		assert.Equal(t, true, reviewResults[api.URL][api.Method].Allowed)
	}

	reviewResults, err = c.Review(guestCtx, apis)
	assert.Nil(t, err)
	for _, api := range deniedAPIs {
		assert.Equal(t, false, reviewResults[api.URL][api.Method].Allowed)
	}
	for _, api := range allowAPIs {
		assert.Equal(t, true, reviewResults[api.URL][api.Method].Allowed)
	}
}

// nolint
func TestController_GetAccesses_Owner(t *testing.T) {
	owner, err := manager.UserManager.Create(ctx, &applicationmodels.User{
		Name: "owner",
	})

	ctx := context.WithValue(ctx, common.UserContextKey(), &userauth.DefaultInfo{
		ID: owner.ID,
	})

	_, err = manager.MemberManager.Create(ctx, &applicationmodels.Member{
		ResourceType: "groups",
		ResourceID:   group.ID,
		Role:         "owner",
		MemberType:   applicationmodels.MemberUser,
		MemberNameID: owner.ID,
	})
	assert.Nil(t, err)

	deniedAPIs := []API{
		{
			URL:    fmt.Sprintf("/apis/core/v1/clusters/%d/templateschematags", cluster.ID),
			Method: "POST",
		},
	}

	allowAPIs := []API{
		{
			URL:    fmt.Sprintf("/apis/core/v1/clusters/%d/shell", cluster.ID),
			Method: "GET",
		},
		{
			URL:    fmt.Sprintf("/apis/core/v1/clusters/%d", cluster.ID),
			Method: "DELETE",
		},
	}
	apis := append(deniedAPIs, allowAPIs...)

	reviewResults, err := c.Review(ctx, apis)
	assert.Nil(t, err)
	for _, api := range deniedAPIs {
		fmt.Println(reviewResults[api.URL][api.Method])
		assert.Equal(t, false, reviewResults[api.URL][api.Method].Allowed)
	}
	for _, api := range allowAPIs {
		fmt.Println(reviewResults[api.URL][api.Method])
		assert.Equal(t, true, reviewResults[api.URL][api.Method].Allowed)
	}
}

// nolint
func TestController_GetAccesses_Admin(t *testing.T) {
	admin, err := manager.UserManager.Create(ctx, &applicationmodels.User{
		Name: "admin",
	})

	ctx := context.WithValue(ctx, common.UserContextKey(), &userauth.DefaultInfo{
		ID:    admin.ID,
		Admin: true,
	})

	apis := []API{
		{
			URL:    fmt.Sprintf("/apis/core/v1/clusters/%d/shell", cluster.ID),
			Method: "GET",
		},
		{
			URL:    fmt.Sprintf("/apis/core/v1/clusters/%d", cluster.ID),
			Method: "DELETE",
		},
		{
			URL:    "/apis/core/v1/groups",
			Method: "POST",
		},
	}

	reviewResults, err := c.Review(ctx, apis)
	assert.Nil(t, err)
	for _, api := range apis {
		fmt.Println(reviewResults[api.URL][api.Method])
		assert.Equal(t, true, reviewResults[api.URL][api.Method].Allowed)
	}
}

const roleConfig = `RolePriorityRankDesc:
  - pe
  - owner
  - maintainer
  - guest
DefaultRole: guest
Roles:
  - name: owner
    desc: 'owner为组/应用/集群的拥有者,拥有最高权限'
    rules:
      - apiGroups:
          - core
        resources:
          - applications
          - groups/applications
          - applications/members
          - applications/envtemplates
        verbs:
          - '*'
        scopes:
          - '*'
        nonResourceURLs:
          - '*'
      - apiGroups:
          - core
        resources:
          - groups
          - groups/members
          - groups/groups
        verbs:
          - '*'
        scopes:
          - '*'
      - apiGroups:
          - core
        resources:
          - applications/clusters
          - clusters
          - clusters/builddeploy
          - clusters/deploy
          - clusters/diffs
          - clusters/next
          - clusters/restart
          - clusters/rollback
          - clusters/status
          - clusters/members
          - clusters/pipelineruns
          - clusters/terminal
          - clusters/containerlog
          - clusters/online
          - clusters/offline
          - clusters/tags
          - pipelineruns
          - pipelineruns/stop
          - pipelineruns/log
          - pipelineruns/diffs
          - clusters/dashboards
          - clusters/pods
          - clusters/free
          - clusters/events
          - clusters/outputs
          - clusters/promote
          - clusters/shell
        verbs:
          - '*'
        scopes:
          - '*'
        nonResourceURLs:
          - '*'
      - apiGroups:
          - core
        resources:
          - clusters/templateschematags
        verbs:
          - get
        scopes:
          - '*'
        nonResourceURLs:
          - '*'
  - name: maintainer
    desc: maintainer为组/应用/集群的管理者，拥有除删除资源之外的其他权限，并且也可以进行成员管理
    rules:
      - apiGroups:
          - core
        resources:
          - applications
          - groups/applications
          - applications/members
          - applications/envtemplates
        verbs:
          - create
          - get
          - update
        scopes:
          - '*'
        nonResourceURLs:
          - '*'
      - apiGroups:
          - core
        resources:
          - groups
          - groups/members
          - groups/groups
        verbs:
          - get
          - create
          - update
        scopes:
          - '*'
      - apiGroups:
          - core
        resources:
          - applications/clusters
          - clusters
          - clusters/builddeploy
          - clusters/deploy
          - clusters/diffs
          - clusters/next
          - clusters/restart
          - clusters/rollback
          - clusters/status
          - clusters/members
          - clusters/pipelineruns
          - clusters/terminal
          - clusters/containerlog
          - clusters/online
          - clusters/offline
          - clusters/tags
          - pipelineruns
          - pipelineruns/stop
          - pipelineruns/log
          - pipelineruns/diffs
          - clusters/dashboards
          - clusters/pods
          - clusters/free
          - clusters/events
          - clusters/outputs
          - clusters/promote
          - clusters/shell
        verbs:
          - create
          - get
          - update
        scopes:
          - '*'
        nonResourceURLs:
          - '*'
      - apiGroups:
          - core
        resources:
          - clusters/templateschematags
        verbs:
          - get
        scopes:
          - '*'
        nonResourceURLs:
          - '*'
  - name: pe
    desc: pe为应用/集群的管理者，拥有除删除资源之外的其他权限，并且也可以进行成员管理。破格修改资源上限等
    rules:
      - apiGroups:
          - core
        resources:
          - applications
          - groups/applications
          - applications/members
          - applications/envtemplates
        verbs:
          - create
          - get
          - update
        scopes:
          - '*'
        nonResourceURLs:
          - '*'
      - apiGroups:
          - core
        resources:
          - groups
          - groups/members
          - groups/groups
        verbs:
          - get
          - create
          - update
        scopes:
          - '*'
      - apiGroups:
          - core
        resources:
          - applications/clusters
          - clusters
          - clusters/builddeploy
          - clusters/deploy
          - clusters/diffs
          - clusters/next
          - clusters/restart
          - clusters/rollback
          - clusters/status
          - clusters/members
          - clusters/pipelineruns
          - clusters/terminal
          - clusters/containerlog
          - clusters/online
          - clusters/offline
          - clusters/tags
          - pipelineruns
          - pipelineruns/stop
          - pipelineruns/log
          - pipelineruns/diffs
          - clusters/dashboards
          - clusters/pods
          - clusters/free
          - clusters/templateschematags
          - clusters/events
          - clusters/outputs
          - clusters/promote
          - clusters/shell
        verbs:
          - create
          - get
          - update
        scopes:
          - '*'
        nonResourceURLs:
          - '*'
  - name: guest
    desc: guest为只读人员，拥有组/应用/项目的只读权限，以及测试环境集群创建的权限
    rules:
      - apiGroups:
          - core
        resources:
          - groups
          - groups/members
          - groups/groups
          - applications
          - groups/applications
          - applications/clusters
          - applications/members
          - applications/envtemplates
          - clusters
          - clusters/diffs
          - clusters/status
          - clusters/members
          - clusters/pipelineruns
          - clusters/containerlog
          - clusters/tags
          - pipelineruns
          - pipelineruns/log
          - pipelineruns/diffs
          - clusters/dashboards
          - clusters/pods
          - clusters/events
          - clusters/outputs
          - clusters/templateschematags
        verbs:
          - get
        scopes:
          - '*'
      - apiGroups:
          - core
        resources:
          - applications/clusters
        verbs:
          - create
          - get
          - update
        scopes:
          - test/*
          - reg/*
          - perf/*
          - pre/*`
