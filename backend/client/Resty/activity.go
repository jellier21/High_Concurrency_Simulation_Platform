package Resty

import (
	"client/Redis"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"log"
)

func GetActivityMessageInService(client *resty.Client) ([]Redis.ActivityMsg, error) {
	resp, err := client.R().Get("http://localhost:8081/api/activity/get")
	if err != nil {
		return nil, err
	}

	// 临时结构体匹配服务端返回格式
	type ServiceActivity struct {
		MslUid    string `json:"msl_uid"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}

	// 外层包装结构体，匹配服务端返回格式
	type ServiceResponse struct {
		Message string            `json:"message"`
		Data    []ServiceActivity `json:"data"`
	}

	var respBody ServiceResponse
	err = json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		log.Printf("解析活动信息失败: %v", err)
		return nil, err
	}
	activities := respBody.Data

	// 转换为客户端期望的格式
	var result []Redis.ActivityMsg
	for _, act := range activities {
		result = append(result, Redis.ActivityMsg{
			Uid:       act.MslUid,
			StartTime: act.StartTime,
			EndTime:   act.EndTime,
		})
	}

	return result, nil
}

func CreateActivityMessageInService(start, end string, client *resty.Client) error {
	resp, err := client.R().
		SetFormData(map[string]string{"start": start, "end": end}).
		Post("http://localhost:8081/api/activity/create")
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return err
	}
	return nil
}
