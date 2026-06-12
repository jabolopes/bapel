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

std::string joinPath(std::string a, std::string b) {
    return (std::filesystem::path(a) / b).string();
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

std::string concat(std::string a, std::string b) {
    return a + b;
}



struct BazelTarget {
    std::string type;
    std::string name;
    std::vector<std::string> srcs;
    std::vector<std::string> hdrs;
    std::vector<std::string> deps;
};
std::vector<BazelTarget> g_targets;

void addTarget(std::string type, std::string name, const std::vector<std::string>& srcs, const std::vector<std::string>& hdrs, const std::vector<std::string>& deps) {
    g_targets.push_back({type, name, srcs, hdrs, deps});
}

bool writeBuildFile() {
    try {
        std::ofstream build("out/BUILD");
        if (!build.is_open()) return false;
        build << "load(\"@rules_cc//cc:defs.bzl\", \"cc_binary\", \"cc_library\")\n\n";
        for (const auto& target : g_targets) {
            build << target.type << "(\n";
            build << "    name = \"" << target.name << "\",\n";
            if (!target.srcs.empty()) {
                build << "    srcs = [\n";
                for (const auto& src : target.srcs) {
                    build << "        \"" << src << "\",\n";
                }
                build << "    ],\n";
            }
            if (!target.hdrs.empty()) {
                build << "    hdrs = [\n";
                for (const auto& hdr : target.hdrs) {
                    build << "        \"" << hdr << "\",\n";
                }
                build << "    ],\n";
            }
            build << "    copts = [\n";
            build << "        \"-std=c++17\",\n";
            build << "        \"-Xassembler\",\n";
            build << "        \"--gsframe=no\",\n";
            build << "    ],\n";
            if (!target.deps.empty()) {
                build << "    deps = [\n";
                for (const auto& dep : target.deps) {
                    build << "        \"" << dep << "\",\n";
                }
                build << "    ],\n";
            }
            build << ")\n\n";
        }
        return true;
    } catch (...) {
        return false;
    }
}

bool ensureWorkspaceSetup() {
    try {
        std::filesystem::create_directories("out");
        
        std::ofstream ws("out/WORKSPACE");
        if (!ws.is_open()) return false;
        ws << "workspace(name = \"bapel_out\")\n";
        ws.close();

        std::ofstream mod("out/MODULE.bazel");
        if (!mod.is_open()) return false;
        mod << "module(name = \"bapel_out\")\n";
        mod << "bazel_dep(name = \"rules_cc\", version = \"0.2.17\")\n";
        mod.close();
        return true;
    } catch (...) {
        return false;
    }
}



std::vector<std::string> getSubArgs(int64_t start) {
    std::vector<std::string> sub;
    if (start < 0 || start >= g_args.size()) return sub;
    sub.assign(g_args.begin() + start, g_args.end());
    return sub;
}


} // namespace cli
