package parser

import (
    "fmt"
    "os"
    "strconv"
    "cpl/main/source/lexer"
)

// Parser structure to hold the tokens and the current position
type Parser struct {
    Tokens []lexer.TokenInfo
    Pos    int
}

// NewParser creates a new parser
func NewParser(Tokens []lexer.TokenInfo) *Parser {
    return &Parser{Tokens: Tokens, Pos: 0}
}

// CurrentToken returns the current token
func (p *Parser) CurrentToken() lexer.TokenInfo {
    if p.Pos < len(p.Tokens) {
        return p.Tokens[p.Pos]
    }
    return lexer.TokenInfo{Token: lexer.INVALID, Lexeme: ""}
}

// Advance moves to the next token
func (p *Parser) Advance() {
    p.Pos++
}

// Match checks if the current token matches the expected token and advances
func (p *Parser) Match(expected lexer.Token) bool {
    if p.CurrentToken().Token == expected {
        p.Advance()
        return true
    }
    return false
}

// Parse parses the input tokens
func (p *Parser) Parse() {
    p.STMT_LIST()
}

// STMT_LIST --> STMT. | STMT; STMT_LIST
func (p *Parser) STMT_LIST() {
    p.STMT()
    for p.Match(lexer.SEMICOLON) {
        p.STMT()
    }
    if !p.Match(lexer.PERIOD) {
        p.Error("Expected '.' at the end of the statement list")
    }
}

// STMT --> POINT_DEF | TEST
func (p *Parser) STMT() {
    if p.CurrentToken().Token == lexer.ID {
        p.POINT_DEF()
    } else if p.CurrentToken().Token == lexer.TEST {
        p.TEST()
    } else {
        p.Error("Expected a statement")
    }
}

// POINT_DEF --> ID = point(NUM, NUM)
func (p *Parser) POINT_DEF() {
    if !p.Match(lexer.ID) {
        p.Error("Expected an identifier")
    }
    if !p.Match(lexer.ASSIGN) {
        p.Error("Expected '='")
    }
    if !p.Match(lexer.POINT) {
        p.Error("Expected 'point'")
    }
    if !p.Match(lexer.LEFT_PAREN) {
        p.Error("Expected '('")
    }
    if !p.Match(lexer.NUM) {
        p.Error("Expected a number")
    }
    if !p.Match(lexer.COMMA) {
        p.Error("Expected ','")
    }
    if !p.Match(lexer.NUM) {
        p.Error("Expected a number")
    }
    if !p.Match(lexer.RIGHT_PAREN) {
        p.Error("Expected ')'")
    }
}

// TEST --> test(OPTION, POINT_LIST)
func (p *Parser) TEST() {
    if !p.Match(lexer.TEST) {
        p.Error("Expected 'test'")
    }
    if !p.Match(lexer.LEFT_PAREN) {
        p.Error("Expected '('")
    }
    if p.CurrentToken().Token != lexer.TRIANGLE && p.CurrentToken().Token != lexer.SQUARE {
        p.Error("Expected 'triangle' or 'square'")
    }
    p.Advance()
    if !p.Match(lexer.COMMA) {
        p.Error("Expected ','")
    }
    p.POINT_LIST()
    if !p.Match(lexer.RIGHT_PAREN) {
        p.Error("Expected ')'")
    }
}

// POINT_LIST --> ID | ID, POINT_LIST
func (p *Parser) POINT_LIST() {
    if !p.Match(lexer.ID) {
        p.Error("Expected an identifier")
    }
    for p.Match(lexer.COMMA) {
        if !p.Match(lexer.ID) {
            p.Error("Expected an identifier")
        }
    }
}

// Error prints an error message and exits
func (p *Parser) Error(message string) {
    fmt.Printf("Syntax error: %s\n", message)
    os.Exit(1)
}

func (p *Parser) ConsumeNumber() int {
    if p.Pos < len(p.Tokens) && p.Tokens[p.Pos].Token == lexer.NUM {
        num, _ := strconv.Atoi(p.Tokens[p.Pos].Lexeme)
        p.Pos++
        return num
    }
    return 0
}
