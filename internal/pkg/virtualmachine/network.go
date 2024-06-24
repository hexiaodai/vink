package virtualmachine

type SubnetConfiguration struct {
	Interface string   `json:"interface"`
	IPv4      []string `json:"ipv4"`
}
