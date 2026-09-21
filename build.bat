@echo off
setlocal
where python >nul 2>nul || (echo Python 3 is required & exit /b 1)
where go >nul 2>nul || (echo Go 1.22+ is required & exit /b 1)
python tools\generate_assets.py || exit /b 1
python tools\audit.py || exit /b 1
echo Built: SNAPWAVE_HORROR_V8_2_AUDITED.exe
