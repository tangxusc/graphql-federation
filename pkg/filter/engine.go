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

func start(ctx context.Context, current *graphqlFederationConfig) {
	ticker := time.NewTicker(current.SchemaRefreshInterval)
	var lastCfg *graphqlFederationConfig
	for {
		select {
		case <-ctx.Done():
			return
		case cfg := <-configUpdateCh:
			lastCfg = cfg
			updateGraphqlEngine(ctx, lastCfg)
			ticker.Reset(lastCfg.SchemaRefreshInterval)
		case <-ticker.C:
			updateGraphqlEngine(ctx, lastCfg)
		}
	}
}

func updateGraphqlEngine(ctx context.Context, cfg *graphqlFederationConfig) {
	api.LogDebugf("[graphql-federation] 更新graphql engine...")
	var subgraphsConfigs = make([]engine.SubgraphConfiguration, len(cfg.SubGraphqlConfig), len(cfg.SubGraphqlConfig))
	for i, configuration := range cfg.SubGraphqlConfig {
		subgraphsConfigs[i] = engine.SubgraphConfiguration{
			Name:                 configuration.ServiceName,
			URL:                  configuration.GraphqlUrl,
			SDL:                  "",
			SubscriptionUrl:      configuration.GraphqlUrl,
			SubscriptionProtocol: engine.SubscriptionProtocolWS,
		}
	}
	client := http.DefaultClient
	client.Timeout = cfg.SchemaRefreshTimeout
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
