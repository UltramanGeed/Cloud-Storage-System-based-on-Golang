package handler

import (
	dblayer "Cloud-Storage-System-based-on-Golang/filestore-server/db"
	"Cloud-Storage-System-based-on-Golang/filestore-server/util"
	"fmt"
	"io/ioutil"

	//"ioutil"
	"net/http"
	"time"
)

const (
	pwdSalt = "*#890"
)

//SignupHandler:处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Println("html")
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}

	fmt.Println("test1")
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("Invalid parameter"))
		return
	}

	fmt.Println("test2")

	encPasswd := util.Sha1([]byte(passwd + pwdSalt))
	suc := dblayer.UserSignup(username, encPasswd)
	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
}

// //SignInHandler:登录接口
// func SignInHandler(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	username := r.Form.Get("username")
// 	password := r.Form.Get("password")

// 	encPasswd := util.Sha1([]byte(password + pwd_salt))

// 	//1.校验用户名及密码
// 	pwdChecked := dblayer.UserSignin(username, encPasswd)
// 	if !pwdChecked {
// 		w.Write([]byte("FAILED"))
// 		return
// 	}

// 	//2.生成访问凭证（token）
// 	token := GenToken(username)
// 	upRes := dblayer.UpdateToken(username, token)
// 	if !upRes {
// 		w.Write([]byte("FAILED"))
// 		return
// 	}
// 	//3.登录成功后重定向到首页
// 	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
// 	resp := util.RespMsg{
// 		Code: 0,
// 		Msg:  "OK",
// 		Data: struct {
// 			Location string
// 			Username string
// 			Token    string
// 		}{
// 			Location: "http://" + r.Host + "/static/view/home.html",
// 			Username: username,
// 			Token:    token,
// 		},
// 	}
// 	w.Write(resp.JSONBytes())
// }

// SignInHandler : 登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// data, err := ioutil.ReadFile("./static/view/signin.html")
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
		// w.Write(data)
		http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	encPasswd := util.Sha1([]byte(password + pwdSalt))

	// 1. 校验用户名及密码
	pwdChecked := dblayer.UserSignin(username, encPasswd)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}

	// 2. 生成访问凭证(token)
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED"))
		return
	}

	// 3. 登录成功后重定向到首页
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))

	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}

//UserInfoHandler:查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	//1.解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	// token := r.Form.Get("token")

	// //2.验证token是否有效
	// isValidToken := IsTokenValid(token)
	// if !isValidToken {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }

	//3.查询用户信息
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//4.组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

//GenToken: 生成token
func GenToken(username string) string {
	//md5(username + timestamp + token_salt)+timestamp[:8]凑够40位
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]

}

//IsTokenValid: token是否有效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	//TODO:判断token的时效性，是否过期
	//TODO:从数据库表tblusertoken查询username对应的token信息
	//TODO:对比两个token是否一致
	return true
}
