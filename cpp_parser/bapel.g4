grammar bapel;

// Entry points
sourceFile: moduleHeader importsSection? implsSection? flagsSection? sources? EOF # baseSourceFile
          | implementsHeader importsSection? implsSection? flagsSection? sources? EOF # implSourceFile
          ;

moduleHeader: MODULE moduleID;
implementsHeader: IMPLEMENTS moduleID;

workspace: WORKSPACE LBRACE packagesSection RBRACE;
packagesSection: PACKAGES LBRACE packageRule+ RBRACE;
packageRule: PREFIX moduleID IN filename
           | MODULE moduleID IN filename
           ;

// Shared sections
importsSection: IMPORTS LBRACE moduleID+ RBRACE;
implsSection: IMPLS LBRACE filename+ RBRACE;
flagsSection: FLAGS LBRACE filename+ RBRACE;

moduleID: IDENTIFIER (DOT IDENTIFIER)*;
filename: STRING_LITERAL;

// Source
sources: source+;
source: declNoExport
      | functionNoExport
      | PUB functionNoExport
      | traitDecl
      | PUB traitDecl
      | implBlock
      ;

traitDecl: TRAIT id LBRACE traitMethod* RBRACE;
traitMethod: FN id functionArgs ARROW type_;

implBlock: IMPL id FOR type_ LBRACE functionNoExport* RBRACE # traitImpl
         | IMPL type_ LBRACE functionNoExport* RBRACE            # inherentImpl
         ;


declNoExport: declNoTerm
            | DECL termDecl
            | PUB termDecl
            ;

functionNoExport: FN id typeAbstraction? functionArgs ARROW type_ blockExpr;

functionArgs: LPAREN (arg (COMMA arg)*)? RPAREN;
arg: IDENTIFIER COLON type_;

// Decl (for annotations)
decl: PUB unexportedDecl
    | unexportedDecl
    ;

unexportedDecl: termDecl
              | typeDecl
              ;

declNoTerm: PUB typeDecl
          | typeDecl
          ;

termDecl: id COLON type_;

typeDecl: TYPE id typeAbstraction? ASSIGN type_
        | TYPE id typeAbstraction?
        ;

typeAbstraction: LBRACKET tvar (COMMA tvar)* RBRACKET;
tvar: SINGLE_QUOTE IDENTIFIER;

// Types
type_: forallType;

forallType: FORALL typeAbstraction functionType
          | functionType
          ;

functionType: ptrType (ARROW functionType)?
            ;

ptrType: AMP ptrType
       | appType
       ;

appType: appType primaryType
       | primaryType
       ;

primaryType: arrayType
           | structType
           | tupleType
           | variantType
           | SINGLE_QUOTE IDENTIFIER
           | id
           | LPAREN type_ RPAREN
           ;

arrayType: LBRACKET type_ COMMA INT_LITERAL RBRACKET
         | LBRACKET type_ RBRACKET
         | LBRACKET type_ COMMA MINUS INT_LITERAL RBRACKET
         ;

structType: STRUCT LBRACE (fields (COMMA?) )? RBRACE;
fields: field (COMMA field)*;
field: id COLON type_;

tupleType: LPAREN RPAREN
         | LPAREN tupleTypeArgs RPAREN
         ;
tupleTypeArgs: type_ (COMMA type_)+;

variantType: VARIANT LBRACE (tags (COMMA?) )? RBRACE;
tags: tag (COMMA tag)*;
tag: id type_;

// Expressions
expression: expressionWithoutBlock
          | expressionWithBlock
          ;

expressionWithoutBlock: assignTerm
                      | operatorExpr
                      | returnTerm
                      ;

expressionWithBlock: blockExpr
                   | ifTerm
                   | forTerm
                   | lambdaTerm
                   | matchTerm
                   | setTerm
                   ;

assignTerm: (id | tupleExpr) LARROW expression;

returnTerm: RETURN expressionWithoutBlock;

ifTerm: IF expressionWithoutBlock blockExpr (ELSE (blockExpr | ifTerm))?;

forTerm: FOR expressionWithoutBlock blockExpr;

lambdaTerm: FN typeAbstraction? functionArgs blockExpr;

matchTerm: MATCH expression LBRACE matchArms (COMMA?) RBRACE;
matchArms: matchArm (COMMA matchArm)*;
matchArm: id IDENTIFIER FAT_ARROW expression;

setTerm: SET expression LBRACE labelValues (COMMA?) RBRACE;

blockExpr: LBRACE blockStatements RBRACE;
blockStatements: statements expressionWithoutBlock?
               | expressionWithoutBlock
               ;
statements: statement+;
statement: letStatement
         | expressionStatement
         ;

letStatement: LET id COLON type_ ASSIGN expression SEMI
            | LET id ASSIGN expression SEMI
            ;

expressionStatement: expressionWithoutBlock SEMI
                   | expressionWithBlock SEMI?
                   ;

// Operators
operatorExpr: logicalOrExpr;

logicalOrExpr: logicalOrExpr OR logicalAndExpr
             | logicalAndExpr
             ;

logicalAndExpr: logicalAndExpr AND equalityExpr
              | equalityExpr
              ;

equalityExpr: equalityExpr (NE | EQ) typeApplicativeArgs? comparisonExpr
            | comparisonExpr
            ;

comparisonExpr: comparisonExpr (GT | GE | LT | LE) typeApplicativeArgs? additiveExpr
              | additiveExpr
              ;

additiveExpr: additiveExpr (PLUS | MINUS) typeApplicativeArgs? multiplicativeExpr
            | multiplicativeExpr
            ;

multiplicativeExpr: multiplicativeExpr (MUL | DIV) typeApplicativeArgs? unaryExpr
                  | unaryExpr
                  ;

unaryExpr: (NOT | MINUS) typeApplicativeArgs? unaryExpr
         | applicativeExpr
         ;

applicativeExpr: applicativeExpr basePrimaryExpr
               | typeApplicativeExpr
               ;

typeApplicativeExpr: primaryExpr typeApplicativeArgs
                   | primaryExpr
                   ;

typeApplicativeArgs: LBRACKET tupleTypeArgs RBRACKET
                   | LBRACKET type_ RBRACKET
                   ;

basePrimaryExpr: AMP id
               | projectionExpr
               | INT_LITERAL
               | FLOAT_LITERAL
               ;

primaryExpr: MUL primaryExpr
           | basePrimaryExpr
           ;

projectionExpr: projectionExpr DOT INT_LITERAL
              | projectionExpr DOT IDENTIFIER
              | derefExpr
              ;

derefExpr: injectionExpr
         | RUNE_LITERAL
         | STRING_LITERAL
         | structExpr
         | tupleExpr
         | varExpr
         | LPAREN expression RPAREN
         ;

injectionExpr: VARIANT LBRACE type_ labelValue RBRACE;

structExpr: STRUCT LBRACE (labelValues (COMMA?))? RBRACE;

labelValues: labelValue (COMMA labelValue)*;
labelValue: id ASSIGN expression
          | INT_LITERAL ASSIGN expression
          ;

tupleExpr: LPAREN RPAREN
         | LPAREN tupleExprArgs RPAREN
         ;
tupleExprArgs: expression (COMMA expression)+;

varExpr: id;

id: idTokens
  | LPAREN ( OR | AND | NE | EQ | GT | GE | LT | LE | PLUS | MINUS | MUL | DIV | NOT ) RPAREN
  ;

idTokens: IDENTIFIER (DOUBLE_COLON (IDENTIFIER | SET))*;

// Lexer rules
WORKSPACE: 'workspace';
PACKAGES: 'packages';
PREFIX: 'prefix';
MODULE: 'module';
IMPLEMENTS: 'implements';
IMPORTS: 'imports';
IMPLS: 'impls';
FLAGS: 'flags';
IN: 'in';
PUB: 'pub';
DECL: 'decl';
FN: 'fn';
TYPE: 'type';
FORALL: 'forall';
STRUCT: 'struct';
VARIANT: 'variant';
MATCH: 'match';
SET: 'set';
LET: 'let';
RETURN: 'return';
IF: 'if';
ELSE: 'else';
FOR: 'for';
TRAIT: 'trait';
IMPL: 'impl';


DOUBLE_COLON: '::';
ARROW: '->';
FAT_ARROW: '=>';
LARROW: '<-';
OR: '||';
AND: '&&';
NE: '!=';
EQ: '==';
GE: '>=';
LE: '<=';
GT: '>';
LT: '<';
PLUS: '+';
MINUS: '-';
MUL: '*';
DIV: '/';
NOT: '!';
AMP: '&';
DOT: '.';
ASSIGN: '=';
COMMA: ',';
SEMI: ';';
COLON: ':';
LBRACE: '{';
RBRACE: '}';
LPAREN: '(';
RPAREN: ')';
LBRACKET: '[';
RBRACKET: ']';
SINGLE_QUOTE: '\'';

IDENTIFIER: [a-zA-Z_][a-zA-Z0-9_]*;
INT_LITERAL: [0-9]+ | '0x' [0-9a-fA-F]+;
FLOAT_LITERAL: [0-9]+ '.' [0-9]+;
RUNE_LITERAL: '\'' ( ~['\\\r\n] | '\\' . ) '\'';
STRING_LITERAL: '"' ( ~["\\\r\n] | '\\' . )* '"';
RAW_STRING_LITERAL: '`' ~[`]* '`';
UNTERMINATED_RUNE_LITERAL: '\'' '\\' ( ~['\\] | '\\' . )*
                         | '\'' ~[a-zA-Z_'\\] ( ~['\\] | '\\' . )*
                         ;
UNTERMINATED_STRING_LITERAL: '"' ( ~["\\] | '\\' . )*;
UNTERMINATED_RAW_STRING_LITERAL: '`' ~[`]*;
UNTERMINATED_BLOCK_COMMENT: '/*' ( ~[*] | '*'+ ~[*/] )* '*'* ;


WS: [ \t\r\n]+ -> skip;
LINE_COMMENT: '//' ~[\r\n]* -> skip;
BLOCK_COMMENT: '/*' .*? '*/' -> skip;
