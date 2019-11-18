package main

import (
    "fmt"
    "net/http"
    "time"

    . "github.com/ajruckman/xlib"
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/plotutil"
    "gonum.org/v1/plot/vg"

    "github.com/ajruckman/ContraCore/internal/db"
    "github.com/ajruckman/ContraCore/internal/schema"
)

func main() {
    http.HandleFunc("/", series)
    fmt.Println("Listening on http://localhost:8080/")
    err := http.ListenAndServe(":8080", nil)
    Err(err)
}

func series(w http.ResponseWriter, r *http.Request) {
    var res []schema.QuestionCountsPerHour
    err := db.XDB.Select(&res, `

SELECT ts_round(time, 3600)
    AS hour,
    count(l.id)
FROM log l
GROUP BY 1
ORDER BY 1;

`)
    Err(err)

    var hours []time.Time
    var counts []float64

    points := make(plotter.XYs, len(res))

    for i, v := range res {
        hours = append(hours, v.Hour)
        counts = append(counts, float64(v.Count))

        points[i].X = float64(i)
        points[i].Y = float64(v.Count)
    }

    p, err := plot.New()
    Err(err)

    p.Title.Text = "DNS query counts per 10 minutes"

    err = plotutil.AddLinePoints(p, "count", points)
    Err(err)

    //f, err := chart.GetDefaultFont()
    //Err(err)

    //series := chart.TimeSeries{
    //    Name:    "===",
    //    XValues: hours,
    //    YValues: counts,
    //}

    c, err := p.WriterTo(16*vg.Inch, 4*vg.Inch, "png")
    Err(err)

    _, err = c.WriteTo(w)
    Err(err)

    //sma := &chart.SMASeries{
    //    InnerSeries: series,
    //    Period:      30,
    //}
    //
    //graph := chart.Chart{
    //    XAxis: chart.XAxis{
    //        Name: "The XAxis",
    //    },
    //    YAxis: chart.YAxis{
    //        Name: "The YAxis",
    //    },
    //
    //    Series: []chart.Series{
    //        series,
    //        sma,
    //    },
    //}
    //
    //fmt.Println(graph.GetFont())
    //
    //f, _ := os.Create("output.png")
    //defer f.Close()
    //graph.Render(chart.PNG, f)
    //
    //w.Header().Set("Content-Type", "image/png")
    //w.Header().Set("Cache-Control", "no-cache")
    //err = graph.Render(chart.PNG, w)
    //Err(err)
}
