package simpline

import (
	"context"
)

type Pipeline struct {
	// Pipeline Context
	context *context.Context
	// Ordered slice of steps to be run
	Pipes []Pipe
	Queue chan context.Context
}

func NewPipeline(pipes ...Pipe) Pipeline {
	ctx := context.Background()

	return Pipeline{
		context: &ctx,
		Pipes:   pipes,
		Queue:   make(chan context.Context),
	}
}

func (p Pipeline) Close() {
	defer close(p.Queue)
	(*p.context).Done()
}

func (p Pipeline) Process(ctx context.Context) (context.Context, error) {
	var err error

	for _, pipe := range p.Pipes {
		ctx, err = pipe.Do(ctx, err)
	}

	return ctx, err
}

func (p Pipeline) WithPipes(pipes ...Pipe) Pipeline {
	p.Pipes = pipes

	return p
}

type DoFn func(context.Context, error) (context.Context, error)

type Pipe struct {
	Do   DoFn
	Name string
}
