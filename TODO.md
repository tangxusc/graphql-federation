# TODO

## graphql-go-tools-execution
因 `graphql.EngineResultWriter`未实现`resolve.SubscriptionResponseWriter`接口,编译出现错误
故修改 [result_writer.go](graphql-go-tools-execution/graphql/result_writer.go) 文件,添加如下代码
```go

// Heartbeat implements resolve.SubscriptionResponseWriter.
func (e *EngineResultWriter) Heartbeat() error {
	panic("unimplemented")
}
```

## graphql-go-tools-v2
因 `graphql-go-tools-v2/pkg/netpoll/netpoll.go`中
```go
func Supported() error {
	// Create an instance of the poller
	poller, err := NewPoller(1, 10*time.Millisecond)
	if err != nil {
		return ErrUnsupported
	}
	defer poller.Close(true)
```
此处`NewPoller`在wasm中未声明,故在[netpoll_unsupported.go](graphql-go-tools-v2/pkg/netpoll/netpoll_unsupported.go)中加入`ErrUnsupported`


```go
//go:build windows || wasm
// +build windows wasm
```
## composition-go
因`[config_factory_federation.go](graphql-go-tools-execution/engine/config_factory_federation.go)中函数`func (f *FederationEngineConfigFactory) BuildEngineConfiguration() (Configuration, error) {`未正确使用已配置的http.Client, 故
修改文件[composition.go](composition-go/composition.go),传递已设置的http.Client
```go
func BuildRouterConfiguration(client *http.Client, subgraphs ...*Subgraph) (string, error) {
	updatedSubgraphs, err := updateSchemas(client, subgraphs)
```
