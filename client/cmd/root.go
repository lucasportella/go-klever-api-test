package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	postpb "github.com/roneycharles/klever/third_party/gen"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var cfgFile string

var client postpb.PostServiceClient
var requestCtx context.Context
var requestOpts grpc.DialOption

var rootCmd = &cobra.Command{
	Use:   "postclient",
	Short: "a gRPC client to communicate with the PostService server",
	Long: `a gRPC client to communicate with the PostService server.
	You can use this client to create and read posts.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.postclient.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	fmt.Println("Starting Post Service Client")

	requestCtx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	requestOpts = grpc.WithInsecure()

	conn, err := grpc.Dial("localhost:50051", requestOpts)
	if err != nil {
		log.Fatalf("Unable to establish client connection to localhost:50051: %v", err)
	}

	client = postpb.NewPostServiceClient(conn)
}

func initConfig() {
	if cfgFile != "" {

		viper.SetConfigFile(cfgFile)
	} else {

		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".postclient")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
