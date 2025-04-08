package execscript

import (
	"encoding/base64"

	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"
)

type Option func(req *session.ExecScriptRequest) error

func (api API) SessionId(sessionId string) Option {
	return func(req *session.ExecScriptRequest) error {
		req.SessionId = &sessionId
		return nil
	}
}

func (api API) ScriptRunner(scriptRunner string) Option {
	return func(req *session.ExecScriptRequest) error {
		req.ScriptRunner = &scriptRunner
		return nil
	}
}

// ScriptContentEncoded base64 encoded value
func (api API) ScriptContentEncoded(scriptContentEncoded string) Option {
	return func(req *session.ExecScriptRequest) error {
		req.ScriptContent = &scriptContentEncoded
		return nil
	}
}

// ScriptContentRaw raw script content
func (api API) ScriptContentRaw(scriptContentRaw string) Option {
	return func(req *session.ExecScriptRequest) error {
		scriptContentEncoded := base64.StdEncoding.EncodeToString([]byte(scriptContentRaw))
		req.ScriptContent = &scriptContentEncoded
		return nil
	}
}

func (api API) WaitTillEnd(waitTillEnd bool) Option {
	return func(req *session.ExecScriptRequest) error {
		req.WaitTillEnd = &waitTillEnd
		return nil
	}
}
