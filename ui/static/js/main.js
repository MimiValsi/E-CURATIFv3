// It's used in info and source view page.
// To avoid missclick, ask to confirm the delete
let del = document.getElementById("check-delete");
if (del) {
  del.onclick = function () {
    return confirm("Êtes-vous sûr(e) de supprimer?");
  };
}

let getFile = document.getElementById("getFile");
let imprt = document.getElementById("importCSV");
imprt.onchange = () => {
  getFile.submit();
  alert("Fichier en cours de transfert...");
};

function search_status() {
  document
    .getElementById("search_status")
    .addEventListener("keyup", function () {
      let input, filter, table, tr, td, i, txtValue;

      // Fetch for the <input> searcg tag
      input = document.getElementById("search_status");
      filter = input.value.toUpperCase(); // It simplifies the research
      table = document.getElementById("myTable");
      tr = table.getElementsByTagName("tr"); // Fetch every <tr> inside table

      /* Inside "sourceView" <table id="myTable"> there's 2 <tr>
       * tr[1] contains 2 <td>. We're going to look inside td[1]
       * which represents "Status" column.
       * And coz td[0] doesn't work... Added to TODO list
       */

      for (i = 1; i < tr.length; i++) {
        // For Info search change to ...("td")[0];
        // For Priority search change to ...("td")[1];
        // For Status search change to ...("td")[2];
        td = tr[i].getElementsByTagName("td")[0];
        //                               here ^

        if (td) {
          txtValue = td.textContent;
          console.log(txtValue);

          if (txtValue.toUpperCase().indexOf(filter) > -1) {
            tr[i].style.display = "";
          } else {
            tr[i].style.display = "none";
          }
        }
      }
    });
}
