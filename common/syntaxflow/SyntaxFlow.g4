grammar SyntaxFlow;

/*
yaklang.io
v1ll4n.a5k@gmail.com

build with `antlr -Dlanguage=Go ./SyntaxFlow.g4 -o sf -package sf -no-listener -visitor`

SyntaxFlow is a search expr can handle some structured data
*/

flow: statements EOF;

statements: statement+;

statement
    : filterStatement eos?              # Filter
    | checkStatement eos?               # Check
    | descriptionStatement eos?         # Description
    | eos                               # Empty
    ;
filterStatement
    : refVariable filterItem*  (As refVariable)? # RefFilterExpr
    | filterExpr  (As refVariable)?              # PureFilterExpr
    ;


// eos means end of statement
// the ';' can be
eos: ';' | line ;
line: '\n';

// descriptionStatement will describe the filterExpr with stringLiteral
descriptionStatement: Desc ('(' descriptionItems? ')') | ('{' descriptionItems? '}');
descriptionItems: descriptionItem (',' descriptionItem)*;
descriptionItem
    : stringLiteral
    | stringLiteral ':' stringLiteral
    ;

// checkStatement will check the filterExpr($params) is true( .len > 0), if not,
// it will record an error with stringLiteral
// if thenExpr is provided, it will be executed(description) after the assertStatement
checkStatement: Check refVariable thenExpr? elseExpr?;
thenExpr: Then stringLiteral;
elseExpr: Else stringLiteral;

refVariable
    :  '$' (identifier | ('(' identifier ')'));


filterItemFirst
    : nameFilter                                 # NamedFilter
    | '.' nameFilter                             # FieldCallFilter
    ;
filterItem
    : filterItemFirst                            #First
    | '(' actualParam? ')'                       # FunctionCallFilter
    | '[' sliceCallItem ']'                      # FieldIndexFilter
    | '?{' conditionExpression '}'               # OptionalFilter
    | '->'                                       # NextFilter
    | '#>'                                       # DefFilter
    | '-->'                                      # DeepNextFilter
    | '-{' (recursiveConfig)? '}->'              # DeepNextConfigFilter
    | '#->'                                      # TopDefFilter
    | '#{' (recursiveConfig)? '}->'              # TopDefConfigFilter
    | '-<' useDefCalcDescription '>-'            # UseDefCalcFilter
    ;

filterExpr: filterItemFirst filterItem* ;

useDefCalcDescription
    : identifier useDefCalcParams?
    ;

useDefCalcParams
    : '{' recursiveConfig? '}'
    | '(' recursiveConfig? ')'
    ;

actualParam
    : singleParam                      # AllParam
    | actualParamFilter+ singleParam?  # EveryParam
    ;

actualParamFilter: singleParam ',' | ',';

singleParam: ( '#>' | '#{' (recursiveConfig)? '}' )? filterStatement ;

recursiveConfig: recursiveConfigItem (',' recursiveConfigItem)* ','? line?;
recursiveConfigItem: line? identifier ':' recursiveConfigItemValue ;
recursiveConfigItemValue
    : (identifier | numberLiteral)
    | '`' filterStatement '`'
    ;

sliceCallItem: nameFilter | numberLiteral;

nameFilter: '*' | identifier | regexpLiteral;

chainFilter
    : '[' ((statements (',' statements)*) | '...') ']'          # Flat
    | '{' ((identifier ':') statements (';' (identifier ':') statements )*)? ';'? '}'  # BuildMap
    ;

stringLiteralWithoutStarGroup: stringLiteralWithoutStar (',' stringLiteralWithoutStar)* ','?;
negativeCondition: (Not | '!');

conditionExpression
    : '(' conditionExpression ')'                                                # ParenCondition
    | filterExpr                                                                 # FilterCondition        // filter dot(.)Member and fields
    |  Opcode ':' opcodes (',' opcodes) * ','?                      # OpcodeTypeCondition    // something like .(call, phi)
    |  Have  ':' stringLiteralWithoutStarGroup       # StringContainHaveCondition // something like .(have: 'a', 'b')
    |  HaveAny ':' stringLiteralWithoutStarGroup       # StringContainAnyCondition // something like .(have: 'a', 'b')
    | negativeCondition conditionExpression                                                    # NotCondition
    | op = (
        '>' | '<' | '=' | '==' | '>='
         | '<=' | '!='
        ) (
            numberLiteral | identifier | boolLiteral
        ) # FilterExpressionCompare
    | op = ( '=~' | '!~') (stringLiteral | regexpLiteral) # FilterExpressionRegexpMatch
    | conditionExpression '&&' conditionExpression      # FilterExpressionAnd
    | conditionExpression '||' conditionExpression      # FilterExpressionOr
    ;

numberLiteral: Number | OctalNumber | BinaryNumber | HexNumber;
stringLiteral: identifier | '*';
stringLiteralWithoutStar: identifier | regexpLiteral ;
regexpLiteral: RegexpLiteral;
identifier: Identifier | keywords | QuotedStringLiteral;

keywords
    : types
    | opcodes
    | Opcode
    // | As
    // | Check
    | Then
    | Desc
    | Else
    | Type
    | In
    | Have
    | HaveAny
    ;

opcodes: Call | Constant | Phi | FormalParam | Return;

types: StringType | NumberType | ListType | DictType | BoolType;
boolLiteral: BoolLiteral;

DeepFilter: '==>';
Deep: '...';

Percent: '%%';
DeepDot: '..';
LtEq: '<=';
GtEq: '>=';
DoubleGt: '>>';
Filter: '=>';
EqEq: '==';
RegexpMatch: '=~';
NotRegexpMatch: '!~';
And: '&&';
Or: '||';
NotEq: '!=';
ConditionStart: '?{';
DeepNextStart: '-{';
DeepNextEnd: '}->';
TopDefStart: '#{';
DefStart: '#>';
TopDef: '#->';
Gt: '>';
Dot: '.';
Lt: '<';
Eq: '=';
Question: '?';
OpenParen: '(';
Comma: ',';
CloseParen: ')';
ListSelectOpen: '[';
ListSelectClose: ']';
MapBuilderOpen: '{';
MapBuilderClose: '}';
ListStart: '#';
DollarOutput: '$';
Colon: ':';
Search: '%';
Bang: '!';
Star: '*';
Minus: '-';
As: 'as';
Backtick: '`';
SingleQuote: '\'';
DoubleQuote: '"';

WhiteSpace: [ \r\n] -> skip;
Number: Digit+;
OctalNumber: '0o' OctalDigit+;
BinaryNumber: '0b' ('0' | '1')+;
HexNumber: '0x' HexDigit+;
//StringLiteral: '`' (~[`])* '`';
StringType: 'str';
ListType: 'list';
DictType: 'dict';
NumberType: 'int' | 'float';
BoolType: 'bool';
BoolLiteral: 'true' | 'false';
Check: 'check';
Then: 'then';
Desc: 'desc' | 'note';
Else: 'else';
Type: 'type';
In: 'in';
Call: 'call';
Constant: 'const' | 'constant';
Phi: 'phi';
FormalParam: 'param' | 'formal_param';
Return: 'return' | 'ret';
Opcode: 'opcode';
Have: 'have';
HaveAny: 'any';
Not: 'not';

Identifier: IdentifierCharStart IdentifierChar*;
IdentifierChar: [0-9] | IdentifierCharStart;

QuotedStringLiteral
    : SingleQuote ( ~['\\\r\n] | ('\\\'') | '\\\\' | '\\')* SingleQuote
    | DoubleQuote ( ~["\\\r\n] | '\\"' | '\\\\' | '\\' )* DoubleQuote
    ;

fragment IdentifierCharStart: '*' | '_' | [a-z] | [A-Z];
fragment HexDigit: [a-fA-F0-9];
fragment Digit: [0-9];
fragment OctalDigit: [0-7];

RegexpLiteral: '/' RegexpLiteralChar+ '/';
fragment RegexpLiteralChar
    : '\\' '/'
    | ~[/]
    ;

WS: [ \t\r]+ -> skip;