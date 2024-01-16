with Ada.Text_IO; use Ada.Text_IO;
procedure Test is begin
    D := 1;
    for I in 0 .. N-1 loop
       if Mem(D, A) and then not Mem(D, B) and then not Mem(D, C) then
          F := F + T(A-D, (B+D) * 2, (C+D) / 2);
       end if;
       D := 2 * D;
    end loop;
end Test;