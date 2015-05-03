package controls

type SwitchNotFoundError string

type ExecError string

func (err SwitchNotFoundError) Error() string {
	return string(err)
}

func (err ExecError) Error() string {
	return string(err)
}
