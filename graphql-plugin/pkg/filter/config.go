package filter

import (
	"context"
	"errors"
	"fmt"

	xds "github.com/cncf/xds/go/xds/type/v3"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/envoyproxy/envoy/contrib/golang/common/go/api"
)

const Name = "graphql-federation"

type graphqlFederationConfig struct {
	graphqlPath        string
	graphqlContextType string
	graphqlMethod      string

	echoBody string
	// other fields
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
	conf := &graphqlFederationConfig{}
	prefix, ok := v.AsMap()["prefix_localreply_body"]
	if !ok {
		return nil, errors.New("missing prefix_localreply_body")
	}
	if str, ok := prefix.(string); ok {
		conf.echoBody = str
	} else {
		return nil, fmt.Errorf("prefix_localreply_body: expect string while got %T", prefix)
	}

	api.LogDebugf("[graphql-federation] 解析配置完成...")
	go start(context.TODO())
	configUpdateCh <- conf
	return conf, nil
}

// Merge configuration from the inherited parent configuration
func (p *GraphqlFederationPluginConfigParser) Merge(parent interface{}, child interface{}) interface{} {
	parentConfig := parent.(*graphqlFederationConfig)
	childConfig := child.(*graphqlFederationConfig)

	// copy one, do not update parentConfig directly.
	newConfig := *parentConfig
	if childConfig.echoBody != "" {
		newConfig.echoBody = childConfig.echoBody
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
