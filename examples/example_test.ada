with Ada.Text_IO; use Ada.Text_IO;
procedure Test is
    S : Integer := 2;
    function Fibonacci(n : Integer) return Integer is
    begin
        if n = 0 then
            return 0;
        elsif n = 1 then
            return 1;
        else
            return Fibonacci(n - 1) + Fibonacci(n - 2);
        end if;
    end Fibonacci;
begin
    S := Fibonacci(10);
end;
