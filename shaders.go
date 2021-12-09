package main

const vertexShader = `#version 300 es
layout (location = 0) in vec2 vert;
layout (location = 1) in vec2 uvs;
layout (location = 2) in vec2 effect;

uniform highp float uTime;
uniform highp float uPulse;

out mediump vec2 uv;
out mediump vec2 pos;
out mediump vec2 eff;

void main() {
	uv = uvs;
	gl_Position = vec4(2.0*vert.x-1.0, 2.0*vert.y-1.0, 0.0, 1.0);
	pos = gl_Position.xy;

	eff = effect;
	if (eff.x == 1.0) {
		gl_Position.x = pos.x - sin(uTime)/30.0;
		gl_Position.y = pos.y - cos(uTime)/20.0;
	} else if(eff.x == 6.0) {
		gl_Position.x = pos.x - sin(uTime)/30.0;
	} else if(eff.x == 7.0) {
		gl_Position.x = pos.x + sin(uTime)/30.0;
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
uniform highp float uPulse;

layout (location = 0) out highp vec4 color;

precision highp float;
precision lowp int;


vec2 random2( vec2 p ) {
    return fract(sin(vec2(dot(p,vec2(127.1,311.7)),dot(p,vec2(269.5,183.3))))*43758.5453);
}

float random (in vec2 _st) {
    return fract(sin(dot(_st.xy,
                         vec2(12.9898,78.233)))*
        43758.5453123);
}

// Based on Morgan McGuire @morgan3d
// https://www.shadertoy.com/view/4dS3Wd
float noise (in vec2 _st) {
    vec2 i = floor(_st);
    vec2 f = fract(_st);

    // Four corners in 2D of a tile
    float a = random(i);
    float b = random(i + vec2(1.0, 0.0));
    float c = random(i + vec2(0.0, 1.0));
    float d = random(i + vec2(1.0, 1.0));

    vec2 u = f * f * (3.0 - 2.0 * f);

    return mix(a, b, u.x) +
            (c - a)* u.y * (1.0 - u.x) +
            (d - b) * u.x * u.y;
}

#define NUM_OCTAVES 5

float fbm ( in vec2 _st) {
    float v = 0.0;
    float a = 0.5;
    vec2 shift = vec2(100.0);
    // Rotate to reduce axial bias
    mat2 rot = mat2(cos(0.5), sin(0.5),
                    -sin(0.5), cos(0.50));
    for (int i = 0; i < NUM_OCTAVES; ++i) {
        v += a * noise(_st);
        _st = rot * _st * 2.0 + shift;
        a *= 0.5;
    }
    return v;
}

void main() {
   	color = texture(image, uv);
   	if (eff.x == 0.0) {
	 	if (color.r > 0.5 && color.g < 0.2) {
			color.r = 0.0;
		}
    } else if(eff.x == 10.0 || eff.x == 11.0) { // Menu bg
       vec2 u_resolution = vec2(1000.0, 1000.0);
       vec2 st = gl_FragCoord.xy/u_resolution.xy*3.;
       //st += st * abs(sin(uTime*0.1)*3.0);
       vec3 c = vec3(0.0);

       vec2 q = vec2(0.);
       q.x = fbm( st + 0.00*uTime);
       q.y = fbm( st + vec2(1.0));

       vec2 r = vec2(0.);
       r.x = fbm( st + 1.0*q + vec2(1.7,9.2)+ 0.15*uTime );
       r.y = fbm( st + 1.0*q + vec2(8.3,2.8)+ 0.126*uTime);

       float f = fbm(st+r);

       c = mix(vec3(0.101961,0.619608,0.666667),
                   vec3(0.666667,0.666667,0.198039),
                   clamp((f*f)*4.0,0.0,1.0));

       c = mix(c,
                   vec3(0,0,sin(uTime*0.164706)),
                   clamp(length(q),0.0,1.0));

       c = mix(c,
                   vec3(sin(uTime*0.666667),1,1),
                   clamp(length(r.x),0.0,1.0));

	   if (eff.x == 10.0) {
       	   color = vec4((f*f*f+.6*f*f+.5*f)*c,1.0);
		} else {
       	   color = vec4((f*f*f+.6*f*f+.5*f)*c,0.2);
		}
 	} else if (eff.x == 1.0 || eff.x == 9.0) { // MetaBalls  + EffectBg
	     vec4 c = color;
		 if (uTime < 4.0) {
		 	c.a += sin(c.a*uTime/3.0)-1.0;
		 } 

         vec2 u_resolution = vec2(1000.0, 1000.0);

         vec2 st = gl_FragCoord.xy/u_resolution.xy;
         st.x *= u_resolution.x/u_resolution.y;
         vec3 color2 = vec3(.0);

         // Scale
		 if (eff.x == 9.0) {
		 	st *= 60.0;
		 } else {
         	st *= 20.0;
		 }

         // Tile the space
         vec2 i_st = floor(st);
         vec2 f_st = fract(st);

         float m_dist = 1.0;  // minimum distance
         for (int j= -1; j <= 1; j++ ) {
             for (int i= -1; i <= 1; i++ ) {
                 // Neighbor place in the grid
                 vec2 neighbor = vec2(float(i),float(j));

   		         // Random position from current + neighbor place in the grid
   		         vec2 offset = random2(i_st + neighbor);

   		         // Animate the offset
   		         offset = 0.5 + 0.5*sin(uTime + 6.2831*offset);

   		         // Position of the cell
        	     vec2 pos2 = neighbor + offset - f_st;

  		         // Cell distance
   		         float dist = length(pos2);

   		         // Metaball it!
				 if (eff.x == 9.0) { 
  		         	//m_dist = min(m_dist, m_dist*dist)*(1.0+(pos.y-uPulse));
  		         	m_dist = min(m_dist, m_dist*dist)*(pos.y-uPulse+1.5);
				 } else {
  		         	m_dist = min(m_dist, m_dist*dist)*(pos.y+0.2);
				 }
   		     }
   		 }

		if (eff.x == 9.0) {
          color2 += step(0.01, m_dist);
       	  color = c * vec4(color2,1.0);
	      if (color.r == 0.0 && color.b == 0.0 && color.g == 0.0 && color.a > 0.0) {
		    c.r = max(0.66, uPulse);
			c.g = 0.0;
			c.b = 0.0;
	      	color = c;
		  } else {
		    c.g += sin(uTime/10.0)/10.0;
		    c.r += cos(uTime/5.0)/10.0;
		    c.b += sin(uTime/20.0)/10.0;
	      	color = c;
		  }
		} else {
          color2 += step(0.05*pos.y, m_dist);
       	  color = c * vec4(color2,1.0);
	      if (color.r == 0.0 && color.b == 0.0 && color.g == 0.0 && color.a > 0.0) {
	      	color = c;
	      	color.a = clamp(sin(pos.x*uTime), 1.0, 0.5);
	     	color.b = 0.8;
	      }
	      if (color.r > 0.5 && color.g < 0.2) {
	   	    color.b = clamp(sin(uTime), 0.0, 0.8);
	   	    color.r = clamp(sin(uTime), 0.0, 0.3);
	   	  }
		} 
   } else if(eff.x == 2.0 || eff.x == 4.0) { // TileTop
	 if (pos.y > 0.75) {
	 	color.a = 1.0-pos.y;
	 }
	 if (eff.x == 4.0) {
	 	if (color.r > 0.5 && color.g < 0.2) {
			color.r = 0.0;
		}
	 }
   } else if(eff.x == 3.0) { // Stats
	  if (color.r > 0.5 && color.g < 0.2) {
	  	color.r = 0.0;
	  }
	  color.b += 0.8;
   } else if (eff.x == 8.0) { // stats blinking
      color.b = sin(uTime*7.0);

   } else if(eff.x == 5.0) { // GameOver (Red)
   		if (color.a > 0.4) {
			color.r = 1.0;
			color.a = 0.2;
	 	}
   } else if(eff.x == 6.0 || eff.x == 7.0) { // Game Over "Logo"
	 	if (color.r > 0.5 && color.g < 0.2) {
			color.a = 0.0;
		}
   } else {
	 	if (color.r > 0.5 && color.g < 0.2) {
			color.r = 0.0;
		}
   }
}
`
