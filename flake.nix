{
  description = "Полноценная DevOps среда разработки для Jash (Go + Python)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux"; # Смените на "aarch64-darwin" для Mac Apple Silicon
      pkgs = import nixpkgs { inherit system; };
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          # --- Основные языки ---
          go_latest
          (python3.withPackages (ps: with ps: [ pip ]))

          # --- Линтеры и Качество Кода ---
          golangci-lint # Главный линтер для Go-кода
          gofumpt        # Более строгий аналог gofmt для форматирования
          python3Packages.flake8 # Проверка синтаксиса для builderexe.py

          # --- Инструменты тестирования и отладки ---
          delve         # Отладчик (debugger) для Go-кода (`dlv`)
          gotestsum     # Красивый вывод тестов в терминале

          # --- Утилиты автоматизации ---
          pre-commit    # Для проверки кода перед каждым git commit
          gnumake       # Если решите написать Makefile для автоматизации скриптов
        ];

        # Переменные окружения для Go-разработчика
        env = {
          CGO_ENABLED = "0"; # Для сборки чистых статически слинкованных бинарников
        };

        shellHook = ''
          echo "========================================================"
          echo "  Добро пожаловать в расширенную среду Jash!"
          echo "  Доступны инструменты качества кода: golangci-lint, gofumpt"
          echo "========================================================"

          # Настройка pre-commit хуков, если они объявлены в проекте
          if [ -f .pre-commit-config.yaml ]; then
            pre-commit install
          fi

          # Полезные алиасы автоматизации
          alias jash-lint="golangci-lint run"
          alias jash-fmt="gofumpt -w ."
          alias jash-test="gotestsum --format short-verbose"
        '';
      };

      # БОНУС: Создание Docker-образа вашего проекта через Nix БЕЗ самого Docker
      packages.${system}.dockerImage = pkgs.dockerTools.buildImage {
        name = "jash-interpreter";
        tag = "latest";
        config = {
          Cmd = [ "${self.packages.${system}.default}/bin/jash" ];
        };
      };
    };
}
