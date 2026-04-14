package blobs

type readingStatus int8

const (
	readingHeader readingStatus = iota
	readingBody
)

type HeaderInformation int8

const (
	HeaderTitle HeaderInformation = iota
	HeaderDescription
	HeaderDatetime
	HeaderURL
)
