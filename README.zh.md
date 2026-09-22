# ingot ABI

> ingot Plugin 与生成式 Runtime Image 之间固定、受版本约束的 Go Contract。

[English](./README.md)

`github.com/ingot-agent/ingot-abi` 定义 ingot Builder 按精确 Go 类型身份识别的
最小公共类型集合。它将 Runtime 独占的宿主行为，与
[`github.com/ingot-agent/sdk`](https://github.com/ingot-agent/sdk) 中可替换的
Agent Contract 明确分离。

这里的 ABI 是用于静态生成 Component wiring 的源码级 Go ABI，不是动态链接
ABI、Plugin loader，也不是通用框架。

## 为什么需要独立模块

Component 之间交换的大多数 Contract 都是普通、可替换的能力，例如 model、
tool、session、prompt、filesystem 和 interaction channel。这些 Contract 可以
位于 Agent SDK、其他领域 SDK，或 Plugin 自己的 module 中。

少数 Contract 的性质不同：生成的 Runtime Image 独占进程生命周期、调用模式、
Plugin state 分配和精确的 Component constructor 形状。Plugin 无法为这些语义
提供等价替换。将它们放入独立模块，可以让这部分特权保持显式、极小、受版本
约束且可审计。

## Package

| Package | 用途 |
|---|---|
| `ingotabi` | Component ABI：`Cleanup`、`Optional[T]`、`Named[T]` 及其 helper。 |
| `invocation` | 只读 Runtime 参数，以及 `ModeRun`/`ModeCheck`。 |
| `lifecycle` | 通过 `Controller` 提交优雅进程关闭请求。 |
| `state` | 绝对路径形式的 Plugin-scoped 持久化位置。 |

本 module 只依赖 Go 标准库。

## Component ABI

每个 Component constructor 都必须具有以下精确形状：

```go
package component

import (
	"context"

	ingotabi "github.com/ingot-agent/ingot-abi"
)

type Dependencies struct{}
type Exports struct{}

func New(
	ctx context.Context,
	deps Dependencies,
) (Exports, ingotabi.Cleanup, error) {
	return Exports{}, nil, nil
}
```

当前 Core Builder 要求上述双参数构造函数。Plugin 通过显式 `state.Scope`
依赖自行加载配置，不再接收 `Config` 参数，也不存在统一 Runtime 配置解码。
参见 [Core 当前文件格式](https://github.com/ingot-agent/ingot/blob/main/docs/FILE_FORMATS.md)
和 [ADR 0003](https://github.com/ingot-agent/ingot/blob/main/docs/adr/0003-plugin-configuration.md)。

Builder 还会把 `ingotabi.Optional[T]` 识别为可选依赖，把
`ingotabi.Named[T]` 识别为稳定的 Runtime Instance Identity。这些 wrapper
属于 Component ABI；其他 module 中形状相同的类型仍按普通 Contract type
处理。

## 宿主 Contract

Component 通过 `Dependencies` 显式声明宿主能力：

```go
import (
	"github.com/ingot-agent/ingot-abi/invocation"
	"github.com/ingot-agent/ingot-abi/lifecycle"
	"github.com/ingot-agent/ingot-abi/state"
)

type Dependencies struct {
	Invocation invocation.Invocation
	Lifecycle  lifecycle.Controller
	State      state.Scope
}
```

生成的 Runtime Image 将这些精确类型作为虚拟宿主 Provider 注入。它们不会增加
Component 之间的图边，Plugin 也不得导出这些类型。

### Invocation

- `Arguments` 返回调用者拥有的副本，并排除 ingot 自有参数。
- `ModeRun` 表示正常执行。
- `ModeCheck` 构造、校验并清理完整 Graph，但不启动交互 loop，也不保留外部资源。

### Lifecycle

- `RequestShutdown(nil)` 表示正常完成意图。
- `RequestShutdown(err)` 记录进程级失败。
- 第一次请求取消 Runtime Context。
- 请求并发安全且不阻塞调用者。
- 非 nil shutdown cause 与 Cleanup error 会被聚合，最终进程结果由生成的
  Runtime 决定。

### State

- `Scope.Dir()` 返回根据 Plugin identity 分配的绝对目录。
- 同一个 Plugin 中的 Component 共享 state directory。
- 文件访问、schema 校验、迁移与持久化由 Plugin 负责。

## Builder 不变量

- Builder 固定 module path 与精确支持版本。
- 选中的 module version 与源码身份进入 lock file 和 Image ID。
- 生产构建拒绝用户指定的 path、version 与 replacement。
- 只有精确的 ingot ABI 类型身份会获得宿主注入或 wrapper 语义。
- Host dependency 必须在 Graph inspection 中保持可见。

## 不应进入本模块的能力

不要加入 HTTP client、filesystem、logger、metrics、tracing、clock、scheduler、
event bus、secret、业务配置、model/tool/session Contract、UI 协议、service
locator 或可扩展 registry。它们属于可替换能力或应用层关注点。

## 兼容性与版本

Builder 会精确锁定本 module。修改导出形状，或修改已承诺的 ownership、并发、
顺序、取消和错误语义时，必须协调发布 ingot ABI 与 Builder。任何新 API 都必须
满足全部宿主 ABI 准入规则，并保持 module 无第三方依赖。

## 开发

参见[贡献指南](CONTRIBUTING.md)中的合同准入与评审规则、
[发布协调说明](RELEASE.md)中的 ABI/Core 协作流程，以及
[安全报告说明](SECURITY.md)。

```sh
go test -race ./...
go vet ./...
```

## License

[Apache License 2.0](./LICENSE)

## 设计历史

[v0.1 提案](docs/design-history/README.md) 已从 Core 迁入，保留原 MIT 许可，
仅供追溯；当前构造函数和配置边界以本 README 及 Core Builder 为准。
