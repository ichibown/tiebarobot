package simsimi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	SIMSIMI_API = "http://app.simsimi.com/app/aicr/request.p"
	ANOTHER_API = "http://www.xiaojo.com/bot/chata.php?chat="
)

// 大坑，json字段首字母必须大写！
type SimResp struct {
	Result           int
	Sentence_link_id int
	Slang            bool
	Msg              string
	Sentence_resp    string
}

// 首字母大写，public方法
func Talk(message string) string {
	resp, err := http.PostForm(SIMSIMI_API,
		url.Values{
			"req":  {message},
			"uid":  {"36140266"},
			"lc":   {"ch"},
			"tz":   {"Asia%2FShanghai"},
			"ft":   {"0"},
			"os":   {"a"},
			"av":   {"1.0"},
			"vkey": {"30cfbb7a6d3b49029caa3c7679865860"},
			"type": {"tta"}})
	if err != nil {
		return "Error When Talking To Simsimi !"
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Error Analyzing Simsimi's Reply"
	}
	simResp := new(SimResp)
	json.Unmarshal(content, simResp)
	return simResp.Sentence_resp
}

func Talk2(message string) string {
	response, err := http.Get(ANOTHER_API + message)
	if err != nil {
		return "Error"
	}
	defer response.Body.Close()
	result, _ := ioutil.ReadAll(response.Body)
	return string(result)
}
