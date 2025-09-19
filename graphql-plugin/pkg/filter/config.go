package filter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	xds "github.com/cncf/xds/go/xds/type/v3"
	"github.com/wundergraph/graphql-go-tools/execution/engine"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/envoyproxy/envoy/contrib/golang/common/go/api"
)

const Name = "graphql-federation"

type SubgraphConfiguration struct {
	ServiceName          string                      `json:"ServiceName"`
	GraphqlUrl           string                      `json:"GraphqlUrl"`
	SubscriptionProtocol engine.SubscriptionProtocol `json:"SubscriptionProtocol"`
	SubscriptionURL      string                      `json:"SubscriptionURL"`
}

type graphqlFederationConfig struct {
	SubGraphqlConfig      []*SubgraphConfiguration
	SchemaRefreshInterval time.Duration
	SchemaRefreshTimeout  time.Duration
}

type GraphqlFederationPluginConfigParser struct {
}

// Parse the GraphqlFederationFilter configuration. We can call the ConfigCallbackHandler to control the GraphqlFederationFilter's
// behavior
func (p *GraphqlFederationPluginConfigParser) Parse(any *anypb.Any, callbacks api.ConfigCallbackHandler) (interface{}, error) {
	configStruct := &xds.TypedStruct{}
	if err := any.UnmarshalTo(configStruct); err != nil {
		return nil, err
	}

	v := configStruct.Value
	conf := &graphqlFederationConfig{
		SchemaRefreshInterval: time.Minute * 5,
		SchemaRefreshTimeout:  time.Minute * 1,
	}

	// 解析SubGraphqlConfig配置
	subGraphqlConfig, ok := v.AsMap()["SubGraphqlConfig"]
	if !ok {
		return nil, errors.New("missing SubGraphqlConfig")
	}

	// 将配置转换为JSON再解析到结构体中
	jsonData, err := json.Marshal(subGraphqlConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SubGraphqlConfig: %v", err)
	}

	err = json.Unmarshal(jsonData, &conf.SubGraphqlConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal SubGraphqlConfig: %v", err)
	}

	refreshInterval, ok := v.AsMap()["SchemaRefreshInterval"]
	if ok {
		refreshIntervalString, ok := refreshInterval.(string)
		if !ok {
			return nil, fmt.Errorf("SchemaRefreshInterval must be a string(like 1s, 1m, 1h)")
		}
		refreshIntervalDuration, err := time.ParseDuration(refreshIntervalString)
		if err != nil {
			return nil, fmt.Errorf("failed to parse SchemaRefreshInterval: %v", err)
		}
		conf.SchemaRefreshInterval = refreshIntervalDuration
	}
	refreshTimeout, ok := v.AsMap()["SchemaRefreshTimeout"]
	if ok {
		refreshTimeoutString, ok := refreshTimeout.(string)
		if !ok {
			return nil, fmt.Errorf("SchemaRefreshTimeout must be a string(like 1s, 1m, 1h)")
		}
		refreshTimeoutDuration, err := time.ParseDuration(refreshTimeoutString)
		if err != nil {
			return nil, fmt.Errorf("failed to parse SchemaRefreshTimeout: %v", err)
		}
		conf.SchemaRefreshTimeout = refreshTimeoutDuration
	}

	api.LogDebugf("[graphql-federation] 解析配置完成...")
	go start(context.TODO(), conf)
	configUpdateCh <- conf
	return conf, nil
}

// Merge configuration from the inherited parent configuration
func (p *GraphqlFederationPluginConfigParser) Merge(parent interface{}, child interface{}) interface{} {
	parentConfig := parent.(*graphqlFederationConfig)
	childConfig := child.(*graphqlFederationConfig)

	// copy one, do not update parentConfig directly.
	newConfig := *parentConfig
	if childConfig.SubGraphqlConfig != nil {
		newConfig.SubGraphqlConfig = childConfig.SubGraphqlConfig
	}
	api.LogDebugf("[graphql-federation] 解析 新的 配置完成...")
	g := &newConfig
	configUpdateCh <- g
	return g
}

func GraphqlFederationPluginFactory(c interface{}, callbacks api.FilterCallbackHandler) api.StreamFilter {
	api.LogDebugf("[graphql-federation] 创建graphql federation plugin factory...")
	conf, ok := c.(*graphqlFederationConfig)
	if !ok {
		api.LogErrorf("[graphql-federation] 配置错误...")
		panic("unexpected graphqlFederationConfig type")
	}

	filter := NewGraphqlFederationFilter(conf, callbacks)
	api.LogDebugf("[graphql-federation] 配置filter完成...")
	return filter
}
