import { execSync } from "child_process";
import path from "path";
import fs from "fs";

// Конфигурация сборки
const CONFIG = {
  appName: "jash",
  srcPath: "./cmd/jash", // Путь к главному go-файлу на основе вашего репозитория
  distDir: "./dist",
  targets: [
    { os: "linux", arch: "amd64", filename: "jash-linux" },
    { os: "windows", arch: "amd64", filename: "jash.exe" },
    { os: "darwin", arch: "arm64", filename: "jash-mac-m1" }, // Опционально: Mac Apple Silicon
  ],
};

function build() {
  console.log("🚀 Запуск кросс-компиляции Jash...");

  // Создаем папку назначения, если её нет
  if (!fs.existsSync(CONFIG.distDir)) {
    fs.mkdirSync(CONFIG.distDir, { recursive: true });
  }

  // Проходим по каждой целевой платформе
  CONFIG.targets.forEach((target) => {
    const outputPath = path.join(CONFIG.distDir, target.filename);
    console.log(
      `📦 Сборка для ${target.os} (${target.arch}) -> ${outputPath}...`,
    );

    // Настраиваем переменные окружения для Go
    const env = {
      ...process.env,
      GOOS: target.os,
      GOARCH: target.arch,
      CGO_ENABLED: "0", // Отключаем CGO для переносимости бинарников
    };

    try {
      // Выполняем сборку
      execSync(`go build -ldflags="-s -w" -o ${outputPath} ${CONFIG.srcPath}`, {
        env,
        stdio: "inherit",
      });
      console.log(`✅ Успешно собрано: ${target.filename}`);
    } catch (error) {
      console.error(`❌ Ошибка сборки под ${target.os}:`, error.message);
      process.exit(1);
    }
  });

  console.log(
    "\n🎉 Все бинарники успешно скомпилированы в папку " + CONFIG.distDir,
  );
}

build();
