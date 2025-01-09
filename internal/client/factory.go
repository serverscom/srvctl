package client

type ClientFactory interface {
	NewClient(token, endpoint string) *Client
}

type DefaultClientFactory struct{}

func (f *DefaultClientFactory) NewClient(token, endpoint string) *Client {
	return NewClient(token, endpoint)
}

type TestClientFactory struct {
	TestClient *Client
}

func (f *TestClientFactory) NewClient(token, endpoint string) *Client {
	return f.TestClient
}
