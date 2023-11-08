package lexer

import (
	"fmt"
	"gada/lexer"
	"gada/reader"
	"gada/token"
	"log"
	"os"
	"strings"
)

type testlexer struct {
	tokens  []lexer.Token
	lexiDic []string
}

func compareTokens(token1 lexer.Token, token2 lexer.Token, lexiDic1 []string, lexiDic2 []string) bool {
	return token1.Position == token2.Position &&
		(token1.Position == 0 || lexiDic1[token1.Position-1] == lexiDic2[token2.Position-1]) &&
		// previous condition check if the literals are equal in case where the token is a literal
		token1.Value == token2.Value &&
		token1.Beginning == token2.Beginning &&
		token1.End == token2.End
}

func AllTest() {
	files, err := os.ReadDir("examples")
	if err != nil {
		log.Fatalf("the directory provided have this error : %s", err)
	}
	expected := getExpected()
	for _, file := range files {
		nameNoExt := strings.Split(file.Name(), ".")[0]
		testPassed := true
		fileLexer := reader.FileLexer("examples/" + file.Name())
		foundTokens, lexicon := fileLexer.Read()
		for ind, token := range foundTokens {
			expecTokens, expecLexi := expected[nameNoExt].tokens, expected[nameNoExt].lexiDic
			if ind >= len(expecTokens) || !compareTokens(token, expecTokens[ind], lexicon, expecLexi) {
				testPassed = false
				if ind >= len(expecTokens) {
					log.Fatalf("\nTest: %s There is more token than expected", file.Name())
				}
				tokenLit1, tokenLit2 := "", ""
				if token.Position != 0 {
					tokenLit1 = lexicon[token.Position-1]
				}
				if expecTokens[ind].Position != 0 {
					tokenLit2 = expecLexi[expecTokens[ind].Position-1]
				}
				// tokenLit1 and tokenLit2 are the literals in case tokens are literals
				// there here for the debug only
				log.Fatalf("\ntoken number: %d token gen: %v %s is different than token expected: %v %s", ind, token, tokenLit1, expecTokens[ind], tokenLit2)
			} else {
				//fmt.Printf("token number: %d token: %v lexi: %v\n", ind, expecTokens, expecLexi)
			}
		}
		if testPassed {
			fmt.Printf("Test %s: passed succesfully \n\n", file.Name())
		} else {
			fmt.Printf("Test %s: not passed\n", file.Name())
		}

	}
}

func DisplayLexer(name string) {
	lexer := reader.FileLexer("examples/" + name)
	foundTokens, lexicon := lexer.Read()
	line := -1
	for _, tok := range foundTokens {
		if tok.Beginning.Line != line {
			line = tok.Beginning.Line
			if tok.Position != 0 {
				fmt.Printf("\nLine : %d (%s:%s:%s from: %d to :%d )", tok.Beginning.Line, tok.Type, token.Tokens[tok.Value], lexicon[tok.Position-1], tok.Beginning.Column, tok.End.Column)
			} else {
				fmt.Printf("\nLine : %d (%s:%s from: %d to :%d )", tok.Beginning.Line, tok.Type, token.Tokens[tok.Value], tok.Beginning.Column, tok.End.Column)
			}
		} else {
			if tok.Position != 0 {
				fmt.Printf("(%s:%s:%s from: %d to :%d )", tok.Type, token.Tokens[tok.Value], lexicon[tok.Position-1], tok.Beginning.Column, tok.End.Column)
			} else {
				fmt.Printf("(%s:%s from: %d to :%d )", tok.Type, token.Tokens[tok.Value], tok.Beginning.Column, tok.End.Column)
			}
		}
	}
	for _, lex := range lexicon {
		fmt.Println(lex)
	}
}
