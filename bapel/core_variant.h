#pragma once

#include <array>
#include <cstddef>
#include <string_view>
#include <utility>
#include <variant>

namespace bapel {

template <size_t N>
struct FixedString {
  char data[N];
  constexpr FixedString(const char (&s)[N]) {
    for (size_t i = 0; i < N; ++i) {
      data[i] = s[i];
    }
  }
  constexpr std::string_view sv() const { return {data, N - 1}; }
};

template <FixedString Tag, typename T>
struct Arm {
  static constexpr auto tag = Tag;
  using type = T;
};

template <typename... Arms>
struct Variant : std::variant<typename Arms::type...> {
  using Base = std::variant<typename Arms::type...>;
  using Base::Base;

  template <FixedString Tag>
  static constexpr size_t index_of() {
    constexpr std::array<std::string_view, sizeof...(Arms)> tags = {Arms::tag.sv()...};
    for (size_t i = 0; i < tags.size(); ++i) {
      if (tags[i] == Tag.sv()) return i;
    }
    return sizeof...(Arms);
  }

  template <FixedString Tag, typename... Args>
  static Variant mk(Args&&... args) {
    constexpr size_t idx = index_of<Tag>();
    static_assert(idx < sizeof...(Arms), "unknown variant tag");
    return Variant(std::in_place_index<idx>, std::forward<Args>(args)...);
  }

  template <FixedString Tag>
  bool is() const {
    constexpr size_t idx = index_of<Tag>();
    static_assert(idx < sizeof...(Arms), "unknown variant tag");
    return this->index() == idx;
  }

  template <FixedString Tag>
  auto& as() {
    constexpr size_t idx = index_of<Tag>();
    static_assert(idx < sizeof...(Arms), "unknown variant tag");
    return std::get<idx>(*this);
  }

  template <FixedString Tag>
  const auto& as() const {
    constexpr size_t idx = index_of<Tag>();
    static_assert(idx < sizeof...(Arms), "unknown variant tag");
    return std::get<idx>(*this);
  }
};

}  // namespace bapel
