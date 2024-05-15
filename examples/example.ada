with Ada.Text_IO; use Ada.Text_IO;
procedure Test is
    type ExampleRecord is
            record
                A : Integer;
                B : Integer;
            end record;
    B : ExampleRecord;

    function EstPremier(N : Integer) return Boolean is
        Racine : Integer;
        B : Boolean := True;
    begin
        if N <= 1 then
            return False;
        elsif N <= 3 then
            return True;
        end if;

        Racine := 1;
        while Racine * Racine <= N loop
            Racine := Racine + 1;
        end loop;
        Racine := Racine - 1;

        for I in 2..Racine loop
            if N rem I = 0 then
                B := False;
            end if;
        end loop;

        return B;
    end EstPremier;
    -- Fonction récursive pour calculer la factorielle
    function Factorielle(N : Integer) return Integer is
    begin
        if N <= 1 then
            return 1;
        else
            return N * Factorielle(N - 1);
        end if;
    end Factorielle;

    -- Procédure pour calculer la factorielle d'un nombre premier
    -- et retourner le résultat dans AdresseResult
    function CalculerFactorielle(X : in Integer; AdresseResult : in out Integer) return ExampleRecord is
        Result : ExampleRecord;
    begin
        return Result;
    end CalculerFactorielle;
    Result : Integer;
begin
    B := CalculerFactorielle(5, Result);
    if Result = -1 then
        Put('L');
        Put('e');
        Put(' ');
        Put('n');
        Put('o');
        Put('m');
        Put('b');
        Put('r');
        Put('e');
        Put(' ');
        Put('n');
        Put(' ');
        Put('e');
        Put('s');
        Put('t');
        Put(' ');
        Put('p');
        Put('a');
        Put('s');
        Put(' ');
        Put('p');
        Put('r');
        Put('e');
        Put('m');
        Put('i');
        Put('e');
        Put('r');
        Put('.');
        New_Line;
    else
        Put('F');
        Put('a');
        Put('c');
        Put('t');
        Put('o');
        Put('r');
        Put('i');
        Put('e');
        Put('l');
        Put('l');
        Put('e');
        Put(' ');
        Put('d');
        Put('e');
        Put(' ');
        Put('5');
        Put(' ');
        Put(':');
        Put(' ');
        New_Line;
    end if;

    Put(Result);
    New_Line;
end Test;