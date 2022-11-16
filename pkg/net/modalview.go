package net

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/openshift-agent-team/tui/pkg/newt"
	"github.com/rivo/tview"
)

func modalNetStateJSONPage(ns *NetState, pages *tview.Pages) (*tview.Modal, error) {
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
