package pfa

type pfaError struct {
	pfa *pfaStruct
	// errVal interface{}
	err   error
	Where string
}

func (err pfaError) Error() string {
	return err.err.Error()
}

func (v pfaError) DataForJSON() interface{} {
	result := map[string]interface{}{}
	result["pfa"] = v.pfa.DataForJSON()
	// result["errVal"] = v.errVal
	result["err"] = v.err
	result["Where"] = v.Where
	return result
}
