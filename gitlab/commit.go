package gitlab

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Commit struct {
	ID          string `json:"id"`
	ShortID     string `json:"short_id"`
	Title       string `json:"title"`
	Message     string `json:"message"`
	AuthorName  string `json:"author_name"`
	AuthoredAt  string `json:"authored_date"`
	CommittedAt string `json:"committed_date"`
}

func GetCommit(projectID, sha, token string) (*Commit, error) {
	u := fmt.Sprintf("%s/projects/%s/repository/commits/%s",
		baseURL, url.PathEscape(projectID), url.PathEscape(sha))

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("PRIVATE-TOKEN", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("commit request failed: %s", resp.Status)
	}

	var c Commit
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
