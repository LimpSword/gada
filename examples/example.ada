with Ada.Text_IO; use Ada.Text_IO;
procedure Test is
    S : Integer := 2;
    D : Integer;
begin
    for I in 1..10 loop
        S := S * 2;
    end loop;
end;
