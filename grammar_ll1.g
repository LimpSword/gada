%token in out true false null
%% /* LL(1) */

fichier
    : 'with' 'Ada.Text_IO;' 'use' 'Ada.Text_IO;'
      'procedure' ident 'is' decl_star
      'begin' instr_plus 'end' ident_opt ';' 'EOF' ;

decl
    : 'type' ident decl2
    | ident_plus_comma ':' 'type' init ';'
    | 'procedure' ident params_opt 'is' decl_star 'begin' instr_plus 'end' ident_opt ';'
    | 'function' ident params_opt 'return' 'type' 'is' decl_star 'begin' instr_plus 'end' ident_opt ';' ;

decl2
    : ';'
    | 'is' decl3 ;

decl3
    : 'access' ident ';'
    | 'record' champs_plus 'end' 'record' ';' ;

init
    : ':=' expr
    | /*eps*/ ;

decl_star
    : decl decl_star
    | /*eps*/ ;

champs
    : ident_plus_comma ':' 'type' ';' ;

champs_plus
    : champs champs_plus2 ;

champs_plus2
    : champs champs_plus2
    | /*eps*/ ;

type
    : ident
    | 'access' ident ;

params
    : '(' param_plus_semicolon ')' ';' ;

params_opt
    : params
    | /*eps*/ ;

param
    : ident_plus_comma ':' mode_opt type ;

param_plus_semicolon
    : param param_plus_semicolon2 ;

param_plus_semicolon2
    : ';' param param_plus_semicolon2
    | /*eps*/ ;

mode
    : in mode2 ;

mode2
    : out
    | /*eps*/ ;

mode_opt
    : mode
    | /*eps*/ ;

expr
    : or_expr ;

or_expr
    : and_expr or_expr_tail
    ;

or_expr_tail
    : 'or' and_expr or_expr_tail
    | /*eps*/ ;

and_expr
    : equality_expr and_expr_tail
    ;

and_expr_tail
    : 'and' equality_expr and_expr_tail
    | /*eps*/ ;

equality_expr
    : relational_expr equality_expr_tail
    ;

equality_expr_tail
    : '=' relational_expr equality_expr_tail
    | '/=' relational_expr equality_expr_tail
    | /*eps*/ ;

relational_expr
    : additive_expr relational_expr_tail
    ;

relational_expr_tail
    : '<' additive_expr relational_expr_tail
    | '<=' additive_expr relational_expr_tail
    | '>' additive_expr relational_expr_tail
    | '>=' additive_expr relational_expr_tail
    | /*eps*/ ;

additive_expr
    : multiplicative_expr additive_expr_tail
    ;

additive_expr_tail
    : '+' multiplicative_expr additive_expr_tail
    | '-' multiplicative_expr additive_expr_tail
    | /*eps*/ ;

multiplicative_expr
    : unary_expr multiplicative_expr_tail
    ;

multiplicative_expr_tail
    : '*' unary_expr multiplicative_expr_tail
    | '/' unary_expr multiplicative_expr_tail
    | 'rem' unary_expr multiplicative_expr_tail
    | /*eps*/ ;

unary_expr
    : '-' unary_expr
    | primary_expr ;

primary_expr
    : entier
    | caractere
    | 'true'
    | 'false'
    | 'null'
    | '(' expr ')'
    | 'not'
    | 'new' ident
    | ident primary_expr2
    | 'character' ''' 'val' '(' expr ')' ;

primary_expr2
    : access2
    | '(' expr_plus_comma ')' primary_expr3 ;

primary_expr3
    : '.' ident access2
    | /*eps*/ ;

access2
    : '.' ident access2
    | /*eps*/ ;

expr_plus_comma
    : expr expr_plus_comma2 ;

expr_plus_comma2
    : ',' expr expr_plus_comma2
    | /*eps*/ ;

expr_opt
    : expr
    | /*eps*/ ;

instr
    : 'access' ':=' expr ';'
    | ident instr2
    | 'return' expr_opt ';'
    | 'begin' instr_plus 'end' ';'
    | 'if' expr 'then' instr_plus else_if_star
      else_instr_opt 'end' 'if' ';'
    | 'for' ident 'in' reverse_instr expr '..' expr
      'loop' instr_plus 'end' 'loop' ';'
    | 'while' expr 'loop' instr_plus 'end' 'loop' ;

instr2
    : ';'
    | '(' expr_plus_comma ')' ';' ;

instr_plus
    : instr instr_plus2 ;

instr_plus2
    : instr instr_plus2
    | /*eps*/ ;

else_if
    : 'elsif' expr 'then' instr_plus ;

else_if_star
    : else_if else_if_star
    | /*eps*/ ;

else_instr
    : 'else' instr_plus ;

else_instr_opt
    : 'else_instr'
    | /*eps*/ ;

reverse_instr
    : 'reverse'
    | /*eps*/ ;

chiffre : '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' ;

chiffre_plus
    : chiffre chiffre_plus2 ;

chiffre_plus2
    : chiffre chiffre_plus2
    | /*eps*/ ;

alpha
    : 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q'
    | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z'
    | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q'
    | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' ;

ident
    : alpha ident_tail_star ;

ident_opt
    : ident
    | /*eps*/ ;

ident_plus_comma
    : ident ident_plus_comma2 ;

ident_plus_comma2
    : ',' ident ident_plus_comma2
    | /*eps*/ ;

ident_tail
    : alpha | chiffre | '_' ;

ident_tail_star
    : ident_tail ident_tail_star
    | /* eps */ ;

entier
    : chiffre_plus ;

caractere
    : ''' printable_caractere ''' ;

printable_caractere
    : ' ' | '!' | '"' | '#' | '$' | '%' | '&' | ''' | '(' | ')' | '*' | '+' | ',' | '-' | '.' | '/'
    | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | ';' | '=' | '?' | '@'
    | alpha | '[' | '\' | ']' | '^' | '_' | '`' | '{' | '|' | '}' | '~' ;
