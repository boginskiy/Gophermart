package tests

type TestLogg struct {
}

func NewTestLogg() *TestLogg {
	return &TestLogg{}
}

func (tl *TestLogg) RaiseError(msg string, err error) {
}

func (tl *TestLogg) RaiseFatal(msg string, err error) {
}

func (tl *TestLogg) RaisePanic(msg string, err error) {
}

func (tl *TestLogg) RaiseInfo(msg string) {
}

func (tl *TestLogg) Clouse() {
}
