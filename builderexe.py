import os
import subprocess
import shutil

# --- КОНФИГУРАЦИЯ ---
APP_NAME = "jash.exe"
MAIN_PATH = "./cmd/jash/main.go"  # Убедись, что тут правильное имя файла (main.go или jash.go)
TARGET_DIR = r"C:\jash-lang"      # Путь к папке на диске C

def build_and_deploy():
    print(f"=== Старт сборки Jash под Windows ===")
    
    # 1. Жестко очищаем целевую папку на Диске C
    if os.path.exists(TARGET_DIR):
        print(f"[...] Папка {TARGET_DIR} существует. Полностью удаляем старые файлы...")
        try:
            shutil.rmtree(TARGET_DIR)
        except Exception as e:
            print(f"[-] Не удалось удалить старые файлы: {e}")
            print("[!] Подсказка: Закрой запущенный jash.exe, если он открыт в другом терминале.")
            return
            
    # Создаем пустую чистую папку заново
    os.makedirs(TARGET_DIR)
    print(f"[+] Папка {TARGET_DIR} успешно очищена и создана заново.")
    
    # 2. Настраиваем окружение для компиляции Go
    env = os.environ.copy()
    env["GOOS"] = "windows"
    env["GOARCH"] = "amd64"
    
    # Флаги -s -w сжимают бинарник, убирая отладку
    out_path = os.path.join(TARGET_DIR, APP_NAME)
    cmd = f'go build -ldflags="-s -w" -o "{out_path}" {MAIN_PATH}'
    
    print(f"[...] Компиляция и перенос {APP_NAME}...")
    result = subprocess.run(cmd, shell=True, env=env)
    
    # 3. Проверяем результат
    if result.returncode == 0:
        print(f"\n[🎉] Успешно скомпилировано!")
        print(f"[+] Свежий файл установлен в: {out_path}")
        print(f"[+] Размер бинарника: {round(os.path.getsize(out_path) / (1024*1024), 2)} МБ")
    else:
        print("\n[-] Ошибка: Сборка не удалась. Проверь синтаксис Go-кода или путь к main.go.")

if __name__ == "__main__":
    build_and_deploy()
