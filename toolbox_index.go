package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var qOption = []*survey.Question{
	{
		Name: "option",
		Prompt: &survey.Select{
			Message: "What are you doing to Firebase Indexes?:",
			Options: []string{"Adding"},
			Default: "Adding",
		},
	},
}

func addingQs() []*survey.Question {
	var addQs = []*survey.Question{
		{
			Name:     "field",
			Prompt:   &survey.Input{Message: "What is the first field?"},
			Validate: survey.Required,
		},
		{
			Name: "order",
			Prompt: &survey.Select{
				Message: "Choose a order:",
				Options: []string{"ASCENDING", "DESCENDING", "ARRAY"},
				Default: "ASCENDING",
			},
		},
		{
			Name: "more",
			Prompt: &survey.Confirm{
				Message: "Add another field to your collection?",
			},
		},
	}

	var newQs = addQs

	return newQs
}

var qColl = []*survey.Question{
	{
		Name:     "collectionID",
		Prompt:   &survey.Input{Message: "What is the collection ID?"},
		Validate: survey.Required,
	},
	{
		Name:     "field",
		Prompt:   &survey.Input{Message: "What is the first field?"},
		Validate: survey.Required,
	},
	{
		Name: "order",
		Prompt: &survey.Select{
			Message: "Choose a order:",
			Options: []string{"ASCENDING", "DESCENDING", "ARRAY"},
			Default: "ASCENDING",
		},
	},
	{
		Name: "more",
		Prompt: &survey.Confirm{
			Message: "Add another field to your collection?",
		},
	},
}

var qCollMore = []*survey.Question{
	{
		Name:     "field",
		Prompt:   &survey.Input{Message: "What is the next field?"},
		Validate: survey.Required,
	},
	{
		Name: "order",
		Prompt: &survey.Select{
			Message: "Choose a order:",
			Options: []string{"ASCENDING", "DESCENDING", "ARRAY"},
			Default: "ASCENDING",
		},
	},
	{
		Name: "more",
		Prompt: &survey.Confirm{
			Message: "Add another field to your collection?",
		},
	},
}

var qCorrect = []*survey.Question{
	{
		Name: "correct",
		Prompt: &survey.Confirm{
			Message: "All info correct?",
		},
	},
}

var qEnv = []*survey.Question{
	{
		Name: "env",
		Prompt: &survey.Select{
			Message: "Which Environment to change Index for?:",
			Options: []string{"prod", "dev"},
			Default: "dev",
		},
	},
}

// the answers will be written to this struct
type addAnswers struct {
	CollectionID string
	Field        string
	Order        string
}

var magenta = color.New(color.FgHiMagenta).SprintFunc()
var cyan = color.New(color.FgHiCyan).SprintFunc()

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Update/Add/Delete Firebase Indexes",
	Long:  `This command is used for adding of Firebase indexes.`,
	Run: func(cmd *cobra.Command, args []string) {
		aOption := struct {
			Option string
		}{}

		aColl := struct {
			CollectionID string
			Field        string
			Order        string
			More         bool
		}{}

		aCorrect := struct {
			Correct bool
		}{}
		aEnv := struct {
			Env string
		}{}

		addAns := []addAnswers{}

		err5 := survey.Ask(qEnv, &aEnv)
		if err5 != nil {
			fmt.Println(err5.Error())
			return
		}

		err := survey.Ask(qOption, &aOption)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if aOption.Option == "Adding" {

			err2 := survey.Ask(qColl, &aColl)
			if err2 != nil {
				fmt.Println(err.Error())
				return
			}

			addAns = append(addAns, addAnswers{aColl.CollectionID, aColl.Field, aColl.Order})

			for aColl.More {
				err3 := survey.Ask(qCollMore, &aColl)
				if err3 != nil {
					fmt.Println(err.Error())
					return
				}

				addAns = append(addAns, addAnswers{aColl.CollectionID, aColl.Field, aColl.Order})
			}

			fmt.Println(magenta("Environment: "), aEnv.Env)
			fmt.Println(magenta("Collection: "), addAns[0].CollectionID)

			for i := 0; i < len(addAns); i++ {
				fmt.Print(magenta("Field ", i+1, ": "), addAns[i].Field, "\n")
				fmt.Println(magenta("Order:  "), addAns[i].Order)
			}

			err4 := survey.Ask(qCorrect, &aCorrect)
			if err4 != nil {
				fmt.Println(err4.Error())
				return
			}

			if aCorrect.Correct {
				currentUser, err := user.Current()
				if err != nil {
					log.Fatalf(err.Error())
				}

				username := currentUser.Username
				now := time.Now()

				cmdGitClone := exec.Command("git", "clone", "https://github.com/pocketrn/terraform-core-functions")
				cmdGitClone.Dir = "/Users/" + username + "/Downloads"
				cmdGitClone.Stdout = os.Stdout
				cmdGitClone.Stderr = os.Stderr
				cmdGitClone.Run()

				if aEnv.Env == "dev" {
					cmdGitBranch := exec.Command("git", "checkout", "-b", "index-dev-"+now.Format("2006-01-02"))
					cmdGitBranch.Dir = "/Users/" + username + "/Downloads/terraform-core-functions"
					cmdGitBranch.Stdout = os.Stdout
					cmdGitBranch.Stderr = os.Stderr
					cmdGitBranch.Run()

					nowTime := now.Format("2006-01-02") + "_" + now.Format("15:04:05")

					file, err := os.OpenFile("/Users/"+username+"/Downloads/terraform-core-functions/modules/firebaseDev/indexes.tf", os.O_APPEND|os.O_WRONLY, 0644)

					if err != nil {
						fmt.Println("Could not open indexes.tf")
						return
					}

					defer file.Close()

					_, err2 := file.WriteString(`
resource "google_firestore_index" "index_fields_` + nowTime + `" {
  project    = var.project
  collection = "` + addAns[0].CollectionID + `"
`)

					if err2 != nil {
						fmt.Println("Could not write text to indexes.tf")

					} else {
						fmt.Println("New Index has been Added.")
					}

					for i := 0; i < len(addAns); i++ {
						file2, err2 := os.OpenFile("/Users/"+username+"/Downloads/terraform-core-functions/modules/firebaseDev/indexes.tf", os.O_APPEND|os.O_WRONLY, 0644)

						if err2 != nil {
							fmt.Println("Could not open indexes.tf")
							return
						}

						defer file2.Close()

						_, err3 := file2.WriteString(`
  fields {
    field_path = "` + addAns[i].Field + `"
    order      = "` + addAns[i].Order + `"
  }
`)

						if err3 != nil {
							fmt.Println("Could not write text to indexes.tf")
						}

					}
					file3, err4 := os.OpenFile("/Users/"+username+"/Downloads/terraform-core-functions/modules/firebaseDev/indexes.tf", os.O_APPEND|os.O_WRONLY, 0644)

					if err4 != nil {
						fmt.Println("Could not open indexes.tf")
						return
					}

					defer file3.Close()

					_, err5 := file3.WriteString(`
  fields {
    field_path = local.default_last_value
    order      = local.default_last_order
  }

  depends_on = [
    google_project_service.firestore_api
  ]
}
`)

					if err5 != nil {
						fmt.Println("Could not write text to indexes.tf")
					}

					cmdGitCommit := exec.Command("git", "commit", "-a", "-m", "Firestore Index Added")
					cmdGitCommit.Dir = "/Users/" + username + "/Downloads/terraform-core-functions"
					cmdGitCommit.Stdout = os.Stdout
					cmdGitCommit.Stderr = os.Stderr
					cmdGitCommit.Run()

					cmdGitPush := exec.Command("git", "push", "-u", "origin", "index-dev-"+now.Format("2006-01-02"))
					cmdGitPush.Dir = "/Users/" + username + "/Downloads/terraform-core-functions"
					cmdGitPush.Run()

					cmdGitPR := exec.Command("gh", "pr", "create", "--title", "Firestore Index Added Dev "+now.Format("2006-01-02"), "--body", "Auto Generated PR by Engineer Toolbox", "--reviewer", "BagelHole", "-a", "@me")
					cmdGitPR.Dir = "/Users/" + username + "/Downloads/terraform-core-functions"
					cmdGitPR.Stdout = os.Stdout
					cmdGitPR.Stderr = os.Stderr
					cmdGitPR.Run()

					cmdDel := exec.Command("rm", "-r", "terraform-core-functions")
					cmdDel.Dir = "/Users/" + username + "/Downloads"
					cmdDel.Run()

					fmt.Println(cyan("Copy the above link into the dev-team Slack channel and tag Toby Miller or Ryan Saunders"))

				} else {
					cmdGitBranch := exec.Command("git", "checkout", "-b", "index-prod-"+now.Format("2006-01-02"))
					cmdGitBranch.Dir = "/Users/" + username + "/Downloads/terraform-core-functions"
					cmdGitBranch.Stdout = os.Stdout
					cmdGitBranch.Stderr = os.Stderr
					cmdGitBranch.Run()

					nowTime := now.Format("2006-01-02") + "_" + now.Format("15:04:05")

					file, err := os.OpenFile("/Users/"+username+"/Downloads/terraform-core-functions/modules/firebase/indexes.tf", os.O_APPEND|os.O_WRONLY, 0644)

					if err != nil {
						fmt.Println("Could not open indexes.tf")
						return
					}

					defer file.Close()

					_, err2 := file.WriteString(`
resource "google_firestore_index" "index_fields_` + nowTime + `" {
  project    = var.project
  collection = "` + addAns[0].CollectionID + `"
`)

					if err2 != nil {
						fmt.Println("Could not write text to indexes.tf")

					} else {
						fmt.Println("New Index has been Added.")
					}

					for i := 0; i < len(addAns); i++ {
						file2, err2 := os.OpenFile("/Users/"+username+"/Downloads/terraform-core-functions/modules/firebase/indexes.tf", os.O_APPEND|os.O_WRONLY, 0644)

						if err2 != nil {
							fmt.Println("Could not open indexes.tf")
							return
						}

						defer file2.Close()

						_, err3 := file2.WriteString(`
  fields {
    field_path = "` + addAns[i].Field + `"
    order      = "` + addAns[i].Order + `"
  }
`)

						if err3 != nil {
							fmt.Println("Could not write text to indexes.tf")
						}

					}
					file3, err4 := os.OpenFile("/Users/"+username+"/Downloads/terraform-core-functions/modules/firebase/indexes.tf", os.O_APPEND|os.O_WRONLY, 0644)

					if err4 != nil {
						fmt.Println("Could not open indexes.tf")
						return
					}

					defer file3.Close()

					_, err5 := file3.WriteString(`
  fields {
    field_path = local.default_last_value
    order      = local.default_last_order
  }

  depends_on = [
    google_project_service.firestore_api
  ]
}
`)

					if err5 != nil {
						fmt.Println("Could not write text to indexes.tf")
					}

					cmdGitCommit := exec.Command("git", "commit", "-a", "-m", "Firestore Index Added")
					cmdGitCommit.Dir = "/Users/" + username + "/Downloads/terraform-core-functions"
					cmdGitCommit.Stdout = os.Stdout
					cmdGitCommit.Stderr = os.Stderr
					cmdGitCommit.Run()

					cmdGitPush := exec.Command("git", "push", "-u", "origin", "index-prod-"+now.Format("2006-01-02"))
					cmdGitPush.Dir = "/Users/" + username + "/Downloads/terraform-core-functions"
					cmdGitPush.Run()

					cmdGitPR := exec.Command("gh", "pr", "create", "--title", "Firestore Index Added Prod "+now.Format("2006-01-02"), "--body", "Auto Generated PR by Engineer Toolbox", "--reviewer", "BagelHole", "-a", "@me")
					cmdGitPR.Dir = "/Users/" + username + "/Downloads/terraform-core-functions"
					cmdGitPR.Stdout = os.Stdout
					cmdGitPR.Stderr = os.Stderr
					cmdGitPR.Run()

					cmdDel := exec.Command("rm", "-r", "terraform-core-functions")
					cmdDel.Dir = "/Users/" + username + "/Downloads"
					cmdDel.Run()

					fmt.Println(cyan("→ Copy the above link into the dev-team Slack channel and tag Toby Miller or Ryan Saunders"))
				}

			}

		}

	},
}

func init() {
	rootCmd.AddCommand(indexCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// indexCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// indexCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
