package presdo

import (
    "testing"
    "path"
)

func TestPathIndexPath(t *testing.T) {
    var expect, result, test string

    setDefaultValues()

    test = "/"
    expect = "markdown"

    if result = paths.IndexPath(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }

    test = "/index" + websiteConfig.Ext
    expect = "markdown"

    if result = paths.IndexPath(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }

    test = "/foo"
    expect = "markdown/foo"

    if result = paths.IndexPath(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }

    test = "/foo/index" + websiteConfig.Ext
    expect = "markdown/foo"

    if result = paths.IndexPath(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }
}

func TestPathPublic(t *testing.T) {
    var expect, result, test string

    setDefaultValues()

    // Test "/"
    // Need to be "{publicDir}"
    test = "/"
    expect = path.Clean(websiteConfig.PublicDir)

    if result = paths.Public(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }

    // Test "/foo"
    // Need to be "{publicDir}/foo"
    test = "/foo"
    expect = websiteConfig.PublicDir + "foo"

    if result = paths.Public(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }

    // Test "/foo/bar/"
    // Need to be "{publicDir}/foo/bar"
    test = "/foo/bar"
    expect = websiteConfig.PublicDir + "foo/bar"

    if result = paths.Public(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }

    // Test "/foo/bar/index.htm"
    // Need to be "{publicDir}/foo/bar/index.htm"
    test = "/foo/bar/index.htm"
    expect = websiteConfig.PublicDir + "foo/bar/index.htm"

    if result = paths.Public(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }
}


func TestPathCache(t *testing.T) {
    var expect, result, test string

    setDefaultValues()

    // Test "/index.htm"
    // Need to be "cache/index.htm"
    test = "/index.htm"
    expect = "cache/index.htm"

    if result = paths.Cache(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }

    // Test "/foo/bar.html"
    // Need to be "cache/foo/bar.html"
    test = "/foo/bar.html"
    expect = "cache/foo/bar.html"

    if result = paths.Cache(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }
}

func TestPathMarkdown(t *testing.T) {
    var expect, result, test string

    setDefaultValues()

    // Test "/"
    // Need to be "{publicDir}"
    test = "/index.htm"
    expect = "markdown/index.md"

    if result = paths.Markdown(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }

    // Test "/foo"
    // Need to be "{publicDir}/foo"
    test = "/foo/bar.html"
    expect = "markdown/foo/bar.md"

    if result = paths.Markdown(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }
}

func TestPathRequest(t *testing.T) {
    var expect, result, test string

    setDefaultValues()

    // Test "/"
    // Need to be "/index.EXT"
    test = "/"
    expect = "/index" + websiteConfig.Ext

    if result = paths.Request(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }

    // Test "/foo"
    // Need to be "/foo/index.EXT"
    test = "/foo"
    expect = "/foo/index" + websiteConfig.Ext

    if result = paths.Request(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }

    // Test "foo"
    // Need to be "foo/index.EXT"
    test = "foo"
    expect = "foo/index" + websiteConfig.Ext

    if result = paths.Request(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }

    // Test "/bar/"
    // Need to be "/bar/index.EXT"
    test = "/bar/"
    expect = "/bar/index" + websiteConfig.Ext

    if result = paths.Request(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }

    // Test "foo/bar/img.gif"
    // Need to be "foo/bar/img.gif"
    test = "foo/bar/img.gif"
    expect = "foo/bar/img.gif"

    if result = paths.Request(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }
}

func TestPathPage(t *testing.T) {
    var expect, result, test string

    setDefaultValues()

    test = "markdown/foo.md"
    expect = "foo" + websiteConfig.Ext

    if result = paths.Page(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }


    test = "markdown/foo/bar.md"
    expect = "foo/bar" + websiteConfig.Ext

    if result = paths.Page(test); expect != result {
        t.Errorf("Path need to be %s, return %s", expect, result )
    }
}