package zone

import "fmt"

const (
	ItemEquipFailed ErrorCodeZone = iota
)

//go:generate stringer -type=ErrorCodeZone
type ErrorCodeZone int

type Err struct {
	Code    ErrorCodeZone
	Message string
	Details ErrDetails
}

type ErrDetails map[string]interface{}

func (e Err) Error() string {
	if len(e.Details) > 0 {
		var res = e.Code.String() + " > "
		if e.Message != "" {
			res += e.Message + " > "
		}
		for k, v := range e.Details {
			res += fmt.Sprintf("%v = %v, ", k, v)
		}
		return res
	}

	if e.Message == "" {
		return e.Code.String()
	}

	return e.Code.String() + " > " + e.Message
}
