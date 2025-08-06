{
  description = "Wails desktop app";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
  };

  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      packages.${system}.default = pkgs.buildGoModule rec {
        pname = "llyfr";
        version = "1.0.0";

        src = ./.;

        vendorHash = null; # or provide if using vendored modules
        buildFlags = [ "-mod=vendor" ];

        nativeBuildInputs = with pkgs; [
          pkg-config
          nodePackages.nodejs
          wails
        ];

        buildInputs = with pkgs; [
          makeWrapper
          webkitgtk_4_1
          gtk3
        ];

        buildPhase = ''
          wails build -clean -o ${pname}
        '';

        installPhase = ''
          mkdir -p $out/bin
          cp build/bin/${pname} $out/bin/
        '';
      };

      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [
          go
          gopls
          gotools
          hugo
          wails
          prettier-plugin-go-template

          pkg-config
          makeWrapper
          gtk3
          webkitgtk_4_1

          nodePackages.nodejs
          vscode-langservers-extracted
          emmet-ls
          typescript-language-server
          taplo
          yaml-language-server
          markdown-oxide
          prettierd
          eslint_d
          mdformat
          svelte-language-server
        ];
        shellHook = ''
          echo "wails dev shell"
        '';
      };
    };
}
