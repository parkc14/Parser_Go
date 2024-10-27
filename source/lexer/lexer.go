package lexer

import (
	"regexp"
	//"fmt"
)

type Token int

//Token types
const (
	ID Token = iota       // starts at 0
	POINT                 
	NUM                   
	ASSIGN               
	SEMICOLON             
	COMMA                
	PERIOD              
	LEFT_PAREN           
	RIGHT_PAREN          
	TRIANGLE            
	SQUARE                
	TEST                  
	INVALID              
)

//TokenNames maps Token types to their string representations.
var tokenNames = map[Token]string{
	ID:         "id",
	POINT:      "point",
	NUM:        "num",
	ASSIGN:     "assign",
	SEMICOLON:  "semicolon",
	COMMA:      "comma",
	PERIOD:     "period",
	LEFT_PAREN: "left_paren",
	RIGHT_PAREN:"right_paren",
	TRIANGLE:   "triangle",
	SQUARE:     "square",
	TEST:       "test",
	INVALID:    "invalid",
}

//TokenInfo holds a token and its corresponding lexeme.
type TokenInfo struct {
	Token  Token
	Lexeme string
}

//String returns the string representation of the Token.
func (TokenKind Token) String() string {
	if name, ok := tokenNames[TokenKind]; ok {
		return name
	}
	return "unknown"
}

//Lexer tokenizes the input string and returns a slice of TokenInfo.
//it uses regular expressions to match predefined token patterns and skips whitespace characters.
func Lexer(input string) []TokenInfo {
    var tokens []TokenInfo
    tokenPatterns := []struct {
        token   Token
        pattern *regexp.Regexp
    }{
        {POINT, regexp.MustCompile(`^point`)},
        {TRIANGLE, regexp.MustCompile(`^triangle`)},
        {SQUARE, regexp.MustCompile(`^square`)},
        {TEST, regexp.MustCompile(`^test`)},
        {NUM, regexp.MustCompile(`^[0-9]+`)},
        {ASSIGN, regexp.MustCompile(`^=`)},
        {SEMICOLON, regexp.MustCompile(`^;`)},
        {COMMA, regexp.MustCompile(`^,`)},
        {PERIOD, regexp.MustCompile(`^\.`)},
        {LEFT_PAREN, regexp.MustCompile(`^\(`)},
        {RIGHT_PAREN, regexp.MustCompile(`^\)`)},
        {ID, regexp.MustCompile(`^[a-z][a-z_]*`)},
    }

	//Initialize the position index
	i := 0
	//loop through the input string until the end 
	for i < len(input) {
		// Skip whitespace characters
		if input[i] == ' ' || input[i] == '\t' || input[i] == '\n' || input[i] == '\r' {
			i++
			continue
		}

		matched := false
		//check each token pattern
		for _, tp := range tokenPatterns {
            if loc := tp.pattern.FindStringIndex(input[i:]); loc != nil && loc[0] == 0 {
                //if a pattern matches, extract the lexeme
                lexeme := input[i : i+loc[1]]
                // Append the token and lexeme to the tokens slice
                tokens = append(tokens, TokenInfo{Token: tp.token, Lexeme: lexeme})
                // Move the index forward by the length of the matched lexeme
                i += loc[1]
                matched = true
                break
			}
		}

		// If no pattern matches, mark the token as INVALID
		if !matched {
			tokens = append(tokens, TokenInfo{Token: INVALID, Lexeme: string(input[i])})
			i++
		}
	}
	//fmt.Println(tokens)
	return tokens
}
