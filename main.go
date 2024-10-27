package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
	"cpl/main/source/lexer"
    "cpl/main/source/parser"
)

// ReadFile reads the content of the given file and returns it as a string.
func ReadFile(filename string) string {
	content, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatal(err)
    }
    return string(content)
}

func main() {
	//checks to see if there are less than 3 arguments in the slice
    if len(os.Args) < 3 {
        fmt.Println("Missing parameter, usage:\ngo run . filename -flag\nflag can be p for prolog generation\nflag can be s for scheme generation")
        return
    }
	//argument in position 0 is the period 
    filename := os.Args[1]  //first argument is filename
    flag := os.Args[2]		//second argument is the flag

    //read the file content
    fileContent := ReadFile(filename)

	
    //tokenize the input
    tokens := lexer.Lexer(fileContent)

    //parse the tokens
    parser := parser.NewParser(tokens)
    parser.Parse()

    //chose format based on the flag
    switch flag {
    case "-s":
        fmt.Println("; Processing Input File", filename)
        fmt.Println("; Lexical and Syntax analysis passed")
        fmt.Println("; Generating Scheme Code")
        GenerateSchemeCode(parser)
    case "-p":
        fmt.Println("/* Processing Input File", filename)
        fmt.Println("   Lexical and Syntax analysis passed")
        fmt.Println("   Generating Prolog Code */")
        GeneratePrologCode(parser)
    default:
        fmt.Println("Invalid flag, usage:\ngo run . filename -flag\nflag can be p for prolog generation\nflag can be s for scheme generation")
    }
}

func GenerateSchemeCode(p *parser.Parser) {
    p.Pos = 0
    points := make(map[string]string)
    for p.Pos < len(p.Tokens) {
        token := p.Tokens[p.Pos]
        switch token.Token {
        case lexer.ID:
            id := token.Lexeme
            p.Pos++
            if p.Match(lexer.ASSIGN) && p.Match(lexer.POINT) && p.Match(lexer.LEFT_PAREN) {
                x := p.ConsumeNumber()
                if !p.Match(lexer.COMMA){
                    p.Pos++
                    continue
                }
                y := p.ConsumeNumber()
                if !p.Match(lexer.RIGHT_PAREN){
                    p.Pos++
                    continue
                }
                points[id] = fmt.Sprintf("(make-point %d %d)", x, y)
            } else{
                p.Pos++
            }
            
        case lexer.TEST:
            p.Pos++ 
            if p.Match(lexer.LEFT_PAREN) {
                option := p.Tokens[p.Pos].Lexeme
                fmt.Printf("(process-%s", option)
                for p.Tokens[p.Pos].Token != lexer.RIGHT_PAREN {
                    id := p.Tokens[p.Pos].Lexeme
                    if !p.Match(lexer.ID){
                        p.Pos++
                        continue
                    }
                    fmt.Printf(" %s", points[id])
                    if p.Tokens[p.Pos].Token == lexer.COMMA {
                        p.Match(lexer.COMMA)
                    }
                }
                if !p.Match(lexer.RIGHT_PAREN) {
                    p.Pos++
                    continue
                }
                fmt.Println(")")
            } else{
                p.Pos++
            }
        default:
            p.Pos++
        }
    }
}


// GeneratePrologCode generates Prolog code based on the parsed tokens
func GeneratePrologCode(p *parser.Parser) {
    p.Pos = 0
    points := make(map[string]string)
    for p.Pos < len(p.Tokens) {
        token := p.Tokens[p.Pos]
        switch token.Token {
        case lexer.ID:
            id := token.Lexeme
            p.Pos++
            if p.Match(lexer.ASSIGN) && p.Match(lexer.POINT) && p.Match(lexer.LEFT_PAREN) {
                x := p.ConsumeNumber()
                if !p.Match(lexer.COMMA) {
                    p.Pos++
                    continue
                }
                y := p.ConsumeNumber()
                if !p.Match(lexer.RIGHT_PAREN) {
                    p.Pos++
                    continue
                }
                points[id] = fmt.Sprintf("point2d(%d, %d)", x, y)
            } else {
                p.Pos++
            }
        case lexer.TEST:
            p.Pos++
            if p.Match(lexer.LEFT_PAREN) {
                option := p.Tokens[p.Pos].Lexeme
                fmt.Printf("\n/* Processing test(%s ", option)
                var pointIDs []string
                for p.Tokens[p.Pos].Token != lexer.RIGHT_PAREN {
                    id := p.Tokens[p.Pos].Lexeme
                    if !p.Match(lexer.ID) {
                        p.Pos++
                        continue
                    }
                    pointIDs = append(pointIDs, id)
                    if p.Tokens[p.Pos].Token == lexer.COMMA {
                        p.Match(lexer.COMMA)
                    }
                }
                if !p.Match(lexer.RIGHT_PAREN) {
                    p.Pos++
                    continue
                }
                fmt.Printf(", %s", pointIDs[0])
                for _, id := range pointIDs[1:] {
                    fmt.Printf(", %s", id)
                }
                fmt.Printf(") */\n")
                fmt.Printf("query(%s(", option)
                for i, id := range pointIDs {
                    if i > 0 {
                        fmt.Print(", ")
                    }
                    fmt.Printf("%s", points[id])
                }
                fmt.Println(")).")
                if option == "triangle" {
                    queries := []string{
                        "line", "triangle", "vertical", "horizontal", "equilateral",
                        "isosceles", "right", "scalene", "acute", "obtuse",
                    }
                    for _, query := range queries {
                        fmt.Printf("query(%s(", query)
                        for i, id := range pointIDs {
                            if i > 0 {
                                fmt.Print(", ")
                            }
                            fmt.Printf("%s", points[id])
                        }
                        fmt.Println(")).")
                    }
                }
            } else {
                p.Pos++
            }
        default:
            p.Pos++
        }
    }
    fmt.Println("\n/* Query Processing */")
    fmt.Println("writeln(T) :- write(T), nl.")
    fmt.Println("main:- forall(query(Q), Q-> (writeln('yes')) ; (writeln('no'))),")
    fmt.Println("       halt.")
}
