import os
import subprocess
import shutil

# --- КОНФИГУРАЦИЯ ---
APP_NAME = "jash.exe"
MAIN_PATH = "./cmd/jash/main.go"  # Точка входа
BUILD_DIR = "./build"             # Папка для готового файла

def build_jash_windows():
    print("=== Старт сборки Jash под Windows ===")
    
    # 1. Очищаем старую сборку, если она была
    if os.path.exists(BUILD_DIR):
        shutil.rmtree(BUILD_DIR)
    os.makedirs(BUILD_DIR)
    
    # 2. Выставляем переменные окружения для Go
    # Архитектура amd64 — это стандарт для 64-битной Windows 10
    env = os.environ.copy()
    env["GOOS"] = "windows"
    env["GOARCH"] = "amd64"
    
    # Флаги -s -w убирают отладочные символы, сжимая бинарник в 2 раза
    out_path = os.path.join(BUILD_DIR, APP_NAME)
    cmd = f'go build -ldflags="-s -w" -o {out_path} {MAIN_PATH}'
    
    print(f"[...] Компиляция файла {APP_NAME}...")
    result = subprocess.run(cmd, shell=True, env=env)
    
    # 3. Проверяем итог
    if result.returncode == 0:
        print(f"\n[🎉] Сборка успешно завершена!")
        print(f"[+] Готовый файл лежит тут: {os.path.abspath(out_path)}")
        print(f"[+] Вес файла: {round(os.path.getsize(out_path) / (1024*1024), 2)} МБ")
    else:
        print("\n[-] Ошибка: Не удалось скомпилировать проект. Проверь код на Go.")

if __name__ == "__main__":
    build_jash_windows()
