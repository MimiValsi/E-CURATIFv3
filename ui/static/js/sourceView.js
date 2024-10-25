// Confirme si l'utilisateur veut vraiment supprimer une source
let del = document.getElementById("check-delete");
if (del) {
  del.onclick = () => {
    return confirm("Êtes-vous sûr(e) de supprimer?");
  };
}

// Lorsque l'utilisateur cherche dans la page .../source/view/xxx
// Ceci permet de chercher dans toutes les colonnes du tableau et changer les rangées selon la recherche.
let input = document.getElementById("search_info");

input.addEventListener("input", function () {
  let table = document.getElementById("info_table");
  let noMatchMessage = document.getElementById("noMatch");
  let rows = table.getElementsByTagName("tr");
  let filter = input.value.toLowerCase();
  let matchFound = false;

  for (let i = 1; i < rows.length; i++) {
    let row = rows[i];
    let cells = row.getElementsByTagName("td");
    let found = false;

    for (let j = 0; j < cells.length; j++) {
      let cell = cells[j];
      if (cell.textContent.toLowerCase().indexOf(filter) > -1) {
        found = true;
        matchFound = true;
        break;
      }
    }

    if (found) {
      row.style.display = "";
    } else {
      row.style.display = "none";
    }
  }

  if (!matchFound) {
    noMatchMessage.style.display = "block";
  } else {
    noMatchMessage.style.display = "none";
  }
});
