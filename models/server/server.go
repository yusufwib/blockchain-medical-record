package server

type Server struct {
	AppName      string
	HTTPPort     int
	JWTSecretKey string
}
