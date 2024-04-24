with Ada.Text_IO; use Ada.Text_IO;
procedure Test is
    S : Integer := 2;
    function Factorial(N : Integer) return Integer is
    begin
        if N = 0 then
            return 1;
        end if;

        return N * Factorial(N - 1);
    end Factorial;
begin
    S := Factorial(12);
    Put(Character'Val(S));
end;
