with Ada.Text_IO; use Ada.Text_IO;

procedure Pascal is

   type Node;
   type Node is record
      Value: Integer;
      Next  : List;
      Father: Integer;
   end record;

   function F return integer is begin
    if 0=1 then return;
    elsif 1=2 or 2=3 or else 3=4 then return 1;
    elsif 4=5 and 5=6 and then 6=7 then return 2;
    else return 3;
    end if;
    end;
begin F; end;