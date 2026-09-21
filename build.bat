@echo off
setlocal
where python >nul 2>nul || (echo Python 3 is required & exit /b 1)
where go >nul 2>nul || (echo Go is required & exit /b 1)
python tools\generate_assets.py || exit /b 1
go build -trimpath -ldflags "-H=windowsgui -s -w" -o SNAPWAVE_HORROR_V8.exe . || exit /b 1
echo Built: SNAPWAVE_HORROR_V8.exe
