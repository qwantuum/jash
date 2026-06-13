{
  description = "Полная DevOps среда разработки для Jash (Go + Python + Node.js)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          # Основные языки
          go_latest
          nodejs_latest
          (python314.withPackages (ps: [ ps.pip ]))

          # Инструменты качества кода и сборки
          golangci-lint
          gofumpt
          gnumake
        ];

        shellHook = ''
          echo "========================================================"
          echo "  Окружение Jash успешно запущено!"
          echo "  Установлены: $(go version)"
          echo "  Node.js:      $(node -v)"
          echo "  Python:       $(python --version)"
          echo "========================================================"
          echo "  Для сборки проекта запустите: node build.js"
          echo "========================================================"
        '';
      };
    };
}
