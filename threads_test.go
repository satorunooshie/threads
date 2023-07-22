package threads

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//go:embed testdata/token.html
var tokenResponse []byte

//go:embed testdata/user_id.html
var userIDResponse []byte

//go:embed testdata/user.json
var userResponse []byte

//go:embed testdata/threads.json
var threadsResponse []byte

//go:embed testdata/replies.json
var repliesResponse []byte

//go:embed testdata/post.json
var postResponse []byte

//go:embed testdata/likers.json
var likersResponse []byte

type fake struct{}

var _ http.RoundTripper = (*fake)(nil)

func (*fake) RoundTrip(req *http.Request) (*http.Response, error) {
	resp := httptest.NewRecorder()
	switch req.URL.String() {
	case baseURL + "/@instagram":
		if req.Header.Get("Referer") == baseURL {
			_, _ = resp.Write(userIDResponse)
		} else {
			_, _ = resp.Write(tokenResponse)
		}
	case apiURL:
		switch req.FormValue("doc_id") {
		case getUserDocID:
			_, _ = resp.Write(userResponse)
		case getUserThreadsDocID:
			_, _ = resp.Write(threadsResponse)
		case getUserRepliesDocID:
			_, _ = resp.Write(repliesResponse)
		case getPostDocID:
			_, _ = resp.Write(postResponse)
		case getLikersDocID:
			_, _ = resp.Write(likersResponse)
		default:
			goto NotFound
		}
	default:
		goto NotFound
	}
NotFound:
	resp.WriteHeader(http.StatusNotFound)
	return resp.Result(), nil
}

func TestNewClient(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		want    *Client
		args    args
		wantErr bool
	}{
		{
			name: "with token option",
			want: &Client{
				client: http.DefaultClient,
				header: http.Header{
					"Authority":       []string{"www.threads.net"},
					"Accept":          []string{"*/*"},
					"Accept-Language": []string{"en-US,en;q=0.9"},
					"Cache-Control":   []string{"no-cache"},
					"Content-Type":    []string{"application/x-www-form-urlencoded"},
					"Connetion":       []string{"keep-alive"},
					"Origin":          []string{"https://www.threads.net"},
					"Pragma":          []string{"no-cache"},
					"Sec-Fetch-Site":  []string{"same-origin"},
					"User-Agent":      []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.44.639.844 Safari/537.36"},
					"X-Ig-Asbd-Id":    []string{"129477"},
					"X-Ig-App-Id":     []string{"238260118697367"},
					"X-Fb-Lsd":        []string{"OqhxMWDlJViVPLZiN5p9Un"},
				},
				token: "OqhxMWDlJViVPLZiN5p9Un",
			},
			args: args{
				opts: []Option{
					WithToken("OqhxMWDlJViVPLZiN5p9Un"),
				},
			},
		},
		{
			name: "with client option",
			want: &Client{
				client: &http.Client{Transport: &fake{}},
				header: http.Header{
					"Authority":       []string{"www.threads.net"},
					"Accept":          []string{"*/*"},
					"Accept-Language": []string{"en-US,en;q=0.9"},
					"Cache-Control":   []string{"no-cache"},
					"Content-Type":    []string{"application/x-www-form-urlencoded"},
					"Connetion":       []string{"keep-alive"},
					"Origin":          []string{"https://www.threads.net"},
					"Pragma":          []string{"no-cache"},
					"Sec-Fetch-Site":  []string{"same-origin"},
					"User-Agent":      []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.44.639.844 Safari/537.36"},
					"X-Ig-Asbd-Id":    []string{"129477"},
					"X-Ig-App-Id":     []string{"238260118697367"},
					"X-Fb-Lsd":        []string{"OqhxMWDlJViVPLZiN5p9Un"},
				},
				token: "OqhxMWDlJViVPLZiN5p9Un",
			},
			args: args{
				opts: []Option{
					WithClient(&http.Client{Transport: &fake{}}),
				},
			},
		},
		{
			name: "with header option",
			want: &Client{
				client: http.DefaultClient,
				header: http.Header{
					"X-Fb-Lsd": []string{"OqhxMWDlJViVPLZiN5p9Un"},
				},
				token: "OqhxMWDlJViVPLZiN5p9Un",
			},
			args: args{
				opts: []Option{
					WithHeader(http.Header{}),
					WithToken("OqhxMWDlJViVPLZiN5p9Un"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(context.Background(), tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(Client{})); diff != "" {
				t.Errorf("NewClient() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestClient_GetUserID(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				name: "instagram",
			},
			want: 25025320,
		},
	}
	ctx := context.Background()
	client, err := NewClient(ctx, WithClient(&http.Client{Transport: &fake{}}), WithToken("OqhxMWDlJViVPLZiN5p9Un"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetUserID(ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserID() got = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestClient_GetUser(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				id: 25025320,
			},
			want: userResponse,
		},
	}
	ctx := context.Background()
	client, err := NewClient(ctx, WithClient(&http.Client{Transport: &fake{}}), WithToken("OqhxMWDlJViVPLZiN5p9Un"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetUser(ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetUser() mismatch (-want +got):\n%s", diff)
			}
			var pretty bytes.Buffer
			_ = json.Indent(&pretty, got, "", "\t")
			t.Log(pretty.String())
		})
	}
}

func TestClient_GetUserThreads(t *testing.T) {
	type args struct {
		userID int
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				userID: 25025320,
			},
			want: threadsResponse,
		},
	}
	ctx := context.Background()
	client, err := NewClient(ctx, WithClient(&http.Client{Transport: &fake{}}), WithToken("OqhxMWDlJViVPLZiN5p9Un"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetUserThreads(ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserThreads() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetUserThreads() mismatch (-want +got):\n%s", diff)
			}
			var pretty bytes.Buffer
			_ = json.Indent(&pretty, got, "", "\t")
			t.Log(pretty.String())
		})
	}
}

func TestClient_GetUserReplies(t *testing.T) {
	type args struct {
		userID int
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				userID: 25025320,
			},
			want: repliesResponse,
		},
	}
	ctx := context.Background()
	client, err := NewClient(ctx, WithClient(&http.Client{Transport: &fake{}}), WithToken("OqhxMWDlJViVPLZiN5p9Un"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetUserReplies(ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserReplies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetUserReplies() mismatch (-want +got):\n%s", diff)
			}
			var pretty bytes.Buffer
			_ = json.Indent(&pretty, got, "", "\t")
			t.Log(pretty.String())
		})
	}
}

func TestClient_GetPost(t *testing.T) {
	type args struct {
		postID int
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				postID: 3152079361880880077,
			},
			want: postResponse,
		},
	}
	ctx := context.Background()
	client, err := NewClient(ctx, WithClient(&http.Client{Transport: &fake{}}), WithToken("OqhxMWDlJViVPLZiN5p9Un"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetPost(ctx, tt.args.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetPost() mismatch (-want +got):\n%s", diff)
			}
			var pretty bytes.Buffer
			_ = json.Indent(&pretty, got, "", "\t")
			t.Log(pretty.String())
		})
	}
}

func TestClient_GetLikers(t *testing.T) {
	type args struct {
		postID int
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				postID: 3152079361880880077,
			},
			want: likersResponse,
		},
	}
	ctx := context.Background()
	client, err := NewClient(ctx, WithClient(&http.Client{Transport: &fake{}}), WithToken("OqhxMWDlJViVPLZiN5p9Un"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetLikers(ctx, tt.args.postID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLikers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetLikers() mismatch (-want +got):\n%s", diff)
			}
			var pretty bytes.Buffer
			_ = json.Indent(&pretty, got, "", "\t")
			t.Log(pretty.String())
		})
	}
}
