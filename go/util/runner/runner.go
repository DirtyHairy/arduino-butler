package runner

type T struct {
	commandChannel chan cmdWrapper
	controlChannel chan int

	Execute func(interface{}) error
}

type cmdWrapper struct {
	cmd           interface{}
	resultChannel chan error
}

func (runner *T) Start() error {
	if runner.commandChannel != nil {
		return Error("already started")
	}

	runner.commandChannel = make(chan cmdWrapper)
	runner.controlChannel = make(chan int)

	go func() {
		for {
			cmd, ok := <-runner.commandChannel

			if !ok {
				break
			}

			cmd.resultChannel <- runner.Execute(cmd.cmd)
		}

		runner.controlChannel <- 1
	}()

	return nil
}

func (runner *T) Stop() error {
	if runner.commandChannel == nil {
		return Error("already stopped")
	}

	close(runner.commandChannel)

	_ = <-runner.controlChannel

	runner.commandChannel = nil
	runner.controlChannel = nil

	return nil
}

func (runner *T) Dispatch(cmd interface{}) error {
	if runner.commandChannel == nil {
		return Error("not running")
	}

	resultChannel := make(chan error)

	runner.commandChannel <- cmdWrapper{cmd, resultChannel}

	return <-resultChannel
}
