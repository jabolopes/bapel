module;

#include <variant>
#include <vector>

#include <SDL3/SDL.h>

export module game:game_material;

import :game_game;

export {

struct MaterialTexture {

};

struct MaterialShape {
  SDL_FRect dst_rect = {0, 0, 0, 0};
  SDL_FColor fill_color = {0, 0, 0, 0};
  SDL_FColor draw_color = {0, 0, 0, 0};
};

struct MaterialLayer {
  bool hidden = false;
  std::variant<MaterialTexture, MaterialShape> value = MaterialShape{};
};

struct Material {
  bool hidden = false;
  std::vector<MaterialLayer> layers;
};

void renderShape(const MaterialShape& shape) {
  auto *renderer = game.renderer();

  if (shape.fill_color.a > 0) {
    SDL_SetRenderDrawColorFloat(renderer, shape.fill_color.r, shape.fill_color.g, shape.fill_color.b, shape.fill_color.a);
    SDL_RenderFillRect(renderer, &shape.dst_rect);
  }

  if (shape.draw_color.a > 0) {
    SDL_SetRenderDrawColorFloat(renderer, shape.draw_color.r, shape.draw_color.g, shape.draw_color.b, shape.draw_color.a);
    SDL_RenderRect(renderer, &shape.dst_rect);
  }
}

void renderTexture(const MaterialTexture& texture) {

}

void renderLayer(const MaterialLayer& layer) {
  if (layer.hidden) {
    return;
  }

  if (auto* texture = std::get_if<MaterialTexture>(&layer.value); texture != nullptr) {
    renderTexture(*texture);
    return;
  }

  renderShape(std::get<MaterialShape>(layer.value));
}

void renderMaterial(const Material& material) {
  if (material.hidden) {
    return;
  }

  for (const auto& layer : material.layers) {
    renderLayer(layer);
  }
}

}  // export
