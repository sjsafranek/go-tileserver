package mapnik

// #include <stdlib.h>
// #include "mapnik_c_api.h"
import "C"

import (
	"errors"
	"fmt"
	"os"
	"unsafe"
)

func init() {
	// register default datasources path and fonts path like the python bindings do
	if err := RegisterDatasources(pluginPath); err != nil {
		fmt.Fprintf(os.Stderr, "MAPNIK: %s\n", err)
	}
	if err := RegisterFonts(fontPath); err != nil {
		fmt.Fprintf(os.Stderr, "MAPNIK: %s\n", err)
	}
}

// RegisterDatasources adds path to the Mapnik plugin search path.
func RegisterDatasources(path string) error {
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))
	if C.mapnik_register_datasources(cs) == 0 {
		e := C.GoString(C.mapnik_register_last_error())
		if e != "" {
			return errors.New("registering datasources: " + e)
		}
		return errors.New("error while registering datasources")
	}
	return nil
}

// RegisterDatasources adds path to the Mapnik fonts search path.
func RegisterFonts(path string) error {
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))
	if C.mapnik_register_fonts(cs) == 0 {
		e := C.GoString(C.mapnik_register_last_error())
		if e != "" {
			return errors.New("registering fonts: " + e)
		}
		return errors.New("error while registering fonts")
	}
	return nil
}

type version struct {
	Numeric int
	Major   int
	Minor   int
	Patch   int
	String  string
}

var Version version

func init() {
	Version.Numeric = int(C.mapnik_version)
	Version.Major = int(C.mapnik_version_major)
	Version.Minor = int(C.mapnik_version_minor)
	Version.Patch = int(C.mapnik_version_patch)
	Version.String = fmt.Sprintf("%d.%d.%d", Version.Major, Version.Minor, Version.Patch)
	fmt.Printf("%v\n", Version.String)
}

// func init() {
// 	// register default datasources path and fonts path like the python bindings do
// 	RegisterDatasources(pluginPath)
// 	RegisterFonts(fontPath)
// }

// func Version() string {
// 	return "Mapnik " + C.GoString(C.mapnik_version_string())
// }

// func RegisterDatasources(path string) {
// 	cs := C.CString(path)
// 	defer C.free(unsafe.Pointer(cs))
// 	err := C.CString(path)
// 	defer C.free(unsafe.Pointer(err))
// 	C.mapnik_register_datasources(cs, &err)
// }
//
// func RegisterFonts(path string) {
// 	cs := C.CString(path)
// 	defer C.free(unsafe.Pointer(cs))
// 	err := C.CString(path)
// 	defer C.free(unsafe.Pointer(err))
// 	C.mapnik_register_fonts(cs, &err)
// }

// Point in 2D space
type Coord struct {
	X, Y float64
}

// Projection from one reference system to the other
type Projection struct {
	p *C.struct__mapnik_projection_t
}

func (p *Projection) Free() {
	C.mapnik_projection_free(p.p)
	p.p = nil
}

func (p Projection) Forward(coord Coord) Coord {
	c := C.mapnik_coord_t{C.double(coord.X), C.double(coord.Y)}
	c = C.mapnik_projection_forward(p.p, c)
	return Coord{float64(c.x), float64(c.y)}
}

// Map base type
type Map struct {
	m *C.struct__mapnik_map_t
}

func NewMap(width, height uint32) *Map {
	return &Map{C.mapnik_map(C.uint(width), C.uint(height))}
}

func (m *Map) lastError() error {
	return errors.New("mapnik: " + C.GoString(C.mapnik_map_last_error(m.m)))
}

// Load initializes the map by loading its stylesheet from stylesheetFile
func (m *Map) Load(stylesheetFile string) error {
	cs := C.CString(stylesheetFile)
	defer C.free(unsafe.Pointer(cs))
	if C.mapnik_map_load(m.m, cs) != 0 {
		return m.lastError()
	}
	return nil
}

// LoadString initializes the map not from a file but from a stylesheet
// provided as a string.
func (m *Map) LoadString(stylesheet string) error {
	cs := C.CString(stylesheet)
	defer C.free(unsafe.Pointer(cs))
	if C.mapnik_map_load_string(m.m, cs) != 0 {
		return m.lastError()
	}
	return nil
}

func (m *Map) Resize(width, height uint32) {
	C.mapnik_map_resize(m.m, C.uint(width), C.uint(height))
}

func (m *Map) Free() {
	C.mapnik_map_free(m.m)
	m.m = nil
}

func (m *Map) SRS() string {
	return C.GoString(C.mapnik_map_get_srs(m.m))
}

func (m *Map) SetSRS(srs string) {
	cs := C.CString(srs)
	defer C.free(unsafe.Pointer(cs))
	C.mapnik_map_set_srs(m.m, cs)
}

func (m *Map) ZoomAll() error {
	if C.mapnik_map_zoom_all(m.m) != 0 {
		return m.lastError()
	}
	return nil
}

func (m *Map) ZoomToMinMax(minx, miny, maxx, maxy float64) {
	bbox := C.mapnik_bbox(C.double(minx), C.double(miny), C.double(maxx), C.double(maxy))
	defer C.mapnik_bbox_free(bbox)
	C.mapnik_map_zoom_to_box(m.m, bbox)
}

func (m *Map) RenderToFile(path string) error {
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))
	if C.mapnik_map_render_to_file(m.m, cs) != 0 {
		return m.lastError()
	}
	return nil
}

func (m *Map) RenderToMemoryPng() ([]byte, error) {
	i := C.mapnik_map_render_to_image(m.m)
	if i == nil {
		return nil, m.lastError()
	}
	defer C.mapnik_image_free(i)
	b := C.mapnik_image_to_png_blob(i)
	defer C.mapnik_image_blob_free(b)
	return C.GoBytes(unsafe.Pointer(b.ptr), C.int(b.len)), nil
}

func (m *Map) Projection() Projection {
	p := Projection{}
	p.p = C.mapnik_map_projection(m.m)
	return p
}

func (m *Map) SetBufferSize(s int) {
	C.mapnik_map_set_buffer_size(m.m, C.int(s))
}
