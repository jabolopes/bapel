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

// @bpl: export type Material
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

// @bpl: export newRect: (i64, i64, i64, i64) -> Material
Material newRect(int64_t x, int64_t y, int64_t w, int64_t h) {
  return Material{
    .layers = {
      MaterialLayer{
        .value = MaterialShape{
          .dst_rect = SDL_FRect{float(x), float(y), float(w), float(h)},
          .fill_color = SDL_FColor{1.0, 0.0, 0.0, 1.0},
        }
      }
    }
  };
}

// @bpl: export addMaterial: Material -> ()
void addMaterial(Material material) {
  auto& ecs = game.ecs();
  const auto entity = ecs.create();
  ecs.emplace<Material>(entity, std::move(material));
}

}  // export
