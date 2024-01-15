with Ada.Text_IO; use Ada.Text_IO;
procedure Test is begin
    x := 1 or else 2 or 1 and then 2 and then 1 or 3;
end Test;