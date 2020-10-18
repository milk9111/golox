package syntax

type (
	Stack struct {
		first  *node
		last   *node
		length int
	}
	node struct {
		next  *node
		prev  *node
		value interface{}
	}
)

func NewStack() *Stack {
	return &Stack{
		first:  nil,
		last:   nil,
		length: 0,
	}
}

func newNode(val interface{}) *node {
	return &node{
		next:  nil,
		prev:  nil,
		value: val,
	}
}

func (stack *Stack) Get(index int) interface{} {
	if stack.IsEmpty() {
		return nil
	}

	n := stack.first
	for i := 0; n != nil; i++ {
		if i == index {
			return n.value
		}

		n = n.next
	}

	return nil
}

func (stack *Stack) Len() int {
	return stack.length
}

func (stack *Stack) IsEmpty() bool {
	return stack.length == 0
}

func (stack *Stack) Push(val interface{}) {
	n := newNode(val)

	n.prev = stack.last
	if stack.last != nil {
		stack.last.next = n
	}
	stack.last = n

	if stack.IsEmpty() {
		stack.first = n
	}

	stack.length++
}

func (stack *Stack) Peek() interface{} {
	if stack.IsEmpty() || stack.last == nil {
		return nil
	}

	return stack.last.value
}

func (stack *Stack) Pop() interface{} {
	if stack.IsEmpty() || stack.last == nil {
		return nil
	}

	n := stack.last

	if stack.first == stack.last {
		stack.first = nil
	}

	stack.last = stack.last.prev
	stack.length--

	return n.value
}
