# VueTorrent manager

Simple utility for keeping vuetorrent up-to-date

## Usage

### Prerequisites
vt-manager uses github api to get information about releases. For this reason you need to obtain github's fine-grained access token with **Repository permission: Contents (read-only)**. [Link to GitHub docs](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token)

You can see available commands by typing
```sh
./vt-manager --help
```

Right now vt-mager support following commands
 - info (prints version of installed vuetorrent)
 - install (get latest or specific version)
 - list (prints all available version for install)
 - revision (prints revision of vt-manager)

### Install new version
This commang will download the latest `vuetorent.zip` from github and unzip it to specified directory (if direcory already exists it will replace all content)
```sh
export GITHUB_ACCESS_TOKEN=xxx
./bin/vt-manager install --dir=./vuetorrent --api-key=$GITHUB_ACCESS_TOKEN
```
To download specific version just add `--version=2.3.0` parameter
```sh 
./bin/vt-manager install --dir=./vuetorrent --api-key=$GITHUB_ACCESS_TOKEN --version=2.3.0
```

### Get installed vuetorrent version
```sh
./bin/vt-manager info --dir=./vuetorrent
```

### List available vuetorrent versions for install
```sh
./bin/vt-manager list --api-key=$GITHUB_ACCESS_TOKEN
```

### Get vt-manger revision
```sh
./bin/vt-manager revision
```

## Build from source

```sh
# Clone this reposiotry and then...
make build

# In the `bin/` directory you should find binaries
# tree bin 
# bin
# ├── linux_amd64
# │   └── vt-manager-amd64
# └── vt-manager
```

