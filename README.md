# Go Multi Channel Distance Field example

This repo demonstrates how to render glyph using the MSDF technique using Go, GL and GLSL.

This code can load a 64x64px glyph and display it anti aliased at any resolution.
It does not take care of generating the MSDF texture from a truetype font, only display is implemented.
The interesting part is the fragment shader.

Please see https://github.com/Chlumsky/msdfgen for more detailed explanations about MSDF and credits.

<img src="https://raw.githubusercontent.com/kivutar/go-msdf-example/master/result.png" />
