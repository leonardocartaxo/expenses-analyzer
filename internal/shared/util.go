package shared

// BaseFindCondition defines a flexible condition with a comparator
type BaseFindCondition struct {
	Field      string
	Comparator string
	Value      interface{}
}
