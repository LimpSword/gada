with Ada.Text_IO; use Ada.Text_IO;

procedure Power is
   function Power(X, N: Integer) return Integer is
      P: Integer;
   begin
      if N = 0 then return 1; end if;
      P := Power(X, N / 2);
      P := P * P;
      if N rem 2 = 1 then P := X * P; end if;
      return P;
   end;

begin
   put(Power(2, 0)); New_Line;
   put(Power(2, 8)); New_Line;
   put(Power(2, 30)); New_Line;
end;

--  Local Variables:
--  compile-command: "gnatmake power.adb && ./power"
--  End:

