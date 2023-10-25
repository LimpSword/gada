grammar gram_c;

<chiffre>
    : 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9

<chiffre+>
    : <chiffre><chiffre+>
    | <chiffre>

<alpha>
    : a | b | c | d | e | f | g | h | i | j | k | l | m | n | o | p | q
    | r | s | t | u | v | w | x | y | z
    | A | B | C | D | E | F | G | H | I | J | K | L | M | N | O | P | Q
    | R | S | T | U | V | W | X | Y | Z 

<ident>
    : <alpha> <ident_tail*>

<ident?>
    : <ident>
    | ^

<ident+_,>
    : <ident>,<ident+_,>
    | <ident>

<ident_tail>
    : <alpha> | <chiffre> | _

<ident_tail*>
    : <ident_tail><ident_tail*>
    | <ident_tail>
    | ^

<entier>
    : <chiffre+>

<caractere>
    : ' <printable_caractere> '

<printable_caractere>
    : ' ' | ! | " | # | $ | % | & | ' | ( | ) | * | + | , | - | . | / 
    | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | : | ; | < | = | > | ? | @ 
    | A | B | C | D | E | F | G | H | I | J | K | L | M | N | O | P | Q 
    | R | S | T | U | V | W | X | Y | Z | [ | \ | ] | ^ | _ | ` 
    | a | b | c | d | e | f | g | h | i | j | k | l | m | n | o | p | q 
    | r | s | t | u | v | w | x | y | z | { | | | } | ~

<fichier>
    : with Ada.Text_IO; use Ada.Text_IO;
      procedure <ident> is <decl*>
      begin <instr+> end <ident?> ; EOF

<decl>
    : type <ident> ;
    | type <ident> is access <ident> ;
    | type <ident> is record <champs+> end record ;
    | <ident+_,> : <type> <init> ;
    | procedure <ident> <params?> is <decl*>
      begin <instr+> end <ident?> ;
    | function <ident> <params?> return <type> is <decl*>
      begin <instr+> end <ident?> ;
    
<init>
    : := <expr>
    | ^

<decl*>
    : <decl><decl*>
    | <decl>
    | ^

<champs>
    : <ident+_,> : <type> ;

<type>
    : <ident>
    | access <ident>

<params>
    : ( <param+_;> )

<params?>
    : <params>
    | ^

<param>
    : <ident+_,> : <mode?> <type>

<param+_;>
    : <param>;<param+_;>
    | <param>

<mode>
    : in
    | in out

<mode?>
    : <mode>
    | ^

<expr>
    : <entier>
    | <caractere>
    | true
    | false
    | null
    | ( <expr> )
    | <access>
    | <expr> <operateur> <expr>
    | not <expr>
    | - <expr>
    | new <ident>
    | <ident> ( <expr+_,> )
    | character ' val ( <expr> )

<expr+_,>
    : <expr>,<expr+_,>
    | <expr>

<expr?>
    : <expr>
    | ^

<instr>
    : <access> := <expr> ;
    | <ident> ;
    | <ident> ( <expr+_,> ) ;
    | return <expr?> ;
    | begin <instr+> end ;
    | if <expr> then <instr+> <else_if*>
      (else <instr+>)? end if ;
    | for <ident> in <reverse> <expr> .. <expr>
      loop <instr+> end loop ;
    | while <expr> loop <instr+> end loop ;

<instr+>
    : <instr><instr+>
    | <instr>

<else_if>
    : elsif <expr> then <instr+>

<else_if*>
    : <else_if><else_if*>
    | <else_if>
    | ^

<operateur>
    : =
    | /=
    | <
    | <=
    | >
    | >=
    | +
    | -
    | *
    | /
    | rem
    | and
    | and then
    | or
    | or else

<access>
    : <ident>
    | <expr> . <ident>

<reverse>
    : reverse
    | ^