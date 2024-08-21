#!/usr/bin/env bash
set -e

function usage(){
    echo """Usage:
-k string
    The SSH key (default "/Users/ahmedelfakharany/.ssh/id_rsa")
-l string
        Path to the chart directory
-p string
        The remote server port (default "22")
-r string
        Path to the remote directory
-s string
        The hostname or IP address
-u string
        The remote server username"""
}

while getopts u:s:k:l:p:r: flag
do
    case "${flag}" in
        u) username=${OPTARG};;
        s) hostname=${OPTARG};;
        k) key=${OPTARG};;
        l) chart_dir=${OPTARG};;
        p) port=${OPTARG};;
        r) remote_dir=${OPTARG};;
    esac
done
if [ -z $username ] || [ -z $hostname ] || [ -z $chart_dir ] || [ -z $remote_dir ]; then
    usage
    exit 1
fi
if [ -z $port ]; then
    port="22"
fi
if [ -z $key ]; then
    key="~/.ssh/id_rsa"
fi
remote_dir=${remote_dir%"/"}
if [ -z $HELM_BIN ]; then
    HELM_BIN=$(which helm)
fi
echo "Packaging chart from ${chart_dir}"
cmd_output=$($HELM_BIN package $chart_dir)
chart_name=$(echo $cmd_output | cut -d ":" -f 2 | tr -d "[:space:]")
chart_base_name=$(basename $chart_name)
scp_cmd=$(which scp || exit 1)
echo "Uploading ${chart_base_name} to ${remote_dir} at ${username}@${hostname}:${port}"
scp -i ${key} -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -P ${port} ${chart_name} ${username}@${hostname}:${remote_dir}/${chart_base_name} 1>/dev/null || exit 1
rm_cmd=$(which rm)
echo "Cleaning up"
$rm_cmd -f $chart_name
echo "Success!"