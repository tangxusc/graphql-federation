package filter

import (
	"context"
	"net/http"
	"time"

	"github.com/envoyproxy/envoy/contrib/golang/common/go/api"
	"github.com/jensneuse/abstractlogger"
	"github.com/wundergraph/graphql-go-tools/execution/engine"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
)

var configUpdateCh = make(chan *graphqlFederationConfig, 1)
var engineConfigFactory *engine.FederationEngineConfigFactory
var logger abstractlogger.Logger = &GraphqlFederationLoggerAdapter{}
var executionEngine *engine.ExecutionEngine

func start(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 5)
	var lastCfg *graphqlFederationConfig
	for {
		select {
		case <-ctx.Done():
			return
		case cfg := <-configUpdateCh:
			updateGraphqlEngine(ctx, cfg)
			lastCfg = cfg
		case <-ticker.C:
			updateGraphqlEngine(ctx, lastCfg)
		}
	}
}

func updateGraphqlEngine(ctx context.Context, cfg *graphqlFederationConfig) {
	api.LogDebugf("[graphql-federation] 更新graphql engine...")
	var subgraphsConfigs = make([]engine.SubgraphConfiguration, 1)
	subgraphsConfigs[0] = engine.SubgraphConfiguration{
		Name:                 "httpbin",
		URL:                  "",
		SDL:                  "",
		SubscriptionUrl:      "",
		SubscriptionProtocol: engine.SubscriptionProtocolWS,
	}
	client := http.DefaultClient
	client.Timeout = time.Second * 1
	engineConfigFactory = engine.NewFederationEngineConfigFactory(
		ctx,
		subgraphsConfigs,
		engine.WithFederationHttpClient(client),
	)
	engineConfig, err := engineConfigFactory.BuildEngineConfiguration()
	if err != nil {
		api.LogErrorf("[graphql-federation]build graphql federation engine configuration error:%v", err)
		return
	}

	executionEngine, err = engine.NewExecutionEngine(ctx, logger, engineConfig, resolve.ResolverOptions{
		MaxConcurrency: 1024,
	})
	if err != nil {
		api.LogErrorf("[graphql-federation]new graphql federation execution engine error:%v", err)
		return
	}
}
