# `dropbox-oneway`

One-way sync of content from Dropbox: download content from Dropbox to a local folder and keep it up to date as files are modified in Dropbox.

## Installation

### Docker

Docker is the preferred way to use `dropbox-oneway`. The [`tantalic/dropbox-oneway`][dockerhub] image can be pulled from Docker Hub as follows:

```shell
docker pull tantalic/dropbox-oneway:v0.1.1
```

### macOS 

Releases for macOS are available through the [Homebrew package manager][homebrew]. If you use Homebrew installation is as easy as:

```
brew tap tantalic/tap
brew install dropbox-oneway
```

### Binary Releases

[Official binaries][releases] are provided for the following platforms:

- Mac OS X (x64)
- Linux (x86, x64, ARMv6, ARMv7, ARM64)
- FreeBSD (x86, x64, ARMv6, ARMv7)
- NetBSD (x86, x64, ARMv6, ARMv7)
- DragonFly BSD (x64)
- Windows (x86, x64)

`dropbox-oneway` is a single executable with no dependencies. Installation is as simple as [downloading the binary for your platform][releases] and placing it in your `PATH`. 

## Usage

All configuration is managed through environment variables (the ability to set options via command line flags and/or arguments may be added in a future release).

| Environment Variable |                                                                          Description                                                                          |
|----------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `DROPBOX_TOKEN`      | The Dropbox API token to use for authentication. You will need to create an app and generate this token from the [Dropbox developer console][db-dev-console]. |
| `DROPBOX_DIRECTORY`  | The path of the content in Dropbox to sync.                                                                                                                   |
| `LOCAL_DIRECTORY`    | The path to download content into. By default this will be the local directory.                                                                               |


### Command Line

To run from the command line set the environment variables before invoking the application:

```bash
DROPBOX_DIRECTORY=/Photos LOCAL_DIRECTORY=~/Pictures/Dropbox DROPBOX_TOKEN=xxxxxxxxxxxxxxxxxxxxxx_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx dropbox-oneway
```

### Docker

To use the Docker image mount the directory to downland contest into as the `/content` volume using the `--volume` flag and set the `DROPBOX_DIRECTORY` and `DROPBOX_TOKEN` variables using the `--env` flag.

```bash
docker run \
    --volume /Users/me/Pictures/Dropbox/output:/dropbox \
    --env DROPBOX_DIRECTORY=/Photos \
    --env DROPBOX_TOKEN=xxxxxxxxxxxxxxxxxxxxxx_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
    tantalic/dropbox-oneway:0.1.1
```

### Kubernetes

The Docker image can be used as a sidecar container to provide content, configuration or other data from dropbox to another container in the pod. Composing a pods in such a way makes it simple to add powerful Dropbox features to an existing application without requiring it to have any knowledge of Dropbox. 

For example you can create a site that automatically updates itself from a folder of Markdown files by combining `dropbox-oneway` with `servemd` as follows (similarly you could combine a folder of HTML files with a standard nginx container):

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: dropbox
type: Opaque
data:
  token: eHh4eHh4eHh4eHh4eHh4eHh4eHh4eF94eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eA==
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: sample-servemd-dropbox
spec:
  replicas: 1
  template: 
    metadata:
      labels:
        name: sample-servemd-dropbox
    spec:
      containers:
        - name: dropbox
          image: tantalic/dropbox-oneway:v0.1.1
          env: 
          - name: DROPBOX_DIRECTORY
            value: /path/to/contents
          - name: DROPBOX_TOKEN
            valueFrom:
              secretKeyRef:
                name: dropbox
                key: token
          volumeMounts:
            - name: dropbox
              mountPath: /dropbox
        - name: servemd
          image: tantalic/servemd:v0.6.0
          ports:
            - containerPort: 3000
              name: http
          env: 
            - name: DOCUMENT_ROOT
              value: /content/path/to/contents
          volumeMounts:
            - name: dropbox
              mountPath: /content
      volumes:
        - name: dropbox
          emptyDir: {}
```

## Contributing

To request a new feature or report a bug open an issue in Github for discussion prior to submitting a pull request. If you're willing to contribute changes/fixes mention that in the issue report or discussion. If you're reporting a bug be as specific as possible about the conditions causing the bug. If possible, attach sample code/files illustrating the issue.

One of the primary objectives of this project is to be simple as possible and composable with other applications (and containers). If a proposed change is not a common use case or introduces undue complexity it may not be accepted.

## License

This is provided under the terms of the [MIT license][license].

[dockerhub]: https://hub.docker.com/r/tantalic/dropbox-oneway/
[db-dev-console]: https://www.dropbox.com/developers/apps
[servemd]: https://github.com/tantalic/servemd
[homebrew]: http://brew.sh/
[releases]: https://github.com/tantalic/dropbox-oneway/releases/latest
[license]: ./LICENSE.txt

