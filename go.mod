module github.com/lmikolajczak/wms-tiles-downloader

go 1.18

require (
	github.com/jarcoal/httpmock v1.2.0
	github.com/schollz/progressbar/v3 v3.8.7
	github.com/spf13/cobra v1.5.0
	github.com/stretchr/testify v1.8.0
)

require (
	github.com/ctessum/geom v0.2.12 // indirect
	github.com/ctessum/polyclip-go v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
	github.com/gonum/floats v0.0.0-20181209220543-c233463c7e82 // indirect
	github.com/gonum/internal v0.0.0-20181124074243-f884aa714029 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/kellydunn/golang-geo v0.7.0 // indirect
	github.com/kylelemons/go-gypsy v1.0.0 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/paulmach/go.geojson v1.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/ziutek/mymysql v1.5.4 // indirect
	golang.org/x/crypto v0.0.0-20220131195533-30dcbda58838 // indirect
	golang.org/x/sys v0.0.0-20220128215802-99c3d69c2c27 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	gonum.org/v1/gonum v0.9.3 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	// Wrongly published versions
	v2.0.0+incompatible
	v1.0.0
)
