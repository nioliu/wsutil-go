# wsutil-go

wsutil-go 是一个 Go 语言实现的 WebSocket 连接管理库。它可以让你轻松创建和管理 WebSocket
连接，包括连接数限制、连接超时、心跳检测、消息广播等功能。WebSocket Group 支持用户自定义升级器和消息处理器，以及使用 Zap
日志库记录日志。

## 安装

使用 Go 模块进行安装：

```go

go get git.woa.com/nioliu/wsutil-go /group

```

## 使用

使用 WebSocket Group 非常简单。以下是一个示例程序：

```go
import (
"fmt"
"net/http"
"time"

"git.woa.com/nioliu/wsutil-go/group"
"git.woa.com/nioliu/wsutil-go/ws"

)

func main() {
package main

// 创建 WebSocket Group 实例
g := group.New()

// 设置 WebSocket 升级器
g.Apply(group.WithUpgrader(&ws.WrappedGorillaUpgrader{}))

// 设置心跳检测
g.Apply(group.WithHeartCheck(time.Minute))

// 设置消息广播失败处理函数
g.Apply(group.WithHandleBroadcastError(func (g *group.Group, conn *ws.SingleConn, err error) error {
fmt.Printf("广播消息失败: %s\n", err.Error())
return nil
}))

// 设置最大连接数
g.Apply(group.WithMaxConnCnt(100))

// 设置连接超时
g.Apply(group.WithMaxConnDuration(time.Hour * 24 * 30))

// 设置 Group ID
g.Apply(group.WithGroupId("my-group"))

// 设置 Group Map
g.Apply(group.WithGroupMap(group.Map{}))

// 设置消息处理器
g.HandleMsg(func (conn *ws.SingleConn, msgType int, data []byte) error {
// 在这里处理接收到的消息
return nil
})

// 启动 HTTP 服务器
http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
// 将 HTTP 请求升级为 WebSocket 连接
conn, err := g.WsUpgrader.Upgrade(w, r, nil)
if err != nil {
// 处理错误
return
}

// 将连接加入 Group 中
if err := g.AddConn(conn); err != nil {
// 处理错误
return
}
})

if err := http.ListenAndServe(":8080", nil); err != nil {
// 处理错误
return
}

}

```

首先，我们创建了一个 WebSocket Group 实例。接下来，我们设置了 WebSocket 升级器、心跳检测、消息广播失败处理函数、最大连接数、连接超时、Group
ID 和 Group Map。然后，我们设置了消息处理器，以便处理接收到的消息。最后，我们启动了一个 HTTP 服务器，并将 HTTP 请求升级为
WebSocket 连接，然后将连接添加到 Group 中。

## 配置Group

Group可以通过Option来进行配置，以下是可用的Option：

- WithMaxConnCnt(cnt int): 设置Group中的最大连接数
- WithHeartCheck(duration time.Duration): 设置心跳检查间隔
- WithHandleBroadcastError(f func(g *Group, conn *ws.SingleConn, err error) error): 设置处理广播消息错误的函数
- WithMaxConnDuration(duration time.Duration): 设置连接最大存活时间
- WithUpgrader(upgrader ws.Upgrader): 设置WebSocket Upgrader
- WithGroupId(id string): 设置Group ID
- WithGroupMap(m Map): 设置Group内部使用的Map
- WithBeforeHandleHookFunc(f ws.HandleMsgFunc): 设置处理消息前的钩子函数
- WithAfterHandleHookFunc(f ws.HandleMsgFunc): 设置处理消息后的钩子函数
- 每个Option都是一个函数，将其作为参数传递给Group的构造函数即可。

```go
g := group.New(
group.WithMaxConnCnt(100),
group.WithHeartCheck(time.Minute),
group.WithUpgrader(&ws.WrappedGorillaUpgrader{}),
// 添加其他Option
)
```
