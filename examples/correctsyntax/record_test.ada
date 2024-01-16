with Ada.Text_IO;
use Ada.Text_IO;

procedure Main is
   type Person is record
      Name : String;
      Age  : Integer;
   end record;
begin
    put(Person);
end Main;
