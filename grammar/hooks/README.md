# hooks

给数据库driver注入sqlhooks.Hooks。

## Usage

```go
package main

import (
	"github.com/yaoapp/yao/cmd"
	"github.com/yaoapp/yao/utils"

	_ "github.com/yaoapp/gou/encoding"
	_ "github.com/yaoapp/yao/aigc"
	_ "github.com/yaoapp/yao/crypto"
	_ "github.com/yaoapp/yao/helper"
	_ "github.com/yaoapp/yao/openai"
	_ "github.com/yaoapp/yao/wework"
	
	_ "github.com/yaoapp/yao/xun/grammar/hooks"
	_ "github.com/yaoapp/yao/xun/grammar/hooks/log"
	// your customized hooks
)

// 主程序
func main() {
	// 1. 注入你的driver
	hooks.RegisterDriver("mysql:log")
	// 2. 此处可定制log hook的字段，例如从context中获取trace_id
	loghook.Default.ContextFields = func(ctx context.Context) log.F {
		return log.F{"trace_id": "-"}
	}
	utils.Init()
	cmd.Execute()
}
```

在.env文件中，指定YAO_DB_DRIVER="mysql:log"，编译启动即可。
