package build

import (
	"fmt"
	"os/exec"
	"path"

	"github.com/golang/glog"
)

func validateExtension(filename, extension string) error {
	if path.Ext(filename) != extension {
		return fmt.Errorf("expected filename with %q extension; got filename %q", extension, filename)
	}
	return nil
}

// Example:
// $ clang++ -std=c++20 -x c++-module -fprebuilt-module-path=out -Ientt/single_include -ISDL/include game_impl.cc --precompile -o out/game-game_impl.pcm
func CompileCCToPCM(inputFilename string, flags []string, prebuiltModulePath, outputFilename string) ([]byte, error) {
	glog.V(1).Infof("Compiling %q to %q...", inputFilename, outputFilename)

	if err := validateExtension(inputFilename, ".cc"); err != nil {
		return nil, err
	}
	if err := validateExtension(outputFilename, ".pcm"); err != nil {
		return nil, err
	}

	args := []string{"-std=c++20", "-x", "c++-module", fmt.Sprintf("-fprebuilt-module-path=%s", prebuiltModulePath), inputFilename, "--precompile", "-o", outputFilename}
	args = append(args, flags...)
	cmd := exec.Command("clang++", args...)

	glog.V(1).Infof("Calling %s", cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run %s: %s", cmd, output)
	}

	return output, nil
}

// Example:
// $ clang++ -std=c++20 -fprebuilt-module-path=out -c out/game-game_impl.pcm -o out/game-game_impl.o
func CompilePCMToObj(inputFilename string, prebuiltModulePath, outputFilename string) ([]byte, error) {
	glog.V(1).Infof("Compiling %q to %q...", inputFilename, outputFilename)

	args := []string{"-std=c++20", fmt.Sprintf("-fprebuilt-module-path=%s", prebuiltModulePath), "-c", inputFilename, "-o", outputFilename}
	cmd := exec.Command("clang++", args...)

	glog.V(1).Infof("Calling %s", cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run %s: %s", cmd, output)
	}

	return output, nil
}

// Example:
//
//	clang++ -std=c++20 -fprebuilt-module-path=out -o out/program \
//	  -Wl,-rpath,SDL/build \
//	  -LSDL/build -lSDL3 \
//	  out/arr-arr_impl.o \
//	  ...
func LinkObjsToExecutable(inputFilenames, flags []string, prebuiltModulePath, outputFilename string) ([]byte, error) {
	if len(inputFilenames) == 0 {
		return nil, fmt.Errorf("no object files (.o) to link")
	}

	args := []string{"-std=c++20", fmt.Sprintf("-fprebuilt-module-path=%s", prebuiltModulePath), "-o", outputFilename}
	args = append(args, flags...)
	args = append(args, inputFilenames...)
	cmd := exec.Command("clang++", args...)

	glog.V(1).Infof("Building executable %q with %s", outputFilename, cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run %s: %s", cmd, output)
	}

	return output, nil
}
