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

package webhook

import (
	"context"
	"log"

	webhookcfg "github.com/horizoncd/horizon/pkg/config/webhook"
	eventhandlersvc "github.com/horizoncd/horizon/pkg/eventhandler"
	"github.com/horizoncd/horizon/pkg/eventhandler/wlgenerator"
	"github.com/horizoncd/horizon/pkg/jobs"
	"github.com/horizoncd/horizon/pkg/param/managerparam"
)

// New runs the agent.
func New(ctx context.Context, eventHandlerService eventhandlersvc.Service,
	webhookCfg webhookcfg.Config, mgrs *managerparam.Manager) (jobs.Job, WebhookService) {
	if err := eventHandlerService.RegisterEventHandler("webhook",
		wlgenerator.NewWebhookLogGenerator(mgrs)); err != nil {
		log.Printf("failed to register event handler, error: %s", err.Error())
		panic(err)
	}
	webhookService := NewWebhookServiceService(ctx, mgrs, webhookCfg)

	return func(ctx context.Context) {
		// start webhook service with multi workers to consume webhook logs and send webhook events
		webhookService.Start()

		<-ctx.Done()
		// graceful exit
		webhookService.StopAndWait()
	}, webhookService
}
