module;

#include <functional>

#include <SDL3/SDL.h>

#include <entt/entt.hpp>

export module game:game_game;

export {

class Game final {
public:
  Game() = default;

  void set_window(SDL_Window *window) { window_ = window; }

  SDL_Renderer* renderer() const { return renderer_; }
  void set_renderer(SDL_Renderer* renderer) { renderer_ = renderer; }

  std::function<void()>& update() { return update_; }
  void set_update(std::function<void()> update) { update_ = std::move(update); }

  entt::registry& ecs() { return ecs_; }

private:
  SDL_Window *window_ = nullptr;
  SDL_Renderer *renderer_ = nullptr;
  std::function<void()> update_;
  entt::registry ecs_;
};

Game game;

}  // export
