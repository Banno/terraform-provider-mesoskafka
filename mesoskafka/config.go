package mesoskafka

type Config struct {
	Url    string
	client *Client
}

func (c *Config) loadAndValidate() error {
	c.client = NewClientForUrl(c.Url)
	return nil
}
