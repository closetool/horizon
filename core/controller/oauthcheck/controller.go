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

package oauthcheck

import (
	"fmt"
	"strings"
	"time"

	tokenmanager "github.com/horizoncd/horizon/pkg/manager"
	"golang.org/x/net/context"

	herrors "github.com/horizoncd/horizon/core/errors"
	"github.com/horizoncd/horizon/pkg/auth"
	rbactype "github.com/horizoncd/horizon/pkg/auth"
	"github.com/horizoncd/horizon/pkg/authentication/user"
	perror "github.com/horizoncd/horizon/pkg/errors"
	"github.com/horizoncd/horizon/pkg/oauth/scope"
	"github.com/horizoncd/horizon/pkg/param"
	"github.com/horizoncd/horizon/pkg/rbac/types"
)

type Controller interface {
	ValidateToken(ctx context.Context, token string) error
	LoadAccessTokenUser(ctx context.Context, token string) (user.User, error)
	CheckScopePermission(ctx context.Context, token string, authInfo auth.RequestInfo) (bool, string, error)
}

type controller struct {
	tokenManager tokenmanager.TokenManager
	userManager  tokenmanager.UserManager
	scopeService scope.Service
}

var _ Controller = &controller{}

func NewOauthChecker(param *param.Param) Controller {
	return &controller{
		tokenManager: param.TokenManager,
		userManager:  param.UserManager,
		scopeService: param.ScopeService,
	}
}

func (c *controller) ValidateToken(ctx context.Context, accessToken string) error {
	token, err := c.tokenManager.LoadTokenByCode(ctx, accessToken)
	if err != nil {
		return err
	}

	isExpired := func() bool {
		return token.CreatedAt.Add(token.ExpiresIn).Before(time.Now())
	}
	neverExpires := func() bool {
		return token.ExpiresIn <= 0
	}

	if neverExpires() {
		return nil
	}

	if isExpired() {
		return perror.Wrap(herrors.ErrOAuthAccessTokenExpired, "")
	}

	return nil
}

func (c *controller) LoadAccessTokenUser(ctx context.Context, accessToken string) (user.User, error) {
	token, err := c.tokenManager.LoadTokenByCode(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	usr, err := c.userManager.GetUserByID(ctx, token.UserID)
	if err != nil {
		return nil, err
	}
	return &user.DefaultInfo{
		Name:     usr.Name,
		FullName: usr.FullName,
		ID:       usr.ID,
		Email:    usr.Email,
		Admin:    usr.Admin,
	}, nil
}

func (c *controller) CheckScopePermission(ctx context.Context, accessToken string,
	requestInfo auth.RequestInfo) (bool, string, error) {
	token, err := c.tokenManager.LoadTokenByCode(ctx, accessToken)
	if err != nil {
		return false, "", err
	}

	usr, err := c.LoadAccessTokenUser(ctx, accessToken)
	if err != nil {
		return false, "", err
	}

	record := rbactype.AttributesRecord{
		User:            usr,
		Verb:            requestInfo.Verb,
		APIGroup:        requestInfo.APIGroup,
		APIVersion:      requestInfo.APIVersion,
		Resource:        requestInfo.Resource,
		SubResource:     requestInfo.Subresource,
		Name:            requestInfo.Name,
		Scope:           requestInfo.Scope,
		ResourceRequest: requestInfo.IsResourceRequest,
		Path:            requestInfo.Path,
	}

	scopeRoles := c.scopeService.GetRulesByScope(strings.Split(token.Scope, " "))
	for _, scopeRule := range scopeRoles {
		for i, policy := range scopeRule.PolicyRules {
			if types.RuleAllow(record, &policy) {
				reason := fmt.Sprintf("user %s allowed by scope(%s) by rule[%d]",
					usr.String(), scopeRule.Name, i)
				return true, reason, nil
			}
		}
	}
	return false, "", nil
}
