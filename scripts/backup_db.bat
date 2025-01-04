@echo off
setlocal

:: Configuration
set PGUSER=postgres
set PGPASSWORD=your_password_here
set BACKUP_DIR=C:\PostgreSQL\backups
set DB_NAME=HE_horse_db

:: Create backup directory if it doesn't exist
if not exist "%BACKUP_DIR%" mkdir "%BACKUP_DIR%"

:: Generate timestamp
for /f "tokens=2-4 delims=/ " %%a in ('date /t') do (set DATE=%%c_%%a_%%b)
for /f "tokens=1-2 delims=/:" %%a in ('time /t') do (set TIME=%%a_%%b)
set TIMESTAMP=%DATE%_%TIME%

:: Perform backup
"C:\Program Files\PostgreSQL\17\bin\pg_dump.exe" -h localhost -U %PGUSER% -F c -b -v -f "%BACKUP_DIR%\%DB_NAME%_%TIMESTAMP%.backup" %DB_NAME%

:: Delete files older than 7 days
forfiles /P "%BACKUP_DIR%" /M *.backup /D -7 /C "cmd /c del @path"

echo Backup completed: %BACKUP_DIR%\%DB_NAME%_%TIMESTAMP%.backup 