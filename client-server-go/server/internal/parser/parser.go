package parser

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	apperrors "github.com/SaiVikrantG/server/internal/errors"
	"github.com/SaiVikrantG/server/internal/models"
)

const maxBodySize = 1 << 20 // 1MB

func Parse(r io.Reader) (*models.Request, error) {
	reader := bufio.NewReader(r)
	req := &models.Request{}

	// Step 1: parse request line
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	requestLine = strings.TrimRight(requestLine, "\r\n")
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return nil, apperrors.NewBadRequest("malformed request line")
	}

	req.HTTPMethod = parts[0]
	req.Path = parts[1]
	req.HTTPVersion = parts[2]

	// Step 2: parse headers
	req.Headers = make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			return nil, apperrors.NewBadRequest("malformed header")
		}

		req.Headers[parts[0]] = parts[1]
	}

	// Step 3: parse body
	contentLengthStr, ok := req.Headers["Content-Length"]
	if !ok {
		return req, nil
	}

	contentLength, err := strconv.Atoi(contentLengthStr)
	if err != nil {
		return nil, apperrors.NewBadRequest("invalid Content-Length")
	}

	if contentLength > maxBodySize {
		return nil, apperrors.NewRequestTooLarge("request body too large")
	}

	body := make([]byte, contentLength)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		return nil, err
	}

	req.Body = body
	return req, nil
}
