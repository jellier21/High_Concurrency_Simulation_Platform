package Resty

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
)

func CreateOrderTask(mslNo string, count int, client *resty.Client) error {
	resp, err := client.R().
		SetFormData(map[string]string{"msl_no": mslNo, "count": strconv.Itoa(count)}).
		Post("http://localhost:8081/api/task/create")
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("service 返回异常状态码: %d, body: %s", resp.StatusCode(), resp.Body())
	}

	return nil
}
