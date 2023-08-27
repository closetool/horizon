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

package global

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Model struct {
	ID        uint                  `gorm:"primarykey" json:"id"`
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
	DeletedTs soft_delete.DeletedAt `json:"-"`
}

type HorizonMetaData struct {
	Application   string
	Cluster       string
	ApplicationID uint
	ClusterID     uint
	Environment   string
	Operator      string
	PipelinerunID uint
	Region        string
	Template      string
	EventID       string
}
