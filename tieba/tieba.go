package tieba

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	TIMEOUT        = 10
	API_LOGIN      = "http://c.tieba.baidu.com/c/s/login"
	API_AT         = "http://c.tieba.baidu.com/c/u/feed/atme"
	API_REPLY_POST = "http://c.tieba.baidu.com/c/c/post/add"
	API_REPLY      = "http://c.tieba.baidu.com/c/u/feed/replyme"
	API_LIKE       = "http://c.tieba.baidu.com/c/c/forum/like"
	API_SIGN       = "http://c.tieba.baidu.com/c/c/forum/sign"
	API_FOLLOW     = "http://c.tieba.baidu.com/c/u/follow/page"
	API_TBS        = "http://tieba.baidu.com/dc/common/tbs"
	FORUM_HOME     = "http://m.tieba.com/f?kw="
)

var client *http.Client

func init() {
	client = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			ResponseHeaderTimeout: time.Second * TIMEOUT,
		},
	}
}

type User struct {
	Uid   string
	Tbs   string
	Bduss string
}

type GetAtJson struct {
	Error_code string
	Error_msg  string
	At_list    []AtNode
}

type AtNode struct {
	Is_floor   string
	Content    string
	Thread_id  string
	Post_id    string
	Fname      string
	Reply      string
	Retry      int
	Quote_user struct {
		Id   string
		Name string
	}
}

type TbsJson struct {
	Tbs      string
	Is_login int
}

type TiebaError struct {
	Error_code string
	Error_msg  string
}

func (e *TiebaError) Error() string {
	return fmt.Sprintf("Error Code: %v,Error Message %v", e.Error_code, e.Error_msg)
}

func doPost(url string, params url.Values, cookie string) (data []byte, err error) {
	request, err := http.NewRequest("POST", url, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", "BaiduTieba for Android 5.1.3")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	if cookie != "" {
		request.Header.Add("Cookie", cookie)
	}

	response, postErr := client.Do(request)

	if postErr != nil {
		return nil, postErr
	}
	defer response.Body.Close()

	content, parseErr := ioutil.ReadAll(response.Body)
	if parseErr != nil {
		return nil, parseErr
	}

	return content, nil
}

func GetFid(fname string) (fid string, err error) {
	response, err := client.Get(FORUM_HOME + fname)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	content, parseErr := ioutil.ReadAll(response.Body)
	if parseErr != nil {
		return "", parseErr
	}
	re := regexp.MustCompile(`fid\D+"(\d+)"`)
	result := re.FindStringSubmatch(string(content))
	return result[1], nil
}

// common post params
func addBasicInfo(params *url.Values) {
	params.Add("_client_type", "1")
	params.Add("_client_version", "4.5.3")
	params.Add("_phone_imei", "860806022136732")
	params.Add("from", "Ad_wandoujia")
	params.Add("net_type", "3")
}

// calculate sign
func addSign(params *url.Values) {
	str, err := url.QueryUnescape(params.Encode())
	if err != nil {
		return
	}
	str = strings.Replace(str, "&", "", -1) + "tiebaclient!!!"
	h := md5.New()
	h.Write([]byte(str))
	str = hex.EncodeToString(h.Sum(nil))
	params.Add("sign", str)
}

func (user *User) getTbs() {
	content, err := doPost(API_TBS, nil, "BDUSS="+user.Bduss)
	if err != nil {
		return
	}
	tbsJson := new(TbsJson)
	json.Unmarshal(content, tbsJson)

	if tbsJson.Is_login != 0 {
		user.Tbs = tbsJson.Tbs
	}
}

// Get At
func (user *User) GetAtMe() (atList []AtNode, err error) {
	params := url.Values{}
	params.Set("BDUSS", user.Bduss)
	params.Add("pn", "1")
	params.Add("uid", user.Uid)
	addBasicInfo(&params)
	addSign(&params)

	content, err := doPost(API_AT, params, "")
	if err != nil {
		return nil, err
	}

	getAtJson := new(GetAtJson)
	json.Unmarshal(content, getAtJson)

	if getAtJson.Error_code != "0" {
		return nil, &TiebaError{getAtJson.Error_code, getAtJson.Error_msg}
	} else {
		return getAtJson.At_list, nil
	}
}

// Reply a floor
func (user *User) ReplyFloor(qid, tid, content, fid, fname string) error {
	params := url.Values{}
	user.getTbs()
	params.Set("BDUSS", user.Bduss)
	params.Add("content", content)
	params.Add("fid", fid)
	params.Add("kw", fname)
	params.Add("quote_id", qid)
	params.Add("tbs", user.Tbs)
	params.Add("tid", tid)
	addBasicInfo(&params)
	addSign(&params)

	replyContent, replyErr := doPost(API_REPLY_POST, params, "")
	if replyErr != nil {
		return replyErr
	}

	result := new(TiebaError)
	json.Unmarshal(replyContent, result)
	if result.Error_code != "0" {
		return result
	} else {
		return nil
	}
}
