package main

import (
	"bytes"
	"github.com/clearcodecn/swaggos"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	oauth2AuthURL  = "https://www.oschina.net/action/oauth2/authorize"
	oauth2TokenURL = "http://localhost:9901/action/openapi/token"
)

var (
	oauth2Scopes = []string{
		"user_api",
		"blog_api",
		"news_api",
		"project_api",
		"tweet_api",
		"post_api",
		"user_mod_api",
		"comment_api",
		"favorite_api",
		"message_api",
		"search_api",
		"notice_api",
	}
)

// docker run -p 1001:8080 swaggerapi/swagger-ui
func main() {

	doc := swaggos.Default()
	// 本地代理(避免跨域)
	doc.HostInfo("localhost:9901", "/action/openapi")
	//doc.HostInfo("www.oschina.net", "/action")
	//oauth2 配置，ui界面存在跨域
	//tokenURL: "https://www.oschina.net/action/openapi/token"
	doc.Oauth2AccessCode(oauth2AuthURL, oauth2TokenURL, oauth2Scopes)

	doc.JWT("access_token")
	doc.Response(400, new(ErrorResp))
	doc.Query("access_token", "token", true)
	doc.Get("/user").Tag("user").Summary("个人信息").JSON(new(User))

	// https://www.oschina.net/openapi/docs/my_information 个人主页详情
	doc.Get("user_information").Summary("个人主页详情").Tag("user").
		Query("user", swaggos.Attribute{
			Type:        swaggos.Integer,
			Description: "查询用户id",
			Required:    true,
		}).
		Query("friend", swaggos.Attribute{
			Type:        swaggos.Integer,
			Description: "被查询用户id（friend和friend_name必须存在一个）",
			Required:    true,
		}).
		Query("friend_name", swaggos.Attribute{
			Type:        swaggos.String,
			Description: "被查询用户ident或名称",
			Required:    false,
		}).
		JSON(new(UserInformation))

	// https://www.oschina.net/openapi/docs/friends_list
	doc.Get("friends_list").Summary("获取好友列表").Tag("user").
		QueryObject(new(FriendListReq)).JSON(new(UserListItem))

	data, _ := doc.Build()
	yml, _ := doc.Yaml()

	// ui 代理
	go func() {
		h := new(oscProxy)
		http.ListenAndServe(":9901", h)
	}()

	http.HandleFunc("/yml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(yml)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(data)
	})
	http.ListenAndServe(":9900", nil)
}

type oscProxy struct{}

func (*oscProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}

	data, _ := ioutil.ReadAll(r.Body)
	r.URL.Scheme = "https"
	r.URL.Host = "www.oschina.net"
	req, err := http.NewRequest(r.Method, r.URL.String(), bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)
		return
	}
	for k := range r.Header {
		req.Header.Set(k, r.Header.Get(k))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	for k := range resp.Header {
		w.Header().Set(k, resp.Header.Get(k))
	}
	io.Copy(w, resp.Body)
}

/**
{
    id: 899**,
    email: "****@gmail.com",
    name: "彭博",
    gender: "male",
    avatar: "http://www.oschina.net/uploads/user/****",
    location: "广东 深圳",
    url: "http://home.oschina.net/****"
}
*/

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Avatar   string `json:"avatar"`
	Location string `json:"location"`
	URL      string `json:"url"`
}

type ErrorResp struct {
	Error   string `json:"error"`
	ErrDesc string `json:"error_description"`
}

/**
{
    "uid": 20,
    "name": "xxx",
    "ident": "xxx",
    "gender": 1,
    "relation": 3,
    "province": "上海",
    "city": "闵行",
    "platforms": [
        "Java EE",
        "PHP",
        ".NET/C#",
        "JavaScript",
        "Delphi/Pascal"
    ],
    "expertise": [
        "WEB开发",
        "服务器端开发",
        "DBA/数据库"
    ],
    "joinTime": "2008-09-18 09:17:15.0",
    "lastLoginTime": "2012-03-13 15:22:58.0",
    "portrait": "http://www.oschina.net/uploads/user/0/20_50.jpg",
    "notice": {
        "replyCount": 0,
        "msgCount": 0,
        "fansCount": 0,
        "referCount": 0
    }
}
*/

type Notice struct {
	ReplyCount int `json:"replyCount"`
	MsgCount   int `json:"msgCount"`
	FansCount  int `json:"fansCount"`
	ReferCount int `json:"referCount"`
}

type UserInformation struct {
	Uid           int      `json:"uid"`
	Name          string   `json:"name"`
	Gender        int      `json:"gender"`
	Relation      int      `json:"relation"`
	Province      string   `json:"province"`
	City          string   `json:"city"`
	Platforms     []string `json:"platforms"`
	Expertise     []string `json:"expertise"`
	JoinTime      string   `json:"joinTime"`
	LastLoginTime string   `json:"lastLoginTime"`
	Portrait      string   `json:"portrait"`
	Notice        Notice   `json:"notice"`
}

type FriendListReq struct {
	Page     int `json:"page" required:"true"`
	PageSize int `json:"pageSize" required:"true"`
	Relation int `json:"relation" enum:"0,1" description:"0-粉丝|1-关注的人" required:"true"`
}

/**
{
    "userList": [
        {
            "expertise": "<无>",
            "name": "test33",
            "userid": 253469,
            "gender": 1,
            "portrait": "http://static.oschina.org/uploads/user/126/253469_100.jpg?t=1366257509000"
        }
    ],
    "notice": {
        "replyCount": 0,
        "msgCount": 0,
        "fansCount": 0,
        "referCount": 0
    }
}
*/

type UserListItem struct {
	Expertise string `json:"expertise"`
	Name      string `json:"name"`
	Userid    int    `json:"userid"`
	Gender    int    `json:"gender"`
	Portrait  string `json:"portrait"`
	Notice    Notice `json:"notice"`
}
