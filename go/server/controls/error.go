package controls

type ControlNotFoundError string

type ExecError string

func (err ControlNotFoundError) Error() string {
	return string(err)
}

func (err ExecError) Error() string {
	return string(err)
}
