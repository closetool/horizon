package cluster

import (
	"github.com/horizoncd/horizon/core/common"
	appmodels "github.com/horizoncd/horizon/pkg/application/models"
	"github.com/horizoncd/horizon/pkg/cluster/code"
	"github.com/horizoncd/horizon/pkg/cluster/models"
	tagmodels "github.com/horizoncd/horizon/pkg/tag/models"
)

type Cluster struct {
	*models.Cluster `json:",inline"`
	*TemplateInput  `json:"templateInput,omitempty"`
	Tags            tagmodels.TagsBasic `json:"tags,omitempty"`
}

type Clusterv2 struct {
	*models.Cluster `json:",inline"`
	*TemplateInput  `json:"templateInput,omitempty"`
	TemplateConfig  map[string]interface{} `json:"templateConfig,omitempty"`
	TemplateInfo    *code.TemplateInfo     `json:"-"`
	MergePatch      bool                   `json:"mergePatch"`
	BuildConfig     map[string]interface{} `json:"buildConfig"`
	Tags            tagmodels.TagsBasic    `json:"tags,omitempty"`
}

func (c *Clusterv2) toClusterModel(application *appmodels.Application) *models.Cluster {
	cluster := &models.Cluster{
		ApplicationID:   c.ApplicationID,
		Name:            c.Name,
		EnvironmentName: c.EnvironmentName,
		RegionName:      c.RegionName,
		Description:     c.Description,
		ExpireSeconds:   c.ExpireSeconds,
		GitURL:          c.GitURL,
		GitSubfolder:    c.GitSubfolder,
		GitRef:          c.GitRef,
		GitRefType:      c.GitRefType,
		Image:           c.Image,
		Template:        c.Template,
		TemplateRelease: c.TemplateRelease,
		Status:          common.ClusterStatusCreating,
	}
	if c.Template == application.Template {
		if c.GitURL == "" {
			c.GitURL = application.GitURL
		}
		if c.GitSubfolder == "" {
			c.GitSubfolder = application.GitSubfolder
		}
		if c.GitRef == "" {
			c.GitRef = application.GitRef
		}
		if c.GitRefType == "" {
			c.GitRefType = application.GitRefType
		}
		if c.Image == "" {
			c.Image = application.Image
		}
	}
	return cluster
}
