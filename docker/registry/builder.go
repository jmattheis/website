package dockerregistry

import (
	"encoding/json"

	"github.com/distribution/distribution/v3"
	"github.com/distribution/distribution/v3/manifest/ocischema"
	"github.com/opencontainers/go-digest"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type Entry struct {
	Digest    digest.Digest
	MediaType string
	Content   []byte
}

type Setter func(string, Entry)

func RegisterLayers(set Setter) {
	var bManifests ocischema.DeserializedImageIndex
	check(bManifests.UnmarshalJSON(read("manifest.json")))

	for _, maniDesc := range bManifests.Manifests {
		m := ocischema.Manifest{}
		check(json.Unmarshal(read(maniDesc.Digest.Encoded()+".manifest.json"), &m))

		for _, layer := range m.Layers {
			registerRaw(set, "application/octet-stream", read(layer.Digest.Encoded()))
		}
	}
}

func DockerStdout(set Setter, content string) Entry {
	var bManifests ocischema.DeserializedImageIndex
	check(bManifests.UnmarshalJSON(read("manifest.json")))

	var newManifestDescriptions []distribution.Descriptor
	for _, maniDesc := range bManifests.Manifests {
		m := ocischema.Manifest{}
		check(json.Unmarshal(read(maniDesc.Digest.Encoded()+".manifest.json"), &m))
		m.Annotations = nil

		cfg := v1.Image{}
		check(json.Unmarshal(read(m.Config.Digest.Encoded()), &cfg))
		cfg.Config.Cmd = []string{"echo", content}
		m.Config.Digest, m.Config.Size, _ = registerJson(set, "application/octet-stream", cfg)

		maniDesc.Digest, maniDesc.Size, _ = registerJson(set, v1.MediaTypeImageManifest, m)

		newManifestDescriptions = append(newManifestDescriptions, maniDesc)
	}
	index, err := ocischema.FromDescriptors(newManifestDescriptions, nil)
	check(err)

	digest, _, indexBytes := registerJson(set, v1.MediaTypeImageIndex, index)

	return Entry{
		Digest:    digest,
		MediaType: v1.MediaTypeImageIndex,
		Content:   indexBytes,
	}
}

func registerJson(set Setter, mediaType string, value any) (digest.Digest, int64, []byte) {
	b, err := json.Marshal(value)
	check(err)
	return registerRaw(set, mediaType, b)
}

func registerRaw(set Setter, mediaType string, b []byte) (digest.Digest, int64, []byte) {
	d := digest.FromBytes(b)
	set(d.String(), Entry{
		Digest:    d,
		MediaType: mediaType,
		Content:   b,
	})
	return d, int64(len(b)), b
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
