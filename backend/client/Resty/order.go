package Resty

import (
	"client/model"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
)

type TaskDB struct {
	model.TaskDB
}
type TaskResponse struct {
	Task TaskDB `json:"task"`
}

func CreateOrderTask(mslNo string, count int, client *resty.Client) (TaskDB, error) {
	resp, err := client.R().
		SetFormData(map[string]string{"msl_no": mslNo, "count": strconv.Itoa(count)}).
		Post("http://localhost:8081/api/task/create")
	if err != nil {
		return TaskDB{}, err
	}
	if resp.StatusCode() != 200 {
		return TaskDB{}, fmt.Errorf("service 返回异常状态码: %d, body: %s", resp.StatusCode(), resp.Body())
	}
	var respBody TaskResponse
	err = json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		return TaskDB{}, err
	}
	return respBody.Task, nil
}

func GetOrderTaskInformationBySQL(client *resty.Client) ([]model.TaskDB, error) {
	resp, err := client.R().
		Get("http://localhost:8081/api/task/get")
	if err != nil {
		return []model.TaskDB{}, err
	}
	var respBody []model.TaskDB
	err = json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		return []model.TaskDB{}, err
	}
	return respBody, nil
}

func GetOrderInformationByUuid(uuid string, client *resty.Client) (model.OrderDB, error) {
	resp, err := client.R().
		SetQueryParam("uuid", uuid).
		Get("http://localhost:8081/api/order/get")
	if err != nil {
		return model.OrderDB{}, err
	}
	var respBody model.OrderDB
	err = json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		return model.OrderDB{}, err
	}
	return respBody, nil
}
