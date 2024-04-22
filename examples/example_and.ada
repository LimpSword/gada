with Ada.Text_IO; use Ada.Text_IO;
procedure Test is
    S : Integer := 2;
    function Example(k : Integer; n : Integer) return Integer is
    begin
        if k = 2 and n = 3 then
            return 1;
        end if;
        return 0;
    end Example;
begin
    S := Example(2, 4);
end;
