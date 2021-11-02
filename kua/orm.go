package kua

import (
	//"encoding/json"

	"errors"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type User struct {
	Id     bson.ObjectId `bson:"_id"`
	Name   string        `bson:"username"`
	Passwd string        `bson:"password"`
	Auth   string        `bason:"auth"`
}

//部分封装--------ok
func conMgo() (*mgo.Session,error) {
	session,err := mgo.Dial("172.17.0.2:27017")
	return session,err
}

func conMgoDb(database string,collection string,session *mgo.Session) *mgo.Collection{
	c := session.DB(database).C(collection)
	return c
}

func findData(c *mgo.Collection,user H) (User,error) {
	res := User{}
	err := c.Find(bson.M{"username":user["username"]}).One(&res)
	return res,err
}

//------ok
func insertData(c *mgo.Collection,user H) error {
	err := c.Insert(user)
	return err
}

func deleteData(c *mgo.Collection,user H) bool {
	err := c.Remove(H{"username":user["username"]})
	if err != nil {
		WriteErrToLogFile(1,err)
		return false
	}
	return true
}

//登陆验证处理函数
func FindMongo(user H,ver V) H {
	//buf, _ := json.Marshal(user)
	//println(string(buf))
	mes := H{"code":0,"err":"","msg":"false","verify":false}
	session,err := conMgo()
	if err!= nil {
		WriteErrToLogFile(1,err)
		mes["err"] = err.Error()
		return mes
	}
	defer session.Close()
	c := conMgoDb("mydb","user",session)
	res,err := findData(c,user)
	if err != nil {
		WriteErrToLogFile(1,err)
		mes["err"] = err.Error()
	} else if res.Passwd != user["password"] {
		mes["err"] = "password error"
	} 
	if mes["err"] == "" {
		if ver["capId"] != "" && ver["bas64"] != "" && captchaVerifyHandle(ver) {
			passInfo := H{"username":user["username"],"password":user["password"],"auth":CreateRandomString(15)}
			c = conMgoDb("mydb","login",session)
			_,err = findData(c,passInfo)
			if err == nil {
				// WriteErrToLogFile(1,err)
				mes["err"] = "user has logined in"
			}else{
				err = insertData(c,passInfo)
				mes = H{"code":1,"auth":passInfo["auth"],"err":err,"msg":"success","verify":true}
			}
		}else{
			mes["err"] = "no verify or verify error"
		}
	}
	return mes
}

// 注册处理函数
func InsertMongo(nu H) H {
	mes := H{"code":0,"err":"","msg":"failed"}
	session,err := conMgo()
	if err != nil {
		WriteErrToLogFile(1,err)
		mes["err"] = err.Error()
	}
	defer session.Close()
	c := conMgoDb("mydb","user",session)
	_,err = findData(c,nu)
	if err == nil {
		mes["err"] = "repeated"
	}	
	if err = insertData(c,nu);err != nil {
		WriteErrToLogFile(1,err)
		mes["err"] = err.Error()
	}
	if mes["err"]== "" {
		mes = H{"code":1,"username":nu["username"],"password":"******","msg":"success"}
	} 
	return mes
}

//注销处理函数
func DeleteMongo(username H) H {
	mes := H{"code":0,"err":"","msg":"failed"}
	session,err := conMgo()
	if err != nil {
		WriteErrToLogFile(1,err)
		mes["err"] = err.Error()
	}
	defer session.Close()
	c := conMgoDb("mydb","login",session)
	res,err := findData(c,username)
	if err != nil {
		err = errors.New("user unlogined in")
		WriteErrToLogFile(1,err)
		mes["err"] = "user unlogined in"
	}else{
		if username["auth"] != res.Auth {
			mes["err"] = "you have no authority"
		}else{
			deleteData(c,username)
			mes = H{"code":1,"data":"goodbye honey","msg":"success"}
		}
	} 
	return mes
}
