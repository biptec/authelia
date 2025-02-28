package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"

	"github.com/authelia/authelia/v4/internal/configuration"
	"github.com/authelia/authelia/v4/internal/utils"
)

func recoverErr(i any) error {
	switch v := i.(type) {
	case nil:
		return nil
	case string:
		return fmt.Errorf("recovered panic: %s", v)
	case error:
		return fmt.Errorf("recovered panic: %w", v)
	default:
		return fmt.Errorf("recovered panic with unknown type: %v", v)
	}
}

func flagsGetUserIdentifiersGenerateOptions(flags *pflag.FlagSet) (users, services, sectors []string, err error) {
	if users, err = flags.GetStringSlice(cmdFlagNameUsers); err != nil {
		return nil, nil, nil, err
	}

	if services, err = flags.GetStringSlice(cmdFlagNameServices); err != nil {
		return nil, nil, nil, err
	}

	if sectors, err = flags.GetStringSlice(cmdFlagNameSectors); err != nil {
		return nil, nil, nil, err
	}

	return users, services, sectors, nil
}

//nolint:gocyclo
func flagsGetRandomCharacters(flags *pflag.FlagSet, flagNameLength, flagNameCharSet, flagNameCharacters string) (r string, err error) {
	var (
		n       int
		charset string
	)

	if n, err = flags.GetInt(flagNameLength); err != nil {
		return "", err
	}

	if n < 1 {
		return "", fmt.Errorf("flag --%s with value '%d' is invalid: must be at least 1", flagNameLength, n)
	}

	useCharSet, useCharacters := flags.Changed(flagNameCharSet), flags.Changed(flagNameCharacters)

	if useCharSet && useCharacters {
		return "", fmt.Errorf("flag --%s and flag --%s are mutually exclusive, only one may be used", flagNameCharSet, flagNameCharacters)
	}

	switch {
	case useCharSet, !useCharSet && !useCharacters:
		var c string

		if c, err = flags.GetString(flagNameCharSet); err != nil {
			return "", err
		}

		switch c {
		case "ascii":
			charset = utils.CharSetASCII
		case "alphanumeric":
			charset = utils.CharSetAlphaNumeric
		case "alphanumeric-lower":
			charset = utils.CharSetAlphabeticLower + utils.CharSetNumeric
		case "alphanumeric-upper":
			charset = utils.CharSetAlphabeticUpper + utils.CharSetNumeric
		case "alphabetic":
			charset = utils.CharSetAlphabetic
		case "alphabetic-lower":
			charset = utils.CharSetAlphabeticLower
		case "alphabetic-upper":
			charset = utils.CharSetAlphabeticUpper
		case "numeric-hex":
			charset = utils.CharSetNumericHex
		case "numeric":
			charset = utils.CharSetNumeric
		case "rfc3986":
			charset = utils.CharSetRFC3986Unreserved
		case "rfc3986-lower":
			charset = utils.CharSetAlphabeticLower + utils.CharSetNumeric + utils.CharSetSymbolicRFC3986Unreserved
		case "rfc3986-upper":
			charset = utils.CharSetAlphabeticUpper + utils.CharSetNumeric + utils.CharSetSymbolicRFC3986Unreserved
		default:
			return "", fmt.Errorf("flag '--%s' with value '%s' is invalid, must be one of 'ascii', 'alphanumeric', 'alphabetic', 'numeric', 'numeric-hex', or 'rfc3986'", flagNameCharSet, c)
		}
	case useCharacters:
		if charset, err = flags.GetString(flagNameCharacters); err != nil {
			return "", err
		}
	}

	return utils.RandomString(n, charset), nil
}

func termReadConfirmation(flags *pflag.FlagSet, name, prompt, confirmation string) (confirmed bool, err error) {
	if confirmed, _ = flags.GetBool(name); confirmed {
		return confirmed, nil
	}

	terminal, fd, state, err := getTerminal(prompt)
	if err != nil {
		return false, err
	}

	defer func(fd int, oldState *term.State) {
		_ = term.Restore(fd, oldState)
	}(fd, state)

	var input string

	if input, err = terminal.ReadLine(); err != nil {
		return false, fmt.Errorf("failed to read from the terminal: %w", err)
	}

	if input != confirmation {
		return false, nil
	}

	return true, nil
}

func getTerminal(prompt string) (terminal *term.Terminal, fd int, state *term.State, err error) {
	fd = int(syscall.Stdin) //nolint:unconvert,nolintlint

	if !term.IsTerminal(fd) {
		return nil, -1, nil, ErrStdinIsNotTerminal
	}

	var width, height int

	if width, height, err = term.GetSize(int(syscall.Stdout)); err != nil { //nolint:unconvert,nolintlint
		return nil, -1, nil, fmt.Errorf("failed to get terminal size: %w", err)
	}

	state, err = term.MakeRaw(fd)
	if err != nil {
		return nil, -1, nil, fmt.Errorf("failed to get terminal state: %w", err)
	}

	c := struct {
		io.Reader
		io.Writer
	}{
		os.Stdin,
		os.Stdout,
	}

	terminal = term.NewTerminal(c, prompt)

	if err = terminal.SetSize(width, height); err != nil {
		return nil, -1, nil, fmt.Errorf("failed to set terminal size: %w", err)
	}

	return terminal, fd, state, nil
}

func termReadPasswordWithPrompt(prompt, flag string) (password string, err error) {
	terminal, fd, state, err := getTerminal("")
	if err != nil {
		if errors.Is(err, ErrStdinIsNotTerminal) {
			switch len(flag) {
			case 0:
				return "", err
			case 1:
				return "", fmt.Errorf("you must either use an interactive terminal or use the -%s flag", flag)
			default:
				return "", fmt.Errorf("you must either use an interactive terminal or use the --%s flag", flag)
			}
		}

		return "", err
	}

	defer func(fd int, oldState *term.State) {
		_ = term.Restore(fd, oldState)
	}(fd, state)

	if password, err = terminal.ReadPassword(prompt); err != nil {
		return "", fmt.Errorf("failed to read the input from the terminal: %w", err)
	}

	return password, nil
}

type XEnvCLIResult int

const (
	XEnvCLIResultCLIExplicit XEnvCLIResult = iota
	XEnvCLIResultCLIImplicit
	XEnvCLIResultEnvironment
)

func loadXEnvCLIConfigValues(cmd *cobra.Command) (configs []string, filters []configuration.FileFilter, err error) {
	var (
		filterNames []string
	)

	if configs, _, err = loadXEnvCLIStringSliceValue(cmd, cmdFlagEnvNameConfig, cmdFlagNameConfig); err != nil {
		return nil, nil, err
	}

	if configs, err = loadXNormalizedPaths(configs); err != nil {
		return nil, nil, err
	}

	if filterNames, _, err = loadXEnvCLIStringSliceValue(cmd, cmdFlagEnvNameConfigFilters, cmdFlagNameConfigExpFilters); err != nil {
		return nil, nil, err
	}

	if filters, err = configuration.NewFileFilters(filterNames); err != nil {
		return nil, nil, fmt.Errorf("error occurred loading configuration: flag '--%s' is invalid: %w", cmdFlagNameConfigExpFilters, err)
	}

	return
}

func loadXNormalizedPaths(paths []string) ([]string, error) {
	var (
		configs, files, dirs []string
		err                  error
	)

	var stat os.FileInfo

	for _, path := range paths {
		if path, err = filepath.Abs(path); err != nil {
			return nil, fmt.Errorf("failed to determine absolute path for '%s': %w", path, err)
		}

		switch stat, err = os.Stat(path); {
		case err == nil && stat.IsDir():
			configs = append(configs, path)
			dirs = append(dirs, path)
		case err == nil:
			configs = append(configs, path)
			files = append(files, path)
		default:
			if os.IsNotExist(err) {
				configs = append(configs, path)
				files = append(files, path)

				continue
			}

			return nil, fmt.Errorf("error occurred stating file at path '%s': %w", path, err)
		}
	}

	for i, file := range files {
		if file, err = filepath.Abs(file); err != nil {
			return nil, fmt.Errorf("failed to determine absolute path for '%s': %w", files[i], err)
		}

		if len(dirs) != 0 {
			filedir := filepath.Dir(file)

			for _, dir := range dirs {
				if filedir == dir {
					return nil, fmt.Errorf("failed to load config directory '%s': the config file '%s' is in that directory which is not supported", dir, file)
				}
			}
		}
	}

	return configs, nil
}

func loadXEnvCLIStringSliceValue(cmd *cobra.Command, envKey, flagName string) (value []string, result XEnvCLIResult, err error) {
	if cmd.Flags().Changed(flagName) {
		value, err = cmd.Flags().GetStringSlice(flagName)

		return value, XEnvCLIResultCLIExplicit, err
	}

	var (
		env string
		ok  bool
	)

	if envKey != "" {
		env, ok = os.LookupEnv(envKey)
	}

	switch {
	case ok && env != "":
		return strings.Split(env, ","), XEnvCLIResultEnvironment, nil
	default:
		value, err = cmd.Flags().GetStringSlice(flagName)

		return value, XEnvCLIResultCLIImplicit, err
	}
}

func newHelpTopic(topic, short, body string) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   topic,
		Short: short,
	}

	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		_ = cmd.Parent().Help()

		fmt.Println()
		fmt.Printf("Help Topic: %s\n\n", topic)
		fmt.Print(body)
		fmt.Print("\n\n")
	})

	return cmd
}
