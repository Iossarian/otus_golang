package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		stageData := getNextStageData(done, in)
		in = stage(stageData)
	}
	return in
}

func getNextStageData(done In, in In) Bi {
	newStageData := make(Bi)
	go func() {
		defer close(newStageData)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				newStageData <- v
			}
		}
	}()

	return newStageData
}
