package mod

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestModInit(t *testing.T) {
	fs := afero.NewOsFs()
	gomod := &goModules{}

	// assumes the test folder (cwd) is not a go module folder
	removeFile(t, fs, "go.sum")
	removeFile(t, fs, "go.mod")

	err := gomod.Init("test")
	assert.NoError(t, err)

	removeFile(t, fs, "go.sum")
	removeFile(t, fs, "go.mod")
}

func TestModInitAlreadyExists(t *testing.T) {
	fs := afero.NewOsFs()
	gomod := &goModules{}

	// assumes the test folder (cwd) is not a go module folder
	removeFile(t, fs, "go.sum")
	removeFile(t, fs, "go.mod")

	err := gomod.Init("test")
	assert.NoError(t, err)

	err = gomod.Init("test")
	assert.Error(t, err)

	removeFile(t, fs, "go.sum")
	removeFile(t, fs, "go.mod")
}

func TestGoModulesGet(t *testing.T) {
	gomod := &goModules{}
	testMods := Modules{}

	mod, err := gomod.Get(RemoteDepsFile, "", &testMods)
	assert.Nil(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)

	mod, err = gomod.Get(RemoteDepsFile, MasterBranch, &testMods)
	assert.Nil(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)

	mod, err = gomod.Get(RemoteDepsFile, "v0.0.1", &testMods)
	assert.Nil(t, err)
	assert.Equal(t, RemoteRepo, mod.Name)
	assert.Equal(t, "v0.0.1", mod.Version)

	mod, err = gomod.Get("github.com/anz-bank/wrongpath", "", &testMods)
	assert.Error(t, err)
	assert.Nil(t, mod)
}

func TestGoModulesFind(t *testing.T) {
	gomod := &goModules{}
	testMods := Modules{}
	local := &Module{Name: "local"}
	mod2 := &Module{Name: "remote", Version: "v0.2.0"}
	testMods.Add(local)
	testMods.Add(mod2)

	assert.Equal(t, local, gomod.Find("local/filename", "", &testMods))
	assert.Equal(t, local, gomod.Find("local/filename2", "", &testMods))
	assert.Equal(t, local, gomod.Find(".//local/filename", "", &testMods))
	assert.Equal(t, local, gomod.Find("local", "", &testMods))
	assert.Nil(t, gomod.Find("local2/filename", "", &testMods))

	assert.Equal(t, local, gomod.Find("local/filename", MasterBranch, &testMods))
	assert.Equal(t, local, gomod.Find("local/filename", "v0.0.1", &testMods))

	assert.Equal(t, mod2, gomod.Find("remote/filename", "v0.2.0", &testMods))
	assert.Nil(t, gomod.Find("remote/filename", "v1.0.0", &testMods))
}

func removeGomodFile(t *testing.T, fs afero.Fs) {
	removeFile(t, fs, "go.mod")
	removeFile(t, fs, "go.sum")
}

func createGomodFile(t *testing.T, fs afero.Fs) {
	gomod, err := fs.Create("go.mod")
	assert.NoError(t, err)
	defer gomod.Close()
	_, err = gomod.WriteString("module github.com/anz-bank/pkg/mod")
	assert.NoError(t, err)
	err = gomod.Sync()
	assert.NoError(t, err)
}

func removeFile(t *testing.T, fs afero.Fs, file string) {
	exists, err := afero.Exists(fs, file)
	assert.NoError(t, err)
	if exists {
		err = fs.Remove(file)
		assert.NoError(t, err)
	}
}
