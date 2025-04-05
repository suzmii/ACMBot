package ctxUtil

import "context"

type ChainContext struct {
	ctx   context.Context
	tasks []Task
}

type Task func(ctx context.Context) (context.Context, error)

func NewChainContext(ctx context.Context) *ChainContext {
	return &ChainContext{ctx: ctx}
}

func (c *ChainContext) Then(task Task) *ChainContext {
	c.tasks = append(c.tasks, task)
	return c
}

func (c *ChainContext) Execute() error {
	ctx := c.ctx
	var err error
	for _, task := range c.tasks {
		ctx, err = task(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
