package lexer

import (
	"fmt"
	"gada/lexer"
	"gada/reader"
	"log"
	"os"
)

type testlexer struct {
	tokens  []lexer.Token
	lexidic []any
}

func compareTokens(token1 lexer.Token, token2 lexer.Token, lexidic1 []any, lexidic2 []any) bool {
	return token1.Position == token2.Position &&
		//(token1.Position == 0 || token2.Position == 0 || lexidic1[token1.Position-1] == lexidic2[token2.Position-1]) &&
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
		fmt.Println(file)
		nameNoExt := file.Name()[0 : len(file.Name())-4]
		fmt.Println(len(expected[nameNoExt].tokens))
		testPassed := true
		fileLexer := reader.FileLexer("examples/" + file.Name())
		foundTokens, lexicon := fileLexer.Read()
		for ind, token := range foundTokens {
			if ind >= len(expected[nameNoExt].tokens) || !compareTokens(token, expected[nameNoExt].tokens[ind], lexicon, expected[nameNoExt].lexidic) {
				testPassed = false
				if ind >= len(expected[nameNoExt].tokens) {
					log.Fatal("\nThere is more token than expected")
				}
				log.Fatalf("\ntoken number: %d token gen: %v is different than token expected: %v", ind, token, expected[nameNoExt].tokens[ind])
			} else {
				fmt.Printf("token number: %d token: %v lexi: %v\n", ind, expected[nameNoExt].tokens, expected[nameNoExt].lexidic)
			}
		}
		if testPassed {
			fmt.Printf("Test %s: passed succesfully \n\n\n", file.Name())
		} else {
			fmt.Println("Test %s: not passed")
		}

	}
}
