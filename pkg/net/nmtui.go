package net

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/nmstate/nmstate/rust/src/go/nmstate/v2"
	"github.com/rivo/tview"
)

func NMTUIRunner(app *tview.Application, pages *tview.Pages) func() {
	return func() {
		app.Suspend(func() {
			cmd := exec.Command("nmtui")
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			err := cmd.Run()
			if err != nil {
				panic(err)
			}
		})
		nm := nmstate.New()
		state, err := nm.RetrieveNetState()
		if err != nil {
			panic(err)
		}

		var netState NetState
		if err := json.Unmarshal([]byte(state), &netState); err != nil {
			panic(err)
		}

		//netStatePage, err := modalNetStateJSONPage(&netState, pages)
		netStatePage, err := ModalTreeView(netState, pages)
		pages.AddPage("netstate", netStatePage, true, true)
	}
}
