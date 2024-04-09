with Ada.Text_IO; use Ada.Text_IO;

procedure Record4 is
   type Person is record
      Name : Character;
      Age  : Integer;
   end record;
   function f(p, pd : Person) return Integer is
   begin
        return p.Age;
   end f;
   null_person : Person;
   i : Integer;
begin
    i := f(null_person, null_person);
end Record4;
