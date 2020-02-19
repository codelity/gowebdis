/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/codelity/gowebdis/api"
	"github.com/codelity/gowebdis/internal/gowebdis"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start gowebdis",
	Run:   runStartCmd,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().String("host", "", "Redis host list seperated by comma")
	startCmd.Flags().String("master-name", "", "master name of sentinel")
	startCmd.Flags().String("sentinel-address", "", "master name of sentinel")
	startCmd.Flags().String("password", "", "Conection password")
	startCmd.Flags().Int("db", 0, "Database number")
	startCmd.Flags().Int("max-retries", -1, "Maximum retries")
	startCmd.Flags().Int("pool-size", 10, "Pool size")
	startCmd.Flags().Int("min-idle-conns", 3, "Minimum pool size")
	startCmd.Flags().Int("min-retry-backoff", -1, "Minimum retry backoff")
	startCmd.Flags().Int("max-retry-backoff", -1, "Maximum retry backoff")
	startCmd.Flags().Int("dial-timeout", -1, "Dial timeout in seconds")
	startCmd.Flags().Int("read-timeout", -1, "Read timeout in seconds")
	startCmd.Flags().Int("write-timeout", -1, "Write timeout")
	startCmd.Flags().Int("max-conn-age", -1, "Maximum connection age")
	startCmd.Flags().Int("pool-timeout", -1, "Pool timeout")
	startCmd.Flags().Int("idle-timeout", 900, "Idle timeout")
	startCmd.Flags().Int("idle-check-frequency", 900, "Idle check frequency")
	viper.BindPFlag("host", startCmd.Flags().Lookup("host"))
	startCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		viper.BindEnv(flag.Name)
		viper.BindPFlag(flag.Name, flag)
	})
}

func runStartCmd(cmd *cobra.Command, args []string) {
	err := gowebdis.InitConnectionSetting(cmd)
	if err != nil {
		log.Error("[ERROR] " + err.Error())
	} else {
		api.StartServer()
	}

}
