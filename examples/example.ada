with Ada.Text_IO; use Ada.Text_IO;
procedure Test is
    S : Integer := 2;
    procedure Pt(I : Integer) is
        D : Integer := 3;
    begin
        S := S + D + I;
    end Pt;
begin
    Pt(4);
    S := S + 1;
end;
