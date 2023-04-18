const xArray = ["Alsace", "Ampère", "Argenteuil", "Billancourt", "Boule", "Courbevoie", "Danton", "La Briche", "Levallois", "Menus", "Nanterre", "Novion", "Puteaux", "Rueil", "St Ouen", "Tilliers"];

const yArray = [0, 5, 0, 0, 0, 1, 6, 10, 0, 0, 0, 0, 0, 0, 0, 0];

const data = [{
  x:nomSources,
  y:nbCuratifs,
  type:"bar"
}];

var options = {
  scrollZoom: true, // lets us scroll to zoom in and out - works
  showLink: false, // removes the link to edit on plotly - works
  modeBarButtonsToRemove: ['toImage', 'zoom2d', 'pan', 'pan2d', 'autoScale2d'],
  //modeBarButtonsToAdd: ['lasso2d'],
  displayModeBar: true, //this one does work
};

var config = {
  modeBarButtonsToRemove: ['sendDataToCloud', 'autoScale2d', 'hoverClosestCartesian', 'hoverCompareCartesian', 'lasso2d', 'select2d', 'toggleSpikelines', 'pan2d', 'zoomIn2d', 'zoomOut2d', 'zoom2d', 'drawclosedpath'],
  displaylogo: false,
  showTips: true,
  responsive: true
};


const layout = {
  title:"Nombre de Curatifs en cours",
  width: 1000,
  height: 400,
  autosize: true
};

Plotly.newPlot("myPlot", data, layout, config);
