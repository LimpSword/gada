with Ada.Text_IO; use Ada.Text_IO;

procedure Record4 is begin
   V.A := '1';
   V.B := '2';
   X.C := '0';
   X.D := V;
   X.E := '3';
end Record4;

--  Local Variables:
--  compile-command: "gnatmake record4.adb && ./record4"
--  End: