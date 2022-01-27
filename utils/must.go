package utils

func Must(res interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}

	return res
}
