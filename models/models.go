package models


type StartResponse struct {
	APIVersion string `json:"apiversion,omitempty"`
	Author string `json:"author,omitempty"`
	Color string `json:"color,omitempty"`
	Head string `json:"head,omitempty"`
	Tail string `json:"tail,omitempty"`
}


type MoveResponse struct {
	Move  string `json:"move"`
	Taunt string `json:"taunt"`
}
