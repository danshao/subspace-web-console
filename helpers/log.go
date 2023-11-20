package helpers

import (
	"encoding/json"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

// SubspaceRawLog basic struct
type SubspaceRawLog struct {
	Timestamp string            `json:"timestamp, string"`
	Level     string            `json:"level, string"`
	RemoteIP  string            `json:"remote_ip, string"`
	Method    string            `json:"method, string"`
	Path      string            `json:"path, string"`
	UserAgent string            `json:"user_agent, string"`
	Params    map[string]string `json:"params"`
	Success   *bool             `json:"success,omitempty"`
}

func NewSubspaceRawLog(ctx *context.Context, success *bool) *SubspaceRawLog {
	if success != nil {
		return &SubspaceRawLog{
			Timestamp: time.Now().String(),
			Level:     "info",
			RemoteIP:  ctx.Input.IP(),
			Method:    ctx.Input.Method(),
			Path:      ctx.Input.URI(),
			UserAgent: ctx.Input.UserAgent(),
			Params:    ctx.Input.Params(),
			Success:   success,
		}
	}

	return &SubspaceRawLog{
		Timestamp: time.Now().String(),
		Level:     "info",
		RemoteIP:  ctx.Input.IP(),
		Method:    ctx.Input.Method(),
		Path:      ctx.Input.URI(),
		UserAgent: ctx.Input.UserAgent(),
		Params:    ctx.Input.Params(),
	}
}

func (srl SubspaceRawLog) String() (rawLogString string) {
	b, err := json.Marshal(srl)
	if err != nil {
		beego.Error("[ERROR] ", err)
	}
	rawLogString = string(b[:])

	return
}
