package config

import "fmt"

const ID_USER = "id"

func GetSkipResults(count, limit, page int) (skip int) {

	res := count / limit
	mod := count % limit

	fmt.Println(res)
	fmt.Println(mod)

	if mod == 0 {
		return limit * (page - 1)
	}

	return 1

}
