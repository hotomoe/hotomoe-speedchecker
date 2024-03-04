for /f "delims=" %%A in ('git rev-parse --show-toplevel') do (cd %%A)
copy /Y assets\syso.json syso.json
%USERPROFILE%\go\bin\syso.exe
del syso.json
SET GOOS=windows& go build  -ldflags "-s -w" -o "dist/hnc-windows-amd64.exe"
SET GOOS=linux& go build  -ldflags "-s -w" -o "dist/hnc-linux-amd64"
SET GOOS=darwin& go build  -ldflags "-s -w" -o "dist/hnc-darwin-amd64"
del out.syso
upx --best --lzma "dist/hnc-windows-amd64.exe"
upx --best --lzma "dist/hnc-linux-amd64"
REM upx --best --lzma "dist/hnc-darwin-amd64"
