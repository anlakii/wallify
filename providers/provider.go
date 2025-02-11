package providers

type Provider interface {
	Update() (bool, error)
}
