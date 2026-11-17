#pragma once

#include <filesystem>
#include <fstream>
#include <sstream>
#include <string>
#include "bapel/stl_optional.h"

// @bpl: pub fs::exists: String -> bool
// @bpl: pub fs::create_directories: String -> bool
// @bpl: pub fs::remove: String -> bool
// @bpl: pub fs::copy: (String, String) -> bool
// @bpl: pub fs::current_path: () -> String
// @bpl: pub fs::set_current_path: String -> bool
// @bpl: pub fs::extension: String -> String
// @bpl: pub fs::parent_path: String -> String
// @bpl: pub fs::stem: String -> String
namespace fs {

inline bool exists(const std::string& path) {
    return std::filesystem::exists(path);
}

inline bool create_directories(const std::string& path) {
    try {
        std::filesystem::create_directories(path);
        return true;
    } catch (...) {
        return false;
    }
}

inline bool remove(const std::string& path) {
    try {
        return std::filesystem::remove(path);
    } catch (...) {
        return false;
    }
}

inline bool copy(const std::string& src, const std::string& dst) {
    try {
        std::filesystem::copy(src, dst, std::filesystem::copy_options::overwrite_existing);
        return true;
    } catch (...) {
        return false;
    }
}

inline std::string current_path() {
    return std::filesystem::current_path().string();
}

inline bool set_current_path(const std::string& path) {
    try {
        std::filesystem::current_path(path);
        return true;
    } catch (...) {
        return false;
    }
}

inline std::string extension(const std::string& path) {
    return std::filesystem::path(path).extension().string();
}

inline std::string parent_path(const std::string& path) {
    return std::filesystem::path(path).parent_path().string();
}

inline std::string stem(const std::string& path) {
    return std::filesystem::path(path).stem().string();
}

// @bpl: pub fs::join: (String, String) -> String
inline std::string join(const std::string& a, const std::string& b) {
    return (std::filesystem::path(a) / b).string();
}

// TODO: read_file, write_file, and lexically_normal are convenience helpers
// rather than direct STL equivalents; eventually remove or move them out of
// stl_filesystem.h.
// @bpl: pub fs::lexically_normal: String -> String
inline std::string lexically_normal(const std::string& path) {
    return std::filesystem::path(path).lexically_normal().string();
}

// @bpl: pub fs::read_file: &String -> Optional String
inline Optional<std::string> read_file(const std::string* path) {
    std::ifstream file(*path);
    if (!file.is_open()) {
        return std::nullopt;
    }
    std::ostringstream buffer;
    buffer << file.rdbuf();
    return std::make_optional(buffer.str());
}

// @bpl: pub fs::write_file: (&String, &String) -> bool
inline bool write_file(const std::string* path, const std::string* content) {
    std::ofstream file(*path);
    if (!file.is_open()) {
        return false;
    }
    file << *content;
    return file.good();
}

} // namespace fs

