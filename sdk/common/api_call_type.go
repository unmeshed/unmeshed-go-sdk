package common

type ApiCallType string

const (
	ApiCallTypeSync   ApiCallType = "SYNC"
	ApiCallTypeAsync  ApiCallType = "ASYNC"
	ApiCallTypeStream ApiCallType = "STREAM"
)

func (a ApiCallType) String() string {
	return string(a)
}

func (a ApiCallType) IsValid() bool {
	switch a {
	case ApiCallTypeSync, ApiCallTypeAsync, ApiCallTypeStream:
		return true
	}
	return false
}
