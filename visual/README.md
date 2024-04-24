# Print avec VisUAL

Le programme présent permet d'exécuter un programme assembleur écrit pour VisUAL et d'afficher les entiers ou
chaînes de caractères en utilisant la procédure `println`.

Ce programme ne sert que d'exécution, il ne modifie pas le programme assembleur. Ainsi, votre programme assembleur doit
être adapté dès le début en suivant les instructions ci-dessous.

## Utilisation

Vous pouvez exécuter votre programme assembleur et afficher son output avec la commande suivante :

```bash
java -jar pcl.jar <programme>
```

Par exemple `java -jar pcl.jar example.s`

## Adapter votre code

Ajoutez par défaut les instructions suivantes à votre programme assembleur.

### Le code à ajouter

Tout d'abord le programme doit commencer par une déclaration de l'espace nécessaire à l'affichage.

```armasm
STR_OUT      FILL    0x1000
```

Ajoutez ensuite la procédure d'affichage à votre programme (par exemple à la fin).

```armasm
println      STMFD   SP!, {LR, R0-R3}
             MOV     R3, R0
             LDR     R1, =STR_OUT ; address of the output buffer
PRINTLN_LOOP LDRB    R2, [R0], #1
             STRB    R2, [R1], #1
             TST     R2, R2
             BNE     PRINTLN_LOOP
             MOV     R2, #10
             STRB    R2, [R1, #-1]
             MOV     R2, #0
             STRB    R2, [R1]


             ;       we need to clear the output buffer
             LDR     R1, =STR_OUT
             MOV     R0, R3
CLEAN        LDRB    R2, [R0], #1
             MOV     R3, #0
             STRB    R3, [R1], #1
             TST     R2, R2
             BNE     CLEAN
             ;       clear 3 more
             STRB    R3, [R1], #1
             STRB    R3, [R1], #1
             STRB    R3, [R1], #1

             LDMFD   SP!, {PC, R0-R3}
```

### Comment appeler la procédure Put ?

Remarque : les explications suivantes peuvent sensiblement différer selon la façon dont vous gérez votre pile.

#### Idée générale

Vous devez stocker en pile la valeur que vous souhaitez afficher précédée d'un '0'. La limite théorique de taille à
afficher est grande (de la taille de STR_OUT).

#### Du code

On suppose que la valeur à afficher se trouve dans R0. L'idée reste la même si vous souhaitez afficher quelque chose de
plus grand que 4 octets.

Pour appeler la procédure `println` :

```armasm
SUB SP, SP, #4   ; réservez 4 octets pour le 0
MOV R1, #0
STR R1, [SP]
SUB SP, SP, #4   : réservez 4 octets pour la valeur (ou plus)
STR R0, [SP]     ; stockez la valeur
MOV R0, SP       ; adresse de la valeur
BL println
ADD SP, SP, #8   ; libérez la pile
```

### Remarque importante

Imaginons que vous disposiez dans la pile de la valeur suivante : 12594 en décimal.
Cette valeur se traduit par 0x3132 en hexadécimal.
Le programme lit ainsi les octets 0x31 et 0x32 et les traduit en ASCII soit '1' et '2'. L'affichage sera ainsi '12'.

Que faut-il en tirer ?
- Si vous souhaitez afficher un caractère : stockez-le en pile précédé d'un '0'.
- Si vous souhaitez ajouter plus d'un caractère : stockez-les en pile un par un.
- Si vous souhaitez afficher un entier : stocker en pile chaque chiffre de l'entier précédé d'un '0'.

## Fonctionnement

La technique utilisée est l'exploitation des breakpoints de VisUAL.

Le programme Java lit d'abord le programme assembleur à la recherche de la ligne `STRB    R2, [R1, #-1]` indiquant la
fin
de l'affichage (ie. le buffer est entièrement stocké). La ligne correspondante (en commençant à 0) est alors stockée et
considérée comme breakpoint.

La ligne utilisée permet d'ajouter un `\n` à la fin de l'affichage. Si vous souhaitez afficher sans saut de ligne, vous
pouvez supprimer la ligne précédente (le MOV de 0) ou recréer une procédure `print` sans le MOV (il faut tout de même le
STRB pour le breakpoint).

On exécute ainsi le programme assembleur à l'aide de la version headless de VisUAL en lui donnant les arguments
nécessaires.
Basiquement, lors de l'exécution, le programme s'arrête à chaque breakpoint et produit du contenu dans un fichier de
log. Le fichier de log est ensuite lu par le programme Java et extrait le contenu de STR_OUT lors de chaque affichage.
Le fichier de log contient la valeur en hexadécimal, convertie en ASCII lors de la lecture des logs.

L'exécution est plutôt verbose, si vous souhaitez simplement afficher le résultat de l'exécution, vous devrez extraire
la bonne information de l'output.