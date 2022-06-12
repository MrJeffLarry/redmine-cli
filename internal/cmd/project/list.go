package project

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/MrJeffLarry/redmine-cli/internal/api"
	"github.com/MrJeffLarry/redmine-cli/internal/config"
	"github.com/MrJeffLarry/redmine-cli/internal/print"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
)

const (
	FLAG_ORDER     = "order"
	FLAG_ORDER_ASC = "asc"
	FLAG_ORDER_DES = "des"

	FLAG_LIMIT = "limit"
)

func countDigi(i int64) (count int) {
	for i > 0 {
		i = i / 10
		count++
	}
	return
}

func parseFlags(cmd *cobra.Command, path string) string {
	//	order, err := cmd.Flags().GetString(FLAG_ORDER)
	return path
}

func displayListGET(r *config.Red_t, cmd *cobra.Command, path string) {
	var err error
	var body []byte
	var status int
	var idLen int
	var nameLen int
	var oldParent string
	var parent string
	var parentLevel int

	projects := projects{}

	path = parseFlags(cmd, path)

	if body, status, err = api.ClientGET(r, path); err != nil {
		fmt.Println(status, "Could not get response from client", err)
		return
	}

	print.PrintDebug(r, status, string(body))

	if err := json.Unmarshal(body, &projects); err != nil {
		fmt.Println(err)
		fmt.Println(status, "Could not parse and read response from server")
		return
	}

	for _, project := range projects.Projects {
		iLen := countDigi(project.ID)
		nLen := len(project.Name)

		if iLen > idLen {
			idLen = iLen
		}

		if nLen > nameLen {
			nameLen = nLen
		}
	}

	fmt.Printf("%s  %s\n",
		"ID "+strings.Repeat(" ", int(math.Abs(float64(idLen-len("ID"))))),
		"NAME"+strings.Repeat(" ", int(math.Abs(float64(nameLen-len("NAME"))))),
	)

	for _, project := range projects.Projects {
		iLeft := idLen - countDigi(project.ID)
		nLeft := nameLen - len(project.Name)

		idPad := strings.Repeat(" ", iLeft)
		name := project.Name + strings.Repeat(" ", nLeft)
		pName := project.Parent.Name

		if len(pName) > 0 && pName == parent {
			// same level do nothing
		} else if len(pName) > 0 && pName != parent {
			if oldParent == pName {
				parentLevel--
			} else {
				parentLevel++
			}
			oldParent = parent
			parent = pName
		} else {
			parent = pName
			parentLevel = 0
		}

		if parentLevel > 0 {
			name = text.FgBlue.Sprint(project.Name + strings.Repeat(" ", nLeft))
		}

		fmt.Printf("%s%s %s %s\n", print.PrintID(project.ID), idPad, strings.Repeat(" ‣", parentLevel), name)
	}
	fmt.Printf("--- projects %d to %d (Total %d) ----\n",
		projects.Offset,
		projects.Limit,
		projects.TotalCount,
	)
}

func cmdProjectList(r *config.Red_t) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List projects",
		Long:  "List all projects",
		Run: func(cmd *cobra.Command, args []string) {
			displayListGET(r, cmd, "/projects.json?limit=1000")
		},
	}

	// All
	cmd.AddCommand(&cobra.Command{
		Use:   "all",
		Short: "List all projects",
		Long:  "List all projects",
		Run: func(cmd *cobra.Command, args []string) {
			displayListGET(r, cmd, "/projects.json")
		},
	})

	cmd.PersistentFlags().String("order", "", "Order on id_ASC or id_DES")
	cmd.PersistentFlags().String("sort", "", "")

	return cmd
}
