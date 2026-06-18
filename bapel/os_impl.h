#pragma once
#include <string>
#include <vector>
#include <tuple>
#include <cstdio>
#include <sys/wait.h>

namespace os {

// @bpl: pub os::exec: (String, Vector String) -> (i64, String)
inline std::tuple<int64_t, std::string> exec(std::string cmd, const std::vector<std::string>& args) {
    for (const auto& arg : args) {
        cmd += " \"" + arg + "\"";
    }
    cmd += " 2>&1";

    std::string result;
    char buffer[128];
    FILE* pipe = popen(cmd.c_str(), "r");
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

} // namespace os
