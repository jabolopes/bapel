#pragma once

#include <cstdint>
#include <iomanip>
#include <sstream>
#include <string>

namespace ir {

class Sha1 {
public:
  Sha1() { reset(); }

  void reset() {
    digest_[0] = 0x67452301;
    digest_[1] = 0xEFCDAB89;
    digest_[2] = 0x98BADCFE;
    digest_[3] = 0x10325476;
    digest_[4] = 0xC3D2E1F0;
    transforms_ = 0;
    buffer_.clear();
  }

  void update(const std::string& s) {
    for (char c : s) {
      buffer_.push_back(static_cast<uint8_t>(c));
      if (buffer_.size() == 64) {
        process_block(buffer_.data());
        transforms_++;
        buffer_.clear();
      }
    }
  }

  std::string final() {
    uint64_t total_bits = (transforms_ * 64 + buffer_.size()) * 8;
    buffer_.push_back(0x80);
    while (buffer_.size() % 64 != 56) {
      buffer_.push_back(0x00);
    }
    for (int i = 7; i >= 0; --i) {
      buffer_.push_back(static_cast<uint8_t>((total_bits >> (i * 8)) & 0xFF));
    }
    for (size_t i = 0; i < buffer_.size(); i += 64) {
      process_block(buffer_.data() + i);
    }

    std::ostringstream result;
    for (uint32_t val : digest_) {
      result << std::hex << std::setfill('0') << std::setw(8) << val;
    }
    reset();
    return result.str();
  }

  static std::string hash(const std::string& s) {
    Sha1 checksum;
    checksum.update(s);
    return checksum.final();
  }

private:
  static uint32_t rol(uint32_t value, size_t bits) {
    return (value << bits) | (value >> (32 - bits));
  }

  static uint32_t blk(const uint32_t block[16], size_t i) {
    return rol(block[(i + 13) & 15] ^ block[(i + 8) & 15] ^ block[(i + 2) & 15] ^ block[i], 1);
  }

  void process_block(const uint8_t block_bytes[64]) {
    uint32_t w[80];
    for (size_t i = 0; i < 16; ++i) {
      w[i] = (static_cast<uint32_t>(block_bytes[i * 4]) << 24) |
             (static_cast<uint32_t>(block_bytes[i * 4 + 1]) << 16) |
             (static_cast<uint32_t>(block_bytes[i * 4 + 2]) << 8) |
             (static_cast<uint32_t>(block_bytes[i * 4 + 3]));
    }
    for (size_t i = 16; i < 80; ++i) {
      w[i] = rol(w[i - 3] ^ w[i - 8] ^ w[i - 14] ^ w[i - 16], 1);
    }

    uint32_t a = digest_[0];
    uint32_t b = digest_[1];
    uint32_t c = digest_[2];
    uint32_t d = digest_[3];
    uint32_t e = digest_[4];

    for (size_t i = 0; i < 20; ++i) {
      uint32_t f = (b & c) | ((~b) & d);
      uint32_t k = 0x5A827999;
      uint32_t temp = rol(a, 5) + f + e + k + w[i];
      e = d; d = c; c = rol(b, 30); b = a; a = temp;
    }
    for (size_t i = 20; i < 40; ++i) {
      uint32_t f = b ^ c ^ d;
      uint32_t k = 0x6ED9EBA1;
      uint32_t temp = rol(a, 5) + f + e + k + w[i];
      e = d; d = c; c = rol(b, 30); b = a; a = temp;
    }
    for (size_t i = 40; i < 60; ++i) {
      uint32_t f = (b & c) | (b & d) | (c & d);
      uint32_t k = 0x8F1BBCDC;
      uint32_t temp = rol(a, 5) + f + e + k + w[i];
      e = d; d = c; c = rol(b, 30); b = a; a = temp;
    }
    for (size_t i = 60; i < 80; ++i) {
      uint32_t f = b ^ c ^ d;
      uint32_t k = 0xCA62C1D6;
      uint32_t temp = rol(a, 5) + f + e + k + w[i];
      e = d; d = c; c = rol(b, 30); b = a; a = temp;
    }

    digest_[0] += a;
    digest_[1] += b;
    digest_[2] += c;
    digest_[3] += d;
    digest_[4] += e;
  }

  uint32_t digest_[5];
  std::vector<uint8_t> buffer_;
  uint64_t transforms_ = 0;
};

} // namespace ir
