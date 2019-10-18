# SpyServer for Go

Create simple, local test servers that records requests and responds with pre-configured results.

## Why use this

Test code making http calls without having a slow, unreliable network in the way.
Simplify testing when making calls to external services.

Go's httptest package is good, but this library gives you a few more benefits:

1. You don't need to point your tests to a different URL.
2. You can test that your client code closes the response body.
3. You can simulate various errors that may happen with the transport.

## How to use

First, create a spyserver like this:
```
response := &http.Response{
    StatusCode: http.StatusOK,
	Body:       body,
}
spy := spyserver.New(response, nil)
```

Next, set the Transport in your http client to the spyserver and call it:
```
client := &http.Client{
    Transport: spy,
}
client.Do(request)
```

You can inspect the request sent to the server in your tests:
```
spy.GetRequest()
```

You can even test that you close the response body in your production code.
```
body := spyserver.NewCloseDetectorFromString(jsonString)

// Use the body the response you configure with your spy server
// Then just inject the spy server as usual and call your code

body.IsClosed() // returns true if you closed the response body
```