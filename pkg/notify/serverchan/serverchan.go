package serverchan

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	httpurl "net/url"

	"github.com/zktnotify/zktnotify/pkg/notify/typed"
)

const (
	NotifyHost = "https://sc.ftqq.com"
)

type ServerChan struct{}

func New() typed.Notifier {
	return &ServerChan{}
}

type responsed struct {
	ErrNo   int    `json:"errno"`
	ErrMsg  string `json:"errmsg"`
	DataSet string `json:"dataset"`
}

func (s *ServerChan) Notify(token string, msg string, receiver ...typed.Receiver) error {
	if s == nil || token == "" {
		return nil
	}

	url := URL(token)
	payload := httpurl.Values{}
	payload.Set("text", msg)
	payload.Set("desp", "# 测试啦")

	rsp, err := http.Get(url + "?" + payload.Encode())
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	response := responsed{}
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}
	if response.ErrNo != 0 {
		return errors.New(response.ErrMsg)
	}

	return nil
}

func URL(token string) string {
	return NotifyHost + "/" + token + ".send"
}
