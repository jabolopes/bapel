#include "bin/main_helper.h"
#include <vector>
#include <string>
#include <tuple>
#include <sstream>
#include <fstream>
#include <filesystem>
#include <iostream>
#include <cstdio>
#include <memory>
#include <stdexcept>
#include <sys/wait.h>

// Forward declaration of Bapel main
int32_t bapel_main();

std::vector<std::string> g_args;

int main(int argc, char** argv) {
    g_args.clear();
    for (int i = 0; i < argc; ++i) {
        g_args.push_back(argv[i]);
    }
    return bapel_main();
}

namespace cli {

int64_t getArgCount() {
    return g_args.size();
}

std::string getArg(int64_t i) {
    if (i < 0 || i >= g_args.size()) {
        return "";
    }
    return g_args[i];
}

std::tuple<int64_t, std::string> exec(std::string cmd, const std::vector<std::string>& args) {
    std::string full_cmd = cmd;
    for (const auto& arg : args) {
        // Simple quoting to handle spaces in arguments.
        // A more robust solution might be needed for complex arguments.
        full_cmd += " \"" + arg + "\"";
    }
    full_cmd += " 2>&1";

    std::string result;
    char buffer[128];
    FILE* pipe = popen(full_cmd.c_str(), "r");
    if (!pipe) {
        return { -1, "popen failed" };
    }
    while (fgets(buffer, sizeof(buffer), pipe) != NULL) {
        result += buffer;
    }
    int status = pclose(pipe);
    int exit_code = WEXITSTATUS(status);
    return { exit_code, result };
}


std::vector<PackageMapping> parseWorkspaceFlat(std::string text) {
    std::vector<PackageMapping> mappings;
    std::istringstream iss(text);
    std::string line;
    while (std::getline(iss, line)) {
        if (line.empty()) continue;
        std::istringstream line_iss(line);
        std::string type, name, path;
        if (line_iss >> type >> name >> path) {
            PackageMapping mapping;
            mapping.is_prefix = (type == "PREFIX");
            mapping.name = name;
            mapping.path = path;
            mappings.push_back(mapping);
        }
    }
    return mappings;
}

SourceFileInfo parseSourceFileFlat(std::string text) {
    SourceFileInfo info;
    std::istringstream iss(text);
    std::string line;
    while (std::getline(iss, line)) {
        if (line.empty()) continue;
        std::istringstream line_iss(line);
        std::string type, value;
        if (line_iss >> type >> value) {
            if (type == "IMPORT") {
                info.importModules.push_back(value);
            } else if (type == "IMPL") {
                info.implFiles.push_back(value);
            }
        }
    }
    return info;
}




std::string replaceSeparator(std::string s, std::string from, std::string to) {
    size_t pos = 0;
    while ((pos = s.find(from, pos)) != std::string::npos) {
         s.replace(pos, from.length(), to);
         pos += to.length();
    }
    return s;
}

bool isPrefixOf(std::string prefix, std::string s) {
    if (s == prefix) return true;
    return s.rfind(prefix + ".", 0) == 0;
}










std::vector<std::string> getSubArgs(int64_t start) {
    std::vector<std::string> sub;
    if (start < 0 || start >= g_args.size()) return sub;
    sub.assign(g_args.begin() + start, g_args.end());
    return sub;
}


} // namespace cli
