// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package sf

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type SyntaxFlowLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var syntaxflowlexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	channelNames           []string
	modeNames              []string
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func syntaxflowlexerLexerInit() {
	staticData := &syntaxflowlexerLexerStaticData
	staticData.channelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.modeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.literalNames = []string{
		"", "'->'", "'-->'", "';'", "'==>'", "'...'", "'%%'", "'..'", "'<='",
		"'>='", "'>>'", "'=>'", "'=='", "'=~'", "'!~'", "'&&'", "'||'", "'!='",
		"'>'", "'.'", "'<'", "'='", "'?'", "'('", "','", "')'", "'['", "']'",
		"'{'", "'}'", "'#'", "'$'", "':'", "'%'", "'!'", "'*'", "'-'", "'as'",
		"", "", "", "", "", "", "'str'", "'list'", "'dict'", "", "'bool'",
	}
	staticData.symbolicNames = []string{
		"", "", "", "", "DeepFilter", "Deep", "Percent", "DeepDot", "LtEq",
		"GtEq", "DoubleGt", "Filter", "EqEq", "RegexpMatch", "NotRegexpMatch",
		"And", "Or", "NotEq", "Gt", "Dot", "Lt", "Eq", "Question", "OpenParen",
		"Comma", "CloseParen", "ListSelectOpen", "ListSelectClose", "MapBuilderOpen",
		"MapBuilderClose", "ListStart", "DollarOutput", "Colon", "Search", "Bang",
		"Star", "Minus", "As", "WhiteSpace", "Number", "OctalNumber", "BinaryNumber",
		"HexNumber", "StringLiteral", "StringType", "ListType", "DictType",
		"NumberType", "BoolType", "BoolLiteral", "Identifier", "IdentifierChar",
		"RegexpLiteral",
	}
	staticData.ruleNames = []string{
		"T__0", "T__1", "T__2", "DeepFilter", "Deep", "Percent", "DeepDot",
		"LtEq", "GtEq", "DoubleGt", "Filter", "EqEq", "RegexpMatch", "NotRegexpMatch",
		"And", "Or", "NotEq", "Gt", "Dot", "Lt", "Eq", "Question", "OpenParen",
		"Comma", "CloseParen", "ListSelectOpen", "ListSelectClose", "MapBuilderOpen",
		"MapBuilderClose", "ListStart", "DollarOutput", "Colon", "Search", "Bang",
		"Star", "Minus", "As", "WhiteSpace", "Number", "OctalNumber", "BinaryNumber",
		"HexNumber", "StringLiteral", "StringType", "ListType", "DictType",
		"NumberType", "BoolType", "BoolLiteral", "Identifier", "IdentifierChar",
		"IdentifierCharStart", "HexDigit", "Digit", "OctalDigit", "RegexpLiteral",
		"RegexpLiteralChar",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 52, 324, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7,
		20, 2, 21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25,
		2, 26, 7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2,
		31, 7, 31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36,
		7, 36, 2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7,
		41, 2, 42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2, 46, 7, 46,
		2, 47, 7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 2, 50, 7, 50, 2, 51, 7, 51, 2,
		52, 7, 52, 2, 53, 7, 53, 2, 54, 7, 54, 2, 55, 7, 55, 2, 56, 7, 56, 1, 0,
		1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 1, 3,
		1, 4, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7,
		1, 7, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 1, 11, 1,
		11, 1, 11, 1, 12, 1, 12, 1, 12, 1, 13, 1, 13, 1, 13, 1, 14, 1, 14, 1, 14,
		1, 15, 1, 15, 1, 15, 1, 16, 1, 16, 1, 16, 1, 17, 1, 17, 1, 18, 1, 18, 1,
		19, 1, 19, 1, 20, 1, 20, 1, 21, 1, 21, 1, 22, 1, 22, 1, 23, 1, 23, 1, 24,
		1, 24, 1, 25, 1, 25, 1, 26, 1, 26, 1, 27, 1, 27, 1, 28, 1, 28, 1, 29, 1,
		29, 1, 30, 1, 30, 1, 31, 1, 31, 1, 32, 1, 32, 1, 33, 1, 33, 1, 34, 1, 34,
		1, 35, 1, 35, 1, 36, 1, 36, 1, 36, 1, 37, 1, 37, 1, 37, 1, 37, 1, 38, 4,
		38, 215, 8, 38, 11, 38, 12, 38, 216, 1, 39, 1, 39, 1, 39, 1, 39, 4, 39,
		223, 8, 39, 11, 39, 12, 39, 224, 1, 40, 1, 40, 1, 40, 1, 40, 4, 40, 231,
		8, 40, 11, 40, 12, 40, 232, 1, 41, 1, 41, 1, 41, 1, 41, 4, 41, 239, 8,
		41, 11, 41, 12, 41, 240, 1, 42, 1, 42, 5, 42, 245, 8, 42, 10, 42, 12, 42,
		248, 9, 42, 1, 42, 1, 42, 1, 43, 1, 43, 1, 43, 1, 43, 1, 44, 1, 44, 1,
		44, 1, 44, 1, 44, 1, 45, 1, 45, 1, 45, 1, 45, 1, 45, 1, 46, 1, 46, 1, 46,
		1, 46, 1, 46, 1, 46, 1, 46, 1, 46, 3, 46, 274, 8, 46, 1, 47, 1, 47, 1,
		47, 1, 47, 1, 47, 1, 48, 1, 48, 1, 48, 1, 48, 1, 48, 1, 48, 1, 48, 1, 48,
		1, 48, 3, 48, 290, 8, 48, 1, 49, 1, 49, 5, 49, 294, 8, 49, 10, 49, 12,
		49, 297, 9, 49, 1, 50, 1, 50, 3, 50, 301, 8, 50, 1, 51, 3, 51, 304, 8,
		51, 1, 52, 1, 52, 1, 53, 1, 53, 1, 54, 1, 54, 1, 55, 1, 55, 4, 55, 314,
		8, 55, 11, 55, 12, 55, 315, 1, 55, 1, 55, 1, 56, 1, 56, 1, 56, 3, 56, 323,
		8, 56, 0, 0, 57, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8, 17,
		9, 19, 10, 21, 11, 23, 12, 25, 13, 27, 14, 29, 15, 31, 16, 33, 17, 35,
		18, 37, 19, 39, 20, 41, 21, 43, 22, 45, 23, 47, 24, 49, 25, 51, 26, 53,
		27, 55, 28, 57, 29, 59, 30, 61, 31, 63, 32, 65, 33, 67, 34, 69, 35, 71,
		36, 73, 37, 75, 38, 77, 39, 79, 40, 81, 41, 83, 42, 85, 43, 87, 44, 89,
		45, 91, 46, 93, 47, 95, 48, 97, 49, 99, 50, 101, 51, 103, 0, 105, 0, 107,
		0, 109, 0, 111, 52, 113, 0, 1, 0, 7, 3, 0, 10, 10, 13, 13, 32, 32, 1, 0,
		96, 96, 1, 0, 48, 57, 4, 0, 42, 42, 65, 90, 95, 95, 97, 122, 3, 0, 48,
		57, 65, 70, 97, 102, 1, 0, 48, 55, 1, 0, 47, 47, 329, 0, 1, 1, 0, 0, 0,
		0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0,
		0, 11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0,
		0, 0, 19, 1, 0, 0, 0, 0, 21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0, 0, 25, 1, 0,
		0, 0, 0, 27, 1, 0, 0, 0, 0, 29, 1, 0, 0, 0, 0, 31, 1, 0, 0, 0, 0, 33, 1,
		0, 0, 0, 0, 35, 1, 0, 0, 0, 0, 37, 1, 0, 0, 0, 0, 39, 1, 0, 0, 0, 0, 41,
		1, 0, 0, 0, 0, 43, 1, 0, 0, 0, 0, 45, 1, 0, 0, 0, 0, 47, 1, 0, 0, 0, 0,
		49, 1, 0, 0, 0, 0, 51, 1, 0, 0, 0, 0, 53, 1, 0, 0, 0, 0, 55, 1, 0, 0, 0,
		0, 57, 1, 0, 0, 0, 0, 59, 1, 0, 0, 0, 0, 61, 1, 0, 0, 0, 0, 63, 1, 0, 0,
		0, 0, 65, 1, 0, 0, 0, 0, 67, 1, 0, 0, 0, 0, 69, 1, 0, 0, 0, 0, 71, 1, 0,
		0, 0, 0, 73, 1, 0, 0, 0, 0, 75, 1, 0, 0, 0, 0, 77, 1, 0, 0, 0, 0, 79, 1,
		0, 0, 0, 0, 81, 1, 0, 0, 0, 0, 83, 1, 0, 0, 0, 0, 85, 1, 0, 0, 0, 0, 87,
		1, 0, 0, 0, 0, 89, 1, 0, 0, 0, 0, 91, 1, 0, 0, 0, 0, 93, 1, 0, 0, 0, 0,
		95, 1, 0, 0, 0, 0, 97, 1, 0, 0, 0, 0, 99, 1, 0, 0, 0, 0, 101, 1, 0, 0,
		0, 0, 111, 1, 0, 0, 0, 1, 115, 1, 0, 0, 0, 3, 118, 1, 0, 0, 0, 5, 122,
		1, 0, 0, 0, 7, 124, 1, 0, 0, 0, 9, 128, 1, 0, 0, 0, 11, 132, 1, 0, 0, 0,
		13, 135, 1, 0, 0, 0, 15, 138, 1, 0, 0, 0, 17, 141, 1, 0, 0, 0, 19, 144,
		1, 0, 0, 0, 21, 147, 1, 0, 0, 0, 23, 150, 1, 0, 0, 0, 25, 153, 1, 0, 0,
		0, 27, 156, 1, 0, 0, 0, 29, 159, 1, 0, 0, 0, 31, 162, 1, 0, 0, 0, 33, 165,
		1, 0, 0, 0, 35, 168, 1, 0, 0, 0, 37, 170, 1, 0, 0, 0, 39, 172, 1, 0, 0,
		0, 41, 174, 1, 0, 0, 0, 43, 176, 1, 0, 0, 0, 45, 178, 1, 0, 0, 0, 47, 180,
		1, 0, 0, 0, 49, 182, 1, 0, 0, 0, 51, 184, 1, 0, 0, 0, 53, 186, 1, 0, 0,
		0, 55, 188, 1, 0, 0, 0, 57, 190, 1, 0, 0, 0, 59, 192, 1, 0, 0, 0, 61, 194,
		1, 0, 0, 0, 63, 196, 1, 0, 0, 0, 65, 198, 1, 0, 0, 0, 67, 200, 1, 0, 0,
		0, 69, 202, 1, 0, 0, 0, 71, 204, 1, 0, 0, 0, 73, 206, 1, 0, 0, 0, 75, 209,
		1, 0, 0, 0, 77, 214, 1, 0, 0, 0, 79, 218, 1, 0, 0, 0, 81, 226, 1, 0, 0,
		0, 83, 234, 1, 0, 0, 0, 85, 242, 1, 0, 0, 0, 87, 251, 1, 0, 0, 0, 89, 255,
		1, 0, 0, 0, 91, 260, 1, 0, 0, 0, 93, 273, 1, 0, 0, 0, 95, 275, 1, 0, 0,
		0, 97, 289, 1, 0, 0, 0, 99, 291, 1, 0, 0, 0, 101, 300, 1, 0, 0, 0, 103,
		303, 1, 0, 0, 0, 105, 305, 1, 0, 0, 0, 107, 307, 1, 0, 0, 0, 109, 309,
		1, 0, 0, 0, 111, 311, 1, 0, 0, 0, 113, 322, 1, 0, 0, 0, 115, 116, 5, 45,
		0, 0, 116, 117, 5, 62, 0, 0, 117, 2, 1, 0, 0, 0, 118, 119, 5, 45, 0, 0,
		119, 120, 5, 45, 0, 0, 120, 121, 5, 62, 0, 0, 121, 4, 1, 0, 0, 0, 122,
		123, 5, 59, 0, 0, 123, 6, 1, 0, 0, 0, 124, 125, 5, 61, 0, 0, 125, 126,
		5, 61, 0, 0, 126, 127, 5, 62, 0, 0, 127, 8, 1, 0, 0, 0, 128, 129, 5, 46,
		0, 0, 129, 130, 5, 46, 0, 0, 130, 131, 5, 46, 0, 0, 131, 10, 1, 0, 0, 0,
		132, 133, 5, 37, 0, 0, 133, 134, 5, 37, 0, 0, 134, 12, 1, 0, 0, 0, 135,
		136, 5, 46, 0, 0, 136, 137, 5, 46, 0, 0, 137, 14, 1, 0, 0, 0, 138, 139,
		5, 60, 0, 0, 139, 140, 5, 61, 0, 0, 140, 16, 1, 0, 0, 0, 141, 142, 5, 62,
		0, 0, 142, 143, 5, 61, 0, 0, 143, 18, 1, 0, 0, 0, 144, 145, 5, 62, 0, 0,
		145, 146, 5, 62, 0, 0, 146, 20, 1, 0, 0, 0, 147, 148, 5, 61, 0, 0, 148,
		149, 5, 62, 0, 0, 149, 22, 1, 0, 0, 0, 150, 151, 5, 61, 0, 0, 151, 152,
		5, 61, 0, 0, 152, 24, 1, 0, 0, 0, 153, 154, 5, 61, 0, 0, 154, 155, 5, 126,
		0, 0, 155, 26, 1, 0, 0, 0, 156, 157, 5, 33, 0, 0, 157, 158, 5, 126, 0,
		0, 158, 28, 1, 0, 0, 0, 159, 160, 5, 38, 0, 0, 160, 161, 5, 38, 0, 0, 161,
		30, 1, 0, 0, 0, 162, 163, 5, 124, 0, 0, 163, 164, 5, 124, 0, 0, 164, 32,
		1, 0, 0, 0, 165, 166, 5, 33, 0, 0, 166, 167, 5, 61, 0, 0, 167, 34, 1, 0,
		0, 0, 168, 169, 5, 62, 0, 0, 169, 36, 1, 0, 0, 0, 170, 171, 5, 46, 0, 0,
		171, 38, 1, 0, 0, 0, 172, 173, 5, 60, 0, 0, 173, 40, 1, 0, 0, 0, 174, 175,
		5, 61, 0, 0, 175, 42, 1, 0, 0, 0, 176, 177, 5, 63, 0, 0, 177, 44, 1, 0,
		0, 0, 178, 179, 5, 40, 0, 0, 179, 46, 1, 0, 0, 0, 180, 181, 5, 44, 0, 0,
		181, 48, 1, 0, 0, 0, 182, 183, 5, 41, 0, 0, 183, 50, 1, 0, 0, 0, 184, 185,
		5, 91, 0, 0, 185, 52, 1, 0, 0, 0, 186, 187, 5, 93, 0, 0, 187, 54, 1, 0,
		0, 0, 188, 189, 5, 123, 0, 0, 189, 56, 1, 0, 0, 0, 190, 191, 5, 125, 0,
		0, 191, 58, 1, 0, 0, 0, 192, 193, 5, 35, 0, 0, 193, 60, 1, 0, 0, 0, 194,
		195, 5, 36, 0, 0, 195, 62, 1, 0, 0, 0, 196, 197, 5, 58, 0, 0, 197, 64,
		1, 0, 0, 0, 198, 199, 5, 37, 0, 0, 199, 66, 1, 0, 0, 0, 200, 201, 5, 33,
		0, 0, 201, 68, 1, 0, 0, 0, 202, 203, 5, 42, 0, 0, 203, 70, 1, 0, 0, 0,
		204, 205, 5, 45, 0, 0, 205, 72, 1, 0, 0, 0, 206, 207, 5, 97, 0, 0, 207,
		208, 5, 115, 0, 0, 208, 74, 1, 0, 0, 0, 209, 210, 7, 0, 0, 0, 210, 211,
		1, 0, 0, 0, 211, 212, 6, 37, 0, 0, 212, 76, 1, 0, 0, 0, 213, 215, 3, 107,
		53, 0, 214, 213, 1, 0, 0, 0, 215, 216, 1, 0, 0, 0, 216, 214, 1, 0, 0, 0,
		216, 217, 1, 0, 0, 0, 217, 78, 1, 0, 0, 0, 218, 219, 5, 48, 0, 0, 219,
		220, 5, 111, 0, 0, 220, 222, 1, 0, 0, 0, 221, 223, 3, 109, 54, 0, 222,
		221, 1, 0, 0, 0, 223, 224, 1, 0, 0, 0, 224, 222, 1, 0, 0, 0, 224, 225,
		1, 0, 0, 0, 225, 80, 1, 0, 0, 0, 226, 227, 5, 48, 0, 0, 227, 228, 5, 98,
		0, 0, 228, 230, 1, 0, 0, 0, 229, 231, 2, 48, 49, 0, 230, 229, 1, 0, 0,
		0, 231, 232, 1, 0, 0, 0, 232, 230, 1, 0, 0, 0, 232, 233, 1, 0, 0, 0, 233,
		82, 1, 0, 0, 0, 234, 235, 5, 48, 0, 0, 235, 236, 5, 120, 0, 0, 236, 238,
		1, 0, 0, 0, 237, 239, 3, 105, 52, 0, 238, 237, 1, 0, 0, 0, 239, 240, 1,
		0, 0, 0, 240, 238, 1, 0, 0, 0, 240, 241, 1, 0, 0, 0, 241, 84, 1, 0, 0,
		0, 242, 246, 5, 96, 0, 0, 243, 245, 8, 1, 0, 0, 244, 243, 1, 0, 0, 0, 245,
		248, 1, 0, 0, 0, 246, 244, 1, 0, 0, 0, 246, 247, 1, 0, 0, 0, 247, 249,
		1, 0, 0, 0, 248, 246, 1, 0, 0, 0, 249, 250, 5, 96, 0, 0, 250, 86, 1, 0,
		0, 0, 251, 252, 5, 115, 0, 0, 252, 253, 5, 116, 0, 0, 253, 254, 5, 114,
		0, 0, 254, 88, 1, 0, 0, 0, 255, 256, 5, 108, 0, 0, 256, 257, 5, 105, 0,
		0, 257, 258, 5, 115, 0, 0, 258, 259, 5, 116, 0, 0, 259, 90, 1, 0, 0, 0,
		260, 261, 5, 100, 0, 0, 261, 262, 5, 105, 0, 0, 262, 263, 5, 99, 0, 0,
		263, 264, 5, 116, 0, 0, 264, 92, 1, 0, 0, 0, 265, 266, 5, 105, 0, 0, 266,
		267, 5, 110, 0, 0, 267, 274, 5, 116, 0, 0, 268, 269, 5, 102, 0, 0, 269,
		270, 5, 108, 0, 0, 270, 271, 5, 111, 0, 0, 271, 272, 5, 97, 0, 0, 272,
		274, 5, 116, 0, 0, 273, 265, 1, 0, 0, 0, 273, 268, 1, 0, 0, 0, 274, 94,
		1, 0, 0, 0, 275, 276, 5, 98, 0, 0, 276, 277, 5, 111, 0, 0, 277, 278, 5,
		111, 0, 0, 278, 279, 5, 108, 0, 0, 279, 96, 1, 0, 0, 0, 280, 281, 5, 116,
		0, 0, 281, 282, 5, 114, 0, 0, 282, 283, 5, 117, 0, 0, 283, 290, 5, 101,
		0, 0, 284, 285, 5, 102, 0, 0, 285, 286, 5, 97, 0, 0, 286, 287, 5, 108,
		0, 0, 287, 288, 5, 115, 0, 0, 288, 290, 5, 101, 0, 0, 289, 280, 1, 0, 0,
		0, 289, 284, 1, 0, 0, 0, 290, 98, 1, 0, 0, 0, 291, 295, 3, 103, 51, 0,
		292, 294, 3, 101, 50, 0, 293, 292, 1, 0, 0, 0, 294, 297, 1, 0, 0, 0, 295,
		293, 1, 0, 0, 0, 295, 296, 1, 0, 0, 0, 296, 100, 1, 0, 0, 0, 297, 295,
		1, 0, 0, 0, 298, 301, 7, 2, 0, 0, 299, 301, 3, 103, 51, 0, 300, 298, 1,
		0, 0, 0, 300, 299, 1, 0, 0, 0, 301, 102, 1, 0, 0, 0, 302, 304, 7, 3, 0,
		0, 303, 302, 1, 0, 0, 0, 304, 104, 1, 0, 0, 0, 305, 306, 7, 4, 0, 0, 306,
		106, 1, 0, 0, 0, 307, 308, 7, 2, 0, 0, 308, 108, 1, 0, 0, 0, 309, 310,
		7, 5, 0, 0, 310, 110, 1, 0, 0, 0, 311, 313, 5, 47, 0, 0, 312, 314, 3, 113,
		56, 0, 313, 312, 1, 0, 0, 0, 314, 315, 1, 0, 0, 0, 315, 313, 1, 0, 0, 0,
		315, 316, 1, 0, 0, 0, 316, 317, 1, 0, 0, 0, 317, 318, 5, 47, 0, 0, 318,
		112, 1, 0, 0, 0, 319, 320, 5, 92, 0, 0, 320, 323, 5, 47, 0, 0, 321, 323,
		8, 6, 0, 0, 322, 319, 1, 0, 0, 0, 322, 321, 1, 0, 0, 0, 323, 114, 1, 0,
		0, 0, 13, 0, 216, 224, 232, 240, 246, 273, 289, 295, 300, 303, 315, 322,
		1, 6, 0, 0,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// SyntaxFlowLexerInit initializes any static state used to implement SyntaxFlowLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewSyntaxFlowLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func SyntaxFlowLexerInit() {
	staticData := &syntaxflowlexerLexerStaticData
	staticData.once.Do(syntaxflowlexerLexerInit)
}

// NewSyntaxFlowLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewSyntaxFlowLexer(input antlr.CharStream) *SyntaxFlowLexer {
	SyntaxFlowLexerInit()
	l := new(SyntaxFlowLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &syntaxflowlexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	l.channelNames = staticData.channelNames
	l.modeNames = staticData.modeNames
	l.RuleNames = staticData.ruleNames
	l.LiteralNames = staticData.literalNames
	l.SymbolicNames = staticData.symbolicNames
	l.GrammarFileName = "SyntaxFlow.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// SyntaxFlowLexer tokens.
const (
	SyntaxFlowLexerT__0            = 1
	SyntaxFlowLexerT__1            = 2
	SyntaxFlowLexerT__2            = 3
	SyntaxFlowLexerDeepFilter      = 4
	SyntaxFlowLexerDeep            = 5
	SyntaxFlowLexerPercent         = 6
	SyntaxFlowLexerDeepDot         = 7
	SyntaxFlowLexerLtEq            = 8
	SyntaxFlowLexerGtEq            = 9
	SyntaxFlowLexerDoubleGt        = 10
	SyntaxFlowLexerFilter          = 11
	SyntaxFlowLexerEqEq            = 12
	SyntaxFlowLexerRegexpMatch     = 13
	SyntaxFlowLexerNotRegexpMatch  = 14
	SyntaxFlowLexerAnd             = 15
	SyntaxFlowLexerOr              = 16
	SyntaxFlowLexerNotEq           = 17
	SyntaxFlowLexerGt              = 18
	SyntaxFlowLexerDot             = 19
	SyntaxFlowLexerLt              = 20
	SyntaxFlowLexerEq              = 21
	SyntaxFlowLexerQuestion        = 22
	SyntaxFlowLexerOpenParen       = 23
	SyntaxFlowLexerComma           = 24
	SyntaxFlowLexerCloseParen      = 25
	SyntaxFlowLexerListSelectOpen  = 26
	SyntaxFlowLexerListSelectClose = 27
	SyntaxFlowLexerMapBuilderOpen  = 28
	SyntaxFlowLexerMapBuilderClose = 29
	SyntaxFlowLexerListStart       = 30
	SyntaxFlowLexerDollarOutput    = 31
	SyntaxFlowLexerColon           = 32
	SyntaxFlowLexerSearch          = 33
	SyntaxFlowLexerBang            = 34
	SyntaxFlowLexerStar            = 35
	SyntaxFlowLexerMinus           = 36
	SyntaxFlowLexerAs              = 37
	SyntaxFlowLexerWhiteSpace      = 38
	SyntaxFlowLexerNumber          = 39
	SyntaxFlowLexerOctalNumber     = 40
	SyntaxFlowLexerBinaryNumber    = 41
	SyntaxFlowLexerHexNumber       = 42
	SyntaxFlowLexerStringLiteral   = 43
	SyntaxFlowLexerStringType      = 44
	SyntaxFlowLexerListType        = 45
	SyntaxFlowLexerDictType        = 46
	SyntaxFlowLexerNumberType      = 47
	SyntaxFlowLexerBoolType        = 48
	SyntaxFlowLexerBoolLiteral     = 49
	SyntaxFlowLexerIdentifier      = 50
	SyntaxFlowLexerIdentifierChar  = 51
	SyntaxFlowLexerRegexpLiteral   = 52
)
