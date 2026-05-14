This directory contains some wrapper scripts to help us run the protocol buffers
compiler and its Go-specific plugins without having to first preinstall them
globally on the host system, and so that we can force depending on specific
versions of those tools.

The `generate.sh` scripts in the sibling `tfplugin5` and `tfplugin6` directories
add this directory to `PATH` before running the protocol buffers compiler, to
ensure that they'll use these wrappers instead of versions that might be
installed elsewhere on the host system.
