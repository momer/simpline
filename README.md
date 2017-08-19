Why not have a simple, leaky, abstraction enforce similarity across projects?

```golang
	db := "A Database, or something"

	pipeline := []simpline.Pipe{
		{StepOne, "Step One"},
		{StepTwo, "Step Two"},
		{NewStepThree(db), "Step Three"}, // a Closure
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
			}
		}

		return
	}(p)

	// Queue up some work
	for i := 0; i < 1000; i++ {
		go func(p simpline.Pipeline) {
			ctx = context.WithValue(context.Background(), "MainZero", "A String!")
			p.Queue <- ctx
		}(p)
	}

	p.Close()
}
```
