package rapid7

type Client struct {
	IDR *IDR
	VM  *VM
}

func New(region, apiKey string) (client *Client, err error) {
	idr, err := newIDR(region, apiKey)
	if err != nil {
		return
	}
	vm, err := newVM(region, apiKey)
	client = &Client{
		IDR: idr,
		VM:  vm,
	}
	return
}
