
const fill = (url, target, callback) => {
    fetch(url)
        .then(function (response) {
            return response.text();
        })
        .then(function (contents) {
            morphdom(document.getElementById(target), `<div id="${target}">${contents}</div>`);
            if (callback) {
                callback();
            }
        });
}

const saveEntity = (entity, id) => {
    const data = document.getElementById("entity-json").innerHTML;
    postData(`/api/${entity}/${id}`, data)
        .then(response => saveResult(response, entity));
}

const deleteEntity = (entity, id) => {
    if (confirm(`Please confirm that you want to delete ${entity} ${id}`)) {
        deleteData(`/api/${entity}/${id}`)
            .then(response => deleteResult(response, entity));
    }
}

// Ajax methods

const postData = (url = ``, data = "") => {
    return fetch(url, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json; charset=utf-8",

        },
        body: data,
    })
};

const deleteData = (url = ``) => {
    return fetch(url, {
        method: "DELETE",
    })
};

// Save Result

// saveResult executes the necessary actions after an entity save
const saveResult = (response, entity) => {
    // For an OK response
    if (response.ok) {
        // Decode the JSON
        response.json().then(
            saved => {
                // Using the returned ID,
                // get the new/edited record
                // and refresh the detail view
                fill(`/${entity}/${saved.ID}`, 'detail-view', alertSaveSuccess)

                //And finally refresh the list view
                fill(`/${entity}`, 'list-view')
            }
        )
    } else {
        response.json().then(
            error => alertSaveFail(error.errorMessage)
        )
    }
}

const deleteResult = (response, entity) => {
    // For an OK response
    if (response.ok) {
        // Decode the JSON
        response.json().then(
            _ => {
                // refresh the list view
                fill(`/${entity}`, 'list-view', alertDeleteSuccess)
            }
        )
    } else {
        response.json().then(
            error => alertDeleteFail(error.errorMessage)
        )
    }
}

// Alerts

const alertSaveSuccess = () => {
    document.getElementById("save-result").classList.add("success");
    document.getElementById("save-result").innerHTML = `<i class="fas fa-check"></i> Saved`;
}

const alertSaveFail = errorMessage => {
    document.getElementById("save-result").classList.add("fail");
    document.getElementById("save-result").innerHTML = `<i class="fas fa-exclamation-triangle"></i> Failed to Save! `;

    document.getElementById("save-error-message").innerHTML = errorMessage;
}

const alertDeleteSuccess = () => {
    console.log("Delete Success");
}

const alertDeleteFail = errorMessage => {
    console.log("Delete Fail");
}

// DOM

const toggleQueryForm = () => {
    toggleDiv("query-form")
}

const toggleSaveNew = () => {
    toggleDiv("detail")
}

const toggleStartEnd = () => {
    var newVal = document.getElementById("match").value;
    if (newVal) {
        document.getElementById("start-end").style.display = "none";
    } else {
        document.getElementById("start-end").style.display = "block";
    }

}


function toggleDiv(id) {
    var div = document.getElementById(id);
    div.style.display = div.style.display == "none" ? "block" : "none";
}