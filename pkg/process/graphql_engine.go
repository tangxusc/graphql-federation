package process

import (
	"net/http"
	"sync"

	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
	"github.com/tangxusc/higress-graphql-federation/pkg/config"
	"github.com/wundergraph/graphql-go-tools/execution/engine"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
)

var HttpClient = &http.Client{
	Transport: &HttpClientRoundTripper{},
}
var EngineConfigFactory *engine.FederationEngineConfigFactory
var ExecutionEngine *engine.ExecutionEngine
var lock sync.RWMutex

func initGraphqlFederationEngine(config config.FederationConfig) error {
	lock.Lock()
	defer lock.Unlock()
	if EngineConfigFactory != nil && ExecutionEngine != nil {
		return nil
	}

	var subgraphsConfigs []engine.SubgraphConfiguration = make([]engine.SubgraphConfiguration, 1)
	subgraphsConfigs[0] = engine.SubgraphConfiguration{
		Name:                 "test",
		URL:                  "http://httpbin",
		SDL:                  "",
		SubscriptionUrl:      "",
		SubscriptionProtocol: engine.SubscriptionProtocolWS,
	}
	EngineConfigFactory = engine.NewFederationEngineConfigFactory(
		config.Ctx,
		subgraphsConfigs, engine.WithFederationHttpClient(HttpClient),
	)
	engineConfig, err := EngineConfigFactory.BuildEngineConfiguration()
	if err != nil {
		proxywasm.LogErrorf("build engine configuration error:%v", err)
		return err
	}
	ExecutionEngine, err = engine.NewExecutionEngine(config.Ctx, config.Logger, engineConfig, resolve.ResolverOptions{
		MaxConcurrency: 1024,
	})
	if err != nil {
		proxywasm.LogErrorf("new execution engine error:%v", err)
		return err
	}
	return nil
}
