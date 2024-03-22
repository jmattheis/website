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

func DockerStdout(content string) map[string]Entry {
	var bManifests ocischema.DeserializedImageIndex
	check(bManifests.UnmarshalJSON(read("manifest.json")))

	store := map[string]Entry{}

	var newManifestDescriptions []distribution.Descriptor
	for _, maniDesc := range bManifests.Manifests {
		m := ocischema.Manifest{}
		check(json.Unmarshal(read(maniDesc.Digest.Encoded()+".manifest.json"), &m))
		m.Annotations = nil

		cfg := v1.Image{}
		check(json.Unmarshal(read(m.Config.Digest.Encoded()), &cfg))
		cfg.Config.Cmd = []string{"echo", content}
		m.Config.Digest, m.Config.Size = registerJson(store, "application/octet-stream", cfg)

		for _, layer := range m.Layers {
			registerRaw(store, "application/octet-stream", read(layer.Digest.Encoded()))
		}

		maniDesc.Digest, maniDesc.Size = registerJson(store, v1.MediaTypeImageManifest, m)

		newManifestDescriptions = append(newManifestDescriptions, maniDesc)
	}
	index, err := ocischema.FromDescriptors(newManifestDescriptions, nil)
	check(err)

	digest, _ := registerJson(store, v1.MediaTypeImageIndex, index)

	store["latest"] = Entry{
		Digest:    digest,
		MediaType: v1.MediaTypeImageIndex,
		Content:   store[digest.String()].Content,
	}
	return store
}

func registerJson(store map[string]Entry, mediaType string, value any) (digest.Digest, int64) {
	b, err := json.Marshal(value)
	check(err)
	return registerRaw(store, mediaType, b)
}

func registerRaw(store map[string]Entry, mediaType string, b []byte) (digest.Digest, int64) {
	d := digest.FromBytes(b)
	store[d.String()] = Entry{
		Digest:    d,
		MediaType: mediaType,
		Content:   b,
	}
	return d, int64(len(b))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
