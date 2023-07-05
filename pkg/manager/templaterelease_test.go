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
	"os"
	"testing"

	"github.com/horizoncd/horizon/core/common"
	herrors "github.com/horizoncd/horizon/core/errors"
	"github.com/horizoncd/horizon/lib/orm"
	userauth "github.com/horizoncd/horizon/pkg/authentication/user"
	perror "github.com/horizoncd/horizon/pkg/errors"
	applicationmodel "github.com/horizoncd/horizon/pkg/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func createTemplateReleaseCtx() (context.Context, *gorm.DB,
	TemplateManager, TemplateReleaseManager, ApplicationManager) {
	var (
		db                 *gorm.DB
		ctx                context.Context
		templateMgr        TemplateManager
		templateReleaseMgr TemplateReleaseManager
		applicationMgr     ApplicationManager
	)

	db, _ = orm.NewSqliteDB("")
	if err := db.AutoMigrate(&applicationmodel.TemplateRelease{},
		&applicationmodel.Application{}, &applicationmodel.Template{},
		&applicationmodel.Member{}); err != nil {
		panic(err)
	}
	ctx = context.TODO()
	// nolint
	ctx = common.WithContext(ctx, &userauth.DefaultInfo{
		ID:   1,
		Name: "Jerry",
	})

	templateMgr = NewTemplateManager(db)
	templateReleaseMgr = NewTemplateReleaseManager(db)
	applicationMgr = NewApplicationManager(db)

	return ctx, db, templateMgr, templateReleaseMgr, applicationMgr
}

func TestTemplateRelease(t *testing.T) {
	ctx, _, templateMgr, templateReleaseMgr, applicationMgr := createTemplateReleaseCtx()
	var (
		templateName = "javaapp"
		name         = "v1.0.0"
		chartVersion = "v1.0.0-test"
		repo         = "repo"
		description  = "javaapp template v1.0.0"
		groupID      = uint(0)
		createdBy    = uint(1)
		updatedBy    = uint(1)
		err          error
	)
	template := &applicationmodel.Template{
		Name:        templateName,
		Description: description,
		Repository:  repo,
		GroupID:     groupID,
		CreatedBy:   createdBy,
		UpdatedBy:   updatedBy,
	}
	template, err = templateMgr.Create(ctx, template)
	assert.Nil(t, err)

	recommend := true

	templateRelease := &applicationmodel.TemplateRelease{
		Template:     template.ID,
		TemplateName: templateName,
		Name:         name,
		ChartVersion: chartVersion,
		Description:  description,
		Recommended:  &recommend,
		CreatedBy:    createdBy,
		UpdatedBy:    updatedBy,
	}
	templateRelease, err = templateReleaseMgr.Create(ctx, templateRelease)
	assert.Nil(t, err)

	assert.Equal(t, name, templateRelease.Name)
	assert.Equal(t, description, templateRelease.Description)
	assert.Equal(t, 1, int(templateRelease.ID))

	b, err := json.Marshal(templateRelease)
	assert.Nil(t, err)
	t.Logf(string(b))

	releases, err := templateReleaseMgr.ListByTemplateName(ctx, templateName)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(releases))
	assert.Equal(t, name, releases[0].Name)
	assert.Equal(t, chartVersion, releases[0].ChartVersion)
	assert.Equal(t, description, releases[0].Description)
	assert.Equal(t, 1, int(releases[0].ID))

	releases, err = templateReleaseMgr.ListByTemplateID(ctx, template.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(releases))
	assert.Equal(t, name, releases[0].Name)
	assert.Equal(t, chartVersion, releases[0].ChartVersion)
	assert.Equal(t, description, releases[0].Description)
	assert.Equal(t, 1, int(releases[0].ID))

	// template release not exists
	templateRelease, err = templateReleaseMgr.GetByTemplateNameAndRelease(ctx, templateName, "not-exist")
	assert.NotNil(t, err)
	_, ok := perror.Cause(err).(*herrors.HorizonErrNotFound)
	assert.True(t, ok)
	assert.Nil(t, templateRelease)

	templateRelease, err = templateReleaseMgr.GetByTemplateNameAndRelease(ctx, templateName, name)
	assert.Nil(t, err)
	assert.NotNil(t, templateRelease)
	assert.Equal(t, chartVersion, templateRelease.ChartVersion)
	assert.Equal(t, name, templateRelease.Name)

	app := &applicationmodel.Application{
		Template:        templateName,
		TemplateRelease: templateRelease.Name,
		Name:            "test",
	}
	_, err = applicationMgr.Create(ctx, app, map[string]string{})
	assert.Nil(t, err)

	apps, _, err := templateReleaseMgr.GetRefOfApplication(ctx, 1)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(apps))
	assert.Equal(t, app.Name, apps[0].Name)

	err = templateReleaseMgr.DeleteByID(ctx, templateRelease.ID)
	assert.Nil(t, err)

	templateRelease, err = templateReleaseMgr.GetByTemplateNameAndRelease(ctx, templateName, name)
	assert.NotNil(t, err)
	assert.Nil(t, templateRelease)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
