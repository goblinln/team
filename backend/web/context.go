package web

import (
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Context hold request information and response object to send
// message back to client.
type Context struct {
	Session *Session

	req         *http.Request
	rsp         *Responser
	routeParams map[string]string
	queryParams url.Values
}

// Method returns the request method
func (c *Context) Method() string {
	return c.req.Method
}

// URL returns request URL information.
func (c *Context) URL() *url.URL {
	return c.req.URL
}

// RemoteIP returns IP address of remote client
func (c *Context) RemoteIP() string {
	if ip := c.RequestHeader().Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ", ")[0]
	}

	if ip := c.RequestHeader().Get("X-Real-IP"); ip != "" {
		return ip
	}

	ip, _, _ := net.SplitHostPort(c.req.RemoteAddr)
	return ip
}

// RemoteAddr returns addr of remote client.
func (c *Context) RemoteAddr() string {
	return c.req.RemoteAddr
}

// Body returns request body
func (c *Context) Body() io.ReadCloser {
	return c.req.Body
}

// RequestHeader returns HTTP header map of request.
func (c *Context) RequestHeader() http.Header {
	return c.req.Header
}

// ResponseHeader returns HTTP header map of reponse to write.
func (c *Context) ResponseHeader() http.Header {
	return c.rsp.Header()
}

// StartSession prepares session data.
func (c *Context) StartSession() {
	if c.Session == nil {
		c.Session = Sessions.Start(c)
	}
}

// EndSession clears session data.
func (c *Context) EndSession() {
	Sessions.End(c)
}

// Cookie returns named cookie value
func (c *Context) Cookie(key string) (*http.Cookie, error) {
	return c.req.Cookie(key)
}

// SetCookie adds a Set-Cookie header to response writer
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.rsp, cookie)
}

// Status of response
func (c *Context) Status() int {
	return c.rsp.Status()
}

// SetStatus for response
func (c *Context) SetStatus(code int) {
	c.rsp.WriteHeader(code)
}

// RouteValue returns named parameter value in route pattern
func (c *Context) RouteValue(name string) string {
	return c.routeParams[name]
}

// QueryValue returns named parameter value in URL query string
func (c *Context) QueryValue(name string) string {
	if c.queryParams == nil {
		c.queryParams = c.req.URL.Query()
	}

	return c.queryParams.Get(name)
}

// QueryArray returns named parameter array in URL query string
func (c *Context) QueryArray(name string) []string {
	if c.queryParams == nil {
		c.queryParams = c.req.URL.Query()
	}

	return c.queryParams[name]
}

// FormValue returns named parameter value in both URL and form
func (c *Context) FormValue(name string) string {
	return c.req.FormValue(name)
}

// FormArray returns named parameter array in both URL and form
func (c *Context) FormArray(name string) []string {
	c.req.FormValue(name) // Just make sure ParseForm has been called.
	return c.req.Form[name]
}

// PostFormValue returns named parameter value only from Form
func (c *Context) PostFormValue(name string) string {
	return c.req.PostFormValue(name)
}

// PostFormArray returns named parameter array only from Form
func (c *Context) PostFormArray(name string) []string {
	c.req.PostFormValue(name) // Make sure ParseForm has been called
	return c.req.PostForm[name]
}

// MultipartForm returns multipart form. Useful to deal with uploads.
func (c *Context) MultipartForm() *multipart.Form {
	if c.req.MultipartForm == nil {
		c.req.ParseMultipartForm(64 << 20)
	}

	return c.req.MultipartForm
}

// Redirect client.
func (c *Context) Redirect(url string) {
	http.Redirect(c.rsp, c.req, url, http.StatusPermanentRedirect)
}

// Stream response
func (c *Context) Stream(status int, contentType string, reader io.Reader) error {
	c.rsp.WriteHeader(status)
	c.rsp.Header().Set("Content-Type", contentType)

	_, err := io.Copy(c.rsp, reader)
	return err
}

// Flush the response writer. Useful when sending response using Stream mode.
func (c *Context) Flush() {
	c.rsp.Flush()
}

// Blob writes buffer of data into response.
func (c *Context) Blob(status int, contentType string, blob []byte) error {
	c.rsp.WriteHeader(status)
	c.rsp.Header().Set("Content-Type", contentType)

	_, err := c.rsp.Write(blob)
	return err
}

// String write string to response
func (c *Context) String(status int, str string) error {
	return c.Blob(status, "text/plain", []byte(str))
}

// HTML writes html page to response.
func (c *Context) HTML(status int, html string) error {
	return c.Blob(status, "text/html", []byte(html))
}

// JSON writes object to json string response.
func (c *Context) JSON(status int, v interface{}) error {
	blob, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return c.Blob(status, "application/json", blob)
}

// File forces client to download a file.
func (c *Context) File(status int, path string) error {
	return c.FileWithName(status, path, "")
}

// FileWithName sames like File but also renames then downloaded file.
func (c *Context) FileWithName(status int, path, name string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		return errors.New(path + " is NOT a file but directory")
	}

	if len(name) == 0 {
		name = info.Name()
	}

	c.rsp.Header().Set("Content-Disposition", "attachement;filename="+name)
	http.ServeContent(c.rsp, c.req, info.Name(), info.ModTime(), file)
	return nil
}
