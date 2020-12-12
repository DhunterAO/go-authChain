package pbftTypes

const (
	RequestInfo    uint8 = 0
	ResponseInfo   uint8 = 1
	CommitInfo     uint8 = 2
	ViewChangeInfo uint8 = 3

	RequestUrl    string = "/pbft/request"
	ResponseUrl   string = "/pbft/response"
	CommitUrl     string = "/pbft/commit"
	ViewChangeUrl string = "/pbft/viewChange"
)
