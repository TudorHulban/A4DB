package a4db

type condition func(o *Object) bool

func (c condition) And(other condition) condition {
	return func(o *Object) bool {
		return c(o) && other(o)
	}
}

func (c condition) Or(other condition) condition {
	return func(o *Object) bool {
		return c(o) || other(o)
	}
}

func (c condition) Not() condition {
	return func(o *Object) bool {
		return !c(o)
	}
}
