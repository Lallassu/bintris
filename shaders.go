package main

const vertexShader = `#version 300 es
layout (location = 0) in vec4 position;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform float uTime;
uniform float uDt;
uniform int effect;

out vec2 TexCoords;
out mediump vec4 vColor;
out mediump vec4 pos;
void main() {
    TexCoords = position.zw;

	vec4 p = position;
	if (effect == 0) { // EffectBlueShine
	    vColor.b += (sin(uTime/2.0)*p.y)/p.x;
	    vColor.b = max(vColor.b, 0.4);
	    if (vColor.b > 0.9 && vColor.r > 0.9 && vColor.g > 0.5) {
	    }
	     	//p.y += sin(p.x + p.y + uTime*2.0);
	}
	gl_Position = projection * view * model * vec4(p.x, p.y,0.0, 1.0); 
	pos = position;
}`

const fragmentShader = `#version 300 es
in mediump vec2 TexCoords;
uniform sampler2D image;
in mediump vec4 pos;

precision highp float;
layout (location = 0) out mediump vec4 color2;
uniform float uTime;


vec2 random2( vec2 p ) {
    return fract(sin(vec2(dot(p,vec2(127.1,311.7)),dot(p,vec2(269.5,183.3))))*43758.5453);
}

void main() {	
	vec4 c =  texture(image, TexCoords);
	vec2 u_resolution = vec2(600.0, 800.0);


    vec2 st = gl_FragCoord.xy/u_resolution.xy;
    st.x *= u_resolution.x/u_resolution.y;
    vec3 color = vec3(.0);

    // Scale
    st *= 30.;

    // Tile the space
    vec2 i_st = floor(st);
    vec2 f_st = fract(st);

    float m_dist = 1.;  // minimum distance
    for (int j= -1; j <= 1; j++ ) {
        for (int i= -1; i <= 1; i++ ) {
            // Neighbor place in the grid
            vec2 neighbor = vec2(float(i),float(j));

            // Random position from current + neighbor place in the grid
            vec2 offset = random2(i_st + neighbor);

            // Animate the offset
            offset = 0.5 + 0.5*sin(uTime + 6.2831*offset);

            // Position of the cell
            vec2 pos = neighbor + offset - f_st;

            // Cell distance
            float dist = length(pos);

            // Metaball it!
            m_dist = min(m_dist, m_dist*dist);
        }
    }

    // Draw cells
    color += step(0.2*pos.y, m_dist);

	if (c.a > 0.0) {
		color.r += 0.5;
		color.g += 0.4;
		color.b += 0.9;
		color2 = c * vec4(color,1.0);
	} else {
		color2 = c * vec4(color,1.0);
	}
}
`
