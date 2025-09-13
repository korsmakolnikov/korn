package templates

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"
)

type UpsertableMock struct {
	Counter int
}

func (u *UpsertableMock) touch() error {
	u.Counter++
	return nil
}

func (u *UpsertableMock) String() string {
	return fmt.Sprint(u.Counter)
}

const (
	BASEPATH = "/home/user/.config/kornvim/builds"
)

func TestNewInitTemplate(t *testing.T) {
	templatable := NewInit(BASEPATH, "testBuildPath")
	err := templatable.Parse()
	if err != nil {
		t.Fatalf("Something wrong creating the template: %+v", err)
	}
}

func TestUpsertFolder(t *testing.T) {
	templatable := NewInit(BASEPATH, "testBuildPath")
	err := templatable.Parse()

	u := &UpsertableMock{Counter: 0}
	err = templatable.Prepare(u)
	if err != nil {
		t.Fatalf("Something wrong preparing the template folder: %+v", err)
	}

	if u.Counter != 1 {
		t.Fatalf("You asserted UpsertableIO.create was called once but it has been called %d times", u.Counter)
	}
}

func TestInitToPath(t *testing.T) {
	templatable := NewInit("/home/user/.config/kornvim/builds", "testBuildPath")
	_ = templatable.Parse()
	actual := templatable.ToPath()
	if "/home/user/.config/kornvim/builds/testBuildPath/lua/init.lua" != actual.String() {
		t.Fatalf("the output path hasn't be set correctly [%s]", actual)
	}

}

func TestInitPointToRightBuildFolder(t *testing.T) {
	wantedBuildName := "testBuildName"
	wantedBuildPath := fmt.Sprintf("/home/user/.config/kornvim/builds/%s", wantedBuildName)
	templatable := NewInit(BASEPATH, wantedBuildName)
	err := templatable.Parse()

	var buf bytes.Buffer
	err = templatable.Execute(&buf)
	if err != nil {
		t.Fatalf("Something executing the template folder: %+v", err)
	}
	actual := buf.String()
	re1 := regexp.MustCompile(fmt.Sprintf(`lazy_root_dir = vim\.fn\.getcwd\(\) \.\. "/" \.\. "%s" \.\. "/lazy"`, wantedBuildName))
	re2 := regexp.MustCompile(fmt.Sprintf(`lazy_root_dir = "%s" \.\. "/lazy"`, wantedBuildPath))

	for _, re := range []*regexp.Regexp{re1, re2} {
		if !re.MatchString(actual) {
			t.Fatalf("the template output doesn't not match with wanted text [%s]", re.String())
		}
	}
}

func TestNewPackageTemplate(t *testing.T) {
	templatable := NewPackage(BASEPATH, "testBuildPath", "korsmakolnikov/kornvim_configurator")
	err := templatable.Parse()
	if err != nil {
		t.Fatalf("Something wrong creating the template: %+v", err)
	}
}

func TestPackageUpsertFolder(t *testing.T) {
	templatable := NewPackage(BASEPATH, "testBuildPath", "korsmakolnikov/kornvim_configurator")
	err := templatable.Parse()

	u := &UpsertableMock{Counter: 0}
	err = templatable.Prepare(u)
	if err != nil {
		t.Fatalf("Something wrong preparing the template folder: %+v", err)
	}

	if u.Counter != 1 {
		t.Fatalf("You asserted UpsertableIO.create was called once but it has been called %d times", u.Counter)
	}
}

func TestPackageToPath(t *testing.T) {
	templatable := NewPackage("/home/user/.config/kornvim/builds", "testBuildPath", "korsmakolnikov/kornvim_configurator")
	_ = templatable.Parse()
	actual := templatable.ToPath()
	if "/home/user/.config/kornvim/builds/testBuildPath/lua/package.lua" != actual.String() {
		t.Fatalf("the output path hasn't be set correctly [%s]", actual)
	}

}

func TestPackageIncludeTheConfiguratorPlugin(t *testing.T) {
	wantedBuildName := "testBuildName"
	templatable := NewPackage(BASEPATH, wantedBuildName, "korsmakolnikov/kornvim_configurator")
	err := templatable.Parse()

	var buf bytes.Buffer
	err = templatable.Execute(&buf)
	if err != nil {
		t.Fatalf("Something executing the template folder: %+v", err)
	}
	actual := buf.String()
	re1 := regexp.MustCompile(`"korsmakolnikov/kornvim_configurator",`)
	re2 := regexp.MustCompile(`require\('kornvim_configurator'\)\.setup\(\)`)

	for _, re := range []*regexp.Regexp{re1, re2} {
		if !re.MatchString(actual) {
			t.Fatalf("the template output doesn't not match with wanted text [%s]", re.String())
		}
	}
}
