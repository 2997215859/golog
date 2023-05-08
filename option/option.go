package option

const (
	optkeyAppName      = "app-name"
	optkeyMaxAge       = "max-age"
	optkeyRotationTime = "rotation-time"
)

type Option interface {
	Name() string
	Value() interface{}
}
