// let navLinks = document.querySelectorAll("nav a")

// for (let i = 0; i < navLinks.length; i++) {
//         let link = navLinks[i]

//         if (link.getAttribute("href") == window.location.pathname) {
//                 link.classList.add("live")
//                 break;
//         }
// }

function searchStatus() {
        let input, filter, table, tr, td, i, txtValue;

        input = document.getElementById("searchStatus"); // Fetch for the <input> searcg tag
        filter = input.value.toUpperCase(); // It simplifies the research
        table = document.getElementById("myTable"); // Created id just for this...
        tr = table.getElementsByTagName("tr") // Fetch every <tr> inside table
        
        /* Inside "sourceView" <table id="myTable"> there's 2 <tr>
         * tr[1] contains 2 <td>. We're going to look inside td[1]
         * which represents "Status" column.
         * And coz td[0] doesn't work...
         */
        
        for (i = 0; i < tr.length; i++) {
                // If a column is added don't forget to change
                //                               here ⬇️
                td = tr[i].getElementsByTagName("td")[1];

                if (td) {
                        txtValue = td.textContent || td.innerText;

                        if (txtValue.toUpperCase().indexOf(filter) > -1) {
                                tr[i].style.display = "";
                        } else {
                                tr[i].style.display = "none";
                        }
                }
        }
}
