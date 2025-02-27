package request

type ReqFindItImageInfo struct {
	ImageInfoList []ImageInfo `json:"imageInfoList"`
}
type ImageInfo struct {
	NormalImage   string     `json:"normalImage"`
	AbnormalImage string     `json:"abnormalImage"`
	Positions     []Position `json:"positions"`
	Level         int        `json:"level"`
}
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
