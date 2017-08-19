Why not have a simple, leaky, abstraction enforce similarity across projects?

```golang
package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/momer/simpline"
)

func main() {
	db := "A Database"
	var wg sync.WaitGroup

	pipeline := []simpline.Pipe{
		{StepOne, "Step One"},
		{StepTwo, "Step Two"},
		{NewStepThree(db), "Step Three"},
	}

	p := simpline.NewPipeline(pipeline...)

	// Send off a goroutine to process the queue
	go func(p simpline.Pipeline) {
	Processing:
		for {
			select {
			case ctx, ok := <-p.Queue:
				if !ok {
					break Processing
				}
				p.Process(ctx)
				wg.Done()
			}
		}

		return
	}(p)

	// Queue up some work
	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func(p simpline.Pipeline) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, "MainZero", "A String!")
			p.Queue <- ctx
		}(p)
	}

	wg.Wait()
	p.Close()
}

func StepOne(ctx context.Context, err error) (context.Context, error) {
	str := fmt.Sprintf("%s, %s", ctx.Value("MainZero").(string), "StepOne")

	return context.WithValue(ctx, "StepOne", str), nil
}

func StepTwo(ctx context.Context, err error) (context.Context, error) {
	str := fmt.Sprintf("%s, %s", ctx.Value("StepOne").(string), "StepTwo")

	return context.WithValue(ctx, "StepTwo", str), nil
}

func NewStepThree(db string) simpline.DoFn {
	return func(ctx context.Context, err error) (context.Context, error) {
		fmt.Printf("%s, %s\n", ctx.Value("StepTwo").(string), db)

		return ctx, nil
	}
}```
