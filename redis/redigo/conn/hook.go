/**
 * @Author raven
 * @Description
 * @Date 2022/7/27
 **/
package conn

import (
	"context"
)

type Hook interface {
	BeforeProcess(ctx context.Context, cmd *baseCmder) (context.Context, error)
	AfterProcess(ctx context.Context, cmd *baseCmder) error

	BeforeProcessPipeline(ctx context.Context, cmds []*baseCmder) (context.Context, error)
	AfterProcessPipeline(ctx context.Context, cmds []*baseCmder) error
}

type DoCommand func(context.Context, *baseCmder) (interface{}, error)
type DoPipeLineCommand func(context.Context, []*baseCmder) ([]interface{}, error)

type hooks struct {
	hooks []Hook
}

func (hs hooks) AddHook(hook Hook) {
	hs.hooks = append(hs.hooks, hook)
}

func (hs hooks) process(ctx context.Context, cmd *baseCmder, do DoCommand) (interface{}, error) {
	if len(hs.hooks) == 0 {
		return do(ctx, cmd)
	}

	var hookIndex int
	var retErr error
	var result interface{}

	for ; hookIndex < len(hs.hooks) && retErr == nil; hookIndex++ {
		ctx, retErr = hs.hooks[hookIndex].BeforeProcess(ctx, cmd)
		if retErr != nil {
			cmd.SetErr(retErr)
		}
	}

	if retErr == nil {
		result, retErr = do(ctx, cmd)
	}

	for hookIndex--; hookIndex >= 0; hookIndex-- {
		if err := hs.hooks[hookIndex].AfterProcess(ctx, cmd); err != nil {
			retErr = err
			cmd.SetErr(retErr)
		}
	}

	return result, retErr
}

func (hs hooks) processPipeline(ctx context.Context, cmds []*baseCmder, do DoPipeLineCommand) ([]interface{}, error) {
	if len(hs.hooks) == 0 {
		return do(ctx, cmds)
	}

	var hookIndex int
	var retErr error
	var result []interface{}

	for ; hookIndex < len(hs.hooks) && retErr == nil; hookIndex++ {
		ctx, retErr = hs.hooks[hookIndex].BeforeProcessPipeline(ctx, cmds)
		if retErr != nil {
			setCmdsErr(cmds, retErr)
		}
	}

	if retErr == nil {
		result, retErr = do(ctx, cmds)
	}

	for hookIndex--; hookIndex >= 0; hookIndex-- {
		if err := hs.hooks[hookIndex].AfterProcessPipeline(ctx, cmds); err != nil {
			retErr = err
			setCmdsErr(cmds, retErr)
		}
	}

	return result, retErr
}

func setCmdsErr(cmds []*baseCmder, e error) {
	for _, cmd := range cmds {
		if cmd.Err() == nil {
			cmd.SetErr(e)
		}
	}
}
