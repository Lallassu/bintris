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

void main() {
   	if (eff.x == 0.0) {
   	    color = texture(image, uv);
	 	if (color.r > 0.5 && color.g < 0.2) {
			color.r = 0.0;
		}
 	} else if (eff.x == 1.0 || eff.x == 9.0) { // MetaBalls  + EffectBg
	     vec4 c = texture(image, uv);
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
		    c.r = max(0.3, uPulse);
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
   	 color = texture(image, uv);
	 if (pos.y > 0.75) {
	 	color.a = 1.0-pos.y;
	 }
	 if (eff.x == 4.0) {
	 	if (color.r > 0.5 && color.g < 0.2) {
			color.r = 0.0;
		}
	 }
   } else if(eff.x == 3.0) { // Stats
   	  color = texture(image, uv);
	  if (color.r > 0.5 && color.g < 0.2) {
	  	color.r = 0.0;
	  }
	  color.b += 0.8;
   } else if (eff.x == 8.0) { // stats blinking
   	  color = texture(image, uv);
	  if (color.r > 0.5 && color.g < 0.2) {
	  	color.b += sin(uTime*10.0);
	  }

   } else if(eff.x == 5.0) { // GameOver (Red)
   	    color = texture(image, uv);
   		if (color.a > 0.4) {
			color.r = 1.0;
			color.a = 0.2;
	 	}
   } else if(eff.x == 6.0 || eff.x == 7.0) { // Game Over "Logo"
   	    color = texture(image, uv);
	 	if (color.r > 0.5 && color.g < 0.2) {
			color.a = 0.0;
		}
   } else {
   	 color = texture(image, uv);
	 	if (color.r > 0.5 && color.g < 0.2) {
			color.r = 0.0;
		}
   }
}
`
