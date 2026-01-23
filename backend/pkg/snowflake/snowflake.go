package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"time"
)

var node *snowflake.Node

// Init 初始化雪花算法节点，用于生成全局唯一 ID。
func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(machineID)
	return
}

// GenIDByInt GenIDByInt返回int64
func GenIDByInt() int64 {
	return node.Generate().Int64()
}

// GenIDByString GenIDByString返回string
func GenIDByString() string {
	return node.Generate().String()
}
