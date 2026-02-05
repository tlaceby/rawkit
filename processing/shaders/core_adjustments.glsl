in vec2 uv;
out vec4 frag_color;

uniform sampler2D u_texture;
uniform float u_exposure;

void main() {
    vec4 color = texture(u_texture, uv);
    
    float stops = u_exposure / 20.0;
    color.rgb *= pow(2.0, stops);
    
    frag_color = vec4(clamp(color.rgb, 0.0, 1.0), color.a);
}