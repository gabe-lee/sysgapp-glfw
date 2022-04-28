package sysgappglfw

import (
	sys "github.com/gabe-lee/sysgapp"
)

// var TextureIndexToGLMap = map[TextureIndex]uint32{
// 	MainTexture: gl.TEXTURE0,
// }

// TriangleStrip = VertexMode(gl.TRIANGLE_STRIP)
// 	TriangleFan   = VertexMode(gl.TRIANGLE_FAN)
// 	Triangles     = VertexMode(gl.TRIANGLES)
// 	Pixels        = VertexMode(gl.POINTS)

var vMode = map[sys.VertexMode]int32{}

var UniversalVertexShader2D = sys.NewShader(VertexShader, `
#version 410 core

uniform vec2 surfaceSize;
uniform vec2 halfSurfaceSize;
uniform vec2 textureSize;

in vec2 vertexPosition;
in int vertexColor;
in vec2 textureUV;

out vec2 fragUV;
out vec4 fragColor;

void main() {
    gl_Position.x = (vertexPosition.x/halfSurfaceSize.x)-1.0;
	gl_Position.y = ((surfaceSize.y-vertexPosition.y)/halfSurfaceSize.y)-1.0;
    gl_Position.w = 1.0;
	fragUV.x = textureUV.x/textureSize.x;
	fragUV.y = textureUV.y/textureSize.y;
	fragColor = unpackUnorm4x8(vertexColor);
}
`+"\x00")
var UniversalFragmentShader2D = NewShader(VertexShader, `
#version 410 core

const vec4 WHITE = vec4(1.0, 1.0, 1.0, 1.0);
const vec2 NO_UV = vec2(-1.0, -1.0);
const int BLEND_SHADOWS = 1 << 0;
const int NO_TEX = 1 << 1;
const int TILE_X = 1 << 2;
const int TILE_Y = 1 << 3;

uniform int renderFlags;
uniform sampler2D tex2D;

in vec2 fragUV;
in vec4 fragColor;

out vec4 pixelColor;

void main() {
	vec4 texColor;
	vec2 adjustedUV = fragUV;
	if ((renderFlags & NO_TEX) == 0 && (renderFlags & TILE_X) > 0) {
		while (adjustedUV.x < 0) {
			adjustedUV.x += 1.0;
		}
		while (adjustedUV.x > 1) {
			adjustedUV.x -= 1.0;
		}
    }
	if ((renderFlags & NO_TEX) == 0 && (renderFlags & TILE_Y) > 0) {
		while (adjustedUV.y < 0) {
			adjustedUV.y += 1.0;
		}
		while (adjustedUV.y > 1) {
			adjustedUV.y -= 1.0;
		}
    }
	if (adjustedUV == NO_UV || adjustedUV.x < 0 || adjustedUV.x > 1 || adjustedUV.y < 0 || adjustedUV.y > 1) {
		texColor = WHITE;
	} else {
		texColor = texture(tex2D, adjustedUV);
	}
	if ((renderFlags & BLEND_SHADOWS) > 0 || texColor.a == 0.0 || texColor.a == 1.0 || texColor.r != texColor.g || texColor.r != texColor.b || texColor.g != texColor.b) {
		pixelColor = texColor * fragColor;
	} else {
		pixelColor = texColor;
	}
}
`+"\x00")

type ShaderIndex int

const (
	UniversalVertex2D ShaderIndex = iota
	UniversalFragment2D
	Primitive2DVertex
	Primitive2DFragment
	Textured2DVertex
	Textured2DFragment
	DitherFadeTexFragment
)

// SHADER LOCATIONS
const (
	// Inputs
	VertexPosition = "vertexPosition\x00"
	VertexColor    = "vertexColor\x00"
	TextureUV      = "textureUV\x00"

	// Outputs
	FragUV     = "fragUV\x00"
	FragColor  = "fragColor\x00"
	PixelColor = "pixelColor\x00"

	// Uniforms
	SurfaceSize     = "surfaceSize\x00"
	HalfSurfaceSize = "halfSurfaceSize\x00"
	RenderFlags     = "renderFlags\x00"
	Texture2D       = "tex2D\x00"
	TextureSize     = "textureSize\x00"

	// Old Locations
	//SourceSize = "sourceSize\x00"
	//QuarterDestSize = "quarterDestSize\x00"
	//DestSize = "destSize\x00"
	//DrawColor = "drawColor\x00"
	//BlendAlpha = "blendAlpha\x00"
	//Sampler2D = "tex2D\x00"
	//DestVertex = "destVert\x00"
	//
	//SourceVertex = "sourceVert\x00"
)
