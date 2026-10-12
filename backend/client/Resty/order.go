package Resty

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
)

type TaskDB struct {
	ID         int
	Uuid       string `json:"uuid"`
	MslNo      string `json:"msl_no"`
	Status     string `json:"status"` //Success|Interrupt|Running
	Total      int64  `json:"total"`
	SuccessInt int64  `json:"success_int"`
	FailInt    int64  `json:"fail_int"`
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
