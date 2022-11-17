package main

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/nmstate/nmstate/rust/src/go/nmstate/v2"
	"github.com/openshift-agent-team/tui/pkg/forms"
	tuiNet "github.com/openshift-agent-team/tui/pkg/net"
	"github.com/openshift-agent-team/tui/pkg/newt"
	"github.com/rivo/tview"
	//"golang.org/x/exp/slices"
)

const (
	QUIT      string = "Quit"
	CONFIGURE string = "Configure"
	YES       string = "Yes"
	NO        string = "No"
)

func doneView(app *tview.Application, pages *tview.Pages) func(int, string) {
	return func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Quit" {
			app.Stop()
		} else {
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

			var netState tuiNet.NetState

			if err := json.Unmarshal([]byte(state), &netState); err != nil {
				panic(err)
			}

			//netStatePage, err := modalNetStateJSONPage(&netState, pages)
			netStatePage, err := tuiNet.ModalTreeView(netState, pages)
			pages.AddPage("netstate", netStatePage, true, true)
		}
	}
}

func node0Handler(app *tview.Application, pages *tview.Pages) func(int, string) {
	return func(buttonIndex int, buttonLabel string) {
		if buttonLabel == YES {
			// TODO: Print addressing, offer to configure, done

			node0Form := forms.Node0Form(app, pages)
			pages.AddPage("node0Form", node0Form, true, true)
		} else {
			regNodeForm := forms.RegNodeModalForm(app, pages)
			pages.AddPage("regNodeConfig", regNodeForm, true, true)
		}
	}
}

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()

	background := tview.NewBox().
		SetBorder(false).
		SetBackgroundColor(newt.ColorBlue)

	node0 := tview.NewModal().
		SetText("Do you wish for this node to be the one that runs the installation service (Only one node may perform this function)?").
		SetTextColor(tcell.ColorBlack).
		SetDoneFunc(node0Handler(app, pages)).
		SetBackgroundColor(newt.ColorGray).
		SetButtonTextColor(tcell.ColorBlack).
		SetButtonBackgroundColor(tcell.ColorDarkGray)

	node0Buttons := []string{YES, NO}
	node0.AddButtons(node0Buttons)

	pages.AddPage("background", background, true, true).
		AddPage("Node0", node0, true, true)

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}
