package audience

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"o2o/common"
	"os"
	"time"
)

// Errors section
var (
	ErrTokenIsNotSet = errors.New("yandex audience token isn't set")
)

//constants
const (
	tokenVariable = "YANDEX_AUDIENCE_TOKEN"
	apiUrl        = "https://api-audience.yandex.ru/"
)

type Client struct {
	apiVersion string
	token      string
	hc         *http.Client
}

type YandexAudience struct {
	token string
}

func NewClient(ctx context.Context) (*Client, error) {
	var client Client
	//Get token from context or env
	if os.Getenv(tokenVariable) == "" {
		if tok, ok := ctx.Value(tokenVariable).(string); ok && tok != "" {
			client.token = tok
		} else {
			return nil, ErrTokenIsNotSet
		}
	} else {
		client.token = os.Getenv(tokenVariable)
	}
	//Creating http client
	client.hc = &http.Client{}
	//Set api version: now available only v1
	client.apiVersion = "v1"
	return &client, nil
}

func (c *Client) Close() error {
	c.hc = nil
	return nil
}

func (c *Client) SegmentsList(pixel ...int) ([]Segment, error) {
	requestPath := fmt.Sprintf("%s%s/management/segments", apiUrl, c.apiVersion)
	if len(pixel) > 0 {
		requestPath += fmt.Sprintf("?pixel=%d", pixel[0])
	}
	requestUrl, err := url.Parse(requestPath)
	if err != nil {
		return nil, err
	}
	resp, err := c.hc.Do(&http.Request{
		Method: http.MethodGet,
		URL:    requestUrl,
		Header: http.Header{"Authorization": {fmt.Sprintf("OAuth %s", c.token)}},
	})
	if err != nil {
		return nil, err
	}
	dt, err := ioutil.ReadAll(resp.Body)
	fmt.Println(requestUrl.String())
	fmt.Println(string(dt))
	return nil, nil
}

func (YandexAudience) New() *YandexAudience {
	ya := YandexAudience{}
	ya.token = os.Getenv("yandex_audience_token")
	return &ya
}

func (ya *YandexAudience) UploadMacs(reportname, fileName string) (id int64, err error) {
	body, statusCode, err := common.UploadFile("https://api-audience.yandex.ru/v1/management/segments/upload_csv_file?",
		fileName, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", ya.token)})

	if err != nil {
		err = errors.New(fmt.Sprintf("file upload error: %s", err.Error()))
		return
	}
	if statusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("bad status: %d, body: %s", statusCode, string(body)))
		return
	}
	response := struct {
		Segment struct {
			Id          int64  `json:"id"`
			Name        string `json:"name"`
			Hashed      bool   `json:"hashed"`
			ContentType string `json:"content_type"`
		} `json:"segment"`
	}{}

	err = json.Unmarshal(body, &response)

	if err != nil {
		return
	}

	response.Segment.Name = reportname
	response.Segment.ContentType = "mac"
	q, err := json.Marshal(response)
	if err != nil {
		return
	}
	//Сохраняем загруженый сегмент
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api-audience.yandex.ru/v1/management/segment/%d/confirm?", response.Segment.Id), bytes.NewBuffer(q))
	if err != nil {
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ya.token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	// Check the response
	if res.StatusCode != http.StatusOK {
		reqBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logrus.Warnf("can't parse unknown error yandex audience %s because %s", string(reqBody), err.Error())
			err = errors.New("внутренняя ошибка")
			return 0, err
		}
		var errResp struct {
			Errors  interface{} `json:"errors"`
			Code    int         `json:"code"`
			Message string      `json:"message"`
		}
		err = json.Unmarshal(reqBody, &errResp)
		if err != nil {
			logrus.Warnf("can't parse unknown error yandex audience %s because %s", string(reqBody), err.Error())
			err = errors.New("внутренняя ошибка")
			return 0, err
		}
		err = errors.New(errResp.Message)
		return 0, err
	}
	return response.Segment.Id, err
}

func (ya *YandexAudience) AppendUser(id int64, email string) (err error) {
	//Выдадим права
	type permission struct {
		Grant struct {
			UserLogin string    `json:"user_login"`
			CreatedAt time.Time `json:"created_at"`
			Comment   string    `json:"comment"`
		} `json:"grant"`
	}
	var perm permission
	perm.Grant.UserLogin = email
	perm.Grant.CreatedAt = time.Now()
	q, err := json.Marshal(perm)
	if err != nil {
		return
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("https://api-audience.yandex.ru/v1/management/segment/%d/grant?", id), bytes.NewBuffer(q))
	if err != nil {
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ya.token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	// Check the response
	if res.StatusCode != http.StatusOK {
		type YandexError struct {
			Errors []struct {
				ErrorType string `json:"error_type"`
				Message   string `json:"message"`
				Location  string `json:"location"`
			} `json:"errors"`
		}
		var errResponse YandexError
		reqBody, _ := ioutil.ReadAll(res.Body)
		err = json.Unmarshal(reqBody, &errResponse)
		if err != nil {
			return
		}
		if len(errResponse.Errors) > 0 {
			err = errors.New(errResponse.Errors[0].Message)
		} else {
			err = errors.New(fmt.Sprintf("bad status: %s, body: %s", res.Status, string(reqBody)))
		}
		return
	}
	return
}

func (ya *YandexAudience) RemoveSegment(id int64) (success bool, err error) {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("https://api-audience.yandex.ru/v1/management/segment/%d", id), nil)
	if err != nil {
		return
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ya.token))
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return
	}
	defer res.Body.Close()
	var response struct {
		Success bool `json:"success"`
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return
	}
	success = response.Success
	return
}

type YandexSegment struct {
	Id              int       `json:"id"`
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	CreateTime      time.Time `json:"create_time"`
	MatchedQuantity int       `json:"matched_quantity"`
}

type yandexAudienceType struct {
	Segments []YandexSegment `json:"segments"`
}

func (ya *YandexAudience) GetUploads() (uploads []YandexSegment, err error) {
	request, err := http.NewRequest("GET", "https://api-audience.yandex.ru/v1/management/segments?", nil)
	if err != nil {
		return
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ya.token))
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return
	}
	defer res.Body.Close()

	var segments yandexAudienceType

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(resBody, &segments)
	if err != nil {
		return
	}
	uploads = segments.Segments
	return
}
