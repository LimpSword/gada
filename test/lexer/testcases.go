package lexer

import (
	"gada/lexer"
	"gada/token"
)

type testlexer struct {
	tokens  []lexer.Token
	lexiDic []string
}

func getExpected() map[string]testlexer {
	expected := make(map[string]testlexer)

	// helloWorld
	tokens := make([]lexer.Token, 0)
	lexi := make([]string, 0)
	tokens = append(tokens, lexer.Token{"", 0, token.WITH, lexer.Position{1, 1}, lexer.Position{1, 5}})
	tokens = append(tokens, lexer.Token{"", 1, token.IDENT, lexer.Position{1, 6}, lexer.Position{1, 13}})
	tokens = append(tokens, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{1, 13}, lexer.Position{1, 14}})
	lexi = append(lexi, "Text_IO")
	expected["helloWorld"] = testlexer{
		tokens:  tokens,
		lexiDic: lexi}

	// errorChar
	tokens5 := make([]lexer.Token, 0)
	lexi5 := make([]string, 0)
	tokens5 = append(tokens5, lexer.Token{"", 1, token.IDENT, lexer.Position{1, 1}, lexer.Position{1, 4}})
	lexi5 = append(lexi5, "hey")
	tokens5 = append(tokens5, lexer.Token{"", 2, token.ILLEGAL, lexer.Position{1, 5}, lexer.Position{1, 12}})
	lexi5 = append(lexi5, "Lexical error: unexpected character 'gl hf' at line 1 between column 5 and 12.")
	expected["errorChar"] = testlexer{
		tokens:  tokens5,
		lexiDic: lexi5}

	// errorIllegalChar
	tokens6 := make([]lexer.Token, 0)
	lexi6 := make([]string, 0)
	tokens6 = append(tokens6, lexer.Token{"", 1, token.IDENT, lexer.Position{1, 1}, lexer.Position{1, 6}})
	lexi6 = append(lexi6, "notan")
	tokens6 = append(tokens6, lexer.Token{"", 2, token.ILLEGAL, lexer.Position{1, 6}, lexer.Position{1, 7}})
	lexi6 = append(lexi6, "Lexical error: unexpected character '$' at line 1 and column 6.")
	tokens6 = append(tokens6, lexer.Token{"", 3, token.IDENT, lexer.Position{1, 7}, lexer.Position{1, 12}})
	lexi6 = append(lexi6, "ident")
	expected["errorIllegalChar"] = testlexer{
		tokens:  tokens6,
		lexiDic: lexi6}

	// errorIllegalChar
	tokens7 := make([]lexer.Token, 0)
	lexi7 := make([]string, 0)
	tokens7 = append(tokens7, lexer.Token{"", 1, token.IDENT, lexer.Position{1, 1}, lexer.Position{1, 6}})
	lexi7 = append(lexi7, "notan")
	tokens7 = append(tokens7, lexer.Token{"", 2, token.ILLEGAL, lexer.Position{1, 6}, lexer.Position{1, 7}})
	lexi7 = append(lexi7, "Lexical error: unexpected character '$' at line 1 and column 6.")
	tokens7 = append(tokens7, lexer.Token{"", 3, token.IDENT, lexer.Position{1, 7}, lexer.Position{1, 12}})
	lexi7 = append(lexi7, "ident")
	expected["errorIllegalChar"] = testlexer{
		tokens:  tokens7,
		lexiDic: lexi7}

	// singlequote2
	tokens8 := make([]lexer.Token, 0)
	lexi8 := make([]string, 0)
	tokens8 = append(tokens8, lexer.Token{"", 1, token.ILLEGAL, lexer.Position{1, 1}, lexer.Position{1, 11}})
	lexi8 = append(lexi8, "Lexical error: new line in rune at line 1 and column 10.")
	tokens8 = append(tokens8, lexer.Token{"", 2, token.IDENT, lexer.Position{2, 1}, lexer.Position{2, 4}})
	lexi8 = append(lexi8, "hey")
	expected["singlequote2"] = testlexer{
		tokens:  tokens8,
		lexiDic: lexi8}

	// singlequote1
	tokens9 := make([]lexer.Token, 0)
	lexi9 := make([]string, 0)
	tokens9 = append(tokens9, lexer.Token{"", 0, token.CHAR_TOK, lexer.Position{1, 1}, lexer.Position{1, 10}})
	tokens9 = append(tokens9, lexer.Token{"", 0, token.CAST, lexer.Position{1, 11}, lexer.Position{1, 12}})
	tokens9 = append(tokens9, lexer.Token{"", 1, token.IDENT, lexer.Position{1, 13}, lexer.Position{1, 16}})
	lexi9 = append(lexi9, "val")
	tokens9 = append(tokens9, lexer.Token{"", 0, token.LPAREN, lexer.Position{1, 17}, lexer.Position{1, 18}})
	tokens9 = append(tokens9, lexer.Token{"", 2, token.INT, lexer.Position{1, 18}, lexer.Position{1, 19}})
	lexi9 = append(lexi9, "3")
	tokens9 = append(tokens9, lexer.Token{"", 0, token.RPAREN, lexer.Position{1, 19}, lexer.Position{1, 20}})
	expected["singlequote1"] = testlexer{
		tokens:  tokens9,
		lexiDic: lexi9}

	// firstline
	tokens3 := make([]lexer.Token, 0)
	lexi3 := make([]string, 0)
	tokens3 = append(tokens3, lexer.Token{"", 0, token.WITH, lexer.Position{1, 1}, lexer.Position{1, 5}})
	tokens3 = append(tokens3, lexer.Token{"", 1, token.IDENT, lexer.Position{1, 6}, lexer.Position{1, 9}})
	lexi3 = append(lexi3, "Ada")
	tokens3 = append(tokens3, lexer.Token{"", 0, token.PERIOD, lexer.Position{1, 9}, lexer.Position{1, 10}})
	tokens3 = append(tokens3, lexer.Token{"", 2, token.IDENT, lexer.Position{1, 10}, lexer.Position{1, 17}})
	lexi3 = append(lexi3, "Text_IO")
	tokens3 = append(tokens3, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{1, 18}, lexer.Position{1, 19}})
	tokens3 = append(tokens3, lexer.Token{"", 0, token.USE, lexer.Position{1, 20}, lexer.Position{1, 23}})
	tokens3 = append(tokens3, lexer.Token{"", 3, token.IDENT, lexer.Position{1, 24}, lexer.Position{1, 27}})
	lexi3 = append(lexi3, "Ada")
	tokens3 = append(tokens3, lexer.Token{"", 0, token.PERIOD, lexer.Position{1, 27}, lexer.Position{1, 28}})
	tokens3 = append(tokens3, lexer.Token{"", 4, token.IDENT, lexer.Position{1, 28}, lexer.Position{1, 35}})
	lexi3 = append(lexi3, "Text_IO")
	tokens3 = append(tokens3, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{1, 36}, lexer.Position{1, 37}})
	expected["firstLine"] = testlexer{
		tokens:  tokens3,
		lexiDic: lexi3}

	// inlineComment
	tokens2 := make([]lexer.Token, 0)
	lexi2 := make([]string, 0)
	// Tokens and positions for LINE 1 "with Text_IO; --use Text_IO;"
	tokens2 = append(tokens2, lexer.Token{"", 0, token.WITH, lexer.Position{1, 1}, lexer.Position{1, 5}})
	tokens2 = append(tokens2, lexer.Token{"", 1, token.IDENT, lexer.Position{1, 6}, lexer.Position{1, 13}})
	lexi2 = append(lexi2, "Text_IO")
	tokens2 = append(tokens2, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{1, 13}, lexer.Position{1, 14}})
	// Tokens and positions for LINE 2 "333 "let's"; -- random -- comment --doing--"
	tokens2 = append(tokens2, lexer.Token{"", 2, token.INT, lexer.Position{2, 1}, lexer.Position{2, 4}})
	lexi2 = append(lexi2, "333")
	tokens2 = append(tokens2, lexer.Token{"", 3, token.STRING, lexer.Position{2, 5}, lexer.Position{2, 12}})
	lexi2 = append(lexi2, "let's")
	tokens2 = append(tokens2, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{2, 12}, lexer.Position{2, 13}})
	// Tokens and positions for LINE 4 "      45.26;"
	tokens2 = append(tokens2, lexer.Token{"", 4, token.INT, lexer.Position{4, 7}, lexer.Position{4, 9}})
	lexi2 = append(lexi2, "45")
	tokens2 = append(tokens2, lexer.Token{"", 0, token.PERIOD, lexer.Position{4, 9}, lexer.Position{4, 10}})
	tokens2 = append(tokens2, lexer.Token{"", 5, token.INT, lexer.Position{4, 10}, lexer.Position{4, 12}})
	lexi2 = append(lexi2, "26")
	tokens2 = append(tokens2, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{4, 12}, lexer.Position{4, 13}})
	// Tokens and positions for LINE 5 "  hwy; ----------"
	tokens2 = append(tokens2, lexer.Token{"", 6, token.IDENT, lexer.Position{5, 3}, lexer.Position{5, 6}})
	lexi2 = append(lexi2, "hwy")
	tokens2 = append(tokens2, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{5, 6}, lexer.Position{5, 7}})

	expected["inlineComment"] = testlexer{
		tokens:  tokens2,
		lexiDic: lexi2}

	// geometry
	tokens1 := make([]lexer.Token, 0)
	lexi1 := make([]string, 0)

	// Tokens and positions for LINE 1 "with Text_IO ; use Text_IO ;"
	tokens1 = append(tokens1, lexer.Token{"", 0, token.WITH, lexer.Position{1, 1}, lexer.Position{1, 5}})
	tokens1 = append(tokens1, lexer.Token{"", 1, token.IDENT, lexer.Position{1, 6}, lexer.Position{1, 9}})
	lexi1 = append(lexi1, "Ada")
	tokens1 = append(tokens1, lexer.Token{"", 0, token.PERIOD, lexer.Position{1, 9}, lexer.Position{1, 10}})
	tokens1 = append(tokens1, lexer.Token{"", 2, token.IDENT, lexer.Position{1, 10}, lexer.Position{1, 17}})
	lexi1 = append(lexi1, "Text_IO")
	tokens1 = append(tokens1, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{1, 18}, lexer.Position{1, 19}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.USE, lexer.Position{1, 20}, lexer.Position{1, 23}})
	tokens1 = append(tokens1, lexer.Token{"", 3, token.IDENT, lexer.Position{1, 24}, lexer.Position{1, 27}})
	lexi1 = append(lexi1, "Ada")
	tokens1 = append(tokens1, lexer.Token{"", 0, token.PERIOD, lexer.Position{1, 27}, lexer.Position{1, 28}})
	tokens1 = append(tokens1, lexer.Token{"", 4, token.IDENT, lexer.Position{1, 28}, lexer.Position{1, 35}})
	lexi1 = append(lexi1, "Text_IO")
	tokens1 = append(tokens1, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{1, 36}, lexer.Position{1, 37}})

	// Tokens and positions for LINE 3 "procedure unDebut is"
	tokens1 = append(tokens1, lexer.Token{"", 0, token.PROCEDURE, lexer.Position{3, 1}, lexer.Position{3, 10}})
	lexi1 = append(lexi1, "unDebut") // Lexical position 2
	tokens1 = append(tokens1, lexer.Token{"", 5, token.IDENT, lexer.Position{3, 11}, lexer.Position{3, 18}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.IS, lexer.Position{3, 19}, lexer.Position{3, 21}})

	// Tokens and positions for line 5 function aireRectangle
	tokens1 = append(tokens1, lexer.Token{"", 0, token.FUNCTION, lexer.Position{5, 4}, lexer.Position{5, 12}})
	lexi1 = append(lexi1, "aireRectangle") // Lexical position 3
	tokens1 = append(tokens1, lexer.Token{"", 6, token.IDENT, lexer.Position{5, 13}, lexer.Position{5, 26}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.LPAREN, lexer.Position{5, 26}, lexer.Position{5, 27}})
	lexi1 = append(lexi1, "larg") // Lexical position 4
	tokens1 = append(tokens1, lexer.Token{"", 7, token.IDENT, lexer.Position{5, 27}, lexer.Position{5, 31}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.COLON, lexer.Position{5, 32}, lexer.Position{5, 33}})
	lexi1 = append(lexi1, "integer") // Lexical position 5
	tokens1 = append(tokens1, lexer.Token{"", 8, token.IDENT, lexer.Position{5, 34}, lexer.Position{5, 41}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{5, 41}, lexer.Position{5, 42}})
	lexi1 = append(lexi1, "long") // Lexical position 6
	tokens1 = append(tokens1, lexer.Token{"", 9, token.IDENT, lexer.Position{5, 43}, lexer.Position{5, 47}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.COLON, lexer.Position{5, 48}, lexer.Position{5, 49}})
	lexi1 = append(lexi1, "integer") // Lexical position 7
	tokens1 = append(tokens1, lexer.Token{"", 10, token.IDENT, lexer.Position{5, 50}, lexer.Position{5, 57}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.RPAREN, lexer.Position{5, 57}, lexer.Position{5, 58}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.RETURN, lexer.Position{5, 59}, lexer.Position{5, 65}})
	lexi1 = append(lexi1, "integer") // Lexical position 8
	tokens1 = append(tokens1, lexer.Token{"", 11, token.IDENT, lexer.Position{5, 66}, lexer.Position{5, 73}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.IS, lexer.Position{5, 74}, lexer.Position{5, 76}})

	// Tokens and positions for line 6 "aire : integer;"
	lexi1 = append(lexi1, "aire") // Lexical position 9
	tokens1 = append(tokens1, lexer.Token{"", 12, token.IDENT, lexer.Position{6, 4}, lexer.Position{6, 8}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.COLON, lexer.Position{6, 8}, lexer.Position{6, 9}})
	lexi1 = append(lexi1, "integer") // Lexical position 10
	tokens1 = append(tokens1, lexer.Token{"", 13, token.IDENT, lexer.Position{6, 10}, lexer.Position{6, 17}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{6, 17}, lexer.Position{6, 18}})

	// Tokens and positions for the line 7 "begin"
	tokens1 = append(tokens1, lexer.Token{"", 0, token.BEGIN, lexer.Position{7, 4}, lexer.Position{7, 9}})

	// Tokens and positions for the line 8 "aire := larg * long;"
	lexi1 = append(lexi1, "aire") // Lexical position 11
	tokens1 = append(tokens1, lexer.Token{"", 14, token.IDENT, lexer.Position{8, 7}, lexer.Position{8, 11}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.COLON, lexer.Position{8, 12}, lexer.Position{8, 13}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.EQL, lexer.Position{8, 13}, lexer.Position{8, 14}})
	lexi1 = append(lexi1, "larg") // Lexical position 12
	tokens1 = append(tokens1, lexer.Token{"", 15, token.IDENT, lexer.Position{8, 15}, lexer.Position{8, 19}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.MUL, lexer.Position{8, 19}, lexer.Position{8, 20}})
	lexi1 = append(lexi1, "long") // Lexical position 13
	tokens1 = append(tokens1, lexer.Token{"", 16, token.IDENT, lexer.Position{8, 20}, lexer.Position{8, 24}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{8, 25}, lexer.Position{8, 26}})

	// Tokens and positions for the line 9 "return aire"
	tokens1 = append(tokens1, lexer.Token{"", 0, token.RETURN, lexer.Position{9, 4}, lexer.Position{9, 10}})
	lexi1 = append(lexi1, "aire") // Lexical position 14
	tokens1 = append(tokens1, lexer.Token{"", 17, token.IDENT, lexer.Position{9, 11}, lexer.Position{9, 15}})

	// Tokens and positions for line 10 "end aireRectangle ;"
	tokens1 = append(tokens1, lexer.Token{"", 0, token.END, lexer.Position{10, 4}, lexer.Position{10, 7}})
	lexi1 = append(lexi1, "aireRectangle") // Lexical position 15
	tokens1 = append(tokens1, lexer.Token{"", 18, token.IDENT, lexer.Position{10, 8}, lexer.Position{10, 21}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{10, 22}, lexer.Position{10, 23}})

	// Tokens and positions for the line 12 "function perimetreRectangle"
	tokens1 = append(tokens1, lexer.Token{"", 0, token.FUNCTION, lexer.Position{12, 4}, lexer.Position{12, 12}})
	lexi1 = append(lexi1, "perimetreRectangle") // Lexical position 16
	tokens1 = append(tokens1, lexer.Token{"", 19, token.IDENT, lexer.Position{12, 13}, lexer.Position{12, 31}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.LPAREN, lexer.Position{12, 31}, lexer.Position{12, 32}})
	lexi1 = append(lexi1, "larg") // Lexical position 17
	tokens1 = append(tokens1, lexer.Token{"", 20, token.IDENT, lexer.Position{12, 32}, lexer.Position{12, 36}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.COLON, lexer.Position{12, 37}, lexer.Position{12, 38}})
	lexi1 = append(lexi1, "integer") // Lexical position 18
	tokens1 = append(tokens1, lexer.Token{"", 21, token.IDENT, lexer.Position{12, 39}, lexer.Position{12, 46}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{12, 46}, lexer.Position{12, 47}})
	lexi1 = append(lexi1, "long") // Lexical position 19
	tokens1 = append(tokens1, lexer.Token{"", 22, token.IDENT, lexer.Position{12, 48}, lexer.Position{12, 52}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.COLON, lexer.Position{12, 53}, lexer.Position{12, 54}})
	lexi1 = append(lexi1, "integer") // Lexical position 20
	tokens1 = append(tokens1, lexer.Token{"", 23, token.IDENT, lexer.Position{12, 55}, lexer.Position{12, 62}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.RPAREN, lexer.Position{12, 62}, lexer.Position{12, 63}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.RETURN, lexer.Position{12, 64}, lexer.Position{12, 70}})
	lexi1 = append(lexi1, "integer") // Lexical position 21
	tokens1 = append(tokens1, lexer.Token{"", 24, token.IDENT, lexer.Position{12, 71}, lexer.Position{12, 78}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.IS, lexer.Position{12, 79}, lexer.Position{12, 81}})

	// Tokens and positions for line 13 "p : integer;"
	lexi1 = append(lexi1, "p") // Lexical position 22
	tokens1 = append(tokens1, lexer.Token{"", 25, token.IDENT, lexer.Position{13, 4}, lexer.Position{13, 5}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.COLON, lexer.Position{13, 6}, lexer.Position{13, 7}})
	lexi1 = append(lexi1, "integer") // Lexical position 23
	tokens1 = append(tokens1, lexer.Token{"", 26, token.IDENT, lexer.Position{13, 8}, lexer.Position{13, 15}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{13, 15}, lexer.Position{13, 16}})

	// Tokens and positions for line 14 "begin"
	tokens1 = append(tokens1, lexer.Token{"", 0, token.BEGIN, lexer.Position{14, 4}, lexer.Position{14, 9}})

	// Tokens and positions for line 15 "p := 2 * (larg + long);"
	lexi1 = append(lexi1, "p") // Lexical position 24
	tokens1 = append(tokens1, lexer.Token{"", 27, token.IDENT, lexer.Position{15, 7}, lexer.Position{15, 8}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.COLON, lexer.Position{15, 9}, lexer.Position{15, 10}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.EQL, lexer.Position{15, 10}, lexer.Position{15, 11}})
	lexi1 = append(lexi1, "larg") // Lexical position 25
	tokens1 = append(tokens1, lexer.Token{"", 28, token.IDENT, lexer.Position{15, 12}, lexer.Position{15, 16}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.MUL, lexer.Position{15, 16}, lexer.Position{15, 17}})
	lexi1 = append(lexi1, "2") // Lexical position 26
	tokens1 = append(tokens1, lexer.Token{"", 29, token.INT, lexer.Position{15, 17}, lexer.Position{15, 18}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.ADD, lexer.Position{15, 19}, lexer.Position{15, 20}})
	lexi1 = append(lexi1, "long") // Lexical position 27
	tokens1 = append(tokens1, lexer.Token{"", 30, token.IDENT, lexer.Position{15, 21}, lexer.Position{15, 25}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.MUL, lexer.Position{15, 25}, lexer.Position{15, 26}})
	lexi1 = append(lexi1, "2") // Lexical position 28
	tokens1 = append(tokens1, lexer.Token{"", 31, token.INT, lexer.Position{15, 26}, lexer.Position{15, 27}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{15, 28}, lexer.Position{15, 29}})

	// Tokens and positions for line 16 "return p"
	tokens1 = append(tokens1, lexer.Token{"", 0, token.RETURN, lexer.Position{16, 4}, lexer.Position{16, 10}})
	lexi1 = append(lexi1, "p") // Lexical position 29
	tokens1 = append(tokens1, lexer.Token{"", 32, token.IDENT, lexer.Position{16, 11}, lexer.Position{16, 12}})

	// Tokens and positions for line 17 "end perimetreRectangle ;"
	tokens1 = append(tokens1, lexer.Token{"", 0, token.END, lexer.Position{17, 4}, lexer.Position{17, 7}})
	lexi1 = append(lexi1, "perimetreRectangle") // Lexical position 30
	tokens1 = append(tokens1, lexer.Token{"", 33, token.IDENT, lexer.Position{17, 8}, lexer.Position{17, 26}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{17, 26}, lexer.Position{17, 27}})

	// Tokens and positions for line 20 "choix : integer;"
	lexi1 = append(lexi1, "choix") // Lexical position 31
	tokens1 = append(tokens1, lexer.Token{"", 34, token.IDENT, lexer.Position{20, 1}, lexer.Position{20, 6}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.COLON, lexer.Position{20, 7}, lexer.Position{20, 8}})
	lexi1 = append(lexi1, "integer") // Lexical position 32
	tokens1 = append(tokens1, lexer.Token{"", 35, token.IDENT, lexer.Position{20, 9}, lexer.Position{20, 16}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{20, 17}, lexer.Position{20, 18}})

	// Tokens and positions for line 24 "begin"
	tokens1 = append(tokens1, lexer.Token{"", 0, token.BEGIN, lexer.Position{24, 1}, lexer.Position{24, 6}})

	// Tokens and positions for line 25 "choix := 2;"
	lexi1 = append(lexi1, "choix") // Lexical position 33
	tokens1 = append(tokens1, lexer.Token{"", 36, token.IDENT, lexer.Position{25, 4}, lexer.Position{25, 9}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.COLON, lexer.Position{25, 10}, lexer.Position{25, 11}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.EQL, lexer.Position{25, 11}, lexer.Position{25, 12}})
	lexi1 = append(lexi1, "2") // Lexical position 34
	tokens1 = append(tokens1, lexer.Token{"", 37, token.INT, lexer.Position{25, 13}, lexer.Position{25, 14}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{25, 14}, lexer.Position{25, 15}})

	// Tokens and positions for line 27 "if choix = 1"
	tokens1 = append(tokens1, lexer.Token{"", 0, token.IF, lexer.Position{27, 4}, lexer.Position{27, 6}})
	lexi1 = append(lexi1, "choix") // Lexical position 35
	tokens1 = append(tokens1, lexer.Token{"", 38, token.IDENT, lexer.Position{27, 7}, lexer.Position{27, 12}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.EQL, lexer.Position{27, 13}, lexer.Position{27, 14}})
	lexi1 = append(lexi1, "1") // Lexical position 36
	tokens1 = append(tokens1, lexer.Token{"", 39, token.INT, lexer.Position{27, 15}, lexer.Position{27, 16}})

	// Tokens and positions for line 28 "then valeur := permetreRectangle(2,3) ;"
	tokens1 = append(tokens1, lexer.Token{"", 0, token.THEN, lexer.Position{28, 7}, lexer.Position{28, 11}})
	lexi1 = append(lexi1, "valeur") // Lexical position 37
	tokens1 = append(tokens1, lexer.Token{"", 40, token.IDENT, lexer.Position{28, 12}, lexer.Position{28, 18}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.COLON, lexer.Position{28, 19}, lexer.Position{28, 20}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.EQL, lexer.Position{28, 20}, lexer.Position{28, 21}})
	lexi1 = append(lexi1, "perimetreRectangle") // Lexical position 38
	tokens1 = append(tokens1, lexer.Token{"", 41, token.IDENT, lexer.Position{28, 22}, lexer.Position{28, 40}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.LPAREN, lexer.Position{28, 40}, lexer.Position{28, 41}})
	lexi1 = append(lexi1, "2") // Lexical position 39
	tokens1 = append(tokens1, lexer.Token{"", 42, token.INT, lexer.Position{28, 41}, lexer.Position{28, 42}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.COMMA, lexer.Position{28, 42}, lexer.Position{28, 43}})
	lexi1 = append(lexi1, "3") // Lexical position 40
	tokens1 = append(tokens1, lexer.Token{"", 43, token.INT, lexer.Position{28, 43}, lexer.Position{28, 44}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.RPAREN, lexer.Position{28, 44}, lexer.Position{28, 45}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{28, 46}, lexer.Position{28, 47}})

	// Tokens and positions for line 29 "put(valeur) ;"
	lexi1 = append(lexi1, "put") // Lexical position 41
	tokens1 = append(tokens1, lexer.Token{"", 44, token.IDENT, lexer.Position{29, 10}, lexer.Position{29, 13}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.LPAREN, lexer.Position{29, 13}, lexer.Position{29, 14}})
	lexi1 = append(lexi1, "valeur") // Lexical position 42
	tokens1 = append(tokens1, lexer.Token{"", 45, token.IDENT, lexer.Position{29, 14}, lexer.Position{29, 20}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.RPAREN, lexer.Position{29, 20}, lexer.Position{29, 21}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{29, 22}, lexer.Position{29, 23}})

	// Tokens and positions for line 30 "else valeur := aireRectangle(2,3) ;"
	tokens1 = append(tokens1, lexer.Token{"", 0, token.ELSE, lexer.Position{30, 7}, lexer.Position{30, 11}})
	lexi1 = append(lexi1, "valeur") // Lexical position 43
	tokens1 = append(tokens1, lexer.Token{"", 46, token.IDENT, lexer.Position{30, 12}, lexer.Position{30, 18}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.COLON, lexer.Position{30, 19}, lexer.Position{30, 20}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.EQL, lexer.Position{30, 20}, lexer.Position{30, 21}})
	lexi1 = append(lexi1, "aireRectangle") // Lexical position 44
	tokens1 = append(tokens1, lexer.Token{"", 47, token.IDENT, lexer.Position{30, 22}, lexer.Position{30, 35}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.LPAREN, lexer.Position{30, 35}, lexer.Position{30, 36}})
	lexi1 = append(lexi1, "2") // Lexical position 45
	tokens1 = append(tokens1, lexer.Token{"", 48, token.INT, lexer.Position{30, 36}, lexer.Position{30, 37}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.COMMA, lexer.Position{30, 37}, lexer.Position{30, 38}})
	lexi1 = append(lexi1, "3") // Lexical position 46
	tokens1 = append(tokens1, lexer.Token{"", 49, token.INT, lexer.Position{30, 38}, lexer.Position{30, 39}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.RPAREN, lexer.Position{30, 39}, lexer.Position{30, 40}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{30, 41}, lexer.Position{30, 42}})

	// Tokens and positions for line 31 "put(valeur) ;"
	lexi1 = append(lexi1, "put") // Lexical position 47
	tokens1 = append(tokens1, lexer.Token{"", 50, token.IDENT, lexer.Position{31, 10}, lexer.Position{31, 13}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.LPAREN, lexer.Position{31, 13}, lexer.Position{31, 14}})
	lexi1 = append(lexi1, "valeur") // Lexical position 48
	tokens1 = append(tokens1, lexer.Token{"", 51, token.IDENT, lexer.Position{31, 14}, lexer.Position{31, 20}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.RPAREN, lexer.Position{31, 20}, lexer.Position{31, 21}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{31, 22}, lexer.Position{31, 23}})

	// Tokens and positions for line 32 "end if;"
	tokens1 = append(tokens1, lexer.Token{"", 0, token.END, lexer.Position{32, 4}, lexer.Position{32, 7}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.IF, lexer.Position{32, 8}, lexer.Position{32, 10}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{32, 10}, lexer.Position{32, 11}})

	// Tokens and positions for line 33 "end unDebut ;"
	tokens1 = append(tokens1, lexer.Token{"", 0, token.END, lexer.Position{33, 1}, lexer.Position{33, 4}})
	lexi1 = append(lexi1, "unDebut") // Lexical position 49
	tokens1 = append(tokens1, lexer.Token{"", 52, token.IDENT, lexer.Position{33, 5}, lexer.Position{33, 12}})
	tokens1 = append(tokens1, lexer.Token{"", 0, token.SEMICOLON, lexer.Position{33, 13}, lexer.Position{33, 14}})

	expected["geometry"] = testlexer{
		tokens:  tokens1,
		lexiDic: lexi1,
	}

	return expected
}

func getExpected2() map[string]testlexer {
	expected := make(map[string]testlexer)

	return expected
}
