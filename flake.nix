{
  description = "connectionInfo - A lightweight service that displays client connection information";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = {
          default = pkgs.buildGoModule {
            pname = "connectionInfo";
            version = "0.1.0";
            src = ./.;
            vendorHash = null; # No external dependencies

            meta = with pkgs.lib; {
              description = "A lightweight service that displays client connection information";
              license = licenses.mit;
              maintainers = [ ];
            };
          };

          connectionInfo = self.packages.${system}.default;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
          ];
        };
      }
    ) // {
      nixosModules.default = { config, lib, pkgs, ... }:
        let
          cfg = config.services.connectionInfo;
        in
        {
          options.services.connectionInfo = {
            enable = lib.mkEnableOption "connectionInfo service";

            port = lib.mkOption {
              type = lib.types.port;
              default = 8080;
              description = "Port for the connectionInfo server to listen on";
            };

            openFirewall = lib.mkOption {
              type = lib.types.bool;
              default = false;
              description = "Whether to open the firewall for the configured port";
            };

            package = lib.mkOption {
              type = lib.types.package;
              default = self.packages.${pkgs.system}.default;
              defaultText = lib.literalExpression "pkgs.connectionInfo";
              description = "The connectionInfo package to use";
            };

            basePath = lib.mkOption {
              type = lib.types.str;
              default = "/connectionInfo";
              description = "Base URL path for the service (e.g., /connectionInfo). Empty means serve at the root of the virtual host.";
              example = "/connectionInfo";
            };

            nginx = {
              enable = lib.mkOption {
                type = lib.types.bool;
                default = true;
                description = "Whether to enable the built-in nginx reverse proxy configuration. Enabled by default so the service is self-contained.";
              };

              virtualHost = lib.mkOption {
                type = lib.types.str;
                default = "localhost";
                description = "nginx virtual host name under which to serve connectionInfo.";
                example = "website.com";
              };

              forceSSL = lib.mkOption {
                type = lib.types.bool;
                default = false;
                description = "Force SSL for the virtual host.";
              };

              enableACME = lib.mkOption {
                type = lib.types.bool;
                default = false;
                description = "Enable ACME (Let's Encrypt) for the virtual host.";
              };
            };
          };

          config = lib.mkIf cfg.enable {
            systemd.services.connectionInfo = {
              description = "Connection Info Service";
              wantedBy = [ "multi-user.target" ];
              after = [ "network.target" ];

              environment = {
                PORT = toString cfg.port;
              };

              serviceConfig = {
                Type = "exec";
                ExecStart = "${cfg.package}/bin/connectionInfo";
                Restart = "on-failure";
                RestartSec = "5s";

                # Security hardening
                DynamicUser = true;
                NoNewPrivileges = true;
                ProtectSystem = "strict";
                ProtectHome = true;
                PrivateTmp = true;
                PrivateDevices = true;
                ProtectKernelTunables = true;
                ProtectKernelModules = true;
                ProtectControlGroups = true;
                RestrictAddressFamilies = [ "AF_INET" "AF_INET6" ];
                RestrictNamespaces = true;
                LockPersonality = true;
                RestrictRealtime = true;
                MemoryDenyWriteExecute = true;
                SystemCallArchitectures = "native";
              };
            };

            networking.firewall.allowedTCPPorts = lib.mkIf cfg.openFirewall [ cfg.port ];

            services.nginx = lib.mkIf cfg.nginx.enable (lib.mkMerge [
              {
                enable = true;
                recommendedProxySettings = true;
              }

              # When basePath is empty: dedicated virtual host mode
              (lib.mkIf (cfg.basePath == "") {
                virtualHosts.${cfg.nginx.virtualHost} = {
                  forceSSL = cfg.nginx.forceSSL;
                  enableACME = cfg.nginx.enableACME;
                  locations."/" = {
                    proxyPass = "http://127.0.0.1:${toString cfg.port}";
                  };
                };
              })

              # When basePath is set: sub-path mode on existing virtual host
              (lib.mkIf (cfg.basePath != "") {
                virtualHosts.${cfg.nginx.virtualHost} = {
                  forceSSL = cfg.nginx.forceSSL;
                  enableACME = cfg.nginx.enableACME;
                  # Redirect bare path to path with trailing slash
                  locations."= ${cfg.basePath}" = {
                    return = "301 ${cfg.basePath}/";
                  };
                  # Proxy with prefix stripping (trailing slash on proxy_pass)
                  locations."${cfg.basePath}/" = {
                    proxyPass = "http://127.0.0.1:${toString cfg.port}/";
                  };
                };
              })
            ]);
          };
        };

      overlays.default = final: prev: {
        connectionInfo = self.packages.${prev.system}.default;
      };
    };
}
