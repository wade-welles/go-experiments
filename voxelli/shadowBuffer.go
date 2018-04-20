package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type ShadowBuffer struct {
	shadowBuffer  uint32
	shadowTexture uint32

	Width  int32
	Height int32
}

func NewShadowBuffer() *ShadowBuffer {
	var currentDrawBuffer int32
	gl.GetIntegerv(gl.DRAW_BUFFER, &currentDrawBuffer)

	var buffer ShadowBuffer

	// maxTextureSize := int32(1024) // opengl.GetGlCaps().MaxTextureSize
	buffer.Width = 800
	buffer.Height = 600

	gl.GenFramebuffers(1, &buffer.shadowBuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, buffer.shadowBuffer)

	gl.GenTextures(1, &buffer.shadowTexture)
	gl.BindTexture(gl.TEXTURE_2D, buffer.shadowTexture)
	gl.TexStorage2D(gl.TEXTURE_2D, 1, gl.DEPTH_COMPONENT16, buffer.Width, buffer.Height)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_COMPARE_MODE, gl.COMPARE_REF_TO_TEXTURE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_COMPARE_FUNC, gl.LEQUAL)

	gl.FramebufferTexture(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, buffer.shadowTexture, 0)
	gl.DrawBuffer(gl.NONE)

	framebufferStatus := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	if framebufferStatus != gl.FRAMEBUFFER_COMPLETE {
		panic(fmt.Sprintf("The framebuffer status is not complete, found: %v\n", framebufferStatus))
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	return &buffer
}

func (r *ShadowBuffer) GetTextureId() uint32 {
	return r.shadowTexture
}

func (r *ShadowBuffer) Delete() {
	gl.DeleteTextures(1, &r.shadowTexture)
	gl.DeleteFramebuffers(1, &r.shadowBuffer)
}

func (r *ShadowBuffer) RenderToBuffer(renderFunction func()) {
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.shadowBuffer)

	renderFunction()

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}
