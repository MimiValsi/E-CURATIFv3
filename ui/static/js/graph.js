// import bb, {bar} from "billboard.js";

let nbCuratifs = [];
for (let i = 0; i < incStruct.length; i++) {
  nbCuratifs.push(incStruct[i].curatifs);
}

let nomSources = [];
for (let i = 0; i < incStruct.length; i++) {
  nomSources.push(incStruct[i].name);
}
// for ESM environment, need to import modules as:
// import bb, {bar} from "billboard.js";

var chart = bb.generate({
  data: {
    names: {
      data1: "Curatifs",
      data2: "Sources"
    },
    columns: [
      ["Nombre de Curatifs en cours", ...nbCuratifs],
    ],
    type: "bar", // for ESM specify as: bar()
  },

  axis: {
    x: {
      type: "category",
      categories: [...nomSources],
      tick: {
        rotate: 75,
        multiline: false
      }
    }
  },

  size: {
    width: 1000
  },

  padding: true,

  zoom: {
    enabled: true,
    type: "drag"
  },

  legend: {
    position: "inset"
  },

  bar: {
    width: {
      ratio: 0.5
    }
    // radius: {
    //   ratio: 0.5,
    // }
  },
  bindto: "#myPlot"
});
