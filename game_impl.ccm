module;

#include <functional>
#include <optional>
#include <string>
#include <variant>
#include <vector>

#include <SDL3/SDL.h>
#include <SDL3/SDL_main.h>

#include <entt/entt.hpp>

export module game:game_impl;

import bapel.core;

import :game_material;

Uint32 pushEventUserCallback(void *userdata, SDL_TimerID timerID, Uint32 interval) {
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

void updateCallback() {
  game.update()();
}

void renderCallback() {
  auto *renderer = game.renderer();
  SDL_SetRenderDrawColor(renderer, 0, 0, 0, 0);
  SDL_RenderClear(renderer);

  for (const auto& [_, material] : game.ecs().view<Material>().each()) {
    renderMaterial(material);
  }

  SDL_RenderPresent(renderer);
}

export {

// @bpl: export type Entity
using Entity = entt::entity;

// @bpl: export add: () -> Entity
Entity add() {
  return game.ecs().create();
}

// @bpl: export init: forall ['a] (Entity, 'a) -> Entity
template <typename A>
Entity init(Entity entity, A a) {
  game.ecs().erase_if(entity, [](auto&, auto&) { return true; });
  game.ecs().emplace<A>(entity, std::move(a));
  return entity;
}

// @bpl: export init2: forall ['a, 'b] (Entity, 'a, 'b) -> Entity
template <typename A, typename B>
Entity init2(Entity entity, A a, B b) {
  game.ecs().erase_if(entity, [](auto&, auto&) { return true; });
  game.ecs().emplace<A>(entity, std::move(a));
  game.ecs().emplace<B>(entity, std::move(b));
  return entity;
}

// @bpl: export iterateAny: forall ['a] () -> std.optional (Entity, 'a)
template <typename A>
std::optional<std::pair<Entity, A>> iterateAny() {
  auto view = game.ecs().template view<A>();
  for (auto [entity, a]: view.each()) {
    return std::make_optional(std::make_pair(entity, a));
  }
  return std::nullopt;
}

// @bpl: export setUpdate: (() -> ()) -> ()
void setUpdate(std::function<void()> update) {
  game.set_update(std::move(update));
}

// @bpl: export gameInit: () -> i64
int gameInit() {
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

  registerUpdateTimer();
  registerRenderTimer();

  Toolkit toolkit;
  toolkit.AddRegistration(std::bind(updateCallback));
  toolkit.AddRegistration(std::bind(renderCallback));
  toolkit.Run();

  SDL_DestroyRenderer(renderer);
  SDL_DestroyWindow(window);
  SDL_Quit();
  return 0;
}

}  // export
