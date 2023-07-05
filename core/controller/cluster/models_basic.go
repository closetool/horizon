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

package cluster

import (
	"time"

	codemodels "github.com/horizoncd/horizon/pkg/cluster/code"
	"github.com/horizoncd/horizon/pkg/models"
	appmodels "github.com/horizoncd/horizon/pkg/models"
)

type Base struct {
	Description   string             `json:"description"`
	Git           *codemodels.Git    `json:"git"`
	Template      *Template          `json:"template"`
	TemplateInput *TemplateInput     `json:"templateInput"`
	Tags          []*models.TagBasic `json:"tags"`
}

type TemplateInput struct {
	Application map[string]interface{} `json:"application"`
	Pipeline    map[string]interface{} `json:"pipeline"`
}

type CreateClusterRequest struct {
	*Base

	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	ExpireTime string `json:"expireTime"`
	// TODO(gjq): remove these two params after migration
	Image        string            `json:"image"`
	ExtraMembers map[string]string `json:"extraMembers"`
}

type UpdateClusterRequest struct {
	*Base
	Environment string `json:"environment"`
	Region      string `json:"region"`
	ExpireTime  string `json:"expireTime"`
}

type GetClusterResponse struct {
	*CreateClusterRequest

	ID                   uint         `json:"id"`
	FullPath             string       `json:"fullPath"`
	Application          *Application `json:"application"`
	Priority             string       `json:"priority"`
	Template             *Template    `json:"template"`
	Scope                *Scope       `json:"scope"`
	LatestDeployedCommit string       `json:"latestDeployedCommit,omitempty"`
	Status               string       `json:"status,omitempty"`
	CreatedAt            time.Time    `json:"createdAt"`
	UpdatedAt            time.Time    `json:"updatedAt"`
	TTLInSeconds         *uint        `json:"ttlInSeconds"`
	CreatedBy            *User        `json:"createdBy,omitempty"`
	UpdatedBy            *User        `json:"updatedBy,omitempty"`
}

type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Application struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Template struct {
	Name    string `json:"name"`
	Release string `json:"release"`
}

type Scope struct {
	Environment       string `json:"environment"`
	Region            string `json:"region"`
	RegionDisplayName string `json:"regionDisplayName,omitempty"`
}

func (r *CreateClusterRequest) toClusterModel(application *appmodels.Application,
	er *appmodels.EnvironmentRegion, expireSeconds uint) (*appmodels.Cluster, []*appmodels.Tag) {
	var (
		// r.Git cannot be nil
		gitURL       = r.Git.URL
		gitSubfolder = r.Git.Subfolder
	)
	// if gitURL or gitSubfolder is empty, use application's gitURL or gitSubfolder
	if gitURL == "" {
		gitURL = application.GitURL
	}
	if gitSubfolder == "" {
		gitSubfolder = application.GitSubfolder
	}
	cluster := &appmodels.Cluster{
		ApplicationID:   application.ID,
		Name:            r.Name,
		EnvironmentName: er.EnvironmentName,
		RegionName:      er.RegionName,
		Description:     r.Description,
		ExpireSeconds:   expireSeconds,
		GitURL:          gitURL,
		GitSubfolder:    gitSubfolder,
		GitRef:          r.Git.Ref(),
		GitRefType:      r.Git.RefType(),
		Template:        r.Template.Name,
		TemplateRelease: r.Template.Release,
	}
	tags := make([]*appmodels.Tag, 0)
	for _, tag := range r.Tags {
		tags = append(tags, &appmodels.Tag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}
	return cluster, tags
}

func (r *UpdateClusterRequest) toClusterModel(cluster *models.Cluster,
	templateRelease string, er *models.EnvironmentRegion) (*models.Cluster, []*models.Tag) {
	var gitURL, gitSubfolder, gitRef, gitRefType string
	if r.Git != nil {
		gitURL, gitSubfolder, gitRefType, gitRef = r.Git.URL,
			r.Git.Subfolder, r.Git.RefType(), r.Git.Ref()
	} else {
		gitURL = cluster.GitURL
		gitSubfolder = cluster.GitSubfolder
		gitRefType = cluster.GitRefType
		gitRef = cluster.GitRef
	}

	tags := make([]*models.Tag, 0)
	for _, tag := range r.Tags {
		tags = append(tags, &models.Tag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}

	return &models.Cluster{
		EnvironmentName: er.EnvironmentName,
		RegionName:      er.RegionName,
		Description:     r.Description,
		GitURL:          gitURL,
		GitSubfolder:    gitSubfolder,
		GitRef:          gitRef,
		GitRefType:      gitRefType,
		Template:        cluster.Template,
		TemplateRelease: templateRelease,
		Status:          cluster.Status,
		ExpireSeconds:   cluster.ExpireSeconds,
	}, tags
}

func getUserFromMap(id uint, userMap map[uint]*appmodels.User) *appmodels.User {
	user, ok := userMap[id]
	if !ok {
		return nil
	}
	return user
}

func toUser(user *appmodels.User) *User {
	if user == nil {
		return nil
	}
	return &User{
		ID:    user.ID,
		Name:  user.FullName,
		Email: user.Email,
	}
}

func ofClusterModel(application *appmodels.Application, cluster *models.Cluster, fullPath, namespace string,
	pipelineJSONBlob, applicationJSONBlob map[string]interface{}, tags ...*models.Tag) *GetClusterResponse {
	expireTime := ""
	if cluster.ExpireSeconds > 0 {
		expireTime = time.Duration(cluster.ExpireSeconds * 1e9).String()
	}

	return &GetClusterResponse{
		CreateClusterRequest: &CreateClusterRequest{
			Base: &Base{
				Description: cluster.Description,
				Tags:        models.Tags(tags).IntoTagsBasic(),
				Git: codemodels.NewGit(cluster.GitURL, cluster.GitSubfolder,
					cluster.GitRefType, cluster.GitRef),
				TemplateInput: &TemplateInput{
					Application: applicationJSONBlob,
					Pipeline:    pipelineJSONBlob,
				},
			},
			Name:       cluster.Name,
			Namespace:  namespace,
			ExpireTime: expireTime,
		},
		ID:       cluster.ID,
		FullPath: fullPath,
		Application: &Application{
			ID:   application.ID,
			Name: application.Name,
		},
		Priority: string(application.Priority),
		Template: &Template{
			Name:    cluster.Template,
			Release: cluster.TemplateRelease,
		},
		Scope: &Scope{
			Environment: cluster.EnvironmentName,
			Region:      cluster.RegionName,
		},
		Status:    cluster.Status,
		CreatedAt: cluster.CreatedAt,
		UpdatedAt: cluster.UpdatedAt,
	}
}

type GitResponse struct {
	GitURL  string `json:"gitURL"`
	HTTPURL string `json:"httpURL"`
}

type ListClusterResponse struct {
	ID          uint                  `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Scope       *Scope                `json:"scope"`
	Template    *Template             `json:"template"`
	Git         *GitResponse          `json:"git"`
	IsFavorite  *bool                 `json:"isFavorite"`
	CreatedAt   time.Time             `json:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt"`
	Tags        []*appmodels.TagBasic `json:"tags,omitempty"`
}

func ofCluster(cluster *appmodels.Cluster) *ListClusterResponse {
	return &ListClusterResponse{
		ID:          cluster.ID,
		Name:        cluster.Name,
		Description: cluster.Description,
		Scope: &Scope{
			Environment: cluster.EnvironmentName,
			Region:      cluster.RegionName,
		},
		Template: &Template{
			Name:    cluster.Template,
			Release: cluster.TemplateRelease,
		},
		Git: &GitResponse{
			GitURL: cluster.GitURL,
		},
		CreatedAt: cluster.CreatedAt,
		UpdatedAt: cluster.UpdatedAt,
	}
}

func ofClusterWithEnvAndRegion(cluster *appmodels.ClusterWithRegion) *ListClusterResponse {
	resp := ofCluster(cluster.Cluster)
	resp.Scope.RegionDisplayName = cluster.RegionDisplayName
	return resp
}

func ofClustersWithEnvRegionTags(clusters []*appmodels.ClusterWithRegion,
	tags []*appmodels.Tag) []*ListClusterResponse {
	tagMap := map[uint][]*appmodels.TagBasic{}
	for _, tag := range tags {
		tagBasic := &appmodels.TagBasic{
			Key:   tag.Key,
			Value: tag.Value,
		}
		tagMap[tag.ResourceID] = append(tagMap[tag.ResourceID], tagBasic)
	}

	respList := make([]*ListClusterResponse, 0)
	for _, c := range clusters {
		cluster := ofClusterWithEnvAndRegion(c)
		cluster.Tags = tagMap[c.ID]
		respList = append(respList, cluster)
	}
	return respList
}

type GetClusterByNameResponse struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Template    *Template       `json:"template"`
	Git         *codemodels.Git `json:"git"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	FullPath    string          `json:"fullPath"`
}

type ListClusterWithFullResponse struct {
	*ListClusterResponse
	IsFavorite *bool  `json:"isFavorite,omitempty"`
	FullName   string `json:"fullName,omitempty"`
	FullPath   string `json:"fullPath,omitempty"`
}

type ListClusterWithExpiryResponse struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	EnvironmentName string    `json:"environmentName"`
	RegionName      string    `json:"regionName"`
	Status          string    `json:"status"`
	ExpireSeconds   uint      `json:"expireSeconds"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func ofClusterWithExpiry(clusters []*appmodels.Cluster) []*ListClusterWithExpiryResponse {
	resList := make([]*ListClusterWithExpiryResponse, 0, len(clusters))
	for _, c := range clusters {
		resList = append(resList, &ListClusterWithExpiryResponse{
			ID:              c.ID,
			Name:            c.Name,
			EnvironmentName: c.EnvironmentName,
			RegionName:      c.RegionName,
			Status:          c.Status,
			ExpireSeconds:   c.ExpireSeconds,
			UpdatedAt:       c.UpdatedAt,
		})
	}
	return resList
}
