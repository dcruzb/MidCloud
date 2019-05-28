package dist

type Message struct {
	Header Header
	Body   Body
}

type Header struct {
	Magic   string
	Version int
	Order   bool
	Type    int
	Size    int
}

type Body struct {
	RequestHeader RequestHeader
	RequestBody   RequestBody
	ReplyHeader   ReplyHeader
	ReplyBody     interface{}
}

type RequestHeader struct {
	Context          string
	RequestId        int
	ResponseExpected bool
	ObjectKey        int
	Operation        string
}

type ReplyHeader struct {
	ServiceContext string
	RequestId      int
	ReplyStatus    int
}

type RequestBody struct {
	Parameters []interface{}
}
