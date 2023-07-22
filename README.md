# github.com/satorunooshie/threads
[![Go Doc](https://godoc.org/github.com/satorunooshie/threads?status.svg)](https://godoc.org/github.com/satorunooshie/threads)

Unofficial, Reverse-Engineered Go API client for Meta's [Threads](threads.net).

## Initialize
Without options, it will fetch access token automatically via http.DefaultClient.
```go
client, err := threads.NewClient(ctx)
```

With access token, it will use it to make requests to the API.
```go
client, err := threads.NewClient(ctx, WithToken(`OqhxMWDlJViVPLZiN5p9Un`))
```

With http client, it will use it to make requests to the API.
```go
client, err := threads.NewClient(ctx, WithClient(&http.Client{}))
```

With custom header, the original header will be modified to make requests to the API.
```go
client, err := threads.NewClient(ctx, WithHeader(http.Header{}))
```

## APIs
### GetUserID
GetUserID returns userID from account name.

```go
id, err := client.GetUserID(ctx, "zuck")
```

### GetUser
GetUser returns user information of user id.

```go
b, err := client.GetUser(ctx, 314216)
```

response example: [user.json](https://github.com/satorunooshie/threads/blob/main/testdata/user.json)

### GetUserThreads
GetUserThreads returns user threads of user id.

```go
b, err := client.GetUserThreads(ctx, 314216)
```

response example: [threads.json](https://github.com/satorunooshie/threads/blob/main/testdata/threads.json)

### GetUserReplies
GetUserReplies returns user replies of post id.

```go
b, err := client.GetUserReplies(ctx, 3141002295235099165)
```

response example: [replies.json](https://github.com/satorunooshie/threads/blob/main/testdata/replies.json)

### GetPost
GetPost returns single post of post id.

```go
b, err := client.GetPost(ctx, 3141002295235099165)
```

response example: [post.json](https://github.com/satorunooshie/threads/blob/main/testdata/post.json)

### GetLikers
GetLikers returns list of users who likes the post of post id.

```go
b, err := client.GetLikers(ctx, 3141002295235099165)
```

response example: [likers.json](https://github.com/satorunooshie/threads/blob/main/testdata/likers.json)

## Special Thanks
[junhoyeo/threads-api](https://github.com/junhoyeo/threads-api)

[m1guelpf/threads-re](https://github.com/m1guelpf/threads-re)

[antonprokopovich/go-threads](https://github.com/antonprokopovich/go-threads)
