package command

const (
	Set = "set"
	Get = "get"
)

type Command interface {
}

type SetCommand struct {
	Key, Val []byte
}

type GetCommand struct {
	Key []byte
}
