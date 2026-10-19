package api

import (
	"context"
	"net/url"
	"strconv"
)

// Page : enveloppe des listes paginées (/services, /invoices, /tickets).
type Page struct {
	Items    []map[string]any `json:"items"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
	Total    int              `json:"total"`
	HasMore  bool             `json:"hasMore"`
}

// pageLimit : maximum accepté par l'API.
const pageLimit = 100

// ListAll : parcourt les pages jusqu'à max éléments (max <= 0 : tout).
func (c *Client) ListAll(ctx context.Context, path string, query url.Values, max int) ([]map[string]any, int, error) {
	q := url.Values{}
	for k, v := range query {
		q[k] = v
	}
	per := pageLimit
	if max > 0 && max < per {
		per = max
	}
	var all []map[string]any
	total := 0
	for page := 1; ; page++ {
		q.Set("page", strconv.Itoa(page))
		q.Set("limit", strconv.Itoa(per))
		var p Page
		if err := c.Get(ctx, path, q, &p); err != nil {
			return nil, 0, err
		}
		total = p.Total
		all = append(all, p.Items...)
		if max > 0 && len(all) >= max {
			return all[:max], total, nil
		}
		if !p.HasMore || len(p.Items) == 0 || page >= 1000 {
			return all, total, nil
		}
	}
}
