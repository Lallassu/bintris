package main

const vertexShader = `#version 300 es
layout (location = 0) in vec2 vert;
layout (location = 1) in vec2 uvs;
layout (location = 2) in vec2 effect;

uniform highp float uTime;

out vec2 uv;
out vec2 pos;
out vec2 eff;

void main() {
	uv = uvs;
	pos = vec2(gl_Position.x, gl_Position.y);
	eff = effect;

	if (eff.x == 1.0) {
		gl_Position = vec4(2.0*vert.x-1.0, 2.0*vert.y-1.0, 0.0, 1.0);
	} else {
		gl_Position = vec4(2.0*vert.x-1.0, 2.0*vert.y-1.0, 0.0, 1.0);
	}
}
`
const fragmentShader = `#version 300 es
in mediump vec2 uv;
in mediump vec2 pos;
in mediump vec2 eff;
uniform sampler2D image;
uniform highp float uTime;
layout (location = 0) out highp vec4 color;

void main() {
	color = texture(image, uv);
	if (eff.x == 1.0) {
		//	color.b = 1.0 * sin(uTime);
		if (fract(sin(pos.y*uTime)) > 0.1) {
			color.b -= sin(pos.y*uTime);
		} else {
			color.b += sin(pos.y*uTime);
		}
	}
}

`

const vertexShader2 = `#version 300 es
layout (location = 0) in vec4 position;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform highp float uTime;
//uniform int effect;

out vec2 TexCoords;
out mediump vec4 vColor;
out mediump vec4 pos;
void main() {
    TexCoords = position.zw;

	vec4 p = position;

	gl_Position = projection * view * model * vec4(p.x, p.y,0.0, 1.0); 
	pos = position;
}`

const fragmentShader2 = `#version 300 es
in mediump vec2 TexCoords;
uniform sampler2D image;
in mediump vec4 pos;

precision highp float;
precision lowp int;
layout (location = 0) out highp vec4 color2;
uniform highp float uTime;
uniform int effect;


vec2 random2( vec2 p ) {
    return fract(sin(vec2(dot(p,vec2(127.1,311.7)),dot(p,vec2(269.5,183.3))))*43758.5453);
}

void main() {	
	vec4 c =  texture(image, TexCoords);

	if (effect == 0) {
		color2 = c;
 	} else if (effect == 1 || effect == 3) { // Metaballs + MetaballsBlue
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


		 color += step(0.4*pos.y, m_dist);
   		 if (color.r > 0.0 && color.b > 0.0 && color.g > 0.0) {
   		 	if (c.a > 0.0) {
   		   		 c.r += 0.1*cos(uTime);
   		    		 c.g += 0.1;
   		   		 c.b += 0.1*sin(uTime);
   		     }
			 if (effect == 3) {
		    	if (c.a > 0.0 && color.r > 0.5) {
		    		c.b += sin(uTime)/5.0+pos.x*pos.y;
		    		c.b = max(0.8, c.b);
		    	}
			}
   		 	color2 = c * vec4(color,1.0);
   		 } else {
   		 	color2 = c;
   		 }
	} else if(effect == 2) {
		color2 = c;
	} else if(effect == 3) {
		if (c.a > 0.0) {
			c.b += sin(uTime)/5.0+pos.x*pos.y;
			c.b = max(0.8, c.b);
		}
		color2 = c;
	}
}
`
