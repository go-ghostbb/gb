package gbview_test

import (
	"context"
	"ghostbb.io/gb/frame/g"
	gbi18n "ghostbb.io/gb/i18n/gb_i18n"
	"ghostbb.io/gb/internal/command"
	gbview "ghostbb.io/gb/os/gb_view"
	gbtest "ghostbb.io/gb/test/gb_test"
	"testing"
)

func Test_Config(t *testing.T) {
	// show error print
	command.Init("-gb.view.errorprint=true")
	gbtest.C(t, func(t *gbtest.T) {
		config := gbview.Config{
			Paths: []string{gbtest.DataPath("config")},
			Data: g.Map{
				"name": "gb",
			},
			DefaultFile: "test.html",
			Delimiters:  []string{"${", "}"},
		}

		view := gbview.New()
		err := view.SetConfig(config)
		t.AssertNil(err)

		view.SetI18n(gbi18n.New())

		str := `hello ${.name},version:${.version}`
		view.Assigns(g.Map{"version": "1.7.0"})
		result, err := view.ParseContent(context.TODO(), str, nil)
		t.AssertNil(err)
		t.Assert(result, "hello gb,version:1.7.0")

		result, err = view.ParseDefault(context.TODO())
		t.AssertNil(err)
		t.Assert(result, "name:gb")

		t.Assert(view.GetDefaultFile(), "test.html")
	})
	// SetConfig path fail: notexist
	gbtest.C(t, func(t *gbtest.T) {
		config := gbview.Config{
			Paths: []string{"notexist", gbtest.DataPath("config/test.html")},
			Data: g.Map{
				"name": "gb",
			},
			DefaultFile: "test.html",
			Delimiters:  []string{"${", "}"},
		}

		view := gbview.New()
		err := view.SetConfig(config)
		t.AssertNE(err, nil)
	})
	// SetConfig path fail: set file path
	gbtest.C(t, func(t *gbtest.T) {
		config := gbview.Config{
			Paths: []string{gbtest.DataPath("config/test.html")},
			Data: g.Map{
				"name": "gb",
			},
			DefaultFile: "test.html",
			Delimiters:  []string{"${", "}"},
		}

		view := gbview.New()
		err := view.SetConfig(config)
		t.AssertNE(err, nil)
	})
}

func Test_ConfigWithMap(t *testing.T) {
	gbtest.C(t, func(t *gbtest.T) {
		view := gbview.New()
		err := view.SetConfigWithMap(g.Map{
			"Paths":       []string{gbtest.DataPath("config")},
			"DefaultFile": "test.html",
			"Delimiters":  []string{"${", "}"},
			"Data": g.Map{
				"name": "gb",
			},
		})
		t.AssertNil(err)

		str := `hello ${.name},version:${.version}`
		view.Assigns(g.Map{"version": "1.7.0"})
		result, err := view.ParseContent(context.TODO(), str, nil)
		t.AssertNil(err)
		t.Assert(result, "hello gb,version:1.7.0")

		result, err = view.ParseDefault(context.TODO())
		t.AssertNil(err)
		t.Assert(result, "name:gb")
	})
	// path as paths
	gbtest.C(t, func(t *gbtest.T) {
		view := gbview.New()
		err := view.SetConfigWithMap(g.Map{
			"Path":        gbtest.DataPath("config"),
			"DefaultFile": "test.html",
			"Delimiters":  []string{"${", "}"},
			"Data": g.Map{
				"name": "gb",
			},
		})
		t.AssertNil(err)

		str := `hello ${.name},version:${.version}`
		view.Assigns(g.Map{"version": "1.7.0"})
		result, err := view.ParseContent(context.TODO(), str, nil)
		t.AssertNil(err)
		t.Assert(result, "hello gb,version:1.7.0")

		result, err = view.ParseDefault(context.TODO())
		t.AssertNil(err)
		t.Assert(result, "name:gb")
	})
	// path as paths
	gbtest.C(t, func(t *gbtest.T) {
		view := gbview.New()
		err := view.SetConfigWithMap(g.Map{
			"Path":        []string{gbtest.DataPath("config")},
			"DefaultFile": "test.html",
			"Delimiters":  []string{"${", "}"},
			"Data": g.Map{
				"name": "gb",
			},
		})
		t.AssertNil(err)

		str := `hello ${.name},version:${.version}`
		view.Assigns(g.Map{"version": "1.7.0"})
		result, err := view.ParseContent(context.TODO(), str, nil)
		t.AssertNil(err)
		t.Assert(result, "hello gb,version:1.7.0")

		result, err = view.ParseDefault(context.TODO())
		t.AssertNil(err)
		t.Assert(result, "name:gb")
	})
	// map is nil
	gbtest.C(t, func(t *gbtest.T) {
		view := gbview.New()
		err := view.SetConfigWithMap(nil)
		t.AssertNE(err, nil)
	})
}
