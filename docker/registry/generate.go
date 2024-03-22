//go:generate rm -rf ./busybox
//go:generate skopeo copy --multi-arch=all docker://index.docker.io/library/busybox:latest dir:./busybox

package dockerregistry
