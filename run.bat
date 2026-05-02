@echo off
title ITSTEP Social - Launcher
echo ========================================
echo Starting ITSTEP Social System (Firestore Mode)
echo ========================================

:: 1. Запуск Бэкенда (Golang) в НОВОМ окне
echo [1/2] Starting Backend...
start "ITSTEP Backend" cmd /k "cd backend && go run cmd/app/main.go"

:: 2. Запуск Фронтенда (Nuxt 3) в ЭТОМ окне
echo [2/2] Starting Frontend...
cd frontend
npm run dev

pause
