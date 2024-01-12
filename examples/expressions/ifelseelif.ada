with Ada.Text_IO; use Ada.Text_IO;
procedure Test is begin
   if 0=1 then Q(0); elsif 1=2 then Q(1); elsif 2=3=4=5=6=7=8=9 then Q(5);else Q(2); end if;
end Test;