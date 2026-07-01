package response

type ResponseWriter interface {
	Header() map[string]string
	Write(statusCode int, body []byte) (int, error)
}
