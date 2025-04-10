package adopter

type ConnAdopter interface {
	ReadMessageFunc() (int, []byte, error)
	WriteMessageFunc(int, []byte) error
	CloseFunc() error
}
