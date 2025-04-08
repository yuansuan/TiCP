package execscript

import (
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

func (api API) ScriptContent(scriptContent string) Option {
	return func(req *session.ExecScriptRequest) error {
		req.ScriptContent = &scriptContent
		return nil
	}
}

func (api API) WaitTillEnd(waitTillEnd bool) Option {
	return func(req *session.ExecScriptRequest) error {
		req.WaitTillEnd = &waitTillEnd
		return nil
	}
}
