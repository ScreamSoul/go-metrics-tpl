package utils

type Closer interface {
	Close() error
}

func CloseForse(closer Closer) {
	if err := closer.Close(); err != nil {
		return
	}
}
