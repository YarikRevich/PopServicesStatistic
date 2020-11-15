function updateTable() {
    request = new XMLHttpRequest()
    request.open("POST", "/usability")
    request.setRequestHeader("Content-Type", "application/json");
    request.addEventListener("readystatechange", () => {
        if (request.readyState == 4) {
            console.log(request.responseText)
            var table = document.getElementById("main-table");
            deleteRows(table);
            getDataFromJsonAndSaveIt(table, request.responseText);
        }
    })
    request.send()
}

function saveDataToTable(table, firstColumn, secondColumn) {
    var row = table.insertRow();
    row.insertCell(0).innerHTML = firstColumn;
    row.insertCell(1).innerHTML = secondColumn;
}

function getDataFromJsonAndSaveIt(table, notParsedJson) {
    var json = JSON.parse(notParsedJson);
    console.log(json);
    for (let key in json) {
        saveDataToTable(table, key, json[key]);
    }
}

function deleteRows(table) {
    table.getElementsByTagName("tbody")[0].innerHTML = table.rows[0].innerHTML;
}

window.onload = () => {
    this.updateTable();
}