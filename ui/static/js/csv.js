// (async () => {
//   const getJsonToCsv = async () => {
//     const response = await fetch("http://localhost:3001/pageTest");
//
//     const data = await response.json();
//     const js = JSON.stringify(data);
//     const jp = await JSON.parse(js);
//
//     return jp;
//   };
//
//   var array =
//     typeof getJsonToCsv != "object" ? JSON.parse(getJsonToCsv) : getJsonToCsv;
//   var str = "";
//
//   for (var i = 0; i < array.length; i++) {
//     var line = "";
//     for (var index in array[i]) {
//       if (line != "") line += ";";
//       line += array[i][index];
//     }
//
//     str += line + "\r\n";
//   }
//
//   console.log(str);
// })();
