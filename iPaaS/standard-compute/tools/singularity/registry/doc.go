/*
Singularity Registry Conventions:

Remote image address: <scheme>://<registry-url>/<image-name>/<tag>-<hash>[.sif]
	e.g. https://static.yuansuan.cn/singularity/starccm/v2021.3-abc123xyz.sif

Local image storage: /<storage-base>/<image-name>/<tag-name>/<hash>.sif
	e.g. /data/singularity/starccm/v2021.3/abc123xyz.sif

Local image storage layout:
	/<image-name>
		- _lock						; save the default tag of the image
		/<tag-name>
			- <hash>.sif
			- _lock					; save the hash of the active version of the image
			- [_<hash>.sif]			; the downloading image file
			- [_<hash>.sif.lock]	; whether the image is downloading and the progress will be saved in it
*/

package registry
