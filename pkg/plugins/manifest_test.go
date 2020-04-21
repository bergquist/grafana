package plugins

import (
	"strings"
	"testing"
	"time"

	"github.com/grafana/grafana/pkg/services/sqlstore"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseAccessForManifestKeys(t *testing.T) {
	testDB := sqlstore.InitTestDB(t)
	mv := newManifestVerifier(testDB)

	nowEpoch := time.Now().Unix()
	testKeys := []ManifestKeys{
		{KeyID: "keyid1", PublicKey: "publicKey1", UpdatedAt: nowEpoch, Since: nowEpoch},
		{KeyID: "keyid2", PublicKey: "publicKey2", UpdatedAt: nowEpoch, Since: nowEpoch},
	}

	t.Run("can write keys to database", func(t *testing.T) {
		session := mv.sqlstore.NewSession()
		defer session.Close()
		affected, err := session.Insert(testKeys)
		assert.Nil(t, err)
		assert.Equal(t, affected, int64(len(testKeys)))
	})

	t.Run("can query keys from database", func(t *testing.T) {
		keys, err := mv.getPublicKey("keyid1")
		assert.Nil(t, err)
		assert.Equal(t, keys, testKeys[0].PublicKey)
	})

	t.Run("should return error if key doesn't exist", func(t *testing.T) {
		_, err := mv.getPublicKey("missing keyId")
		assert.NotNil(t, err)
	})
}

func TestManifestParsing(t *testing.T) {
	txt := `-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA512

{
  "plugin": "grafana-googlesheets-datasource",
  "version": "1.0.0-dev",
  "files": {
    "LICENSE": "7df059597099bb7dcf25d2a9aedfaf4465f72d8d",
    "README.md": "08ec6d704b6115bef57710f6d7e866c050cb50ee",
    "gfx_sheets_darwin_amd64": "1b8ae92c6e80e502bb0bf2d0ae9d7223805993ab",
    "gfx_sheets_linux_amd64": "f39e0cc7344d3186b1052e6d356eecaf54d75b49",
    "gfx_sheets_windows_amd64.exe": "c8825dfec512c1c235244f7998ee95182f9968de",
    "module.js": "aaec6f51a995b7b843b843cd14041925274d960d",
    "module.js.LICENSE.txt": "7f822fe9341af8f82ad1b0c69aba957822a377cf",
    "module.js.map": "c5a524f5c4237f6ed6a016d43cd46938efeadb45",
    "plugin.json": "55556b845e91935cc48fae3aa67baf0f22694c3f"
  },
  "time": 1586817677115,
  "keyId": "7e4d0c6a708866e7"
}
-----BEGIN PGP SIGNATURE-----
Version: OpenPGP.js v4.10.1
Comment: https://openpgpjs.org

wqEEARMKAAYFAl6U6o0ACgkQfk0ManCIZuevWAIHSvcxOy1SvvL5gC+HpYyG
VbSsUvF2FsCoXUCTQflK6VdJfSPNzm8YdCdx7gNrBdly6HEs06ZaRp44F/ve
NR7DnB0CCQHO+4FlSPtXFTzNepoc+CytQyDAeOLMLmf2Tqhk2YShk+G/YlVX
74uuP5UXZxwK2YKJovdSknDIU7MhfuvvQIP/og==
=hBea
-----END PGP SIGNATURE-----`

	testDB := sqlstore.InitTestDB(t)
	mv := newManifestVerifier(testDB)

	t.Run("valid manifest", func(t *testing.T) {
		manifest, err := mv.readPluginManifest([]byte(txt))

		assert.Nil(t, err)
		assert.NotNil(t, manifest)
		assert.Equal(t, manifest.Plugin, "grafana-googlesheets-datasource")
	})

	t.Run("invalid manifest", func(t *testing.T) {
		modified := strings.ReplaceAll(txt, "README.md", "xxxxxxxxxx")
		manifest, err := mv.readPluginManifest([]byte(modified))
		assert.NotNil(t, err)
		assert.Nil(t, manifest)
	})
}
