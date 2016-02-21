package plots;

import (
  "github.com/gonum/plot";
  "github.com/gonum/plot/plotter";
  "github.com/gonum/plot/plotutil";
  "github.com/gonum/plot/vg";
);

func GeneratePlot(x, y []float64, title, xLabel, yLabel, legendLabel, fileName string) {
  outPlotPoints := make(plotter.XYs, len(x));

  for i, _ := range x {
    outPlotPoints[i].X = x[i];
    outPlotPoints[i].Y = y[i];
  }

  outPlot, err := plot.New();
  if err != nil {
    panic(err);
  }

  outPlot.Title.Text = title;
  outPlot.X.Label.Text = xLabel;
  outPlot.Y.Label.Text = yLabel;

  err = plotutil.AddLines(outPlot,
    legendLabel, outPlotPoints);
  if err != nil {
    panic(err);
  }

  if err := outPlot.Save(6 * vg.Inch, 6 * vg.Inch, fileName); err != nil {
    panic(err);
  }
};
