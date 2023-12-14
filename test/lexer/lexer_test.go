package lexer

import (
	"gada/lexer"
	"gada/reader"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strings"
	"testing"
)

func compareTokens(token1 lexer.Token, token2 lexer.Token, lexiDic1 []string, lexiDic2 []string) bool {
	return token1.Position == token2.Position &&
		(token1.Position == 0 || lexiDic1[token1.Position-1] == lexiDic2[token2.Position-1]) &&
		// previous condition check if the literals are equal in case where the token is a literal
		token1.Value == token2.Value &&
		token1.Beginning == token2.Beginning &&
		token1.End == token2.End
}

func TestAll(t *testing.T) {
	files, err := os.ReadDir("tests")
	if err != nil {
		log.Fatalf("the directory provided have this error : %s", err)
	}
	expected := getExpected()
	for _, file := range files {
		t.Logf("Test %s beginning\n", file.Name())
		nameNoExt := strings.Split(file.Name(), ".")[0]
		fileLexer := reader.FileLexer("tests/" + file.Name())
		foundTokens, lexicon := fileLexer.Read()
		expecTokens, expecLexi := expected[nameNoExt].tokens, expected[nameNoExt].lexiDic

		assert.Equalf(t, len(foundTokens), len(expecTokens), "The token count doesn't match in file %s", file.Name())

		for ind, tok := range foundTokens {

			tokenLit1, tokenLit2 := "", ""
			if tok.Position != 0 {
				tokenLit1 = lexicon[tok.Position-1]
			}
			if expecTokens[ind].Position != 0 {
				tokenLit2 = expecLexi[expecTokens[ind].Position-1]
			}

			assert.Truef(t, compareTokens(tok, expecTokens[ind], lexicon, expecLexi),
				"The token doesn't match in file %s, token number: %d, token gen: %v %s is different than token expected: %v %s",
				file.Name(), ind, tok, tokenLit1, expecTokens[ind], tokenLit2)
		}

		t.Logf("Test %s ending\n", file.Name())
	}
}
