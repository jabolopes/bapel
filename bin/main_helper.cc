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


namespace cli {


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















} // namespace cli
