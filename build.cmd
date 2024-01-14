@echo off
setlocal enabledelayedexpansion

set VERSION=1.2.0
set OUTPUT_DIR=bin
set BINARY_NAME=generate-geoip-geosite

:: Определение пути к 7-Zip
set SEVENZIP_PATH="D:\Program Files\7-Zip\7z.exe"

:: Поддерживаемые ОС
set OS_LIST=linux windows

:: Поддерживаемые архитектуры
set ARCH_LIST=amd64 arm64 386

:: Создаем выходную директорию, если её нет
if not exist %OUTPUT_DIR% mkdir %OUTPUT_DIR%

:: Компилируем и упаковываем для каждой комбинации ОС и архитектуры
for %%o in (%OS_LIST%) do (
  for %%a in (%ARCH_LIST%) do (

    set ARCHIVE_NAME=%CD%\%OUTPUT_DIR%\%BINARY_NAME%-!VERSION!-%%o-%%a

    :: Компиляция
    set GOOS=%%o
	set GOARCH=%%a
	
    echo Compiling for %%o-%%a:
	
    :: Упаковка в архив
    if "%%o"=="windows" (
		go build -o %CD%\%OUTPUT_DIR%\!BINARY_NAME!.exe
		%SEVENZIP_PATH% a !ARCHIVE_NAME!.zip %CD%\%OUTPUT_DIR%\!BINARY_NAME!.exe
		:: Удаление временного бинарного файла
		del %CD%\%OUTPUT_DIR%\!BINARY_NAME!.exe
    ) else (
		go build -o %CD%\%OUTPUT_DIR%\!BINARY_NAME!
		tar -czvf !ARCHIVE_NAME!.tar.gz -C %CD%\%OUTPUT_DIR% %BINARY_NAME%
		:: Удаление временного бинарного файла
		del %CD%\%OUTPUT_DIR%\!BINARY_NAME!
    )

  )
)

echo Compilation and packaging completed.
pause

REM sing-box-1.7.1-linux-386.tar.gz
REM sing-box-1.7.1-linux-amd64.tar.gz
REM sing-box-1.7.1-linux-amd64v3.tar.gz
REM sing-box-1.7.1-linux-arm64.tar.gz
REM sing-box-1.7.1-linux-armv7.tar.gz
REM sing-box-1.7.1-linux-s390x.tar.gz
REM sing-box-1.7.1-windows-386-legacy.zip
REM sing-box-1.7.1-windows-386.zip
REM sing-box-1.7.1-windows-amd64-legacy.zip
REM sing-box-1.7.1-windows-amd64.zip
REM sing-box-1.7.1-windows-amd64v3.zip
REM sing-box-1.7.1-windows-arm64.zip

REM generate-geoip-geosite-{version}-{OS}-{arch}.{zip/tar.gz}