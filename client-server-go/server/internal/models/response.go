package models

type Response struct {
	HTTPVersion   string
	StatusCode    int
	Status        string
	ContentType   string
	ContentLength int
	Connection    string
	ResponseBody  string
}
