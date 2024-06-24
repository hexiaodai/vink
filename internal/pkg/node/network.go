package node

type NetworkInterface struct {
	Name    string `json:"name"`
	IP      string `json:"ip"`
	State   string `json:"state"`
	Subnet  string `json:"subnet"`
	Gateway string `json:"gateway"`
}
