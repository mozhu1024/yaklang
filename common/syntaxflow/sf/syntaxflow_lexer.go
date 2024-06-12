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
		"", "';'", "'->'", "'-->'", "'-<'", "'>-'", "'==>'", "'...'", "'%%'",
		"'..'", "'<='", "'>='", "'>>'", "'=>'", "'=='", "'=~'", "'!~'", "'&&'",
		"'||'", "'!='", "'?{'", "'-{'", "'}->'", "'#{'", "'#>'", "'#->'", "'>'",
		"'.'", "'<'", "'='", "'?'", "'('", "','", "')'", "'['", "']'", "'{'",
		"'}'", "'#'", "'$'", "':'", "'%'", "'!'", "'*'", "'-'", "'as'", "'`'",
		"", "", "", "", "", "'str'", "'list'", "'dict'", "", "'bool'", "", "'assert'",
		"'then'", "", "'else'",
	}
	staticData.symbolicNames = []string{
		"", "", "", "", "", "", "DeepFilter", "Deep", "Percent", "DeepDot",
		"LtEq", "GtEq", "DoubleGt", "Filter", "EqEq", "RegexpMatch", "NotRegexpMatch",
		"And", "Or", "NotEq", "ConditionStart", "DeepNextStart", "DeepNextEnd",
		"TopDefStart", "DefStart", "TopDef", "Gt", "Dot", "Lt", "Eq", "Question",
		"OpenParen", "Comma", "CloseParen", "ListSelectOpen", "ListSelectClose",
		"MapBuilderOpen", "MapBuilderClose", "ListStart", "DollarOutput", "Colon",
		"Search", "Bang", "Star", "Minus", "As", "Backtick", "WhiteSpace", "Number",
		"OctalNumber", "BinaryNumber", "HexNumber", "StringType", "ListType",
		"DictType", "NumberType", "BoolType", "BoolLiteral", "Assert", "Then",
		"Desc", "Else", "Identifier", "IdentifierChar", "RegexpLiteral", "WS",
	}
	staticData.ruleNames = []string{
		"T__0", "T__1", "T__2", "T__3", "T__4", "DeepFilter", "Deep", "Percent",
		"DeepDot", "LtEq", "GtEq", "DoubleGt", "Filter", "EqEq", "RegexpMatch",
		"NotRegexpMatch", "And", "Or", "NotEq", "ConditionStart", "DeepNextStart",
		"DeepNextEnd", "TopDefStart", "DefStart", "TopDef", "Gt", "Dot", "Lt",
		"Eq", "Question", "OpenParen", "Comma", "CloseParen", "ListSelectOpen",
		"ListSelectClose", "MapBuilderOpen", "MapBuilderClose", "ListStart",
		"DollarOutput", "Colon", "Search", "Bang", "Star", "Minus", "As", "Backtick",
		"WhiteSpace", "Number", "OctalNumber", "BinaryNumber", "HexNumber",
		"StringType", "ListType", "DictType", "NumberType", "BoolType", "BoolLiteral",
		"Assert", "Then", "Desc", "Else", "Identifier", "IdentifierChar", "IdentifierCharStart",
		"HexDigit", "Digit", "OctalDigit", "RegexpLiteral", "RegexpLiteralChar",
		"WS",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 65, 403, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7,
		20, 2, 21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25,
		2, 26, 7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2,
		31, 7, 31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36,
		7, 36, 2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7,
		41, 2, 42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2, 46, 7, 46,
		2, 47, 7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 2, 50, 7, 50, 2, 51, 7, 51, 2,
		52, 7, 52, 2, 53, 7, 53, 2, 54, 7, 54, 2, 55, 7, 55, 2, 56, 7, 56, 2, 57,
		7, 57, 2, 58, 7, 58, 2, 59, 7, 59, 2, 60, 7, 60, 2, 61, 7, 61, 2, 62, 7,
		62, 2, 63, 7, 63, 2, 64, 7, 64, 2, 65, 7, 65, 2, 66, 7, 66, 2, 67, 7, 67,
		2, 68, 7, 68, 2, 69, 7, 69, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1,
		2, 1, 2, 1, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1,
		6, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1,
		9, 1, 10, 1, 10, 1, 10, 1, 11, 1, 11, 1, 11, 1, 12, 1, 12, 1, 12, 1, 13,
		1, 13, 1, 13, 1, 14, 1, 14, 1, 14, 1, 15, 1, 15, 1, 15, 1, 16, 1, 16, 1,
		16, 1, 17, 1, 17, 1, 17, 1, 18, 1, 18, 1, 18, 1, 19, 1, 19, 1, 19, 1, 20,
		1, 20, 1, 20, 1, 21, 1, 21, 1, 21, 1, 21, 1, 22, 1, 22, 1, 22, 1, 23, 1,
		23, 1, 23, 1, 24, 1, 24, 1, 24, 1, 24, 1, 25, 1, 25, 1, 26, 1, 26, 1, 27,
		1, 27, 1, 28, 1, 28, 1, 29, 1, 29, 1, 30, 1, 30, 1, 31, 1, 31, 1, 32, 1,
		32, 1, 33, 1, 33, 1, 34, 1, 34, 1, 35, 1, 35, 1, 36, 1, 36, 1, 37, 1, 37,
		1, 38, 1, 38, 1, 39, 1, 39, 1, 40, 1, 40, 1, 41, 1, 41, 1, 42, 1, 42, 1,
		43, 1, 43, 1, 44, 1, 44, 1, 44, 1, 45, 1, 45, 1, 46, 1, 46, 1, 46, 1, 46,
		1, 47, 4, 47, 269, 8, 47, 11, 47, 12, 47, 270, 1, 48, 1, 48, 1, 48, 1,
		48, 4, 48, 277, 8, 48, 11, 48, 12, 48, 278, 1, 49, 1, 49, 1, 49, 1, 49,
		4, 49, 285, 8, 49, 11, 49, 12, 49, 286, 1, 50, 1, 50, 1, 50, 1, 50, 4,
		50, 293, 8, 50, 11, 50, 12, 50, 294, 1, 51, 1, 51, 1, 51, 1, 51, 1, 52,
		1, 52, 1, 52, 1, 52, 1, 52, 1, 53, 1, 53, 1, 53, 1, 53, 1, 53, 1, 54, 1,
		54, 1, 54, 1, 54, 1, 54, 1, 54, 1, 54, 1, 54, 3, 54, 319, 8, 54, 1, 55,
		1, 55, 1, 55, 1, 55, 1, 55, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 1,
		56, 1, 56, 1, 56, 3, 56, 335, 8, 56, 1, 57, 1, 57, 1, 57, 1, 57, 1, 57,
		1, 57, 1, 57, 1, 58, 1, 58, 1, 58, 1, 58, 1, 58, 1, 59, 1, 59, 1, 59, 1,
		59, 1, 59, 1, 59, 1, 59, 1, 59, 3, 59, 357, 8, 59, 1, 60, 1, 60, 1, 60,
		1, 60, 1, 60, 1, 61, 1, 61, 5, 61, 366, 8, 61, 10, 61, 12, 61, 369, 9,
		61, 1, 62, 1, 62, 3, 62, 373, 8, 62, 1, 63, 3, 63, 376, 8, 63, 1, 64, 1,
		64, 1, 65, 1, 65, 1, 66, 1, 66, 1, 67, 1, 67, 4, 67, 386, 8, 67, 11, 67,
		12, 67, 387, 1, 67, 1, 67, 1, 68, 1, 68, 1, 68, 3, 68, 395, 8, 68, 1, 69,
		4, 69, 398, 8, 69, 11, 69, 12, 69, 399, 1, 69, 1, 69, 0, 0, 70, 1, 1, 3,
		2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15, 8, 17, 9, 19, 10, 21, 11, 23, 12,
		25, 13, 27, 14, 29, 15, 31, 16, 33, 17, 35, 18, 37, 19, 39, 20, 41, 21,
		43, 22, 45, 23, 47, 24, 49, 25, 51, 26, 53, 27, 55, 28, 57, 29, 59, 30,
		61, 31, 63, 32, 65, 33, 67, 34, 69, 35, 71, 36, 73, 37, 75, 38, 77, 39,
		79, 40, 81, 41, 83, 42, 85, 43, 87, 44, 89, 45, 91, 46, 93, 47, 95, 48,
		97, 49, 99, 50, 101, 51, 103, 52, 105, 53, 107, 54, 109, 55, 111, 56, 113,
		57, 115, 58, 117, 59, 119, 60, 121, 61, 123, 62, 125, 63, 127, 0, 129,
		0, 131, 0, 133, 0, 135, 64, 137, 0, 139, 65, 1, 0, 7, 3, 0, 10, 10, 13,
		13, 32, 32, 1, 0, 48, 57, 4, 0, 42, 42, 65, 90, 95, 95, 97, 122, 3, 0,
		48, 57, 65, 70, 97, 102, 1, 0, 48, 55, 1, 0, 47, 47, 3, 0, 9, 9, 13, 13,
		32, 32, 409, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7,
		1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0,
		15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1, 0, 0, 0, 0, 21, 1, 0, 0, 0,
		0, 23, 1, 0, 0, 0, 0, 25, 1, 0, 0, 0, 0, 27, 1, 0, 0, 0, 0, 29, 1, 0, 0,
		0, 0, 31, 1, 0, 0, 0, 0, 33, 1, 0, 0, 0, 0, 35, 1, 0, 0, 0, 0, 37, 1, 0,
		0, 0, 0, 39, 1, 0, 0, 0, 0, 41, 1, 0, 0, 0, 0, 43, 1, 0, 0, 0, 0, 45, 1,
		0, 0, 0, 0, 47, 1, 0, 0, 0, 0, 49, 1, 0, 0, 0, 0, 51, 1, 0, 0, 0, 0, 53,
		1, 0, 0, 0, 0, 55, 1, 0, 0, 0, 0, 57, 1, 0, 0, 0, 0, 59, 1, 0, 0, 0, 0,
		61, 1, 0, 0, 0, 0, 63, 1, 0, 0, 0, 0, 65, 1, 0, 0, 0, 0, 67, 1, 0, 0, 0,
		0, 69, 1, 0, 0, 0, 0, 71, 1, 0, 0, 0, 0, 73, 1, 0, 0, 0, 0, 75, 1, 0, 0,
		0, 0, 77, 1, 0, 0, 0, 0, 79, 1, 0, 0, 0, 0, 81, 1, 0, 0, 0, 0, 83, 1, 0,
		0, 0, 0, 85, 1, 0, 0, 0, 0, 87, 1, 0, 0, 0, 0, 89, 1, 0, 0, 0, 0, 91, 1,
		0, 0, 0, 0, 93, 1, 0, 0, 0, 0, 95, 1, 0, 0, 0, 0, 97, 1, 0, 0, 0, 0, 99,
		1, 0, 0, 0, 0, 101, 1, 0, 0, 0, 0, 103, 1, 0, 0, 0, 0, 105, 1, 0, 0, 0,
		0, 107, 1, 0, 0, 0, 0, 109, 1, 0, 0, 0, 0, 111, 1, 0, 0, 0, 0, 113, 1,
		0, 0, 0, 0, 115, 1, 0, 0, 0, 0, 117, 1, 0, 0, 0, 0, 119, 1, 0, 0, 0, 0,
		121, 1, 0, 0, 0, 0, 123, 1, 0, 0, 0, 0, 125, 1, 0, 0, 0, 0, 135, 1, 0,
		0, 0, 0, 139, 1, 0, 0, 0, 1, 141, 1, 0, 0, 0, 3, 143, 1, 0, 0, 0, 5, 146,
		1, 0, 0, 0, 7, 150, 1, 0, 0, 0, 9, 153, 1, 0, 0, 0, 11, 156, 1, 0, 0, 0,
		13, 160, 1, 0, 0, 0, 15, 164, 1, 0, 0, 0, 17, 167, 1, 0, 0, 0, 19, 170,
		1, 0, 0, 0, 21, 173, 1, 0, 0, 0, 23, 176, 1, 0, 0, 0, 25, 179, 1, 0, 0,
		0, 27, 182, 1, 0, 0, 0, 29, 185, 1, 0, 0, 0, 31, 188, 1, 0, 0, 0, 33, 191,
		1, 0, 0, 0, 35, 194, 1, 0, 0, 0, 37, 197, 1, 0, 0, 0, 39, 200, 1, 0, 0,
		0, 41, 203, 1, 0, 0, 0, 43, 206, 1, 0, 0, 0, 45, 210, 1, 0, 0, 0, 47, 213,
		1, 0, 0, 0, 49, 216, 1, 0, 0, 0, 51, 220, 1, 0, 0, 0, 53, 222, 1, 0, 0,
		0, 55, 224, 1, 0, 0, 0, 57, 226, 1, 0, 0, 0, 59, 228, 1, 0, 0, 0, 61, 230,
		1, 0, 0, 0, 63, 232, 1, 0, 0, 0, 65, 234, 1, 0, 0, 0, 67, 236, 1, 0, 0,
		0, 69, 238, 1, 0, 0, 0, 71, 240, 1, 0, 0, 0, 73, 242, 1, 0, 0, 0, 75, 244,
		1, 0, 0, 0, 77, 246, 1, 0, 0, 0, 79, 248, 1, 0, 0, 0, 81, 250, 1, 0, 0,
		0, 83, 252, 1, 0, 0, 0, 85, 254, 1, 0, 0, 0, 87, 256, 1, 0, 0, 0, 89, 258,
		1, 0, 0, 0, 91, 261, 1, 0, 0, 0, 93, 263, 1, 0, 0, 0, 95, 268, 1, 0, 0,
		0, 97, 272, 1, 0, 0, 0, 99, 280, 1, 0, 0, 0, 101, 288, 1, 0, 0, 0, 103,
		296, 1, 0, 0, 0, 105, 300, 1, 0, 0, 0, 107, 305, 1, 0, 0, 0, 109, 318,
		1, 0, 0, 0, 111, 320, 1, 0, 0, 0, 113, 334, 1, 0, 0, 0, 115, 336, 1, 0,
		0, 0, 117, 343, 1, 0, 0, 0, 119, 356, 1, 0, 0, 0, 121, 358, 1, 0, 0, 0,
		123, 363, 1, 0, 0, 0, 125, 372, 1, 0, 0, 0, 127, 375, 1, 0, 0, 0, 129,
		377, 1, 0, 0, 0, 131, 379, 1, 0, 0, 0, 133, 381, 1, 0, 0, 0, 135, 383,
		1, 0, 0, 0, 137, 394, 1, 0, 0, 0, 139, 397, 1, 0, 0, 0, 141, 142, 5, 59,
		0, 0, 142, 2, 1, 0, 0, 0, 143, 144, 5, 45, 0, 0, 144, 145, 5, 62, 0, 0,
		145, 4, 1, 0, 0, 0, 146, 147, 5, 45, 0, 0, 147, 148, 5, 45, 0, 0, 148,
		149, 5, 62, 0, 0, 149, 6, 1, 0, 0, 0, 150, 151, 5, 45, 0, 0, 151, 152,
		5, 60, 0, 0, 152, 8, 1, 0, 0, 0, 153, 154, 5, 62, 0, 0, 154, 155, 5, 45,
		0, 0, 155, 10, 1, 0, 0, 0, 156, 157, 5, 61, 0, 0, 157, 158, 5, 61, 0, 0,
		158, 159, 5, 62, 0, 0, 159, 12, 1, 0, 0, 0, 160, 161, 5, 46, 0, 0, 161,
		162, 5, 46, 0, 0, 162, 163, 5, 46, 0, 0, 163, 14, 1, 0, 0, 0, 164, 165,
		5, 37, 0, 0, 165, 166, 5, 37, 0, 0, 166, 16, 1, 0, 0, 0, 167, 168, 5, 46,
		0, 0, 168, 169, 5, 46, 0, 0, 169, 18, 1, 0, 0, 0, 170, 171, 5, 60, 0, 0,
		171, 172, 5, 61, 0, 0, 172, 20, 1, 0, 0, 0, 173, 174, 5, 62, 0, 0, 174,
		175, 5, 61, 0, 0, 175, 22, 1, 0, 0, 0, 176, 177, 5, 62, 0, 0, 177, 178,
		5, 62, 0, 0, 178, 24, 1, 0, 0, 0, 179, 180, 5, 61, 0, 0, 180, 181, 5, 62,
		0, 0, 181, 26, 1, 0, 0, 0, 182, 183, 5, 61, 0, 0, 183, 184, 5, 61, 0, 0,
		184, 28, 1, 0, 0, 0, 185, 186, 5, 61, 0, 0, 186, 187, 5, 126, 0, 0, 187,
		30, 1, 0, 0, 0, 188, 189, 5, 33, 0, 0, 189, 190, 5, 126, 0, 0, 190, 32,
		1, 0, 0, 0, 191, 192, 5, 38, 0, 0, 192, 193, 5, 38, 0, 0, 193, 34, 1, 0,
		0, 0, 194, 195, 5, 124, 0, 0, 195, 196, 5, 124, 0, 0, 196, 36, 1, 0, 0,
		0, 197, 198, 5, 33, 0, 0, 198, 199, 5, 61, 0, 0, 199, 38, 1, 0, 0, 0, 200,
		201, 5, 63, 0, 0, 201, 202, 5, 123, 0, 0, 202, 40, 1, 0, 0, 0, 203, 204,
		5, 45, 0, 0, 204, 205, 5, 123, 0, 0, 205, 42, 1, 0, 0, 0, 206, 207, 5,
		125, 0, 0, 207, 208, 5, 45, 0, 0, 208, 209, 5, 62, 0, 0, 209, 44, 1, 0,
		0, 0, 210, 211, 5, 35, 0, 0, 211, 212, 5, 123, 0, 0, 212, 46, 1, 0, 0,
		0, 213, 214, 5, 35, 0, 0, 214, 215, 5, 62, 0, 0, 215, 48, 1, 0, 0, 0, 216,
		217, 5, 35, 0, 0, 217, 218, 5, 45, 0, 0, 218, 219, 5, 62, 0, 0, 219, 50,
		1, 0, 0, 0, 220, 221, 5, 62, 0, 0, 221, 52, 1, 0, 0, 0, 222, 223, 5, 46,
		0, 0, 223, 54, 1, 0, 0, 0, 224, 225, 5, 60, 0, 0, 225, 56, 1, 0, 0, 0,
		226, 227, 5, 61, 0, 0, 227, 58, 1, 0, 0, 0, 228, 229, 5, 63, 0, 0, 229,
		60, 1, 0, 0, 0, 230, 231, 5, 40, 0, 0, 231, 62, 1, 0, 0, 0, 232, 233, 5,
		44, 0, 0, 233, 64, 1, 0, 0, 0, 234, 235, 5, 41, 0, 0, 235, 66, 1, 0, 0,
		0, 236, 237, 5, 91, 0, 0, 237, 68, 1, 0, 0, 0, 238, 239, 5, 93, 0, 0, 239,
		70, 1, 0, 0, 0, 240, 241, 5, 123, 0, 0, 241, 72, 1, 0, 0, 0, 242, 243,
		5, 125, 0, 0, 243, 74, 1, 0, 0, 0, 244, 245, 5, 35, 0, 0, 245, 76, 1, 0,
		0, 0, 246, 247, 5, 36, 0, 0, 247, 78, 1, 0, 0, 0, 248, 249, 5, 58, 0, 0,
		249, 80, 1, 0, 0, 0, 250, 251, 5, 37, 0, 0, 251, 82, 1, 0, 0, 0, 252, 253,
		5, 33, 0, 0, 253, 84, 1, 0, 0, 0, 254, 255, 5, 42, 0, 0, 255, 86, 1, 0,
		0, 0, 256, 257, 5, 45, 0, 0, 257, 88, 1, 0, 0, 0, 258, 259, 5, 97, 0, 0,
		259, 260, 5, 115, 0, 0, 260, 90, 1, 0, 0, 0, 261, 262, 5, 96, 0, 0, 262,
		92, 1, 0, 0, 0, 263, 264, 7, 0, 0, 0, 264, 265, 1, 0, 0, 0, 265, 266, 6,
		46, 0, 0, 266, 94, 1, 0, 0, 0, 267, 269, 3, 131, 65, 0, 268, 267, 1, 0,
		0, 0, 269, 270, 1, 0, 0, 0, 270, 268, 1, 0, 0, 0, 270, 271, 1, 0, 0, 0,
		271, 96, 1, 0, 0, 0, 272, 273, 5, 48, 0, 0, 273, 274, 5, 111, 0, 0, 274,
		276, 1, 0, 0, 0, 275, 277, 3, 133, 66, 0, 276, 275, 1, 0, 0, 0, 277, 278,
		1, 0, 0, 0, 278, 276, 1, 0, 0, 0, 278, 279, 1, 0, 0, 0, 279, 98, 1, 0,
		0, 0, 280, 281, 5, 48, 0, 0, 281, 282, 5, 98, 0, 0, 282, 284, 1, 0, 0,
		0, 283, 285, 2, 48, 49, 0, 284, 283, 1, 0, 0, 0, 285, 286, 1, 0, 0, 0,
		286, 284, 1, 0, 0, 0, 286, 287, 1, 0, 0, 0, 287, 100, 1, 0, 0, 0, 288,
		289, 5, 48, 0, 0, 289, 290, 5, 120, 0, 0, 290, 292, 1, 0, 0, 0, 291, 293,
		3, 129, 64, 0, 292, 291, 1, 0, 0, 0, 293, 294, 1, 0, 0, 0, 294, 292, 1,
		0, 0, 0, 294, 295, 1, 0, 0, 0, 295, 102, 1, 0, 0, 0, 296, 297, 5, 115,
		0, 0, 297, 298, 5, 116, 0, 0, 298, 299, 5, 114, 0, 0, 299, 104, 1, 0, 0,
		0, 300, 301, 5, 108, 0, 0, 301, 302, 5, 105, 0, 0, 302, 303, 5, 115, 0,
		0, 303, 304, 5, 116, 0, 0, 304, 106, 1, 0, 0, 0, 305, 306, 5, 100, 0, 0,
		306, 307, 5, 105, 0, 0, 307, 308, 5, 99, 0, 0, 308, 309, 5, 116, 0, 0,
		309, 108, 1, 0, 0, 0, 310, 311, 5, 105, 0, 0, 311, 312, 5, 110, 0, 0, 312,
		319, 5, 116, 0, 0, 313, 314, 5, 102, 0, 0, 314, 315, 5, 108, 0, 0, 315,
		316, 5, 111, 0, 0, 316, 317, 5, 97, 0, 0, 317, 319, 5, 116, 0, 0, 318,
		310, 1, 0, 0, 0, 318, 313, 1, 0, 0, 0, 319, 110, 1, 0, 0, 0, 320, 321,
		5, 98, 0, 0, 321, 322, 5, 111, 0, 0, 322, 323, 5, 111, 0, 0, 323, 324,
		5, 108, 0, 0, 324, 112, 1, 0, 0, 0, 325, 326, 5, 116, 0, 0, 326, 327, 5,
		114, 0, 0, 327, 328, 5, 117, 0, 0, 328, 335, 5, 101, 0, 0, 329, 330, 5,
		102, 0, 0, 330, 331, 5, 97, 0, 0, 331, 332, 5, 108, 0, 0, 332, 333, 5,
		115, 0, 0, 333, 335, 5, 101, 0, 0, 334, 325, 1, 0, 0, 0, 334, 329, 1, 0,
		0, 0, 335, 114, 1, 0, 0, 0, 336, 337, 5, 97, 0, 0, 337, 338, 5, 115, 0,
		0, 338, 339, 5, 115, 0, 0, 339, 340, 5, 101, 0, 0, 340, 341, 5, 114, 0,
		0, 341, 342, 5, 116, 0, 0, 342, 116, 1, 0, 0, 0, 343, 344, 5, 116, 0, 0,
		344, 345, 5, 104, 0, 0, 345, 346, 5, 101, 0, 0, 346, 347, 5, 110, 0, 0,
		347, 118, 1, 0, 0, 0, 348, 349, 5, 100, 0, 0, 349, 350, 5, 101, 0, 0, 350,
		351, 5, 115, 0, 0, 351, 357, 5, 99, 0, 0, 352, 353, 5, 110, 0, 0, 353,
		354, 5, 111, 0, 0, 354, 355, 5, 116, 0, 0, 355, 357, 5, 101, 0, 0, 356,
		348, 1, 0, 0, 0, 356, 352, 1, 0, 0, 0, 357, 120, 1, 0, 0, 0, 358, 359,
		5, 101, 0, 0, 359, 360, 5, 108, 0, 0, 360, 361, 5, 115, 0, 0, 361, 362,
		5, 101, 0, 0, 362, 122, 1, 0, 0, 0, 363, 367, 3, 127, 63, 0, 364, 366,
		3, 125, 62, 0, 365, 364, 1, 0, 0, 0, 366, 369, 1, 0, 0, 0, 367, 365, 1,
		0, 0, 0, 367, 368, 1, 0, 0, 0, 368, 124, 1, 0, 0, 0, 369, 367, 1, 0, 0,
		0, 370, 373, 7, 1, 0, 0, 371, 373, 3, 127, 63, 0, 372, 370, 1, 0, 0, 0,
		372, 371, 1, 0, 0, 0, 373, 126, 1, 0, 0, 0, 374, 376, 7, 2, 0, 0, 375,
		374, 1, 0, 0, 0, 376, 128, 1, 0, 0, 0, 377, 378, 7, 3, 0, 0, 378, 130,
		1, 0, 0, 0, 379, 380, 7, 1, 0, 0, 380, 132, 1, 0, 0, 0, 381, 382, 7, 4,
		0, 0, 382, 134, 1, 0, 0, 0, 383, 385, 5, 47, 0, 0, 384, 386, 3, 137, 68,
		0, 385, 384, 1, 0, 0, 0, 386, 387, 1, 0, 0, 0, 387, 385, 1, 0, 0, 0, 387,
		388, 1, 0, 0, 0, 388, 389, 1, 0, 0, 0, 389, 390, 5, 47, 0, 0, 390, 136,
		1, 0, 0, 0, 391, 392, 5, 92, 0, 0, 392, 395, 5, 47, 0, 0, 393, 395, 8,
		5, 0, 0, 394, 391, 1, 0, 0, 0, 394, 393, 1, 0, 0, 0, 395, 138, 1, 0, 0,
		0, 396, 398, 7, 6, 0, 0, 397, 396, 1, 0, 0, 0, 398, 399, 1, 0, 0, 0, 399,
		397, 1, 0, 0, 0, 399, 400, 1, 0, 0, 0, 400, 401, 1, 0, 0, 0, 401, 402,
		6, 69, 0, 0, 402, 140, 1, 0, 0, 0, 14, 0, 270, 278, 286, 294, 318, 334,
		356, 367, 372, 375, 387, 394, 399, 1, 6, 0, 0,
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
	SyntaxFlowLexerT__3            = 4
	SyntaxFlowLexerT__4            = 5
	SyntaxFlowLexerDeepFilter      = 6
	SyntaxFlowLexerDeep            = 7
	SyntaxFlowLexerPercent         = 8
	SyntaxFlowLexerDeepDot         = 9
	SyntaxFlowLexerLtEq            = 10
	SyntaxFlowLexerGtEq            = 11
	SyntaxFlowLexerDoubleGt        = 12
	SyntaxFlowLexerFilter          = 13
	SyntaxFlowLexerEqEq            = 14
	SyntaxFlowLexerRegexpMatch     = 15
	SyntaxFlowLexerNotRegexpMatch  = 16
	SyntaxFlowLexerAnd             = 17
	SyntaxFlowLexerOr              = 18
	SyntaxFlowLexerNotEq           = 19
	SyntaxFlowLexerConditionStart  = 20
	SyntaxFlowLexerDeepNextStart   = 21
	SyntaxFlowLexerDeepNextEnd     = 22
	SyntaxFlowLexerTopDefStart     = 23
	SyntaxFlowLexerDefStart        = 24
	SyntaxFlowLexerTopDef          = 25
	SyntaxFlowLexerGt              = 26
	SyntaxFlowLexerDot             = 27
	SyntaxFlowLexerLt              = 28
	SyntaxFlowLexerEq              = 29
	SyntaxFlowLexerQuestion        = 30
	SyntaxFlowLexerOpenParen       = 31
	SyntaxFlowLexerComma           = 32
	SyntaxFlowLexerCloseParen      = 33
	SyntaxFlowLexerListSelectOpen  = 34
	SyntaxFlowLexerListSelectClose = 35
	SyntaxFlowLexerMapBuilderOpen  = 36
	SyntaxFlowLexerMapBuilderClose = 37
	SyntaxFlowLexerListStart       = 38
	SyntaxFlowLexerDollarOutput    = 39
	SyntaxFlowLexerColon           = 40
	SyntaxFlowLexerSearch          = 41
	SyntaxFlowLexerBang            = 42
	SyntaxFlowLexerStar            = 43
	SyntaxFlowLexerMinus           = 44
	SyntaxFlowLexerAs              = 45
	SyntaxFlowLexerBacktick        = 46
	SyntaxFlowLexerWhiteSpace      = 47
	SyntaxFlowLexerNumber          = 48
	SyntaxFlowLexerOctalNumber     = 49
	SyntaxFlowLexerBinaryNumber    = 50
	SyntaxFlowLexerHexNumber       = 51
	SyntaxFlowLexerStringType      = 52
	SyntaxFlowLexerListType        = 53
	SyntaxFlowLexerDictType        = 54
	SyntaxFlowLexerNumberType      = 55
	SyntaxFlowLexerBoolType        = 56
	SyntaxFlowLexerBoolLiteral     = 57
	SyntaxFlowLexerAssert          = 58
	SyntaxFlowLexerThen            = 59
	SyntaxFlowLexerDesc            = 60
	SyntaxFlowLexerElse            = 61
	SyntaxFlowLexerIdentifier      = 62
	SyntaxFlowLexerIdentifierChar  = 63
	SyntaxFlowLexerRegexpLiteral   = 64
	SyntaxFlowLexerWS              = 65
)
