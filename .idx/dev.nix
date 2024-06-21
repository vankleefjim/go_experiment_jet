# To learn more about how to use Nix to configure your environment
# see: https://developers.google.com/idx/guides/customize-idx-env
{ pkgs, ... }: {
  # Which nixpkgs channel to use.
  channel = "unstable"; # "stable-23.11"; # or "unstable"
  # needs unstable to have latest go version
  # Use https://search.nixos.org/packages to find packages
  packages = [
    pkgs.go
    pkgs.nodejs_20
    pkgs.nodePackages.nodemon
    pkgs.gnumake
    pkgs.docker-compose
  ];
  # Sets environment variables in the workspace
  env = {};
  # virtualisation.containers.enable = true;
  #https://nixos.wiki/wiki/Podman
  # virtualisation  = {
  #   podman = {
  #     enable = true;
  #     dockerCompat=true;
  #      defaultNetwork.settings.dns_enabled = true;
  #   };
  # };
  services.docker.enable=true;
  # services.podman.enable=true;
  idx = {
    # Search for the extensions you want on https://open-vsx.org/ and use "publisher.id"
    extensions = [
      "golang.go"
      "tanhakabir.rest-book"
      "cmoog.sqlnotebook"
    ];
    # Enable previews and customize configuration
    previews = {
      enable = true;
      previews = {
        web = {
          command = [
            "nodemon"
            "--signal" "SIGHUP"
            "-w" "."
            "-e" "go,html"
            "-x" "go run server.go -http localhost:$PORT"
          ];
          manager = "web";
        };
      };
    };
  };
}
