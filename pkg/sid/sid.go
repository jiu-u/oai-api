package sid

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"strconv"
	"time"
)

type Sid struct {
	*snowflake.Node
}

func NewSid() *Sid {
	startTime := "2025-01-01"
	fmt.Println(startTime)
	machineID := "12"
	id64, err := strconv.ParseInt(machineID, 10, 64)
	if err != nil {
		panic(err)
	}
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		panic(err)
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err := snowflake.NewNode(id64)
	if err != nil {
		panic(err)
	}
	return &Sid{node}
}

func (s *Sid) GenString() string {
	return s.Generate().String()
}

func (s *Sid) GenInt64() int64 {
	return s.Generate().Int64()
}

func (s *Sid) GenUint64() uint64 {
	return uint64(s.GenInt64())
}
