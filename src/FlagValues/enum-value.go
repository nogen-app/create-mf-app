package flagvalues

import (
	"fmt"
	"strings"
)

//EnumValue is used by the cli flag parser to allow certain values to flag
type EnumValue struct {
	Enum     []string
	Default  string
	Selected string
}

//Set sets the value if its an allowed enum, otherwise it returns error
func (e *EnumValue) Set(value string) error {
	for _, enum := range e.Enum {
		if enum == value {
			e.Selected = value
			return nil
		}
	}

	return fmt.Errorf("Allowed values are %s", strings.Join(e.Enum, ", "))
}

//String returns the selected enum value as a string
func (e EnumValue) String() string {
	if e.Selected == "" {
		return e.Default
	}
	return e.Selected
}
