with Ada.Text_IO; use Ada.Text_IO;
procedure Test_a is
   procedure R(N: Integer) is
   begin
      if N = 0 then return; end if;
      Put('a');
      R(N - 1);
      for I in reverse 1..10 loop
            if not (I rem 2 = 0) then -- Vérification si I n'est pas pair (not pair)
               Put_Line('2');
            end if;
      end loop;
   end;
begin
   R(42);
   New_Line;
end TEST_a;