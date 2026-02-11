package internal

type RandomGenerator interface {
	GenerateRandomAlias() string
}

type HashGenerator interface {
	GenerateHashedAlias(str string) string
}


