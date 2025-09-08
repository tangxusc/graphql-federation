package config

import (
	"context"
	"net/http"
	"time"

	"github.com/wundergraph/graphql-go-tools/execution/engine"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"

	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
	"github.com/higress-group/wasm-go/pkg/wrapper"
	"github.com/tidwall/gjson"
)

func init() {
	RegisterCtxOption(wrapper.ParseConfig(parseFederationConfig))
}

type FederationConfig struct {
	EnableQueryPlan     bool          `json:"enableQueryPlanning"`
	EnableCaching       bool          `json:"enableCaching"`
	MaxQueryDepth       int64         `json:"maxQueryDepth"`
	QueryTimeout        time.Duration `json:"queryTimeout"`
	EnableIntrospect    bool          `json:"enableIntrospection"`
	GraphqlAddress      string        `json:"graphqlAddress"`
	DebugMode           bool          `json:"debugMode"`
	EngineConfigFactory *engine.FederationEngineConfigFactory
	ExecutionEngine     *engine.ExecutionEngine
	Logger              WasmLogger
	HttpClient          *http.Client
	CancelFunc          context.CancelFunc
	Ctx                 context.Context
}

func parseFederationConfig(json gjson.Result, config *FederationConfig) error {
	config.EnableQueryPlan = json.Get("enableQueryPlanning").Bool()
	config.EnableCaching = json.Get("enableCaching").Bool()
	config.MaxQueryDepth = json.Get("maxQueryDepth").Int()
	config.QueryTimeout = time.Duration(json.Get("queryTimeout").Int()) * time.Millisecond
	config.EnableIntrospect = json.Get("enableIntrospection").Bool()
	config.GraphqlAddress = json.Get("graphqlAddress").String()
	config.DebugMode = json.Get("debugMode").Bool()

	config.Logger = WasmLogger{}
	config.HttpClient = &http.Client{
		Transport: &HttpClientRoundTripper{},
	}
	config.Ctx, config.CancelFunc = context.WithCancel(context.TODO())

	var subgraphsConfigs []engine.SubgraphConfiguration = make([]engine.SubgraphConfiguration, 1)
	subgraphsConfigs[0] = engine.SubgraphConfiguration{
		Name:                 "test",
		URL:                  "http://httpbin",
		SDL:                  "",
		SubscriptionUrl:      "",
		SubscriptionProtocol: engine.SubscriptionProtocolWS,
	}

	//TODO: http client
	//clusterClient := wrapper.NewClusterClient(wrapper.FQDNCluster{
	//	FQDN: serviceName,
	//	Port: servicePort,
	//})
	// config.EngineConfigFactory = engine.NewFederationEngineConfigFactory(config.Ctx, subgraphsConfigs)
	config.EngineConfigFactory = engine.NewFederationEngineConfigFactory(
		config.Ctx,
		subgraphsConfigs, engine.WithFederationHttpClient(config.HttpClient),
	)
	engineConfig, err := config.EngineConfigFactory.BuildEngineConfiguration()
	if err != nil {
		proxywasm.LogErrorf("build engine configuration error:%v", err)
		return err
	}
	config.ExecutionEngine, err = engine.NewExecutionEngine(config.Ctx, config.Logger, engineConfig, resolve.ResolverOptions{
		MaxConcurrency: 1024,
	})
	if err != nil {
		proxywasm.LogErrorf("new execution engine error:%v", err)
		return err
	}
	return nil
}
