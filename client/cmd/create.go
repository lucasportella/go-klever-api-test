package cmd

import (
	"context"
	"fmt"

	postpb "github.com/roneycharles/klever/third_party/gen"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new post",
	Long: `Create a new post on the server through gRPC. 
	
	A post requires an Title and Content.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		title, err := cmd.Flags().GetString("title")
		content, err := cmd.Flags().GetString("content")
		if err != nil {
			return err
		}
		post := &postpb.Post{
			Title:   title,
			Content: content,
		}
		res, err := client.CreatePost(
			context.TODO(),
			&postpb.CreatePostRequest{
				Post: post,
			},
		)
		if err != nil {
			return err
		}
		fmt.Printf("Post created: %s\n", res.Post.Id)
		return nil
	},
}

func init() {
	createCmd.Flags().StringP("title", "t", "", "A title for the post")
	createCmd.Flags().StringP("content", "c", "", "The content for the post")
	createCmd.MarkFlagRequired("title")
	createCmd.MarkFlagRequired("content")
	rootCmd.AddCommand(createCmd)
}
