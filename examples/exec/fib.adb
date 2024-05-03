with Ada.Text_IO; use Ada.Text_IO;

procedure Fib is

   procedure PrintInt(N: Integer) is
      C: Integer := N rem 10;
   begin
      if N > 9 then PrintInt(N / 10); end if;
      Put(Character'Val(48 + C));
   end;

   f : integer := -1;

   procedure fib(n	: integer) is
      procedure somme is
         tmp : integer;
      begin fib(n-2); tmp := f; fib(n-1); f := f + Tmp; return; end;
   begin
      if n <= 1 then f := n; else somme; end if; return;
   end;

begin fib(13); printint(f); end;
