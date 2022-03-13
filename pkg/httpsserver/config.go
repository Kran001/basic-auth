package httpsserver

type Config interface {
	Host() string
	PortHTTPS() string
	HTTPSConnectionString() string
}
