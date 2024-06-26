package vulinboxagentclient

func WithOnClose(f func()) Option {
	return func(c *Client) {
		c.onClose = f
	}
}
