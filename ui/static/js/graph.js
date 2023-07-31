(async () => {
        const getJson = async () => {
                const response = await fetch("http://localhost:3001/jsonGraph");
                
                const data = await response.json();
                const js = JSON.stringify(data);
                const jp = await JSON.parse(js)

                return jp;
        };
        
        const getJsonCD = async () => {
                const curatifDone = await fetch("http://localhost:3001/curatifDone");
                const dataCD = await curatifDone.json();
                const jsCD = JSON.stringify(dataCD);
                const jpCD = await JSON.parse(jsCD);

                return jpCD;
        }
        // const prioJson = async () => {
        //   const response = await fetch("http://localhost:3001/prioData");
  //
  //   const data = await response.json();
  //   const js = JSON.stringify(data);
  //   const jp = await JSON.parse(js);
  //
  //   return jp;
  // };

  // Fetch every active 'Curatif'
  const jsData = await getJson();

  let nbCuratifs = [];
  for (let i = 0; i < jsData.length; i++) {
    nbCuratifs.push(jsData[i].curatifs);
  }

  let nomSources = [];
  for (let i = 0; i < jsData.length; i++) {
    nomSources.push(jsData[i].name);
  }

  let codeGMAO = []
  for (let i = 0; i < jsData.length; i++) {
    codeGMAO.push(jsData[i].code_GMAO)
  }

  // Fetch only already done 'Curatif'
  const jsDataCD = await getJsonCD();
  let nbCuratifsDone = [];
        for (let i = 0; i < jsDataCD.length; i++) {
                nbCuratifsDone.push(jsDataCD[i].curatifs);
        }

  var chart = bb.generate({
    bindto: "#myPlot",

    data: {
      names: {
        data1: "Curatifs",
        data2: "Sources"
      },
      columns: [
        ["Curatifs en cours", ...nbCuratifs],
        ["Curatifs résolus", ...nbCuratifsDone],
      ],
      type: "bar",
    },

    axis: {
      x: {
        type: "category",
        categories: [...codeGMAO],
        height: 50,
        tick: {
          rotate: 75,
          multiline: false,
        }
      }
    },

    size: {
      width: 1000,
      height: 400
    },

    padding: true,

    resize: true,

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
    }
  });
})();
