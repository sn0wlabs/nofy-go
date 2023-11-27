package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

const (
	// API Version
	API_Version string = "0.0.1"

	// API URL
	API_URL string = "https://api.nofy.io"
)

type Backend struct {
	// encoded API Key
	Key string
}

// params is not used yet, can be implemented later for header params
func (c *Backend) NewRequest(method, path, contentType string, params interface{}) (*http.Request, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = API_URL + path

	// Body is set later by `Do`.
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return nil, err
	}

	authorization := "Basic " + c.Key

	req.Header.Add("Authorization", authorization)
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Nofy-Version", API_Version)

	ua := "nofy-go/" + API_Version

	req.Header.Add("User-Agent", ua)

	return req, nil
}

func (c *Backend) Do(req *http.Request, body *bytes.Buffer, v interface{}) (*http.Response, error) {
	// In Stripes API documentation they say that there are problems
	// with http2 transports but we're living a risky life so we're
	// going to use it anyway.
	//
	// (17.11.2023)
	client := &http.Client{}

	if body != nil {
		// We can safely reuse the same buffer that we used to encode our body,
		// but return a new reader to it everytime so that each read is from
		// the beginning.
		reader := bytes.NewReader(body.Bytes())

		req.Body = nopReadCloser{reader}

		// And also add the same thing to `Request.GetBody`, which allows
		// `net/http` to get a new body in cases like a redirect. This is
		// usually not used, but it doesn't hurt to set it in case it's
		// needed. See:
		//
		//     https://github.com/stripe/stripe-go/issues/710
		//
		req.GetBody = func() (io.ReadCloser, error) {
			reader := bytes.NewReader(body.Bytes())
			return nopReadCloser{reader}, nil
		}
	}

	return client.Do(req)
}

// UnmarshalJSONVerbose unmarshals JSON, but in case of a failure logs and
// produces a more descriptive error. (also stolen froms stripe)
func (c *Backend) UnmarshalJSONVerbose(statusCode int, body []byte, v interface{}) error {
	err := json.Unmarshal(body, v)
	if err != nil {
		// If we got invalid JSON back then something totally unexpected is
		// happening (caused by a bug on the server side). Put a sample of the
		// response body into the error message so we can get a better feel for
		// what the problem was.
		bodySample := string(body)
		if len(bodySample) > 500 {
			bodySample = bodySample[0:500] + " ..."
		}

		// Make sure a multi-line response ends up all on one line
		bodySample = strings.Replace(bodySample, "\n", "\\n", -1)

		newErr := fmt.Errorf("Couldn't deserialize JSON (response status: %v, body sample: '%s'): %v",
			statusCode, bodySample, err)
		return newErr
	}

	return nil
}

func (c *Backend) Call(method, path string, params, v interface{}) *Error {
	var body *bytes.Buffer

	if params != nil {
		body = &bytes.Buffer{}

		encoder := json.NewEncoder(body)
		encoder.SetEscapeHTML(false)

		err := encoder.Encode(params)
		if err != nil {
			return &ParamEncodeError
		}
	}

	req, err := c.NewRequest(method, path, "application/json", params)
	if err != nil {
		return &RequestError
	}

	return c.rawCall(req, v, body)
}

func (c *Backend) CallMultipart(method, path, fileName string, file *os.File, v interface{}) *Error {
	var buf bytes.Buffer

	// Create a new multipart writer with the buffer
	w := multipart.NewWriter(&buf)

	fw, err := w.CreateFormFile("file", fileName)
	if err != nil {
		return &RequestError
	}

	// Copy the contents of the file to the form field
	if _, err := io.Copy(fw, file); err != nil {
		return &RequestError
	}

	// Close the multipart writer to finalize the request
	w.Close()

	req, err := c.NewRequest(method, path, w.FormDataContentType(), nil)
	if err != nil {
		return &RequestError
	}

	return c.rawCall(req, v, &buf)
}

func (c *Backend) rawCall(req *http.Request, v interface{}, body *bytes.Buffer) *Error {
	res, err := c.Do(req, body, v)
	if err != nil {
		return &RequestError
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return &Error{
			Code:    "readingBodyError",
			Message: fmt.Sprintf("Error reading body - statuscode: %v", res.StatusCode),
		}
	}

	if res.StatusCode >= 400 {
		return &Error{
			Code:    "statusCodeError",
			Message: fmt.Sprintf("Statuscode: %v", res.StatusCode),
		}
	}

	if v == nil {
		return &NoResponseTypeError
	}

	if err := c.UnmarshalJSONVerbose(res.StatusCode, resBody, v); err != nil {
		return &Error{
			Code:    "decodeError",
			Message: fmt.Sprintf("Error decoding response: %v", err.Error()),
		}
	}

	return nil
}
