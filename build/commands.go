package build

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/golang/glog"
)

func validateExtension(filename, extension string) error {
	if path.Ext(filename) != extension {
		return fmt.Errorf("expected filename with %q extension; got filename %q", extension, filename)
	}
	return nil
}

func hashFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open %q: %v", filename, err)
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

type Command struct {
	InputFilenames []string
	OutputFilename string
	Cmd            *exec.Cmd
}

func runCommand(outputDirectory string, command Command) ([]byte, error) {
	hashes := make([]string, 0, len(command.InputFilenames))
	for _, inputFilename := range command.InputFilenames {
		hash, err := hashFile(inputFilename)
		if err != nil {
			return nil, err
		}

		hashes = append(hashes, hash)
	}

	cachedCommand := fmt.Sprintf("env BPL_INPUT_HASHES=%s %s", strings.Join(hashes, ","), command.Cmd)
	stampBase := fmt.Sprintf("%x.stamp", sha1.Sum([]byte(cachedCommand)))
	stampFilename := path.Join(outputDirectory, stampBase)

	cached := false
	{
		outputFileHash, err := os.ReadFile(stampFilename)
		if os.IsNotExist(err) {
			// Nothing to do.
		} else if err != nil {
			return nil, err
		} else {
			gotHash, err := hashFile(command.OutputFilename)
			if err != nil {
				return nil, err
			}

			cached = gotHash == string(outputFileHash)
		}
	}

	if cached {
		glog.V(1).Infof("Already cached %s", command.Cmd)
		return nil, nil
	}

	glog.V(1).Infof("Calling %s", command.Cmd)

	output, err := command.Cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run %s: %s", command.Cmd, output)
	}

	{
		outputFileHash, err := hashFile(command.OutputFilename)
		if err != nil {
			return nil, err
		}

		if err := os.WriteFile(stampFilename, []byte(outputFileHash), 0644); err != nil {
			return nil, err
		}
	}

	return output, nil
}

// Example:
// $ clang++ -std=c++20 -fprebuilt-module-path=out -Ientt/single_include -ISDL/include game_impl.ccm --precompile -o out/game-game_impl.pcm
func CompileCCMToPCMCommand(inputFilename string, flags []string, outputFilename string) (Command, error) {
	if err := validateExtension(inputFilename, ".ccm"); err != nil {
		return Command{}, err
	}
	if err := validateExtension(outputFilename, ".pcm"); err != nil {
		return Command{}, err
	}

	prebuiltModulePath := path.Dir(outputFilename)

	args := []string{"-std=c++20", fmt.Sprintf("-fprebuilt-module-path=%s", prebuiltModulePath), inputFilename, "--precompile", "-o", outputFilename}
	args = append(args, flags...)
	return Command{[]string{inputFilename}, outputFilename, exec.Command("clang++", args...)}, nil
}

// Example:
// $ clang++ -std=c++20 -fprebuilt-module-path=out -c out/game-game_impl.pcm -o out/game-game_impl.o
func CompilePCMToObjCommand(inputFilename string, outputFilename string) (Command, error) {
	if err := validateExtension(inputFilename, ".pcm"); err != nil {
		return Command{}, err
	}
	if err := validateExtension(outputFilename, ".o"); err != nil {
		return Command{}, err
	}

	prebuiltModulePath := path.Dir(outputFilename)

	args := []string{"-std=c++20", fmt.Sprintf("-fprebuilt-module-path=%s", prebuiltModulePath), "-c", inputFilename, "-o", outputFilename}
	return Command{[]string{inputFilename}, outputFilename, exec.Command("clang++", args...)}, nil
}

// Example:
//
//	clang++ -std=c++20 -o out/program \
//	  -Wl,-rpath,SDL/build \
//	  -LSDL/build -lSDL3 \
//	  out/arr-arr_impl.o \
//	  ...
func LinkObjsToExecutable(inputFilenames, flags []string, outputFilename string) (Command, error) {
	if len(inputFilenames) == 0 {
		return Command{}, fmt.Errorf("no object files (.o) to link")
	}

	for _, inputFilename := range inputFilenames {
		if err := validateExtension(inputFilename, ".o"); err != nil {
			return Command{}, err
		}
	}

	args := []string{"-std=c++20", "-o", outputFilename}
	args = append(args, flags...)
	args = append(args, inputFilenames...)
	return Command{inputFilenames, outputFilename, exec.Command("clang++", args...)}, nil
}
