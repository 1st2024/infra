# E2B API

Monorepo with backend services for handling VM sessions, environment pipelines, and the API for them.

## Deployment

Run `make version` and commit the changes the command generates.
Then push to `main` to deploy these changes. Changed packages will be automatically deployed.

**If the deployment fails don't run the previous commands again, just fix the error and push to `main`.**

### Errors in the env-instance-task-driver
If you are not developing env-instance-task-driver module on a unix machine you won't be able to complile the module because you will miss constants like "netlink.SCOPE_UNIVERSE".
Use the ubuntu devcontainer to develop this module.

### Resizing fc-envs disk (XFS)
After modifying https://console.cloud.google.com/compute/disksDetail/zones/us-central1-a/disks/fc-envs disk size in the dashboard you need to resize the partition by running the `xfs_growfs -d /dev/sdb`.

### Creating XFS filesystem
```
mkfs.xfs  /dev/sdb
```

> You need to remount the fs after formatting it.

### Convert disk to XFS
- https://access.redhat.com/discussions/6134431


### FC API Client
For generating the FC client use go swagger (https://goswagger.io/install.html).
Right now you need to manually check and update the [firecracker yml](./internal/client/firecracker.yml) version. It should match the version of the firecracker binary you are using specified in the [terraform file](../cluster-disk-image/variables.pkr.hcl).

### Remote development setup (VSCode)


1. Run `gcloud compute config-ssh`
2. Install "Remote SSH" extension
3. Go to the "Tunnels/SSH" tab in "Remote Explorer" extension
4. Click on config button in the SSH section
5. Select the (usually first) option: `/Users/<username>/.ssh/config`
6. Refresh the list of SSHs (you may need to reload VSCode)
7. Click on the "Connect in new window" next to the orch-client instance. This will open the remote VSCode
8. Pull this repo from GitHub via the Source control VSCode tab
*9. You may need to install `sudo apt install make` and `sudo snap install go --classic` until we add this to the default instance image*
10. The Go debugger should be attachable to the remote instance too
11. Configure git identifiers
- `git config --global user.email "MY_NAME@example.com"`
- `git config --global user.name "NAME"`


### Shell script best practices

https://sharats.me/posts/shell-script-best-practices/
