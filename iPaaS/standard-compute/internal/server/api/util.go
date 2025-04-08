package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/state"
)

func getState(c *gin.Context) (*state.State, error) {
	stateI, exist := c.Get("state")
	if !exist {
		return nil, fmt.Errorf("ginCtx['state'] not exist")
	}

	s, ok := stateI.(*state.State)
	if !ok {
		return nil, fmt.Errorf("ginCtx['state'] cannot convert to *state.State")
	}

	return s, nil
}
