async function FetchJSON() {
  const response = await fetch("https://localhost:8080/jsonGraph");

  if (!response.ok) {
    const msg = `An error has occured: ${response.status}`;
    throw new Error(msg);
  }

  const data = await response.json();
  return data;
}

FetchJSON().then(
  function (js) {
    let sources = [];
    let codeGMAO = [];
    let aRealiser = [];
    let enCours = [];
    let done = [];
    let total = [];

    for (let i = 0; i < js.length; i++) {
      aRealiser.push(js[i].a_realiser);
      enCours.push(js[i].en_cours);
      done.push(js[i].done);
      sources.push(js[i].name);
      codeGMAO.push(js[i].code_GMAO);
      total.push(js[i].curatifs);
    }

    bb.generate({
      bindto: "#myPlot",

      data: {
        names: {
          data1: "Curatifs",
          data2: "Sources",
        },
        columns: [
          ["Curatifs à réaliser", ...aRealiser],
          ["Curatifs en cours", ...enCours],
          ["Curatifs réalisés", ...done],
          ["Total", ...total],
        ],
        type: "bar",

        groups: [
          ["Curatifs à réaliser", "Curatifs en cours", "Curatifs réalisés"],
        ],
      },
      color: {
        pattern: ["#cc1111", "#0080ff", "#99cc33", "#999999"],
      },

      axis: {
        x: {
          type: "category",
          categories: [...codeGMAO],
          height: 50,
        },
      },

      size: {
        width: 1000,
        height: 400,
      },

      padding: true,

      zoom: {
        enabled: true,
        type: "drag",
      },

      legend: {
        position: "bottom",
      },

      bar: {
        width: {
          ratio: 0.5,
        },
      },
    });
  },
  function () {
    let plot = document.getElementById("myPlot");
    plot.textContent = "DB not found";
  }
);
