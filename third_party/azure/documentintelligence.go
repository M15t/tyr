package azure

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"tyr/internal/types"

	contextutil "tyr/internal/api/context"

	"gorm.io/datatypes"
)

// AnalyzeDocument get payload then send to Azure Document Intelligence
func (s *Service) AnalyzeDocument(c contextutil.Context, modelID, apiVersion string, payload io.Reader) (*ResponseHeaders, error) {
	start := time.Now()
	url := s.cfg.Endpoint + "/formrecognizer/documentModels/" + modelID + ":analyze?api-version=" + apiVersion

	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Ocp-Apim-Subscription-Key", s.cfg.Secret)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	reqHeader, err := parseHeadersToMap(req.Header)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusAccepted {
		data := new(ErrorResponse)
		json.Unmarshal(resData, &data)

		if err := s.repo.ActivityLog.Create(c.GetContext(), &types.ActivityLog{
			RequestURL:     url,
			RequestMethod:  method,
			RequestHeaders: reqHeader,
			ResponseCode:   res.StatusCode,
			ResponseBody:   datatypes.JSON(resData),
			DurationMS:     time.Since(start).Milliseconds(),
			IPAddress:      c.RealIP(),
		}); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf(data.Error.Message)
	}

	resHeaders, err := parseHeadersToMap(res.Header)
	if err != nil {
		return nil, err
	}

	if err := s.repo.ActivityLog.Create(c.GetContext(), &types.ActivityLog{
		RequestURL:      url,
		RequestMethod:   method,
		RequestHeaders:  reqHeader,
		ResponseCode:    res.StatusCode,
		ResponseHeaders: resHeaders,
		DurationMS:      time.Since(start).Milliseconds(),
		IPAddress:       c.RealIP(),
	}); err != nil {
		return nil, err
	}

	return toResponseHeaders(res.Header), nil
}

// GetAnalyzeDocument get request id the request the result of document
func (s *Service) GetAnalyzeDocument(c contextutil.Context, url string) (*ResultAnalyzeResponse, error) {
	start := time.Now()
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Ocp-Apim-Subscription-Key", s.cfg.Secret)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	reqHeaders, err := parseHeadersToMap(req.Header)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		data := new(ErrorResponse)
		json.Unmarshal(resData, &data)

		if err := s.repo.ActivityLog.Create(c.GetContext(), &types.ActivityLog{
			RequestURL:     url,
			RequestMethod:  method,
			RequestHeaders: reqHeaders,
			ResponseCode:   res.StatusCode,
			ResponseBody:   datatypes.JSON(resData),
			DurationMS:     time.Since(start).Milliseconds(),
			IPAddress:      c.RealIP(),
		}); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf(data.Error.Message)
	}

	data := new(ResultAnalyzeResponse)
	json.Unmarshal(resData, &data)

	if err := s.repo.ActivityLog.Create(c.GetContext(), &types.ActivityLog{
		RequestURL:     url,
		RequestMethod:  method,
		RequestHeaders: reqHeaders,
		ResponseCode:   res.StatusCode,
		ResponseBody:   datatypes.JSON(resData),
		DurationMS:     time.Since(start).Milliseconds(),
		IPAddress:      c.RealIP(),
		APIMRequestID:  getAPIMRequestID(url),
	}); err != nil {
		return nil, err
	}

	return data, nil
}
