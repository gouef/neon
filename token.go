package neon

type TokenType int

const (
	String     TokenType = 1
	Literal    TokenType = 2
	Char       TokenType = 0
	Comment    TokenType = 3
	Newline    TokenType = 4
	Whitespace TokenType = 5
	End        TokenType = -1
)

type Token struct {
	Type     any
	Text     string
	Position Position
}

func NewToken(tokenType any, text string, position Position) Token {
	return Token{
		Type:     tokenType,
		Text:     text,
		Position: position,
	}
}

func (t Token) Is(kinds ...any) bool {
	for _, kind := range kinds {
		if t.Type == kind {
			return true
		}
	}
	return false
}
