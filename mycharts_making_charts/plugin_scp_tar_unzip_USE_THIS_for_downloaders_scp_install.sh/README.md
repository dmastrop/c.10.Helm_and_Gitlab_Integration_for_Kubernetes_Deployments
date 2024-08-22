This project is one of the labs presented in my course: Helm - The Kubernetes Package Manager Hands-On. 
It is a Helm plugin that is used to automatically package and upload a chart to a remote server over SSH
Use the binary as 
```bash
./helmscp push /path/to/chart scp://username@hostname[:port]/path/to/repo
```
It uses secure copy (scp) behind the scenes to upload charts to the remote server. It assumes that you are using `~/.ssh/id_rsa` as the private key. 
But you can override this value as follows:
```bash
export SCP_KEY=~/pth/to/key
```
The remote repo is a directory that you have sufficient permissions to upload and modify files. For example, your home directory.
## As a Helm plugin
Install the plugin by running 
```bash
helm plugin install https://github.com/abohmeed/helmscpplugin
```
Alternatively you can clone the repo and run 
```bash
helm plugin install /path/to/cloned/repo
```
It can be used as a downloader plugin as follows:
```bash
helm repo add helmscprepo scp://username@hostname[:port]/path/to/repo
helm repo update
helm scp push /path/to/chart scp://username@hostname[:port]/path/to/repo
```
Then you can use it like any other repo.
### Note
Helm 3 must be installed on the remote server as the plugin uses it to reindex the charts once a new chart is added
### Deleting charts
You can also delete charts from the remote repo as follows:
```bash
helm scp delete scp://username@hostname[:port]/path/to/chart.tgz
``` 