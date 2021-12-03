CLIExample2
| - word     flag string 
| - myint      flag int64

subCmd : foo  Alias ( fc ,fo )
  | -flag name word string 

subCmd : bar
  | - flag name floatnumber flaot64

CLIexample2 cant execute on its OWN ..SubCmds foo or bar is must
So it has no RUN funciton

SubCommands ... 
add subcommand in init of subCommand command

if no URL module then for local use:
go mod init cliexample2
cobra init --pkg-name cliexample2

OR use below if there is a http or ssh repo with URL
go mod init github.com/user/cliexample2
cobra init --pkg-name github.com/user/cliexample2


# To run :
./cliexample2 -h 
./cliexample2 -i 20 -w "testword" 

To see validations try giving Operations value 0 or -1 or same for workers.
