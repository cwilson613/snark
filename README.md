# dkpswitch

`dkpswitch` is a project designed to simplify managing various versions and releases of DKP.

## Usage

Because the DKP releases are gated in private D2iQ repositories, you will need to set a Personal Access Token with the appropriate permissions to read repositories and packages.

You can provide this token via the command line (if `dkpswitch` does not find an environment variable).

Alternately, for a more permanent solution, you can set the environment variable `GITHUB_PERSONAL_ACCESS_TOKEN` in your shell environment file (e.g. `~/.bashrc` or `~/.zshrc`).

### Initialize in directory

If you have used this tool previously to set your version of DKP/Konvoy, a `.dkp` file will have been written in the directory.

`dkpswitch init` will initialize `dkp`/`konvoy` in the present directory.  If no `.dkp` file is found, you will be provided a list of versions to choose from.  If a `.dkp` file is present, `dkpswitch` will initialize the version found within that file (i.e. the version that was previously used)

### Pick version from a list

`dkpswitch list` will reach out and query the GitHub API for all GA releases of DKP, and present them to the user in an interactive list.  Items may be selected by navigating with the arrow keys and pressing the `Enter` key once their desired version is highlighted.

This will download and inflate the associated binaries, and create a symlink from `/usr/local/bin/{binary}` to the selected version of DKP.  

This works for both 1.x (`konvoy` binary) and 2.x (`dkp` binary).
Note that these will both be independently set and managed, so if you run this tool and select a 1.x, followed by a 2.x, both `dkp` and `konvoy` will be symlinked to their respective last selection.

#### Non-GA Releases

Simply include the `all` keyword when listing versions: `dkpswitch list all` to fetch ALL available releases.  

If you do not explicitly specify `all`, only GA releases will be shown.

### Specify a specific version

If you already know the desired DKP version you wish to use, simply skip the list and specify it directly with `dkpswitch <version>`.  To use DKP v2.1.1, either:
- `dkpswitch 2.1.1`
- `dkpswitch v2.1.1`

may be provided.

## Building from source

With the appropriate Go environment in place, the Makefile provides the necessary commands for build, run, etc.