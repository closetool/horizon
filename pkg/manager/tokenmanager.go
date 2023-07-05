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

	"github.com/horizoncd/horizon/pkg/models"
	"github.com/horizoncd/horizon/pkg/token/storage"
	"gorm.io/gorm"
)

type TokenManager interface {
	CreateToken(context.Context, *models.Token) (*models.Token, error)
	LoadTokenByID(context.Context, uint) (*models.Token, error)
	LoadTokenByCode(ctx context.Context, code string) (*models.Token, error)
	RevokeTokenByID(context.Context, uint) error
	RevokeTokenByClientID(ctx context.Context, clientID string) error
}

func NewTokenManager(db *gorm.DB) TokenManager {
	return &tokenManager{storage: storage.NewStorage(db)}
}

type tokenManager struct {
	storage storage.Storage
}

func (m *tokenManager) CreateToken(ctx context.Context, token *models.Token) (*models.Token, error) {
	return m.storage.Create(ctx, token)
}

func (m *tokenManager) LoadTokenByID(ctx context.Context, id uint) (*models.Token, error) {
	return m.storage.GetByID(ctx, id)
}

func (m *tokenManager) LoadTokenByCode(ctx context.Context, code string) (*models.Token, error) {
	return m.storage.GetByCode(ctx, code)
}

func (m *tokenManager) RevokeTokenByID(ctx context.Context, id uint) error {
	return m.storage.DeleteByID(ctx, id)
}

func (m *tokenManager) RevokeTokenByClientID(ctx context.Context, clientID string) error {
	return m.storage.DeleteByClientID(ctx, clientID)
}
