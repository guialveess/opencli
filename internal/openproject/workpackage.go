package openproject

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WorkPackageListResponse struct {
	Total    int `json:"total"`
	Count    int `json:"count"`
	Embedded struct {
		Elements []WorkPackage `json:"elements"`
	} `json:"_embedded"`
	Links struct {
		Next *struct {
			Href string `json:"href"`
		} `json:"nextByOffset"`
	} `json:"_links"`
}

type WorkPackage struct {
	ID          int    `json:"id"`
	LockVersion int    `json:"lockVersion"`
	Subject     string `json:"subject"`
	Description struct {
		Raw string `json:"raw"`
	} `json:"description"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Links     struct {
		Status struct {
			Title string `json:"title"`
		} `json:"status"`
		Type struct {
			Title string `json:"title"`
		} `json:"type"`
		Priority struct {
			Title string `json:"title"`
		} `json:"priority"`
		Assignee struct {
			Title string `json:"title"`
		} `json:"assignee"`
	} `json:"_links"`
}

type WorkPackagePage struct {
	Items       []WorkPackage
	Total       int
	Page        int
	PageSize    int
	TotalPages  int
	HasNextPage bool
}

func (c *Client) ListWorkPackages(page, pageSize int) (*WorkPackagePage, error) {

	all, err := c.ListAllWorkPackages()
	if err != nil {
		return nil, err
	}

	total := len(all)
	totalPages := (total + pageSize - 1) / pageSize

	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	return &WorkPackagePage{
		Items:       all[start:end],
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		HasNextPage: page < totalPages,
	}, nil
}

func (c *Client) ListAllWorkPackages() ([]WorkPackage, error) {

	path := fmt.Sprintf("/api/v3/projects/%s/work_packages?pageSize=500", c.Project)

	req, err := c.newRequest(http.MethodGet, path)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result WorkPackageListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Embedded.Elements, nil
}

func (c *Client) GetWorkPackage(id int) (*WorkPackage, error) {
	path := fmt.Sprintf("/api/v3/work_packages/%d", id)

	req, err := c.newRequest(http.MethodGet, path)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("work package #%d não encontrado", id)
	}

	var wp WorkPackage
	if err := json.NewDecoder(resp.Body).Decode(&wp); err != nil {
		return nil, err
	}

	return &wp, nil
}

type CreateWorkPackageRequest struct {
	Subject     string
	Description string
	Type        string // opcional: Task, Bug, Feature, etc.
}

type CreateWorkPackageResponse struct {
	ID      int    `json:"id"`
	Subject string `json:"subject"`
}

func (c *Client) CreateWorkPackage(req *CreateWorkPackageRequest) (*CreateWorkPackageResponse, error) {
	path := fmt.Sprintf("/api/v3/projects/%s/work_packages", c.Project)

	payload := map[string]interface{}{
		"subject": req.Subject,
		"description": map[string]string{
			"format": "markdown",
			"raw":    req.Description,
		},
	}

	if req.Type != "" {
		payload["_links"] = map[string]interface{}{
			"type": map[string]string{
				"href": fmt.Sprintf("/api/v3/types?name=%s", req.Type),
			},
		}
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := c.newRequest(http.MethodPost, path)
	if err != nil {
		return nil, err
	}

	httpReq.Body = io.NopCloser(bytes.NewReader(payloadBytes))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("falha ao criar work package: %s (status %d)", string(body), resp.StatusCode)
	}

	var result CreateWorkPackageResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) AssignTaskForMe(id int, assigneeID int) error {
	wp, err := c.GetWorkPackage(id)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/api/v3/work_packages/%d", id)

	payload := map[string]interface{}{
		"lockVersion": wp.LockVersion,
		"_links": map[string]interface{}{
			"assignee": map[string]string{
				"href": fmt.Sprintf("/api/v3/users/%d", assigneeID),
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := c.newRequest(http.MethodPatch, path)
	if err != nil {
		return err
	}

	req.Body = io.NopCloser(bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to assign work package #%d to user %d", id, assigneeID)
	}

	return nil
}
