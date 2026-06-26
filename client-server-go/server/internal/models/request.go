package models

type Request struct {
	HTTPMethod  string
	FilePath    string
	HTTPVersion string
	Host        string
	Connection  string
}
