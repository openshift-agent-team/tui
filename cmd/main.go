package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/nmstate/nmstate/rust/src/go/nmstate/v2"
	tuiNet "github.com/openshift-agent-team/tui/pkg/net"
	"github.com/openshift-agent-team/tui/pkg/newt"
	"github.com/rivo/tview"
)

func modalNetStateJSONPage(ns *tuiNet.NetState, pages *tview.Pages) (*tview.Modal, error) {
	if pages == nil {
		return nil, fmt.Errorf("Can't add modal NetState page to nil pages")
	}

	modal := tview.NewModal().
		SetText(fmt.Sprintf("%+v", *ns)).
		SetTextColor(tcell.ColorBlack).
		SetBackgroundColor(newt.ColorGrey)
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Key() == tcell.KeyESC {
			pages.HidePage("netstate")
		}
		return event
	})

	return modal, nil
}

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

			var filteredNetState tuiNet.NetState
			if err := json.Unmarshal([]byte(state), &filteredNetState); err != nil {
				panic(err)
			}

			//netStatePage, err := modalNetStateJSONPage(&filteredNetState, pages)
			netStatePage, err := tuiNet.TreeView(filteredNetState, pages)
			pages.AddPage("netstate", netStatePage, true, true)
		}
	}
}

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()

	background := tview.NewBox().
		SetBorder(false).
		SetBackgroundColor(newt.ColorBlue)

	shouldConfigure := tview.NewModal().
		SetText("Do you wish for this node to be the one that runs the installation service (Only one node may perform this function)?").
		SetTextColor(tcell.ColorBlack).
		AddButtons([]string{"Quit", "Configure"}).
		SetDoneFunc(doneView(app, pages)).
		SetBackgroundColor(newt.ColorGrey)

	shouldConfigure.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Key() == tcell.KeyESC {
			app.Stop()
		}
		return event
	})

	pages.AddPage("background", background, true, true).
		AddPage("ShouldConfigure", shouldConfigure, true, true)

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}
