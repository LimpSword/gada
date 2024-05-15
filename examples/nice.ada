with Ada.Text_IO; use Ada.Text_IO;
procedure Test is
    S : Integer := 2;
    procedure Test(K : Integer; N : in out Integer; J : Integer) is
    begin
        Put('N'); Put(' '); Put(N); New_Line;
        N := 3 + K + J;
        Put('N'); Put(' '); Put(N); New_Line;
        Put('K'); Put(' '); Put(K); New_Line;
        Put('J'); Put(' '); Put(J); New_Line;
    end Test;
begin
    Put(S); New_Line;
    Test(1, S, 5);
    Put(S);
end;
