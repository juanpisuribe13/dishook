// Copyright © 2021 Juanpis
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jochasinga/requests"
	"github.com/spf13/cobra"
)

var (
	avatar_url string
	username   string
	message    string
	tts        bool

	author_id        bool
	username_content bool
	discriminator    bool
	components       bool
	message_content  bool
	channel_id       bool
	edited_timestamp bool
	embeds           bool
	flags            bool
	message_id       bool

	webhook_type bool
	webhook_id   bool
	timestamp    bool

	has_avatar_url    bool
	is_bot            bool
	mentions_everyone bool
	mentions_roles    bool
	is_pinned         bool
	has_tts           bool
)

func init() {
	// Execute Commands
	execute_cmd.Flags().StringVarP(&avatar_url, "avatar-url", "a", "", "sets webhook's profile picture")
	execute_cmd.Flags().StringVarP(&message, "message", "m", "", "sets message")
	execute_cmd.Flags().StringVarP(&username, "username", "u", "", "sets username of webhook")
	execute_cmd.Flags().BoolVarP(&tts, "tts", "t", false, "sets if tts should be enabled or not")

	// Get Commands
	// json info | author:map
	get_cmd.Flags().BoolVarP(&has_avatar_url, "avatar-url", "a", false, "avatar link of the webhook")
	get_cmd.Flags().BoolVarP(&is_bot, "bot", "b", false, "returns if obtained webhook is bot")
	get_cmd.Flags().BoolVarP(&discriminator, "discriminator", "d", false, "returns discriminator")
	get_cmd.Flags().BoolVarP(&author_id, "author-id", "", false, "returns ID of webhook")
	get_cmd.Flags().BoolVarP(&username_content, "username", "u", false, "name used for webhook")

	// message info
	get_cmd.Flags().BoolVarP(&message_content, "message", "m", false, "message sent")
	get_cmd.Flags().BoolVarP(&message_id, "message-id", "s", false, "message ID")
	get_cmd.Flags().BoolVarP(&channel_id, "channel-id", "c", false, "channel ID")

	get_cmd.Flags().BoolVarP(&mentions_everyone, "mentions-everyone", "e", false, "returns if everyone is mentioned")
	get_cmd.Flags().BoolVarP(&mentions_roles, "mention-roles", "r", false, "returns if roles are mentioned")
	get_cmd.Flags().BoolVarP(&is_pinned, "pinned", "p", false, "returns if message is pinned")
	get_cmd.Flags().BoolVarP(&timestamp, "timestamp", "", false, "returns the time webhook was executed")
	get_cmd.Flags().BoolVarP(&has_tts, "tts", "t", false, "returns if TTS was used")

	// webhook info
	get_cmd.Flags().BoolVarP(&webhook_id, "webhook-id", "", false, "webhook ID")
	get_cmd.Flags().BoolVarP(&webhook_type, "webhook-type", "", false, "Webhook type")

	// misc
	get_cmd.Flags().BoolVarP(&components, "components", "", false, "components included with the message")
	get_cmd.Flags().BoolVarP(&edited_timestamp, "edited-timestamp", "", false, "time when message was edited")
	get_cmd.Flags().BoolVarP(&embeds, "embeds", "", false, "array of message embeds/components")
	get_cmd.Flags().BoolVarP(&flags, "flags", "", false, "name of webhook")

	// Edit Command (lol)
	edit_cmd.Flags().StringVarP(&message, "message", "m", "", "sets message you wanna edit")

	root_cmd.AddCommand(get_cmd, execute_cmd, edit_cmd, delete_cmd)
}

// Allocates all strings from array starting from
// given argument position into one variable with
// blank spaces
//
// Used when input has no flags set.
func merge_strings(args []string, arg_pos int) string {
	var str string
	for i := arg_pos; i < len(args); i++ {
		str = fmt.Sprintf("%s %s", str, args[i])
	}
	return strings.TrimSpace(str)
}

// Semds an HTTP request method to Discord's webhook
// with provided URL and JSON map.
// Automatically marshalls the JSON map.
//
// Supported HTTP Methods: POST, PATCH
func request_HTTP(http_method string, URL string, json_map map[string]string) {
	json_value, _ := json.Marshal(json_map)

	switch http_method {
	case "POST":
		resp, err := requests.Post(URL, "application/json", bytes.NewBuffer(json_value))
		_, _ = resp, err
	case "PATCH":
		resp, err := requests.Patch(URL, "application/json", bytes.NewBuffer(json_value))
		_, _ = resp, err
	}
}

// Checks if provided URL matches
// Discord's webhook API URL.
func is_token_valid(url string) bool {
	if url[0:33] == "https://discord.com/api/webhooks/" {
		url_r, err := requests.Get(url)
		url_code := url_r.StatusCode
		if url_code != 401 || err != nil {
			return true
		}
	}

	// Passed thru all checks without any true statement
	return false

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		// what
	// 		// i might tickle around and see what this does
	// 		fmt.Printf("'%s' is not a valid webhook URL.", url)
	// 		os.Exit(0)
	// 	}
	// }()
}

// Checks if given value doesn't pass set limit (2000).
func is_max(msg string) bool { return len(msg) >= 2000 }

// Panics when an non-nil error is given.
//
// Error value is first parsed by fmt.Errorf(), then a red-colored
// ERROR string is merged with the given value.
//
// Example: return fmt.Errorf("message flag required") -> ERROR: message flag required
func ManageError(err error) {
	if err != nil {
		red := color.New(color.FgRed).Sprintf("ERROR:")
		log.SetFlags(0)
		log.Fatal(fmt.Errorf("%s %s", red, err))
	}
}

// Automatically refers to execute() or help if no argument is parsed.
var root_cmd = &cobra.Command{
	Use:  "dishook [url] [message]\n  dishook [url] [flags]",
	Args: cobra.MinimumNArgs(2),

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			// return fmt.Errorf("no arguments given")
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		if !is_token_valid(url) {
			ManageError(fmt.Errorf("'%s' not a valid webhook token", args[0]))
		}

		err := execute(args)
		ManageError(err)

		return nil
	},
}

func Execute() {
	if err := root_cmd.Execute(); err != nil {
		// fmt.Fprintln(os.Stderr, err)
		os.Exit(0)
	}
}
