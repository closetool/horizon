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

package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	NeverExpire     = "never"
	ExpiresAtFormat = "2006-01-02"
)

const ClaimsIssuer = "horizon"

type Claims struct {
	PipelinerunID *uint
	jwt.RegisteredClaims
}

type ClaimsOption func(*Claims)

func WithPipelinerunID(pipelinerunID uint) ClaimsOption {
	return func(claims *Claims) {
		claims.PipelinerunID = &pipelinerunID
	}
}

type Token struct {
	ID uint `gorm:"primarykey"`

	// grant client info
	ClientID    string `gorm:"column:client_id"`
	RedirectURI string `gorm:"column:redirect_uri"`
	State       string `gorm:"column:state"`

	// token basic info
	Name string `gorm:"column:name"`
	// Code authorize_code/access_token/refresh_token
	Code      string        `gorm:"column:code"`
	CreatedAt time.Time     `gorm:"column:created_at"`
	CreatedBy uint          `gorm:"column:created_by"`
	ExpiresIn time.Duration `gorm:"column:expires_in"`
	Scope     string        `gorm:"column:scope"`

	UserID uint `gorm:"column:user_id"`
}
