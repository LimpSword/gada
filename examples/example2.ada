with Ada.Text_IO; use Ada.Text_IO;
procedure Test is
    S : Integer := 2;
    procedure Pt(V: Integer) is
    begin
        S := V;
        if V > 0 then
            Pt(V-1);
        end if;
    end Pt;
begin
    Pt(3);
end;
