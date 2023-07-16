package threads

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	baseURL = "https://www.threads.net/@"
	apiURL  = "https://www.threads.net/api/graphql"

	getUserDocID        = "23996318473300828"
	getUserThreadsDocID = "6232751443445612"
	getUserRepliesDocID = "6307072669391286"
	getPostDocID        = "5587632691339264"
	getLikersDocID      = "9360915773983802"
)

var userIDRegex = regexp.MustCompile(`"user_id":"(\d+)"`)

type Client struct {
	client *http.Client
	header http.Header
	token  string
}

type Option func(*Client)

func NewClient(ctx context.Context, opts ...Option) (*Client, error) {
	c := Client{
		client: http.DefaultClient,
		header: make(http.Header),
	}
	c.header.Add("Authority", "www.threads.net")
	c.header.Add("Accept", "*/*")
	c.header.Add("Accept-Language", "en-US,en;q=0.9")
	c.header.Add("Cache-Control", "no-cache")
	c.header.Add("Content-Type", "application/x-www-form-urlencoded")
	c.header.Add("Connetion", "keep-alive")
	c.header.Add("Origin", "https://www.threads.net")
	c.header.Add("Pragma", "no-cache")
	c.header.Add("Sec-Fetch-Site", "same-origin")
	c.header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.44.639.844 Safari/537.36")
	c.header.Add("X-IG-ASBD-ID", "129477")
	c.header.Add("X-IG-App-ID", "238260118697367")

	for _, opt := range opts {
		opt(&c)
	}

	if c.token == "" {
		token, err := c.getToken(ctx)
		if err != nil {
			return nil, err
		}
		c.token = token
	}
	c.header.Add("X-FB-LSD", c.token)

	return &c, nil
}

func WithToken(token string) Option {
	return func(c *Client) {
		c.token = token
	}
}

func WithClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}

func WithHeader(header http.Header) Option {
	return func(c *Client) {
		c.header = header
	}
}

// getToken returns the token used to make requests to the API.
func (c *Client) getToken(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"instagram", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.44.639.844 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	pos := bytes.Index(body, []byte("\"token\""))
	return string(body[pos+9 : pos+31]), nil
}

// GetUserID returns the user ID of the given username.
func (c *Client) GetUserID(ctx context.Context, name string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+name, nil)
	if err != nil {
		return 0, err
	}
	req.Header = c.header.Clone()
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Referer", baseURL)
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Upgrade-Insecure-Requests", "1")

	req.Header.Del("X-ASBD-ID")
	req.Header.Del("X-FB-LSD")
	req.Header.Del("X-IG-App-ID")

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var userID int
	if err := binary.Read(bytes.NewReader(userIDRegex.FindSubmatch(body)[1]), binary.BigEndian, &userID); err != nil {
		return 0, err
	}
	return userID, nil
}

// GetUser returns the profile of the given user ID.
func (c *Client) GetUser(ctx context.Context, userID int) ([]byte, error) {
	h := c.header.Clone()
	h.Add("X-FB-Friendly-Name", "BarcelonaProfileRootQuery")
	return sendRequest(ctx, c.client, h, c.token,
		getUserDocID,
		map[string]int{"userID": userID},
	)
}

// GetUserThreads returns the threads posted by the given user ID.
func (c *Client) GetUserThreads(ctx context.Context, userID int) ([]byte, error) {
	h := c.header.Clone()
	h.Add("X-FB-Friendly-Name", "BarcelonaProfileThreadsTabQuery")
	return sendRequest(ctx, c.client, h, c.token,
		getUserThreadsDocID,
		map[string]int{"userID": userID},
	)
}

// GetUserReplies returns the replies posted by the given user ID.
func (c *Client) GetUserReplies(ctx context.Context, userID int) ([]byte, error) {
	h := c.header.Clone()
	h.Add("X-FB-Friendly-Name", "BarcelonaProfileRepliesTabQuery")
	return sendRequest(ctx, c.client, h, c.token,
		getUserRepliesDocID,
		map[string]int{"userID": userID},
	)
}

// GetPost returns the post of the given post ID.
func (c *Client) GetPost(ctx context.Context, postID int) ([]byte, error) {
	h := c.header.Clone()
	h.Add("X-FB-Friendly-Name", "BarcelonaPostPageQuery")
	return sendRequest(ctx, c.client, h, c.token,
		getPostDocID,
		map[string]int{"postID": postID},
	)
}

// GetLikers returns the liker list of the given post ID.
func (c *Client) GetLikers(ctx context.Context, postID int) ([]byte, error) {
	return sendRequest(ctx, c.client, c.header, c.token,
		getLikersDocID,
		map[string]int{"mediaID": postID},
	)
}

func sendRequest(ctx context.Context, c *http.Client, headers http.Header, token, docID string, variables map[string]int) ([]byte, error) {
	b, err := json.Marshal(variables)
	if err != nil {
		return nil, err
	}

	data := url.Values{}
	data.Set("lsd", token)
	data.Set("doc_id", docID)
	data.Set("variables", string(b))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header = headers

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
