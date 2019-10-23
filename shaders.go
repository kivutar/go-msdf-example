package main

var fragmentShader = `
#version 140
#if __VERSION__ >= 130
#define COMPAT_VARYING in
#define COMPAT_ATTRIBUTE in
#define COMPAT_TEXTURE texture
#define COMPAT_FRAGCOLOR FragColor
out vec4 COMPAT_FRAGCOLOR;
#else
#define COMPAT_VARYING varying
#define COMPAT_ATTRIBUTE attribute
#define COMPAT_TEXTURE texture2D
#define COMPAT_FRAGCOLOR gl_FragColor
#endif
uniform sampler2D Texture;
COMPAT_VARYING vec2 fragTexCoord;

float median(float r, float g, float b) {
	return max(min(r, g), min(max(r, g), b));
}

void main() {
	vec3 sample = COMPAT_TEXTURE(Texture, fragTexCoord).rgb;
	ivec2 sz = textureSize(Texture, 0).xy;
	float dx = dFdx(fragTexCoord.x) * sz.x;
	float dy = dFdy(fragTexCoord.y) * sz.y;
	float toPixels = 8.0 * inversesqrt(dx * dx + dy * dy);
	float sigDist = median(sample.r, sample.g, sample.b);
	float w = fwidth(sigDist);
	float opacity = smoothstep(0.5 - w, 0.5 + w, sigDist);
	COMPAT_FRAGCOLOR = vec4(vec3(1.0,0.0,0.0), opacity);
}
` + "\x00"

var vertexShader = `
#version 140
#if __VERSION__ >= 130
#define COMPAT_VARYING out
#define COMPAT_ATTRIBUTE in
#define COMPAT_TEXTURE texture
#define COMPAT_FRAGCOLOR FragColor
out vec4 COMPAT_FRAGCOLOR;
#else
#define COMPAT_VARYING varying
#define COMPAT_ATTRIBUTE attribute
#define COMPAT_TEXTURE texture2D
#define COMPAT_FRAGCOLOR gl_FragColor
#endif
COMPAT_ATTRIBUTE vec2 vert;
COMPAT_ATTRIBUTE vec2 vertTexCoord;
COMPAT_VARYING vec2 fragTexCoord;
void main() {
  fragTexCoord = vertTexCoord;
  gl_Position = vec4(vert, 0.0, 1.0);
}
` + "\x00"
