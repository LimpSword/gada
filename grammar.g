%token in out true false null
%% /* LL(1) */

chiffre : '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' ;

chiffre_plus
    : chiffre chiffre_plus
    | chiffre ;

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
    : ident ',' ident_plus_comma
    | ident ;

ident_tail
    : alpha | chiffre | '_' ;

ident_tail_star
    : ident_tail ident_tail_star
    | ident_tail
    | /* eps */ ;

entier
    : chiffre_plus ;

caractere
    : ''' printable_caractere ''' ;

printable_caractere
    : ' ' | '!' | '"' | '#' | '$' | '%' | '&' | ''' | '(' | ')' | '*' | '+' | ',' | '-' | '.' | '/' 
    | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | ':' | ';' | '=' | '?' | '@' 
    | alpha | '[' | '\' | ']' | '^' | '_' | '`' | '{' | '|' | '}' | '~' ;

fichier
    : 'with' 'Ada.Text_IO;' 'use' 'Ada.Text_IO;'
      'procedure' ident 'is' decl_star
      'begin' instr_plus 'end' ident_opt ';' 'EOF' ;

decl
    : 'type' ident ';'
    | 'type' ident 'is' 'access' ident ';'
    | 'type' ident 'is' 'record' champs_plus 'end' 'record' ';'
    | ident_plus_comma ':' 'type' init ';'
    | 'procedure' ident params_opt 'is' decl_star
      'begin' instr_plus 'end' ident_opt ';'
    | 'function' ident params_opt 'return' 'type' 'is' ;

decl_star
    : 'begin' instr_plus 'end' ident_opt ';' ;
    
init
    : ':=' expr
    | /*eps*/ ;

decl_star
    : decl decl_star
    | decl
    | /*eps*/ ;

champs
    : ident_plus_comma ':' 'type' ';' ;

champs_plus
    : champs champs_plus
    | champs ;

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
    : param ';' param_plus_semicolon
    | param ;

mode
    : in
    | in out ;

mode_opt
    : mode
    | /*eps*/ ;

expr
    : entier expr_operator
    | caractere expr_operator
    | true expr_operator
    | false expr_operator
    | null expr_operator
    | '(' expr ')' expr_operator
    | access expr_operator
    | 'not' expr expr_operator
    | '-' expr expr_operator
    | 'new' ident expr_operator
    | ident '(' expr_plus_comma ')' expr_operator
    | 'character' ''' 'val' '(' expr ')' expr_operator;

expr_operator
    : operateur expr expr_operator ;

expr_plus_comma
    : expr ',' expr_plus_comma
    | expr ;

expr_opt
    : expr
    | /*eps*/ ;

instr
    : 'access' ':=' expr ';'
    | ident ';'
    | ident '(' expr_plus_comma ')' ';'
    | 'return' expr_opt ';'
    | 'begin' instr_plus 'end' ';'
    | 'if' expr 'then' instr_plus else_if_star
      else_instr_opt 'end' 'if' ';'
    | 'for' ident 'in' reverse_instr expr '..' expr
      'loop' instr_plus 'end' 'loop' ';'
    | 'while' expr 'loop' instr_plus 'end' 'loop' ;

instr_plus
    : instr instr_plus
    | instr ;

else_if
    : 'elsif' expr 'then' instr_plus ;

else_if_star
    : else_if else_if_star
    | else_if
    | /*eps*/ ;

else_instr
    : 'else' instr_plus ;

else_instr_opt
    : 'else_instr'
    | /*eps*/ ;

operateur
    : '='
    | '/='
    | '>'
    | '>='
    | '<'
    | '<='
    | '+'
    | '-'
    | '*'
    | '/'
    | 'rem'
    | 'and'
    | 'and' 'then'
    | 'or'
    | 'or' 'else' ;

access
    : ident
    | ident '.' expr ; /* A vérifier */

reverse_instr
    : 'reverse'
    | /*eps*/ ;
