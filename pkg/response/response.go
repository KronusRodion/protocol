package response



type Response uint8


const (
	OK Response = 0
	Error Response = 1
	NotFound Response = 3
)