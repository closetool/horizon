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

package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/horizoncd/horizon/core/common"
	herrors "github.com/horizoncd/horizon/core/errors"
	"github.com/horizoncd/horizon/lib/orm"
	userauth "github.com/horizoncd/horizon/pkg/authentication/user"
	applicationdao "github.com/horizoncd/horizon/pkg/dao"
	perror "github.com/horizoncd/horizon/pkg/errors"
	"github.com/horizoncd/horizon/pkg/models"
	appmodels "github.com/horizoncd/horizon/pkg/models"
	"github.com/horizoncd/horizon/pkg/server/global"
	callbacks "github.com/horizoncd/horizon/pkg/util/ormcallbacks"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	// use tmp sqlite
	notExistID = uint(100)
)

func createGroupCtx() (*gorm.DB, context.Context, GroupManager, RegionManager,
	EnvironmentManager, EnvironmentRegionManager, TagManager) {
	var (
		db, _        = orm.NewSqliteDB("")
		ctx          context.Context
		groupMgr     = NewGroupManager(db)
		regionMgr    = NewRegionManager(db)
		envMgr       = NewEnvironmentManager(db)
		envregionMgr = NewEnvironmentRegionManager(db)
		tagMgr       = NewTagManager(db)
	)
	// nolint
	ctx = context.WithValue(context.Background(), common.UserContextKey(), &userauth.DefaultInfo{
		Name: "tony",
		ID:   110,
	})
	callbacks.RegisterCustomCallbacks(db)
	// create table
	err := db.AutoMigrate(&appmodels.Group{})
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
	err = db.AutoMigrate(&appmodels.Application{})
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
	err = db.AutoMigrate(&appmodels.Member{})
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
	err = db.AutoMigrate(&appmodels.EnvironmentRegion{})
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
	err = db.AutoMigrate(&appmodels.Region{})
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
	err = db.AutoMigrate(&appmodels.Environment{})
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
	err = db.AutoMigrate(&appmodels.Tag{})
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
	return db, ctx, groupMgr, regionMgr, envMgr, envregionMgr, tagMgr
}

func TestUint(t *testing.T) {
	g := appmodels.Group{
		ParentID: 0,
	}

	_, err := json.Marshal(g)
	assert.Nil(t, err)
}

func getGroup(parentID uint, name, path string) *appmodels.Group {
	return &appmodels.Group{
		Name:            name,
		Path:            path,
		VisibilityLevel: "private",
		ParentID:        parentID,
		CreatedBy:       1,
		UpdatedBy:       1,
	}
}

func TestCreate(t *testing.T) {
	db, ctx, groupMgr, _, _, _, _ := createGroupCtx()
	// normal create, parentID is nil
	g1, err := groupMgr.Create(ctx, getGroup(0, "1", "a"))
	assert.Nil(t, err)
	get, _ := groupMgr.GetByID(ctx, g1.ID)
	assert.Equal(t, fmt.Sprintf("%d", g1.ID), get.TraversalIDs)

	// name conflict, parentID is nil
	_, err = groupMgr.Create(ctx, getGroup(0, "1", "b"))
	assert.Equal(t, herrors.ErrNameConflict, perror.Cause(err))

	// path conflict, with parentID is nil
	_, err = groupMgr.Create(ctx, getGroup(0, "2", "a"))
	assert.Equal(t, herrors.ErrPathConflict, perror.Cause(err))

	// name conflict with application
	name := "app"
	_, err = applicationdao.NewApplicationDAO(db).Create(ctx, &appmodels.Application{
		Name: name,
	}, nil)
	assert.Nil(t, err)
	_, err = groupMgr.Create(ctx, getGroup(0, name, "a"))
	assert.Equal(t, perror.Cause(err), herrors.ErrGroupConflictWithApplication)

	// normal create, parentID: not nil
	group2 := getGroup(g1.ID, "2", "b")
	g2, err := groupMgr.Create(ctx, group2)
	assert.Nil(t, err)
	get, _ = groupMgr.GetByID(ctx, g2.ID)
	assert.Equal(t, fmt.Sprintf("%d,%d", g1.ID, g2.ID), get.TraversalIDs)

	// drop table
	res := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&appmodels.Group{})
	assert.Nil(t, res.Error)
}

func TestDelete(t *testing.T) {
	db, ctx, groupMgr, _, _, _, _ := createGroupCtx()
	g1, err := groupMgr.Create(ctx, getGroup(0, "1", "a"))
	assert.Nil(t, err)

	// delete exist record
	_, err = groupMgr.Delete(ctx, g1.ID)
	assert.Nil(t, err)

	_, err = groupMgr.GetByID(ctx, g1.ID)
	_, ok := perror.Cause(err).(*herrors.HorizonErrNotFound)
	assert.True(t, ok)

	// delete not exist record
	var count int64
	count, err = groupMgr.Delete(ctx, notExistID)
	assert.Equal(t, 0, int(count))
	assert.Nil(t, err)

	// drop table
	res := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&appmodels.Group{})
	assert.Nil(t, res.Error)
}

func TestGetByID(t *testing.T) {
	db, ctx, groupMgr, _, _, _, _ := createGroupCtx()
	g1, err := groupMgr.Create(ctx, getGroup(0, "1", "a"))
	assert.Nil(t, err)

	// query exist record
	group1, err := groupMgr.GetByID(ctx, g1.ID)
	assert.Nil(t, err)
	assert.NotNil(t, group1.ID)

	// query not exist record
	_, err = groupMgr.GetByID(ctx, notExistID)
	_, ok := perror.Cause(err).(*herrors.HorizonErrNotFound)
	assert.True(t, ok)

	// drop table
	res := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&appmodels.Group{})
	assert.Nil(t, res.Error)
}

func TestGetByIDs(t *testing.T) {
	db, ctx, groupMgr, _, _, _, _ := createGroupCtx()
	g1, err := groupMgr.Create(ctx, getGroup(0, "1", "a"))
	assert.Nil(t, err)
	g2, _ := groupMgr.Create(ctx, getGroup(0, "2", "b"))

	groups, err := groupMgr.GetByIDs(ctx, []uint{g1.ID, g2.ID})
	assert.Nil(t, err)
	assert.Equal(t, g1.ID, groups[0].ID)
	assert.Equal(t, g2.ID, groups[1].ID)

	// drop table
	res := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&appmodels.Group{})
	assert.Nil(t, res.Error)
}

func TestGetAll(t *testing.T) {
	db, ctx, groupMgr, _, _, _, _ := createGroupCtx()
	g1, err := groupMgr.Create(ctx, getGroup(0, "1", "a"))
	assert.Nil(t, err)
	g2, err := groupMgr.Create(ctx, getGroup(0, "2", "b"))
	assert.Nil(t, err)

	groups, err := groupMgr.GetAll(ctx)
	assert.Nil(t, err)
	assert.Equal(t, g1.ID, groups[0].ID)
	assert.Equal(t, g2.ID, groups[1].ID)

	// drop table
	res := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&appmodels.Group{})
	assert.Nil(t, res.Error)
}

func TestGetByPaths(t *testing.T) {
	db, ctx, groupMgr, _, _, _, _ := createGroupCtx()
	id, err := groupMgr.Create(ctx, getGroup(0, "1", "a"))
	assert.Nil(t, err)
	id2, _ := groupMgr.Create(ctx, getGroup(0, "2", "b"))

	groups, err := groupMgr.GetByPaths(ctx, []string{"a", "b"})
	assert.Nil(t, err)
	assert.Equal(t, id.ID, groups[0].ID)
	assert.Equal(t, id2.ID, groups[1].ID)

	// test GetByNameOrPathUnderParent
	groups, err = groupMgr.GetByNameOrPathUnderParent(ctx, "1", "b", 0)
	assert.Nil(t, err)
	assert.Equal(t, len(groups), 2)
	assert.Equal(t, groups[0].Path, "a")
	assert.Equal(t, groups[1].Name, "2")

	// drop table
	res := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&appmodels.Group{})
	assert.Nil(t, res.Error)
}

func TestGetByNameFuzzily(t *testing.T) {
	db, ctx, groupMgr, _, _, _, _ := createGroupCtx()
	id, err := groupMgr.Create(ctx, getGroup(0, "1", "a"))
	assert.Nil(t, err)
	id2, _ := groupMgr.Create(ctx, getGroup(0, "21", "b"))

	groups, err := groupMgr.GetByNameFuzzily(ctx, "1")
	assert.Nil(t, err)
	assert.Equal(t, id.ID, groups[0].ID)
	assert.Equal(t, id2.ID, groups[1].ID)

	// drop table
	res := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&appmodels.Group{})
	assert.Nil(t, res.Error)
}

func TestUpdateBasic(t *testing.T) {
	db, ctx, groupMgr, _, _, _, _ := createGroupCtx()
	group1 := getGroup(0, "1", "a")
	g1, err := groupMgr.Create(ctx, group1)
	assert.Nil(t, err)

	// update exist record
	group1.ID = g1.ID
	group1.Name = "update1"
	err = groupMgr.UpdateBasic(ctx, group1)
	assert.Nil(t, err)
	group, err := groupMgr.GetByID(ctx, g1.ID)
	assert.Nil(t, err)
	assert.Equal(t, "update1", group.Name)

	// update fail because of conflict
	group2 := getGroup(0, "2", "b")
	g2, err := groupMgr.Create(ctx, group2)
	assert.Nil(t, err)
	group2.ID = g2.ID
	group2.Name = "update1"
	err = groupMgr.UpdateBasic(ctx, group2)
	assert.Equal(t, herrors.ErrNameConflict, perror.Cause(err))

	// update regionSelector
	err = groupMgr.UpdateRegionSelector(ctx, g1.ID, "XXX")
	assert.Nil(t, err)
	group, _ = groupMgr.GetByID(ctx, g1.ID)
	assert.Equal(t, group.RegionSelector, "XXX")

	// drop table
	res := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&appmodels.Group{})
	assert.Nil(t, res.Error)
}

func TestTransferGroup(t *testing.T) {
	_, ctx, groupMgr, _, _, _, _ := createGroupCtx()
	g1, err := groupMgr.Create(ctx, getGroup(0, "1", "a"))
	assert.Nil(t, err)
	g2, err := groupMgr.Create(ctx, getGroup(g1.ID, "2", "b"))
	assert.Nil(t, err)
	g3, err := groupMgr.Create(ctx, getGroup(0, "3", "c"))
	assert.Nil(t, err)
	_, err = groupMgr.Create(ctx, getGroup(g3.ID, "2", "d"))
	assert.Nil(t, err)

	// not valid transfer: name conflict
	err = groupMgr.Transfer(ctx, g2.ID, g3.ID)
	assert.True(t, perror.Cause(err) == herrors.ErrNameConflict)

	// valid transfer
	err = groupMgr.Transfer(ctx, g1.ID, g3.ID)
	assert.Nil(t, err)

	group, err := groupMgr.GetByID(ctx, g2.ID)
	assert.Nil(t, err)

	expect := []string{
		strconv.Itoa(int(g3.ID)),
		strconv.Itoa(int(g1.ID)),
		strconv.Itoa(int(g2.ID)),
	}
	assert.Equal(t, strings.Join(expect, ","), group.TraversalIDs)
}

func TestManagerGetChildren(t *testing.T) {
	db, ctx, groupMgr, _, _, _, _ := createGroupCtx()
	g0, err := groupMgr.Create(ctx, getGroup(0, "0", "0"))
	assert.Nil(t, err)
	g1, err := groupMgr.Create(ctx, getGroup(g0.ID, "1", "a"))
	assert.Nil(t, err)
	g2, err := groupMgr.Create(ctx, getGroup(g0.ID, "2", "b"))
	assert.Nil(t, err)
	a1, err := applicationdao.NewApplicationDAO(db).Create(ctx, &appmodels.Application{
		Name:    "3",
		GroupID: g0.ID,
	}, nil)
	assert.Nil(t, err)

	type args struct {
		parentID   uint
		pageNumber int
		pageSize   int
	}
	tests := []struct {
		name    string
		args    args
		want    []*appmodels.GroupOrApplication
		want1   int64
		wantErr bool
	}{
		{
			name: "firstPage",
			args: args{
				parentID:   g0.ID,
				pageNumber: 1,
				pageSize:   2,
			},
			want: []*appmodels.GroupOrApplication{
				{
					Model: global.Model{
						ID:        g2.ID,
						UpdatedAt: g2.UpdatedAt,
					},
					Name:        "2",
					Path:        "b",
					Type:        "group",
					Description: "",
				},
				{
					Model: global.Model{
						ID:        g1.ID,
						UpdatedAt: g1.UpdatedAt,
					},
					Name:        "1",
					Path:        "a",
					Type:        "group",
					Description: "",
				},
			},
			want1: 3,
		},
		{
			name: "secondPage",
			args: args{
				parentID:   g0.ID,
				pageNumber: 2,
				pageSize:   2,
			},
			want: []*appmodels.GroupOrApplication{
				{
					Model: global.Model{
						ID:        a1.ID,
						UpdatedAt: a1.UpdatedAt,
					},
					Name:        "3",
					Path:        "3",
					Type:        "application",
					Description: "",
				},
			},
			want1: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := groupMgr.GetChildren(ctx, tt.args.parentID, tt.args.pageNumber, tt.args.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChildren() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i, val := range got {
				val.UpdatedAt = tt.want[i].UpdatedAt
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChildren() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetChildren() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetSubGroupsByGroupIDs(t *testing.T) {
	_, ctx, groupMgr, _, _, _, _ := createGroupCtx()
	g1, err := groupMgr.Create(ctx, getGroup(0, "a", "a"))
	assert.Nil(t, err)
	get, _ := groupMgr.GetByID(ctx, g1.ID)
	assert.Equal(t, fmt.Sprintf("%d", g1.ID), get.TraversalIDs)

	g2, err := groupMgr.Create(ctx, getGroup(0, "b", "b"))
	assert.Nil(t, err)
	get2, _ := groupMgr.GetByID(ctx, g2.ID)
	assert.Equal(t, fmt.Sprintf("%d", g2.ID), get2.TraversalIDs)

	g3, err := groupMgr.Create(ctx, getGroup(g1.ID, "c", "c"))
	assert.Nil(t, err)
	get3, _ := groupMgr.GetByID(ctx, g3.ID)
	assert.Equal(t, fmt.Sprintf("%d,%d", g1.ID, g3.ID), get3.TraversalIDs)

	g4, err := groupMgr.Create(ctx, getGroup(g2.ID, "c", "c"))
	assert.Nil(t, err)
	get4, _ := groupMgr.GetByID(ctx, g4.ID)
	assert.Equal(t, fmt.Sprintf("%d,%d", g2.ID, g4.ID), get4.TraversalIDs)

	ids := []uint{g1.ID, g2.ID}
	groups, err := groupMgr.GetSubGroupsByGroupIDs(ctx, ids)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(groups))
	for _, group := range groups {
		t.Logf("group: %v", group)
	}

	ids2 := []uint{g2.ID}
	groups2, err := groupMgr.GetSubGroupsByGroupIDs(ctx, ids2)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(groups2))
	for _, group := range groups2 {
		t.Logf("group: %v", group)
	}

	ids3 := []uint{g3.ID}
	groups3, err := groupMgr.GetSubGroupsByGroupIDs(ctx, ids3)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(groups3))
	assert.Equal(t, g3.ID, groups3[0].ID)
	for _, group := range groups3 {
		t.Logf("group: %v", group)
	}
}

func Test_manager_GetSelectableRegionsByEnv(t *testing.T) {
	_, ctx, groupMgr, regionMgr, envMgr, envregionMgr, tagMgr := createGroupCtx()
	// initializing data
	r1, _ := regionMgr.Create(ctx, &appmodels.Region{
		Name:        "hz",
		DisplayName: "HZ",
	})
	_, _ = regionMgr.Create(ctx, &appmodels.Region{
		Name:        "hz-disabled",
		DisplayName: "HZ",
		Disabled:    true,
	})
	r3, _ := regionMgr.Create(ctx, &appmodels.Region{
		Name:        "hz3",
		DisplayName: "HZ",
	})
	devEnv, _ := envMgr.CreateEnvironment(ctx, &appmodels.Environment{
		Name:        "dev",
		DisplayName: "开发",
	})
	_, _ = envregionMgr.CreateEnvironmentRegion(ctx, &appmodels.EnvironmentRegion{
		EnvironmentName: devEnv.Name,
		RegionName:      "hz",
		IsDefault:       true,
	})
	_, _ = envregionMgr.CreateEnvironmentRegion(ctx, &appmodels.EnvironmentRegion{
		EnvironmentName: devEnv.Name,
		RegionName:      "hz-disabled",
	})
	_, _ = envregionMgr.CreateEnvironmentRegion(ctx, &appmodels.EnvironmentRegion{
		EnvironmentName: devEnv.Name,
		RegionName:      "hz3",
	})
	_ = tagMgr.UpsertByResourceTypeID(ctx, common.ResourceRegion, r1.ID, []*models.TagBasic{
		{
			Key:   "a",
			Value: "1",
		}, {
			Key:   "b",
			Value: "1",
		},
	})
	_ = tagMgr.UpsertByResourceTypeID(ctx, common.ResourceRegion, r3.ID, []*models.TagBasic{
		{
			Key:   "a",
			Value: "1",
		}, {
			Key:   "c",
			Value: "1",
		},
	})
	g1, err := groupMgr.Create(ctx, &appmodels.Group{
		Name: "11",
		Path: "pp",
		RegionSelector: `- key: "a"
  operator: "in"
  values: 
    - "1"
- key: "b"
  operator: "in"
  values: 
    - "1"
`,
	})
	assert.Nil(t, err)
	// get regionSelector from parent group
	g2, _ := groupMgr.Create(ctx, &appmodels.Group{
		Name:     "22",
		Path:     "p2",
		ParentID: g1.ID,
	})

	type args struct {
		id  uint
		env string
	}
	tests := []struct {
		name    string
		args    args
		want    appmodels.RegionParts
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				id:  g2.ID,
				env: "dev",
			},
			want: appmodels.RegionParts{
				{
					Name:        "hz",
					DisplayName: "HZ",
					IsDefault:   true,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := groupMgr.GetSelectableRegionsByEnv(ctx, tt.args.id, tt.args.env)
			if (err != nil) != tt.wantErr {
				t.Errorf(fmt.Sprintf("GetSelectableRegionsByEnv(%v, %v, %v)", ctx, tt.args.id, tt.args.env))
			}
			assert.Equalf(t, tt.want, got, "GetSelectableRegionsByEnv(%v, %v, %v)", ctx, tt.args.id, tt.args.env)
		})
	}

	defaultRegions, err := groupMgr.GetDefaultRegions(ctx, g2.ID)
	assert.Nil(t, err)
	assert.Equal(t, len(defaultRegions), 1)
	assert.Equal(t, defaultRegions[0].RegionName, "hz")
	assert.Equal(t, defaultRegions[0].EnvironmentName, "dev")
}

func Test_manager_GetSelectableRegions(t *testing.T) {
	_, ctx, groupMgr, regionMgr, _, _, tagMgr := createGroupCtx()
	r1, _ := regionMgr.Create(ctx, &appmodels.Region{
		Name:        "hz",
		DisplayName: "HZ",
	})
	_, _ = regionMgr.Create(ctx, &appmodels.Region{
		Name:        "hz-disabled",
		DisplayName: "HZ",
		Disabled:    true,
	})
	r3, _ := regionMgr.Create(ctx, &appmodels.Region{
		Name:        "hz3",
		DisplayName: "HZ",
	})

	_ = tagMgr.UpsertByResourceTypeID(ctx, common.ResourceRegion, r1.ID, []*models.TagBasic{
		{
			Key:   "a",
			Value: "11",
		},
	})
	_ = tagMgr.UpsertByResourceTypeID(ctx, common.ResourceRegion, r3.ID, []*models.TagBasic{
		{
			Key:   "a",
			Value: "11",
		},
	})

	g1, err := groupMgr.Create(ctx, &appmodels.Group{
		Name: "112",
		Path: "pp2",
		RegionSelector: `- key: "a"
  operator: "in"
  values: 
    - "11"
`,
	})
	assert.Nil(t, err)
	type args struct {
		id uint
	}
	tests := []struct {
		name    string
		args    args
		want    appmodels.RegionParts
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				id: g1.ID,
			},
			want: appmodels.RegionParts{
				{
					Name:        "hz",
					DisplayName: "HZ",
				},
				{
					Name:        "hz3",
					DisplayName: "HZ",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := groupMgr.GetSelectableRegions(ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf(fmt.Sprintf("GetSelectableRegions(%v, %v)", ctx, tt.args.id))
			}
			assert.Equalf(t, tt.want, got, "GetSelectableRegions(%v, %v)", tt.args.id)
		})
	}
}
