package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func moduleAction(command string, cmd *cobra.Command, args []string) {
	c, err := getSession(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := c.CompleteSession(); err != nil {
		fmt.Println(err.Error())
		return
	}

	model := "ir.module.module"

	domain := []interface{}{
		[]interface{}{"name", "in", args},
	}

	var ids []int64
	if err := c.Search(model,
		[]interface{}{domain},
		nil, &ids); err != nil {
		panic(err)
	}

	var result interface{}
	if err := c.DoRequest(command, model,
		[]interface{}{ids},
		map[string]interface{}{}, &result); err != nil {
		panic(err)
	}

	fmt.Println(result)
}
