package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	
	"github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/spf13/cobra"
)

type Config struct {
	OktaDomain     string
	ConfigFilePath string
	OrgName        string
	Client         *okta.APIClient
}

var (
	configFile  string
	outputDir   string
	inputDir    string
)

var rootCmd = &cobra.Command{
	Use:   "envsync",
	Short: "A tool for backing up and restoring Okta developer environments",
	Long: `envsync is a tool for backing up and restoring Okta developer environments.
It is only designed and tested for use with Okta developer accounts.`,
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup an Okta developer environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := LoadConfig(configFile)
		if err != nil {
			return err
		}
		
		return PerformBackup(cfg, outputDir)
	},
}

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore an Okta developer environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := LoadConfig(configFile)
		if err != nil {
			return err
		}
		
		return PerformRestore(cfg, inputDir)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(restoreCmd)
	
	backupCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to Okta config file")
	backupCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Directory to store backup files")
	
	restoreCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to Okta config file")
	restoreCmd.Flags().StringVarP(&inputDir, "input", "i", "", "Directory containing backup files")
	restoreCmd.MarkFlagRequired("input")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".okta", "okta.yaml")
}

// PrepareOktaCliArgs prepares arguments for okta-cli-client with config flag if specified
func PrepareOktaCliArgs(cfg *Config, args ...string) []string {
	if cfg.ConfigFilePath != "" {
		return append([]string{"--config", cfg.ConfigFilePath}, args...)
	}
	return args
}

var devOrgPattern = regexp.MustCompile(`(?i)((?:dev|trial)-\d+)\.okta`)

func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = DefaultConfigPath()
	}

	fmt.Printf("Using configuration from %s\n", configPath)
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file %s does not exist", configPath)
	}
	
	domain, orgName, err := scanConfigForDevDomain(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	
	if domain == "" {
		return nil, fmt.Errorf("this tool is only designed for Okta developer accounts (dev-*.okta)")
	}
	
	config := &Config{
		ConfigFilePath: configPath,
		OktaDomain:     domain,
		OrgName:        orgName,
	}
	
	return config, nil
}

// scanConfigForDevDomain scans a config file to find an Okta developer domain
// Returns the full domain and the org name (dev-XXXXX)
func scanConfigForDevDomain(filePath string) (string, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		
		matches := devOrgPattern.FindStringSubmatch(line)
		if len(matches) > 1 {
			// matches[0] contains the full match, e.g., "dev-123456.okta.com"
			// matches[1] contains the org name part, e.g., "dev-123456"
			domain := matches[0]
			orgName := matches[1]
			return domain, orgName, nil
		}
	}
	
	if err := scanner.Err(); err != nil {
		return "", "", err
	}
	
	return "", "", nil
}
