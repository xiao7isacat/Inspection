package checkctl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Post(url string, body []byte) error {

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
		}
		return fmt.Errorf("请求响应报错:%v\t报错内容为:%v", resp.StatusCode, string(body))
	}
	return nil
}

func Put(url string, body []byte) error {

	// 创建请求
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
		}
		return fmt.Errorf("请求响应报错:%v\t报错内容为:%v", resp.StatusCode, string(body))
	}
	return nil
}

func Get(url string) error {
	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("请求响应报错:%v\t报错内容为:%v", resp.StatusCode, string(body))
	}

	fmt.Printf("Body: %s\n", string(body))
	return nil
}
