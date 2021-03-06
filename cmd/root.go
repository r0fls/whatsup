// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

var cfgFile string
var route string
var period string
var method string
var pdkey string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "whatsup",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: rootRun,
}

func rootRun(cmd *cobra.Command, args []string) {
	var err error
	var resp *http.Response

	request, err := http.NewRequest(strings.ToUpper(method), route, nil)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	duration, err := time.ParseDuration(period)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Track triggered incidents.
	incidents := make(map[string]bool)

	for _ = range time.Tick(duration) {
		// TODO: this should run in a goroutine
		resp, err = http.DefaultClient.Do(request)
		if err != nil || resp.StatusCode >= 400 {
			log.Printf("Failed to reach route: %s", route)
			event := pagerduty.Event{
				Type:        "trigger",
				ServiceKey:  pdkey,
				Description: route,
				IncidentKey: route,
			}
			pdresp, err := pagerduty.CreateEvent(event)
			if err != nil {
				fmt.Println(pdresp)
				fmt.Println("ERROR:", err.Error())
			}
			incidents[route] = true
		} else {
			// This defaults to false
			if incidents[route] {
				// Resolve incident
				event := pagerduty.Event{
					Type:        "resolve",
					ServiceKey:  pdkey,
					IncidentKey: route,
				}
				pdresp, err := pagerduty.CreateEvent(event)
				if err != nil {
					fmt.Println(pdresp)
					fmt.Println("ERROR:", err.Error())
				}
				incidents[route] = false
			}
		}
	}
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "", "config file (default is $HOME/.whatsup.yaml)")
	// TODO: this should take a list of routes
	RootCmd.PersistentFlags().StringVarP(&route, "route", "r", "", "route to check, e.g. www.example.com")
	RootCmd.PersistentFlags().StringVarP(&period, "period", "p", "", "Period of frequency.")
	RootCmd.PersistentFlags().StringVarP(&method, "method", "m", "", "Period of frequency.")
	RootCmd.PersistentFlags().StringVarP(&pdkey, "pdkey", "k", "", "Period of frequency.")

	viper.BindPFlag("route", RootCmd.PersistentFlags().Lookup("route"))
	viper.BindPFlag("period", RootCmd.PersistentFlags().Lookup("period"))
	viper.BindPFlag("method", RootCmd.PersistentFlags().Lookup("method"))
	viper.BindPFlag("pdkey", RootCmd.PersistentFlags().Lookup("pdkey"))

	if period == "" {
		period = "60s"
	}

	if method == "" {
		method = "GET"
	}

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".whatsup" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".whatsup")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
