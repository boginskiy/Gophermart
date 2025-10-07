package config

type Args struct {
}

func NewArgs() *Args {
	return &Args{}
}

func (a *Args) GetSomeParam() any {
	return ""
}
