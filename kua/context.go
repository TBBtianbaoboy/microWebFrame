package kua

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type H map[string]interface{}  //for JSON

type Context struct {   //for w,r
	W http.ResponseWriter
	R *http.Request
	Path string
	Method string
	Params map[string]string //路由参数
	Sc int
	Body io.ReadCloser
}

//private for create Context object
func newContext(w http.ResponseWriter,r *http.Request) *Context{
	return &Context{
		W : w,
		R : r,
		Path: r.URL.Path,
		Method: r.Method,
		Body: r.Body,
	}
}

//获取路由参数
func (c *Context)Param(key string) string {
	value,_ := c.Params[key]
	return value
}

//read form -------------where?
func (c *Context)PostForm(key string) string {
	return c.R.FormValue(key)
}

func (c *Context)Query(key string) string {
	return c.R.URL.Query().Get(key)
}

//set response header statuscode
func (c *Context)Status (code int){
	c.Sc = code
	c.W.WriteHeader(code)//do in status 
}

//set response header 
func (c *Context)SetHeader (key string,value string){
	c.W.Header().Set(key,value)
}


func (c *Context)String(code int,format string,values ...interface{}){
	c.SetHeader("Content-Type","text/plain")
	c.Status(code)
	c.W.Write([]byte(fmt.Sprintf(format,values...)))
}

func (c *Context)JSON(code int,obj interface{}){
	c.SetHeader("Content-Type","application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.W)
	if err := encoder.Encode(obj);err != nil {
		WriteErrToLogFile(2,err)
		http.Error(c.W,err.Error(),500)
	}
}

func (c *Context)Data(code int,data []byte){
	c.Status(code)
	c.W.Write(data)
}

func (c *Context)HTML(code int,html string){
	c.SetHeader("Content-Type","text/html")
	c.Status(code)
	c.W.Write([]byte(html))
}













































