package clients

type BmpClient interface {
	Info() (InfoResponse, error)
}
