package proxy

type CreateProxyRequest struct {
	Domain     string `json:"domain"`
	TargetHost string `json:"target_host"`
	TargetPort int    `json:"target_port"`
	SSL        bool   `json:"ssl"`
}

type Proxy struct {
	ID         string `json:"id"`
	Domain     string `json:"domain"`
	TargetHost string `json:"target_host"`
	TargetPort int    `json:"target_port"`
	SSL        bool   `json:"ssl"`
	ConfigFile string `json:"config_file"`
}