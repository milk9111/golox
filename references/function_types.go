package references

type FunctionType int

const (
	None FunctionType = iota
	Function
	Method
	Klass
	Property
	Initializer
)

func GetFunctionTypeName(t FunctionType) string {
	switch t {
	case None:
		return "Variable"
	case Function:
		return "Function"
	case Method:
		return "Method"
	case Klass:
		return "Class"
	case Property:
		return "Property"
	case Initializer:
		return "Initializer"
	}

	return "Variable"
}
