package server

type closeGroup struct {
	count   int
	closeCh chan struct{}
}

func (cg *closeGroup) Add(fn func() error) func() {
	if cg.count == 0 {
		cg.closeCh = make(chan struct{})
	}
	cg.count++
	return func() {
		err := fn()
		if err != nil {
			panic(err.Error())
		}
		cg.count--
		if cg.count == 0 {
			close(cg.closeCh)
			cg.count = -1
		}
	}
}

func (cg *closeGroup) Done() chan struct{} {
	return cg.closeCh
}
