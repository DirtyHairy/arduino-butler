package controls

type MarshalledBackend struct {
    Id string `json:"id"`
    Type string `json:"type"`
    Host string `json:"host"`
}

type MarshalledSwitch struct {
    Id string `json:"id"`
    Type string `json:"type"`
    BackendId string `json:"backendId"`
    BackendIndex uint `json:"backendIndex"`
}

type MarshalledControlSet struct {
    Backends []MarshalledBackend `json:"backends"`
    Switches []MarshalledSwitch `json:"switches"`
}
