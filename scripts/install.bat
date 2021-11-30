@echo ===============================   MongoDB Local Backup Service for Windows   =============================
@echo ==============================  Version：V1.0.0   ============================
@echo ===========================   Created：2022-01-12  ===========================
@echo ======================= ServiceName：MongoDBLocalBackupService ======================
@echo ================ Copyright: @2013-2022 Will All Rights Reserved. =============
@echo.
@echo.
@echo.

@echo off

mode con cols=100 lines=20
color 3f
set name=MongoDBLocalBackupService
set display_name=MongoDB Local Backup Service
set path=%~dp0%startup.bat
setlocal
set uac=~uac_permission_tmp_%random%
md "%SystemRoot%\system32\%uac%" 2>nul
if %errorlevel%==0 ( rd "%SystemRoot%\system32\%uac%" >nul 2>nul ) else (
    echo set uac = CreateObject^("Shell.Application"^)>"%temp%\%uac%.vbs"
    echo uac.ShellExecute "%~s0","","","runas",1 >>"%temp%\%uac%.vbs"
    echo WScript.Quit >>"%temp%\%uac%.vbs"
    "%temp%\%uac%.vbs" /f
    del /f /q "%temp%\%uac%.vbs" & exit )
endlocal

sc create %name% binPath=%path% start= AUTO displayname=%name%
sc description %name%  "Export MongoDB data into files in JSON format"
net start %name%
 
@echo.