package commands

type Command interface {
	Serialize() string
	ToHuman() string
}
