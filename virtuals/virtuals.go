package virtuals

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GetPageResponse struct {
	Data            []Data `json:"data"`
	Page            int
	HasNextPage     bool
	HasPreviousPage bool
}

type Data struct {
	UID           string  `json:"uid"`
	Name          string  `json:"name"`
	Status        string  `json:"status"`
	Category      string  `json:"category"`
	McapInVirtual float64 `json:"mcapInVirtual"`
}

func GetPage(page int, minMcap int) (*GetPageResponse, error) {
	pageSize := 100
	fmt.Printf("getting page %d from virtuals (pageSize=%d)\n", page, pageSize)

	url := fmt.Sprintf("https://api.virtuals.io/api/virtuals?filters[status][$in][0]=AVAILABLE&filters[status][$in][1]=ACTIVATING&filters[priority][$ne]=-1&filters[mcapInVirtual][$gt]=%d&sort[0]=mcapInVirtual:asc&sort[1]=createdAt:desc&pagination[page]=%d&pagination[pageSize]=%d", minMcap, page, pageSize)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result struct {
		Data []Data `json:"data"`
		Meta struct {
			Pagination struct {
				Page      int `json:"page"`
				PageSize  int `json:"pageSize"`
				PageCount int `json:"pageCount"`
				Total     int `json:"total"`
			} `json:"pagination"`
		} `json:"meta"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &GetPageResponse{
		Data:            result.Data,
		Page:            page,
		HasNextPage:     page < result.Meta.Pagination.PageCount,
		HasPreviousPage: page >= 1,
	}, nil
}

func GetPrice() (float64, error) {
	fmt.Printf("getting virtuals price\n")

	url := "https://api.coingecko.com/api/v3/simple/price?ids=virtual-protocol%2Cethereum&vs_currencies=USD"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := client.Do(req)
	if err != nil {
		return -1, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, fmt.Errorf("failed to read response body: %w", err)
	}

	var result map[string]map[string]float64
	if err := json.Unmarshal(body, &result); err != nil {
		return -1, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if virtualProtocol, ok := result["virtual-protocol"]; ok {
		if usdPrice, ok := virtualProtocol["usd"]; ok {
			return usdPrice, nil
		}
	}

	return -1, fmt.Errorf("virtual-protocol price not found in response")
}
