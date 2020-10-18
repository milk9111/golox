package references

type FunctionType int

const (
	None FunctionType = iota
	Function
	Method
	Klass
	Property
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
	}

	return "Variable"
}
