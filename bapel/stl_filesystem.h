#pragma once

#include <filesystem>
#include <string>

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

inline bool exists(std::string path) {
    return std::filesystem::exists(path);
}

inline bool create_directories(std::string path) {
    try {
        std::filesystem::create_directories(path);
        return true;
    } catch (...) {
        return false;
    }
}

inline bool remove(std::string path) {
    try {
        return std::filesystem::remove(path);
    } catch (...) {
        return false;
    }
}

inline bool copy(std::string src, std::string dst) {
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

inline bool set_current_path(std::string path) {
    try {
        std::filesystem::current_path(path);
        return true;
    } catch (...) {
        return false;
    }
}

inline std::string extension(std::string path) {
    return std::filesystem::path(path).extension().string();
}

inline std::string parent_path(std::string path) {
    return std::filesystem::path(path).parent_path().string();
}

inline std::string stem(std::string path) {
    return std::filesystem::path(path).stem().string();
}

} // namespace fs
