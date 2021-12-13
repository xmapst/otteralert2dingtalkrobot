package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"otteralter/engine"
	"otteralter/utils"
	"syscall"
)

var rootCmd = &cobra.Command{
	Use:               os.Args[0],
	Version:           VersionIfo(),
	Short:             "A simple otter monitor dingtalk robot alert",
	DisableAutoGenTag: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Title, VersionIfo())
		engine.Run()
	},
}

func init() {
	registerSignalHandlers()
	logrus.SetFormatter(&utils.ConsoleFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	rootCmd.PersistentFlags().StringSliceVarP(&engine.Addr, "connect", "c", nil, "connection rabbitmq address string\n-c 172.16.2.29:2181,172.16.2.29:2182")
	_ = rootCmd.MarkPersistentFlagRequired("connect")
	rootCmd.PersistentFlags().StringVarP(&engine.DingTalkUrl, "url", "u", "", "dingtalk robot url\n-u https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxxxxxx")
	_ = rootCmd.MarkPersistentFlagRequired("url")
	rootCmd.PersistentFlags().StringVarP(&engine.DingTalkSecret, "secret", "s", "", "dingtalk robot secret. optional\n-s SEC000000000000000000000")
	rootCmd.PersistentFlags().StringVarP(&engine.Interval, "interval", "t", "5m", "monitoring interval. optional\n-i 5m")
}

func registerSignalHandlers() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigs
		os.Exit(0)
	}()
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
