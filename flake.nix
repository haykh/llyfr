{
  description = "Llyfr flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };

        pname = "llyfr";
        version = "1.0.0";

        frontend = pkgs.buildNpmPackage {
          pname = "${pname}-frontend";
          inherit version;

          src = ./frontend;

          npmDepsHash = "sha256-7SoU4WzwScTF8cUIZw3u0V8VSCJre2LkRpmE3q9L+WQ=";

          npmBuildScript = "build";

          installPhase = ''
            runHook preInstall
            mkdir -p $out
            cp -r dist $out/dist
            runHook postInstall
          '';
        };
      in
      {
        packages.default = pkgs.buildGoModule {
          inherit pname version;

          src = ./.;

          vendorHash = "sha256-r6WnXz8Low1xsElJTZgTuS/57ezEGhCw81ba8G00P7o=";

          tags = [
            "desktop"
            "production"
            "webkit2_41"
          ];
          ldflags = [
            "-w"
            "-s"
          ];

          nativeBuildInputs = with pkgs; [
            pkg-config
            wrapGAppsHook3
          ];

          buildInputs = with pkgs; [
            gtk3
            webkitgtk_4_1
          ];

          preBuild = ''
            rm -rf frontend/dist
            mkdir -p frontend
            cp -r ${frontend}/dist frontend/dist
          '';

        };

        devShells.default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            go
            nodejs
            wails
            pkg-config
          ];

          buildInputs = with pkgs; [
            gtk3
            webkitgtk_4_1
          ];
        };
      }
    );
}
