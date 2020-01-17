SET CORPUS=corpus
SET GO111MODULE=off
del /Q fuzz-build.zip
del /Q %CORPUS%\crashers\*.*
del /Q %CORPUS%\suppressions\*.*
SET /a PROCS=%NUMBER_OF_PROCESSORS%

go-fuzz-build -o=fuzz-build.zip -func=Fuzz .

:LOOP
go run timeout.go -duration=5m go-fuzz -minimize=5s -bin=fuzz-build.zip -workdir=%CORPUS% -procs=%PROCS%
GOTO LOOP
