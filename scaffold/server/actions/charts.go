package actions

import (
	"errors"
	"net/http"
	"time"

	"github.com/nerdynz/helpers"

	flow "github.com/nerdynz/flow"
)

type anal struct {
	UniqueID     float64   `db:"unique_id" json:"UniqueID"`
	Date         string    `db:"date" json:"DateCreated"`
	DateCreated  time.Time `db:"date_created" json:"DateCreated"`
	DateModified time.Time `db:"date_modified" json:"DateModified"`
	Browser      string    `db:"browser" json:"Browser"`
	Device       string    `db:"device" json:"Device"`
	Version      int       `db:"version" json:"Version"`
}

func Views(ctx *flow.Context) {
	dayMonthYear := "01-" + helpers.PadLeft(ctx.URLParam("month"), 2, "0") + "-" + ctx.URLParam("year")
	start, err := time.Parse("02-01-2006", dayMonthYear)
	end := start.AddDate(0, 1, 0)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "bad month year combo", errors.New("bad month year combo"))
		return
	}
	as := make([]*anal, 0)
	err = ctx.DB.SQL(`select * from (
		select unique_id, device, browser, to_char(date_created, 'DD Mon YYYY') as date, max(date_created) as date_created from analytics
		where date_created between $1 and $2
		group by unique_id, device, browser, to_char(date_created, 'DD Mon YYYY')) a
	order by date_created
	`, start, end).QueryStructs(&as)

	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "failed to load analytics data", err)
		return
	}

	labels := make([]string, 0)
	all := make([]int, 0)
	mobile := make([]int, 0)
	desktop := make([]int, 0)

	for d := start; d.Month() == start.Month(); d = d.AddDate(0, 0, 1) {
		allCount := 0
		mobileCount := 0
		desktopCount := 0
		labels = append(labels, d.Format("02 Jan"))
		for _, a := range as {
			if d.Day() == a.DateCreated.Day() {
				if a.Device == "Phone" {
					mobileCount++
				} else if a.Device == "Computer" {
					desktopCount++
				}
				allCount++
			}
		}
		all = append(all, allCount)
		desktop = append(desktop, desktopCount)
		mobile = append(mobile, mobileCount)
	}

	series := make([][]int, 0)
	series = append(series, all)
	series = append(series, desktop)
	series = append(series, mobile)

	data := struct {
		Labels []string `json:"labels"`
		Series [][]int  `json:"series"`
	}{
		Labels: labels,
		Series: series,
	}
	ctx.JSON(http.StatusOK, data)
}
