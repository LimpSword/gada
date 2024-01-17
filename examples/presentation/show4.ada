with Ada.Text_IO; use Ada.Text_IO;
procedure Test is begin
   if not not (0=1) then y := 3;y := Character'val(y);
   else
   y :=1;
   for i in 0..N loop
        y := y + -1;
        Put(y);
   end loop;
   end if;
end Test;