with Ada.Text_IO ; use Ada.Text_IO ;

procedure unDebut is

   function aireRectangle(larg : integer; long : integer) return integer is
   aire: integer;
   begin
      aire := larg*long ;
   return aire
   end aireRectangle ;

   function perimetreRectangle(larg : integer; long : integer) return integer is
   p : integer;
   begin
      p := larg*2 + long*2 ;
   return p
   end perimetreRectangle;

      -- VARIABLES
choix : integer ;

      -- PROCEDURE PRINCIPALE

begin
   choix := 2;

   if choix = 1
      then valeur := perimetreRectangle(2,3) ;
         put(valeur) ;
      else valeur := aireRectangle(2,3) ;
         put(valeur) ;
   end if;
end unDebut ;