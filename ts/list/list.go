package list

type node[T any] struct {
	parent List[T]
	value  T
}

type List[T any] struct {
	node *node[T]
	size int
}

func New[T any]() List[T] {
	return List[T]{nil, 0}
}

func (l List[T]) Add(value T) List[T] {
	return List[T]{
		&node[T]{l, value},
		l.size + 1,
	}
}

func (l List[T]) Remove() List[T] {
	if l.node == nil {
		return l
	}

	return l.node.parent
}

func (l List[T]) Value() (T, bool) {
	if l.node == nil {
		var t T
		return t, false
	}

	return l.node.value, true
}

func (l List[T]) Empty() bool {
	return l.node == nil
}

func (l List[T]) Size() int {
	return l.size
}

func (l List[T]) Iterate() *Iterator[T] {
	return &Iterator[T]{l}
}

type Iterator[T any] struct {
	list List[T]
}

func (i *Iterator[T]) Size() int {
	return i.list.size
}

func (i *Iterator[T]) Next() (T, bool) {
	if i.list.node == nil {
		var t T
		return t, false
	}

	value := i.list.node.value
	i.list = i.list.node.parent
	return value, true
}

func (i *Iterator[T]) Collect() []T {
	values := make([]T, i.Size())
	for j := len(values) - 1; j >= 0; j-- {
		value, ok := i.Next()
		if !ok {
			break
		}

		values[j] = value
	}
	return values
}

func (i *Iterator[T]) CollectReverse() []T {
	values := make([]T, 0, i.Size())
	for {
		value, ok := i.Next()
		if !ok {
			break
		}

		values = append(values, value)
	}
	return values
}
