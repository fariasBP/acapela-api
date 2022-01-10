package config

// type Error interface {
// 	Error() string
// }

type (
	Res struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	ResErr struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Err  string `json:"err"`
	}
	resJson struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}
)

func SetRes(code int, msg string) *Res {
	d := &Res{
		Code: code,
		Msg:  msg,
	}
	return d
}

func SetResError(code int, msg string, err string) *ResErr {
	d := &ResErr{
		Code: code,
		Msg:  msg,
		Err:  err,
	}
	return d
}

func SetResJson(code int, msg string, json interface{}) *resJson {
	d := &resJson{
		Code: code,
		Msg:  msg,
		Data: json,
	}
	return d
}
