package neon

type Encoder struct {
	BLOCK       bool
	blockMode   bool
	indentation string
}

func NewEncoder() Encoder {
	return Encoder{
		BLOCK:       true,
		blockMode:   false,
		indentation: "\t",
	}
}

func (e Encoder) encode(val) string {

}
