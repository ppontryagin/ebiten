package opengl

// #cgo LDFLAGS: -framework OpenGL
//
// #include <stdlib.h>
// #include <OpenGL/gl.h>
import "C"
import (
	"github.com/hajimehoshi/go-ebiten/graphics"
	"github.com/hajimehoshi/go-ebiten/graphics/matrix"
	"github.com/hajimehoshi/go-ebiten/graphics/opengl/offscreen"
	"math"
)

type Context struct {
	screenId    graphics.RenderTargetId
	ids         *ids
	offscreen   *offscreen.Offscreen
	screenScale int
}

func newContext(ids *ids, screenWidth, screenHeight, screenScale int) *Context {
	context := &Context{
		ids:         ids,
		offscreen:   offscreen.New(screenWidth, screenHeight, screenScale),
		screenScale: screenScale,
	}

	var err error
	context.screenId, err = ids.CreateRenderTarget(
		screenWidth, screenHeight, graphics.FilterNearest)
	if err != nil {
		panic("initializing the offscreen failed: " + err.Error())
	}
	context.Clear()

	return context
}

func (context *Context) update(draw func(graphics.Context)) {
	C.glEnable(C.GL_TEXTURE_2D)
	C.glEnable(C.GL_BLEND)

	context.ResetOffscreen()
	context.Clear()

	draw(context)

	C.glFlush()
	context.offscreen.SetMainFramebuffer()
	context.Clear()

	scale := float64(context.screenScale)
	geometryMatrix := matrix.IdentityGeometry()
	geometryMatrix.Scale(scale, scale)
	context.DrawRenderTarget(context.screenId,
		geometryMatrix, matrix.IdentityColor())
	C.glFlush()
}

func (context *Context) Clear() {
	context.Fill(0, 0, 0)
}

func (context *Context) Fill(r, g, b uint8) {
	const max = float64(math.MaxUint8)
	C.glClearColor(
		C.GLclampf(float64(r)/max),
		C.GLclampf(float64(g)/max),
		C.GLclampf(float64(b)/max),
		1)
	C.glClear(C.GL_COLOR_BUFFER_BIT)
}

func (context *Context) DrawTexture(
	id graphics.TextureId,
	geometryMatrix matrix.Geometry, colorMatrix matrix.Color) {
	tex := context.ids.TextureAt(id)
	context.offscreen.DrawTexture(tex, geometryMatrix, colorMatrix)
}

func (context *Context) DrawRenderTarget(
	id graphics.RenderTargetId,
	geometryMatrix matrix.Geometry, colorMatrix matrix.Color) {
	context.DrawTexture(context.ids.ToTexture(id), geometryMatrix, colorMatrix)
}

func (context *Context) DrawTextureParts(
	id graphics.TextureId, parts []graphics.TexturePart,
	geometryMatrix matrix.Geometry, colorMatrix matrix.Color) {
	tex := context.ids.TextureAt(id)
	context.offscreen.DrawTextureParts(tex, parts, geometryMatrix, colorMatrix)
}

func (context *Context) DrawRenderTargetParts(
	id graphics.RenderTargetId, parts []graphics.TexturePart,
	geometryMatrix matrix.Geometry, colorMatrix matrix.Color) {
	context.DrawTextureParts(context.ids.ToTexture(id), parts, geometryMatrix, colorMatrix)
}

func (context *Context) ResetOffscreen() {
	context.SetOffscreen(context.screenId)
}

func (context *Context) SetOffscreen(renderTargetId graphics.RenderTargetId) {
	renderTarget := context.ids.RenderTargetAt(renderTargetId)
	context.offscreen.Set(renderTarget)
}