module;

#include <functional>
#include <string>
#include <vector>

#include <SDL3/SDL.h>
#include <SDL3/SDL_main.h>

#include <entt/entt.hpp>

export module game:game_impl;

// Needed because of import<vector> results in Bad file data:
// https://stackoverflow.com/questions/70456868/vector-in-c-module-causes-useless-bad-file-data-gcc-output
namespace std _GLIBCXX_VISIBILITY(default){}

Uint32 pushEventUserCallback(void *userdata, SDL_TimerID timerID,
                             Uint32 interval) {
  SDL_Event event;
  event.type = SDL_EVENT_USER;
  event.user.type = SDL_EVENT_USER;
  event.user.code = *(int *)userdata;
  event.user.data1 = nullptr;
  event.user.data2 = nullptr;

  SDL_PushEvent(&event);

  return interval;
}

void registerUpdateTimer() {
  static int code = 0;
  const uint32_t delay = 16; // milliseconds
  if (!SDL_AddTimer(delay, &pushEventUserCallback, /*userdata=*/&code)) {
    SDL_Log("Failed to create update timer: %s", SDL_GetError());
    return;
  }
  SDL_Log("Created update timer");
}

void registerRenderTimer() {
  static int code = 1;
  const uint32_t delay = 16; // milliseconds
  if (!SDL_AddTimer(delay, &pushEventUserCallback, /*userdata=*/&code)) {
    SDL_Log("Failed to create render timer: %s", SDL_GetError());
    return;
  }
  SDL_Log("Created render timer");
}

class Toolkit final {
public:
  Toolkit() = default;

  void AddRegistration(std::function<void()> registration);
  void Run();

private:
  std::vector<std::function<void()>> registrations_;
};

void Toolkit::AddRegistration(std::function<void()> registration) {
  registrations_.push_back(std::move(registration));
}

void Toolkit::Run() {
  while (true) {
    SDL_Event event;
    if (!SDL_WaitEvent(&event)) {
      return;
    }

    switch (event.type) {
    case SDL_EVENT_QUIT:
      return;

    case SDL_EVENT_USER:
      if (event.user.code >= 0 && event.user.code < registrations_.size()) {
        registrations_[event.user.code]();
      }
    }
  }
}

export struct Material {
  SDL_FRect dst_rect;
};

class Game final {
public:
  Game() = default;

  void set_window(SDL_Window *window) { window_ = window; }

  SDL_Renderer* renderer() const { return renderer_; }
  void set_renderer(SDL_Renderer* renderer) { renderer_ = renderer; }

  entt::registry& ecs() { return ecs_; }

private:
  SDL_Window *window_ = nullptr;
  SDL_Renderer *renderer_ = nullptr;
  entt::registry ecs_;
};

Game game;

void update() {}

void render() {
  auto *renderer = game.renderer();
  SDL_SetRenderDrawColor(renderer, 0, 0, 0, 0);
  SDL_RenderClear(renderer);

  SDL_SetRenderDrawColor(renderer, 255, 0, 0, 255);
  for (const auto& [_, material] : game.ecs().view<Material>().each()) {
    SDL_RenderFillRect(renderer, &material.dst_rect);
  }

  SDL_RenderPresent(renderer);
}

export void addMaterial(Material material) {
  auto& ecs = game.ecs();
  const auto entity = ecs.create();
  ecs.emplace<Material>(entity, std::move(material));
}

export Material newRect(int64_t x, int64_t y, int64_t w, int64_t h) {
  return Material{SDL_FRect{float(x), float(y), float(w), float(h)}};
}

export int gameInit() {
  if (!SDL_Init(SDL_INIT_VIDEO)) {
    SDL_Log("Failed to initialize SDL: %s\n", SDL_GetError());
    return 1;
  }

  SDL_Window* window = nullptr;
  SDL_Renderer* renderer = nullptr;
  if (!SDL_CreateWindowAndRenderer("bapel", 800, 600, 0, &window, &renderer)) {
    SDL_Log("Failed to create window: %s\n", SDL_GetError());
    return 1;
  }

  game.set_window(window);
  game.set_renderer(renderer);

  addMaterial(Material{SDL_FRect{0, 0, 100, 100}});

  registerUpdateTimer();
  registerRenderTimer();

  Toolkit toolkit;
  toolkit.AddRegistration(std::bind(update));
  toolkit.AddRegistration(std::bind(render));
  toolkit.Run();

  SDL_DestroyRenderer(renderer);
  SDL_DestroyWindow(window);
  SDL_Quit();
  return 0;
}
