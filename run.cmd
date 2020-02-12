SET CORPUS=corpus
SET GO111MODULE=off
SET ARCHIVE=fuzz-build.zip
del /Q %ARCHIVE%
del /Q %CORPUS%\crashers\*.*
del /Q %CORPUS%\suppressions\*.*
SET /a PROCS=%NUMBER_OF_PROCESSORS%

go-fuzz-build -o=%ARCHIVE% -func=Fuzz .

:LOOP
go run timeout.go -duration=5m go-fuzz -minimize=5s -bin=%ARCHIVE% -workdir=%CORPUS% -procs=%PROCS%
GOTO LOOP
