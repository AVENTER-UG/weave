with import <nixpkgs> {};

stdenv.mkDerivation {
name = "go-env";

buildInputs = [
    syft
    grype
    docker
    trivy
		stdenv.cc.cc
		libpcap
		docker-credential-helpers
		dbus
		go
		glibc
		glibc.static
		gcc
		flex
		bison
];

SOURCE_DATE_EPOCH = 315532800;
PROJDIR = "${toString ./.}";
S_NETWORK="host";

shellHook = ''
    export LD_LIBRARY_PATH="${pkgs.stdenv.cc.cc.lib}/lib"
    export PATH=/tmp/bin:$PATH
    export GOTMPDIR=/tmp
    export TMPDIR=/tmp

    mkdir /tmp/bin
    '';
}
